package controllers

import (
	"gowes/services"
	"net/http"
)

// SystemController handles HTTP requests for database metadata
type SystemController struct {
	service services.SystemService
}

// NewSystemController creates a new SystemController
func NewSystemController(service services.SystemService) *SystemController {
	return &SystemController{service: service}
}

// DatabaseTablesHandler handles /api/database/tables
func (c *SystemController) DatabaseTablesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tables, err := c.service.ListTables()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "DATABASE_ERROR", "gagal mengambil daftar tabel")
			return
		}
		meta := map[string]any{"count": len(tables)}
		writeSuccess(w, http.StatusOK, tables, "daftar tabel database", meta)
	default:
		w.Header().Set("Allow", "GET")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method tidak diizinkan")
	}
}

// TableColumnsHandler handles /api/database/columns
func (c *SystemController) TableColumnsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method tidak diizinkan")
		return
	}

	schema := r.URL.Query().Get("schema")
	table := r.URL.Query().Get("table")

	if schema == "" || table == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "schema dan table parameter wajib diisi")
		return
	}

	columns, err := c.service.GetTableColumns(schema, table)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "DATABASE_ERROR", "gagal mengambil struktur kolom")
		return
	}

	meta := map[string]any{"schema": schema, "table": table, "count": len(columns)}
	writeSuccess(w, http.StatusOK, columns, "struktur kolom tabel", meta)
}
