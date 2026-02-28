package service

import (
	"github.com/Ravierin/BudgetTracker/backend/internal/model"
	"github.com/Ravierin/BudgetTracker/backend/internal/repository"
	"context"
	"time"
)

type MonthlyIncomeService struct {
	repo *repository.MonthlyIncomeRepository
}

func NewMonthlyIncomeService(repo *repository.MonthlyIncomeRepository) *MonthlyIncomeService {
	return &MonthlyIncomeService{repo: repo}
}

func (s *MonthlyIncomeService) SaveMonthlyIncome(ctx context.Context, income model.MonthlyIncome) error {
	return s.repo.SaveMonthlyIncome(ctx, income)
}

func (s *MonthlyIncomeService) GetAllMonthlyIncomes(ctx context.Context) ([]model.MonthlyIncome, error) {
	return s.repo.GetAllMonthlyIncomes(ctx)
}

func (s *MonthlyIncomeService) GetIncomesByExchange(ctx context.Context, exchange string) ([]model.MonthlyIncome, error) {
	return s.repo.GetIncomesByExchange(ctx, exchange)
}

func (s *MonthlyIncomeService) GetIncomesByDateRange(ctx context.Context, start, end time.Time) ([]model.MonthlyIncome, error) {
	return s.repo.GetIncomesByDateRange(ctx, start, end)
}

func (s *MonthlyIncomeService) DeleteMonthlyIncome(ctx context.Context, id int) error {
	return s.repo.DeleteMonthlyIncome(ctx, id)
}

func (s *MonthlyIncomeService) CalculateTotalIncome(ctx context.Context, exchange string) (float64, error) {
	incomes, err := s.repo.GetAllMonthlyIncomes(ctx)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, i := range incomes {
		if exchange == "" || i.Exchange == exchange {
			total += i.PNL
		}
	}

	return total, nil
}

func (s *MonthlyIncomeService) CalculateMonthlyTotal(ctx context.Context, year int, month time.Month, exchange string) (float64, error) {
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	incomes, err := s.repo.GetIncomesByDateRange(ctx, start, end)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, i := range incomes {
		if exchange == "" || i.Exchange == exchange {
			total += i.PNL
		}
	}

	return total, nil
}
