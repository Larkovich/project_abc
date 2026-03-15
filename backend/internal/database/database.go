package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const maxRetries = 5
const retryInterval = 5 * time.Second

func Connect() (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	for i := range maxRetries {
		if err = db.Ping(); err == nil {
			break
		}
		slog.Warn("database not ready, retrying...",
			"attempt", i+1,
			"max", maxRetries,
			"error", err,
		)
		time.Sleep(retryInterval)
	}
	if err != nil {
		return nil, fmt.Errorf("db.Ping after %d retries: %w", maxRetries, err)
	}

	slog.Info("database connected",
		"host", os.Getenv("DB_HOST"),
		"name", os.Getenv("DB_NAME"),
	)

	return db, nil
}
