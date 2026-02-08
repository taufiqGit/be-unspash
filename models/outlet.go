package models

import (
	"time"
)

type Outlet struct {
	ID         string    `json:"id"`
	CompanyID  string    `json:"company_id"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	Supervisor string    `json:"supervisor"`
	Address    string    `json:"address"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
