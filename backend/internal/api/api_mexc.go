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
	"log"
	"net/http"
	"sort"
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
		baseURL:   "https://api.mexc.com",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// signV1 generates signature for /api/v1/private/* endpoints
// Signature = HMAC-SHA256(ApiKey + Request-Time + queryParams, apiSecret)
func (m *MEXClient) signV1(query string, timestamp int64) string {
	payload := fmt.Sprintf("%s%d%s", m.apiKey, timestamp, query)
	h := hmac.New(sha256.New, []byte(m.apiSecret))
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

// doRequestV1 makes request to /api/v1/private/* endpoints
// Uses ApiKey/Signature/Request-Time headers for authentication
func (m *MEXClient) doRequestV1(ctx context.Context, endpoint string, params map[string]string) ([]byte, error) {
	timestamp := time.Now().UnixMilli()

	// Build query string in ALPHABETICAL order for signature (required by MEXC)
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys) // Sort alphabetically

	query := ""
	for _, k := range keys {
		if query != "" {
			query += "&"
		}
		query += fmt.Sprintf("%s=%s", k, params[k])
	}

	signature := m.signV1(query, timestamp)

	req, err := http.NewRequestWithContext(ctx, "GET", m.baseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("ApiKey", m.apiKey)
	req.Header.Set("Request-Time", strconv.FormatInt(timestamp, 10))
	req.Header.Set("Signature", signature)
	req.Header.Set("Content-Type", "application/json")

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
		return nil, fmt.Errorf("MEXC API v1 error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	return body, nil
}

func (m *MEXClient) GetPositions() ([]model.Position, error) {
	return m.GetPositionsWithContext(context.Background())
}

func (m *MEXClient) GetPositionsWithContext(ctx context.Context) ([]model.Position, error) {
	var allPositions []model.Position
	page := 1

	for {
		// Parameters must be in alphabetical order for signature
		params := map[string]string{
			"page_num":  strconv.Itoa(page),
			"page_size": "100",
		}

		body, err := m.doRequestV1(ctx, "/api/v1/private/position/list/history_positions", params)
		if err != nil {
			return nil, err
		}

		// Response format for v1 API - data is an array directly
		var resp struct {
			Success bool            `json:"success"`
			Code    int             `json:"code"`
			Data    json.RawMessage `json:"data"`
		}

		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, err
		}

		if !resp.Success || resp.Code != 0 {
			return nil, fmt.Errorf("MEXC API v1 error: success=%v, code=%d", resp.Success, resp.Code)
		}

		// Try to unmarshal data as array first (MEXC v1 returns array directly)
		var positionsList []struct {
			PositionID      int64   `json:"positionId"`
			Symbol          string  `json:"symbol"`
			PositionType    int     `json:"positionType"` // 1=Buy, 2=Sell
			CloseVol        float64 `json:"closeVol"`
			CloseAvgPrice   float64 `json:"closeAvgPrice"`
			OpenAvgPrice    float64 `json:"openAvgPrice"`
			Leverage        int     `json:"leverage"`
			CloseProfitLoss float64 `json:"closeProfitLoss"`
			CreateTime      int64   `json:"createTime"`
			UpdateTime      int64   `json:"updateTime"`
		}

		// Try as array first
		if err := json.Unmarshal(resp.Data, &positionsList); err != nil {
			// Try as object with list field
			var dataObj struct {
				List []struct {
					OrderID      string `json:"orderId"`
					Symbol       string `json:"symbol"`
					Side         int    `json:"side"`
					DealAvgPrice string `json:"dealAvgPrice"`
					DealQty      string `json:"dealQty"`
					Leverage     int    `json:"leverage"`
					Profit       string `json:"profit"`
					CreateTime   int64  `json:"createTime"`
				} `json:"list"`
			}
			if err := json.Unmarshal(resp.Data, &dataObj); err != nil {
				return nil, fmt.Errorf("failed to parse MEXC response: %w", err)
			}
			positionsList = nil // Not used in this case
		}

		if len(positionsList) == 0 {
			break // No more positions
		}

		for _, pos := range positionsList {
			side := "Buy"
			if pos.PositionType == 2 {
				side = "Sell"
			}

			// Volume = CloseVol (closed quantity)
			volume := pos.CloseVol
			// Margin = Volume / Leverage (initial margin)
			margin := volume / float64(pos.Leverage)

			allPositions = append(allPositions, model.Position{
				OrderID:   fmt.Sprintf("%d", pos.PositionID),
				Exchange:  "mexc",
				Symbol:    pos.Symbol,
				Volume:    volume,
				Margin:    margin,
				Leverage:  pos.Leverage,
				ClosedPnl: pos.CloseProfitLoss,
				Side:      side,
				UpdatedAt: time.UnixMilli(pos.UpdateTime),
			})
		}

		// If we got less than page_size, we've reached the end
		if len(positionsList) < 100 {
			break
		}

		page++
		time.Sleep(50 * time.Millisecond) // Rate limiting
	}

	return allPositions, nil
}

// GetBalance returns total futures account balance in USDT
func (m *MEXClient) GetBalance(ctx context.Context) (float64, error) {
	// MEXC futures balance endpoint
	params := map[string]string{}

	body, err := m.doRequestV1(ctx, "/api/v1/private/account/balance", params)
	if err != nil {
		return 0, err
	}

	var resp struct {
		Success bool `json:"success"`
		Code    int  `json:"code"`
		Data    struct {
			Balance string `json:"balance"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return 0, err
	}

	if !resp.Success || resp.Code != 0 {
		return 0, fmt.Errorf("MEXC API error: code=%d", resp.Code)
	}

	balance, _ := strconv.ParseFloat(resp.Data.Balance, 64)
	log.Printf("[mexc] Balance: %.2f USDT", balance)
	return balance, nil
}

var _ ExchangeClient = (*MEXClient)(nil)
