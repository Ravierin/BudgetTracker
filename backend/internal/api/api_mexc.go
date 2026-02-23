package api

import (
	"BudgetTracker/backend/internal/model"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type MEXClient struct {
	apiKey    string
	apiSecret string
	baseURL   string
	client    *http.Client
}

func NewMEXClient(apiKey, apiSecret string) *MEXClient {
	return &MEXClient{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		baseURL:   "https://contract.mexc.com",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (m *MEXClient) sign(query string, timestamp int64) string {
	payload := fmt.Sprintf("%s%s", query, strconv.FormatInt(timestamp, 10))
	h := hmac.New(sha256.New, []byte(m.apiSecret))
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

func (m *MEXClient) doRequest(ctx context.Context, endpoint string, params map[string]string) ([]byte, error) {
	timestamp := time.Now().UnixMilli()

	query := ""
	for k, v := range params {
		if query != "" {
			query += "&"
		}
		query += fmt.Sprintf("%s=%s", k, v)
	}

	signature := m.sign(query, timestamp)

	req, err := http.NewRequestWithContext(ctx, "GET", m.baseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("ApiKey", m.apiKey)
	req.Header.Set("Request-Time", strconv.FormatInt(timestamp, 10))
	req.Header.Set("Signature", signature)

	if query != "" {
		req.URL.RawQuery = query
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MEXC API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	return body, nil
}

func (m *MEXClient) GetPositions() ([]model.Position, error) {
	ctx := context.Background()

	params := map[string]string{
		"page_size": "100",
		"page_num":  "1",
	}

	body, err := m.doRequest(ctx, "/api/v1/contract/trade/hisorders", params)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Code int `json:"code"`
		Data struct {
			HisOrders []struct {
				OrderID      string `json:"orderId"`
				Symbol       string `json:"symbol"`
				Side         int    `json:"side"`
				OrderType    int    `json:"orderType"`
				DealAvgPrice string `json:"dealAvgPrice"`
				DealQty      string `json:"dealQty"`
				Leverage     string `json:"leverage"`
				Profit       string `json:"profit"`
				CreateTime   int64  `json:"createTime"`
			} `json:"hisOrders"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("MEXC API error: code=%d", resp.Code)
	}

	var positions []model.Position
	for _, order := range resp.Data.HisOrders {
		avgPrice, _ := strconv.ParseFloat(order.DealAvgPrice, 64)
		qty, _ := strconv.ParseFloat(order.DealQty, 64)
		leverage, _ := strconv.Atoi(order.Leverage)
		profit, _ := strconv.ParseFloat(order.Profit, 64)

		side := "Buy"
		if order.Side == 2 {
			side = "Sell"
		}

		cumExitValue := avgPrice * qty

		positions = append(positions, model.Position{
			OrderID:      order.OrderID,
			Exchange:     "mexc",
			Symbol:       order.Symbol,
			CumExitValue: cumExitValue,
			Quantity:     qty,
			Leverage:     leverage,
			ClosedPnl:    profit,
			Side:         side,
			UpdatedAt:    time.UnixMilli(order.CreateTime),
		})
	}

	return positions, nil
}

var _ ExchangeClient = (*MEXClient)(nil)
