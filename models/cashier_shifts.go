package models

import "time"

type CashierShift struct {
	ID           string    `json:"id"`
	CompanyID    string    `json:"company_id"`
	OutletID     string    `json:"outlet_id"`
	UserID       string    `json:"user_id"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Status       string    `json:"status"`
	OpeningCash  float64   `json:"opening_cash"`
	ClosingCash  float64   `json:"closing_cash"`
	ExpectedCash float64   `json:"expected_cash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type StartCashierShiftInput struct {
	UserID       string    `json:"user_id"`
	OutletID     string    `json:"outlet_id"`
	StartTime    time.Time `json:"start_time"`
	Status       string    `json:"status"`
	OpeningCash  float64   `json:"opening_cash"`
	ClosingCash  float64   `json:"closing_cash"`
	ExpectedCash float64   `json:"expected_cash"`
}

type EndCashierShiftInput struct {
	EndTime      time.Time `json:"end_time"`
	Status       string    `json:"status"`
	ClosingCash  float64   `json:"closing_cash"`
	ExpectedCash float64   `json:"expected_cash"`
}
