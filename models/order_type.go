package models

import (
	"time"
)

type OrderType struct {
	ID                      string    `json:"id"`
	CompanyID               string    `json:"company_id"`
	Name                    string    `json:"name"`
	IsActivePriceAdjustment bool      `json:"is_active_price_adjustment"`
	PriceIncrease           float64   `json:"price_increase"`
	PriceDecrease           float64   `json:"price_decrease"`
	IncreaseType            string    `json:"increase_type"`
	DecreaseType            string    `json:"decrease_type"`
	IncreaseValue           float64   `json:"increase_value"`
	DecreaseValue           float64   `json:"decrease_value"`
	IsActive                bool      `json:"is_active"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type OrderTypeInput struct {
	Name                    string  `json:"name"`
	IsActivePriceAdjustment bool    `json:"is_active_price_adjustment"`
	PriceIncrease           float64 `json:"price_increase"`
	PriceDecrease           float64 `json:"price_decrease"`
	IncreaseType            string  `json:"increase_type"`
	DecreaseType            string  `json:"decrease_type"`
	IncreaseValue           float64 `json:"increase_value"`
	DecreaseValue           float64 `json:"decrease_value"`
	IsActive                bool    `json:"is_active"`
}
