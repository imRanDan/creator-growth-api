package database

import (
	"log"
)

func RunMigrations() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);
	`

	_, err := DB.Exec(query)
	if err != nil {
		return err
	}

	log.Println("✅ Database migrations completed")
	return nil
}
