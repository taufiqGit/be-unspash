package repositories

import (
	"database/sql"
	"fmt"
	"gowes/models"
)

type TaxRepository interface {
	FindAll(companyID string, params models.PaginationParams) ([]models.Tax, int, error)
	FindByID(id string) (models.Tax, error)
	Create(tax models.Tax) (models.Tax, error)
	Update(tax models.Tax) (models.Tax, error)
	Delete(id string) error
}

type taxRepository struct {
	db *sql.DB
}

func NewTaxRepository(db *sql.DB) TaxRepository {
	return &taxRepository{db: db}
}

func (r *taxRepository) FindAll(companyID string, params models.PaginationParams) ([]models.Tax, int, error) {
	baseQuery := " FROM taxes WHERE company_id = $1"
	args := []interface{}{companyID}
	argIdx := 2

	if params.Search != "" {
		baseQuery += fmt.Sprintf(" AND name ILIKE $%d", argIdx)
		args = append(args, "%"+params.Search+"%")
		argIdx++
	}

	// Count total
	var total int
	countQuery := "SELECT COUNT(*)" + baseQuery
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Validate sort column
	allowedSorts := map[string]bool{"name": true, "rate": true, "created_at": true, "updated_at": true}
	sortBy := "created_at"
	if allowedSorts[params.SortBy] {
		sortBy = params.SortBy
	}

	sortOrder := "DESC"
	if params.SortOrder == "ASC" {
		sortOrder = "ASC"
	}

	query := "SELECT id, name, rate, company_id, created_at, updated_at" + baseQuery
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	offset := (params.Page - 1) * params.Limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, params.Limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	taxes := []models.Tax{}
	for rows.Next() {
		var t models.Tax
		if err := rows.Scan(&t.ID, &t.Name, &t.Rate, &t.CompanyID, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, 0, err
		}
		taxes = append(taxes, t)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return taxes, total, nil
}

func (r *taxRepository) FindByID(id string) (models.Tax, error) {
	row := r.db.QueryRow(`
		SELECT id, name, rate, company_id, created_at, updated_at
		FROM taxes
		WHERE id = $1
	`, id)

	var t models.Tax
	if err := row.Scan(&t.ID, &t.Name, &t.Rate, &t.CompanyID, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return models.Tax{}, err
	}
	return t, nil
}

func (r *taxRepository) Create(tax models.Tax) (models.Tax, error) {
	err := r.db.QueryRow(`
		INSERT INTO taxes (company_id, name, rate, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, tax.CompanyID, tax.Name, tax.Rate, tax.CreatedAt, tax.UpdatedAt).Scan(&tax.ID)

	if err != nil {
		return models.Tax{}, err
	}
	return tax, nil
}

func (r *taxRepository) Update(tax models.Tax) (models.Tax, error) {
	err := r.db.QueryRow(`
		UPDATE taxes
		SET name = $1, rate = $2, updated_at = $3
		WHERE id = $4
		RETURNING id
	`, tax.Name, tax.Rate, tax.UpdatedAt, tax.ID).Scan(&tax.ID)

	if err != nil {
		return models.Tax{}, err
	}
	return tax, nil
}

func (r *taxRepository) Delete(id string) error {
	res, err := r.db.Exec(`DELETE FROM taxes WHERE id = $1`, id)
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
