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

type CustomerHandler struct {
	service services.CustomerService
}

func NewCustomerHandler(service services.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

func (h *CustomerHandler) ListOrCreate(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(UserContextKey).(models.User)
	if !ok || user.CompanyID == nil {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "company info missing")
		return
	}

	switch r.Method {
	case http.MethodGet:
		params := utils.ParsePaginationParams(r)
		customers, total, err := h.service.ListCustomers(*user.CompanyID, params)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list customers")
			return
		}
		meta := utils.CalculateMeta(total, params)
		writeSuccess(w, http.StatusOK, customers, "customer list", meta)

	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "failed to read body")
			return
		}
		var input models.CustomerInput
		if err := json.Unmarshal(body, &input); err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON format")
			return
		}
		if strings.TrimSpace(input.Name) == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "name cannot be empty")
			return
		}

		customer, err := h.service.CreateCustomer(input, *user.CompanyID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create customer")
			return
		}
		writeSuccess(w, http.StatusCreated, customer, "customer created", nil)

	default:
		w.Header().Set("Allow", "GET, POST")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	}
}

func (h *CustomerHandler) HandleByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "id cannot be empty")
		return
	}

	switch r.Method {
	case http.MethodGet:
		customer, err := h.service.GetCustomer(id)
		if err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "customer not found")
			return
		}
		writeSuccess(w, http.StatusOK, customer, "customer detail", nil)

	case http.MethodPut:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "failed to read body")
			return
		}
		var input models.CustomerInput
		if err := json.Unmarshal(body, &input); err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON format")
			return
		}
		if strings.TrimSpace(input.Name) == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "name cannot be empty")
			return
		}

		updated, err := h.service.UpdateCustomer(id, input)
		if err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "customer not found")
			return
		}
		writeSuccess(w, http.StatusOK, updated, "customer updated", nil)

	case http.MethodDelete:
		if err := h.service.DeleteCustomer(id); err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "customer not found")
			return
		}
		writeSuccess(w, http.StatusOK, nil, "customer deleted", nil)

	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	}
}
