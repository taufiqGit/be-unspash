package repositories

import (
	"database/sql"
	"fmt"
	"gowes/models"
)

type SupplierRepository interface {
	FindAll(companyID string, params models.PaginationParams) ([]models.Supplier, int, error)
	FindByID(id string, companyID string) (models.Supplier, error)
	Create(supplier models.Supplier) (models.Supplier, error)
	Update(supplier models.Supplier) (models.Supplier, error)
	Delete(id string, companyID string) error
}

type supplierRepository struct {
	db *sql.DB
}

func NewSupplierRepository(db *sql.DB) SupplierRepository {
	return &supplierRepository{db: db}
}

func (r *supplierRepository) FindAll(companyID string, params models.PaginationParams) ([]models.Supplier, int, error) {
	baseQuery := " FROM suppliers WHERE company_id = $1"
	args := []interface{}{companyID}
	argIdx := 2

	if params.Search != "" {
		baseQuery += fmt.Sprintf(" AND (name ILIKE $%d OR email ILIKE $%d OR phone ILIKE $%d OR company_name ILIKE $%d)", argIdx, argIdx, argIdx, argIdx)
		args = append(args, "%"+params.Search+"%")
		argIdx++
	}

	var total int
	countQuery := "SELECT COUNT(*)" + baseQuery
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	allowedSorts := map[string]bool{"name": true, "company_name": true, "created_at": true, "updated_at": true}
	sortBy := "created_at"
	if allowedSorts[params.SortBy] {
		sortBy = params.SortBy
	}

	sortOrder := "DESC"
	if params.SortOrder == "ASC" || params.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	query := "SELECT id, name, company_id, address, phone, email, company_name, tax_number, is_active, created_by, created_at, updated_at, updated_by" + baseQuery
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	offset := (params.Page - 1) * params.Limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, params.Limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	suppliers := []models.Supplier{}
	for rows.Next() {
		var supplier models.Supplier
		if err := rows.Scan(&supplier.ID, &supplier.Name, &supplier.CompanyID, &supplier.Address, &supplier.Phone, &supplier.Email, &supplier.CompanyName, &supplier.TaxNumber, &supplier.IsActive, &supplier.CreatedBy, &supplier.CreatedAt, &supplier.UpdatedAt, &supplier.UpdatedBy); err != nil {
			return nil, 0, err
		}
		suppliers = append(suppliers, supplier)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return suppliers, total, nil
}

func (r *supplierRepository) FindByID(id string, companyID string) (models.Supplier, error) {
	row := r.db.QueryRow(`
		SELECT id, name, company_id, address, phone, email, company_name, tax_number, is_active, created_by, created_at, updated_at, updated_by
		FROM suppliers
		WHERE id = $1 AND company_id = $2
	`, id, companyID)

	var supplier models.Supplier
	if err := row.Scan(&supplier.ID, &supplier.Name, &supplier.CompanyID, &supplier.Address, &supplier.Phone, &supplier.Email, &supplier.CompanyName, &supplier.TaxNumber, &supplier.IsActive, &supplier.CreatedBy, &supplier.CreatedAt, &supplier.UpdatedAt, &supplier.UpdatedBy); err != nil {
		return models.Supplier{}, err
	}

	return supplier, nil
}

func (r *supplierRepository) Create(supplier models.Supplier) (models.Supplier, error) {
	err := r.db.QueryRow(`
		INSERT INTO suppliers (name, company_id, address, phone, email, company_name, tax_number, is_active, created_by, created_at, updated_at, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`, supplier.Name, supplier.CompanyID, supplier.Address, supplier.Phone, supplier.Email, supplier.CompanyName, supplier.TaxNumber, supplier.IsActive, supplier.CreatedBy, supplier.CreatedAt, supplier.UpdatedAt, supplier.UpdatedBy).Scan(&supplier.ID)
	if err != nil {
		return models.Supplier{}, err
	}

	return supplier, nil
}

func (r *supplierRepository) Update(supplier models.Supplier) (models.Supplier, error) {
	err := r.db.QueryRow(`
		UPDATE suppliers
		SET name = $1, address = $2, phone = $3, email = $4, company_name = $5, tax_number = $6, is_active = $7, updated_at = $8, updated_by = $9
		WHERE id = $10 AND company_id = $11
		RETURNING id
	`, supplier.Name, supplier.Address, supplier.Phone, supplier.Email, supplier.CompanyName, supplier.TaxNumber, supplier.IsActive, supplier.UpdatedAt, supplier.UpdatedBy, supplier.ID, supplier.CompanyID).Scan(&supplier.ID)
	if err != nil {
		return models.Supplier{}, err
	}

	return supplier, nil
}

func (r *supplierRepository) Delete(id string, companyID string) error {
	res, err := r.db.Exec(`DELETE FROM suppliers WHERE id = $1 AND company_id = $2`, id, companyID)
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
