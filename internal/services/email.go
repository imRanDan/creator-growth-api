package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// ResendEmail sends an email using Resend API (free tier: 3,000 emails/month)
func ResendEmail(to, subject, htmlBody string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		// If no API key, just log and continue (graceful degradation)
		fmt.Println("⚠️  RESEND_API_KEY not set - skipping email send")
		return nil
	}

	fromEmail := os.Getenv("RESEND_FROM_EMAIL")
	if fromEmail == "" {
		fromEmail = "onboarding@resend.dev" // Resend default
	}

	payload := map[string]interface{}{
		"from":    fromEmail,
		"to":      []string{to},
		"subject": subject,
		"html":    htmlBody,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling email payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("resend API error (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// SendWelcomeEmail sends a welcome email to new waitlist signups
func SendWelcomeEmail(email string) error {
	subject := "Welcome to Creator Growth! 🎉"
	htmlBody := `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="utf-8">
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px; text-align: center; border-radius: 10px 10px 0 0; }
				.content { background: #f9f9f9; padding: 30px; border-radius: 0 0 10px 10px; }
				.button { display: inline-block; background: #667eea; color: white; padding: 12px 24px; text-decoration: none; border-radius: 5px; margin-top: 20px; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>📈 Welcome to Creator Growth!</h1>
				</div>
				<div class="content">
					<p>Hey there! 👋</p>
					<p>Thanks for joining the waitlist! We're excited to have you on board.</p>
					<p>We'll notify you as soon as we launch. In the meantime, get ready to:</p>
					<ul>
						<li>📊 Track your Instagram engagement in real-time</li>
						<li>🚀 Get smart insights to grow your audience</li>
						<li>⚡ See which posts perform best</li>
					</ul>
					<p>Talk soon!</p>
					<p>— The Creator Growth Team</p>
				</div>
			</div>
		</body>
		</html>
	`

	return ResendEmail(email, subject, htmlBody)
}


