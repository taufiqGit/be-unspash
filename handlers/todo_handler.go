package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"gowes/models"
	"gowes/services"
)

// TodoHandler handles HTTP requests for todos
type TodoHandler struct {
	service services.TodoService
}

// NewTodoHandler creates a new TodoHandler
func NewTodoHandler(service services.TodoService) *TodoHandler {
	return &TodoHandler{service: service}
}

// ListOrCreate handles /api/todos (GET for list, POST for create)
func (c *TodoHandler) ListOrCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		todos, err := c.service.ListTodos()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "gagal mengambil data todos")
			return
		}
		meta := map[string]any{"count": len(todos)}
		writeSuccess(w, http.StatusOK, todos, "list data todos aseli", meta)
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "gagal membaca bodies")
			return
		}
		var in models.TodoInput
		if err := json.Unmarshal(body, &in); err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "format JSON tidak valid!!")
			return
		}
		if strings.TrimSpace(in.Title) == "" || strings.TrimSpace(in.ImageURL) == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "titleies dan image url tidak boleh kosonggnn")
			return
		}
		created, err := c.service.CreateTodo(in)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "gagal membuat todo")
			return
		}
		writeSuccess(w, http.StatusCreated, created, "todo created", nil)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method tidak diizinkan")
	}
}

// HandleByID handles /api/todos/{id} (GET, PUT, DELETE)
func (c *TodoHandler) HandleByID(w http.ResponseWriter, r *http.Request) {
	// Ekstrak ID dari path
	idStr := strings.TrimPrefix(r.URL.Path, "/api/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "ID tidak valid")
		return
	}

	switch r.Method {
	case http.MethodGet:
		todo, err := c.service.GetTodo(id)
		if err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "todo tidak ditemukan")
			return
		}
		writeSuccess(w, http.StatusOK, todo, "todo detail", nil)
	case http.MethodPut:
		var in models.TodoInput
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "format JSON tidak valid")
			return
		}
		if strings.TrimSpace(in.Title) == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "title tidak boleh kosong")
			return
		}
		updated, err := c.service.UpdateTodo(id, in)
		if err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "todo tidak ditemukan")
			return
		}
		writeSuccess(w, http.StatusOK, updated, "todo updated", nil)
	case http.MethodDelete:
		if err := c.service.DeleteTodo(id); err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "todo tidak ditemukan")
			return
		}
		writeSuccess(w, http.StatusOK, map[string]int{"id": id}, "todo deleted", nil)
	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method tidak diizinkan")
	}
}
