package handler

import (
	"ledger/internal/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CreateAccountRequest struct {
	OwnerName      string  `json:"owner_name" validate:"required,min=3,max=100"`
	InitialBalance float64 `json:"initial_balance" validate:"required,gte=0"`
}

func (h *LedgerHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	data := &CreateAccountRequest{}

	if utils.DecodeAndValidate(w, r, h.Validate, data) {
		return
	}

	account, err := h.LedgerService.CreateAccount(data.OwnerName, data.InitialBalance)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(w, r, http.StatusCreated, map[string]interface{}{
		"message":         "account created",
		"id":              account.ID,
		"owner_name":      account.OwnerName,
		"initial_balance": account.Balance,
		"created_at":      account.CreatedAt,
	})
}

func (h *LedgerHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "accountID")

	if accountID == "" {
		utils.ErrorResponse(w, r, http.StatusBadRequest, "account_id is required")
		return
	}

	balance, err := h.LedgerService.GetBalance(accountID)
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(w, r, http.StatusOK, map[string]interface{}{
		"account_id": accountID,
		"balance":    balance,
	})
}
