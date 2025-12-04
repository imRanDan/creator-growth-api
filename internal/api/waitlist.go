package api

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/imRanDan/creator-growth-api/internal/database"
)

// WaitlistSignup handles email signup for the waitlist
func WaitlistSignup(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	// Basic email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Insert email into waitlist (ON CONFLICT DO NOTHING to handle duplicates gracefully)
	_, err := database.DB.Exec(
		"INSERT INTO waitlist (email) VALUES ($1) ON CONFLICT (email) DO NOTHING",
		req.Email,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add email to waitlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully added to waitlist!",
		"email":   req.Email,
	})
}

