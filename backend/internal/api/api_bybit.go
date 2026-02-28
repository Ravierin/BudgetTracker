package api

import (
	"BudgetTracker/backend/internal/model"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	bybit "github.com/bybit-exchange/bybit.go.api"
)

type BybitClient struct {
	bybit     *bybit.Client
	apiKey    string
	apiSecret string
}

func NewBybitClient(apiKey, apiSecretKey string) *BybitClient {
	bybit := bybit.NewBybitHttpClient(apiKey, apiSecretKey)
	return &BybitClient{
		bybit:     bybit,
		apiKey:    apiKey,
		apiSecret: apiSecretKey,
	}
}

func (b *BybitClient) GetPositions() ([]model.Position, error) {
	return b.GetPositionsWithContext(context.Background())
}

func (b *BybitClient) GetPositionsWithContext(ctx context.Context) ([]model.Position, error) {
	// Try direct API call to /v5/position/get-closed-positions first
	positions, err := b.getDirectClosedPositions(ctx)
	if err == nil && len(positions) > 0 {
		return positions, nil
	}
	
	// Fallback to GetClosePnl for UTA accounts
	return b.getClosePnl(ctx)
}

// getDirectClosedPositions calls /v5/position/get-closed-positions directly
func (b *BybitClient) getDirectClosedPositions(ctx context.Context) ([]model.Position, error) {
	log.Printf("[bybit] Trying /v5/position/get-closed-positions...")
	
	baseURL := "https://api.bybit.com"
	endpoint := "/v5/position/get-closed-positions"
	
	params := url.Values{}
	params.Set("category", "linear")
	
	queryString := params.Encode()
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signature := b.signV5("GET", endpoint, queryString, timestamp)
	
	reqURL := baseURL + endpoint + "?" + queryString
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("X-BAPI-API-KEY", b.apiKey)
	req.Header.Set("X-BAPI-SIGN", signature)
	req.Header.Set("X-BAPI-TIMESTAMP", timestamp)
	req.Header.Set("X-BAPI-RECV-WINDOW", "30000")
	
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[bybit] HTTP request error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	
	var apiResp struct {
		RetCode int         `json:"retCode"`
		RetMsg  string      `json:"retMsg"`
		Result  interface{} `json:"result"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}
	
	log.Printf("[bybit] GetClosedPositions RetCode: %d, RetMsg: %s", apiResp.RetCode, apiResp.RetMsg)
	
	if apiResp.RetCode != 0 {
		log.Printf("[bybit] API error: %s", apiResp.RetMsg)
		return nil, fmt.Errorf("API error: %s", apiResp.RetMsg)
	}
	
	result, ok := apiResp.Result.(map[string]interface{})
	if !ok {
		return nil, nil
	}
	
	items, ok := result["list"].([]interface{})
	if !ok {
		return nil, nil
	}
	
	log.Printf("[bybit] GetClosedPositions retrieved %d positions", len(items))
	
	var positions []model.Position
	for _, item := range items {
		posMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		
		orderID, _ := posMap["orderId"].(string)
		symbol, _ := posMap["symbol"].(string)
		side, _ := posMap["side"].(string)
		
		// Volume = cumEntryValue (position value at entry in USDT)
		volume, _ := strconv.ParseFloat(posMap["cumEntryValue"].(string), 64)
		leverage, _ := strconv.Atoi(posMap["leverage"].(string))
		margin := volume / float64(leverage)
		closedPnl, _ := strconv.ParseFloat(posMap["closedPnl"].(string), 64)
		
		updatedTime, _ := strconv.ParseFloat(posMap["updatedTime"].(string), 64)
		date := time.UnixMilli(int64(updatedTime))
		
		positions = append(positions, model.Position{
			OrderID:      orderID,
			Exchange:     "bybit",
			Symbol:       symbol,
			Volume:       volume,
			Margin:       margin,
			Leverage:     leverage,
			ClosedPnl:    closedPnl,
			Side:         side,
			UpdatedAt:    date,
		})
	}
	
	return positions, nil
}

// getClosePnl uses the GetClosePnl endpoint (for UTA accounts)
func (b *BybitClient) getClosePnl(ctx context.Context) ([]model.Position, error) {
	params := map[string]interface{}{
		"category": "linear",
		"limit":    100,
	}

	result, err := b.bybit.NewClassicalBybitServiceWithParams(params).GetClosePnl(ctx)
	if err != nil {
		log.Printf("[bybit] GetClosePnl error: %v", err)
		return nil, err
	}

	log.Printf("[bybit] GetClosePnl RetCode: %d, RetMsg: %s", result.RetCode, result.RetMsg)

	list, ok := result.Result.(map[string]interface{})
	if !ok {
		return nil, nil
	}

	items, ok := list["list"].([]interface{})
	if !ok {
		return nil, nil
	}

	log.Printf("[bybit] GetClosePnl retrieved %d positions", len(items))

	if len(items) == 0 {
		// Try execution history as last resort
		return b.getExecutionHistory(ctx)
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

		volume, _ := strconv.ParseFloat(posMap["cumEntryValue"].(string), 64)
		leverage, _ := strconv.Atoi(posMap["leverage"].(string))
		margin := volume / float64(leverage)
		closedPnl, _ := strconv.ParseFloat(posMap["closedPnl"].(string), 64)

		updatedTime, _ := strconv.ParseFloat(posMap["updatedTime"].(string), 64)
		date := time.UnixMilli(int64(updatedTime))

		positions = append(positions, model.Position{
			OrderID:      orderID,
			Exchange:     "bybit",
			Symbol:       symbol,
			Volume:       volume,
			Margin:       margin,
			Leverage:     leverage,
			ClosedPnl:    closedPnl,
			Side:         side,
			UpdatedAt:    date,
		})
	}

	return positions, nil
}

// getExecutionHistory fetches closed positions from execution history via direct HTTP API
func (b *BybitClient) getExecutionHistory(ctx context.Context) ([]model.Position, error) {
	log.Printf("[bybit] Trying direct API call to execution history...")

	// Bybit V5 API: /v5/execution/list
	baseURL := "https://api.bybit.com"
	endpoint := "/v5/execution/list"
	
	var allExecutions []map[string]interface{}
	var cursor string
	
	// Paginate through all history
	page := 0
	for {
		page++
		params := url.Values{}
		params.Set("category", "linear")
		params.Set("limit", "1000")
		
		if cursor != "" {
			params.Set("cursor", cursor)
		}
		
		queryString := params.Encode()
		timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		
		// Generate signature
		signature := b.signV5("GET", endpoint, queryString, timestamp)
		
		reqURL := baseURL + endpoint + "?" + queryString
		req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
		if err != nil {
			return nil, err
		}
		
		req.Header.Set("X-BAPI-API-KEY", b.apiKey)
		req.Header.Set("X-BAPI-SIGN", signature)
		req.Header.Set("X-BAPI-TIMESTAMP", timestamp)
		req.Header.Set("X-BAPI-RECV-WINDOW", "30000")
		
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[bybit] HTTP request error: %v", err)
			return nil, err
		}
		
		var apiResp struct {
			RetCode int         `json:"retCode"`
			RetMsg  string      `json:"retMsg"`
			Result  interface{} `json:"result"`
		}
		
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()
		
		log.Printf("[bybit] Page %d - RetCode: %d, RetMsg: %s", page, apiResp.RetCode, apiResp.RetMsg)
		
		if apiResp.RetCode != 0 {
			log.Printf("[bybit] API error: %s", apiResp.RetMsg)
			break
		}
		
		result, ok := apiResp.Result.(map[string]interface{})
		if !ok {
			break
		}
		
		execList, ok := result["list"].([]interface{})
		if !ok {
			break
		}
		
		if len(execList) == 0 {
			log.Printf("[bybit] No more executions")
			break
		}
		
		log.Printf("[bybit] Retrieved %d executions on page %d", len(execList), page)
		
		for _, item := range execList {
			execMap, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			allExecutions = append(allExecutions, execMap)
		}
		
		// Get next page cursor
		nextCursor, _ := result["nextPageCursor"].(string)
		if nextCursor == "" {
			break
		}
		cursor = nextCursor
		
		// Rate limiting - wait between requests
		time.Sleep(200 * time.Millisecond)
	}
	
	log.Printf("[bybit] Total executions retrieved: %d", len(allExecutions))
	
	// Group by orderId and find closing executions
	orderMap := make(map[string]map[string]interface{})
	
	for _, execMap := range allExecutions {
		orderID, _ := execMap["orderId"].(string)
		if orderID == "" {
			continue
		}
		
		// Look for executions with closedSize (closing executions)
		closedSize, _ := execMap["closedSize"].(string)
		if closedSize == "0" || closedSize == "" {
			continue
		}
		
		orderMap[orderID] = execMap
	}
	
	log.Printf("[bybit] Found %d closed positions", len(orderMap))
	
	var positions []model.Position
	for _, execMap := range orderMap {
		orderID, _ := execMap["orderId"].(string)
		symbol, _ := execMap["symbol"].(string)
		side, _ := execMap["side"].(string)
		
		closedSize, _ := strconv.ParseFloat(execMap["closedSize"].(string), 64)
		execPrice, _ := strconv.ParseFloat(execMap["execPrice"].(string), 64)
		volume := closedSize * execPrice
		
		leverageStr, _ := execMap["leverage"].(string)
		leverage, _ := strconv.Atoi(leverageStr)
		if leverage == 0 {
			leverage = 1
		}
		margin := volume / float64(leverage)
		
		// closedPnl might be in execFee or closedPnl field
		closedPnl, _ := strconv.ParseFloat(execMap["closedPnl"].(string), 64)
		
		execTime, _ := strconv.ParseFloat(execMap["execTime"].(string), 64)
		date := time.UnixMilli(int64(execTime))
		
		positions = append(positions, model.Position{
			OrderID:      orderID,
			Exchange:     "bybit",
			Symbol:       symbol,
			Volume:       volume,
			Margin:       margin,
			Leverage:     leverage,
			ClosedPnl:    closedPnl,
			Side:         side,
			UpdatedAt:    date,
		})
	}
	
	return positions, nil
}

// signV5 generates signature for Bybit V5 API
func (b *BybitClient) signV5(method, path, queryString, timestamp string) string {
	// V5 signature: timestamp + apiKey + recvWindow + payload
	// For GET: payload = queryString
	// For POST: payload = body
	recvWindow := "30000"
	payload := timestamp + b.apiKey + recvWindow
	
	if method == "GET" && queryString != "" {
		payload += queryString
	}
	
	h := hmac.New(sha256.New, []byte(b.apiSecret))
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}
