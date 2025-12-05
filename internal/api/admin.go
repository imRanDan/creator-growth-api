package api

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imRanDan/creator-growth-api/internal/database"
)

// GetWaitlistEntries returns all waitlist entries (admin only)
// Simple password-based auth for now (you can upgrade to JWT later)
func GetWaitlistEntries(c *gin.Context) {
	// Simple password check via query param or header
	// In production, use proper JWT auth or API key
	adminPassword := c.GetHeader("X-Admin-Password")
	if adminPassword == "" {
		adminPassword = c.Query("password")
	}

	expectedPassword := os.Getenv("ADMIN_PASSWORD")
	if expectedPassword == "" {
		expectedPassword = "changeme" // Default for dev
	}

	if adminPassword != expectedPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get pagination params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}
	offset := (page - 1) * limit

	// Get total count
	var total int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM waitlist").Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count entries"})
		return
	}

	// Get entries
	rows, err := database.DB.Query(
		"SELECT id, email, created_at FROM waitlist ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch entries"})
		return
	}
	defer rows.Close()

	var entries []map[string]interface{}
	for rows.Next() {
		var id, email, createdAt string
		if err := rows.Scan(&id, &email, &createdAt); err != nil {
			continue
		}
		entries = append(entries, map[string]interface{}{
			"id":         id,
			"email":      email,
			"created_at": createdAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"entries": entries,
		"total":   total,
		"page":    page,
		"limit":   limit,
		"pages":   (total + limit - 1) / limit, // Ceiling division
	})
}

