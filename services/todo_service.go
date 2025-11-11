package services

import (
    "sync"
    "time"

    "gowes/models"
)

var (
    todos  = make([]models.Todo, 0)
    nextID = 1
    mu     sync.RWMutex
)

func ListTodos() []models.Todo {
    mu.RLock()
    defer mu.RUnlock()
    // return a copy to avoid external mutation
    cp := make([]models.Todo, len(todos))
    copy(cp, todos)
    return cp
}

func GetTodo(id int) (models.Todo, bool) {
    mu.RLock()
    defer mu.RUnlock()
    for _, t := range todos {
        if t.ID == id {
            return t, true
        }
    }
    return models.Todo{}, false
}

func CreateTodo(in models.TodoInput) models.Todo {
    mu.Lock()
    defer mu.Unlock()
    now := time.Now().UTC()
    t := models.Todo{
        ID:        nextID,
        Title:     in.Title,
        Done:      in.Done,
        ImageURL:  in.ImageURL,
        CreatedAt: now,
        UpdatedAt: now,
    }
    nextID++
    todos = append(todos, t)
    return t
}

func UpdateTodo(id int, in models.TodoInput) (models.Todo, bool) {
    mu.Lock()
    defer mu.Unlock()
    for i, t := range todos {
        if t.ID == id {
            t.Title = in.Title
            t.Done = in.Done
            t.ImageURL = in.ImageURL
            t.UpdatedAt = time.Now().UTC()
            todos[i] = t
            return t, true
        }
    }
    return models.Todo{}, false
}

func DeleteTodo(id int) bool {
    mu.Lock()
    defer mu.Unlock()
    for i, t := range todos {
        if t.ID == id {
            todos = append(todos[:i], todos[i+1:]...)
            return true
        }
    }
    return false
}