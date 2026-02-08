package repositories

import (
	"context"
	"database/sql"
	"gowes/models"
)

type CompanyRepository interface {
	Create(ctx context.Context, tx *sql.Tx, company models.Company) (models.Company, error)
}

type companyRepository struct {
	db *sql.DB
}

func NewCompanyRepository(db *sql.DB) CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) Create(ctx context.Context, tx *sql.Tx, company models.Company) (models.Company, error) {
	query := `
		INSERT INTO company (id, name, created_at, updated_at)
		VALUES (uuid_generate_v4(), $1, $2, $3)
		RETURNING id
	`
	// Use tx if provided, otherwise db
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRowContext(ctx, query, company.Name, company.CreatedAt, company.UpdatedAt)
	} else {
		row = r.db.QueryRowContext(ctx, query, company.Name, company.CreatedAt, company.UpdatedAt)
	}

	err := row.Scan(&company.ID)
	if err != nil {
		return models.Company{}, err
	}
	return company, nil
}
