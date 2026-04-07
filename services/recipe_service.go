package services

import (
	"errors"
	"gowes/models"
	"gowes/repositories"
	"strings"
	"time"
)

var (
	ErrRecipeProductRequired    = errors.New("recipe product_id is required")
	ErrRecipeIngredientRequired = errors.New("recipe ingredient_id is required")
)

type RecipeService interface {
	ListRecipes(companyID string, params models.PaginationParams) ([]models.Recipe, int, error)
	GetRecipe(id string, companyID string) (models.Recipe, error)
	CreateRecipe(companyID string, userID string, in models.RecipeInput) (models.Recipe, error)
	UpdateRecipe(id string, companyID string, userID string, in models.RecipeInput) (models.Recipe, error)
	DeleteRecipe(id string, companyID string) error
}

type recipeService struct {
	repo repositories.RecipeRepository
}

func NewRecipeService(repo repositories.RecipeRepository) RecipeService {
	return &recipeService{repo: repo}
}

func (s *recipeService) ListRecipes(companyID string, params models.PaginationParams) ([]models.Recipe, int, error) {
	return s.repo.FindAll(companyID, params)
}

func (s *recipeService) GetRecipe(id string, companyID string) (models.Recipe, error) {
	return s.repo.FindByID(id, companyID)
}

func (s *recipeService) CreateRecipe(companyID string, userID string, in models.RecipeInput) (models.Recipe, error) {
	if err := validateRecipeInput(in); err != nil {
		return models.Recipe{}, err
	}

	now := time.Now().UTC()
	isActive := true
	if in.IsActive != nil {
		isActive = *in.IsActive
	}

	recipe := models.Recipe{
		CompanyID:    companyID,
		ProductID:    strings.TrimSpace(in.ProductID),
		IngredientID: strings.TrimSpace(in.IngredientID),
		IsActive:     isActive,
		CreatedBy:    userID,
		CreatedAt:    now,
		UpdatedAt:    now,
		UpdatedBy:    userID,
	}

	return s.repo.Create(recipe)
}

func (s *recipeService) UpdateRecipe(id string, companyID string, userID string, in models.RecipeInput) (models.Recipe, error) {
	if err := validateRecipeInput(in); err != nil {
		return models.Recipe{}, err
	}

	recipe, err := s.repo.FindByID(id, companyID)
	if err != nil {
		return models.Recipe{}, err
	}

	recipe.ProductID = strings.TrimSpace(in.ProductID)
	recipe.IngredientID = strings.TrimSpace(in.IngredientID)
	if in.IsActive != nil {
		recipe.IsActive = *in.IsActive
	}
	recipe.UpdatedBy = userID
	recipe.UpdatedAt = time.Now().UTC()

	return s.repo.Update(recipe)
}

func (s *recipeService) DeleteRecipe(id string, companyID string) error {
	return s.repo.Delete(id, companyID)
}

func validateRecipeInput(in models.RecipeInput) error {
	if strings.TrimSpace(in.ProductID) == "" {
		return ErrRecipeProductRequired
	}
	if strings.TrimSpace(in.IngredientID) == "" {
		return ErrRecipeIngredientRequired
	}
	return nil
}
