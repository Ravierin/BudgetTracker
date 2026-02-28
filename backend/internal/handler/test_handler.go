package handler

import (
	"github.com/Ravierin/BudgetTracker/backend/internal/model"
	"github.com/Ravierin/BudgetTracker/backend/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type TestHandler struct {
	positionRepo *repository.PositionRepository
}

func NewTestHandler(positionRepo *repository.PositionRepository) *TestHandler {
	return &TestHandler{positionRepo: positionRepo}
}

func (h *TestHandler) GenerateTestPositions(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()

	// Generate 100 test positions
	var positions []model.Position
	exchanges := []string{"bybit", "mexc"}
	symbols := []string{"BTCUSDT", "ETHUSDT", "SOLUSDT", "XRPUSDT", "BNBUSDT"}
	sides := []string{"Buy", "Sell"}

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 100; i++ {
		exchange := exchanges[rand.Intn(len(exchanges))]
		symbol := symbols[rand.Intn(len(symbols))]
		side := sides[rand.Intn(len(sides))]
		
		// Random volume between 100 and 10000 USDT
		volume := float64(rand.Intn(9900) + 100)
		
		// Random leverage between 1 and 20
		leverage := rand.Intn(20) + 1
		
		// Random PnL between -500 and 1000
		pnl := float64(rand.Intn(1500) - 500)
		
		// Random date within last 6 months
		daysAgo := rand.Intn(180)
		date := time.Now().AddDate(0, 0, -daysAgo)

		positions = append(positions, model.Position{
			OrderID:   fmt.Sprintf("test_position_%d_%d", time.Now().UnixNano(), i),
			Exchange:  exchange,
			Symbol:    symbol,
			Volume:    volume,
			Leverage:  leverage,
			ClosedPnl: pnl,
			Side:      side,
			UpdatedAt: date,
		})
	}

	// Save all positions
	if err := h.positionRepo.SavePositionBatch(ctx, positions); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"count":  len(positions),
		"message": fmt.Sprintf("Generated %d test positions", len(positions)),
	})
}

func (h *TestHandler) DeleteTestPositions(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()

	// Get all positions
	positions, err := h.positionRepo.GetAllPositions(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Debug: log all order IDs
	for _, p := range positions {
		fmt.Printf("Position ID: %d, OrderID: %s\n", p.ID, p.OrderID)
	}

	// Delete test positions (those with order_id starting with "test_position")
	deleted := 0
	for _, p := range positions {
		// Check if order_id starts with "test_position"
		if len(p.OrderID) >= 14 && p.OrderID[:14] == "test_position" {
			if err := h.positionRepo.DeletePosition(ctx, p.ID); err != nil {
				continue
			}
			deleted++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"count":   deleted,
		"message": fmt.Sprintf("Deleted %d test positions", deleted),
	})
}
