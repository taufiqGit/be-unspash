package services

import (
	"errors"
	"gowes/models"
	"gowes/repositories"
	"strings"
	"time"
)

var ErrPurchaseOutletRequired = errors.New("outlet_id is required")
var ErrPurchaseDetailsRequired = errors.New("purchase details are required")
var ErrPurchaseDetailProductRequired = errors.New("product_id is required for every purchase detail")
var ErrPurchaseDetailQtyInvalid = errors.New("quantity must be greater than zero")
var ErrPurchaseDetailPriceInvalid = errors.New("price cannot be negative")

type PurchaseService interface {
	ListPurchases(companyID string, params models.PaginationParams) ([]models.Purchase, int, error)
	GetPurchase(id string, companyID string) (models.Purchase, error)
	CreatePurchase(companyID string, userID string, input models.PurchaseInput) (models.Purchase, error)
}

type purchaseService struct {
	repo repositories.PurchaseRepository
}

func NewPurchaseService(repo repositories.PurchaseRepository) PurchaseService {
	return &purchaseService{repo: repo}
}

func (s *purchaseService) ListPurchases(companyID string, params models.PaginationParams) ([]models.Purchase, int, error) {
	return s.repo.FindAll(companyID, params)
}

func (s *purchaseService) GetPurchase(id string, companyID string) (models.Purchase, error) {
	return s.repo.FindByID(id, companyID)
}

func (s *purchaseService) CreatePurchase(companyID string, userID string, input models.PurchaseInput) (models.Purchase, error) {
	if strings.TrimSpace(input.OutletID) == "" {
		return models.Purchase{}, ErrPurchaseOutletRequired
	}
	if len(input.Details) == 0 {
		return models.Purchase{}, ErrPurchaseDetailsRequired
	}

	subtotal := 0.0
	details := make([]models.PurchaseDetail, 0, len(input.Details))
	for _, d := range input.Details {
		if strings.TrimSpace(d.ProductID) == "" {
			return models.Purchase{}, ErrPurchaseDetailProductRequired
		}
		if d.Quantity <= 0 {
			return models.Purchase{}, ErrPurchaseDetailQtyInvalid
		}
		if d.Price < 0 {
			return models.Purchase{}, ErrPurchaseDetailPriceInvalid
		}

		total := float64(d.Quantity) * d.Price
		subtotal += total
		details = append(details, models.PurchaseDetail{
			ProductID: strings.TrimSpace(d.ProductID),
			Quantity:  d.Quantity,
			Price:     d.Price,
			Total:     total,
		})
	}

	grandTotal := subtotal + input.TaxValue - input.DiscountBill
	if grandTotal < 0 {
		grandTotal = 0
	}

	changeAmount := 0.0
	if input.PaidAmount > grandTotal {
		changeAmount = input.PaidAmount - grandTotal
	}

	paymentMethod := strings.TrimSpace(input.PaymentMethod)
	if paymentMethod == "" {
		paymentMethod = "cash"
	}
	status := strings.TrimSpace(input.Status)
	if status == "" {
		status = "completed"
	}

	now := time.Now().UTC()
	purchase := models.Purchase{
		CompanyID:     companyID,
		UserID:        userID,
		OutletID:      strings.TrimSpace(input.OutletID),
		PaymentMethod: paymentMethod,
		GrandTotal:    grandTotal,
		TaxValue:      input.TaxValue,
		PaidAmount:    input.PaidAmount,
		ChangeAmount:  changeAmount,
		Status:        status,
		DiscountBill:  input.DiscountBill,
		Details:       details,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	return s.repo.CreateWithStockMovement(purchase)
}
