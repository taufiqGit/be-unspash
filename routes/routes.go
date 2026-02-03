package routes

import (
	"net/http"

	"gowes/controllers"
)

func RegisterTodoRoutes(mux *http.ServeMux, c *controllers.TodoController) {
	mux.HandleFunc("/api/todos", c.ListOrCreate)
	mux.HandleFunc("/api/todos/", c.HandleByID)
}

func RegisterSystemRoutes(mux *http.ServeMux, c *controllers.SystemController) {
	mux.HandleFunc("/api/database/tables", c.DatabaseTablesHandler)
	mux.HandleFunc("/api/database/columns", c.TableColumnsHandler)
}

func RegisterAuthRoutes(mux *http.ServeMux, c *controllers.AuthController) {
	mux.HandleFunc("/api/register", c.Register)
	mux.HandleFunc("/api/login", c.Login)
}

func RegisterCategoryRoutes(mux *http.ServeMux, c *controllers.CategoryController) {
	// Protected routes wrapped with AuthMiddleware
	mux.Handle("/api/categories", controllers.AuthMiddleware(http.HandlerFunc(c.ListOrCreate)))
	mux.Handle("/api/categories/", controllers.AuthMiddleware(http.HandlerFunc(c.HandleByID)))
}
