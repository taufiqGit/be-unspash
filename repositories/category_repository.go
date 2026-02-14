package repositories

import (
	"database/sql"
	"fmt"
	"gowes/models"
	"log"
)

type CategoryRepository interface {
	FindAll(companyID string, params models.PaginationParams) ([]models.Category, int, error)
	FindByID(id string) (models.Category, error)
	Create(category models.Category) (models.Category, error)
	Update(category models.Category) (models.Category, error)
	Delete(id string) error
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindAll(companyID string, params models.PaginationParams) ([]models.Category, int, error) {
	// 1. Base Query for Count and Data
	baseQuery := " FROM categories WHERE company_id = $1"
	args := []interface{}{companyID}
	argIdx := 2

	if params.Search != "" {
		baseQuery += fmt.Sprintf(" AND name ILIKE $%d", argIdx)
		args = append(args, "%"+params.Search+"%")
		argIdx++
	}

	// 2. Get Total Count
	var total int
	countQuery := "SELECT COUNT(*)" + baseQuery
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// 3. Get Data with Pagination
	// Validate sort column to prevent SQL Injection
	allowedSorts := map[string]bool{"name": true, "created_at": true, "updated_at": true}
	sortBy := "created_at"
	if allowedSorts[params.SortBy] {
		sortBy = params.SortBy
	}

	query := "SELECT id, company_id, name, description, created_at, updated_at" + baseQuery
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, params.SortOrder)

	offset := (params.Page - 1) * params.Limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, params.Limit, offset)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	fmt.Println(params, "asuw")

	var categories = []models.Category{}
	for rows.Next() {
		fmt.Println("Scanning category...")
		var c models.Category
		if err := rows.Scan(&c.ID, &c.CompanyID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			log.Println("Error scanning category:", err)
			continue
		}
		categories = append(categories, c)
	}
	fmt.Println("Total categories:", len(categories))
	return categories, total, nil
}

func (r *categoryRepository) FindByID(id string) (models.Category, error) {
	row := r.db.QueryRow(`
		SELECT id, company_id, name, description, created_at, updated_at 
		FROM categories 
		WHERE id = $1
	`, id)

	var c models.Category
	if err := row.Scan(&c.ID, &c.CompanyID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return models.Category{}, err
	}
	return c, nil
}

func (r *categoryRepository) Create(category models.Category) (models.Category, error) {
	err := r.db.QueryRow(`
		INSERT INTO categories (company_id, name, description, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id
	`, category.CompanyID, category.Name, category.Description, category.CreatedAt, category.UpdatedAt).Scan(&category.ID)

	if err != nil {
		return models.Category{}, err
	}
	return category, nil
}

func (r *categoryRepository) Update(category models.Category) (models.Category, error) {
	err := r.db.QueryRow(`
		UPDATE categories 
		SET name = $1, description = $2, updated_at = $3 
		WHERE id = $4 
		RETURNING id
	`, category.Name, category.Description, category.UpdatedAt, category.ID).Scan(&category.ID)

	if err != nil {
		return models.Category{}, err
	}
	return category, nil
}

func (r *categoryRepository) Delete(id string) error {
	res, err := r.db.Exec(`DELETE FROM categories WHERE id = $1`, id)
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
