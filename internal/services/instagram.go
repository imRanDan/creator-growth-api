package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/imRanDan/creator-growth-api/internal/database"
)

// === SHORT-LIVED TOKEN RESPONSE ===
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	UserID      int64  `json:"user_id"`
}

type InstagramPost struct {
	ID           string    `db:"id"`
	IGPostID     string    `db:"ig_post_id"`
	AccountID    string    `db:"account_id"`
	Caption      string    `db:"caption"`
	MediaType    string    `db:"media_type"`
	MediaURL     string    `db:"media_url"`
	LikeCount    int       `db:"like_count"`
	CommentCount int       `db:"comment_count"`
	PostedAt     time.Time `db:"posted_at"`
	FetchedAt    time.Time `db:"fetched_at"`
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
	account, err := GetInstagramAccountByID(accountID)
	if err != nil {
		return fmt.Errorf("failed to load account: %w", err)
	}

	if account.IGUserID == "" || account.AccessToken == "" {
		return fmt.Errorf("account missing ig_user_id or access token")
	}

	//2) fetch recent media from IG Graph API
	media, err := FetchUserMedia(account.AccessToken, 50) //fetches the last 50 posts
	if err != nil {
		return fmt.Errorf("failed to fetch user media: %w", err)
	}

	if len(media) == 0 {
		log.Printf("no media returned for account %s", accountID)
		return nil
	}

	// 3) upsert (update + insert) each post into instagram_posts table
	for _, m := range media {
		post := &InstagramPost{
			IGPostID:     m.ID,
			AccountID:    accountID,
			Caption:      m.Caption,
			MediaType:    m.MediaType,
			MediaURL:     m.MediaURL,
			LikeCount:    m.LikeCount,
			CommentCount: m.CommentCount,
			FetchedAt:    time.Now().UTC(),
		}
		// parse timestamp from IG (ISO 8601)
		if m.Timestamp != "" {
			if t, err := time.Parse(time.RFC3339, m.Timestamp); err == nil {
				post.PostedAt = t
			}
		}

		if err := UpsertInstagramPost(post); err != nil {
			log.Printf("failed to upsert post %s: %v", m.ID, err)
			// continue on error so we don't fail entire batch
		}
	}
	log.Printf("✅ fetched and stored %d posts for account %s", len(media), accountID)
	return nil
}

// UpsertInstagramPost inserts or updates a post in instagram_posts
func UpsertInstagramPost(post *InstagramPost) error {
	if post == nil {
		return errors.New("nil post")
	}
	if post.ID == "" {
		post.ID = uuid.New().String()
	}
	if post.FetchedAt.IsZero() {
		post.FetchedAt = time.Now().UTC()
	}

	query := `
		INSERT INTO instagram_posts
          (id, ig_post_id, account_id, caption, media_type, media_url, like_count, comments_count, posted_at, fetched_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
        ON CONFLICT (account_id, ig_post_id) DO UPDATE
          SET caption = EXCLUDED.caption,
              media_type = EXCLUDED.media_type,
              media_url = EXCLUDED.media_url,
              like_count = EXCLUDED.like_count,
              comments_count = EXCLUDED.comments_count,
              posted_at = EXCLUDED.posted_at,
              fetched_at = EXCLUDED.fetched_at
        RETURNING id
	`

	var id string
	err := database.DB.QueryRow(
		query,
		post.ID,
		post.IGPostID,
		post.AccountID,
		post.Caption,
		post.MediaType,
		post.MediaURL,
		post.LikeCount,
		post.CommentCount,
		post.PostedAt,
		post.FetchedAt,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("failed to upsert post: %w", err)
	}
	post.ID = id
	return nil
}

// GetPostsBy AccountID returns the recent posts for an account
func GetPostsByAccountID(accountID string, limit int) ([]InstagramPost, error) {
	if limit <= 0 {
		limit = 25
	}

	query := `
		SELECT id, ig_post_id, account_id, caption, media_type, media_url, like_count, comments_count, posted_at, fetched_at
		FROM instagram_posts
		WHERE account_id = $1
		ORDER BY posted_at DESC
		LIMIT $2
	`

	rows, err := database.DB.Query(query, accountID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []InstagramPost
	for rows.Next() {
		var p InstagramPost
		if err := rows.Scan(
			&p.ID, &p.IGPostID, &p.AccountID, &p.Caption, &p.MediaType, &p.MediaURL,
			&p.LikeCount, &p.CommentCount, &p.PostedAt, &p.FetchedAt,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, rows.Err()
}

// ExtractHashtags extracts hashtags from a caption string
func ExtractHashtags(caption string) []string {
	var hashtags []string
	words := strings.Fields(caption)
	for _, word := range words {
		if strings.HasPrefix(word, "#") {
			//clean punctuation
			tag := strings.Trim(word, ".,!?;:")
			if len(tag) > 1 {
				hashtags = append(hashtags, strings.ToLower(tag))
			}
		}
	}
	return hashtags
}
