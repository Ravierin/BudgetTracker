package api

import (
	"BudgetTracker/backend/internal/model"
	"context"
	"strconv"
	"time"

	bybit "github.com/bybit-exchange/bybit.go.api"
)

type BybitClient struct {
	bybit *bybit.Client
}

func NewBybitClient(apiKey, apiSecretKey string) *BybitClient {
	bybit := bybit.NewBybitHttpClient(apiKey, apiSecretKey)
	return &BybitClient{bybit: bybit}
}

func (b *BybitClient) GetPositions() ([]model.Position, error) {
	return b.GetPositionsWithContext(context.Background())
}

func (b *BybitClient) GetPositionsWithContext(ctx context.Context) ([]model.Position, error) {
	params := map[string]interface{}{
		"category": "linear",
		"limit":    100,
	}

	result, err := b.bybit.NewClassicalBybitServiceWithParams(params).GetClosePnl(ctx)
	if err != nil {
		return nil, err
	}

	list, ok := result.Result.(map[string]interface{})
	if !ok {
		return nil, nil
	}

	items, ok := list["list"].([]interface{})
	if !ok {
		return nil, nil
	}

	var positions []model.Position
	for _, item := range items {
		posMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		orderID, _ := posMap["orderId"].(string)
		symbol, _ := posMap["symbol"].(string)
		side, _ := posMap["side"].(string)

		cumExitValue, _ := strconv.ParseFloat(posMap["cumExitValue"].(string), 64)
		qty, _ := strconv.ParseFloat(posMap["closedSize"].(string), 64)
		leverage, _ := strconv.Atoi(posMap["leverage"].(string))
		closedPnl, _ := strconv.ParseFloat(posMap["closedPnl"].(string), 64)

		updatedTime, _ := strconv.ParseFloat(posMap["updatedTime"].(string), 64)
		date := time.UnixMilli(int64(updatedTime))

		positions = append(positions, model.Position{
			OrderID:      orderID,
			Exchange:     "bybit",
			Symbol:       symbol,
			CumExitValue: cumExitValue,
			Quantity:     qty,
			Leverage:     leverage,
			ClosedPnl:    closedPnl,
			Side:         side,
			UpdatedAt:    date,
		})
	}

	return positions, nil
}
