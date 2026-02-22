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

type WithdrawalHandler struct {
	service *service.WithdrawalService
	wsHub   *websocket.Hub
}

func NewWithdrawalHandler(service *service.WithdrawalService, wsHub *websocket.Hub) *WithdrawalHandler {
	return &WithdrawalHandler{
		service: service,
		wsHub:   wsHub,
	}
}

func (h *WithdrawalHandler) GetAllWithdrawals(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	exchange := r.URL.Query().Get("exchange")
	var withdrawals []model.Withdrawal
	var err error

	if exchange != "" {
		withdrawals, err = h.service.GetWithdrawalsByExchange(ctx, exchange)
	} else {
		withdrawals, err = h.service.GetAllWithdrawals(ctx)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if withdrawals == nil {
		withdrawals = []model.Withdrawal{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(withdrawals)
}

func (h *WithdrawalHandler) CreateWithdrawal(w http.ResponseWriter, r *http.Request) {
	var withdrawal model.Withdrawal
	if err := json.NewDecoder(r.Body).Decode(&withdrawal); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if withdrawal.CreatedAt.IsZero() {
		withdrawal.CreatedAt = time.Now()
	}

	ctx := r.Context()
	if err := h.service.SaveWithdrawal(ctx, withdrawal); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.wsHub.Broadcast(map[string]interface{}{
		"type": "withdrawal_created",
		"data": withdrawal,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(withdrawal)
}

func (h *WithdrawalHandler) DeleteWithdrawal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.service.DeleteWithdrawal(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.wsHub.Broadcast(map[string]interface{}{
		"type":         "withdrawal_deleted",
		"withdrawalId": id,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
