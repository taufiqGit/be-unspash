package services

import (
	"gowes/models"
	"gowes/repositories"
	"time"
)

type TodoService interface {
	ListTodos() ([]models.Todo, error)
	GetTodo(id int) (models.Todo, error)
	CreateTodo(in models.TodoInput) (models.Todo, error)
	UpdateTodo(id int, in models.TodoInput) (models.Todo, error)
	DeleteTodo(id int) error
}

type todoService struct {
	repo repositories.TodoRepository
}

func NewTodoService(repo repositories.TodoRepository) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) ListTodos() ([]models.Todo, error) {
	return s.repo.FindAll()
}

func (s *todoService) GetTodo(id int) (models.Todo, error) {
	return s.repo.FindByID(id)
}

func (s *todoService) CreateTodo(in models.TodoInput) (models.Todo, error) {
	t := models.Todo{
		Title:     in.Title,
		Done:      in.Done,
		ImageURL:  in.ImageURL,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	return s.repo.Create(t)
}

func (s *todoService) UpdateTodo(id int, in models.TodoInput) (models.Todo, error) {
	// Check if exists
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return models.Todo{}, err
	}

	existing.Title = in.Title
	existing.Done = in.Done
	existing.ImageURL = in.ImageURL
	existing.UpdatedAt = time.Now().UTC()

	return s.repo.Update(existing)
}

func (s *todoService) DeleteTodo(id int) error {
	return s.repo.Delete(id)
}
