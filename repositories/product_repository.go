package repositories

import (
	"database/sql"
	"fmt"
	"gowes/models"
)

type ProductRepository interface {
	FindAll(companyID string, params models.PaginationParams) ([]models.ProductList, int, error)
	Create(companyID string, payload models.ProductInput) (models.Product, error)
	FindByID(productID string) (models.Product, error)
	Update(productID string, payload models.ProductInput) (models.Product, error)
	DeleteById(productID string) (string, error)
	UpdateAddOnsByProductID(addOnIDs []string, productID string, companyID string) ([]models.AddOnProduct, error)
	FindAllMobile(companyID string) ([]models.ProductList, error)
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindAll(companyID string, params models.PaginationParams) ([]models.ProductList, int, error) {
	baseQuery := `
		FROM products p
		JOIN categories c ON p.category_id = c.id
		WHERE p.company_id = $1
	`
	args := []interface{}{companyID}
	argIndex := 2

	if params.Search != "" {
		baseQuery += fmt.Sprintf(" AND (p.name ILIKE $%d OR p.sku ILIKE $%d)", argIndex, argIndex+1)
		args = append(args, "%"+params.Search+"%", "%"+params.Search+"%")
		argIndex += 2
	}

	var total int
	countQuery := "SELECT COUNT(*) " + baseQuery
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	allowedSorts := map[string]string{"name": "p.name", "created_at": "p.created_at", "updated_at": "p.updated_at"}
	sortedBy := "p.created_at"
	if col, ok := allowedSorts[params.SortBy]; ok {
		sortedBy = col
	}
	query := "SELECT p.id, p.name, p.sku, p.unit, p.unit_id, p.cost, p.price, p.image_url, c.name AS category " + baseQuery
	query += fmt.Sprintf(" ORDER BY %s %s", sortedBy, params.SortOrder)

	offset := (params.Page - 1) * params.Limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, params.Limit, offset)
	fmt.Println(query, params.Page, params.Limit, offset, "min")
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products = []models.ProductList{}
	for rows.Next() {
		var product models.ProductList
		var unitID sql.NullString
		var imageURL sql.NullString
		var category sql.NullString
		if err := rows.Scan(&product.ID, &product.Name, &product.SKU, &product.Unit, &unitID, &product.Cost, &product.Price, &imageURL, &category); err != nil {
			return nil, 0, err
		}
		if unitID.Valid {
			product.UnitID = unitID.String
		} else {
			product.UnitID = ""
		}
		if imageURL.Valid {
			product.ImageURL = imageURL.String
		} else {
			product.ImageURL = ""
		}
		if category.Valid {
			product.Category = category.String
		} else {
			product.Category = ""
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	fmt.Println(args, "args")
	return products, total, nil
}

func (r *productRepository) Create(companyID string, payload models.ProductInput) (models.Product, error) {
	var createProduct models.Product
	queryInsert := "INSERT INTO products (name, sku, unit, unit_id, cost, price, image_url, company_id, category_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, name, sku, unit, unit_id, cost, price, image_url, company_id, category_id, created_at, updated_at"
	args := []interface{}{payload.Name, payload.SKU, payload.Unit, payload.UnitID, payload.Cost, payload.Price, payload.ImageURL, companyID, payload.CategoryID}
	var row *sql.Row
	row = r.db.QueryRow(queryInsert, args...)
	if err := row.Scan(&createProduct.ID, &createProduct.Name, &createProduct.SKU, &createProduct.Unit, &createProduct.UnitID, &createProduct.Cost, &createProduct.Price, &createProduct.ImageURL, &createProduct.CompanyID, &createProduct.CategoryID, &createProduct.CreatedAt, &createProduct.UpdatedAt); err != nil {
		return models.Product{}, err
	}
	return createProduct, nil
}

func (r *productRepository) FindByID(productID string) (models.Product, error) {
	var product models.Product
	query := "SELECT id, name, sku, unit, unit_id, cost, price, image_url, company_id, category_id, created_at, updated_at FROM products WHERE id = $1"
	row := r.db.QueryRow(query, productID)
	if err := row.Scan(&product.ID, &product.Name, &product.SKU, &product.Unit, &product.UnitID, &product.Cost, &product.Price, &product.ImageURL, &product.CompanyID, &product.CategoryID, &product.CreatedAt, &product.UpdatedAt); err != nil {
		return models.Product{}, err
	}
	return product, nil
}

func (r *productRepository) Update(productID string, payload models.ProductInput) (models.Product, error) {
	var product models.Product
	query := `
		UPDATE products
		SET name = $1, sku = $2, unit = $3, unit_id = $4, cost = $5, price = $6, image_url = $7, category_id = $8, updated_at = NOW()
		WHERE id = $9
		RETURNING id, name, sku, unit, unit_id, cost, price, image_url, company_id, category_id, created_at, updated_at
	`
	err := r.db.QueryRow(query,
		payload.Name,
		payload.SKU,
		payload.Unit,
		payload.UnitID,
		payload.Cost,
		payload.Price,
		payload.ImageURL,
		payload.CategoryID,
		productID,
	).Scan(
		&product.ID, &product.Name, &product.SKU, &product.Unit, &product.UnitID,
		&product.Cost, &product.Price, &product.ImageURL,
		&product.CompanyID, &product.CategoryID,
		&product.CreatedAt, &product.UpdatedAt,
	)
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}

func (r *productRepository) DeleteById(productID string) (string, error) {
	var imageURL string
	query := "DELETE FROM products WHERE id = $1 RETURNING COALESCE(image_url, '')"
	err := r.db.QueryRow(query, productID).Scan(&imageURL)
	if err != nil {
		return "", err
	}
	return imageURL, nil
}

func (r *productRepository) UpdateAddOnsByProductID(addOnIDs []string, productID string, companyID string) ([]models.AddOnProduct, error) {
	var addOns []models.AddOnProduct
	queryDelete := "DELETE FROM add_on_products WHERE product_id = $1 AND company_id = $2"
	_, err := r.db.Exec(queryDelete, productID, companyID)
	if err != nil {
		return nil, err
	}
	queryInsert := "INSERT INTO add_on_products (add_on_id, product_id, company_id) VALUES ($1, $2, $3)"
	for _, addOnID := range addOnIDs {
		_, err := r.db.Exec(queryInsert, addOnID, productID, companyID)
		if err != nil {
			return nil, err
		}
	}
	return addOns, nil
}

func (r *productRepository) FindAllMobile(companyID string) ([]models.ProductList, error) {
	query := `
		SELECT id, name, sku, unit, unit_id, cost, price, image_url, category_id
		FROM products
		WHERE company_id = $1
	`
	rows, err := r.db.Query(query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products = []models.ProductList{}
	for rows.Next() {
		var product models.ProductList
		var unitID sql.NullString
		var imageURL sql.NullString
		var category sql.NullString
		if err := rows.Scan(&product.ID, &product.Name, &product.SKU, &product.Unit, &unitID, &product.Cost, &product.Price, &imageURL, &category); err != nil {
			return nil, err
		}
		if unitID.Valid {
			product.UnitID = unitID.String
		} else {
			product.UnitID = ""
		}
		if imageURL.Valid {
			product.ImageURL = imageURL.String
		} else {
			product.ImageURL = ""
		}
		if category.Valid {
			product.Category = category.String
		} else {
			product.Category = ""
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
