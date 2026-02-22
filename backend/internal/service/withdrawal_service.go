package service

import (
	"BudgetTracker/backend/internal/model"
	"BudgetTracker/backend/internal/repository"
	"context"
	"time"
)

type WithdrawalService struct {
	repo *repository.WithdrawalRepository
}

func NewWithdrawalService(repo *repository.WithdrawalRepository) *WithdrawalService {
	return &WithdrawalService{repo: repo}
}

func (s *WithdrawalService) SaveWithdrawal(ctx context.Context, withdrawal model.Withdrawal) error {
	return s.repo.SaveWithdrawal(ctx, withdrawal)
}

func (s *WithdrawalService) GetAllWithdrawals(ctx context.Context) ([]model.Withdrawal, error) {
	return s.repo.GetAllWithdrawals(ctx)
}

func (s *WithdrawalService) GetWithdrawalsByExchange(ctx context.Context, exchange string) ([]model.Withdrawal, error) {
	return s.repo.GetWithdrawalsByExchange(ctx, exchange)
}

func (s *WithdrawalService) GetWithdrawalsByDateRange(ctx context.Context, start, end time.Time) ([]model.Withdrawal, error) {
	return s.repo.GetWithdrawalsByDateRange(ctx, start, end)
}

func (s *WithdrawalService) DeleteWithdrawal(ctx context.Context, id int) error {
	return s.repo.DeleteWithdrawal(ctx, id)
}

func (s *WithdrawalService) CalculateTotalWithdrawals(ctx context.Context, exchange string) (float64, error) {
	withdrawals, err := s.repo.GetAllWithdrawals(ctx)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, w := range withdrawals {
		if exchange == "" || w.Exchange == exchange {
			total += w.Amount
		}
	}

	return total, nil
}
