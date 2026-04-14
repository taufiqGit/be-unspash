package repositories

import (
	"database/sql"
	"fmt"
	"gowes/models"
	"strings"
)

type StockMovementRepository interface {
	FindAll(companyID string, params models.PaginationParams, outletID string, productID string, movementType string, referenceType string) ([]models.StockMovement, int, error)
	FindByID(companyID string, id string) (models.StockMovement, error)
}

type stockMovementRepository struct {
	db *sql.DB
}

func NewStockMovementRepository(db *sql.DB) StockMovementRepository {
	return &stockMovementRepository{db: db}
}

func (r *stockMovementRepository) FindAll(companyID string, params models.PaginationParams, outletID string, productID string, movementType string, referenceType string) ([]models.StockMovement, int, error) {
	baseQuery := `
		FROM stock_movements sm
		JOIN products p ON sm.product_id = p.id
		JOIN outlets o ON sm.outlet_id = o.id
		WHERE o.company_id = $1
	`
	args := []interface{}{companyID}
	argIdx := 2

	if strings.TrimSpace(outletID) != "" {
		baseQuery += fmt.Sprintf(" AND sm.outlet_id = $%d", argIdx)
		args = append(args, strings.TrimSpace(outletID))
		argIdx++
	}
	if strings.TrimSpace(productID) != "" {
		baseQuery += fmt.Sprintf(" AND sm.product_id = $%d", argIdx)
		args = append(args, strings.TrimSpace(productID))
		argIdx++
	}
	if strings.TrimSpace(movementType) != "" {
		baseQuery += fmt.Sprintf(" AND sm.type = $%d", argIdx)
		args = append(args, strings.TrimSpace(movementType))
		argIdx++
	}
	if strings.TrimSpace(referenceType) != "" {
		baseQuery += fmt.Sprintf(" AND sm.reference_type = $%d", argIdx)
		args = append(args, strings.TrimSpace(referenceType))
		argIdx++
	}
	if strings.TrimSpace(params.Search) != "" {
		baseQuery += fmt.Sprintf(" AND (p.name ILIKE $%d OR p.sku ILIKE $%d OR o.name ILIKE $%d OR COALESCE(sm.note, '') ILIKE $%d)", argIdx, argIdx, argIdx, argIdx)
		args = append(args, "%"+strings.TrimSpace(params.Search)+"%")
		argIdx++
	}

	var total int
	if err := r.db.QueryRow("SELECT COUNT(*)"+baseQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	allowedSorts := map[string]string{
		"created_at": "sm.created_at",
		"qty":        "sm.qty",
		"type":       "sm.type",
		"product":    "p.name",
		"outlet":     "o.name",
	}
	sortBy := "sm.created_at"
	if col, ok := allowedSorts[params.SortBy]; ok {
		sortBy = col
	}

	sortOrder := "DESC"
	if params.SortOrder == "ASC" {
		sortOrder = "ASC"
	}

	query := `SELECT sm.id, sm.product_id, p.name, p.sku, sm.outlet_id, o.name, sm.type, sm.qty, sm.reference_type, sm.reference_id, sm.note, sm.created_at` + baseQuery
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	offset := (params.Page - 1) * params.Limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, params.Limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	movements := []models.StockMovement{}
	for rows.Next() {
		var movement models.StockMovement
		var referenceTypeNull sql.NullString
		var referenceIDNull sql.NullString
		var noteNull sql.NullString

		if err := rows.Scan(
			&movement.ID,
			&movement.ProductID,
			&movement.ProductName,
			&movement.ProductSKU,
			&movement.OutletID,
			&movement.OutletName,
			&movement.Type,
			&movement.Qty,
			&referenceTypeNull,
			&referenceIDNull,
			&noteNull,
			&movement.CreatedAt,
		); err != nil {
			return nil, 0, err
		}

		if referenceTypeNull.Valid {
			movement.ReferenceType = models.StockReferenceType(referenceTypeNull.String)
		}
		if referenceIDNull.Valid {
			movement.ReferenceID = referenceIDNull.String
		}
		if noteNull.Valid {
			movement.Note = noteNull.String
		}

		movements = append(movements, movement)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return movements, total, nil
}

func (r *stockMovementRepository) FindByID(companyID string, id string) (models.StockMovement, error) {
	row := r.db.QueryRow(`
		SELECT sm.id, sm.product_id, p.name, p.sku, sm.outlet_id, o.name, sm.type, sm.qty, sm.reference_type, sm.reference_id, sm.note, sm.created_at
		FROM stock_movements sm
		JOIN products p ON sm.product_id = p.id
		JOIN outlets o ON sm.outlet_id = o.id
		WHERE o.company_id = $1 AND sm.id = $2
	`, companyID, id)

	var movement models.StockMovement
	var referenceTypeNull sql.NullString
	var referenceIDNull sql.NullString
	var noteNull sql.NullString
	if err := row.Scan(
		&movement.ID,
		&movement.ProductID,
		&movement.ProductName,
		&movement.ProductSKU,
		&movement.OutletID,
		&movement.OutletName,
		&movement.Type,
		&movement.Qty,
		&referenceTypeNull,
		&referenceIDNull,
		&noteNull,
		&movement.CreatedAt,
	); err != nil {
		return models.StockMovement{}, err
	}

	if referenceTypeNull.Valid {
		movement.ReferenceType = models.StockReferenceType(referenceTypeNull.String)
	}
	if referenceIDNull.Valid {
		movement.ReferenceID = referenceIDNull.String
	}
	if noteNull.Valid {
		movement.Note = noteNull.String
	}

	return movement, nil
}
