package api

import (
	"BudgetTracker/backend/internal/model"
	"errors"
)

type MEXClient struct {
	apiKey    string
	apiSecret string
}

func NewMEXClient(apiKey, apiSecret string) *MEXClient {
	return &MEXClient{
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
}

func (m *MEXClient) GetPositions() ([]model.Position, error) {
	return nil, errors.New("MEXC integration not implemented yet")
}

var _ ExchangeClient = (*MEXClient)(nil)
