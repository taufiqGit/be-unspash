package services

import (
	"gowes/models"
	"gowes/repositories"
)

type OutletService interface {
	FindAll(companyID string, params models.PaginationParams) ([]models.Outlet, int, error)
	Create(companyID string, outlet models.OutletInput) (models.Outlet, error)
	//	Update(outlet models.OutletInput, id string) (models.Outlet, error)
	FindByID(id string) (models.Outlet, error)
	// Delete(id string) error
}

type outletService struct {
	repo repositories.OutletRepository
}

func NewOutletService(repo repositories.OutletRepository) OutletService {
	return &outletService{repo: repo}
}

func (s *outletService) FindAll(companyID string, params models.PaginationParams) ([]models.Outlet, int, error) {
	outlets, total, err := s.repo.FindAll(companyID, params)
	if err != nil {
		return []models.Outlet{}, 0, err
	}
	return outlets, total, nil
}

func (s *outletService) Create(companyID string, outlet models.OutletInput) (models.Outlet, error) {
	return s.repo.Create(&outlet, companyID)
}

func (s *outletService) FindByID(id string) (models.Outlet, error) {
	return s.repo.FindByID(id)
}
