package services

import (
	"errors"
	"gowes/models"
	"gowes/repositories"
	"strings"
	"time"
)

// ─── Sentinel errors ──────────────────────────────────────────────────────────

var (
	ErrDiscountNameRequired               = errors.New("discount name is required")
	ErrDiscountValueInvalid               = errors.New("discount_value must be greater than 0")
	ErrDiscountTypeInvalid                = errors.New("invalid discount type, must be one of: product_rp, product_pct, receipt_rp, receipt_pct")
	ErrDiscountPctExceeded                = errors.New("discount_value cannot exceed 100 for percentage type")
	ErrDiscountTargetNotApplicable        = errors.New("target_type is only applicable for product discount types")
	ErrDiscountSpecificTargetNotApplicable = errors.New("specific_target_type requires target_type to be 'specific'")
)

// ─── Interface ────────────────────────────────────────────────────────────────

type DiscountService interface {
	ListDiscounts(companyID string, params models.PaginationParams) ([]models.Discount, int, error)
	GetDiscount(id string) (models.Discount, error)
	CreateDiscount(companyID string, in models.DiscountInput) (models.Discount, error)
	UpdateDiscount(id string, in models.DiscountInput) (models.Discount, error)
	DeleteDiscount(id string) error
}

// ─── Implementation ───────────────────────────────────────────────────────────

type discountService struct {
	repo repositories.DiscountRepository
}

func NewDiscountService(repo repositories.DiscountRepository) DiscountService {
	return &discountService{repo: repo}
}

func (s *discountService) ListDiscounts(companyID string, params models.PaginationParams) ([]models.Discount, int, error) {
	return s.repo.FindAll(companyID, params)
}

func (s *discountService) GetDiscount(id string) (models.Discount, error) {
	return s.repo.FindByID(id)
}

func (s *discountService) CreateDiscount(companyID string, in models.DiscountInput) (models.Discount, error) {
	if err := validateDiscountInput(in); err != nil {
		return models.Discount{}, err
	}

	now := time.Now().UTC()
	d := models.Discount{
		CompanyID:          companyID,
		Name:               strings.TrimSpace(in.Name),
		Type:               in.Type,
		DiscountValue:      in.DiscountValue,
		MaxAmount:          in.MaxAmount,
		MinPurchase:        in.MinPurchase,
		TargetType:         in.TargetType,
		SpecificTargetType: in.SpecificTargetType,
		ApplyToOrderTypes:  in.ApplyToOrderTypes,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	outletIDs, categoryIDs, productIDs, orderTypeIDs := resolveRelationIDs(in)

	return s.repo.Create(d, outletIDs, categoryIDs, productIDs, orderTypeIDs)
}

func (s *discountService) UpdateDiscount(id string, in models.DiscountInput) (models.Discount, error) {
	if err := validateDiscountInput(in); err != nil {
		return models.Discount{}, err
	}

	existing, err := s.repo.FindByID(id)
	if err != nil {
		return models.Discount{}, err
	}

	existing.Name               = strings.TrimSpace(in.Name)
	existing.Type               = in.Type
	existing.DiscountValue      = in.DiscountValue
	existing.MaxAmount          = in.MaxAmount
	existing.MinPurchase        = in.MinPurchase
	existing.TargetType         = in.TargetType
	existing.SpecificTargetType = in.SpecificTargetType
	existing.ApplyToOrderTypes  = in.ApplyToOrderTypes
	existing.UpdatedAt          = time.Now().UTC()

	outletIDs, categoryIDs, productIDs, orderTypeIDs := resolveRelationIDs(in)

	return s.repo.Update(existing, outletIDs, categoryIDs, productIDs, orderTypeIDs)
}

func (s *discountService) DeleteDiscount(id string) error {
	return s.repo.Delete(id)
}

// ─── helpers ──────────────────────────────────────────────────────────────────

// validateDiscountInput memvalidasi field wajib dan konsistensi antar field.
func validateDiscountInput(in models.DiscountInput) error {
	if strings.TrimSpace(in.Name) == "" {
		return ErrDiscountNameRequired
	}
	if in.DiscountValue <= 0 {
		return ErrDiscountValueInvalid
	}

	switch in.Type {
	case models.DiscountTypeProductRp, models.DiscountTypeProductPct,
		models.DiscountTypeReceiptRp, models.DiscountTypeReceiptPct:
		// valid
	default:
		return ErrDiscountTypeInvalid
	}

	// Tipe persen: nilai diskon tidak boleh lebih dari 100
	if (in.Type == models.DiscountTypeProductPct || in.Type == models.DiscountTypeReceiptPct) &&
		in.DiscountValue > 100 {
		return ErrDiscountPctExceeded
	}

	// target_type hanya relevan untuk Diskon Barang
	if in.TargetType != nil {
		if in.Type != models.DiscountTypeProductRp && in.Type != models.DiscountTypeProductPct {
			return ErrDiscountTargetNotApplicable
		}
	}

	// specific_target_type hanya relevan saat target_type = "specific"
	if in.SpecificTargetType != nil {
		if in.TargetType == nil || *in.TargetType != models.DiscountTargetSpecific {
			return ErrDiscountSpecificTargetNotApplicable
		}
	}

	return nil
}

// resolveRelationIDs membersihkan dan memfilter IDs relasi berdasarkan tipe dan
// target diskon, lalu mengembalikan 4 slice siap pakai.
func resolveRelationIDs(in models.DiscountInput) (outletIDs, categoryIDs, productIDs, orderTypeIDs []string) {
	// Outlet — berlaku untuk semua tipe
	outletIDs = deduplicateIDs(in.OutletIDs)

	// Category / Product — hanya untuk Diskon Barang dengan target = "specific"
	isProductType := in.Type == models.DiscountTypeProductRp || in.Type == models.DiscountTypeProductPct
	if isProductType && in.TargetType != nil && *in.TargetType == models.DiscountTargetSpecific &&
		in.SpecificTargetType != nil {
		switch *in.SpecificTargetType {
		case models.DiscountSpecificTargetCategory:
			categoryIDs = deduplicateIDs(in.TargetCategoryIDs)
		case models.DiscountSpecificTargetProduct:
			productIDs = deduplicateIDs(in.TargetProductIDs)
		}
	}

	// Order types — hanya aktif jika toggle apply_to_order_types = true
	if in.ApplyToOrderTypes {
		orderTypeIDs = deduplicateIDs(in.OrderTypeIDs)
	}

	return
}

// deduplicateIDs menghapus duplikat dan string kosong dari slice ID.
func deduplicateIDs(ids []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id != "" && !seen[id] {
			seen[id] = true
			result = append(result, id)
		}
	}
	return result
}
