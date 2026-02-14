package services

import (
	"gowes/models"
	"gowes/repositories"
)

type AddOnService interface {
	FindAll(companyID string, params models.PaginationParams) ([]models.AddOn, int, error)
	Create(addOn *models.AddOnInput, companyID string) (models.AddOn, error)
	Update(addOn *models.AddOnInput, id string) (models.AddOn, error)
	FindById(id string) (models.AddOn, error)
	Delete(id string) error
}

type addOnService struct {
	repo repositories.AddOnRepository
}

func NewAddOnService(repo repositories.AddOnRepository) AddOnService {
	return &addOnService{repo: repo}
}

func (s *addOnService) FindAll(companyID string, params models.PaginationParams) ([]models.AddOn, int, error) {
	addOns, total, err := s.repo.FindAll(companyID, params)

	return addOns, total, err
}

func (s *addOnService) Create(addOn *models.AddOnInput, companyID string) (models.AddOn, error) {
	return s.repo.Create(addOn, companyID)
}

func (s *addOnService) Update(addOn *models.AddOnInput, id string) (models.AddOn, error) {
	return s.repo.Update(addOn, id)
}

func (s *addOnService) FindById(id string) (models.AddOn, error) {
	return s.repo.FindByID(id)
}

func (s *addOnService) Delete(id string) error {
	return s.repo.Delete(id)
}
