package routes

import (
	"net/http"

	"gowes/controllers"
)

func RegisterTodoRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/todos", controllers.TodosHandler)
	mux.HandleFunc("/api/todos/", controllers.TodoByIDHandler)
	mux.HandleFunc("/api/database/tables", controllers.DatabaseTablesHandler)
	mux.HandleFunc("/api/database/columns", controllers.TableColumnsHandler)
}

func RegisterCategoryRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/categories", controllers.CategoriesHandler)
	mux.HandleFunc("/api/categories/", controllers.CategoryByIDHandler)
}
