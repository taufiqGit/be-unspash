package repositories

import (
	"database/sql"
	"fmt"
	"gowes/models"
)

type UnitRepository interface {
	FindAll(companyID string, params models.PaginationParams) ([]models.Unit, int, error)
	FindByID(id string, companyID string) (models.Unit, error)
	Create(unit models.Unit) (models.Unit, error)
	Update(unit models.Unit) (models.Unit, error)
	Delete(id string, companyID string) error
}

type unitRepository struct {
	db *sql.DB
}

func NewUnitRepository(db *sql.DB) UnitRepository {
	return &unitRepository{db: db}
}

func (r *unitRepository) FindAll(companyID string, params models.PaginationParams) ([]models.Unit, int, error) {
	baseQuery := " FROM units WHERE company_id = $1"
	args := []interface{}{companyID}
	argIdx := 2

	if params.Search != "" {
		baseQuery += fmt.Sprintf(" AND (name ILIKE $%d OR symbol ILIKE $%d OR type ILIKE $%d)", argIdx, argIdx+1, argIdx+2)
		args = append(args, "%"+params.Search+"%", "%"+params.Search+"%", "%"+params.Search+"%")
		argIdx += 3
	}

	var total int
	countQuery := "SELECT COUNT(*)" + baseQuery
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	allowedSorts := map[string]bool{"name": true, "symbol": true, "type": true, "created_at": true, "updated_at": true}
	sortBy := "created_at"
	if allowedSorts[params.SortBy] {
		sortBy = params.SortBy
	}

	sortOrder := "DESC"
	if params.SortOrder == "ASC" {
		sortOrder = "ASC"
	}

	query := "SELECT id, name, symbol, type, company_id, created_by, updated_by, created_at, updated_at" + baseQuery
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	offset := (params.Page - 1) * params.Limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, params.Limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	units := []models.Unit{}
	for rows.Next() {
		var unit models.Unit
		if err := rows.Scan(&unit.ID, &unit.Name, &unit.Symbol, &unit.Type, &unit.CompanyID, &unit.CreatedBy, &unit.UpdatedBy, &unit.CreatedAt, &unit.UpdatedAt); err != nil {
			return nil, 0, err
		}
		units = append(units, unit)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return units, total, nil
}

func (r *unitRepository) FindByID(id string, companyID string) (models.Unit, error) {
	row := r.db.QueryRow(`
		SELECT id, name, symbol, type, company_id, created_by, updated_by, created_at, updated_at
		FROM units
		WHERE id = $1 AND company_id = $2
	`, id, companyID)

	var unit models.Unit
	if err := row.Scan(&unit.ID, &unit.Name, &unit.Symbol, &unit.Type, &unit.CompanyID, &unit.CreatedBy, &unit.UpdatedBy, &unit.CreatedAt, &unit.UpdatedAt); err != nil {
		return models.Unit{}, err
	}

	return unit, nil
}

func (r *unitRepository) Create(unit models.Unit) (models.Unit, error) {
	err := r.db.QueryRow(`
		INSERT INTO units (name, symbol, type, company_id, created_by, updated_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`, unit.Name, unit.Symbol, unit.Type, unit.CompanyID, unit.CreatedBy, unit.UpdatedBy, unit.CreatedAt, unit.UpdatedAt).Scan(&unit.ID)
	if err != nil {
		return models.Unit{}, err
	}

	return unit, nil
}

func (r *unitRepository) Update(unit models.Unit) (models.Unit, error) {
	err := r.db.QueryRow(`
		UPDATE units
		SET name = $1, symbol = $2, type = $3, company_id = $4, updated_by = $5, updated_at = $6
		WHERE id = $7 AND company_id = $8
		RETURNING id
	`, unit.Name, unit.Symbol, unit.Type, unit.CompanyID, unit.UpdatedBy, unit.UpdatedAt, unit.ID, unit.CompanyID).Scan(&unit.ID)
	if err != nil {
		return models.Unit{}, err
	}

	return unit, nil
}

func (r *unitRepository) Delete(id string, companyID string) error {
	res, err := r.db.Exec(`DELETE FROM units WHERE id = $1 AND company_id = $2`, id, companyID)
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
