package services

import (
	"gowes/models"
	"gowes/repositories"
)

type OrderTypeService interface {
	FindAll(companyID string, params models.PaginationParams) ([]models.OrderType, int, error)
	Create(companyID string, orderType models.OrderTypeInput) (models.OrderType, error)
	Update(orderType models.OrderTypeInput, id string) (models.OrderType, error)
	FindByID(id string) (models.OrderType, error)
	Delete(id string) error
}

type orderTypeService struct {
	repo repositories.OrderTypeRepository
}

func NewOrderTypeService(repo repositories.OrderTypeRepository) OrderTypeService {
	return &orderTypeService{repo: repo}
}

func (s *orderTypeService) Create(companyID string, orderType models.OrderTypeInput) (models.OrderType, error) {
	return s.repo.Create(companyID, orderType)
}

func (s *orderTypeService) FindAll(companyID string, params models.PaginationParams) ([]models.OrderType, int, error) {
	orderTypes, total, err := s.repo.FindAll(companyID, params)

	return orderTypes, total, err
}

func (s *orderTypeService) Update(orderType models.OrderTypeInput, id string) (models.OrderType, error) {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return models.OrderType{}, err
	}
	return s.repo.Update(orderType, id)
}

func (s *orderTypeService) FindByID(id string) (models.OrderType, error) {
	return s.repo.FindByID(id)
}

func (s *orderTypeService) Delete(id string) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}
