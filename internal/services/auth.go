package services

import (
	"errors"

	"github.com/imRanDan/creator-growth-api/internal/database"
	"github.com/imRanDan/creator-growth-api/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateUser(email, password string) (*models.User, error) {
	//Hash password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	//Insert user into database
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
		return nil, err
	}

	return &user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
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
		return nil, err
	}

	return &user, nil
}

func AuthenticateUser(email, password string) (*models.User, error) {
	user, err := GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}
