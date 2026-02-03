package repositories

import (
	"database/sql"
	"gowes/models"
	"log"
)

type CategoryRepository interface {
	FindAll() ([]models.Category, error)
	FindByID(id int) (models.Category, error)
	Create(category models.Category) (models.Category, error)
	Update(category models.Category) (models.Category, error)
	Delete(id int) error
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindAll() ([]models.Category, error) {
	rows, err := r.db.Query(`
		SELECT id, name, description, created_at, updated_at 
		FROM categories 
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			log.Println("Error scanning category:", err)
			continue
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *categoryRepository) FindByID(id int) (models.Category, error) {
	row := r.db.QueryRow(`
		SELECT id, name, description, created_at, updated_at 
		FROM categories 
		WHERE id = $1
	`, id)

	var c models.Category
	if err := row.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return models.Category{}, err
	}
	return c, nil
}

func (r *categoryRepository) Create(category models.Category) (models.Category, error) {
	err := r.db.QueryRow(`
		INSERT INTO categories (name, description, created_at, updated_at) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id
	`, category.Name, category.Description, category.CreatedAt, category.UpdatedAt).Scan(&category.ID)

	if err != nil {
		return models.Category{}, err
	}
	return category, nil
}

func (r *categoryRepository) Update(category models.Category) (models.Category, error) {
	err := r.db.QueryRow(`
		UPDATE categories 
		SET name = $1, description = $2, updated_at = $3 
		WHERE id = $4 
		RETURNING id
	`, category.Name, category.Description, category.UpdatedAt, category.ID).Scan(&category.ID)

	if err != nil {
		return models.Category{}, err
	}
	return category, nil
}

func (r *categoryRepository) Delete(id int) error {
	res, err := r.db.Exec(`DELETE FROM categories WHERE id = $1`, id)
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
