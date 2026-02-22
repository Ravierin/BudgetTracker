package api

import "BudgetTracker/backend/internal/model"

type ExchangeClient interface {
	GetPositions() ([]model.Position, error)
}
