package handlers

import (
	"encoding/json"
	"gowes/models"
	"gowes/services"
	"gowes/utils"
	"io"
	"net/http"
	"strings"
)

type RecipeHandler struct {
	service services.RecipeService
}

func NewRecipeHandler(service services.RecipeService) *RecipeHandler {
	return &RecipeHandler{service: service}
}

func (h *RecipeHandler) ListOrCreate(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(UserContextKey).(models.User)
	if !ok || user.ID == "" || user.CompanyID == nil || *user.CompanyID == "" {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "user or company info missing")
		return
	}

	switch r.Method {
	case http.MethodGet:
		params := utils.ParsePaginationParams(r)
		recipes, total, err := h.service.ListRecipes(*user.CompanyID, params)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list recipes")
			return
		}
		meta := utils.CalculateMeta(total, params)
		writeSuccess(w, http.StatusOK, recipes, "recipe list", meta)
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "failed to read body")
			return
		}
		var input models.RecipeInput
		if err := json.Unmarshal(body, &input); err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON format")
			return
		}
		if strings.TrimSpace(input.ProductID) == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "product_id cannot be empty")
			return
		}
		if strings.TrimSpace(input.IngredientID) == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "ingredient_id cannot be empty")
			return
		}

		recipe, err := h.service.CreateRecipe(*user.CompanyID, user.ID, input)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create recipe")
			return
		}
		writeSuccess(w, http.StatusCreated, recipe, "recipe created", nil)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	}
}

func (h *RecipeHandler) HandleByID(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(UserContextKey).(models.User)
	if !ok || user.ID == "" || user.CompanyID == nil || *user.CompanyID == "" {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "user or company info missing")
		return
	}

	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "id cannot be empty")
		return
	}

	switch r.Method {
	case http.MethodGet:
		recipe, err := h.service.GetRecipe(id, *user.CompanyID)
		if err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "recipe not found")
			return
		}
		writeSuccess(w, http.StatusOK, recipe, "recipe detail", nil)
	case http.MethodPut:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "failed to read body")
			return
		}
		var input models.RecipeInput
		if err := json.Unmarshal(body, &input); err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON format")
			return
		}
		if strings.TrimSpace(input.ProductID) == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "product_id cannot be empty")
			return
		}
		if strings.TrimSpace(input.IngredientID) == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "ingredient_id cannot be empty")
			return
		}

		recipe, err := h.service.UpdateRecipe(id, *user.CompanyID, user.ID, input)
		if err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "recipe not found")
			return
		}
		writeSuccess(w, http.StatusOK, recipe, "recipe updated", nil)
	case http.MethodDelete:
		if err := h.service.DeleteRecipe(id, *user.CompanyID); err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "recipe not found")
			return
		}
		writeSuccess(w, http.StatusOK, nil, "recipe deleted", nil)
	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	}
}
