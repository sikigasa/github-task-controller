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

type githubPATRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewGithubPATRepository は新しいGithubPATRepositoryを作成する
func NewGithubPATRepository(db *sql.DB, logger *slog.Logger) repository.GithubPATRepository {
	return &githubPATRepository{
		db:     db,
		logger: logger,
	}
}

func (r *githubPATRepository) Create(ctx context.Context, pat *model.GithubPAT) error {
	query := `
		INSERT INTO github_pat (id, user_id, token_encrypted, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) DO UPDATE SET token_encrypted = $3, updated_at = $5
	`

	_, err := r.db.ExecContext(ctx, query,
		pat.ID, pat.UserID, pat.TokenEncrypted, pat.CreatedAt, pat.UpdatedAt,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to create github pat", "error", err)
		return fmt.Errorf("failed to create github pat: %w", err)
	}

	r.logger.InfoContext(ctx, "github pat created/updated", "user_id", pat.UserID)
	return nil
}

func (r *githubPATRepository) FindByUserID(ctx context.Context, userID string) (*model.GithubPAT, error) {
	query := `
		SELECT id, user_id, token_encrypted, created_at, updated_at
		FROM github_pat
		WHERE user_id = $1
	`

	var pat model.GithubPAT
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&pat.ID, &pat.UserID, &pat.TokenEncrypted, &pat.CreatedAt, &pat.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // PATが存在しない場合はnilを返す
	}
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to find github pat", "error", err)
		return nil, fmt.Errorf("failed to find github pat: %w", err)
	}

	return &pat, nil
}

func (r *githubPATRepository) Update(ctx context.Context, pat *model.GithubPAT) error {
	query := `
		UPDATE github_pat
		SET token_encrypted = $1, updated_at = $2
		WHERE user_id = $3
	`

	result, err := r.db.ExecContext(ctx, query, pat.TokenEncrypted, time.Now(), pat.UserID)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to update github pat", "error", err)
		return fmt.Errorf("failed to update github pat: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("github pat not found")
	}

	r.logger.InfoContext(ctx, "github pat updated", "user_id", pat.UserID)
	return nil
}

func (r *githubPATRepository) Delete(ctx context.Context, userID string) error {
	query := `DELETE FROM github_pat WHERE user_id = $1`

	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to delete github pat", "error", err)
		return fmt.Errorf("failed to delete github pat: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("github pat not found")
	}

	r.logger.InfoContext(ctx, "github pat deleted", "user_id", userID)
	return nil
}
