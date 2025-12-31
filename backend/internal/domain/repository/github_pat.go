package repository

import (
	"context"

	"github.com/sikigasa/github-task-controller/backend/internal/domain/model"
)

// GithubPATRepository はGitHub PATのリポジトリインターフェース
type GithubPATRepository interface {
	Create(ctx context.Context, pat *model.GithubPAT) error
	FindByUserID(ctx context.Context, userID string) (*model.GithubPAT, error)
	Update(ctx context.Context, pat *model.GithubPAT) error
	Delete(ctx context.Context, userID string) error
}
