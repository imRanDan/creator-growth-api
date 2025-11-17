package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/imRanDan/creator-growth-api/internal/database"
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

// === DATABASE TYPES & OPERATIONS ===

// InstagramAccount represents a connected IG account
type InstagramAccount struct {
	ID            string
	UserID        string
	IGUserID      string
	Username      string
	AccessToken   string
	TokenExpireAt time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Followers     sql.NullInt64
}

// SaveInstagramAccount inserts or updates an instagram_accounts row
func SaveInstagramAccount(a *InstagramAccount) error {
	if a == nil {
		return errors.New("nil account")
	}
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	now := time.Now().UTC()
	a.UpdatedAt = now
	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
	}

	query := `
        INSERT INTO instagram_accounts
          (id, user_id, ig_user_id, username, access_token, token_expires_at, created_at, updated_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
        ON CONFLICT (ig_user_id) DO UPDATE
          SET username = EXCLUDED.username,
              access_token = EXCLUDED.access_token,
              token_expires_at = EXCLUDED.token_expires_at,
              updated_at = EXCLUDED.updated_at
        RETURNING id
    `
	var id string
	err := database.DB.QueryRow(
		query,
		a.ID,
		a.UserID,
		a.IGUserID,
		a.Username,
		a.AccessToken,
		a.TokenExpireAt,
		a.CreatedAt,
		a.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return err
	}
	a.ID = id
	return nil
}

// GetInstagramAccountByID returns account by internal id
func GetInstagramAccountByID(id string) (*InstagramAccount, error) {
	var a InstagramAccount
	query := `SELECT id, user_id, ig_user_id, username, access_token, token_expires_at, created_at, updated_at FROM instagram_accounts WHERE id = $1`
	err := database.DB.QueryRow(query, id).Scan(
		&a.ID, &a.UserID, &a.IGUserID, &a.Username, &a.AccessToken, &a.TokenExpireAt, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// GetInstagramAccountByUserID returns the first IG account for a given app user
func GetInstagramAccountByUserID(userID string) (*InstagramAccount, error) {
	var a InstagramAccount
	query := `SELECT id, user_id, ig_user_id, username, access_token, token_expires_at, created_at, updated_at FROM instagram_accounts WHERE user_id = $1 LIMIT 1`
	err := database.DB.QueryRow(query, userID).Scan(
		&a.ID, &a.UserID, &a.IGUserID, &a.Username, &a.AccessToken, &a.TokenExpireAt, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// DeleteInstagramAccountByUserID removes an account for the given app user
func DeleteInstagramAccountByUserID(userID string) error {
	_, err := database.DB.Exec("DELETE FROM instagram_accounts WHERE user_id = $1", userID)
	return err
}

// FetchAndStorePosts: skeleton that will call IG Graph API and upsert posts
func FetchAndStorePosts(accountID string) error {
	// TODO:
	// 1) load instagram_accounts by id, get IGUserID & AccessToken
	// 2) call GET https://graph.instagram.com/{ig_user_id}/media?fields=id,caption,media_type,media_url,timestamp,like_count,comments_count&access_token=...
	// 3) upsert into instagram_posts table
	// 4) return
	return nil
}
