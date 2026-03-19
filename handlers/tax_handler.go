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

type TaxHandler struct {
	service services.TaxService
}

func NewTaxHandler(service services.TaxService) *TaxHandler {
	return &TaxHandler{service: service}
}

func (h *TaxHandler) ListOrCreate(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(UserContextKey).(models.User)
	if !ok || user.CompanyID == nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "company info missing")
		return
	}

	switch r.Method {
	case http.MethodGet:
		params := utils.ParsePaginationParams(r)
		taxes, total, err := h.service.ListTaxes(*user.CompanyID, params)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list taxes")
			return
		}
		meta := utils.CalculateMeta(total, params)
		writeSuccess(w, http.StatusOK, taxes, "tax list", meta)

	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "failed to read body")
			return
		}
		var input models.TaxInput
		if err := json.Unmarshal(body, &input); err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON format")
			return
		}

		if strings.TrimSpace(input.Name) == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "name cannot be empty")
			return
		}
		if input.Rate < 0 {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "rate cannot be negative")
			return
		}
		if input.Rate > 100 {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "rate cannot exceed 100%")
			return
		}

		tax, err := h.service.CreateTax(*user.CompanyID, input)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create tax")
			return
		}
		writeSuccess(w, http.StatusCreated, tax, "tax created", nil)

	default:
		w.Header().Set("Allow", "GET, POST")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	}
}

func (h *TaxHandler) HandleByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "id cannot be empty")
		return
	}

	switch r.Method {
	case http.MethodGet:
		tax, err := h.service.GetTax(id)
		if err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "tax not found")
			return
		}
		writeSuccess(w, http.StatusOK, tax, "tax detail", nil)

	case http.MethodPut:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "failed to read body")
			return
		}
		var input models.TaxInput
		if err := json.Unmarshal(body, &input); err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON format")
			return
		}

		if strings.TrimSpace(input.Name) == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "name cannot be empty")
			return
		}
		if input.Rate < 0 {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "rate cannot be negative")
			return
		}
		if input.Rate > 100 {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "rate cannot exceed 100%")
			return
		}

		updated, err := h.service.UpdateTax(id, input)
		if err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "tax not found")
			return
		}
		writeSuccess(w, http.StatusOK, updated, "tax updated", nil)

	case http.MethodDelete:
		if err := h.service.DeleteTax(id); err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "tax not found")
			return
		}
		writeSuccess(w, http.StatusOK, nil, "tax deleted", nil)

	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	}
}
