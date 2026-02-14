package repositories

import (
	"database/sql"
	"gowes/models"
	"log"
)

type TodoRepository interface {
	FindAll() ([]models.Todo, error)
	FindByID(id int) (models.Todo, error)
	Create(todo models.Todo) (models.Todo, error)
	Update(todo models.Todo) (models.Todo, error)
	Delete(id int) error
}

type todoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) FindAll() ([]models.Todo, error) {
	rows, err := r.db.Query(`
		SELECT id, title, done, image_url, created_at, updated_at 
		FROM todos 
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos = []models.Todo{}
	for rows.Next() {
		var t models.Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.ImageURL, &t.CreatedAt, &t.UpdatedAt); err != nil {
			log.Println("Error scanning todo:", err)
			continue
		}
		todos = append(todos, t)
	}
	return todos, nil
}

func (r *todoRepository) FindByID(id int) (models.Todo, error) {
	row := r.db.QueryRow(`
		SELECT id, title, done, image_url, created_at, updated_at 
		FROM todos 
		WHERE id = $1
	`, id)

	var t models.Todo
	if err := row.Scan(&t.ID, &t.Title, &t.Done, &t.ImageURL, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return models.Todo{}, err
	}
	return t, nil
}

func (r *todoRepository) Create(todo models.Todo) (models.Todo, error) {
	err := r.db.QueryRow(`
		INSERT INTO todos (title, done, image_url, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id
	`, todo.Title, todo.Done, todo.ImageURL, todo.CreatedAt, todo.UpdatedAt).Scan(&todo.ID)

	if err != nil {
		return models.Todo{}, err
	}
	return todo, nil
}

func (r *todoRepository) Update(todo models.Todo) (models.Todo, error) {
	err := r.db.QueryRow(`
		UPDATE todos 
		SET title = $1, done = $2, image_url = $3, updated_at = $4 
		WHERE id = $5 
		RETURNING id
	`, todo.Title, todo.Done, todo.ImageURL, todo.UpdatedAt, todo.ID).Scan(&todo.ID)

	if err != nil {
		return models.Todo{}, err
	}
	return todo, nil
}

func (r *todoRepository) Delete(id int) error {
	res, err := r.db.Exec(`DELETE FROM todos WHERE id = $1`, id)
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
