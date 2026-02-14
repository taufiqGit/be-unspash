package repositories

import (
	"database/sql"
	"fmt"
	"gowes/models"
)

type OutletRepository interface {
	FindAll(companyID string, params models.PaginationParams) ([]models.Outlet, int, error)
	FindByID(id string) (models.Outlet, error)
	Create(outlet *models.OutletInput, companyID string) (models.Outlet, error)
	// Update(outlet *models.OutletInput, id string) (models.Outlet, error)
	// Delete(id string) error
}

type outletRepository struct {
	db *sql.DB
}

func NewOutletRepository(db *sql.DB) OutletRepository {
	return &outletRepository{db: db}
}

func (r *outletRepository) FindAll(companyID string, params models.PaginationParams) ([]models.Outlet, int, error) {
	// 1. Base Query for Count and Data
	baseQuery := " FROM outlets WHERE company_id = $1"
	args := []interface{}{companyID}
	argIdx := 2

	if params.Search != "" {
		baseQuery += fmt.Sprintf(" AND name ILIKE $%d", argIdx)
		args = append(args, "%"+params.Search+"%")
		argIdx++
	}

	var total int
	countQuery := "SELECT COUNT(*)" + baseQuery
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	allowedSorts := map[string]bool{"name": true, "created_at": true, "updated_at": true}
	sortBy := "created_at"
	if allowedSorts[params.SortBy] {
		sortBy = params.SortBy
	}

	query := "SELECT id, code, name, supervisor, address, phone, email, is_active" + baseQuery
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, params.SortOrder)

	offset := (params.Page - 1) * params.Limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, params.Limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var outlets = []models.Outlet{}
	for rows.Next() {
		var outlet models.Outlet
		if err := rows.Scan(&outlet.ID, &outlet.Code, &outlet.Name, &outlet.Supervisor, &outlet.Address, &outlet.Phone, &outlet.Email, &outlet.IsActive); err != nil {
			return nil, 0, err
		}
		outlets = append(outlets, outlet)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return outlets, total, nil
}

func (r *outletRepository) FindByID(id string) (models.Outlet, error) {
	var outlet models.Outlet
	query := "SELECT id, code, name, supervisor, address, phone, email, is_active FROM outlets WHERE id = $1"
	if err := r.db.QueryRow(query, id).Scan(&outlet.ID, &outlet.Code, &outlet.Name, &outlet.Supervisor, &outlet.Address, &outlet.Phone, &outlet.Email, &outlet.IsActive); err != nil {
		return models.Outlet{}, err
	}
	return outlet, nil
}

func (r *outletRepository) Create(outlet *models.OutletInput, companyID string) (models.Outlet, error) {
	var createdOutlet models.Outlet
	query := "INSERT INTO outlets (company_id, code, name, supervisor, address, phone, email, is_active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, code, name, supervisor, address, phone, email, is_active"
	if err := r.db.QueryRow(query, companyID, outlet.Code, outlet.Name, outlet.Supervisor, outlet.Address, outlet.Phone, outlet.Email, outlet.IsActive).Scan(&createdOutlet.ID, &createdOutlet.Code, &createdOutlet.Name, &createdOutlet.Supervisor, &createdOutlet.Address, &createdOutlet.Phone, &createdOutlet.Email, &createdOutlet.IsActive); err != nil {
		return models.Outlet{}, err
	}
	return createdOutlet, nil
}
