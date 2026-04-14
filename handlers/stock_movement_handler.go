package handlers

import (
	"database/sql"
	"errors"
	"gowes/models"
	"gowes/services"
	"gowes/utils"
	"net/http"
	"strings"
)

type StockMovementHandler struct {
	service services.StockMovementService
}

func NewStockMovementHandler(service services.StockMovementService) *StockMovementHandler {
	return &StockMovementHandler{service: service}
}

func (h *StockMovementHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	user, ok := r.Context().Value(UserContextKey).(models.User)
	if !ok || user.ID == "" || user.CompanyID == nil || *user.CompanyID == "" {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "user or company info missing")
		return
	}

	params := utils.ParsePaginationParams(r)
	outletID := strings.TrimSpace(r.URL.Query().Get("outlet_id"))
	productID := strings.TrimSpace(r.URL.Query().Get("product_id"))
	movementType := strings.TrimSpace(r.URL.Query().Get("type"))
	referenceType := strings.TrimSpace(r.URL.Query().Get("reference_type"))

	movements, total, err := h.service.ListStockMovements(*user.CompanyID, params, outletID, productID, movementType, referenceType)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list stock movements")
		return
	}

	meta := utils.CalculateMeta(total, params)
	writeSuccess(w, http.StatusOK, movements, "stock movement list", meta)
}

func (h *StockMovementHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

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

	movement, err := h.service.GetStockMovement(*user.CompanyID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "stock movement not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get stock movement")
		return
	}

	writeSuccess(w, http.StatusOK, movement, "stock movement detail", nil)
}
