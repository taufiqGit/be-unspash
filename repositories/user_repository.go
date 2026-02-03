package repositories

import (
	"database/sql"
	"errors"
	"gowes/models"
)

type UserRepository interface {
	Create(user models.User) (models.User, error)
	FindByEmail(email string) (models.User, error)
	FindByUsername(username string) (models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user models.User) (models.User, error) {
	// Note: We use DEFAULT uuid_generate_v4() for ID in SQL, so we scan it back
	query := `
		INSERT INTO users (username, email, password_hash, role, pos_pin, company_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	err := r.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.PosPIN,
		user.CompanyID,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *userRepository) FindByEmail(email string) (models.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, pos_pin, company_id, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	var user models.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.PosPIN,
		&user.CompanyID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, nil // Return empty user if not found
		}
		return models.User{}, err
	}
	return user, nil
}

func (r *userRepository) FindByUsername(username string) (models.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, pos_pin, company_id, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	var user models.User
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.PosPIN,
		&user.CompanyID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, nil
		}
		return models.User{}, err
	}
	return user, nil
}
