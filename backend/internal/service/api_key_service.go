package service

import (
	"github.com/Ravierin/BudgetTracker/backend/internal/model"
	"github.com/Ravierin/BudgetTracker/backend/internal/repository"
	"context"
)

type APIKeyService struct {
	repo *repository.APIKeyRepository
}

func NewAPIKeyService(repo *repository.APIKeyRepository) *APIKeyService {
	return &APIKeyService{repo: repo}
}

func (s *APIKeyService) GetAPIKey(ctx context.Context, exchange string) (*model.APIKey, error) {
	return s.repo.GetByExchange(ctx, exchange)
}

func (s *APIKeyService) SaveAPIKey(ctx context.Context, apiKey *model.APIKey) error {
	// Allow empty keys - user can configure exchanges gradually
	return s.repo.Upsert(ctx, apiKey)
}

func (s *APIKeyService) GetAllAPIKeys(ctx context.Context) ([]model.APIKey, error) {
	return s.repo.GetAll(ctx)
}
