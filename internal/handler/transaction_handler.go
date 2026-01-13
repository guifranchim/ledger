package handler

import (
	"ledger/internal/utils"
	"net/http"
)

type CreateTransactionRequest struct {
	FromAccountID string  `json:"from_account_id" validate:"required"`
	ToAccountID   string  `json:"to_account_id" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	Description   string  `json:"description" validate:"max=255"`
}

func (h *LedgerHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	data := &CreateTransactionRequest{}

	if utils.DecodeAndValidate(w, r, h.Validate, data) {
		return
	}

	utils.SuccessResponse(w, r, http.StatusCreated, map[string]interface{}{
		"message":         "transaction created",
		"from_account_id": data.FromAccountID,
		"to_account_id":   data.ToAccountID,
		"amount":          data.Amount,
	})
}

func (h *LedgerHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"transactions": []}`))
}

func (h *LedgerHandler) ReverseTransaction(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "transaction reversed"}`))
}
