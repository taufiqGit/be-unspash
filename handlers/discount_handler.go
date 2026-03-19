package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"gowes/models"
	"gowes/services"
	"gowes/utils"
)

type DiscountHandler struct {
	service services.DiscountService
}

func NewDiscountHandler(service services.DiscountService) *DiscountHandler {
	return &DiscountHandler{service: service}
}

func (h *DiscountHandler) ListOrCreate(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(UserContextKey).(models.User)
	if !ok || user.CompanyID == nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "company info missing")
		return
	}

	switch r.Method {
	case http.MethodGet:
		params := utils.ParsePaginationParams(r)
		discounts, total, err := h.service.ListDiscounts(*user.CompanyID, params)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list discounts")
			return
		}
		meta := utils.CalculateMeta(total, params)
		writeSuccess(w, http.StatusOK, discounts, "discount list", meta)

	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "failed to read body")
			return
		}
		var input models.DiscountInput
		if err := json.Unmarshal(body, &input); err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON format")
			return
		}

		if strings.TrimSpace(input.Name) == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "name cannot be empty")
			return
		}
		if input.Type == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "type cannot be empty")
			return
		}
		if input.DiscountValue <= 0 {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "discount_value must be greater than 0")
			return
		}

		discount, err := h.service.CreateDiscount(*user.CompanyID, input)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
			return
		}
		writeSuccess(w, http.StatusCreated, discount, "discount created", nil)

	default:
		w.Header().Set("Allow", "GET, POST")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	}
}

func (h *DiscountHandler) HandleByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "id cannot be empty")
		return
	}

	switch r.Method {
	case http.MethodGet:
		discount, err := h.service.GetDiscount(id)
		if err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "discount not found")
			return
		}
		writeSuccess(w, http.StatusOK, discount, "discount detail", nil)

	case http.MethodPut:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "failed to read body")
			return
		}
		var input models.DiscountInput
		if err := json.Unmarshal(body, &input); err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON format")
			return
		}

		if strings.TrimSpace(input.Name) == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "name cannot be empty")
			return
		}
		if input.Type == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "type cannot be empty")
			return
		}
		if input.DiscountValue <= 0 {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "discount_value must be greater than 0")
			return
		}

		updated, err := h.service.UpdateDiscount(id, input)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
			return
		}
		writeSuccess(w, http.StatusOK, updated, "discount updated", nil)

	case http.MethodDelete:
		if err := h.service.DeleteDiscount(id); err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "discount not found")
			return
		}
		writeSuccess(w, http.StatusOK, nil, "discount deleted", nil)

	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	}
}
