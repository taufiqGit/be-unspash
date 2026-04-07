package models

import "time"

type Supplier struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CompanyID   string    `json:"company_id"`
	Address     string    `json:"address"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	CompanyName string    `json:"company_name"`
	TaxNumber   string    `json:"tax_number"`
	IsActive    bool      `json:"is_active"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
}

type SupplierInput struct {
	CompanyID   string `json:"company_id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	CompanyName string `json:"company_name"`
	TaxNumber   string `json:"tax_number"`
	IsActive    *bool  `json:"is_active"`
}
