package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"gowes/models"
	"gowes/services"
	"gowes/utils"
)

// CategoryHandler handles HTTP requests for categories
type CategoryHandler struct {
	service services.CategoryService
}

// NewCategoryHandler creates a new CategoryHandler
func NewCategoryHandler(service services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// ListOrCreate handles /api/categories (GET list, POST create)
func (c *CategoryHandler) ListOrCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Contoh: Ambil user dari context (dari AuthMiddleware)
		user, ok := r.Context().Value(UserContextKey).(models.User)
		if !ok || user.CompanyID == nil {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "company info missing")
			return
		}

		// Get pagination params
		params := utils.ParsePaginationParams(r)
		fmt.Println(params, "ikih")

		categories, err := c.service.ListCategories(*user.CompanyID, params)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list categories")
			return
		}

		meta := utils.CalculateMeta(len(categories), params)
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

		// Get user from context to extract CompanyID
		user, ok := r.Context().Value(UserContextKey).(models.User)
		if !ok || user.CompanyID == nil {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "company info missing")
			return
		}

		created, err := c.service.CreateCategory(in, *user.CompanyID)
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
func (c *CategoryHandler) HandleByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	id := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	if id == "" {
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
		writeSuccess(w, http.StatusOK, nil, "category deleted", nil)
	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	}
}
