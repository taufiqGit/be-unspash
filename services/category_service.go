package services

import (
	"gowes/models"
	"gowes/repositories"
	"time"
)

type CategoryService interface {
	ListCategories() ([]models.Category, error)
	GetCategory(id int) (models.Category, error)
	CreateCategory(in models.CategoryInput) (models.Category, error)
	UpdateCategory(id int, in models.CategoryInput) (models.Category, error)
	DeleteCategory(id int) error
}

type categoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) ListCategories() ([]models.Category, error) {
	return s.repo.FindAll()
}

func (s *categoryService) GetCategory(id int) (models.Category, error) {
	return s.repo.FindByID(id)
}

func (s *categoryService) CreateCategory(in models.CategoryInput) (models.Category, error) {
	c := models.Category{
		Name:        in.Name,
		Description: in.Description,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	return s.repo.Create(c)
}

func (s *categoryService) UpdateCategory(id int, in models.CategoryInput) (models.Category, error) {
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

func (s *categoryService) DeleteCategory(id int) error {
	return s.repo.Delete(id)
}
