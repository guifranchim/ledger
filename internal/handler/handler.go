package handler

import (
	"ledger/internal/services"

	"github.com/go-playground/validator/v10"
)

type LedgerHandler struct {
	Validate      *validator.Validate
	LedgerService *services.LedgerService
}

func NewLedgerHandler(LedgerService *services.LedgerService) *LedgerHandler {
	return &LedgerHandler{
		Validate:      validator.New(),
		LedgerService: LedgerService,
	}
}
