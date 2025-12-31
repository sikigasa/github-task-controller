package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/sikigasa/github-task-controller/backend/internal/domain/repository"
	"github.com/sikigasa/github-task-controller/backend/internal/model"
)

// UserRepositoryImpl はUserRepositoryの実装
type UserRepositoryImpl struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewUserRepository は新しいUserRepositoryImplを作成する
func NewUserRepository(db *sql.DB, logger *slog.Logger) repository.UserRepository {
	return &UserRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

// Create は新しいユーザーをデータベースに保存する
func (r *UserRepositoryImpl) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, email, name, picture, google_id, refresh_token, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Name,
		user.Picture,
		user.GoogleID,
		user.RefreshToken,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to insert user", "error", err)
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

// FindByID はIDでユーザーを取得する
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, email, name, picture, google_id, refresh_token, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user model.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Picture,
		&user.GoogleID,
		&user.RefreshToken,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		r.logger.ErrorContext(ctx, "failed to query user", "id", id, "error", err)
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return &user, nil
}

// FindByEmail はメールアドレスでユーザーを取得する
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, name, picture, google_id, refresh_token, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user model.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Picture,
		&user.GoogleID,
		&user.RefreshToken,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		r.logger.ErrorContext(ctx, "failed to query user by email", "email", email, "error", err)
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return &user, nil
}

// FindByGoogleID はGoogle IDでユーザーを取得する
func (r *UserRepositoryImpl) FindByGoogleID(ctx context.Context, googleID string) (*model.User, error) {
	query := `
		SELECT id, email, name, picture, google_id, refresh_token, created_at, updated_at
		FROM users
		WHERE google_id = $1
	`

	var user model.User
	err := r.db.QueryRowContext(ctx, query, googleID).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Picture,
		&user.GoogleID,
		&user.RefreshToken,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		r.logger.ErrorContext(ctx, "failed to query user by google_id", "google_id", googleID, "error", err)
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return &user, nil
}

// Update はユーザー情報を更新する
func (r *UserRepositoryImpl) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET email = $2, name = $3, picture = $4, refresh_token = $5, updated_at = $6
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Name,
		user.Picture,
		user.RefreshToken,
		user.UpdatedAt,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to update user", "id", user.ID, "error", err)
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to get rows affected", "error", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return model.ErrNotFound
	}

	return nil
}

// Delete はユーザーを削除する
func (r *UserRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to delete user", "id", id, "error", err)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to get rows affected", "error", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return model.ErrNotFound
	}

	return nil
}
