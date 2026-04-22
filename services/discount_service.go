package services

import (
	"errors"
	"fmt"
	"gowes/models"
	"gowes/repositories"
	"strings"
	"time"
)

// ─── Sentinel errors ──────────────────────────────────────────────────────────

var (
	ErrDiscountNameRequired        = errors.New("discount name is required")
	ErrDiscountValueInvalid        = errors.New("discount_value must be greater than 0")
	ErrDiscountTypeInvalid         = errors.New("invalid discount type, must be one of: product_rp, product_pct, receipt_rp, receipt_pct")
	ErrDiscountPctExceeded         = errors.New("discount_value cannot exceed 100 for percentage type")
	ErrDiscountTargetNotApplicable = errors.New("target_type is only applicable for product discount types")
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
		fmt.Println(err, "err validate discount input")
		return models.Discount{}, err
	}

	now := time.Now().UTC()
	d := models.Discount{
		CompanyID:         companyID,
		Name:              strings.TrimSpace(in.Name),
		Type:              in.Type,
		DiscountValue:     in.DiscountValue,
		MaxAmount:         in.MaxAmount,
		MinPurchase:       in.MinPurchase,
		TargetType:        in.TargetType,
		Priority:          in.Priority,
		ApplyToOrderTypes: in.ApplyToOrderTypes,
		CreatedAt:         now,
		UpdatedAt:         now,
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

	existing.Name = strings.TrimSpace(in.Name)
	existing.Type = in.Type
	existing.DiscountValue = in.DiscountValue
	existing.MaxAmount = in.MaxAmount
	existing.MinPurchase = in.MinPurchase
	existing.TargetType = in.TargetType
	existing.Priority = in.Priority
	existing.ApplyToOrderTypes = in.ApplyToOrderTypes
	existing.UpdatedAt = time.Now().UTC()

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
	// if in.TargetType != nil {
	// 	fmt.Println(in.Type != models.DiscountTypeProductRp && in.Type != models.DiscountTypeProductPct, "target type")
	// 	if in.Type != models.DiscountTypeProductRp && in.Type != models.DiscountTypeProductPct {
	// 		return ErrDiscountTargetNotApplicable
	// 	}
	// }

	return nil
}

// resolveRelationIDs membersihkan dan memfilter IDs relasi berdasarkan tipe dan
// target diskon, lalu mengembalikan 4 slice siap pakai.
func resolveRelationIDs(in models.DiscountInput) (outletIDs, categoryIDs, productIDs, orderTypeIDs []string) {
	// Outlet — berlaku untuk semua tipe
	outletIDs = deduplicateIDs(in.OutletIDs)

	// Category / Product — hanya relevan untuk Diskon Barang dengan target = "category" atau "product"
	isProductType := in.Type == models.DiscountTypeProductRp || in.Type == models.DiscountTypeProductPct
	if isProductType && in.TargetType != nil && (*in.TargetType == models.DiscountTargetTypeCategory || *in.TargetType == models.DiscountTargetTypeProduct) {
		categoryIDs = deduplicateIDs(in.TargetCategoryIDs)
		productIDs = deduplicateIDs(in.TargetProductIDs)
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
