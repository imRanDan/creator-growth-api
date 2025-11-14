package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// === SHORT-LIVED TOKEN RESPONSE ===
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	UserID      int64  `json:"user_id"`
}

// ExchangeCodeForToken exchanges Instagram auth code for short-lived access token
func ExchangeCodeForToken(code string) (*TokenResponse, error) {
	clientID := os.Getenv("INSTAGRAM_CLIENT_ID")
	clientSecret := os.Getenv("INSTAGRAM_CLIENT_SECRET")
	redirectURI := os.Getenv("INSTAGRAM_REDIRECT_URI")

	if clientID == "" || clientSecret == "" || redirectURI == "" {
		return nil, fmt.Errorf("missing Instagram OAuth environment variables")
	}

	data := url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {redirectURI},
		"code":          {code},
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.PostForm("https://api.instagram.com/oauth/access_token", data)
	if err != nil {
		return nil, fmt.Errorf("failed to send token request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("instagram token error %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokenResp, nil
}

// === LONG-LIVED TOKEN RESPONSE ===
type LongLivedTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"` // seconds until expiry
}

// ExchangeForLongLivedToken converts short-lived token to long-lived (60 days)
func ExchangeForLongLivedToken(shortToken string) (string, int64, error) {
	clientSecret := os.Getenv("INSTAGRAM_CLIENT_SECRET")
	if clientSecret == "" {
		return "", 0, fmt.Errorf("INSTAGRAM_CLIENT_SECRET not set")
	}

	url := fmt.Sprintf(
		"https://graph.instagram.com/access_token?grant_type=ig_exchange_token&client_secret=%s&access_token=%s",
		clientSecret, shortToken,
	)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", 0, fmt.Errorf("long-lived token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("instagram long-lived error %d: %s", resp.StatusCode, string(body))
	}

	var result LongLivedTokenResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", 0, fmt.Errorf("failed to parse long-lived token: %w", err)
	}

	return result.AccessToken, result.ExpiresIn, nil
}

// === PROFILE & MEDIA STRUCTS ===
type InstagramProfile struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type InstagramMedia struct {
	ID           string `json:"id"`
	Caption      string `json:"caption"`
	MediaType    string `json:"media_type"`
	MediaURL     string `json:"media_url"`
	Timestamp    string `json:"timestamp"`
	LikeCount    int    `json:"like_count,omitempty"`
	CommentCount int    `json:"comments_count,omitempty"`
}

type MediaResponse struct {
	Data []InstagramMedia `json:"data"`
}

// GetUserProfile fetches Instagram user profile using long-lived token
func GetUserProfile(accessToken string) (*InstagramProfile, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token required")
	}

	url := fmt.Sprintf("https://graph.instagram.com/me?fields=id,username&access_token=%s", accessToken)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("profile request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("instagram profile error %d: %s", resp.StatusCode, string(body))
	}

	var profile InstagramProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("failed to decode profile: %w", err)
	}

	return &profile, nil
}

// FetchUserMedia gets recent media with engagement metrics
func FetchUserMedia(accessToken string, limit int) ([]InstagramMedia, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("access token required")
	}
	if limit <= 0 {
		limit = 25
	}

	fields := "id,caption,media_type,media_url,timestamp,like_count,comments_count"
	url := fmt.Sprintf("https://graph.instagram.com/me/media?fields=%s&access_token=%s&limit=%d", fields, accessToken, limit)

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("media request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("instagram media error %d: %s", resp.StatusCode, string(body))
	}

	var mediaResp MediaResponse
	if err := json.NewDecoder(resp.Body).Decode(&mediaResp); err != nil {
		return nil, fmt.Errorf("failed to decode media: %w", err)
	}

	return mediaResp.Data, nil
}

// === DATABASE OPERATIONS (TO BE IMPLEMENTED) ===
func SaveInstagramAccount(userID, igUserID, username, accessToken string) error {
	// TODO: Save to DB with upsert
	// Store: user_id, instagram_user_id, username, access_token, expires_at
	fmt.Printf("Saving IG account: user=%s, ig=%s, @%s\n", userID, igUserID, username)
	return nil
}

func GetInstagramAccount(userID string) (*struct {
	AccessToken string
	ExpiresAt   time.Time
}, error) {
	// TODO: Fetch from DB
	return nil, fmt.Errorf("not implemented")
}

func DeleteInstagramAccount(userID string) error {
	// TODO: Delete from DB
	fmt.Printf("Deleting IG account for user: %s\n", userID)
	return nil
}
