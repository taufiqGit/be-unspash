package models

import (
	"time"
)

type Customer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Phone     *string   `json:"phone,omitempty"`
	Email     *string   `json:"email,omitempty"`
	CompanyID string    `json:"company_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CustomerInput struct {
	Name      string  `json:"name"`
	Phone     *string `json:"phone,omitempty"`
	Email     *string `json:"email,omitempty"`
	CompanyID string  `json:"company_id"`
}
