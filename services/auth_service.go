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
	VerifyEmail(token string) error
}

type authService struct {
	userRepo    repositories.UserRepository
	outletRepo  repositories.OutletRepository
	companyRepo repositories.CompanyRepository
	emailRepo   repositories.EmailRepository
	db          *sql.DB
}

func NewAuthService(userRepo repositories.UserRepository, companyRepo repositories.CompanyRepository, outletRepo repositories.OutletRepository, emailRepo repositories.EmailRepository, db *sql.DB) AuthService {
	return &authService{
		userRepo:    userRepo,
		companyRepo: companyRepo,
		outletRepo:  outletRepo,
		emailRepo:   emailRepo,
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
		IsOwner:      true,
		CompanyID:    &createdCompany.ID,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	if input.PosPIN != "" {
		pin := input.PosPIN
		newUser.PosPIN = &pin
	}

	newOutlet := models.OutletInput{
		Code:       fmt.Sprintf("%s-001", input.BussinessName),
		Name:       input.BussinessName,
		Supervisor: "-",
		Address:    "-",
		Phone:      input.Phone,
		Email:      input.Email,
		IsActive:   true,
	}

	createdUser, errUser := s.userRepo.Create(ctx, tx, newUser)
	if errUser != nil {
		return models.User{}, errUser
	}

	_, errOutlet := s.outletRepo.Create(&newOutlet, createdCompany.ID, ctx, tx)
	if errOutlet != nil {
		return models.User{}, errOutlet
	}

	// 7. Commit Transaction
	if err := tx.Commit(); err != nil {
		return models.User{}, err
	}

	// 8. Generate JWT for verify email
	token, err := generateJWT(createdUser, true)
	if err != nil {
		return models.User{}, err
	}

	// 9. Send verification email
	if err := s.emailRepo.SendVerificationEmail(ctx, createdUser.Email, createdUser.Username, os.Getenv("FE_VERIFY_MAIL"), token); err != nil {
		return models.User{}, err
	}

	return createdUser, nil
}

func (s *authService) Login(input models.LoginInput) (models.AuthResponse, error) {
	// 1. Find User
	user, err := s.userRepo.FindByEmail(input.Identifier)
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

	if !user.Active {
		return models.AuthResponse{}, errors.New("account not active")
	}

	// 4. Generate JWT
	token, err := generateJWT(user, false)
	if err != nil {
		return models.AuthResponse{}, err
	}

	return models.AuthResponse{
		AccessToken: token,
		User:        user,
	}, nil
}

func (s *authService) VerifyEmail(token string) error {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return err
	}

	user_id, ok := claims["sub"].(string)
	if !ok {
		return errors.New("invalid token")
	}

	_, err = s.userRepo.FindByID(user_id)
	if err != nil {
		return err
	}

	_, err = s.userRepo.ChangeActivateUser(user_id)
	if err != nil {
		return err
	}

	return nil
}

func generateJWT(user models.User, for_verified bool) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-me" // Fallback for dev
	}
	fmt.Println(*user.CompanyID)
	claims := jwt.MapClaims{
		"sub":          user.ID,
		"role":         user.Role,
		"company_id":   user.CompanyID,
		"exp":          time.Now().Add(time.Hour * 24).Unix(), // 24 hours
		"for_verified": for_verified,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
