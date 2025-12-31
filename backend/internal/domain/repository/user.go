package repository

import (
	"context"

	"github.com/sikigasa/github-task-controller/backend/internal/model"
)

// UserRepository はユーザーのデータアクセスを抽象化するインターフェース
type UserRepository interface {
	// Create は新しいユーザーを作成する
	Create(ctx context.Context, user *model.User) error

	// FindByID はIDでユーザーを取得する
	FindByID(ctx context.Context, id string) (*model.User, error)

	// FindByEmail はメールアドレスでユーザーを取得する
	FindByEmail(ctx context.Context, email string) (*model.User, error)

	// FindByGoogleID はGoogle IDでユーザーを取得する
	FindByGoogleID(ctx context.Context, googleID string) (*model.User, error)

	// Update はユーザー情報を更新する
	Update(ctx context.Context, user *model.User) error

	// Delete はユーザーを削除する
	Delete(ctx context.Context, id string) error
}
