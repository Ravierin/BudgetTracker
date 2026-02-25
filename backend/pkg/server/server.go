package server

import (
	"BudgetTracker/backend/internal/api"
	"BudgetTracker/backend/internal/handler"
	"BudgetTracker/backend/internal/model"
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
	apiKeyService     *service.APIKeyService
	bybitClient       *api.BybitClient
	mexcClient        *api.MEXClient
	wsHub             *websocket.Hub
}

func NewServer(
	positionService *service.PositionService,
	withdrawalService *service.WithdrawalService,
	incomeService *service.MonthlyIncomeService,
	apiKeyService *service.APIKeyService,
	bybitClient *api.BybitClient,
	mexcClient *api.MEXClient,
	gateClient *api.GateClient,
	bitgetClient *api.BitgetClient,
) *Server {
	hub := websocket.NewHub()
	go hub.Run()

	s := &Server{
		router:            mux.NewRouter(),
		positionService:   positionService,
		withdrawalService: withdrawalService,
		incomeService:     incomeService,
		apiKeyService:     apiKeyService,
		bybitClient:       bybitClient,
		mexcClient:        mexcClient,
		wsHub:             hub,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.Use(loggingMiddleware)
	s.router.Use(corsMiddleware)

	// Handle CORS preflight for all routes
	s.router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	api := s.router.PathPrefix("/api/v1").Subrouter()

	positionHandler := handler.NewPositionHandler(s.positionService, s.wsHub)
	api.HandleFunc("/positions", positionHandler.GetAllPositions).Methods("GET")
	api.HandleFunc("/positions/{id}", positionHandler.GetPosition).Methods("GET")
	api.HandleFunc("/positions", positionHandler.CreatePosition).Methods("POST")
	api.HandleFunc("/positions/{id}", positionHandler.DeletePosition).Methods("DELETE")
	api.HandleFunc("/positions/sync", positionHandler.SyncPositions).Methods("POST")

	withdrawalHandler := handler.NewWithdrawalHandler(s.withdrawalService, s.wsHub)
	api.HandleFunc("/withdrawals", withdrawalHandler.GetAllWithdrawals).Methods("GET")
	api.HandleFunc("/withdrawals", withdrawalHandler.CreateWithdrawal).Methods("POST")
	api.HandleFunc("/withdrawals/{id}", withdrawalHandler.DeleteWithdrawal).Methods("DELETE")

	incomeHandler := handler.NewMonthlyIncomeHandler(s.positionService, s.wsHub)
	api.HandleFunc("/monthly-income", incomeHandler.GetAllMonthlyIncomes).Methods("GET")
	api.HandleFunc("/monthly-income/{id}", incomeHandler.GetMonthlyIncome).Methods("GET")

	apiKeyHandler := handler.NewAPIKeyHandler(s.apiKeyService)
	api.HandleFunc("/api-keys", apiKeyHandler.GetAPIKeys).Methods("GET")
	api.HandleFunc("/api-keys", apiKeyHandler.SaveAPIKeys).Methods("POST")

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

func (s *Server) GetHandler() http.Handler {
	return s.router
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
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type SyncService struct {
	positionService *service.PositionService
	apiKeyService   *service.APIKeyService
	wsHub           *websocket.Hub
	interval        time.Duration
	stopChan        chan struct{}
	exchangeName    string
}

func NewSyncService(
	positionService *service.PositionService,
	apiKeyService *service.APIKeyService,
	wsHub *websocket.Hub,
	interval time.Duration,
	exchangeName string,
) *SyncService {
	return &SyncService{
		positionService: positionService,
		apiKeyService:   apiKeyService,
		wsHub:           wsHub,
		interval:        interval,
		stopChan:        make(chan struct{}),
		exchangeName:    exchangeName,
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
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	apiKey, err := s.apiKeyService.GetAPIKey(ctx, s.exchangeName)
	if err != nil {
		log.Printf("[%s] Failed to get API key: %v", s.exchangeName, err)
		return
	}

	if apiKey.APIKey == "" || apiKey.APISecret == "" {
		log.Printf("[%s] API keys not configured", s.exchangeName)
		return
	}

	log.Printf("[%s] Syncing positions...", s.exchangeName)

	var positions []model.Position
	switch s.exchangeName {
	case "bybit":
		client := api.NewBybitClient(apiKey.APIKey, apiKey.APISecret)
		positions, err = client.GetPositionsWithContext(ctx)
	case "mexc":
		client := api.NewMEXClient(apiKey.APIKey, apiKey.APISecret)
		positions, err = client.GetPositionsWithContext(ctx)
	case "gate":
		client := api.NewGateClient(apiKey.APIKey, apiKey.APISecret)
		positions, err = client.GetPositionsWithContext(ctx)
	case "bitget":
		client := api.NewBitgetClient(apiKey.APIKey, apiKey.APISecret)
		positions, err = client.GetPositionsWithContext(ctx)
	default:
		log.Printf("Unknown exchange: %s", s.exchangeName)
		return
	}

	if err != nil {
		log.Printf("[%s] Sync error: %v", s.exchangeName, err)
		return
	}

	log.Printf("[%s] Retrieved %d positions", s.exchangeName, len(positions))

	if err := s.positionService.SavePositionsBatch(ctx, positions); err != nil {
		log.Printf("[%s] Failed to save positions: %v", s.exchangeName, err)
		return
	}

	log.Printf("[%s] Saved %d positions successfully", s.exchangeName, len(positions))

	message := map[string]interface{}{
		"type":      "positions_update",
		"positions": positions,
		"count":     len(positions),
		"exchange":  s.exchangeName,
	}
	s.wsHub.Broadcast(message)
}
