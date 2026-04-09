package repositories

import (
	"database/sql"
	"fmt"
	"gowes/models"
)

type PurchaseRepository interface {
	FindAll(companyID string, params models.PaginationParams) ([]models.Purchase, int, error)
	FindByID(id string, companyID string) (models.Purchase, error)
	CreateWithStockMovement(purchase models.Purchase) (models.Purchase, error)
}

type purchaseRepository struct {
	db *sql.DB
}

func NewPurchaseRepository(db *sql.DB) PurchaseRepository {
	return &purchaseRepository{db: db}
}

func (r *purchaseRepository) FindAll(companyID string, params models.PaginationParams) ([]models.Purchase, int, error) {
	baseQuery := " FROM purchases WHERE company_id = $1"
	args := []interface{}{companyID}
	argIdx := 2

	var total int
	if err := r.db.QueryRow("SELECT COUNT(*)"+baseQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	allowedSorts := map[string]bool{
		"created_at":  true,
		"updated_at":  true,
		"grand_total": true,
		"status":      true,
	}
	sortBy := "created_at"
	if allowedSorts[params.SortBy] {
		sortBy = params.SortBy
	}

	sortOrder := "DESC"
	if params.SortOrder == "ASC" {
		sortOrder = "ASC"
	}

	query := `SELECT id, company_id, user_id, outlet_id, payment_method, grand_total, tax_value, paid_amount, change_amount, status, discount_bill, created_at, updated_at` + baseQuery
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	offset := (params.Page - 1) * params.Limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, params.Limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	purchases := []models.Purchase{}
	for rows.Next() {
		var purchase models.Purchase
		if err := rows.Scan(
			&purchase.ID,
			&purchase.CompanyID,
			&purchase.UserID,
			&purchase.OutletID,
			&purchase.PaymentMethod,
			&purchase.GrandTotal,
			&purchase.TaxValue,
			&purchase.PaidAmount,
			&purchase.ChangeAmount,
			&purchase.Status,
			&purchase.DiscountBill,
			&purchase.CreatedAt,
			&purchase.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		purchases = append(purchases, purchase)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return purchases, total, nil
}

func (r *purchaseRepository) FindByID(id string, companyID string) (models.Purchase, error) {
	row := r.db.QueryRow(`
		SELECT id, company_id, user_id, outlet_id, payment_method, grand_total, tax_value, paid_amount, change_amount, status, discount_bill, created_at, updated_at
		FROM purchases
		WHERE id = $1 AND company_id = $2
	`, id, companyID)

	var purchase models.Purchase
	if err := row.Scan(
		&purchase.ID,
		&purchase.CompanyID,
		&purchase.UserID,
		&purchase.OutletID,
		&purchase.PaymentMethod,
		&purchase.GrandTotal,
		&purchase.TaxValue,
		&purchase.PaidAmount,
		&purchase.ChangeAmount,
		&purchase.Status,
		&purchase.DiscountBill,
		&purchase.CreatedAt,
		&purchase.UpdatedAt,
	); err != nil {
		return models.Purchase{}, err
	}

	detailRows, err := r.db.Query(`
		SELECT id, purchase_id, product_id, quantity, price, total
		FROM purchase_details
		WHERE purchase_id = $1
		ORDER BY id ASC
	`, purchase.ID)
	if err != nil {
		return models.Purchase{}, err
	}
	defer detailRows.Close()

	purchase.Details = []models.PurchaseDetail{}
	for detailRows.Next() {
		var detail models.PurchaseDetail
		if err := detailRows.Scan(
			&detail.ID,
			&detail.PurchaseID,
			&detail.ProductID,
			&detail.Quantity,
			&detail.Price,
			&detail.Total,
		); err != nil {
			return models.Purchase{}, err
		}
		purchase.Details = append(purchase.Details, detail)
	}
	if err := detailRows.Err(); err != nil {
		return models.Purchase{}, err
	}

	return purchase, nil
}

func (r *purchaseRepository) CreateWithStockMovement(purchase models.Purchase) (models.Purchase, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return models.Purchase{}, err
	}
	defer tx.Rollback()

	if err := tx.QueryRow(`
		INSERT INTO purchases (
			company_id, user_id, outlet_id, payment_method, grand_total, tax_value, paid_amount, change_amount, status, discount_bill, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, company_id, user_id, outlet_id, payment_method, grand_total, tax_value, paid_amount, change_amount, status, discount_bill, created_at, updated_at
	`,
		purchase.CompanyID,
		purchase.UserID,
		purchase.OutletID,
		purchase.PaymentMethod,
		purchase.GrandTotal,
		purchase.TaxValue,
		purchase.PaidAmount,
		purchase.ChangeAmount,
		purchase.Status,
		purchase.DiscountBill,
		purchase.CreatedAt,
		purchase.UpdatedAt,
	).Scan(
		&purchase.ID,
		&purchase.CompanyID,
		&purchase.UserID,
		&purchase.OutletID,
		&purchase.PaymentMethod,
		&purchase.GrandTotal,
		&purchase.TaxValue,
		&purchase.PaidAmount,
		&purchase.ChangeAmount,
		&purchase.Status,
		&purchase.DiscountBill,
		&purchase.CreatedAt,
		&purchase.UpdatedAt,
	); err != nil {
		return models.Purchase{}, err
	}

	details := make([]models.PurchaseDetail, 0, len(purchase.Details))
	for _, detail := range purchase.Details {
		var insertedDetail models.PurchaseDetail
		if err := tx.QueryRow(`
			INSERT INTO purchase_details (purchase_id, product_id, quantity, price, total)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, purchase_id, product_id, quantity, price, total
		`,
			purchase.ID,
			detail.ProductID,
			detail.Quantity,
			detail.Price,
			detail.Total,
		).Scan(
			&insertedDetail.ID,
			&insertedDetail.PurchaseID,
			&insertedDetail.ProductID,
			&insertedDetail.Quantity,
			&insertedDetail.Price,
			&insertedDetail.Total,
		); err != nil {
			return models.Purchase{}, err
		}

		if _, err := tx.Exec(`
			INSERT INTO stocks (product_id, outlet_id, qty)
			VALUES ($1, $2, $3)
			ON CONFLICT (product_id, outlet_id)
			DO UPDATE SET qty = stocks.qty + EXCLUDED.qty
		`, insertedDetail.ProductID, purchase.OutletID, insertedDetail.Quantity); err != nil {
			return models.Purchase{}, err
		}

		if _, err := tx.Exec(`
			INSERT INTO stock_movements (product_id, outlet_id, type, qty, reference_type, reference_id, note, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`,
			insertedDetail.ProductID,
			purchase.OutletID,
			"IN",
			insertedDetail.Quantity,
			"purchase",
			purchase.ID,
			"purchase stock in",
			purchase.CreatedAt,
		); err != nil {
			return models.Purchase{}, err
		}

		details = append(details, insertedDetail)
	}
	purchase.Details = details

	if err := tx.Commit(); err != nil {
		return models.Purchase{}, err
	}

	return purchase, nil
}
