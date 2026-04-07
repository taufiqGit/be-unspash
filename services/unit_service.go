package services

import (
	"errors"
	"gowes/models"
	"gowes/repositories"
	"strings"
	"time"
)

var (
	ErrUnitNameRequired   = errors.New("unit name is required")
	ErrUnitSymbolRequired = errors.New("unit symbol is required")
	ErrUnitTypeRequired   = errors.New("unit type is required")
)

type UnitService interface {
	ListUnits(companyID string, params models.PaginationParams) ([]models.Unit, int, error)
	GetUnit(id string, companyID string) (models.Unit, error)
	CreateUnit(companyID string, userID string, input models.UnitInput) (models.Unit, error)
	UpdateUnit(id string, companyID string, userID string, input models.UnitInput) (models.Unit, error)
	DeleteUnit(id string, companyID string) error
}

type unitService struct {
	repo repositories.UnitRepository
}

func NewUnitService(repo repositories.UnitRepository) UnitService {
	return &unitService{repo: repo}
}

func (s *unitService) ListUnits(companyID string, params models.PaginationParams) ([]models.Unit, int, error) {
	return s.repo.FindAll(companyID, params)
}

func (s *unitService) GetUnit(id string, companyID string) (models.Unit, error) {
	return s.repo.FindByID(id, companyID)
}

func (s *unitService) CreateUnit(companyID string, userID string, input models.UnitInput) (models.Unit, error) {
	if err := validateUnitInput(input); err != nil {
		return models.Unit{}, err
	}

	now := time.Now().UTC()
	unit := models.Unit{
		Name:      strings.TrimSpace(input.Name),
		Symbol:    strings.TrimSpace(input.Symbol),
		Type:      strings.TrimSpace(input.Type),
		CompanyID: companyID,
		CreatedBy: userID,
		UpdatedBy: userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return s.repo.Create(unit)
}

func (s *unitService) UpdateUnit(id string, companyID string, userID string, input models.UnitInput) (models.Unit, error) {
	if err := validateUnitInput(input); err != nil {
		return models.Unit{}, err
	}

	unit, err := s.repo.FindByID(id, companyID)
	if err != nil {
		return models.Unit{}, err
	}

	unit.Name = strings.TrimSpace(input.Name)
	unit.Symbol = strings.TrimSpace(input.Symbol)
	unit.Type = strings.TrimSpace(input.Type)
	unit.CompanyID = companyID
	unit.UpdatedBy = userID
	unit.UpdatedAt = time.Now().UTC()

	return s.repo.Update(unit)
}

func (s *unitService) DeleteUnit(id string, companyID string) error {
	return s.repo.Delete(id, companyID)
}

func validateUnitInput(input models.UnitInput) error {
	if strings.TrimSpace(input.Name) == "" {
		return ErrUnitNameRequired
	}
	if strings.TrimSpace(input.Symbol) == "" {
		return ErrUnitSymbolRequired
	}
	if strings.TrimSpace(input.Type) == "" {
		return ErrUnitTypeRequired
	}

	return nil
}
