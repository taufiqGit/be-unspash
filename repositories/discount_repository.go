package repositories

import (
	"database/sql"
	"fmt"
	"gowes/models"
)

type DiscountRepository interface {
	FindAll(companyID string, params models.PaginationParams) ([]models.Discount, int, error)
	FindByID(id string) (models.Discount, error)
	Create(discount models.Discount, outletIDs, categoryIDs, productIDs, orderTypeIDs []string) (models.Discount, error)
	Update(discount models.Discount, outletIDs, categoryIDs, productIDs, orderTypeIDs []string) (models.Discount, error)
	Delete(id string) error
}

type discountRepository struct {
	db *sql.DB
}

func NewDiscountRepository(db *sql.DB) DiscountRepository {
	return &discountRepository{db: db}
}

// ─── scan helpers ─────────────────────────────────────────────────────────────

func scanDiscountFromRow(row *sql.Row) (models.Discount, error) {
	var d models.Discount
	var maxAmount, minPurchase sql.NullFloat64
	var targetType, specificTargetType sql.NullString

	err := row.Scan(
		&d.ID, &d.CompanyID, &d.Name, &d.Type, &d.DiscountValue,
		&maxAmount, &minPurchase,
		&targetType, &specificTargetType,
		&d.ApplyToOrderTypes,
		&d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		return models.Discount{}, err
	}

	if maxAmount.Valid {
		d.MaxAmount = &maxAmount.Float64
	}
	if minPurchase.Valid {
		d.MinPurchase = &minPurchase.Float64
	}
	if targetType.Valid {
		t := models.DiscountTarget(targetType.String)
		d.TargetType = &t
	}
	if specificTargetType.Valid {
		s := models.DiscountSpecificTarget(specificTargetType.String)
		d.SpecificTargetType = &s
	}
	return d, nil
}

func scanDiscountFromRows(rows *sql.Rows) (models.Discount, error) {
	var d models.Discount
	var maxAmount, minPurchase sql.NullFloat64
	var targetType, specificTargetType sql.NullString

	err := rows.Scan(
		&d.ID, &d.CompanyID, &d.Name, &d.Type, &d.DiscountValue,
		&maxAmount, &minPurchase,
		&targetType, &specificTargetType,
		&d.ApplyToOrderTypes,
		&d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		return models.Discount{}, err
	}

	if maxAmount.Valid {
		d.MaxAmount = &maxAmount.Float64
	}
	if minPurchase.Valid {
		d.MinPurchase = &minPurchase.Float64
	}
	if targetType.Valid {
		t := models.DiscountTarget(targetType.String)
		d.TargetType = &t
	}
	if specificTargetType.Valid {
		s := models.DiscountSpecificTarget(specificTargetType.String)
		d.SpecificTargetType = &s
	}
	return d, nil
}

// ─── relation helpers ────────────────────────────────────────────────────────

// loadRelations mengisi semua slice relasi (junction tables) ke dalam Discount.
func (r *discountRepository) loadRelations(d *models.Discount) error {
	// 1. Outlets
	outletRows, err := r.db.Query(
		`SELECT id, discount_id, outlet_id FROM discount_outlets WHERE discount_id = $1`, d.ID,
	)
	if err != nil {
		return fmt.Errorf("load discount_outlets: %w", err)
	}
	defer outletRows.Close()
	d.OutletIDs = []models.DiscountTargetOutlet{}
	for outletRows.Next() {
		var o models.DiscountTargetOutlet
		if err := outletRows.Scan(&o.ID, &o.ParentID, &o.OutletID); err != nil {
			return err
		}
		d.OutletIDs = append(d.OutletIDs, o)
	}

	// 2. Target categories
	catRows, err := r.db.Query(
		`SELECT id, discount_id, category_id FROM discount_target_categories WHERE discount_id = $1`, d.ID,
	)
	if err != nil {
		return fmt.Errorf("load discount_target_categories: %w", err)
	}
	defer catRows.Close()
	d.TargetCategoryIDs = []models.DiscountTargetCategory{}
	for catRows.Next() {
		var c models.DiscountTargetCategory
		if err := catRows.Scan(&c.ID, &c.ParentID, &c.CategoryId); err != nil {
			return err
		}
		d.TargetCategoryIDs = append(d.TargetCategoryIDs, c)
	}

	// 3. Target products
	prodRows, err := r.db.Query(
		`SELECT id, discount_id, product_id FROM discount_target_products WHERE discount_id = $1`, d.ID,
	)
	if err != nil {
		return fmt.Errorf("load discount_target_products: %w", err)
	}
	defer prodRows.Close()
	d.TargetProductIDs = []models.DiscountTargetProduct{}
	for prodRows.Next() {
		var p models.DiscountTargetProduct
		if err := prodRows.Scan(&p.ID, &p.ParentID, &p.ProductId); err != nil {
			return err
		}
		d.TargetProductIDs = append(d.TargetProductIDs, p)
	}

	// 4. Order types
	otRows, err := r.db.Query(
		`SELECT id, discount_id, order_type_id FROM discount_order_types WHERE discount_id = $1`, d.ID,
	)
	if err != nil {
		return fmt.Errorf("load discount_order_types: %w", err)
	}
	defer otRows.Close()
	d.OrderTypeIDs = []models.DiscountTargetOrderType{}
	for otRows.Next() {
		var o models.DiscountTargetOrderType
		if err := otRows.Scan(&o.ID, &o.ParentID, &o.OrderTypeID); err != nil {
			return err
		}
		d.OrderTypeIDs = append(d.OrderTypeIDs, o)
	}

	return nil
}

// insertJunctions menyisipkan semua data junction table dalam satu transaksi.
func insertJunctions(tx *sql.Tx, discountID string, d *models.Discount, outletIDs, categoryIDs, productIDs, orderTypeIDs []string) error {
	// Outlets
	for _, outletID := range outletIDs {
		var o models.DiscountTargetOutlet
		o.OutletID = outletID
		o.ParentID = discountID
		if err := tx.QueryRow(
			`INSERT INTO discount_outlets (discount_id, outlet_id) VALUES ($1, $2) RETURNING id`,
			discountID, outletID,
		).Scan(&o.ID); err != nil {
			return fmt.Errorf("insert discount_outlets: %w", err)
		}
		d.OutletIDs = append(d.OutletIDs, o)
	}

	// Target categories
	for _, categoryID := range categoryIDs {
		var c models.DiscountTargetCategory
		c.CategoryId = categoryID
		c.ParentID = discountID
		if err := tx.QueryRow(
			`INSERT INTO discount_target_categories (discount_id, category_id) VALUES ($1, $2) RETURNING id`,
			discountID, categoryID,
		).Scan(&c.ID); err != nil {
			return fmt.Errorf("insert discount_target_categories: %w", err)
		}
		d.TargetCategoryIDs = append(d.TargetCategoryIDs, c)
	}

	// Target products
	for _, productID := range productIDs {
		var p models.DiscountTargetProduct
		p.ProductId = productID
		p.ParentID = discountID
		if err := tx.QueryRow(
			`INSERT INTO discount_target_products (discount_id, product_id) VALUES ($1, $2) RETURNING id`,
			discountID, productID,
		).Scan(&p.ID); err != nil {
			return fmt.Errorf("insert discount_target_products: %w", err)
		}
		d.TargetProductIDs = append(d.TargetProductIDs, p)
	}

	// Order types
	for _, orderTypeID := range orderTypeIDs {
		var o models.DiscountTargetOrderType
		o.OrderTypeID = orderTypeID
		o.ParentID = discountID
		if err := tx.QueryRow(
			`INSERT INTO discount_order_types (discount_id, order_type_id) VALUES ($1, $2) RETURNING id`,
			discountID, orderTypeID,
		).Scan(&o.ID); err != nil {
			return fmt.Errorf("insert discount_order_types: %w", err)
		}
		d.OrderTypeIDs = append(d.OrderTypeIDs, o)
	}

	return nil
}

// deleteJunctions menghapus semua junction table rows berdasarkan discount_id.
func deleteJunctions(tx *sql.Tx, discountID string) error {
	tables := []string{
		"discount_outlets",
		"discount_target_categories",
		"discount_target_products",
		"discount_order_types",
	}
	for _, table := range tables {
		if _, err := tx.Exec(
			fmt.Sprintf(`DELETE FROM %s WHERE discount_id = $1`, table), discountID,
		); err != nil {
			return fmt.Errorf("delete %s: %w", table, err)
		}
	}
	return nil
}

// ─── FindAll ─────────────────────────────────────────────────────────────────

func (r *discountRepository) FindAll(companyID string, params models.PaginationParams) ([]models.Discount, int, error) {
	baseQuery := " FROM discounts WHERE company_id = $1"
	args := []interface{}{companyID}
	argIdx := 2

	if params.Search != "" {
		baseQuery += fmt.Sprintf(" AND name ILIKE $%d", argIdx)
		args = append(args, "%"+params.Search+"%")
		argIdx++
	}

	// Count total
	var total int
	if err := r.db.QueryRow("SELECT COUNT(*)" + baseQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Validate sort column
	allowedSorts := map[string]bool{"name": true, "created_at": true, "updated_at": true, "type": true}
	sortBy := "created_at"
	if allowedSorts[params.SortBy] {
		sortBy = params.SortBy
	}
	sortOrder := "DESC"
	if params.SortOrder == "ASC" {
		sortOrder = "ASC"
	}

	selectClause := `SELECT id, company_id, name, type, discount_value, max_amount, min_purchase,
		target_type, specific_target_type, apply_to_order_types, created_at, updated_at`

	query := selectClause + baseQuery
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	offset := (params.Page - 1) * params.Limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, params.Limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	discounts := []models.Discount{}
	for rows.Next() {
		d, err := scanDiscountFromRows(rows)
		if err != nil {
			return nil, 0, err
		}
		if err := r.loadRelations(&d); err != nil {
			return nil, 0, err
		}
		discounts = append(discounts, d)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return discounts, total, nil
}

// ─── FindByID ────────────────────────────────────────────────────────────────

func (r *discountRepository) FindByID(id string) (models.Discount, error) {
	row := r.db.QueryRow(`
		SELECT id, company_id, name, type, discount_value, max_amount, min_purchase,
		       target_type, specific_target_type, apply_to_order_types, created_at, updated_at
		FROM discounts
		WHERE id = $1
	`, id)

	d, err := scanDiscountFromRow(row)
	if err != nil {
		return models.Discount{}, err
	}

	if err := r.loadRelations(&d); err != nil {
		return models.Discount{}, err
	}

	return d, nil
}

// ─── Create ──────────────────────────────────────────────────────────────────

func (r *discountRepository) Create(discount models.Discount, outletIDs, categoryIDs, productIDs, orderTypeIDs []string) (models.Discount, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return models.Discount{}, err
	}
	defer tx.Rollback()

	err = tx.QueryRow(`
		INSERT INTO discounts
			(company_id, name, type, discount_value, max_amount, min_purchase,
			 target_type, specific_target_type, apply_to_order_types, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`,
		discount.CompanyID,
		discount.Name,
		discount.Type,
		discount.DiscountValue,
		discount.MaxAmount,
		discount.MinPurchase,
		discount.TargetType,
		discount.SpecificTargetType,
		discount.ApplyToOrderTypes,
		discount.CreatedAt,
		discount.UpdatedAt,
	).Scan(&discount.ID)
	if err != nil {
		return models.Discount{}, fmt.Errorf("insert discount: %w", err)
	}

	// Reset junction slices before populating
	discount.OutletIDs = []models.DiscountTargetOutlet{}
	discount.TargetCategoryIDs = []models.DiscountTargetCategory{}
	discount.TargetProductIDs = []models.DiscountTargetProduct{}
	discount.OrderTypeIDs = []models.DiscountTargetOrderType{}

	if err := insertJunctions(tx, discount.ID, &discount, outletIDs, categoryIDs, productIDs, orderTypeIDs); err != nil {
		return models.Discount{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.Discount{}, err
	}

	return discount, nil
}

// ─── Update ──────────────────────────────────────────────────────────────────

func (r *discountRepository) Update(discount models.Discount, outletIDs, categoryIDs, productIDs, orderTypeIDs []string) (models.Discount, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return models.Discount{}, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE discounts
		SET name = $1, type = $2, discount_value = $3, max_amount = $4, min_purchase = $5,
		    target_type = $6, specific_target_type = $7, apply_to_order_types = $8, updated_at = $9
		WHERE id = $10
	`,
		discount.Name,
		discount.Type,
		discount.DiscountValue,
		discount.MaxAmount,
		discount.MinPurchase,
		discount.TargetType,
		discount.SpecificTargetType,
		discount.ApplyToOrderTypes,
		discount.UpdatedAt,
		discount.ID,
	)
	if err != nil {
		return models.Discount{}, fmt.Errorf("update discount: %w", err)
	}

	// Replace junction table data: delete all then re-insert
	if err := deleteJunctions(tx, discount.ID); err != nil {
		return models.Discount{}, err
	}

	// Reset junction slices before re-populating
	discount.OutletIDs = []models.DiscountTargetOutlet{}
	discount.TargetCategoryIDs = []models.DiscountTargetCategory{}
	discount.TargetProductIDs = []models.DiscountTargetProduct{}
	discount.OrderTypeIDs = []models.DiscountTargetOrderType{}

	if err := insertJunctions(tx, discount.ID, &discount, outletIDs, categoryIDs, productIDs, orderTypeIDs); err != nil {
		return models.Discount{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.Discount{}, err
	}

	return discount, nil
}

// ─── Delete ──────────────────────────────────────────────────────────────────

func (r *discountRepository) Delete(id string) error {
	// Junction tables akan terhapus otomatis via ON DELETE CASCADE
	res, err := r.db.Exec(`DELETE FROM discounts WHERE id = $1`, id)
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
