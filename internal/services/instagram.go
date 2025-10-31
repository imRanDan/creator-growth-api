package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/imRanDan/creator-growth-api/internal/database"
)

type InstagramTokenResponse struct {
	AccessToken string `json:"access_token"`
	user_id     int64  `json:"user_id"`
}

type InstagramUserProfile struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func ExchangeCodeForToken(code string) (*InstagramTokenResponse, error) {
	clientID := os.Getenv("INSTAGRAM_CLIENT_ID")
	clientSecret := os.Getenv("INSTAGRAM_CLIENT_SECRET")
	redirectURI := os.Getenv("INSTAGRAM_REDIRECT_URI")

	tokenURL := "https://api.instagram.com/oauth/access_token"

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", redirectURI)
	data.Set("code", code)

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("instagram api error: %s", string(body))
	}

	var tokenResp InstagramTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

func GetUserProfile(accessToken string) (*InstagramUserProfile, error) {
	url := fmt.Sprintf("https://graph.instagram.com/me?fields=id,username&access_token=%s", accessToken)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("instagram api error: %s", string(body))
	}

	var profile InstagramUserProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

func SaveInstagramAccount(userID, instagramUserID, username, accessToken string) error {
	query := `
	INSERT INTO instagram_accounts (user_id, instagram_user_id, username, access_token, connected_at)
	VALUES ($1, $2, $3, $4, NOW())
	ON CONFLICT (instagram_user_id)
	DO UPDATE SET
		access_token = EXCLUDED.access_token,
		username = EXCLUDED.username,
		updated_at = NOW()
	`

	_, err := database.DB.Exec(query, userID, instagramUserID, username, accessToken)
	return err
}
