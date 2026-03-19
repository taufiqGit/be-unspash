package services

import (
	"errors"
	"gowes/models"
	"gowes/repositories"
	"strings"
	"time"
)

var (
	ErrTaxNameRequired = errors.New("tax name is required")
	ErrTaxRateInvalid  = errors.New("tax rate must be between 0 and 100")
)

type TaxService interface {
	ListTaxes(companyID string, params models.PaginationParams) ([]models.Tax, int, error)
	GetTax(id string) (models.Tax, error)
	CreateTax(companyID string, in models.TaxInput) (models.Tax, error)
	UpdateTax(id string, in models.TaxInput) (models.Tax, error)
	DeleteTax(id string) error
}

type taxService struct {
	repo repositories.TaxRepository
}

func NewTaxService(repo repositories.TaxRepository) TaxService {
	return &taxService{repo: repo}
}

func (s *taxService) ListTaxes(companyID string, params models.PaginationParams) ([]models.Tax, int, error) {
	return s.repo.FindAll(companyID, params)
}

func (s *taxService) GetTax(id string) (models.Tax, error) {
	return s.repo.FindByID(id)
}

func (s *taxService) CreateTax(companyID string, in models.TaxInput) (models.Tax, error) {
	if err := validateTaxInput(in); err != nil {
		return models.Tax{}, err
	}

	t := models.Tax{
		CompanyID: companyID,
		Name:      strings.TrimSpace(in.Name),
		Rate:      in.Rate,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	return s.repo.Create(t)
}

func (s *taxService) UpdateTax(id string, in models.TaxInput) (models.Tax, error) {
	if err := validateTaxInput(in); err != nil {
		return models.Tax{}, err
	}

	existing, err := s.repo.FindByID(id)
	if err != nil {
		return models.Tax{}, err
	}

	existing.Name = strings.TrimSpace(in.Name)
	existing.Rate = in.Rate
	existing.UpdatedAt = time.Now().UTC()

	return s.repo.Update(existing)
}

func (s *taxService) DeleteTax(id string) error {
	return s.repo.Delete(id)
}

func validateTaxInput(in models.TaxInput) error {
	if strings.TrimSpace(in.Name) == "" {
		return ErrTaxNameRequired
	}
	if in.Rate < 0 || in.Rate > 100 {
		return ErrTaxRateInvalid
	}
	return nil
}
