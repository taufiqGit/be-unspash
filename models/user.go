package models

import (
	"time"
)

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleCashier UserRole = "cashier"
	RoleWaiter  UserRole = "waiter"
)

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never return password hash
	PosPIN       *string   `json:"pos_pin,omitempty"`
	Phone        *string   `json:"phone,omitempty"`
	CompanyID    *string   `json:"company_id,omitempty"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserInput struct {
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Role     UserRole `json:"role"`
	PosPIN   string   `json:"pos_pin"`
}

type UserRegisterInput struct {
	Username      string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Phone         string `json:"phone"`
	BussinessName string `json:"bussiness_name"`
	PosPIN        string `json:"pos_pin"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
