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

type AddOnHandler struct {
	service services.AddOnService
}

func NewAddOnHandler(service services.AddOnService) *AddOnHandler {
	return &AddOnHandler{service: service}
}

func (h *AddOnHandler) ListOrCreateAddOn(w http.ResponseWriter, r *http.Request) {
	// Implementation for listing or creating add-ons
	switch r.Method {
	case http.MethodGet:
		user, ok := r.Context().Value(UserContextKey).(models.User)
		if !ok || user.CompanyID == nil {
			writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "company info missing")
			return
		}

		// Get pagination params
		params := utils.ParsePaginationParams(r)

		addons, total, err := h.service.FindAll(*user.CompanyID, params)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
			return
		}

		meta := utils.CalculateMeta(total, params)
		writeSuccess(w, http.StatusOK, addons, "list add-ons", meta)
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
			return
		}

		var in models.AddOnInput
		if err := json.Unmarshal(body, &in); err != nil {
			writeError(w, http.StatusBadRequest, "bad_request", "invalid JSON format")
			return
		}
		if strings.TrimSpace(in.Name) == "" {
			writeError(w, http.StatusBadRequest, "validation_error", "name cannot be empty")
			return
		}

		// Get user from context to extract CompanyID
		user, ok := r.Context().Value(UserContextKey).(models.User)
		if !ok || user.CompanyID == nil {
			writeError(w, http.StatusUnauthorized, "unauthorized", "company info missing")
			return
		}

		created, err := h.service.Create(&in, *user.CompanyID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
			return
		}
		writeSuccess(w, http.StatusCreated, created, "add-on created", nil)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
	}
}

func (h *AddOnHandler) HandleByIdAddOn(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "id cannot be empty")
		return
	}

	switch r.Method {
	case http.MethodGet:
		addOn, err := h.service.FindById(id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
			return
		}

		writeSuccess(w, http.StatusOK, addOn, "add-on found", nil)
	case http.MethodPut:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
			return
		}

		var in models.AddOnInput
		if err := json.Unmarshal(body, &in); err != nil {
			writeError(w, http.StatusBadRequest, "bad_request", "invalid JSON format")
			return
		}
		if strings.TrimSpace(in.Name) == "" {
			writeError(w, http.StatusBadRequest, "validation_error", "name cannot be empty")
			return
		}

		updated, err := h.service.Update(&in, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
			return
		}

		writeSuccess(w, http.StatusOK, updated, "add-on updated", nil)
	case http.MethodDelete:
		err := h.service.Delete(id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_server_error", err.Error())
			return
		}
		writeSuccess(w, http.StatusOK, nil, "add-on deleted", nil)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
	}
}
