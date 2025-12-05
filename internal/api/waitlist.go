package api

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/imRanDan/creator-growth-api/internal/database"
	"github.com/imRanDan/creator-growth-api/internal/services"
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
	result, err := database.DB.Exec(
		"INSERT INTO waitlist (email) VALUES ($1) ON CONFLICT (email) DO NOTHING",
		req.Email,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add email to waitlist"})
		return
	}

	// Check if a row was actually inserted (new signup vs duplicate)
	rowsAffected, _ := result.RowsAffected()
	isNewSignup := rowsAffected > 0

	// Send welcome email only for new signups (in background, don't block response)
	if isNewSignup {
		go func() {
			if err := services.SendWelcomeEmail(req.Email); err != nil {
				// Log error but don't fail the request
				fmt.Printf("Failed to send welcome email to %s: %v\n", req.Email, err)
			}
		}()
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully added to waitlist!",
		"email":   req.Email,
	})
}


