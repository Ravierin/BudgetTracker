package server

import (
	"github.com/Ravierin/BudgetTracker/backend/internal/api"
	"github.com/Ravierin/BudgetTracker/backend/internal/handler"
	"github.com/Ravierin/BudgetTracker/backend/internal/model"
	"github.com/Ravierin/BudgetTracker/backend/internal/repository"
	"github.com/Ravierin/BudgetTracker/backend/internal/service"
	"github.com/Ravierin/BudgetTracker/backend/pkg/websocket"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	router            *mux.Router
	positionService   *service.PositionService
	withdrawalService *service.WithdrawalService
	incomeService     *service.MonthlyIncomeService
	apiKeyService     *service.APIKeyService
	balanceService    *service.BalanceService
	positionRepo      *repository.PositionRepository
	bybitClient       *api.BybitClient
	mexcClient        *api.MEXClient
	wsHub             *websocket.Hub
}

func NewServer(
	positionService *service.PositionService,
	withdrawalService *service.WithdrawalService,
	incomeService *service.MonthlyIncomeService,
	apiKeyService *service.APIKeyService,
	balanceService *service.BalanceService,
	positionRepo *repository.PositionRepository,
	bybitClient *api.BybitClient,
	mexcClient *api.MEXClient,
) *Server {
	hub := websocket.NewHub()
	go hub.Run()

	s := &Server{
		router:            mux.NewRouter(),
		positionService:   positionService,
		withdrawalService: withdrawalService,
		incomeService:     incomeService,
		apiKeyService:     apiKeyService,
		balanceService:    balanceService,
		positionRepo:      positionRepo,
		bybitClient:       bybitClient,
		mexcClient:        mexcClient,
		wsHub:             hub,
	}

	s.setupRoutes()
	
	// Initial sync of positions from all exchanges
	go s.initialSync()
	
	return s
}

// initialSync performs initial synchronization from all exchanges
func (s *Server) initialSync() {
	// Wait a bit for server to start
	time.Sleep(2 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Get all active API keys
	apiKeys, err := s.apiKeyService.GetAllAPIKeys(ctx)
	if err != nil {
		log.Printf("Failed to get API keys for initial sync: %v", err)
		return
	}

	log.Printf("Found %d API keys configured", len(apiKeys))

	totalSynced := 0
	for _, key := range apiKeys {
		if !key.IsActive {
			log.Printf("[%s] Skipping inactive key", key.Exchange)
			continue
		}

		log.Printf("[%s] Starting initial sync...", key.Exchange)
		synced := s.syncExchange(ctx, key.Exchange, key.APIKey, key.APISecret)
		totalSynced += synced
		log.Printf("[%s] Initial sync completed: %d positions", key.Exchange, synced)
	}

	log.Printf("Initial sync completed. Total positions synced: %d", totalSynced)
}

// syncExchange syncs positions for a single exchange
func (s *Server) syncExchange(ctx context.Context, exchangeName, apiKey, apiSecret string) int {
	var positions []model.Position
	var err error

	switch exchangeName {
	case "bybit":
		client := api.NewBybitClient(apiKey, apiSecret)
		positions, err = client.GetPositionsWithContext(ctx)
	case "mexc":
		client := api.NewMEXClient(apiKey, apiSecret)
		positions, err = client.GetPositionsWithContext(ctx)
	default:
		return 0
	}

	if err != nil {
		// Skip temporary errors silently
		if !containsTemporaryError(err.Error()) {
			log.Printf("[%s] Sync error: %v", exchangeName, err)
		}
		return 0
	}

	if len(positions) > 0 {
		if err := s.positionService.SavePositionsBatch(ctx, positions); err != nil {
			log.Printf("[%s] Failed to save positions: %v", exchangeName, err)
			return 0
		}
		log.Printf("[%s] Synced %d positions", exchangeName, len(positions))
	}

	return len(positions)
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

	balanceHandler := handler.NewBalanceHandler(s.balanceService)
	api.HandleFunc("/balance", balanceHandler.GetBalance).Methods("GET")

	testHandler := handler.NewTestHandler(s.positionRepo)
	api.HandleFunc("/test/generate", testHandler.GenerateTestPositions).Methods("POST")
	api.HandleFunc("/test/clear", testHandler.DeleteTestPositions).Methods("POST")

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
	log.Printf("[%s] Starting sync...", s.exchangeName)
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	apiKey, err := s.apiKeyService.GetAPIKey(ctx, s.exchangeName)
	if err != nil {
		log.Printf("[%s] Failed to get API key: %v", s.exchangeName, err)
		return // Skip silently
	}

	if apiKey.APIKey == "" || apiKey.APISecret == "" {
		log.Printf("[%s] API keys not configured", s.exchangeName)
		return // Keys not configured, skip silently
	}

	var positions []model.Position
	switch s.exchangeName {
	case "bybit":
		client := api.NewBybitClient(apiKey.APIKey, apiKey.APISecret)
		positions, err = client.GetPositionsWithContext(ctx)
	case "mexc":
		client := api.NewMEXClient(apiKey.APIKey, apiKey.APISecret)
		positions, err = client.GetPositionsWithContext(ctx)
	default:
		return
	}

	if err != nil {
		// Log only significant errors, skip rate limits/temporary issues
		errMsg := err.Error()
		if containsTemporaryError(errMsg) {
			return // Don't log temporary errors
		}
		log.Printf("[%s] Sync error: %v", s.exchangeName, err)
		return
	}

	if len(positions) > 0 {
		if err := s.positionService.SavePositionsBatch(ctx, positions); err != nil {
			log.Printf("[%s] Failed to save positions: %v", s.exchangeName, err)
			return
		}
		log.Printf("[%s] Synced %d positions", s.exchangeName, len(positions))
	} else {
		log.Printf("[%s] No positions found", s.exchangeName)
	}

	// Broadcast update
	message := map[string]interface{}{
		"type":      "positions_update",
		"positions": positions,
		"count":     len(positions),
		"exchange":  s.exchangeName,
	}
	s.wsHub.Broadcast(message)
}

// containsTemporaryError checks if error is temporary (rate limit, timeout, etc.)
func containsTemporaryError(errMsg string) bool {
	temporaryErrors := []string{
		"9999", // MEXC generic error (often rate limit)
		"10001", // Bybit rate limit
		"10006", // Bybit rate limit
		"10014", // Bybit rate limit
		"timeout",
		"deadline exceeded",
		"earlier than 2 years", // Bybit time range error
		"rate limit",
	}
	for _, temp := range temporaryErrors {
		if strings.Contains(errMsg, temp) {
			return true
		}
	}
	return false
}
