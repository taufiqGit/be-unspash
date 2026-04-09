package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"gowes/models"
	"gowes/services"
	"io"
	"net/http"
	"strings"
)

type CashierShiftHandler struct {
	service services.CashierShiftService
}

func NewCashierShiftHandler(service services.CashierShiftService) *CashierShiftHandler {
	return &CashierShiftHandler{service: service}
}

func (h *CashierShiftHandler) StartShift(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	user, ok := r.Context().Value(UserContextKey).(models.User)
	if !ok || user.ID == "" || user.CompanyID == nil || *user.CompanyID == "" {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "user or company info missing")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "failed to read body")
		return
	}

	var input models.StartCashierShiftInput
	if err := json.Unmarshal(body, &input); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON format")
		return
	}

	if strings.TrimSpace(input.OutletID) == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "outlet_id cannot be empty")
		return
	}

	shift, err := h.service.StartShift(*user.CompanyID, user.ID, input)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrCashierShiftOutletRequired), errors.Is(err, services.ErrCashierShiftUserRequired):
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		case errors.Is(err, services.ErrCashierShiftAlreadyActive):
			writeError(w, http.StatusConflict, "CONFLICT", err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to start cashier shift")
		}
		return
	}

	writeSuccess(w, http.StatusCreated, shift, "cashier shift started", nil)
}

func (h *CashierShiftHandler) EndShift(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	user, ok := r.Context().Value(UserContextKey).(models.User)
	if !ok || user.ID == "" || user.CompanyID == nil || *user.CompanyID == "" {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "user or company info missing")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "failed to read body")
		return
	}

	var input models.EndCashierShiftInput
	if err := json.Unmarshal(body, &input); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON format")
		return
	}

	shift, err := h.service.EndShift(*user.CompanyID, user.ID, input)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrCashierShiftUserRequired):
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		case errors.Is(err, sql.ErrNoRows):
			writeError(w, http.StatusNotFound, "NOT_FOUND", "active cashier shift not found")
		default:
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to end cashier shift")
		}
		return
	}

	writeSuccess(w, http.StatusOK, shift, "cashier shift ended", nil)
}
