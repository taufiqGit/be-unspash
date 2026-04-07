package models

import "time"

type Role struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CompanyID string    `json:"company_id"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RoleInput struct {
	Name      string `json:"name"`
	CompanyID string `json:"company_id"`
}
