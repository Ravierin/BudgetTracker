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

type PositionHandler struct {
	service *service.PositionService
	wsHub   *websocket.Hub
}

func NewPositionHandler(service *service.PositionService, wsHub *websocket.Hub) *PositionHandler {
	return &PositionHandler{
		service: service,
		wsHub:   wsHub,
	}
}

func (h *PositionHandler) GetAllPositions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	exchange := r.URL.Query().Get("exchange")
	var positions []model.Position
	var err error

	if exchange != "" {
		positions, err = h.service.GetPositionsByExchange(ctx, exchange)
	} else {
		positions, err = h.service.GetAllPositions(ctx)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if positions == nil {
		positions = []model.Position{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(positions)
}

func (h *PositionHandler) GetPosition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	positions, err := h.service.GetAllPositions(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, p := range positions {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "Position not found", http.StatusNotFound)
}

func (h *PositionHandler) CreatePosition(w http.ResponseWriter, r *http.Request) {
	var position model.Position
	if err := json.NewDecoder(r.Body).Decode(&position); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if position.UpdatedAt.IsZero() {
		position.UpdatedAt = time.Now()
	}

	ctx := r.Context()
	if err := h.service.SavePosition(ctx, position); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.wsHub.Broadcast(map[string]interface{}{
		"type": "position_created",
		"data": position,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(position)
}

func (h *PositionHandler) DeletePosition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.service.DeletePosition(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.wsHub.Broadcast(map[string]interface{}{
		"type":   "position_deleted",
		"postId": id,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

func (h *PositionHandler) SyncPositions(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Use automatic sync every 5 seconds", http.StatusNotImplemented)
}
