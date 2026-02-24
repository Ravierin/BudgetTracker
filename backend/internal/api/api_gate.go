package api

import (
	"BudgetTracker/backend/internal/model"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type GateClient struct {
	apiKey    string
	apiSecret string
	baseURL   string
	client    *http.Client
}

func NewGateClient(apiKey, apiSecret string) *GateClient {
	return &GateClient{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		baseURL:   "https://api.gateio.ws/api/v4",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (g *GateClient) sign(method, urlPath, queryString, body, timestamp string) string {
	payload := method + "\n" + urlPath + "\n" + queryString + "\n" + body + "\n" + timestamp
	h := hmac.New(sha256.New, []byte(g.apiSecret))
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

func (g *GateClient) doRequest(ctx context.Context, method, endpoint, queryString string) ([]byte, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := g.sign(method, endpoint, queryString, "", timestamp)

	url := g.baseURL + endpoint
	if queryString != "" {
		url += "?" + queryString
	}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("KEY", g.apiKey)
	req.Header.Set("SIGN", signature)
	req.Header.Set("Timestamp", timestamp)
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Gate.io API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	return body, nil
}

func (g *GateClient) GetPositions() ([]model.Position, error) {
	return g.GetPositionsWithContext(context.Background())
}

func (g *GateClient) GetPositionsWithContext(ctx context.Context) ([]model.Position, error) {
	return []model.Position{}, nil
}

var _ ExchangeClient = (*GateClient)(nil)
