package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func NewDB(dbPath string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = initSchema(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize database schema: %w", err)
	}

	return db, nil
}

func initSchema(db *sql.DB) error {
	rocketTableSQL := `
	CREATE TABLE IF NOT EXISTS rockets (
		channel TEXT PRIMARY KEY,
		type TEXT NOT NULL,
		speed INTEGER NOT NULL,
		mission TEXT NOT NULL,
		launch_time TIMESTAMP NOT NULL,
		status TEXT NOT NULL,
		exploded_at TIMESTAMP,
		reason TEXT,
		last_updated TIMESTAMP NOT NULL,
		last_message INTEGER NOT NULL
	);`

	if _, err := db.Exec(rocketTableSQL); err != nil {
		return fmt.Errorf("failed to create rockets table: %w", err)
	}

	messagesTableSQL := `
	CREATE TABLE IF NOT EXISTS processed_messages (
		channel TEXT NOT NULL,
		message_number INTEGER NOT NULL,
		processed_at TIMESTAMP NOT NULL,
		PRIMARY KEY (channel, message_number)
	);`

	if _, err := db.Exec(messagesTableSQL); err != nil {
		return fmt.Errorf("failed to create processed_messages table: %w", err)
	}

	return nil
}
