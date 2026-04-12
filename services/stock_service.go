package services

import (
	"errors"
	"gowes/models"
	"gowes/repositories"
	"strings"
)

var ErrStockOutletRequired = errors.New("outlet_id is required")
var ErrStockProductRequired = errors.New("product_id is required")

type StockService interface {
	ListStocks(companyID string, params models.PaginationParams, outletID string, productID string) ([]models.StockPerOutlet, int, error)
	GetStock(companyID string, outletID string, productID string) (models.StockPerOutlet, error)
}

type stockService struct {
	repo repositories.StockRepository
}

func NewStockService(repo repositories.StockRepository) StockService {
	return &stockService{repo: repo}
}

func (s *stockService) ListStocks(companyID string, params models.PaginationParams, outletID string, productID string) ([]models.StockPerOutlet, int, error) {
	return s.repo.FindAll(companyID, params, outletID, productID)
}

func (s *stockService) GetStock(companyID string, outletID string, productID string) (models.StockPerOutlet, error) {
	if strings.TrimSpace(outletID) == "" {
		return models.StockPerOutlet{}, ErrStockOutletRequired
	}
	if strings.TrimSpace(productID) == "" {
		return models.StockPerOutlet{}, ErrStockProductRequired
	}
	return s.repo.FindByOutletAndProduct(companyID, strings.TrimSpace(outletID), strings.TrimSpace(productID))
}
