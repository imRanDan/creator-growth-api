package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imRanDan/creator-growth-api/internal/models"
	"github.com/imRanDan/creator-growth-api/internal/services"
)

type RegisterRequest struct {
	Email		string		`json:"email" binding:"required,email"`
	Password	string		`json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email		string		`json:"email" binding:"required,email"`
	Password	string		`json:"password" binding:"required`
}

type AuthResponse struct {
	Token	string				 `json:"token"`
	User	models.UserResponse	 `json:"user"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) 
		return
	}

	//Create User
	user, err := services.CreateUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H("error": "Failed to create user"))
		return
	}

	//Generate token
	token, err := services.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	//Return response
	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User: models.UserResponse{
			ID:			user.ID,
			Email:		user.Email,
			CreatedAt:	user.CreatedAt,
		},
	})
}

