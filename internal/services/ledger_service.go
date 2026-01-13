package services

import (
	"context"
	"errors"
	"ledger/internal/models"
	"ledger/internal/repository"
	"sync"

	"gorm.io/gorm"
)

const NumWorkers = 10

type TransactionJob struct {
	FromAccountID string
	ToAccountID   string
	Amount        float64
	Description   string
	ResultChan    chan error
}

type LedgerService struct {
	repo       *repository.LedgerRepository
	db         *gorm.DB
	mu         sync.Mutex
	jobQueue   chan *TransactionJob
	workerPool *sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewLedgerService(repo *repository.LedgerRepository, db *gorm.DB) *LedgerService {
	ctx, cancel := context.WithCancel(context.Background())

	service := &LedgerService{
		repo:       repo,
		db:         db,
		jobQueue:   make(chan *TransactionJob, 100),
		workerPool: &sync.WaitGroup{},
		ctx:        ctx,
		cancel:     cancel,
	}

	service.startWorkers()

	return service
}

func (s *LedgerService) startWorkers() {
	for i := 0; i < NumWorkers; i++ {
		s.workerPool.Add(1)
		go s.worker(i)
	}
}

func (s *LedgerService) worker(id int) {
	defer s.workerPool.Done()

	for {
		select {
		case <-s.ctx.Done():
			return
		case job, ok := <-s.jobQueue:
			if !ok {
				return
			}

			err := s.processTransaction(job.FromAccountID, job.ToAccountID, job.Amount, job.Description)
			job.ResultChan <- err
			close(job.ResultChan)
		}
	}
}

func (s *LedgerService) Shutdown() {
	s.cancel()
	close(s.jobQueue)
	s.workerPool.Wait()
}

func (s *LedgerService) CreateAccount(ownerName string, initialBalance float64) (*models.Account, error) {
	if initialBalance < 0 {
		return nil, errors.New("initial balance cannot be negative")
	}

	account, err := s.repo.CreateAccount(ownerName, initialBalance)
	if err != nil {
		return nil, err
	}

	if initialBalance > 0 {
		_, err = s.repo.CreateTransaction(
			account.ID,
			models.TransactionTypeCredit,
			initialBalance,
			"Initial balance",
		)
		if err != nil {
			return nil, err
		}
	}

	return account, nil
}

func (s *LedgerService) GetBalance(accountID string) (float64, error) {
	account, err := s.repo.GetAccountByID(accountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("account not found")
		}
		return 0, err
	}
	return account.Balance, nil
}

func (s *LedgerService) CreateTransaction(fromAccountID, toAccountID string, amount float64, description string) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	job := &TransactionJob{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        amount,
		Description:   description,
		ResultChan:    make(chan error, 1),
	}

	select {
	case s.jobQueue <- job:

		return <-job.ResultChan
	case <-s.ctx.Done():
		return errors.New("service is shutting down")
	}
}

func (s *LedgerService) processTransaction(fromAccountID, toAccountID string, amount float64, description string) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	return s.db.Transaction(func(tx *gorm.DB) error {
		var fromAccount, toAccount *models.Account
		var err error

		if fromAccountID < toAccountID {
			fromAccount, err = s.repo.GetAccountByIDForUpdate(tx, fromAccountID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("from account not found")
				}
				return err
			}

			toAccount, err = s.repo.GetAccountByIDForUpdate(tx, toAccountID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("to account not found")
				}
				return err
			}
		} else {
			toAccount, err = s.repo.GetAccountByIDForUpdate(tx, toAccountID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("to account not found")
				}
				return err
			}

			fromAccount, err = s.repo.GetAccountByIDForUpdate(tx, fromAccountID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("from account not found")
				}
				return err
			}
		}

		if fromAccount.Balance < amount {
			return errors.New("insufficient balance")
		}

		_, err = s.repo.CreateTransactionInTx(tx, fromAccountID, models.TransactionTypeDebit, amount, description)
		if err != nil {
			return err
		}

		if err := s.repo.UpdateAccountBalanceInTx(tx, fromAccountID, fromAccount.Balance-amount); err != nil {
			return err
		}

		_, err = s.repo.CreateTransactionInTx(tx, toAccountID, models.TransactionTypeCredit, amount, description)
		if err != nil {
			return err
		}

		if err := s.repo.UpdateAccountBalanceInTx(tx, toAccountID, toAccount.Balance+amount); err != nil {
			return err
		}

		return nil
	})
}

func (s *LedgerService) ListTransactions(accountID string, limit, offset int) ([]models.Transaction, error) {
	if accountID != "" {
		return s.repo.GetTransactionsByAccountID(accountID, limit, offset)
	}
	return s.repo.GetAllTransactions(limit, offset)
}

func (s *LedgerService) ReverseTransaction(transactionID string) error {
	tx, err := s.repo.GetTransactionByID(transactionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("transaction not found")
		}
		return err
	}

	account, err := s.repo.GetAccountByID(tx.AccountID)
	if err != nil {
		return err
	}

	var newBalance float64
	var reverseType models.TransactionType

	if tx.Type == models.TransactionTypeDebit {
		newBalance = account.Balance + tx.Amount
		reverseType = models.TransactionTypeCredit
	} else if tx.Type == models.TransactionTypeCredit {
		if account.Balance < tx.Amount {
			return errors.New("insufficient balance to reverse")
		}
		newBalance = account.Balance - tx.Amount
		reverseType = models.TransactionTypeDebit
	} else {
		return errors.New("cannot reverse a reversal transaction")
	}

	_, err = s.repo.CreateTransaction(
		tx.AccountID,
		reverseType,
		tx.Amount,
		"Reversal of transaction "+transactionID,
	)
	if err != nil {
		return err
	}

	return s.repo.UpdateAccountBalance(tx.AccountID, newBalance)
}
