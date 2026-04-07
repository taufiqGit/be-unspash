package repositories

import (
	"database/sql"
	"fmt"
	"gowes/models"
)

type RoleRepository interface {
	FindAll(companyID string, params models.PaginationParams) ([]models.Role, int, error)
	FindByID(id string, companyID string) (models.Role, error)
	Create(role models.Role) (models.Role, error)
	Update(role models.Role) (models.Role, error)
	Delete(id string, companyID string) error
}

type roleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindAll(companyID string, params models.PaginationParams) ([]models.Role, int, error) {
	baseQuery := " FROM roles WHERE company_id = $1"
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

	sortOrder := "DESC"
	if params.SortOrder == "ASC" {
		sortOrder = "ASC"
	}

	query := "SELECT id, name, company_id, created_by, updated_by, created_at, updated_at" + baseQuery
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	offset := (params.Page - 1) * params.Limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, params.Limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	roles := []models.Role{}
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.CompanyID, &role.CreatedBy, &role.UpdatedBy, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, 0, err
		}
		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (r *roleRepository) FindByID(id string, companyID string) (models.Role, error) {
	row := r.db.QueryRow(`
		SELECT id, name, company_id, created_by, updated_by, created_at, updated_at
		FROM roles
		WHERE id = $1 AND company_id = $2
	`, id, companyID)

	var role models.Role
	if err := row.Scan(&role.ID, &role.Name, &role.CompanyID, &role.CreatedBy, &role.UpdatedBy, &role.CreatedAt, &role.UpdatedAt); err != nil {
		return models.Role{}, err
	}

	return role, nil
}

func (r *roleRepository) Create(role models.Role) (models.Role, error) {
	err := r.db.QueryRow(`
		INSERT INTO roles (name, company_id, created_by, updated_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, role.Name, role.CompanyID, role.CreatedBy, role.UpdatedBy, role.CreatedAt, role.UpdatedAt).Scan(&role.ID)
	if err != nil {
		return models.Role{}, err
	}

	return role, nil
}

func (r *roleRepository) Update(role models.Role) (models.Role, error) {
	err := r.db.QueryRow(`
		UPDATE roles
		SET name = $1, company_id = $2, updated_by = $3, updated_at = $4
		WHERE id = $5 AND company_id = $6
		RETURNING id
	`, role.Name, role.CompanyID, role.UpdatedBy, role.UpdatedAt, role.ID, role.CompanyID).Scan(&role.ID)
	if err != nil {
		return models.Role{}, err
	}

	return role, nil
}

func (r *roleRepository) Delete(id string, companyID string) error {
	res, err := r.db.Exec(`DELETE FROM roles WHERE id = $1 AND company_id = $2`, id, companyID)
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
