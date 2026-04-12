package handlers

import (
	"fmt"
	"gowes/models"
	"gowes/services"
	"gowes/utils"
	"net/http"
	"strings"
)

type ProductHandler struct {
	service services.ProductService
}

func NewProductHandler(service services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) ListOrCreate(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(UserContextKey).(models.User)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Unauthorized", "User not found in context")
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Get pagination params
		params := utils.ParsePaginationParams(r)
		products, total, err := h.service.FindAll(*user.CompanyID, params)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error(), "Failed to get products")
			return
		}
		meta := utils.CalculateMeta(total, params)
		writeSuccess(w, http.StatusOK, products, "Product list", meta)
	case http.MethodPost:
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			writeError(w, http.StatusBadRequest, "bad_request", "invalid multipart form")
			return
		}
		name := r.FormValue("name")
		price := r.FormValue("price")
		sku := r.FormValue("sku")
		unit := r.FormValue("unit")
		unitID := r.FormValue("unit_id")
		cost := r.FormValue("cost")
		categoryID := r.FormValue("category_id")
		addOnIDs := r.FormValue("add_on_ids")
		var addOnIDList []string
		if addOnIDs != "" {
			addOnIDList = strings.Split(addOnIDs, ",")
			fmt.Println(addOnIDList)
			for _, addOnID := range addOnIDList {
				fmt.Println(addOnID)
			}
		} else {
			addOnIDList = []string{}
		}
		file, header, err := r.FormFile("image")
		if err != nil {
			writeError(w, http.StatusBadRequest, "bad_request", "invalid image file")
			return
		}
		defer file.Close()

		if name == "" || price == "" || sku == "" || unit == "" || unitID == "" || cost == "" || categoryID == "" {
			writeError(w, http.StatusBadRequest, "bad_request", "missing required fields")
			return
		}

		payload := models.ProductInput{
			Name:       name,
			Price:      utils.ParseFloat64(price),
			SKU:        sku,
			Unit:       unit,
			UnitID:     unitID,
			Cost:       utils.ParseFloat64(cost),
			CategoryID: categoryID,
			CompanyID:  *user.CompanyID,
			ImageURL:   "",
		}

		product, err := h.service.Create(*user.CompanyID, payload, file, header, addOnIDList)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error(), "Failed to create product")
			return
		}
		writeSuccess(w, http.StatusCreated, product, "Product created", nil)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}
}

func (h *ProductHandler) HandleByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "id cannot be empty")
		return
	}

	user, ok := r.Context().Value(UserContextKey).(models.User)
	if !ok {
		writeError(w, http.StatusUnauthorized, "Unauthorized", "User not found in context")
		return
	}

	switch r.Method {
	case http.MethodGet:
		product, err := h.service.FindByID(id)
		if err != nil {
			writeError(w, http.StatusNotFound, "not_found", "Product not found")
			return
		}
		writeSuccess(w, http.StatusOK, product, "Product detail", nil)
	case http.MethodDelete:
		err := h.service.DeleteById(id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error(), "Failed to delete product")
			return
		}
		writeSuccess(w, http.StatusOK, nil, "Product deleted", nil)
	case http.MethodPut:
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			writeError(w, http.StatusBadRequest, "bad_request", "invalid multipart form")
			return
		}
		name := r.FormValue("name")
		price := r.FormValue("price")
		sku := r.FormValue("sku")
		unit := r.FormValue("unit")
		unitID := r.FormValue("unit_id")
		cost := r.FormValue("cost")
		categoryID := r.FormValue("category_id")
		if name == "" || price == "" || sku == "" || unit == "" || unitID == "" || cost == "" || categoryID == "" {
			writeError(w, http.StatusBadRequest, "bad_request", "missing required fields")
			return
		}

		payload := models.ProductInput{
			Name:       name,
			Price:      utils.ParseFloat64(price),
			SKU:        sku,
			Unit:       unit,
			UnitID:     unitID,
			Cost:       utils.ParseFloat64(cost),
			CategoryID: categoryID,
			CompanyID:  *user.CompanyID,
		}

		// Gambar bersifat opsional saat update — jika tidak dikirim, tetap pakai gambar lama
		imageFile, imageHeader, err := r.FormFile("image")
		if err == nil {
			defer imageFile.Close()
		} else {
			imageFile = nil
			imageHeader = nil
		}

		updated, err := h.service.Update(id, payload, imageFile, imageHeader)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error(), "Failed to update product")
			return
		}
		writeSuccess(w, http.StatusOK, updated, "Product updated", nil)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}
}
