package models

import "time"

type Unit struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Symbol    string    `json:"symbol"`
	Type      string    `json:"type"`
	CompanyID string    `json:"company_id"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
}

type UnitInput struct {
	Name      string `json:"name"`
	Symbol    string `json:"symbol"`
	Type      string `json:"type"` // weight / volume / qty
	CompanyID string `json:"company_id"`
}
