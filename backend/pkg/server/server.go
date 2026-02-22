package server

import (
	"BudgetTracker/backend/internal/api"
	"BudgetTracker/backend/internal/handler"
	"BudgetTracker/backend/internal/service"
	"BudgetTracker/backend/pkg/websocket"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	router            *mux.Router
	positionService   *service.PositionService
	withdrawalService *service.WithdrawalService
	incomeService     *service.MonthlyIncomeService
	bybitClient       *api.BybitClient
	wsHub             *websocket.Hub
}

func NewServer(
	positionService *service.PositionService,
	withdrawalService *service.WithdrawalService,
	incomeService *service.MonthlyIncomeService,
	bybitClient *api.BybitClient,
) *Server {
	hub := websocket.NewHub()
	go hub.Run()

	s := &Server{
		router:            mux.NewRouter(),
		positionService:   positionService,
		withdrawalService: withdrawalService,
		incomeService:     incomeService,
		bybitClient:       bybitClient,
		wsHub:             hub,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.Use(loggingMiddleware)
	s.router.Use(corsMiddleware)

	api := s.router.PathPrefix("/api/v1").Subrouter()

	positionHandler := handler.NewPositionHandler(s.positionService, s.bybitClient, s.wsHub)
	api.HandleFunc("/positions", positionHandler.GetAllPositions).Methods("GET")
	api.HandleFunc("/positions/{id}", positionHandler.GetPosition).Methods("GET")
	api.HandleFunc("/positions", positionHandler.CreatePosition).Methods("POST")
	api.HandleFunc("/positions/{id}", positionHandler.DeletePosition).Methods("DELETE")
	api.HandleFunc("/positions/sync", positionHandler.SyncPositions).Methods("POST")

	withdrawalHandler := handler.NewWithdrawalHandler(s.withdrawalService, s.wsHub)
	api.HandleFunc("/withdrawals", withdrawalHandler.GetAllWithdrawals).Methods("GET")
	api.HandleFunc("/withdrawals", withdrawalHandler.CreateWithdrawal).Methods("POST")
	api.HandleFunc("/withdrawals/{id}", withdrawalHandler.DeleteWithdrawal).Methods("DELETE")

	incomeHandler := handler.NewMonthlyIncomeHandler(s.incomeService, s.wsHub)
	api.HandleFunc("/monthly-income", incomeHandler.GetAllMonthlyIncomes).Methods("GET")
	api.HandleFunc("/monthly-income", incomeHandler.CreateMonthlyIncome).Methods("POST")
	api.HandleFunc("/monthly-income/{id}", incomeHandler.DeleteMonthlyIncome).Methods("DELETE")

	api.HandleFunc("/ws", s.wsHub.HandleWebSocket)

	s.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
}

func (s *Server) Start(port string) error {
	log.Printf("Starting server on port %s", port)
	return http.ListenAndServe(":"+port, s.router)
}

func (s *Server) GetWSHub() *websocket.Hub {
	return s.wsHub
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed in %v", time.Since(start))
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type SyncService struct {
	positionService *service.PositionService
	bybitClient     *api.BybitClient
	wsHub           *websocket.Hub
	interval        time.Duration
	stopChan        chan struct{}
}

func NewSyncService(
	positionService *service.PositionService,
	bybitClient *api.BybitClient,
	wsHub *websocket.Hub,
	interval time.Duration,
) *SyncService {
	return &SyncService{
		positionService: positionService,
		bybitClient:     bybitClient,
		wsHub:           wsHub,
		interval:        interval,
		stopChan:        make(chan struct{}),
	}
}

func (s *SyncService) Start() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.sync()
		case <-s.stopChan:
			return
		}
	}
}

func (s *SyncService) Stop() {
	close(s.stopChan)
}

func (s *SyncService) sync() {
	ctx := context.Background()

	positions, err := s.bybitClient.GetPositions()
	if err != nil {
		return
	}

	if err := s.positionService.SavePositionsBatch(ctx, positions); err != nil {
		return
	}

	message := map[string]interface{}{
		"type":      "positions_update",
		"positions": positions,
		"count":     len(positions),
	}
	s.wsHub.Broadcast(message)
}
