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
