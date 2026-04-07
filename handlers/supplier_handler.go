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

type SupplierHandler struct {
	service services.SupplierService
}

func NewSupplierHandler(service services.SupplierService) *SupplierHandler {
	return &SupplierHandler{service: service}
}

func (h *SupplierHandler) ListOrCreate(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(UserContextKey).(models.User)
	if !ok || user.ID == "" || user.CompanyID == nil || *user.CompanyID == "" {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "user or company info missing")
		return
	}

	switch r.Method {
	case http.MethodGet:
		params := utils.ParsePaginationParams(r)
		suppliers, total, err := h.service.ListSuppliers(*user.CompanyID, params)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list suppliers")
			return
		}
		meta := utils.CalculateMeta(total, params)
		writeSuccess(w, http.StatusOK, suppliers, "supplier list", meta)
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "failed to read body")
			return
		}
		var input models.SupplierInput
		if err := json.Unmarshal(body, &input); err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON format")
			return
		}
		if strings.TrimSpace(input.Name) == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "name cannot be empty")
			return
		}

		supplier, err := h.service.CreateSupplier(*user.CompanyID, user.ID, input)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create supplier")
			return
		}
		writeSuccess(w, http.StatusCreated, supplier, "supplier created", nil)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	}
}

func (h *SupplierHandler) HandleByID(w http.ResponseWriter, r *http.Request) {
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
		supplier, err := h.service.GetSupplier(id, *user.CompanyID)
		if err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "supplier not found")
			return
		}
		writeSuccess(w, http.StatusOK, supplier, "supplier detail", nil)
	case http.MethodPut:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "failed to read body")
			return
		}
		var input models.SupplierInput
		if err := json.Unmarshal(body, &input); err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON format")
			return
		}
		if strings.TrimSpace(input.Name) == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "name cannot be empty")
			return
		}

		supplier, err := h.service.UpdateSupplier(id, *user.CompanyID, user.ID, input)
		if err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "supplier not found")
			return
		}
		writeSuccess(w, http.StatusOK, supplier, "supplier updated", nil)
	case http.MethodDelete:
		if err := h.service.DeleteSupplier(id, *user.CompanyID); err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "supplier not found")
			return
		}
		writeSuccess(w, http.StatusOK, nil, "supplier deleted", nil)
	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	}
}
