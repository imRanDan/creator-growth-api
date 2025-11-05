package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// Token response structure
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	UserID      int64  `json:"user_id"`
}

// Profile structure
type InstagramProfile struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// Media structure
type InstagramMedia struct {
	ID        string `json:"id"`
	Caption   string `json:"caption"`
	MediaType string `json:"media_type"`
	MediaURL  string `json:"media_url"`
	Timestamp string `json:"timestamp"`
}

type MediaResponse struct {
	Data []InstagramMedia `json:"data"`
}

// Exchange authorization code for access token
func ExchangeCodeForToken(code string) (*TokenResponse, error) {
	clientID := os.Getenv("INSTAGRAM_CLIENT_ID")
	clientSecret := os.Getenv("INSTAGRAM_CLIENT_SECRET")
	redirectURI := os.Getenv("INSTAGRAM_REDIRECT_URI")

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", redirectURI)
	data.Set("code", code)

	resp, err := http.PostForm("https://api.instagram.com/oauth/access_token", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// Get user profile from Instagram
func GetUserProfile(accessToken string) (*InstagramProfile, error) {
	url := fmt.Sprintf("https://graph.instagram.com/me?fields=id,username&access_token=%s", accessToken)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var profile InstagramProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

// Fetch user's media from Instagram
func FetchUserMedia(accessToken string) ([]InstagramMedia, error) {
	url := fmt.Sprintf("https://graph.instagram.com/me/media?fields=id,caption,media_type,media_url,timestamp&access_token=%s", accessToken)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var mediaResp MediaResponse
	if err := json.NewDecoder(resp.Body).Decode(&mediaResp); err != nil {
		return nil, err
	}

	return mediaResp.Data, nil
}

// Database functions (stub - implement with your DB)
func SaveInstagramAccount(userID, igUserID, username, accessToken string) error {
	// TODO: Implement database save
	// INSERT INTO instagram_accounts (user_id, ig_user_id, username, access_token)
	fmt.Println("Saving Instagram account:", userID, igUserID, username)
	return nil
}

func GetInstagramAccount(userID string) (*struct{ AccessToken string }, error) {
	// TODO: Implement database query
	// SELECT access_token FROM instagram_accounts WHERE user_id = ?
	return &struct{ AccessToken string }{AccessToken: "dummy_token"}, nil
}

func DeleteInstagramAccount(userID string) error {
	// TODO: Implement database delete
	// DELETE FROM instagram_accounts WHERE user_id = ?
	fmt.Println("Deleting Instagram account:", userID)
	return nil
}
