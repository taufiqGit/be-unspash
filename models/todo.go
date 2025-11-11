package models

import "time"

type Todo struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Done      bool      `json:"done"`
    ImageURL  string    `json:"image_url"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type TodoInput struct {
    Title string `json:"title"`
    Done  bool   `json:"done"`
    ImageURL string `json:"image_url"`
}