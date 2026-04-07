package models

import "time"

type Recipe struct {
	ID           string    `json:"id"`
	CompanyID    string    `json:"company_id"`
	ProductID    string    `json:"product_id"`
	IngredientID string    `json:"ingredient_id"`
	IsActive     bool      `json:"is_active"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
}

type RecipeInput struct {
	CompanyID    string `json:"company_id"`
	ProductID    string `json:"product_id"`
	IngredientID string `json:"ingredient_id"`
	IsActive     *bool  `json:"is_active"`
}
