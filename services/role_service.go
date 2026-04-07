package services

import (
	"errors"
	"gowes/models"
	"gowes/repositories"
	"strings"
	"time"
)

var ErrRoleNameRequired = errors.New("role name is required")

type RoleService interface {
	ListRoles(companyID string, params models.PaginationParams) ([]models.Role, int, error)
	GetRole(id string, companyID string) (models.Role, error)
	CreateRole(companyID string, userID string, input models.RoleInput) (models.Role, error)
	UpdateRole(id string, companyID string, userID string, input models.RoleInput) (models.Role, error)
	DeleteRole(id string, companyID string) error
}

type roleService struct {
	repo repositories.RoleRepository
}

func NewRoleService(repo repositories.RoleRepository) RoleService {
	return &roleService{repo: repo}
}

func (s *roleService) ListRoles(companyID string, params models.PaginationParams) ([]models.Role, int, error) {
	return s.repo.FindAll(companyID, params)
}

func (s *roleService) GetRole(id string, companyID string) (models.Role, error) {
	return s.repo.FindByID(id, companyID)
}

func (s *roleService) CreateRole(companyID string, userID string, input models.RoleInput) (models.Role, error) {
	if err := validateRoleInput(input); err != nil {
		return models.Role{}, err
	}

	now := time.Now().UTC()
	role := models.Role{
		Name:      strings.TrimSpace(input.Name),
		CompanyID: companyID,
		CreatedBy: userID,
		UpdatedBy: userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return s.repo.Create(role)
}

func (s *roleService) UpdateRole(id string, companyID string, userID string, input models.RoleInput) (models.Role, error) {
	if err := validateRoleInput(input); err != nil {
		return models.Role{}, err
	}

	role, err := s.repo.FindByID(id, companyID)
	if err != nil {
		return models.Role{}, err
	}

	role.Name = strings.TrimSpace(input.Name)
	role.CompanyID = companyID
	role.UpdatedBy = userID
	role.UpdatedAt = time.Now().UTC()

	return s.repo.Update(role)
}

func (s *roleService) DeleteRole(id string, companyID string) error {
	return s.repo.Delete(id, companyID)
}

func validateRoleInput(input models.RoleInput) error {
	if strings.TrimSpace(input.Name) == "" {
		return ErrRoleNameRequired
	}

	return nil
}
