package repositories

import (
	"database/sql"
	"fmt"
	"gowes/models"
	"strings"
)

type StockRepository interface {
	FindAll(companyID string, params models.PaginationParams, outletID string, productID string) ([]models.StockPerOutlet, int, error)
	FindByOutletAndProduct(companyID string, outletID string, productID string) (models.StockPerOutlet, error)
}

type stockRepository struct {
	db *sql.DB
}

func NewStockRepository(db *sql.DB) StockRepository {
	return &stockRepository{db: db}
}

func (r *stockRepository) FindAll(companyID string, params models.PaginationParams, outletID string, productID string) ([]models.StockPerOutlet, int, error) {
	baseQuery := `
		FROM stocks s
		JOIN products p ON s.product_id = p.id
		JOIN outlets o ON s.outlet_id = o.id
		WHERE o.company_id = $1
	`
	args := []interface{}{companyID}
	argIdx := 2

	if strings.TrimSpace(outletID) != "" {
		baseQuery += fmt.Sprintf(" AND s.outlet_id = $%d", argIdx)
		args = append(args, strings.TrimSpace(outletID))
		argIdx++
	}

	if strings.TrimSpace(productID) != "" {
		baseQuery += fmt.Sprintf(" AND s.product_id = $%d", argIdx)
		args = append(args, strings.TrimSpace(productID))
		argIdx++
	}

	if strings.TrimSpace(params.Search) != "" {
		baseQuery += fmt.Sprintf(" AND (p.name ILIKE $%d OR p.sku ILIKE $%d OR o.name ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, "%"+strings.TrimSpace(params.Search)+"%")
		argIdx++
	}

	var total int
	if err := r.db.QueryRow("SELECT COUNT(*)"+baseQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	allowedSorts := map[string]string{
		"created_at": "s.id",
		"qty":        "s.qty",
		"product":    "p.name",
		"outlet":     "o.name",
	}
	sortBy := "s.id"
	if col, ok := allowedSorts[params.SortBy]; ok {
		sortBy = col
	}

	sortOrder := "DESC"
	if params.SortOrder == "ASC" {
		sortOrder = "ASC"
	}

	query := `SELECT s.id, s.product_id, p.name, p.sku, s.outlet_id, o.name, s.qty` + baseQuery
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	offset := (params.Page - 1) * params.Limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, params.Limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	stocks := []models.StockPerOutlet{}
	for rows.Next() {
		var s models.StockPerOutlet
		if err := rows.Scan(
			&s.ID,
			&s.ProductID,
			&s.ProductName,
			&s.ProductSKU,
			&s.OutletID,
			&s.OutletName,
			&s.Qty,
		); err != nil {
			return nil, 0, err
		}
		stocks = append(stocks, s)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return stocks, total, nil
}

func (r *stockRepository) FindByOutletAndProduct(companyID string, outletID string, productID string) (models.StockPerOutlet, error) {
	row := r.db.QueryRow(`
		SELECT s.id, s.product_id, p.name, p.sku, s.outlet_id, o.name, s.qty
		FROM stocks s
		JOIN products p ON s.product_id = p.id
		JOIN outlets o ON s.outlet_id = o.id
		WHERE o.company_id = $1 AND s.outlet_id = $2 AND s.product_id = $3
	`, companyID, outletID, productID)

	var s models.StockPerOutlet
	if err := row.Scan(
		&s.ID,
		&s.ProductID,
		&s.ProductName,
		&s.ProductSKU,
		&s.OutletID,
		&s.OutletName,
		&s.Qty,
	); err != nil {
		return models.StockPerOutlet{}, err
	}

	return s, nil
}
