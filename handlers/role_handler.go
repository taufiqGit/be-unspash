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

type RoleHandler struct {
	service services.RoleService
}

func NewRoleHandler(service services.RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

func (h *RoleHandler) ListOrCreate(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(UserContextKey).(models.User)
	if !ok || user.ID == "" || user.CompanyID == nil || *user.CompanyID == "" {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "user or company info missing")
		return
	}

	switch r.Method {
	case http.MethodGet:
		params := utils.ParsePaginationParams(r)
		roles, total, err := h.service.ListRoles(*user.CompanyID, params)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list roles")
			return
		}
		meta := utils.CalculateMeta(total, params)
		writeSuccess(w, http.StatusOK, roles, "role list", meta)
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "failed to read body")
			return
		}

		var input models.RoleInput
		if err := json.Unmarshal(body, &input); err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON format")
			return
		}

		if strings.TrimSpace(input.Name) == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "name cannot be empty")
			return
		}

		role, err := h.service.CreateRole(*user.CompanyID, user.ID, input)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
			return
		}

		writeSuccess(w, http.StatusCreated, role, "role created", nil)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	}
}

func (h *RoleHandler) HandleByID(w http.ResponseWriter, r *http.Request) {
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
		role, err := h.service.GetRole(id, *user.CompanyID)
		if err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "role not found")
			return
		}
		writeSuccess(w, http.StatusOK, role, "role detail", nil)
	case http.MethodPut:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "failed to read body")
			return
		}

		var input models.RoleInput
		if err := json.Unmarshal(body, &input); err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON format")
			return
		}

		if strings.TrimSpace(input.Name) == "" {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "name cannot be empty")
			return
		}

		role, err := h.service.UpdateRole(id, *user.CompanyID, user.ID, input)
		if err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "role not found")
			return
		}
		writeSuccess(w, http.StatusOK, role, "role updated", nil)
	case http.MethodDelete:
		if err := h.service.DeleteRole(id, *user.CompanyID); err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "role not found")
			return
		}
		writeSuccess(w, http.StatusOK, nil, "role deleted", nil)
	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	}
}
