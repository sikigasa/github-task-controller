package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
)

// DBConfig はデータベース接続設定
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewDB は新しいデータベース接続を作成する
func NewDB(ctx context.Context, cfg DBConfig, logger *slog.Logger) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.ErrorContext(ctx, "failed to open database", "error", err)
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 接続設定
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// 接続確認
	if err := db.PingContext(ctx); err != nil {
		logger.ErrorContext(ctx, "failed to ping database", "error", err)
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.InfoContext(ctx, "database connection established")
	return db, nil
}

// InitSchema はデータベーススキーマを初期化する
func InitSchema(ctx context.Context, db *sql.DB, logger *slog.Logger) error {
	schema := `
		CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(36) PRIMARY KEY,
			email VARCHAR(255) NOT NULL UNIQUE,
			name VARCHAR(255) NOT NULL,
			picture TEXT,
			google_id VARCHAR(255) NOT NULL UNIQUE,
			refresh_token TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
		CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id);

		CREATE TABLE IF NOT EXISTS todos (
			id VARCHAR(36) PRIMARY KEY,
			title VARCHAR(200) NOT NULL,
			description TEXT,
			completed BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_todos_created_at ON todos(created_at);
		CREATE INDEX IF NOT EXISTS idx_todos_completed ON todos(completed);
	`

	_, err := db.ExecContext(ctx, schema)
	if err != nil {
		logger.ErrorContext(ctx, "failed to initialize schema", "error", err)
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	logger.InfoContext(ctx, "database schema initialized")
	return nil
}
