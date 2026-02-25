package handler

import (
	"BudgetTracker/backend/internal/model"
	"BudgetTracker/backend/internal/service"
	"BudgetTracker/backend/pkg/websocket"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type MonthlyIncomeHandler struct {
	positionService *service.PositionService
	wsHub           *websocket.Hub
}

func NewMonthlyIncomeHandler(positionService *service.PositionService, wsHub *websocket.Hub) *MonthlyIncomeHandler {
	return &MonthlyIncomeHandler{
		positionService: positionService,
		wsHub:           wsHub,
	}
}

func (h *MonthlyIncomeHandler) GetAllMonthlyIncomes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	exchange := r.URL.Query().Get("exchange")

	// Aggregate monthly PnL from positions
	incomes, err := h.positionService.AggregateMonthlyPnl(ctx, exchange)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if incomes == nil {
		incomes = []model.MonthlyIncome{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(incomes)
}

func (h *MonthlyIncomeHandler) GetMonthlyIncome(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := time.Parse("2006-01", vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	incomes, err := h.positionService.AggregateMonthlyPnl(ctx, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, income := range incomes {
		if income.CreatedAt.Year() == id.Year() && income.CreatedAt.Month() == id.Month() {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(income)
			return
		}
	}

	http.Error(w, "Monthly income not found", http.StatusNotFound)
}
