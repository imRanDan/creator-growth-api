package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/imRanDan/creator-growth-api/internal/services"
)

// ConnectInstagram handles the redirect to Instagram's OAuth URL
func ConnectInstagram(c *gin.Context) {
	userID := c.GetString("user_id")
	userEmail := c.GetString("user_email")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	clientID := os.Getenv("INSTAGRAM_CLIENT_ID")
	redirectURI := os.Getenv("INSTAGRAM_REDIRECT_URI")
	if clientID == "" || redirectURI == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Instagram  app configuration missing"})
		return
	}

	// Generate state token bound to this user (short-lived)
	state, err := services.GenerateStateToken(userID, userEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state token"})
		return
	}

	authURL := fmt.Sprintf(
		"https://api.instagram.com/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user_profile,user_media&response_type=code&state=%s",
		clientID, redirectURI, state,
	)

	c.JSON(http.StatusOK, gin.H{
		"url": authURL,
	})
}

// InstagramCallback handles the redirect from Instagram after authorization
func InstagramCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing authorization code "})
		return
	}

	if state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing state parameter"})
	}

	// Validate state token to get userID (CSRF + user binding)
	claims, err := services.ValidateToken(state)
	if err != nil || claims.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired state token"})
		return
	}
	userID := claims.UserID

	// Step 1: Exchange code for short-lived token
	shortTokenResp, err := services.ExchangeCodeForToken(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to exchange code for token",
			"details": err.Error(),
		})
		return
	}

	// Step 2: Exchange for long-lived token (60 days)
	longLivedToken, expiresIn, err := services.ExchangeForLongLivedToken(shortTokenResp.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get long-lived token",
			"details": err.Error(),
		})
		return
	}

	// Step 3: Get Instagram profile
	profile, err := services.GetUserProfile(longLivedToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch Instagram profile",
			"details": err.Error(),
		})
		return
	}

	// Step 4: Save to database
	err = services.SaveInstagramAccount(userID, profile.ID, profile.Username, longLivedToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to save Instagram account",
			"details": err.Error(),
		})
		return
	}

	// Success! Optionally redirect to frontend
	c.JSON(http.StatusOK, gin.H{
		"message":    "Instagram account connected successfully!",
		"username":   profile.Username,
		"ig_user_id": profile.ID,
		"expires_in": expiresIn,
		"connected":  true,
	})
}
