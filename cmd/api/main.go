package main

import (
	"ledger/internal/handler"
	"ledger/internal/router"
	"ledger/internal/services"
	"log/slog"
	"net/http"
)

func main() {
	port := ":8080"
	slog.Info("Starting API server", "port", port)

	// repositories

	// services
	ledgerService := services.NewLedgerService()

	// handlers
	ledgerHandler := handler.NewLedgerHandler(ledgerService)

	r := router.SetupRoutes(ledgerHandler)

	slog.Info("API server is running", "port", port)
	if err := http.ListenAndServe(port, r); err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}
