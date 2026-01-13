package main

import (
	"context"
	"ledger/internal/config"
	"ledger/internal/handler"
	"ledger/internal/models"
	"ledger/internal/repository"
	"ledger/internal/router"
	"ledger/internal/services"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	port := ":8080"
	slog.Info("Starting API server", "port", port)

	dbConfig := config.GetDatabaseConfig()
	db, err := config.ConnectDatabase(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&models.Account{}, &models.Transaction{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	slog.Info("Database migrations completed")

	ledgerRepo := repository.NewLedgerRepository(db)
	ledgerService := services.NewLedgerService(ledgerRepo, db)
	ledgerHandler := handler.NewLedgerHandler(ledgerService)

	r := router.SetupRoutes(ledgerHandler)

	// Configurar servidor HTTP
	server := &http.Server{
		Addr:    port,
		Handler: r,
	}

	// Canal para capturar sinais de shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Iniciar servidor em goroutine
	go func() {
		slog.Info("API server is running", "port", port, "workers", 10)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Aguardar sinal de shutdown
	<-shutdown
	slog.Info("Shutting down server gracefully...")

	// Shutdown dos workers
	ledgerService.Shutdown()
	slog.Info("Workers stopped")

	// Shutdown do servidor HTTP
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server stopped")
}
