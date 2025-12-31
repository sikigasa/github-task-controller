package repository

import (
	"context"

	"github.com/sikigasa/github-task-controller/backend/internal/domain/model"
)

// GithubAccountRepository はGitHubアカウントのリポジトリインターフェース
type GithubAccountRepository interface {
	// Create は新しいGitHubアカウント情報を作成する
	Create(ctx context.Context, account *model.GithubAccount) error
	// FindByProviderAccountID はプロバイダーアカウントIDで検索する
	FindByProviderAccountID(ctx context.Context, provider, providerAccountID string) (*model.GithubAccount, error)
	// FindByUserID はユーザーIDで検索する
	FindByUserID(ctx context.Context, userID string) (*model.GithubAccount, error)
	// Update はGitHubアカウント情報を更新する
	Update(ctx context.Context, account *model.GithubAccount) error
	// Delete はGitHubアカウント情報を削除する
	Delete(ctx context.Context, provider, providerAccountID string) error
}

// GoogleAccountRepository はGoogleアカウントのリポジトリインターフェース
type GoogleAccountRepository interface {
	// Create は新しいGoogleアカウント情報を作成する
	Create(ctx context.Context, account *model.GoogleAccount) error
	// FindByProviderAccountID はプロバイダーアカウントIDで検索する
	FindByProviderAccountID(ctx context.Context, provider, providerAccountID string) (*model.GoogleAccount, error)
	// FindByUserID はユーザーIDで検索する
	FindByUserID(ctx context.Context, userID string) (*model.GoogleAccount, error)
	// Update はGoogleアカウント情報を更新する
	Update(ctx context.Context, account *model.GoogleAccount) error
	// Delete はGoogleアカウント情報を削除する
	Delete(ctx context.Context, provider, providerAccountID string) error
}
