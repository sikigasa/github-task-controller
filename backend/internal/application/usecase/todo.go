package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/sikigasa/github-task-controller/backend/internal/domain/model"
	"github.com/sikigasa/github-task-controller/backend/internal/domain/repository"
)

// TodoUsecase はTODOに関するビジネスロジックを実装する
type TodoUsecase struct {
	repo   repository.TodoRepository
	logger *slog.Logger
}

// NewTodoUsecase は新しいTodoUsecaseを作成する
func NewTodoUsecase(repo repository.TodoRepository, logger *slog.Logger) *TodoUsecase {
	return &TodoUsecase{
		repo:   repo,
		logger: logger,
	}
}

// Create は新しいTODOを作成する
func (u *TodoUsecase) Create(ctx context.Context, req *model.CreateTodoRequest) (*model.Todo, error) {
	u.logger.InfoContext(ctx, "creating new todo", "title", req.Title)

	todo := &model.Todo{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := u.repo.Create(ctx, todo); err != nil {
		u.logger.ErrorContext(ctx, "failed to create todo", "error", err)
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	u.logger.InfoContext(ctx, "todo created successfully", "id", todo.ID)
	return todo, nil
}

// GetByID はIDでTODOを取得する
func (u *TodoUsecase) GetByID(ctx context.Context, id string) (*model.Todo, error) {
	u.logger.InfoContext(ctx, "getting todo by id", "id", id)

	todo, err := u.repo.FindByID(ctx, id)
	if err != nil {
		u.logger.ErrorContext(ctx, "failed to get todo", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	return todo, nil
}

// GetAll はすべてのTODOを取得する
func (u *TodoUsecase) GetAll(ctx context.Context) ([]*model.Todo, error) {
	u.logger.InfoContext(ctx, "getting all todos")

	todos, err := u.repo.FindAll(ctx)
	if err != nil {
		u.logger.ErrorContext(ctx, "failed to get todos", "error", err)
		return nil, fmt.Errorf("failed to get todos: %w", err)
	}

	u.logger.InfoContext(ctx, "retrieved todos", "count", len(todos))
	return todos, nil
}

// Update はTODOを更新する
func (u *TodoUsecase) Update(ctx context.Context, id string, req *model.UpdateTodoRequest) (*model.Todo, error) {
	u.logger.InfoContext(ctx, "updating todo", "id", id)

	// 既存のTODOを取得
	todo, err := u.repo.FindByID(ctx, id)
	if err != nil {
		u.logger.ErrorContext(ctx, "failed to find todo for update", "id", id, "error", err)
		return nil, fmt.Errorf("failed to find todo: %w", err)
	}

	// リクエストに含まれるフィールドのみ更新
	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}
	todo.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, todo); err != nil {
		u.logger.ErrorContext(ctx, "failed to update todo", "id", id, "error", err)
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	u.logger.InfoContext(ctx, "todo updated successfully", "id", id)
	return todo, nil
}

// Delete はTODOを削除する
func (u *TodoUsecase) Delete(ctx context.Context, id string) error {
	u.logger.InfoContext(ctx, "deleting todo", "id", id)

	// 削除前に存在確認
	if _, err := u.repo.FindByID(ctx, id); err != nil {
		u.logger.ErrorContext(ctx, "failed to find todo for deletion", "id", id, "error", err)
		return fmt.Errorf("failed to find todo: %w", err)
	}

	if err := u.repo.Delete(ctx, id); err != nil {
		u.logger.ErrorContext(ctx, "failed to delete todo", "id", id, "error", err)
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	u.logger.InfoContext(ctx, "todo deleted successfully", "id", id)
	return nil
}
