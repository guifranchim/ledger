package handler

import (
	"ledger/internal/utils"
	"net/http"
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

	utils.SuccessResponse(w, r, http.StatusCreated, map[string]interface{}{
		"message":         "account created",
		"owner_name":      data.OwnerName,
		"initial_balance": data.InitialBalance,
	})
}

func (h *LedgerHandler) GetBalance(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"balance": 0}`))
}
