package api

import (
	"BudgetTracker/backend/internal/model"
	"context"
)

type ExchangeClient interface {
	GetPositions() ([]model.Position, error)
	GetPositionsWithContext(ctx context.Context) ([]model.Position, error)
}
