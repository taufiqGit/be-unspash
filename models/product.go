package models

import (
	"time"
)

type ProductType string

const (
	ProductTypeRawMaterial   ProductType = "raw_material"
	ProductTypeFinishedGoods ProductType = "finished_goods"
)

type Product struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	SKU        string      `json:"sku"`
	Unit       string      `json:"unit"`
	UnitID     string      `json:"unit_id"`
	Cost       float64     `json:"cost"`
	Price      float64     `json:"price"`
	ImageURL   string      `json:"image_url"`
	CompanyID  string      `json:"company_id"`
	CategoryID string      `json:"category_id"`
	Type       ProductType `json:"type"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

// raw_material
// finished_goods
type ProductList struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	SKU      string  `json:"sku"`
	Unit     string  `json:"unit"`
	UnitID   string  `json:"unit_id"`
	Cost     float64 `json:"cost"`
	Price    float64 `json:"price"`
	ImageURL string  `json:"image_url"`
	Category string  `json:"category"`
}
type ProductInput struct {
	Name       string  `json:"name"`
	SKU        string  `json:"sku"`
	Unit       string  `json:"unit"`
	UnitID     string  `json:"unit_id"`
	Cost       float64 `json:"cost"`
	Price      float64 `json:"price"`
	ImageURL   string  `json:"image_url"`
	CompanyID  string  `json:"company_id"`
	CategoryID string  `json:"category_id"`
}
