package service

import (
	"BudgetTracker/backend/internal/model"
	"BudgetTracker/backend/internal/repository"
	"context"
	"time"
)

type PositionService struct {
	repo *repository.PositionRepository
}

func NewPositionService(repo *repository.PositionRepository) *PositionService {
	return &PositionService{repo: repo}
}

func (s *PositionService) SavePosition(ctx context.Context, position model.Position) error {
	return s.repo.SavePosition(ctx, position)
}

func (s *PositionService) SavePositionsBatch(ctx context.Context, positions []model.Position) error {
	return s.repo.SavePositionBatch(ctx, positions)
}

func (s *PositionService) GetAllPositions(ctx context.Context) ([]model.Position, error) {
	return s.repo.GetAllPositions(ctx)
}

func (s *PositionService) GetPositionsByExchange(ctx context.Context, exchange string) ([]model.Position, error) {
	return s.repo.GetPositionsByExchange(ctx, exchange)
}

func (s *PositionService) GetPositionsByDateRange(ctx context.Context, start, end time.Time) ([]model.Position, error) {
	return s.repo.GetPositionsByDateRange(ctx, start, end)
}

func (s *PositionService) DeletePosition(ctx context.Context, id int) error {
	return s.repo.DeletePosition(ctx, id)
}

func (s *PositionService) CalculateTotalPnl(ctx context.Context, exchange string) (float64, error) {
	positions, err := s.repo.GetAllPositions(ctx)
	if err != nil {
		return 0, err
	}

	var totalPnl float64
	for _, p := range positions {
		if exchange == "" || p.Exchange == exchange {
			totalPnl += p.ClosedPnl
		}
	}

	return totalPnl, nil
}

func (s *PositionService) CalculateMonthlyPnl(ctx context.Context, year int, month time.Month, exchange string) (float64, error) {
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	positions, err := s.repo.GetPositionsByDateRange(ctx, start, end)
	if err != nil {
		return 0, err
	}

	var totalPnl float64
	for _, p := range positions {
		if exchange == "" || p.Exchange == exchange {
			totalPnl += p.ClosedPnl
		}
	}

	return totalPnl, nil
}

// AggregateMonthlyPnl aggregates PnL by month from positions
// Returns monthly income data grouped by year-month
func (s *PositionService) AggregateMonthlyPnl(ctx context.Context, exchange string) ([]model.MonthlyIncome, error) {
	positions, err := s.repo.GetAllPositions(ctx)
	if err != nil {
		return nil, err
	}

	// Group by year-month and exchange
	monthlyMap := make(map[string]*model.MonthlyIncome)
	
	for _, p := range positions {
		if exchange != "" && p.Exchange != exchange {
			continue
		}

		// Get year-month key (e.g., "2026-01-mexc")
		var key string
		if exchange == "" {
			key = p.UpdatedAt.Format("2006-01") + "-" + p.Exchange
		} else {
			key = p.UpdatedAt.Format("2006-01")
		}
		
		monthStart := time.Date(p.UpdatedAt.Year(), p.UpdatedAt.Month(), 1, 0, 0, 0, 0, p.UpdatedAt.Location())
		
		if existing, ok := monthlyMap[key]; ok {
			existing.PNL += p.ClosedPnl
			existing.Amount += p.Volume
		} else {
			monthlyMap[key] = &model.MonthlyIncome{
				Exchange:  p.Exchange,
				PNL:       p.ClosedPnl,
				Amount:    p.Volume,
				CreatedAt: monthStart,
			}
		}
	}

	// Convert map to slice
	result := make([]model.MonthlyIncome, 0, len(monthlyMap))
	for _, income := range monthlyMap {
		result = append(result, *income)
	}

	return result, nil
}
