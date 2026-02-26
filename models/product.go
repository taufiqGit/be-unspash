package models

import (
	"time"
)

type Product struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	SKU        string    `json:"sku"`
	Unit       string    `json:"unit"`
	Cost       float64   `json:"cost"`
	Price      float64   `json:"price"`
	ImageURL   string    `json:"image_url"`
	CompanyID  string    `json:"company_id"`
	CategoryID string    `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ProductList struct {
	Name     string  `json:"name"`
	SKU      string  `json:"sku"`
	Unit     string  `json:"unit"`
	Cost     float64 `json:"cost"`
	Price    float64 `json:"price"`
	ImageURL string  `json:"image_url"`
	Category string  `json:"category"`
}
type ProductInput struct {
	Name       string  `json:"name"`
	SKU        string  `json:"sku"`
	Unit       string  `json:"unit"`
	Cost       float64 `json:"cost"`
	Price      float64 `json:"price"`
	ImageURL   string  `json:"image_url"`
	CompanyID  string  `json:"company_id"`
	CategoryID string  `json:"category_id"`
}
