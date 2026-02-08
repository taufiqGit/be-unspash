package models

import (
	"time"
)

type AddOnProduct struct {
	ID        string    `json:"id"`
	CompanyID string    `json:"company_id"`
	AddOnID   string    `json:"add_on_id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
