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

type StockHandler struct {
	service services.StockService
}

func NewStockHandler(service services.StockService) *StockHandler {
	return &StockHandler{service: service}
}

func (h *StockHandler) List(w http.ResponseWriter, r *http.Request) {
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

	stocks, total, err := h.service.ListStocks(*user.CompanyID, params, outletID, productID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list stocks")
		return
	}

	meta := utils.CalculateMeta(total, params)
	writeSuccess(w, http.StatusOK, stocks, "stock list", meta)
}

func (h *StockHandler) GetByOutletAndProduct(w http.ResponseWriter, r *http.Request) {
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

	outletID := strings.TrimSpace(r.PathValue("outlet_id"))
	productID := strings.TrimSpace(r.PathValue("product_id"))

	stock, err := h.service.GetStock(*user.CompanyID, outletID, productID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrStockOutletRequired), errors.Is(err, services.ErrStockProductRequired):
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		case errors.Is(err, sql.ErrNoRows):
			writeError(w, http.StatusNotFound, "NOT_FOUND", "stock not found")
		default:
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get stock")
		}
		return
	}

	writeSuccess(w, http.StatusOK, stock, "stock detail", nil)
}
