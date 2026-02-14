package repositories

import (
	"database/sql"
	"fmt"
	"gowes/models"
)

type AddOnRepository interface {
	FindAll(companyID string, params models.PaginationParams) ([]models.AddOn, int, error)
	FindByID(id string) (models.AddOn, error)
	Create(addOn *models.AddOnInput, companyID string) (models.AddOn, error)
	Update(addOn *models.AddOnInput, id string) (models.AddOn, error)
	Delete(id string) error
}

type addOnRepository struct {
	db *sql.DB
}

func NewAddOnRepository(db *sql.DB) AddOnRepository {
	return &addOnRepository{db: db}
}

func (r *addOnRepository) FindAll(companyID string, params models.PaginationParams) ([]models.AddOn, int, error) {
	// 1. Base Query for Count and Data
	baseQuery := " FROM add_ons WHERE company_id = $1"
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

	// 4. Build Data Query
	offset := (params.Page - 1) * params.Limit
	query := "SELECT id, company_id, name, price, is_active, created_at, updated_at" + baseQuery
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, params.SortOrder)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, params.Limit, offset)
	fmt.Println(argIdx, argIdx+1)
	// 5. Execute Data Query
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	addOns := []models.AddOn{}
	for rows.Next() {
		var addOn models.AddOn
		if err := rows.Scan(&addOn.ID, &addOn.CompanyID, &addOn.Name, &addOn.Price, &addOn.IsActive, &addOn.CreatedAt, &addOn.UpdatedAt); err != nil {
			return nil, 0, err
		}
		addOns = append(addOns, addOn)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return addOns, total, nil
}

func (r *addOnRepository) Create(input *models.AddOnInput, companyID string) (models.AddOn, error) {
	query := "INSERT INTO add_ons (company_id, name, price, is_active) VALUES ($1, $2, $3, $4) RETURNING id, company_id, name, price, created_at, updated_at"
	args := []interface{}{companyID, input.Name, input.Price, true}

	var newAddOn models.AddOn
	if err := r.db.QueryRow(query, args...).Scan(&newAddOn.ID, &newAddOn.CompanyID, &newAddOn.Name, &newAddOn.Price, &newAddOn.CreatedAt, &newAddOn.UpdatedAt); err != nil {
		return models.AddOn{}, err
	}

	return newAddOn, nil
}

func (r *addOnRepository) FindByID(id string) (models.AddOn, error) {
	query := "SELECT id, company_id, name, price, is_active, created_at, updated_at FROM add_ons WHERE id = $1"
	var addOn models.AddOn
	if err := r.db.QueryRow(query, id).Scan(&addOn.ID, &addOn.CompanyID, &addOn.Name, &addOn.Price, &addOn.IsActive, &addOn.CreatedAt, &addOn.UpdatedAt); err != nil {
		return models.AddOn{}, err
	}
	return addOn, nil
}

func (r *addOnRepository) Update(input *models.AddOnInput, id string) (models.AddOn, error) {
	query := "UPDATE add_ons SET name = $1, price = $2, is_active = $3 WHERE id = $4 RETURNING id, company_id, name, price, is_active, created_at, updated_at"
	args := []interface{}{input.Name, input.Price, input.IsActive, id}

	var updatedAddOn models.AddOn
	if err := r.db.QueryRow(query, args...).Scan(&updatedAddOn.ID, &updatedAddOn.CompanyID, &updatedAddOn.Name, &updatedAddOn.Price, &updatedAddOn.IsActive, &updatedAddOn.CreatedAt, &updatedAddOn.UpdatedAt); err != nil {
		return models.AddOn{}, err
	}

	return updatedAddOn, nil
}

func (r *addOnRepository) Delete(id string) error {
	query := "DELETE FROM add_ons WHERE id = $1"
	res, err := r.db.Exec(query, id)
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
