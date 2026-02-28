package handler

import (
	"github.com/Ravierin/BudgetTracker/backend/internal/model"
	"github.com/Ravierin/BudgetTracker/backend/internal/service"
	"encoding/json"
	"net/http"
)

type BalanceHandler struct {
	service *service.BalanceService
}

func NewBalanceHandler(service *service.BalanceService) *BalanceHandler {
	return &BalanceHandler{service: service}
}

func (h *BalanceHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	totalBalance, exchangeBalances, err := h.service.GetTotalBalance(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"totalBalance":     totalBalance,
		"exchangeBalances": exchangeBalances,
	}

	if exchangeBalances == nil {
		response["exchangeBalances"] = []model.ExchangeBalance{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
