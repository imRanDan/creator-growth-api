package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/imRanDan/creator-growth-api/internal/services"
)

// Start OAuth flow
func ConnectInstagram(c *gin.Context) {
	clientID := os.Getenv("INSTAGRAM_CLIENT_ID")
	redirectURI := os.Getenv("INSTAGRAM_REDIRECT_URI")

	authURL := fmt.Sprintf(
		"https://api.instagram.com/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user_profile,user_media&response_type=code",
		clientID,
		redirectURI,
	)

	c.JSON(http.StatusOK, gin.H{
		"auth_URL": authURL,
	})
}

// OAuth callback
func InstagramCallback(c *gin.Context) {
	code := c.Query("code")
	if code == " " {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No authorization code provided"})
		return
	}

	//get user id from context (set by middleware)
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	//get user profile
	profile, err := services.GetUserProfile(tokenResp.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user profile", "details": err.Error()})
		return
	}

	// Save Instagram account
	err = services.SaveInstagramAccount(userId.(string), profile.ID, profile.Username, tokenResp.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save Instagram account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Instagram account connected successfully",
		"username": profile.Username,
	})
}

func GetInstagramMedia(c *gin.Context) {
	userId, _ := c.Get("user_id")

	//Get Instagram user from the DB
	account, err := services.GetInstagramAccount(userId.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Instagram account not connected"})
		return
	}

	//Fetch media from Instagram API
	media, err := services.FetchUserMedia(account.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch media"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"media": media,
	})
}

// Disconnect Instagram account
func DisconnectInstagram(c *gin.Context) {
	userId, _ := c.Get("user_id")

	err := services.DeleteInstagramAccount(userId.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disconnect account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Instagram account disconnected",
	})
}
