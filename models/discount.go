package models

import "time"

// DiscountType merepresentasikan 4 jenis diskon
type DiscountType string

const (
	DiscountTypeProductRp  DiscountType = "product_rp"  // Diskon Barang (Rp)
	DiscountTypeProductPct DiscountType = "product_pct" // Diskon Barang (%)
	DiscountTypeReceiptRp  DiscountType = "receipt_rp"  // Diskon Struk (Rp)
	DiscountTypeReceiptPct DiscountType = "receipt_pct" // Diskon Struk (%)
)

// DiscountTarget merepresentasikan target diskon untuk tipe Diskon Barang
type DiscountTarget string

const (
	DiscountTargetAll          DiscountTarget = "all"      // Semua Barang
	DiscountTargetTypeProduct  DiscountTarget = "product"  // Barang Tertentu
	DiscountTargetTypeCategory DiscountTarget = "category" // Kategori Tertentu
)

type DiscountTargetCategory struct {
	ID         string `json:"id"`
	CategoryId string `json:"category_id"`
	ParentID   string `json:"parent_id"`
}

type DiscountTargetProduct struct {
	ID        string `json:"id"`
	ProductId string `json:"product_id"`
	ParentID  string `json:"parent_id"`
}

type DiscountTargetOutlet struct {
	ID       string `json:"id"`
	OutletID string `json:"outlet_id"`
	ParentID string `json:"parent_id"`
}

type DiscountTargetOrderType struct {
	ID          string `json:"id"`
	OrderTypeID string `json:"order_type_id"`
	ParentID    string `json:"parent_id"`
}

// DiscountSpecificTarget merepresentasikan sub-target ketika target = "specific"
type DiscountSpecificTarget string

const (
	DiscountSpecificTargetCategory DiscountSpecificTarget = "category" // Kategori
	DiscountSpecificTargetProduct  DiscountSpecificTarget = "product"  // Barang Tertentu
)

// Discount adalah model utama yang merepresentasikan data di database
type Discount struct {
	ID        string                 `json:"id"`
	CompanyID string                 `json:"company_id"`
	Name      string                 `json:"name"`
	Type      DiscountType           `json:"type"`
	OutletIDs []DiscountTargetOutlet `json:"outlet_ids"`

	// Nilai diskon: nominal Rp atau persentase (%)
	DiscountValue float64 `json:"discount_value"`

	// Maksimal Nominal (Rp) — hanya untuk: Diskon Barang (%) & Diskon Struk (%)
	MaxAmount *float64 `json:"max_amount,omitempty"`

	// Minimum Pembelanjaan (Rp) — hanya untuk: Diskon Struk (Rp) & Diskon Struk (%)
	MinPurchase *float64 `json:"min_purchase,omitempty"`

	// Target Diskon — hanya untuk: Diskon Barang (Rp) & Diskon Barang (%)
	// Nilai: "all" (Semua Barang) | "specific" (Spesifik)
	TargetType *DiscountTarget `json:"target_type,omitempty"`

	// Prioritas diskon (semakin kecil semakin diprioritaskan)
	Priority int `json:"priority"`

	// ID kategori yang menjadi target
	TargetCategoryIDs []DiscountTargetCategory `json:"target_category_ids,omitempty"`

	// ID produk yang menjadi target
	TargetProductIDs []DiscountTargetProduct `json:"target_product_ids,omitempty"`

	// Target Diskon Lainnya (opsional) — berlaku untuk semua tipe
	ApplyToOrderTypes bool                      `json:"apply_to_order_types"`
	OrderTypeIDs      []DiscountTargetOrderType `json:"order_type_ids,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DiscountInput digunakan untuk menerima payload dari request (Create & Update)
type DiscountInput struct {
	Name          string       `json:"name"`
	Type          DiscountType `json:"type"`
	CompanyID     string       `json:"company_id"`
	DiscountValue float64      `json:"discount_value"`
	OutletIDs     []string     `json:"outlet_ids,omitempty"`
	// Opsional tergantung tipe
	MaxAmount   *float64 `json:"max_amount,omitempty"`
	MinPurchase *float64 `json:"min_purchase,omitempty"`

	// Target Diskon — untuk Diskon Barang saja
	TargetType        *DiscountTarget `json:"target_type,omitempty"`
	Priority          int             `json:"priority"`
	TargetCategoryIDs []string        `json:"target_category_ids,omitempty"`
	TargetProductIDs  []string        `json:"target_product_ids,omitempty"`

	// Target Diskon Lainnya — opsional untuk semua tipe
	ApplyToOrderTypes bool     `json:"apply_to_order_types"`
	OrderTypeIDs      []string `json:"order_type_ids,omitempty"`
}
