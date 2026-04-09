package models

import "time"

type Purchase struct {
	ID            string           `json:"id"`
	CompanyID     string           `json:"company_id"`
	UserID        string           `json:"user_id"`
	OutletID      string           `json:"outlet_id"`
	PaymentMethod string           `json:"payment_method"`
	GrandTotal    float64          `json:"grand_total"`
	TaxValue      float64          `json:"tax_value"`
	PaidAmount    float64          `json:"paid_amount"`
	ChangeAmount  float64          `json:"change_amount"`
	Status        string           `json:"status"`
	DiscountBill  float64          `json:"discount_bill"`
	Details       []PurchaseDetail `json:"details,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

type PurchaseDetail struct {
	ID         string  `json:"id"`
	PurchaseID string  `json:"purchase_id"`
	ProductID  string  `json:"product_id"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	Total      float64 `json:"total"`
}

type PurchaseDetailInput struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type PurchaseInput struct {
	OutletID      string                `json:"outlet_id"`
	PaymentMethod string                `json:"payment_method"`
	TaxValue      float64               `json:"tax_value"`
	PaidAmount    float64               `json:"paid_amount"`
	Status        string                `json:"status"`
	DiscountBill  float64               `json:"discount_bill"`
	Details       []PurchaseDetailInput `json:"details"`
}
