package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

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
	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing code or state"})
		return
	}

	claims, err := services.ValidateToken(state)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state token"})
		return
	}
	userID := claims.UserID
	userEmail := claims.Email

	clientID := os.Getenv("INSTAGRAM_CLIENT_ID")
	clientSecret := os.Getenv("INSTAGRAM_CLIENT_SECRET")
	redirectURI := os.Getenv("INSTAGRAM_REDIRECT_URI")

	// Exchange code for short-lived token
	shortURL := "https://api.instagram.com/oauth/access_token"
	resp, err := http.PostForm(shortURL, map[string][]string{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {redirectURI},
		"code":          {code},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token exchange failed"})
		return
	}
	defer resp.Body.Close()

	var shortResp struct {
		AccessToken string `json:"access_token"`
		UserID      string `json:"user_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&shortResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid token response"})
		return
	}

	// Exchange short token for long-lived token
	longURL := fmt.Sprintf(
		"https://graph.instagram.com/access_token?grant_type=ig_exchange_token&client_secret=%s&access_token=%s",
		clientSecret, shortResp.AccessToken,
	)
	longRespRaw, err := http.Get(longURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get long-lived token"})
		return
	}
	defer longRespRaw.Body.Close()

	var longResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(longRespRaw.Body).Decode(&longResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid long token response"})
		return
	}

	// Get IG user info
	meURL := fmt.Sprintf("https://graph.instagram.com/me?fields=id,username&access_token=%s", longResp.AccessToken)
	meResp, err := http.Get(meURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch ig user"})
		return
	}
	defer meResp.Body.Close()
	var me struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}
	if err := json.NewDecoder(meResp.Body).Decode(&me); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid ig user response"})
		return
	}

	// Save instagram account associated with our user
	account := services.InstagramAccount{
		UserID:        userID,
		IGUserID:      me.ID,
		Username:      me.Username,
		AccessToken:   longResp.AccessToken,
		TokenExpireAt: time.Now().Add(time.Duration(longResp.ExpiresIn) * time.Second),
	}
	if err := services.SaveInstagramAccount(&account); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save instagram account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "connected",
		"username": me.Username,
		"user":     userEmail,
	})
}
