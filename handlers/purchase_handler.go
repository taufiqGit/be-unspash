package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"gowes/models"
	"gowes/services"
	"gowes/utils"
	"io"
	"net/http"
	"strings"
)

type PurchaseHandler struct {
	service services.PurchaseService
}

func NewPurchaseHandler(service services.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{service: service}
}

func (h *PurchaseHandler) ListOrCreate(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(UserContextKey).(models.User)
	if !ok || user.ID == "" || user.CompanyID == nil || *user.CompanyID == "" {
		writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "user or company info missing")
		return
	}

	switch r.Method {
	case http.MethodGet:
		params := utils.ParsePaginationParams(r)
		purchases, total, err := h.service.ListPurchases(*user.CompanyID, params)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list purchases")
			return
		}
		meta := utils.CalculateMeta(total, params)
		writeSuccess(w, http.StatusOK, purchases, "purchase list", meta)
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "failed to read body")
			return
		}

		var input models.PurchaseInput
		if err := json.Unmarshal(body, &input); err != nil {
			writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON format")
			return
		}

		purchase, err := h.service.CreatePurchase(*user.CompanyID, user.ID, input)
		if err != nil {
			switch {
			case errors.Is(err, services.ErrPurchaseOutletRequired),
				errors.Is(err, services.ErrPurchaseDetailsRequired),
				errors.Is(err, services.ErrPurchaseDetailProductRequired),
				errors.Is(err, services.ErrPurchaseDetailQtyInvalid),
				errors.Is(err, services.ErrPurchaseDetailPriceInvalid):
				writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
			default:
				writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create purchase")
			}
			return
		}

		writeSuccess(w, http.StatusCreated, purchase, "purchase created", nil)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	}
}

func (h *PurchaseHandler) HandleByID(w http.ResponseWriter, r *http.Request) {
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
		purchase, err := h.service.GetPurchase(id, *user.CompanyID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeError(w, http.StatusNotFound, "NOT_FOUND", "purchase not found")
				return
			}
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get purchase")
			return
		}
		writeSuccess(w, http.StatusOK, purchase, "purchase detail", nil)
	default:
		w.Header().Set("Allow", "GET")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
	}
}
