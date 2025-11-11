package routes

import (
    "net/http"

    "gowes/controllers"
)

func RegisterTodoRoutes(mux *http.ServeMux) {
    mux.HandleFunc("/api/todos", controllers.TodosHandler)
    mux.HandleFunc("/api/todos/", controllers.TodoByIDHandler)
}