package model

import "time"

type Position struct {
	ID           int       `json:"id"`
	OrderID      string    `json:"orderId"`
	Exchange     string    `json:"exchange"`
	Symbol       string    `json:"symbol"`
	CumExitValue float64   `json:"cumExitValue"`
	Quantity     float64   `json:"qty"`
	Leverage     int       `json:"leverage"`
	ClosedPnl    float64   `json:"closedPnl"`
	Side         string    `json:"side"`
	UpdatedAt    time.Time `json:"date"`
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
