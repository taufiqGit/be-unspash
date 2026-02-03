package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"gowes/models"
	"gowes/services"
)

// CategoryController handles HTTP requests for categories
type CategoryController struct {
	service services.CategoryService
}

// NewCategoryController creates a new CategoryController
func NewCategoryController(service services.CategoryService) *CategoryController {
	return &CategoryController{service: service}
}

// ListOrCreate handles /api/categories (GET list, POST create)
func (c *CategoryController) ListOrCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		categories, err := c.service.ListCategories()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list categories")
			return
		}
		meta := map[string]any{"count": len(categories)}
		writeSuccess(w, http.StatusOK, categories, "list categories iki", meta)
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "failed to read body")
			return
		}
		var in models.CategoryInput
		if err := json.Unmarshal(body, &in); err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON format")
			return
		}
		if strings.TrimSpace(in.Name) == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "name cannot be empty")
			return
		}
		created, err := c.service.CreateCategory(in)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create category")
			return
		}
		writeSuccess(w, http.StatusCreated, created, "category created", nil)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	}
}

// HandleByID handles /api/categories/{id} (GET detail, PUT update, DELETE delete)
func (c *CategoryController) HandleByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "invalid ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		category, err := c.service.GetCategory(id)
		if err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "category not found")
			return
		}
		writeSuccess(w, http.StatusOK, category, "category detail", nil)
	case http.MethodPut:
		var in models.CategoryInput
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON format")
			return
		}
		if strings.TrimSpace(in.Name) == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "name cannot be empty")
			return
		}
		updated, err := c.service.UpdateCategory(id, in)
		if err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "category not found")
			return
		}
		writeSuccess(w, http.StatusOK, updated, "category updated", nil)
	case http.MethodDelete:
		if err := c.service.DeleteCategory(id); err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "category not found")
			return
		}
		writeSuccess(w, http.StatusOK, map[string]int{"id": id}, "category deleted", nil)
	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	}
}
