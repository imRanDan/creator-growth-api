package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/imRanDan/creator-growth-api/internal/models"
	"github.com/imRanDan/creator-growth-api/internal/services"
)

// === REQUEST STRUCTS WITH VALIDATION ===
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string              `json:"token"`
	User  models.UserResponse `json:"user"`
}

// Register handles user registration
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": validationError(err),
		})
		return
	}

	// Create user
	user, err := services.CreateUser(req.Email, req.Password)
	if err != nil {
		if err.Error() == "email already in use" {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
		return
	}

	// Generate JWT
	token, err := services.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User: models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	})
}

// Login handles user login
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": validationError(err),
		})
		return
	}

	// Authenticate
	user, err := services.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT
	token, err := services.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User: models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	})
}

// GetCurrentUser returns authenticated user info
func GetCurrentUser(c *gin.Context) {
	email, exists := c.Get("user_email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := services.GetUserByEmail(email.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	})
}

// validationError formats Gin validation errors
func validationError(err error) []string {
	if err == nil {
		return nil
	}
	// If it's validator.ValidationErrors, extract all messages
	if ve, ok := err.(validator.ValidationErrors); ok {
		out := make([]string, 0, len(ve))
		for _, fe := range ve {
			out = append(out, fe.Error())
		}
		return out
	}
	// Fallback to a single error string
	return []string{err.Error()}
}
