package main

import (
	"BudgetTracker/backend/internal/api"
	"BudgetTracker/backend/internal/repository"
	"BudgetTracker/backend/internal/service"
	"BudgetTracker/backend/pkg/config"
	"BudgetTracker/backend/pkg/database"
	"BudgetTracker/backend/pkg/server"
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
	bybitClient := api.NewBybitClient(cfg.BybitAPIKey, cfg.BybitAPISecretKey)
	mexcClient := api.NewMEXClient(cfg.MEXCAPIKey, cfg.MEXCAPISecretKey)

	srv := server.NewServer(positionService, withdrawalService, incomeService, bybitClient, mexcClient)

	bybitSyncService := server.NewSyncService(positionService, bybitClient, srv.GetWSHub(), 30*time.Second, "bybit")
	mexcSyncService := server.NewSyncService(positionService, mexcClient, srv.GetWSHub(), 30*time.Second, "mexc")

	go bybitSyncService.Start()
	go mexcSyncService.Start()

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
		bybitSyncService.Stop()
		mexcSyncService.Stop()

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
