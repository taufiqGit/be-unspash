// package repositories

// import (
// 	"database/sql"
// 	"fmt"
// 	"gowes/models"
// )

// type ProductRepository interface {
// 	FindAll(companyID string, params models.PaginationParams) ([]models.Product, int, error)
// 	Create(companyID string, payload models.ProductInput) (models.Product, error)
// }

// type productRepository struct {
// 	db *sql.DB
// }

// func NewProductRepository(db *sql.DB) ProductRepository {
// 	return &productRepository{db: db}
// }

// func (r *productRepository) FindAll(companyID string, params models.PaginationParams) ([]models.ProductList, int, error) {
// 	query := `
// 		SELECT
// 			p.name,
// 			p.sku,
// 			p.unit,
// 			p.cost,
// 			p.price,
// 			p.image_url,
// 			c.name AS category
// 		FROM products p
// 		JOIN categories c ON p.category_id = c.id
// 		WHERE p.company_id = $1
// 	`
// 	args := []interface{}{companyID}
// 	argIndex := 2

// 	if params.Search != "" {
// 		query += fmt.Sprintf(" AND (p.name ILIKE $%d OR p.sku ILIKE $%d)", argIndex, argIndex+1)
// 		args = append(args, "%"+params.Search+"%", "%"+params.Search+"%")
// 		argIndex += 2
// 	}