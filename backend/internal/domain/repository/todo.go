package repository

import (
	"context"

	"github.com/sikigasa/github-task-controller/backend/internal/domain/model"
)

// TodoRepository はTODOのデータアクセスを抽象化するインターフェース
type TodoRepository interface {
	// Create は新しいTODOを作成する
	Create(ctx context.Context, todo *model.Todo) error

	// FindByID はIDでTODOを取得する
	FindByID(ctx context.Context, id string) (*model.Todo, error)

	// FindAll はすべてのTODOを取得する
	FindAll(ctx context.Context) ([]*model.Todo, error)

	// Update はTODOを更新する
	Update(ctx context.Context, todo *model.Todo) error

	// Delete はTODOを削除する
	Delete(ctx context.Context, id string) error
}
