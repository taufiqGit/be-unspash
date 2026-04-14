package services

import (
	"gowes/models"
	"gowes/repositories"
	"strings"
)

type StockMovementService interface {
	ListStockMovements(companyID string, params models.PaginationParams, outletID string, productID string, movementType string, referenceType string) ([]models.StockMovement, int, error)
	GetStockMovement(companyID string, id string) (models.StockMovement, error)
}

type stockMovementService struct {
	repo repositories.StockMovementRepository
}

func NewStockMovementService(repo repositories.StockMovementRepository) StockMovementService {
	return &stockMovementService{repo: repo}
}

func (s *stockMovementService) ListStockMovements(companyID string, params models.PaginationParams, outletID string, productID string, movementType string, referenceType string) ([]models.StockMovement, int, error) {
	return s.repo.FindAll(companyID, params, strings.TrimSpace(outletID), strings.TrimSpace(productID), strings.TrimSpace(movementType), strings.TrimSpace(referenceType))
}

func (s *stockMovementService) GetStockMovement(companyID string, id string) (models.StockMovement, error) {
	return s.repo.FindByID(companyID, strings.TrimSpace(id))
}
