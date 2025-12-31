package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/sikigasa/github-task-controller/backend/internal/domain/model"
	"github.com/sikigasa/github-task-controller/backend/internal/domain/repository"
)

type googleAccountRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewGoogleAccountRepository は新しいGoogleAccountRepositoryを作成する
func NewGoogleAccountRepository(db *sql.DB, logger *slog.Logger) repository.GoogleAccountRepository {
	return &googleAccountRepository{
		db:     db,
		logger: logger,
	}
}

func (r *googleAccountRepository) Create(ctx context.Context, account *model.GoogleAccount) error {
	query := `
		INSERT INTO google_account (user_id, provider, provider_account_id, access_token, refresh_token, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	// expires_at を Unix timestamp に変換
	var expiresAt *int64
	if account.ExpiresAt != nil {
		ts := account.ExpiresAt.Unix()
		expiresAt = &ts
	}

	_, err := r.db.ExecContext(ctx, query,
		account.UserID, account.Provider, account.ProviderAccountID,
		account.AccessToken, account.RefreshToken, expiresAt,
		account.CreatedAt, account.UpdatedAt,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to create google account", "error", err)
		return fmt.Errorf("failed to create google account: %w", err)
	}

	r.logger.InfoContext(ctx, "google account created", "user_id", account.UserID)
	return nil
}

func (r *googleAccountRepository) FindByProviderAccountID(ctx context.Context, provider, providerAccountID string) (*model.GoogleAccount, error) {
	query := `
		SELECT user_id, provider, provider_account_id, access_token, refresh_token, expires_at, created_at, updated_at
		FROM google_account
		WHERE provider = $1 AND provider_account_id = $2
	`

	var account model.GoogleAccount
	var expiresAt sql.NullInt64
	err := r.db.QueryRowContext(ctx, query, provider, providerAccountID).Scan(
		&account.UserID, &account.Provider, &account.ProviderAccountID,
		&account.AccessToken, &account.RefreshToken, &expiresAt,
		&account.CreatedAt, &account.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("google account not found: %s", providerAccountID)
	}
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to find google account", "error", err)
		return nil, fmt.Errorf("failed to find google account: %w", err)
	}

	if expiresAt.Valid {
		t := time.Unix(expiresAt.Int64, 0)
		account.ExpiresAt = &t
	}

	return &account, nil
}

func (r *googleAccountRepository) FindByUserID(ctx context.Context, userID string) (*model.GoogleAccount, error) {
	query := `
		SELECT user_id, provider, provider_account_id, access_token, refresh_token, expires_at, created_at, updated_at
		FROM google_account
		WHERE user_id = $1
	`

	var account model.GoogleAccount
	var expiresAt sql.NullInt64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&account.UserID, &account.Provider, &account.ProviderAccountID,
		&account.AccessToken, &account.RefreshToken, &expiresAt,
		&account.CreatedAt, &account.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("google account not found for user: %s", userID)
	}
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to find google account by user_id", "error", err)
		return nil, fmt.Errorf("failed to find google account: %w", err)
	}

	if expiresAt.Valid {
		t := time.Unix(expiresAt.Int64, 0)
		account.ExpiresAt = &t
	}

	return &account, nil
}

func (r *googleAccountRepository) Update(ctx context.Context, account *model.GoogleAccount) error {
	query := `
		UPDATE google_account
		SET access_token = $1, refresh_token = $2, expires_at = $3, updated_at = $4
		WHERE provider = $5 AND provider_account_id = $6
	`

	var expiresAt *int64
	if account.ExpiresAt != nil {
		ts := account.ExpiresAt.Unix()
		expiresAt = &ts
	}

	result, err := r.db.ExecContext(ctx, query,
		account.AccessToken, account.RefreshToken, expiresAt, time.Now(),
		account.Provider, account.ProviderAccountID,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to update google account", "error", err)
		return fmt.Errorf("failed to update google account: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("google account not found")
	}

	r.logger.InfoContext(ctx, "google account updated")
	return nil
}

func (r *googleAccountRepository) Delete(ctx context.Context, provider, providerAccountID string) error {
	query := `DELETE FROM google_account WHERE provider = $1 AND provider_account_id = $2`

	result, err := r.db.ExecContext(ctx, query, provider, providerAccountID)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to delete google account", "error", err)
		return fmt.Errorf("failed to delete google account: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("google account not found")
	}

	r.logger.InfoContext(ctx, "google account deleted")
	return nil
}

type githubAccountRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewGithubAccountRepository は新しいGithubAccountRepositoryを作成する
func NewGithubAccountRepository(db *sql.DB, logger *slog.Logger) repository.GithubAccountRepository {
	return &githubAccountRepository{
		db:     db,
		logger: logger,
	}
}

func (r *githubAccountRepository) Create(ctx context.Context, account *model.GithubAccount) error {
	query := `
		INSERT INTO github_account (user_id, provider, provider_account_id, access_token, refresh_token, expires_at, pat_encrypted, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	var expiresAt *int64
	if account.ExpiresAt != nil {
		ts := account.ExpiresAt.Unix()
		expiresAt = &ts
	}

	_, err := r.db.ExecContext(ctx, query,
		account.UserID, account.Provider, account.ProviderAccountID,
		account.AccessToken, account.RefreshToken, expiresAt, account.PATEncrypted,
		account.CreatedAt, account.UpdatedAt,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to create github account", "error", err)
		return fmt.Errorf("failed to create github account: %w", err)
	}

	r.logger.InfoContext(ctx, "github account created", "user_id", account.UserID)
	return nil
}

func (r *githubAccountRepository) FindByProviderAccountID(ctx context.Context, provider, providerAccountID string) (*model.GithubAccount, error) {
	query := `
		SELECT user_id, provider, provider_account_id, access_token, refresh_token, expires_at, pat_encrypted, created_at, updated_at
		FROM github_account
		WHERE provider = $1 AND provider_account_id = $2
	`

	var account model.GithubAccount
	var expiresAt sql.NullInt64
	var patEncrypted sql.NullString
	err := r.db.QueryRowContext(ctx, query, provider, providerAccountID).Scan(
		&account.UserID, &account.Provider, &account.ProviderAccountID,
		&account.AccessToken, &account.RefreshToken, &expiresAt, &patEncrypted,
		&account.CreatedAt, &account.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("github account not found: %s", providerAccountID)
	}
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to find github account", "error", err)
		return nil, fmt.Errorf("failed to find github account: %w", err)
	}

	if expiresAt.Valid {
		t := time.Unix(expiresAt.Int64, 0)
		account.ExpiresAt = &t
	}
	if patEncrypted.Valid {
		account.PATEncrypted = &patEncrypted.String
	}

	return &account, nil
}

func (r *githubAccountRepository) FindByUserID(ctx context.Context, userID string) (*model.GithubAccount, error) {
	query := `
		SELECT user_id, provider, provider_account_id, access_token, refresh_token, expires_at, pat_encrypted, created_at, updated_at
		FROM github_account
		WHERE user_id = $1
	`

	var account model.GithubAccount
	var expiresAt sql.NullInt64
	var patEncrypted sql.NullString
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&account.UserID, &account.Provider, &account.ProviderAccountID,
		&account.AccessToken, &account.RefreshToken, &expiresAt, &patEncrypted,
		&account.CreatedAt, &account.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // アカウントが存在しない場合はnilを返す
	}
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to find github account by user_id", "error", err)
		return nil, fmt.Errorf("failed to find github account: %w", err)
	}

	if expiresAt.Valid {
		t := time.Unix(expiresAt.Int64, 0)
		account.ExpiresAt = &t
	}
	if patEncrypted.Valid {
		account.PATEncrypted = &patEncrypted.String
	}

	return &account, nil
}

func (r *githubAccountRepository) Update(ctx context.Context, account *model.GithubAccount) error {
	query := `
		UPDATE github_account
		SET access_token = $1, refresh_token = $2, expires_at = $3, pat_encrypted = $4, updated_at = $5
		WHERE provider = $6 AND provider_account_id = $7
	`

	var expiresAt *int64
	if account.ExpiresAt != nil {
		ts := account.ExpiresAt.Unix()
		expiresAt = &ts
	}

	result, err := r.db.ExecContext(ctx, query,
		account.AccessToken, account.RefreshToken, expiresAt, account.PATEncrypted, time.Now(),
		account.Provider, account.ProviderAccountID,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to update github account", "error", err)
		return fmt.Errorf("failed to update github account: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("github account not found")
	}

	r.logger.InfoContext(ctx, "github account updated")
	return nil
}

func (r *githubAccountRepository) Delete(ctx context.Context, provider, providerAccountID string) error {
	query := `DELETE FROM github_account WHERE provider = $1 AND provider_account_id = $2`

	result, err := r.db.ExecContext(ctx, query, provider, providerAccountID)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to delete github account", "error", err)
		return fmt.Errorf("failed to delete github account: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("github account not found")
	}

	r.logger.InfoContext(ctx, "github account deleted")
	return nil
}
