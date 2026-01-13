package repository

import (
	"ledger/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LedgerRepository struct {
	db *gorm.DB
}

func NewLedgerRepository(db *gorm.DB) *LedgerRepository {
	return &LedgerRepository{
		db: db,
	}
}

func (r *LedgerRepository) CreateAccount(ownerName string, initialBalance float64) (*models.Account, error) {
	account := &models.Account{
		OwnerName: ownerName,
		Balance:   initialBalance,
	}

	if err := r.db.Create(account).Error; err != nil {
		return nil, err
	}

	return account, nil
}

func (r *LedgerRepository) GetAccountByID(id string) (*models.Account, error) {
	var account models.Account
	if err := r.db.First(&account, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *LedgerRepository) GetAccountByIDForUpdate(tx *gorm.DB, id string) (*models.Account, error) {
	var account models.Account
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&account, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *LedgerRepository) UpdateAccountBalance(id string, newBalance float64) error {
	return r.db.Model(&models.Account{}).Where("id = ?", id).Update("balance", newBalance).Error
}

func (r *LedgerRepository) UpdateAccountBalanceInTx(tx *gorm.DB, id string, newBalance float64) error {
	return tx.Model(&models.Account{}).Where("id = ?", id).Update("balance", newBalance).Error
}

func (r *LedgerRepository) CreateTransaction(accountID string, txType models.TransactionType, amount float64, description string) (*models.Transaction, error) {
	tx := &models.Transaction{
		AccountID:   accountID,
		Type:        txType,
		Amount:      amount,
		Description: description,
	}

	if err := r.db.Create(tx).Error; err != nil {
		return nil, err
	}

	return tx, nil
}

func (r *LedgerRepository) CreateTransactionInTx(tx *gorm.DB, accountID string, txType models.TransactionType, amount float64, description string) (*models.Transaction, error) {
	transaction := &models.Transaction{
		AccountID:   accountID,
		Type:        txType,
		Amount:      amount,
		Description: description,
	}

	if err := tx.Create(transaction).Error; err != nil {
		return nil, err
	}

	return transaction, nil
}

func (r *LedgerRepository) GetTransactionsByAccountID(accountID string, limit, offset int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("account_id = ?", accountID).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error

	return transactions, err
}

func (r *LedgerRepository) GetAllTransactions(limit, offset int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error

	return transactions, err
}

func (r *LedgerRepository) GetTransactionByID(id string) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := r.db.First(&transaction, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}
