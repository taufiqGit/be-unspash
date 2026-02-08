package services

import (
	"gowes/models"
	"gowes/repositories"
	"time"
)

type CategoryService interface {
	ListCategories(companyID string, params models.PaginationParams) ([]models.Category, error)
	GetCategory(id string) (models.Category, error)
	CreateCategory(in models.CategoryInput, companyID string) (models.Category, error)
	UpdateCategory(id string, in models.CategoryInput) (models.Category, error)
	DeleteCategory(id string) error
}

type categoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) ListCategories(companyID string, params models.PaginationParams) ([]models.Category, error) {
	categories, _, err := s.repo.FindAll(companyID, params)
	return categories, err
}

func (s *categoryService) GetCategory(id string) (models.Category, error) {
	return s.repo.FindByID(id)
}

func (s *categoryService) CreateCategory(in models.CategoryInput, companyID string) (models.Category, error) {
	c := models.Category{
		CompanyID:   companyID,
		Name:        in.Name,
		Description: in.Description,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	return s.repo.Create(c)
}

func (s *categoryService) UpdateCategory(id string, in models.CategoryInput) (models.Category, error) {
	// Check if exists
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return models.Category{}, err
	}

	existing.Name = in.Name
	existing.Description = in.Description
	existing.UpdatedAt = time.Now().UTC()

	return s.repo.Update(existing)
}

func (s *categoryService) DeleteCategory(id string) error {
	return s.repo.Delete(id)
}
