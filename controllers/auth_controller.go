package controllers

import (
	"encoding/json"
	"gowes/models"
	"gowes/services"
	"net/http"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	var input models.UserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid JSON format")
		return
	}

	user, err := c.authService.Register(input)
	if err != nil {
		// Differentiate errors if needed (e.g. duplicate vs system error)
		// For now, returning 400 for business logic errors is common
		writeError(w, http.StatusBadRequest, "REGISTRATION_FAILED", err.Error())
		return
	}

	writeSuccess(w, http.StatusCreated, user, "registration successful", nil)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
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
