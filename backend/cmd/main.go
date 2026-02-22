package main

import (
	"BudgetTracker/backend/internal/api"
	"BudgetTracker/backend/internal/repository"
	"BudgetTracker/backend/internal/service"
	"BudgetTracker/backend/pkg/config"
	"BudgetTracker/backend/pkg/database"
	"BudgetTracker/backend/pkg/server"
	"log"
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

	srv := server.NewServer(positionService, withdrawalService, incomeService, bybitClient)

	syncService := server.NewSyncService(positionService, bybitClient, srv.GetWSHub(), 30*time.Second)
	go syncService.Start()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down...")
		syncService.Stop()
	}()

	// Запуск сервера
	if err := srv.Start("8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
