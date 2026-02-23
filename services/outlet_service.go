package services

import (
	"context"
	"gowes/models"
	"gowes/repositories"
)

type OutletService interface {
	FindAll(companyID string, params models.PaginationParams) ([]models.Outlet, int, error)
	Create(companyID string, outlet models.OutletInput) (models.Outlet, error)
	Update(outlet models.OutletInput, id string) (models.Outlet, error)
	FindByID(id string) (models.Outlet, error)
	Delete(id string) error
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
	ctx := context.Background()
	return s.repo.Create(&outlet, companyID, ctx, nil)
}

func (s *outletService) FindByID(id string) (models.Outlet, error) {
	return s.repo.FindByID(id)
}

func (s *outletService) Update(payload models.OutletInput, id string) (models.Outlet, error) {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return models.Outlet{}, err
	}

	return s.repo.Update(&payload, id)
}

func (s *outletService) Delete(id string) error {
	return s.repo.Delete(id)
}
