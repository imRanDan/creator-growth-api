package services

import (
	"database/sql"
	"errors"

	"github.com/imRanDan/creator-growth-api/internal/database"
	"github.com/imRanDan/creator-growth-api/internal/models"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compares a password with its hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateUser registers a new user with hashed password
func CreateUser(email, password string) (*models.User, error) {
	// Validate input
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	// Hash password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Insert into DB
	query := `
        INSERT INTO users (email, password)
        VALUES ($1, $2)
        RETURNING id, email, created_at, updated_at
    `

	var user models.User
	err = database.DB.QueryRow(query, email, hashedPassword).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		// Handle duplicate email (Postgres unique violation code 23505)
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, errors.New("email already in use")
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail fetches a user by email (with password)
func GetUserByEmail(email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email required")
	}

	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1`

	var user models.User
	err := database.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// AuthenticateUser verifies email + password and returns user
func AuthenticateUser(email, password string) (*models.User, error) {
	user, err := GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !CheckPassword(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Return user WITHOUT password
	user.Password = ""
	return user, nil
}
