package database

import (
	"log"
)

// RunMigrations runs all required database schema migrations
func RunMigrations() error {
	// 1️⃣ Create users table
	userTable := `
	CREATE EXTENSION IF NOT EXISTS "pgcrypto";

	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);
	`

	_, err := DB.Exec(userTable)
	if err != nil {
		return err
	}

	// 2️⃣ Create instagram_accounts table
	instagramTable := `
	CREATE TABLE IF NOT EXISTS instagram_accounts (
		id SERIAL PRIMARY KEY,
		user_id UUID REFERENCES users(id) ON DELETE CASCADE,
		instagram_user_id TEXT UNIQUE NOT NULL,
		username TEXT,
		access_token TEXT NOT NULL,
		connected_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);
	`

	_, err = DB.Exec(instagramTable)
	if err != nil {
		return err
	}

	log.Println("✅ Database migrations completed")
	return nil
}
