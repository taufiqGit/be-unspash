package repositories

import (
	"database/sql"
	"fmt"
	"gowes/models"
)

type RecipeRepository interface {
	FindAll(companyID string, params models.PaginationParams) ([]models.Recipe, int, error)
	FindByID(id string, companyID string) (models.Recipe, error)
	Create(recipe models.Recipe) (models.Recipe, error)
	Update(recipe models.Recipe) (models.Recipe, error)
	Delete(id string, companyID string) error
}

type recipeRepository struct {
	db *sql.DB
}

func NewRecipeRepository(db *sql.DB) RecipeRepository {
	return &recipeRepository{db: db}
}

func (r *recipeRepository) FindAll(companyID string, params models.PaginationParams) ([]models.Recipe, int, error) {
	baseQuery := " FROM recipes WHERE company_id = $1"
	args := []interface{}{companyID}
	argIdx := 2

	if params.Search != "" {
		baseQuery += fmt.Sprintf(" AND (product_id::text ILIKE $%d OR ingredient_id::text ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+params.Search+"%")
		argIdx++
	}

	var total int
	countQuery := "SELECT COUNT(*)" + baseQuery
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	allowedSorts := map[string]bool{"product_id": true, "ingredient_id": true, "created_at": true, "updated_at": true}
	sortBy := "created_at"
	if allowedSorts[params.SortBy] {
		sortBy = params.SortBy
	}

	sortOrder := "DESC"
	if params.SortOrder == "ASC" || params.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	query := "SELECT id, company_id, product_id, ingredient_id, is_active, created_by, created_at, updated_at, updated_by" + baseQuery
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	offset := (params.Page - 1) * params.Limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, params.Limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	recipes := []models.Recipe{}
	for rows.Next() {
		var recipe models.Recipe
		if err := rows.Scan(&recipe.ID, &recipe.CompanyID, &recipe.ProductID, &recipe.IngredientID, &recipe.IsActive, &recipe.CreatedBy, &recipe.CreatedAt, &recipe.UpdatedAt, &recipe.UpdatedBy); err != nil {
			return nil, 0, err
		}
		recipes = append(recipes, recipe)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return recipes, total, nil
}

func (r *recipeRepository) FindByID(id string, companyID string) (models.Recipe, error) {
	row := r.db.QueryRow(`
		SELECT id, company_id, product_id, ingredient_id, is_active, created_by, created_at, updated_at, updated_by
		FROM recipes
		WHERE id = $1 AND company_id = $2
	`, id, companyID)

	var recipe models.Recipe
	if err := row.Scan(&recipe.ID, &recipe.CompanyID, &recipe.ProductID, &recipe.IngredientID, &recipe.IsActive, &recipe.CreatedBy, &recipe.CreatedAt, &recipe.UpdatedAt, &recipe.UpdatedBy); err != nil {
		return models.Recipe{}, err
	}

	return recipe, nil
}

func (r *recipeRepository) Create(recipe models.Recipe) (models.Recipe, error) {
	err := r.db.QueryRow(`
		INSERT INTO recipes (company_id, product_id, ingredient_id, is_active, created_by, created_at, updated_at, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`, recipe.CompanyID, recipe.ProductID, recipe.IngredientID, recipe.IsActive, recipe.CreatedBy, recipe.CreatedAt, recipe.UpdatedAt, recipe.UpdatedBy).Scan(&recipe.ID)
	if err != nil {
		return models.Recipe{}, err
	}

	return recipe, nil
}

func (r *recipeRepository) Update(recipe models.Recipe) (models.Recipe, error) {
	err := r.db.QueryRow(`
		UPDATE recipes
		SET product_id = $1, ingredient_id = $2, is_active = $3, updated_at = $4, updated_by = $5
		WHERE id = $6 AND company_id = $7
		RETURNING id
	`, recipe.ProductID, recipe.IngredientID, recipe.IsActive, recipe.UpdatedAt, recipe.UpdatedBy, recipe.ID, recipe.CompanyID).Scan(&recipe.ID)
	if err != nil {
		return models.Recipe{}, err
	}

	return recipe, nil
}

func (r *recipeRepository) Delete(id string, companyID string) error {
	res, err := r.db.Exec(`DELETE FROM recipes WHERE id = $1 AND company_id = $2`, id, companyID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}
