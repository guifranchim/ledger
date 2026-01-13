package handler

import (
	"ledger/internal/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

	err := h.LedgerService.CreateTransaction(
		data.FromAccountID,
		data.ToAccountID,
		data.Amount,
		data.Description,
	)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
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
	accountID := r.URL.Query().Get("account_id")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	transactions, err := h.LedgerService.ListTransactions(accountID, limit, offset)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, r, http.StatusOK, map[string]interface{}{
		"transactions": transactions,
		"limit":        limit,
		"offset":       offset,
	})
}

func (h *LedgerHandler) ReverseTransaction(w http.ResponseWriter, r *http.Request) {
	transactionID := chi.URLParam(r, "transactionID")

	if transactionID == "" {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "transaction_id is required")
		return
	}

	err := h.LedgerService.ReverseTransaction(transactionID)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(w, r, http.StatusOK, map[string]interface{}{
		"message":        "transaction reversed",
		"transaction_id": transactionID,
	})
}
