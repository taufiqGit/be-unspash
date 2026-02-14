package handlers

import (
	"encoding/json"
	"fmt"
	"gowes/models"
	"gowes/services"
	"net/http"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (c *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	var input models.UserRegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON format")
		return
	}
	// Debug input value if needed, but r.Body is already drained here
	fmt.Printf("Received input: %+v\n", input)

	user, err := c.authService.Register(input)
	if err != nil {
		// Differentiate errors if needed (e.g. duplicate vs system error)
		// For now, returning 400 for business logic errors is common
		writeError(w, http.StatusBadRequest, "REGISTRATION_FAILED", err.Error())
		return
	}

	writeSuccess(w, http.StatusCreated, user, "registration successful", nil)
}

func (c *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	var input models.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON format")
		return
	}

	response, err := c.authService.Login(input)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "LOGIN_FAILED", err.Error())
		return
	}

	writeSuccess(w, http.StatusOK, response, "login successful", nil)
}
