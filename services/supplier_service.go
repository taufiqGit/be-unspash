package services

import (
	"errors"
	"gowes/models"
	"gowes/repositories"
	"strings"
	"time"
)

var ErrSupplierNameRequired = errors.New("supplier name is required")

type SupplierService interface {
	ListSuppliers(companyID string, params models.PaginationParams) ([]models.Supplier, int, error)
	GetSupplier(id string, companyID string) (models.Supplier, error)
	CreateSupplier(companyID string, userID string, in models.SupplierInput) (models.Supplier, error)
	UpdateSupplier(id string, companyID string, userID string, in models.SupplierInput) (models.Supplier, error)
	DeleteSupplier(id string, companyID string) error
}

type supplierService struct {
	repo repositories.SupplierRepository
}

func NewSupplierService(repo repositories.SupplierRepository) SupplierService {
	return &supplierService{repo: repo}
}

func (s *supplierService) ListSuppliers(companyID string, params models.PaginationParams) ([]models.Supplier, int, error) {
	return s.repo.FindAll(companyID, params)
}

func (s *supplierService) GetSupplier(id string, companyID string) (models.Supplier, error) {
	return s.repo.FindByID(id, companyID)
}

func (s *supplierService) CreateSupplier(companyID string, userID string, in models.SupplierInput) (models.Supplier, error) {
	if strings.TrimSpace(in.Name) == "" {
		return models.Supplier{}, ErrSupplierNameRequired
	}

	now := time.Now().UTC()
	isActive := true
	if in.IsActive != nil {
		isActive = *in.IsActive
	}

	supplier := models.Supplier{
		Name:        strings.TrimSpace(in.Name),
		CompanyID:   companyID,
		Address:     strings.TrimSpace(in.Address),
		Phone:       strings.TrimSpace(in.Phone),
		Email:       strings.TrimSpace(in.Email),
		CompanyName: strings.TrimSpace(in.CompanyName),
		TaxNumber:   strings.TrimSpace(in.TaxNumber),
		IsActive:    isActive,
		CreatedBy:   userID,
		CreatedAt:   now,
		UpdatedAt:   now,
		UpdatedBy:   userID,
	}

	return s.repo.Create(supplier)
}

func (s *supplierService) UpdateSupplier(id string, companyID string, userID string, in models.SupplierInput) (models.Supplier, error) {
	if strings.TrimSpace(in.Name) == "" {
		return models.Supplier{}, ErrSupplierNameRequired
	}

	supplier, err := s.repo.FindByID(id, companyID)
	if err != nil {
		return models.Supplier{}, err
	}

	supplier.Name = strings.TrimSpace(in.Name)
	supplier.Address = strings.TrimSpace(in.Address)
	supplier.Phone = strings.TrimSpace(in.Phone)
	supplier.Email = strings.TrimSpace(in.Email)
	supplier.CompanyName = strings.TrimSpace(in.CompanyName)
	supplier.TaxNumber = strings.TrimSpace(in.TaxNumber)
	if in.IsActive != nil {
		supplier.IsActive = *in.IsActive
	}
	supplier.UpdatedAt = time.Now().UTC()
	supplier.UpdatedBy = userID

	return s.repo.Update(supplier)
}

func (s *supplierService) DeleteSupplier(id string, companyID string) error {
	return s.repo.Delete(id, companyID)
}
