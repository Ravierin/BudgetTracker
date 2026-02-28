package main

import (
	"github.com/Ravierin/BudgetTracker/backend/internal/api"
	"github.com/Ravierin/BudgetTracker/backend/internal/repository"
	"github.com/Ravierin/BudgetTracker/backend/internal/service"
	"github.com/Ravierin/BudgetTracker/backend/pkg/config"
	"github.com/Ravierin/BudgetTracker/backend/pkg/database"
	"github.com/Ravierin/BudgetTracker/backend/pkg/server"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewDatabase(cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	positionRepo := repository.NewPositionRepository(db)
	positionService := service.NewPositionService(positionRepo)
	withdrawalRepo := repository.NewWithdrawalRepository(db)
	withdrawalService := service.NewWithdrawalService(withdrawalRepo)
	incomeRepo := repository.NewMonthlyIncomeRepository(db)
	incomeService := service.NewMonthlyIncomeService(incomeRepo)
	apiKeyRepo := repository.NewAPIKeyRepository(db)
	apiKeyService := service.NewAPIKeyService(apiKeyRepo)
	balanceService := service.NewBalanceService(apiKeyService)

	// Create clients with empty keys - will be populated dynamically from DB
	bybitClient := api.NewBybitClient("", "")
	mexcClient := api.NewMEXClient("", "")

	srv := server.NewServer(positionService, withdrawalService, incomeService, apiKeyService, balanceService, positionRepo, bybitClient, mexcClient)

	// Create sync services for all exchanges
	exchanges := []string{"bybit", "mexc"}
	for _, exchangeName := range exchanges {
		syncService := server.NewSyncService(positionService, apiKeyService, srv.GetWSHub(), 30*time.Second, exchangeName)
		go syncService.Start()
	}

	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      srv.GetHandler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}()

	log.Println("Starting server on port 8080")
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Println("Server stopped")
}
