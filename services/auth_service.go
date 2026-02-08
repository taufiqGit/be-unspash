package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gowes/models"
	"gowes/repositories"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(input models.UserRegisterInput) (models.User, error)
	Login(input models.LoginInput) (models.AuthResponse, error)
}

type authService struct {
	userRepo    repositories.UserRepository
	companyRepo repositories.CompanyRepository
	db          *sql.DB
}

func NewAuthService(userRepo repositories.UserRepository, companyRepo repositories.CompanyRepository, db *sql.DB) AuthService {
	return &authService{
		userRepo:    userRepo,
		companyRepo: companyRepo,
		db:          db,
	}
}

func (s *authService) Register(input models.UserRegisterInput) (models.User, error) {
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
	if strings.TrimSpace(input.BussinessName) == "" {
		return models.User{}, errors.New("business name cannot be empty")
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

	// 4. Start Transaction
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.User{}, err
	}
	defer tx.Rollback()

	// 5. Create Company
	newCompany := models.Company{
		Name:      input.BussinessName,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	createdCompany, err := s.companyRepo.Create(ctx, tx, newCompany)
	if err != nil {
		return models.User{}, err
	}

	// 6. Create User (Admin) linked to Company
	// First user is always Admin
	newUser := models.User{
		Username:     input.Username,
		Email:        input.Email,
		Phone:        &input.Phone,
		PasswordHash: string(hashedPassword),
		Role:         models.RoleAdmin,
		CompanyID:    &createdCompany.ID,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	if input.PosPIN != "" {
		pin := input.PosPIN
		newUser.PosPIN = &pin
	}

	createdUser, err := s.userRepo.Create(ctx, tx, newUser)
	if err != nil {
		return models.User{}, err
	}

	// 7. Commit Transaction
	if err := tx.Commit(); err != nil {
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
	fmt.Println(*user.CompanyID)
	claims := jwt.MapClaims{
		"sub":        user.ID,
		"role":       user.Role,
		"company_id": user.CompanyID,
		"exp":        time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
