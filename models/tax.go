package models

import (
	"time"
)

type Tax struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Rate        float64   `json:"rate"`
	CompanyID   string    `json:"company_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaxInput struct {
	Name string  `json:"name"`
	Rate float64 `json:"rate"`
}
