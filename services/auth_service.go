package services

import (
	"errors"
	"gowes/models"
	"gowes/repositories"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(input models.UserInput) (models.User, error)
	Login(input models.LoginInput) (models.AuthResponse, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(input models.UserInput) (models.User, error) {
	// 1. Validation
	if strings.TrimSpace(input.Username) == "" {
		return models.User{}, errors.New("username cannot be empty")
	}
	if strings.TrimSpace(input.Email) == "" {
		return models.User{}, errors.New("email cannot be empty")
	}
	if len(input.Password) < 6 {
		return models.User{}, errors.New("password must be at least 6 characters")
	}

	// 2. Check Duplicates
	existingUser, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return models.User{}, err
	}
	if existingUser.ID != "" {
		return models.User{}, errors.New("email already registered")
	}

	existingUser, err = s.userRepo.FindByUsername(input.Username)
	if err != nil {
		return models.User{}, err
	}
	if existingUser.ID != "" {
		return models.User{}, errors.New("username already taken")
	}

	// 3. Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	// 4. Create User
	// Default role is Waiter if not specified
	role := input.Role
	if role == "" {
		role = models.RoleWaiter
	}

	newUser := models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		Role:         role,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	if input.PosPIN != "" {
		pin := input.PosPIN
		newUser.PosPIN = &pin
	}

	// Note: CompanyID is left nil for now as per simple register flow.
	// In a real app, we might create a company here or link to one.

	createdUser, err := s.userRepo.Create(newUser)
	if err != nil {
		return models.User{}, err
	}

	return createdUser, nil
}

func (s *authService) Login(input models.LoginInput) (models.AuthResponse, error) {
	// 1. Find User
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return models.AuthResponse{}, err
	}
	if user.ID == "" {
		return models.AuthResponse{}, errors.New("invalid email or password")
	}

	// 2. Check Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return models.AuthResponse{}, errors.New("invalid email or password")
	}

	// 3. Generate JWT
	token, err := generateJWT(user)
	if err != nil {
		return models.AuthResponse{}, err
	}

	return models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func generateJWT(user models.User) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-me" // Fallback for dev
	}

	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
