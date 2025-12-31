package repository

import (
	"context"

	"github.com/sikigasa/github-task-controller/backend/internal/domain/model"
)

// UserRepository はユーザーのリポジトリインターフェース
type UserRepository interface {
	// Create は新しいユーザーを作成する
	Create(ctx context.Context, user *model.User) error
	// FindByID はIDでユーザーを検索する
	FindByID(ctx context.Context, id string) (*model.User, error)
	// FindByEmail はメールアドレスでユーザーを検索する
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	// Update はユーザー情報を更新する
	Update(ctx context.Context, user *model.User) error
	// Delete はユーザーを削除する
	Delete(ctx context.Context, id string) error
}
