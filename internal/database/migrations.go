package database

import (
	"log"
)

// RunMigrations runs all required database schema migrations
func RunMigrations() error {
	// 1️⃣ Create users table and enable gen_random_uuid
	userTable := `
    CREATE EXTENSION IF NOT EXISTS "pgcrypto";

    CREATE TABLE IF NOT EXISTS users (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        email VARCHAR(255) UNIQUE NOT NULL,
        password VARCHAR(255) NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    );
    `

	_, err := DB.Exec(userTable)
	if err != nil {
		return err
	}

	// 2️⃣ Create instagram_accounts table (UUID PK, token expiry, followers)
	instagramTable := `
    CREATE TABLE IF NOT EXISTS instagram_accounts (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        ig_user_id TEXT UNIQUE NOT NULL,
        username TEXT,
        access_token TEXT NOT NULL,
        token_expires_at TIMESTAMP WITH TIME ZONE,
        followers BIGINT,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    );
    `

	_, err = DB.Exec(instagramTable)
	if err != nil {
		return err
	}

	// 3️⃣ Create instagram_posts table
	instagramPosts := `
    CREATE TABLE IF NOT EXISTS instagram_posts (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        ig_post_id TEXT NOT NULL,
        account_id UUID NOT NULL REFERENCES instagram_accounts(id) ON DELETE CASCADE,
        caption TEXT,
        media_type TEXT,
        media_url TEXT,
        like_count INT,
        comments_count INT,
        posted_at TIMESTAMP WITH TIME ZONE,
        fetched_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        UNIQUE(account_id, ig_post_id)
    );
    `

	_, err = DB.Exec(instagramPosts)
	if err != nil {
		return err
	}

	log.Println("✅ Database migrations completed")
	return nil
}
