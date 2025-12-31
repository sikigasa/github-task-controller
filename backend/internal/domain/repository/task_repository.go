package repository

import (
	"context"

	"github.com/sikigasa/github-task-controller/backend/internal/domain/model"
)

// TaskRepository はタスクのリポジトリインターフェース
type TaskRepository interface {
	// Create は新しいタスクを作成する
	Create(ctx context.Context, task *model.Task) error
	// FindByID はIDでタスクを検索する
	FindByID(ctx context.Context, id string) (*model.Task, error)
	// FindByProjectID はプロジェクトIDで全タスクを検索する
	FindByProjectID(ctx context.Context, projectID string) ([]*model.Task, error)
	// Update はタスク情報を更新する
	Update(ctx context.Context, task *model.Task) error
	// Delete はタスクを削除する
	Delete(ctx context.Context, id string) error
}
