package services

import (
	"time"
)

// StartTokenRefreshJob runs a background goroutine to refresh IG tokens periodically.
// TODO: implement actual refresh logic that iterates instagram_accounts and refreshes tokens.
func StartTokenRefreshJob() {
	go func() {
		ticker := time.NewTicker(12 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			// TODO: query DB for accounts and refresh long-lived tokens if needed
		}
	}()
}
