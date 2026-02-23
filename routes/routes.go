package routes

import (
	"net/http"

	"gowes/handlers"
)

func RegisterTodoRoutes(mux *http.ServeMux, h *handlers.TodoHandler) {
	mux.HandleFunc("/api/todos", h.ListOrCreate)
	mux.HandleFunc("/api/todos/", h.HandleByID)
}

func RegisterSystemRoutes(mux *http.ServeMux, h *handlers.SystemHandler) {
	mux.HandleFunc("/api/database/tables", h.DatabaseTablesHandler)
	mux.HandleFunc("/api/database/columns", h.TableColumnsHandler)
}

func RegisterAuthRoutes(mux *http.ServeMux, h *handlers.AuthHandler) {
	mux.HandleFunc("/api/register", h.Register)
	mux.HandleFunc("/api/login", h.Login)
}

func RegisterCategoryRoutes(mux *http.ServeMux, h *handlers.CategoryHandler) {
	// Protected routes wrapped with AuthMiddleware
	mux.Handle("/api/categories", handlers.AuthMiddleware(http.HandlerFunc(h.ListOrCreate)))
	mux.Handle("/api/categories/", handlers.AuthMiddleware(http.HandlerFunc(h.HandleByID)))
}

func RegisterAddOnRoutes(mux *http.ServeMux, h *handlers.AddOnHandler) {
	mux.Handle("/api/add-on", handlers.AuthMiddleware(http.HandlerFunc(h.ListOrCreateAddOn)))
	// Protected routes wrapped with AuthMiddleware
	mux.Handle("/api/add-on/{id}", handlers.AuthMiddleware(http.HandlerFunc(h.HandleByIdAddOn)))
}

func RegisterOrderTypesRoutes(mux *http.ServeMux, h *handlers.OrderTypeHandler) {
	mux.Handle("/api/order-types", handlers.AuthMiddleware(http.HandlerFunc(h.ListOrCreate)))
	mux.Handle("/api/order-types/{id}", handlers.AuthMiddleware(http.HandlerFunc(h.HandleByID)))
}

func RegisterOutletRoutes(mux *http.ServeMux, h *handlers.OutletHandler) {
	mux.Handle("/api/outlets", handlers.AuthMiddleware(http.HandlerFunc(h.ListOrCreate)))
	mux.Handle("/api/outlets/{id}", handlers.AuthMiddleware(http.HandlerFunc(h.HandlerById)))
}
