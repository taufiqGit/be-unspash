package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"gowes/models"
	"gowes/services"
	"gowes/utils"
)

type OutletHandler struct {
	service services.OutletService
}

func NewOutletHandler(service services.OutletService) *OutletHandler {
	return &OutletHandler{service: service}
}

func (h *OutletHandler) ListOrCreate(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(UserContextKey).(models.User)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Unauthorized", "User not found in context")
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Get pagination params
		params := utils.ParsePaginationParams(r)
		outlets, total, err := h.service.FindAll(*user.CompanyID, params)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error(), "Failed to get outlets")
			return
		}
		meta := utils.CalculateMeta(total, params)
		writeSuccess(w, http.StatusOK, outlets, "Outlet list", meta)
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "bad_request", "invalid JSON format")
			return
		}
		var input models.OutletInput
		if err := json.Unmarshal(body, &input); err != nil {
			writeError(w, http.StatusBadRequest, "bad_request", "invalid JSON format")
			return
		}
		outlet, err := h.service.Create(*user.CompanyID, input)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error(), "Failed to create outlet")
			return
		}
		writeSuccess(w, http.StatusCreated, outlet, "Outlet created successfully", nil)

	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed", "Method not allowed")
		return
	}
}
