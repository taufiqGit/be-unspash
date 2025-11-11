package controllers

import (
    "encoding/json"
    "io"
    "net/http"
    "strconv"
    "strings"

    "gowes/models"
    "gowes/services"
)

// TodosHandler menangani /api/todos (GET untuk list, POST untuk create)
func TodosHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        todos := services.ListTodos()
        meta := map[string]any{"count": len(todos)}
        writeSuccess(w, http.StatusOK, todos, "list todos", meta)
    case http.MethodPost:
        body, err := io.ReadAll(r.Body)
        if err != nil {
            writeError(w, http.StatusBadRequest, "BAD_REQUEST", "gagal membaca body")
            return
        }
        var in models.TodoInput
        if err := json.Unmarshal(body, &in); err != nil {
            writeError(w, http.StatusBadRequest, "BAD_REQUEST", "format JSON tidak valid")
            return
        }
        if strings.TrimSpace(in.Title) == "" || strings.TrimSpace(in.ImageURL) == "" {
            writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "title dan image url tidak boleh kosong")
            return
        }
        created := services.CreateTodo(in)
        writeSuccess(w, http.StatusCreated, created, "todo created", nil)
    default:
        w.Header().Set("Allow", "GET, POST")
        writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method tidak diizinkan")
    }
}

// TodoByIDHandler menangani /api/todos/{id} (GET, PUT, DELETE)
func TodoByIDHandler(w http.ResponseWriter, r *http.Request) {
    // Ekstrak ID dari path
    idStr := strings.TrimPrefix(r.URL.Path, "/api/todos/")
    id, err := strconv.Atoi(idStr)
    if err != nil || id <= 0 {
        writeError(w, http.StatusNotFound, "NOT_FOUND", "ID tidak valid")
        return
    }

    switch r.Method {
    case http.MethodGet:
        todo, ok := services.GetTodo(id)
        if !ok {
            writeError(w, http.StatusNotFound, "NOT_FOUND", "todo tidak ditemukan")
            return
        }
        writeSuccess(w, http.StatusOK, todo, "todo detail", nil)
    case http.MethodPut:
        var in models.TodoInput
        if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
            writeError(w, http.StatusBadRequest, "BAD_REQUEST", "format JSON tidak valid")
            return
        }
        if strings.TrimSpace(in.Title) == "" {
            writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "title tidak boleh kosong")
            return
        }
        updated, ok := services.UpdateTodo(id, in)
        if !ok {
            writeError(w, http.StatusNotFound, "NOT_FOUND", "todo tidak ditemukan")
            return
        }
        writeSuccess(w, http.StatusOK, updated, "todo updated", nil)
    case http.MethodDelete:
        if ok := services.DeleteTodo(id); !ok {
            writeError(w, http.StatusNotFound, "NOT_FOUND", "todo tidak ditemukan")
            return
        }
        writeSuccess(w, http.StatusOK, map[string]int{"id": id}, "todo deleted", nil)
    default:
        w.Header().Set("Allow", "GET, PUT, DELETE")
        writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method tidak diizinkan")
    }
}