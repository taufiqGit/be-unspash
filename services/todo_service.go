package services

import (
	"log"
	"time"

	"gowes/db"
	"gowes/models"
)

func ListTodos() []models.Todo {
	rows, err := db.DB.Query(`
		SELECT id, title, done, image_url, created_at, updated_at 
		FROM todos 
		ORDER BY id ASC
	`)
	if err != nil {
		log.Println("Error listing todos:", err)
		return []models.Todo{}
	}
	defer rows.Close()

	todos := []models.Todo{}
	for rows.Next() {
		var t models.Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.ImageURL, &t.CreatedAt, &t.UpdatedAt); err != nil {
			log.Println("Error scanning todo:", err)
			continue
		}
		todos = append(todos, t)
	}
	return todos
}

func GetTodo(id int) (models.Todo, bool) {
	row := db.DB.QueryRow(`
		SELECT id, title, done, image_url, created_at, updated_at 
		FROM todos 
		WHERE id = $1
	`, id)

	var t models.Todo
	if err := row.Scan(&t.ID, &t.Title, &t.Done, &t.ImageURL, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return models.Todo{}, false
	}
	return t, true
}

func CreateTodo(in models.TodoInput) models.Todo {
	var t models.Todo
	err := db.DB.QueryRow(`
		INSERT INTO todos (title, done, image_url, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, title, done, image_url, created_at, updated_at
	`, in.Title, in.Done, in.ImageURL, time.Now().UTC(), time.Now().UTC()).Scan(&t.ID, &t.Title, &t.Done, &t.ImageURL, &t.CreatedAt, &t.UpdatedAt)

	if err != nil {
		log.Println("Error creating todo:", err)
		return models.Todo{}
	}
	return t
}

func UpdateTodo(id int, in models.TodoInput) (models.Todo, bool) {
	var t models.Todo
	err := db.DB.QueryRow(`
		UPDATE todos 
		SET title = $1, done = $2, image_url = $3, updated_at = $4 
		WHERE id = $5 
		RETURNING id, title, done, image_url, created_at, updated_at
	`, in.Title, in.Done, in.ImageURL, time.Now().UTC(), id).Scan(&t.ID, &t.Title, &t.Done, &t.ImageURL, &t.CreatedAt, &t.UpdatedAt)

	if err != nil {
		log.Println("Error updating todo:", err)
		return models.Todo{}, false
	}
	return t, true
}

func DeleteTodo(id int) bool {
	res, err := db.DB.Exec(`DELETE FROM todos WHERE id = $1`, id)
	if err != nil {
		log.Println("Error deleting todo:", err)
		return false
	}
	count, err := res.RowsAffected()
	if err != nil {
		return false
	}
	return count > 0
}
