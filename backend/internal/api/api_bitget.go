package api

import (
	"BudgetTracker/backend/internal/model"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type BitgetClient struct {
	apiKey    string
	apiSecret string
	baseURL   string
	client    *http.Client
}

func NewBitgetClient(apiKey, apiSecret string) *BitgetClient {
	return &BitgetClient{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		baseURL:   "https://api.bitget.com",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (b *BitgetClient) sign(timestamp, method, requestPath, body string) string {
	payload := timestamp + method + requestPath + body
	h := hmac.New(sha256.New, []byte(b.apiSecret))
	h.Write([]byte(payload))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (b *BitgetClient) doRequest(ctx context.Context, method, endpoint, queryString, body string) ([]byte, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := b.sign(timestamp, method, endpoint, body)

	url := b.baseURL + endpoint
	if queryString != "" {
		url += "?" + queryString
	}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("ACCESS-KEY", b.apiKey)
	req.Header.Set("ACCESS-SIGN", signature)
	req.Header.Set("ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("Content-Type", "application/json")

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bitget API error: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (b *BitgetClient) GetPositions() ([]model.Position, error) {
	return b.GetPositionsWithContext(context.Background())
}

func (b *BitgetClient) GetPositionsWithContext(ctx context.Context) ([]model.Position, error) {
	return []model.Position{}, nil
}

var _ ExchangeClient = (*BitgetClient)(nil)
