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
	mux.HandleFunc("/api/verify-email", h.VerifyEmail)
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

func RegisterProductRoutes(mux *http.ServeMux, h *handlers.ProductHandler) {
	mux.Handle("/api/products", handlers.AuthMiddleware(http.HandlerFunc(h.ListOrCreate)))
	// Protected routes wrapped with AuthMiddleware
	mux.Handle("/api/products/{id}", handlers.AuthMiddleware(http.HandlerFunc(h.HandleByID)))
}

func RegisterCustomerRoutes(mux *http.ServeMux, h *handlers.CustomerHandler) {
	mux.Handle("/api/customers", handlers.AuthMiddleware(http.HandlerFunc(h.ListOrCreate)))
	mux.Handle("/api/customers/{id}", handlers.AuthMiddleware(http.HandlerFunc(h.HandleByID)))
}

func RegisterDiscountRoutes(mux *http.ServeMux, h *handlers.DiscountHandler) {
	mux.Handle("/api/discounts", handlers.AuthMiddleware(http.HandlerFunc(h.ListOrCreate)))
	mux.Handle("/api/discounts/{id}", handlers.AuthMiddleware(http.HandlerFunc(h.HandleByID)))
}

func RegisterTaxRoutes(mux *http.ServeMux, h *handlers.TaxHandler) {
	mux.Handle("/api/taxes", handlers.AuthMiddleware(http.HandlerFunc(h.ListOrCreate)))
	mux.Handle("/api/taxes/{id}", handlers.AuthMiddleware(http.HandlerFunc(h.HandleByID)))
}

func RegisterRoleRoutes(mux *http.ServeMux, h *handlers.RoleHandler) {
	mux.Handle("/api/roles", handlers.AuthMiddleware(http.HandlerFunc(h.ListOrCreate)))
	mux.Handle("/api/roles/{id}", handlers.AuthMiddleware(http.HandlerFunc(h.HandleByID)))
}

func RegisterUnitRoutes(mux *http.ServeMux, h *handlers.UnitHandler) {
	mux.Handle("/api/units", handlers.AuthMiddleware(http.HandlerFunc(h.ListOrCreate)))
	mux.Handle("/api/units/{id}", handlers.AuthMiddleware(http.HandlerFunc(h.HandleByID)))
}

func RegisterSupplierRoutes(mux *http.ServeMux, h *handlers.SupplierHandler) {
	mux.Handle("/api/suppliers", handlers.AuthMiddleware(http.HandlerFunc(h.ListOrCreate)))
	mux.Handle("/api/suppliers/{id}", handlers.AuthMiddleware(http.HandlerFunc(h.HandleByID)))
}

func RegisterRecipeRoutes(mux *http.ServeMux, h *handlers.RecipeHandler) {
	mux.Handle("/api/recipes", handlers.AuthMiddleware(http.HandlerFunc(h.ListOrCreate)))
	mux.Handle("/api/recipes/{id}", handlers.AuthMiddleware(http.HandlerFunc(h.HandleByID)))
}
