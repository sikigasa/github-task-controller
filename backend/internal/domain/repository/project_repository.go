package repository

import (
	"context"

	"github.com/sikigasa/github-task-controller/backend/internal/domain/model"
)

// ProjectRepository はプロジェクトのリポジトリインターフェース
type ProjectRepository interface {
	// Create は新しいプロジェクトを作成する
	Create(ctx context.Context, project *model.Project) error
	// FindByID はIDでプロジェクトを検索する
	FindByID(ctx context.Context, id string) (*model.Project, error)
	// FindByUserID はユーザーIDで全プロジェクトを検索する
	FindByUserID(ctx context.Context, userID string) ([]*model.Project, error)
	// Update はプロジェクト情報を更新する
	Update(ctx context.Context, project *model.Project) error
	// Delete はプロジェクトを削除する
	Delete(ctx context.Context, id string) error
}
