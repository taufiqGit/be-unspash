package services

import (
	"log"
	"time"

	"gowes/db"
	"gowes/models"
)

func ListCategories() []models.Category {
	rows, err := db.DB.Query(`
		SELECT id, name, description, created_at, updated_at 
		FROM categories 
		ORDER BY id ASC
	`)
	if err != nil {
		log.Println("Error listing categories:", err)
		return []models.Category{}
	}
	defer rows.Close()

	categories := []models.Category{}
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			log.Println("Error scanning category:", err)
			continue
		}
		categories = append(categories, c)
	}
	return categories
}

func GetCategory(id int) (models.Category, bool) {
	row := db.DB.QueryRow(`
		SELECT id, name, description, created_at, updated_at 
		FROM categories 
		WHERE id = $1
	`, id)

	var c models.Category
	if err := row.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return models.Category{}, false
	}
	return c, true
}

func CreateCategory(in models.CategoryInput) models.Category {
	var c models.Category
	err := db.DB.QueryRow(`
		INSERT INTO categories (name, description, created_at, updated_at) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, name, description, created_at, updated_at
	`, in.Name, in.Description, time.Now().UTC(), time.Now().UTC()).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)

	if err != nil {
		log.Println("Error creating category:", err)
		return models.Category{}
	}
	return c
}

func UpdateCategory(id int, in models.CategoryInput) (models.Category, bool) {
	var c models.Category
	err := db.DB.QueryRow(`
		UPDATE categories 
		SET name = $1, description = $2, updated_at = $3 
		WHERE id = $4 
		RETURNING id, name, description, created_at, updated_at
	`, in.Name, in.Description, time.Now().UTC(), id).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)

	if err != nil {
		log.Println("Error updating category:", err)
		return models.Category{}, false
	}
	return c, true
}

func DeleteCategory(id int) bool {
	res, err := db.DB.Exec(`DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		log.Println("Error deleting category:", err)
		return false
	}
	count, err := res.RowsAffected()
	if err != nil {
		return false
	}
	return count > 0
}
