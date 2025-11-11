package controllers

import (
    "encoding/json"
    "net/http"
)

type APIError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

type APIResponse struct {
    Success bool              `json:"success"`
    Message string            `json:"message,omitempty"`
    Data    any               `json:"data,omitempty"`
    Error   *APIError         `json:"error,omitempty"`
    Meta    map[string]any    `json:"meta,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(v)
}

func writeSuccess(w http.ResponseWriter, status int, data any, message string, meta map[string]any) {
    writeJSON(w, status, APIResponse{Success: true, Message: message, Data: data, Meta: meta})
}

func writeError(w http.ResponseWriter, status int, code, message string) {
    writeJSON(w, status, APIResponse{Success: false, Error: &APIError{Code: code, Message: message}})
}