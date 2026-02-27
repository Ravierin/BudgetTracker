package service

import (
	"BudgetTracker/backend/internal/api"
	"BudgetTracker/backend/internal/model"
	"context"
	"log"
)

type BalanceService struct {
	apiKeyService *APIKeyService
}

func NewBalanceService(apiKeyService *APIKeyService) *BalanceService {
	return &BalanceService{apiKeyService: apiKeyService}
}

// GetTotalBalance returns total balance across all configured exchanges
func (s *BalanceService) GetTotalBalance(ctx context.Context) (float64, []model.ExchangeBalance, error) {
	apiKeys, err := s.apiKeyService.GetAllAPIKeys(ctx)
	if err != nil {
		return 0, nil, err
	}

	var totalBalance float64
	var exchangeBalances []model.ExchangeBalance

	for _, key := range apiKeys {
		if !key.IsActive || key.APIKey == "" || key.APISecret == "" {
			continue
		}

		balance, err := s.getExchangeBalance(ctx, key.Exchange, key.APIKey, key.APISecret)
		if err != nil {
			log.Printf("[%s] Failed to get balance: %v", key.Exchange, err)
			continue
		}

		if balance > 0 {
			exchangeBalances = append(exchangeBalances, model.ExchangeBalance{
				Exchange: key.Exchange,
				Balance:  balance,
			})
			totalBalance += balance
		}
	}

	return totalBalance, exchangeBalances, nil
}

func (s *BalanceService) getExchangeBalance(ctx context.Context, exchange, apiKey, apiSecret string) (float64, error) {
	var client api.ExchangeClient

	switch exchange {
	case "bybit":
		client = api.NewBybitClient(apiKey, apiSecret)
	case "mexc":
		client = api.NewMEXClient(apiKey, apiSecret)
	case "gate":
		client = api.NewGateClient(apiKey, apiSecret)
	case "bitget":
		client = api.NewBitgetClient(apiKey, apiSecret)
	default:
		return 0, nil
	}

	return client.GetBalance(ctx)
}
