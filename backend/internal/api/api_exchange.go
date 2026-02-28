package api

import (
	"github.com/Ravierin/BudgetTracker/backend/internal/model"
	"context"
)

type ExchangeClient interface {
	GetPositions() ([]model.Position, error)
	GetPositionsWithContext(ctx context.Context) ([]model.Position, error)
	GetBalance(ctx context.Context) (float64, error)
}
