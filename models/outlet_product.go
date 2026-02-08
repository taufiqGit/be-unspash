package models

import (
	"time"
)

type OutletProduct struct {
	ID        string    `json:"id"`
	CompanyID string    `json:"company_id"`
	OutletID  string    `json:"outlet_id"`
	ProductID string    `json:"product_id"`
	Stock     int       `json:"stock"`
	Price     float64   `json:"price"`
	Cost      float64   `json:"cost"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
