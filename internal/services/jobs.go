package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/imRanDan/creator-growth-api/internal/database"
)

// StartTokenRefreshJob runs a background goroutine to refresh IG tokens periodically.
// Checks every 12 hours for tokens expiring within 7 days and refreshes them.
func StartTokenRefreshJob() {
	go func() {
		// Run once on startup after a short delay
		time.Sleep(30 * time.Second)
		refreshExpiringTokens()

		// Then run every 12 hours
		ticker := time.NewTicker(12 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			refreshExpiringTokens()
		}
	}()
	log.Println("✅ Token refresh job started (runs every 12 hours)")
}

// refreshExpiringTokens finds accounts with tokens expiring within 7 days and refreshes them
func refreshExpiringTokens() {
	log.Println("🔄 Checking for expiring Instagram tokens...")

	// Find tokens expiring in the next 7 days
	query := `
		SELECT id, ig_user_id, username, access_token, token_expires_at
		FROM instagram_accounts
		WHERE token_expires_at < NOW() + INTERVAL '7 days'
		  AND token_expires_at > NOW()
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("❌ Error querying expiring tokens: %v", err)
		return
	}
	defer rows.Close()

	var refreshed, failed int
	for rows.Next() {
		var id, igUserID, username, accessToken string
		var expiresAt time.Time

		if err := rows.Scan(&id, &igUserID, &username, &accessToken, &expiresAt); err != nil {
			log.Printf("❌ Error scanning row: %v", err)
			failed++
			continue
		}

		log.Printf("🔄 Refreshing token for @%s (expires: %s)", username, expiresAt.Format("Jan 2"))

		// Refresh the token
		newToken, newExpiry, err := RefreshLongLivedToken(accessToken)
		if err != nil {
			log.Printf("❌ Failed to refresh token for @%s: %v", username, err)
			failed++
			continue
		}

		// Update the database
		if err := updateAccountToken(id, newToken, newExpiry); err != nil {
			log.Printf("❌ Failed to save new token for @%s: %v", username, err)
			failed++
			continue
		}

		log.Printf("✅ Refreshed token for @%s (new expiry: %s)", username, newExpiry.Format("Jan 2"))
		refreshed++
	}

	if refreshed > 0 || failed > 0 {
		log.Printf("🔄 Token refresh complete: %d refreshed, %d failed", refreshed, failed)
	} else {
		log.Println("✅ No tokens need refreshing")
	}
}

// RefreshLongLivedToken calls Instagram API to refresh a long-lived token
func RefreshLongLivedToken(currentToken string) (string, time.Time, error) {
	url := fmt.Sprintf(
		"https://graph.instagram.com/refresh_access_token?grant_type=ig_refresh_token&access_token=%s",
		currentToken,
	)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", time.Time{}, fmt.Errorf("instagram refresh error: status %d", resp.StatusCode)
	}

	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int64  `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", time.Time{}, fmt.Errorf("failed to parse refresh response: %w", err)
	}

	newExpiry := time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)
	return result.AccessToken, newExpiry, nil
}

// updateAccountToken updates the token and expiry in the database
func updateAccountToken(accountID, newToken string, newExpiry time.Time) error {
	query := `
		UPDATE instagram_accounts
		SET access_token = $1, token_expires_at = $2, updated_at = NOW()
		WHERE id = $3
	`
	_, err := database.DB.Exec(query, newToken, newExpiry, accountID)
	return err
}
