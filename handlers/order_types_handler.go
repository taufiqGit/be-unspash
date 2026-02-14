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

type OrderTypeHandler struct {
	service services.OrderTypeService
}

func NewOrderTypeHandler(service services.OrderTypeService) *OrderTypeHandler {
	return &OrderTypeHandler{service: service}
}

func (h *OrderTypeHandler) ListOrCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "bad_request", "invalid JSON format")
			return
		}
		var in models.OrderTypeInput
		if err := json.Unmarshal(body, &in); err != nil {
			writeError(w, http.StatusBadRequest, "bad_request", "invalid JSON format")
			return
		}
		if strings.TrimSpace(in.Name) == "" {
			writeError(w, http.StatusBadRequest, "validation_error", "name cannot be empty")
			return
		}

		user, ok := r.Context().Value(UserContextKey).(models.User)
		if !ok || user.CompanyID == nil {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "company info missing")
			return
		}

		orderType, err := h.service.Create(*user.CompanyID, in)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
			return
		}
		writeSuccess(w, http.StatusCreated, orderType, "order type created", nil)
	case http.MethodGet:
		user, ok := r.Context().Value(UserContextKey).(models.User)
		if !ok || user.CompanyID == nil {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "company info missing")
			return
		}

		// Get pagination params
		params := utils.ParsePaginationParams(r)

		orderTypes, total, err := h.service.FindAll(*user.CompanyID, params)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
			return
		}

		meta := utils.CalculateMeta(total, params)
		writeSuccess(w, http.StatusOK, orderTypes, "list order types", meta)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
	}
}

func (h *OrderTypeHandler) HandleByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "id cannot be empty")
		return
	}
	switch r.Method {
	case http.MethodGet:
		orderType, err := h.service.FindByID(id)
		if err != nil {
			writeError(w, http.StatusNotFound, "not_found", err.Error())
			return
		}

		writeSuccess(w, http.StatusOK, orderType, "order type detail", nil)
	case http.MethodDelete:
		err := h.service.Delete(id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
			return
		}

		writeSuccess(w, http.StatusOK, nil, "order type deleted", nil)
	case http.MethodPut:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
			return
		}
		var in models.OrderTypeInput
		if err := json.Unmarshal(body, &in); err != nil {
			writeError(w, http.StatusBadRequest, "bad_request", "invalid JSON format")
			return
		}
		if strings.TrimSpace(in.Name) == "" {
			writeError(w, http.StatusBadRequest, "validation_error", "name cannot be empty")
			return
		}

		updated, err := h.service.Update(in, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
			return
		}

		writeSuccess(w, http.StatusOK, updated, "order type updated", nil)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
	}
}
