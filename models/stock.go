package models

import "time"

type Stock struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	OutletID  string `json:"outlet_id"`
	Qty       int    `json:"qty"`
}

type StockPerOutlet struct {
	ID          string `json:"id"`
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	ProductSKU  string `json:"product_sku"`
	OutletID    string `json:"outlet_id"`
	OutletName  string `json:"outlet_name"`
	Qty         int    `json:"qty"`
}

type StockMovementType string

const (
	StockMovementTypeIn         StockMovementType = "IN"
	StockMovementTypeOut        StockMovementType = "OUT"
	StockMovementTypeAdjustment StockMovementType = "ADJUSTMENT"
	StockMovementTypeTransfer   StockMovementType = "TRANSFER"
)

type StockReferenceType string

const (
	StockReferenceTypePurchase   StockReferenceType = "purchase"
	StockReferenceTypeSale       StockReferenceType = "sale"
	StockReferenceTypeAdjustment StockReferenceType = "adjustment"
	StockReferenceTypeTransfer   StockReferenceType = "transfer"
)

type StockMovement struct {
	ID            string             `json:"id"`
	ProductID     string             `json:"product_id"`
	ProductName   string             `json:"product_name,omitempty"`
	ProductSKU    string             `json:"product_sku,omitempty"`
	OutletID      string             `json:"outlet_id"`
	OutletName    string             `json:"outlet_name,omitempty"`
	Type          StockMovementType  `json:"type"`
	Qty           int                `json:"qty"`
	ReferenceType StockReferenceType `json:"reference_type"`
	ReferenceID   string             `json:"reference_id"`
	Note          string             `json:"note"`
	CreatedAt     time.Time          `json:"created_at"`
}
