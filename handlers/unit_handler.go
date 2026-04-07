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

type UnitHandler struct {
	service services.UnitService
}

func NewUnitHandler(service services.UnitService) *UnitHandler {
	return &UnitHandler{service: service}
}

func (h *UnitHandler) ListOrCreate(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(UserContextKey).(models.User)
	if !ok || user.ID == "" || user.CompanyID == nil || *user.CompanyID == "" {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "user or company info missing")
		return
	}

	switch r.Method {
	case http.MethodGet:
		params := utils.ParsePaginationParams(r)
		units, total, err := h.service.ListUnits(*user.CompanyID, params)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list units")
			return
		}
		meta := utils.CalculateMeta(total, params)
		writeSuccess(w, http.StatusOK, units, "unit list", meta)
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "failed to read body")
			return
		}

		var input models.UnitInput
		if err := json.Unmarshal(body, &input); err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON format")
			return
		}

		if strings.TrimSpace(input.Name) == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "name cannot be empty")
			return
		}
		if strings.TrimSpace(input.Symbol) == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "symbol cannot be empty")
			return
		}
		if strings.TrimSpace(input.Type) == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "type cannot be empty")
			return
		}

		unit, err := h.service.CreateUnit(*user.CompanyID, user.ID, input)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
			return
		}

		writeSuccess(w, http.StatusCreated, unit, "unit created", nil)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	}
}

func (h *UnitHandler) HandleByID(w http.ResponseWriter, r *http.Request) {
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
		unit, err := h.service.GetUnit(id, *user.CompanyID)
		if err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "unit not found")
			return
		}
		writeSuccess(w, http.StatusOK, unit, "unit detail", nil)
	case http.MethodPut:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "failed to read body")
			return
		}

		var input models.UnitInput
		if err := json.Unmarshal(body, &input); err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON format")
			return
		}

		if strings.TrimSpace(input.Name) == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "name cannot be empty")
			return
		}
		if strings.TrimSpace(input.Symbol) == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "symbol cannot be empty")
			return
		}
		if strings.TrimSpace(input.Type) == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "type cannot be empty")
			return
		}

		unit, err := h.service.UpdateUnit(id, *user.CompanyID, user.ID, input)
		if err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "unit not found")
			return
		}
		writeSuccess(w, http.StatusOK, unit, "unit updated", nil)
	case http.MethodDelete:
		if err := h.service.DeleteUnit(id, *user.CompanyID); err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "unit not found")
			return
		}
		writeSuccess(w, http.StatusOK, nil, "unit deleted", nil)
	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	}
}
