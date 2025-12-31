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
// 注意: 本番環境ではマイグレーションツール（golang-migrate等）を使用してください
func InitSchema(ctx context.Context, db *sql.DB, logger *slog.Logger) error {
	schema := `
		-- pg_uuidv7拡張が必要な場合はマイグレーションで実行
		-- CREATE EXTENSION IF NOT EXISTS pg_uuidv7;

		CREATE TABLE IF NOT EXISTS users (
			id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) UNIQUE NOT NULL,
			name VARCHAR(255),
			image_url TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS github_account (
			user_id uuid NOT NULL,
			provider VARCHAR NOT NULL,
			provider_account_id VARCHAR NOT NULL,
			access_token VARCHAR,
			refresh_token VARCHAR,
			expires_at BIGINT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT github_account_pk PRIMARY KEY (provider, provider_account_id),
			CONSTRAINT github_account_user_fk
				FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS google_account (
			user_id uuid NOT NULL,
			provider VARCHAR NOT NULL,
			provider_account_id VARCHAR NOT NULL,
			access_token VARCHAR,
			refresh_token VARCHAR,
			expires_at BIGINT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT google_account_pk PRIMARY KEY (provider, provider_account_id),
			CONSTRAINT google_account_user_fk
				FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS project (
			id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id uuid NOT NULL,
			title VARCHAR NOT NULL,
			description TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT project_user_fk
				FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS task (
			id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			project_id uuid NOT NULL,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			status INT NOT NULL,
			end_date TIMESTAMP,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT task_project_fk
				FOREIGN KEY (project_id) REFERENCES project(id) ON DELETE CASCADE
		);

		-- インデックス
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
		CREATE INDEX IF NOT EXISTS idx_github_account_user_id ON github_account(user_id);
		CREATE INDEX IF NOT EXISTS idx_google_account_user_id ON google_account(user_id);
		CREATE INDEX IF NOT EXISTS idx_project_user_id ON project(user_id);
		CREATE INDEX IF NOT EXISTS idx_task_project_id ON task(project_id);
		CREATE INDEX IF NOT EXISTS idx_task_status ON task(status);
	`

	_, err := db.ExecContext(ctx, schema)
	if err != nil {
		logger.ErrorContext(ctx, "failed to initialize schema", "error", err)
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	logger.InfoContext(ctx, "database schema initialized")
	return nil
}
