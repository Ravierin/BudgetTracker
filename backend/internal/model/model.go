package model

import "time"

type Position struct {
	ID           int       `json:"id"`
	OrderID      string    `json:"orderId"`
	Exchange     string    `json:"exchange"`
	Symbol       string    `json:"symbol"`
	Volume       float64   `json:"volume"`
	Leverage     int       `json:"leverage"`
	ClosedPnl    float64   `json:"closedPnl"`
	Side         string    `json:"side"`
	UpdatedAt    time.Time `json:"date"`
}

type ExchangeBalance struct {
	Exchange string  `json:"exchange"`
	Balance  float64 `json:"balance"`
}

type Withdrawal struct {
	ID        int       `json:"id"`
	Exchange  string    `json:"exchange"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"date"`
}

type MonthlyIncome struct {
	ID        int       `json:"id"`
	Exchange  string    `json:"exchange"`
	Amount    float64   `json:"amount"`
	PNL       float64   `json:"pnl"`
	CreatedAt time.Time `json:"date"`
}

type APIKey struct {
	ID        int       `json:"id"`
	Exchange  string    `json:"exchange"`
	APIKey    string    `json:"apiKey"`
	APISecret string    `json:"apiSecret"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
