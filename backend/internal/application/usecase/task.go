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

// TaskUsecase はタスクに関するユースケース
type TaskUsecase struct {
	taskRepo repository.TaskRepository
	logger   *slog.Logger
}

// NewTaskUsecase は新しいTaskUsecaseを作成する
func NewTaskUsecase(taskRepo repository.TaskRepository, logger *slog.Logger) *TaskUsecase {
	return &TaskUsecase{
		taskRepo: taskRepo,
		logger:   logger,
	}
}

// CreateTask は新しいタスクを作成する
func (u *TaskUsecase) CreateTask(ctx context.Context, projectID, title, description string, status model.TaskStatus, endDate *time.Time) (*model.Task, error) {
	now := time.Now()
	task := &model.Task{
		ID:          uuid.New().String(),
		ProjectID:   projectID,
		Title:       title,
		Description: description,
		Status:      status,
		EndDate:     endDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := u.taskRepo.Create(ctx, task); err != nil {
		u.logger.ErrorContext(ctx, "failed to create task", "error", err)
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	u.logger.InfoContext(ctx, "task created", "task_id", task.ID, "project_id", projectID)
	return task, nil
}

// GetTask はIDでタスクを取得する
func (u *TaskUsecase) GetTask(ctx context.Context, id string) (*model.Task, error) {
	task, err := u.taskRepo.FindByID(ctx, id)
	if err != nil {
		u.logger.ErrorContext(ctx, "failed to get task", "error", err, "task_id", id)
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return task, nil
}

// ListTasksByProjectID はプロジェクトIDで全タスクを取得する
func (u *TaskUsecase) ListTasksByProjectID(ctx context.Context, projectID string) ([]*model.Task, error) {
	tasks, err := u.taskRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		u.logger.ErrorContext(ctx, "failed to list tasks", "error", err, "project_id", projectID)
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	return tasks, nil
}

// UpdateTask はタスク情報を更新する
func (u *TaskUsecase) UpdateTask(ctx context.Context, id, title, description string, status model.TaskStatus, endDate *time.Time) (*model.Task, error) {
	task, err := u.taskRepo.FindByID(ctx, id)
	if err != nil {
		u.logger.ErrorContext(ctx, "failed to find task", "error", err, "task_id", id)
		return nil, fmt.Errorf("failed to find task: %w", err)
	}

	task.Title = title
	task.Description = description
	task.Status = status
	task.EndDate = endDate
	task.UpdatedAt = time.Now()

	if err := u.taskRepo.Update(ctx, task); err != nil {
		u.logger.ErrorContext(ctx, "failed to update task", "error", err, "task_id", id)
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	u.logger.InfoContext(ctx, "task updated", "task_id", id)
	return task, nil
}

// DeleteTask はタスクを削除する
func (u *TaskUsecase) DeleteTask(ctx context.Context, id string) error {
	if err := u.taskRepo.Delete(ctx, id); err != nil {
		u.logger.ErrorContext(ctx, "failed to delete task", "error", err, "task_id", id)
		return fmt.Errorf("failed to delete task: %w", err)
	}

	u.logger.InfoContext(ctx, "task deleted", "task_id", id)
	return nil
}
