package repositories

import (
	"database/sql"
	"fmt"
	"gowes/models"
)

type CustomerRepository interface {
	FindAll(companyID string, params models.PaginationParams) ([]models.Customer, int, error)
	FindByID(id string) (models.Customer, error)
	Create(customer models.Customer) (models.Customer, error)
	Update(customer models.Customer) (models.Customer, error)
	Delete(id string) error
}

type customerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) FindAll(companyID string, params models.PaginationParams) ([]models.Customer, int, error) {
	baseQuery := " FROM customers WHERE company_id = $1"
	args := []interface{}{companyID}
	argIdx := 2

	if params.Search != "" {
		baseQuery += fmt.Sprintf(" AND (name ILIKE $%d OR email ILIKE $%d OR phone ILIKE $%d)", argIdx, argIdx, argIdx)
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
	allowedSorts := map[string]bool{"name": true, "created_at": true, "updated_at": true}
	sortBy := "created_at"
	if allowedSorts[params.SortBy] {
		sortBy = params.SortBy
	}

	sortOrder := "DESC"
	if params.SortOrder == "asc" || params.SortOrder == "ASC" {
		sortOrder = "ASC"
	}

	query := "SELECT id, company_id, name, phone, email, created_at, updated_at" + baseQuery
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	offset := (params.Page - 1) * params.Limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, params.Limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	customers := []models.Customer{}
	for rows.Next() {
		var c models.Customer
		if err := rows.Scan(&c.ID, &c.CompanyID, &c.Name, &c.Phone, &c.Email, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, 0, err
		}
		customers = append(customers, c)
	}

	return customers, total, nil
}

func (r *customerRepository) FindByID(id string) (models.Customer, error) {
	row := r.db.QueryRow(`
		SELECT id, company_id, name, phone, email, created_at, updated_at
		FROM customers
		WHERE id = $1
	`, id)

	var c models.Customer
	if err := row.Scan(&c.ID, &c.CompanyID, &c.Name, &c.Phone, &c.Email, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return models.Customer{}, err
	}
	return c, nil
}

func (r *customerRepository) Create(customer models.Customer) (models.Customer, error) {
	err := r.db.QueryRow(`
		INSERT INTO customers (company_id, name, phone, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, customer.CompanyID, customer.Name, customer.Phone, customer.Email, customer.CreatedAt, customer.UpdatedAt).Scan(&customer.ID)

	if err != nil {
		return models.Customer{}, err
	}
	return customer, nil
}

func (r *customerRepository) Update(customer models.Customer) (models.Customer, error) {
	err := r.db.QueryRow(`
		UPDATE customers
		SET name = $1, phone = $2, email = $3, updated_at = $4
		WHERE id = $5
		RETURNING id
	`, customer.Name, customer.Phone, customer.Email, customer.UpdatedAt, customer.ID).Scan(&customer.ID)

	if err != nil {
		return models.Customer{}, err
	}
	return customer, nil
}

func (r *customerRepository) Delete(id string) error {
	res, err := r.db.Exec(`DELETE FROM customers WHERE id = $1`, id)
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
