package router

import (
	"ledger/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes(h *handler.LedgerHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", h.HealthCheck)

	r.Route("/v1", func(r chi.Router) {

		r.Post("/accounts", h.CreateAccount)
		r.Get("/accounts/{accountID}/balance", h.GetBalance)

		r.Post("/transactions", h.CreateTransaction)
		r.Get("/transactions", h.ListTransactions)
		r.Post("/transactions/{transactionID}/reverse", h.ReverseTransaction)
	})

	return r
}
