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

// CategoriesHandler handles /api/categories (GET list, POST create)
func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		categories := services.ListCategories()
		meta := map[string]any{"count": len(categories)}
		writeSuccess(w, http.StatusOK, categories, "list categories", meta)
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
		created := services.CreateCategory(in)
		writeSuccess(w, http.StatusCreated, created, "category created", nil)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	}
}

// CategoryByIDHandler handles /api/categories/{id} (GET detail, PUT update, DELETE delete)
func CategoryByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "invalid ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		category, ok := services.GetCategory(id)
		if !ok {
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
		updated, ok := services.UpdateCategory(id, in)
		if !ok {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "category not found")
			return
		}
		writeSuccess(w, http.StatusOK, updated, "category updated", nil)
	case http.MethodDelete:
		if ok := services.DeleteCategory(id); !ok {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "category not found")
			return
		}
		writeSuccess(w, http.StatusOK, map[string]int{"id": id}, "category deleted", nil)
	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	}
}
