package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/imRanDan/creator-growth-api/internal/services"
)

// ConnectInstagram handles the redirect to Instagram's OAuth URL
func ConnectInstagram(c *gin.Context) {
	clientID := os.Getenv("INSTAGRAM_CLIENT_ID")
	redirectURI := os.Getenv("INSTAGRAM_REDIRECT_URI")

	if clientID == "" || redirectURI == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Instagram environment variables not set"})
		return
	}

	authURL := "https://api.instagram.com/oauth/authorize" +
		"?client_id=" + clientID +
		"&redirect_uri=" + redirectURI +
		"&scope=user_profile,user_media" +
		"&response_type=code"

	c.JSON(http.StatusOK, gin.H{
		"url": authURL,
	})
}

// InstagramCallback handles the redirect from Instagram after authorization
func InstagramCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code parameter"})
		return
	}

	// Step 1: Exchange code for token
	tokenResp, err := services.ExchangeCodeForToken(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code", "details": err.Error()})
		return
	}

	// Step 2: Get user profile
	profile, err := services.GetUserProfile(tokenResp.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile", "details": err.Error()})
		return
	}

	// Step 3: Save to database (requires JWT middleware)
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user context"})
		return
	}

	err = services.SaveInstagramAccount(userID, profile.ID, profile.Username, tokenResp.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save Instagram account", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Instagram account connected successfully",
		"username": profile.Username,
	})
}
