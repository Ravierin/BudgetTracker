package handler

import (
	"BudgetTracker/backend/internal/model"
	"BudgetTracker/backend/internal/service"
	"BudgetTracker/backend/pkg/websocket"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type MonthlyIncomeHandler struct {
	service *service.MonthlyIncomeService
	wsHub   *websocket.Hub
}

func NewMonthlyIncomeHandler(service *service.MonthlyIncomeService, wsHub *websocket.Hub) *MonthlyIncomeHandler {
	return &MonthlyIncomeHandler{
		service: service,
		wsHub:   wsHub,
	}
}

func (h *MonthlyIncomeHandler) GetAllMonthlyIncomes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	exchange := r.URL.Query().Get("exchange")
	var incomes []model.MonthlyIncome
	var err error

	if exchange != "" {
		incomes, err = h.service.GetIncomesByExchange(ctx, exchange)
	} else {
		incomes, err = h.service.GetAllMonthlyIncomes(ctx)
	}

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

func (h *MonthlyIncomeHandler) CreateMonthlyIncome(w http.ResponseWriter, r *http.Request) {
	var income model.MonthlyIncome
	if err := json.NewDecoder(r.Body).Decode(&income); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if income.CreatedAt.IsZero() {
		income.CreatedAt = time.Now()
	}

	ctx := r.Context()
	if err := h.service.SaveMonthlyIncome(ctx, income); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.wsHub.Broadcast(map[string]interface{}{
		"type": "monthly_income_created",
		"data": income,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(income)
}

func (h *MonthlyIncomeHandler) DeleteMonthlyIncome(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.service.DeleteMonthlyIncome(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.wsHub.Broadcast(map[string]interface{}{
		"type":     "monthly_income_deleted",
		"incomeId": id,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
