package services

import (
	"gowes/models"
	"gowes/repositories"
	"time"
)

type CustomerService interface {
	ListCustomers(companyID string, params models.PaginationParams) ([]models.Customer, int, error)
	GetCustomer(id string) (models.Customer, error)
	CreateCustomer(in models.CustomerInput, companyID string) (models.Customer, error)
	UpdateCustomer(id string, in models.CustomerInput) (models.Customer, error)
	DeleteCustomer(id string) error
}

type customerService struct {
	repo repositories.CustomerRepository
}

func NewCustomerService(repo repositories.CustomerRepository) CustomerService {
	return &customerService{repo: repo}
}

func (s *customerService) ListCustomers(companyID string, params models.PaginationParams) ([]models.Customer, int, error) {
	return s.repo.FindAll(companyID, params)
}

func (s *customerService) GetCustomer(id string) (models.Customer, error) {
	return s.repo.FindByID(id)
}

func (s *customerService) CreateCustomer(in models.CustomerInput, companyID string) (models.Customer, error) {
	c := models.Customer{
		CompanyID: companyID,
		Name:      in.Name,
		Phone:     in.Phone,
		Email:     in.Email,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	return s.repo.Create(c)
}

func (s *customerService) UpdateCustomer(id string, in models.CustomerInput) (models.Customer, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return models.Customer{}, err
	}

	existing.Name = in.Name
	existing.Phone = in.Phone
	existing.Email = in.Email
	existing.UpdatedAt = time.Now().UTC()

	return s.repo.Update(existing)
}

func (s *customerService) DeleteCustomer(id string) error {
	return s.repo.Delete(id)
}
