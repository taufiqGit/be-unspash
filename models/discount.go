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
	DiscountTargetAll      DiscountTarget = "all"      // Semua Barang
	DiscountTargetSpecific DiscountTarget = "specific" // Spesifik
)

// DiscountSpecificTarget merepresentasikan sub-target ketika target = "specific"
type DiscountSpecificTarget string

const (
	DiscountSpecificTargetCategory DiscountSpecificTarget = "category" // Kategori
	DiscountSpecificTargetProduct  DiscountSpecificTarget = "product"  // Barang Tertentu
)

// Discount adalah model utama yang merepresentasikan data di database
type Discount struct {
	ID        string       `json:"id"`
	CompanyID string       `json:"company_id"`
	OutletID  string       `json:"outlet_id"`
	Name      string       `json:"name"`
	Type      DiscountType `json:"type"`

	// Nilai diskon: nominal Rp atau persentase (%)
	DiscountValue float64 `json:"discount_value"`

	// Maksimal Nominal (Rp) — hanya untuk: Diskon Barang (%) & Diskon Struk (%)
	MaxAmount *float64 `json:"max_amount,omitempty"`

	// Minimum Pembelanjaan (Rp) — hanya untuk: Diskon Struk (Rp) & Diskon Struk (%)
	MinPurchase *float64 `json:"min_purchase,omitempty"`

	// Target Diskon — hanya untuk: Diskon Barang (Rp) & Diskon Barang (%)
	// Nilai: "all" (Semua Barang) | "specific" (Spesifik)
	TargetType *DiscountTarget `json:"target_type,omitempty"`

	// Target Spesifik — hanya aktif ketika TargetType = "specific"
	// Nilai: "category" (Kategori) | "product" (Barang Tertentu)
	SpecificTargetType *DiscountSpecificTarget `json:"specific_target_type,omitempty"`

	// ID kategori yang menjadi target — aktif ketika SpecificTargetType = "category"
	TargetCategoryIDs []string `json:"target_category_ids,omitempty"`

	// ID produk yang menjadi target — aktif ketika SpecificTargetType = "product"
	TargetProductIDs []string `json:"target_product_ids,omitempty"`

	// Target Diskon Lainnya (opsional) — berlaku untuk semua tipe
	ApplyToOrderTypes bool     `json:"apply_to_order_types"`
	OrderTypeIDs      []string `json:"order_type_ids,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DiscountInput digunakan untuk menerima payload dari request (Create & Update)
type DiscountInput struct {
	Name          string       `json:"name"`
	Type          DiscountType `json:"type"`
	OutletID      string       `json:"outlet_id"`
	DiscountValue float64      `json:"discount_value"`

	// Opsional tergantung tipe
	MaxAmount   *float64 `json:"max_amount,omitempty"`
	MinPurchase *float64 `json:"min_purchase,omitempty"`

	// Target Diskon — untuk Diskon Barang saja
	TargetType         *DiscountTarget         `json:"target_type,omitempty"`
	SpecificTargetType *DiscountSpecificTarget  `json:"specific_target_type,omitempty"`
	TargetCategoryIDs  []string                `json:"target_category_ids,omitempty"`
	TargetProductIDs   []string                `json:"target_product_ids,omitempty"`

	// Target Diskon Lainnya — opsional untuk semua tipe
	ApplyToOrderTypes bool     `json:"apply_to_order_types"`
	OrderTypeIDs      []string `json:"order_type_ids,omitempty"`
}
