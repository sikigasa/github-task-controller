package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/sikigasa/github-task-controller/backend/internal/domain/model"
	"github.com/sikigasa/github-task-controller/backend/internal/domain/repository"
)

type taskRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewTaskRepository は新しいTaskRepositoryを作成する
func NewTaskRepository(db *sql.DB, logger *slog.Logger) repository.TaskRepository {
	return &taskRepository{
		db:     db,
		logger: logger,
	}
}

func (r *taskRepository) Create(ctx context.Context, task *model.Task) error {
	query := `
		INSERT INTO task (id, project_id, title, description, status, end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		task.ID, task.ProjectID, task.Title, task.Description,
		task.Status, task.EndDate, task.CreatedAt, task.UpdatedAt,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to create task", "error", err)
		return fmt.Errorf("failed to create task: %w", err)
	}

	r.logger.InfoContext(ctx, "task created", "task_id", task.ID)
	return nil
}

func (r *taskRepository) FindByID(ctx context.Context, id string) (*model.Task, error) {
	query := `
		SELECT id, project_id, title, description, status, end_date, created_at, updated_at
		FROM task
		WHERE id = $1
	`

	var task model.Task
	var endDate sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID, &task.ProjectID, &task.Title, &task.Description,
		&task.Status, &endDate, &task.CreatedAt, &task.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task not found: %s", id)
	}
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to find task by id", "error", err, "id", id)
		return nil, fmt.Errorf("failed to find task by id: %w", err)
	}

	if endDate.Valid {
		task.EndDate = &endDate.Time
	}

	return &task, nil
}

func (r *taskRepository) FindByProjectID(ctx context.Context, projectID string) ([]*model.Task, error) {
	query := `
		SELECT id, project_id, title, description, status, end_date, created_at, updated_at
		FROM task
		WHERE project_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to find tasks by project_id", "error", err, "project_id", projectID)
		return nil, fmt.Errorf("failed to find tasks by project_id: %w", err)
	}
	defer rows.Close()

	var tasks []*model.Task
	for rows.Next() {
		var task model.Task
		var endDate sql.NullTime
		err := rows.Scan(
			&task.ID, &task.ProjectID, &task.Title, &task.Description,
			&task.Status, &endDate, &task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			r.logger.ErrorContext(ctx, "failed to scan task", "error", err)
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		if endDate.Valid {
			task.EndDate = &endDate.Time
		}

		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		r.logger.ErrorContext(ctx, "error iterating tasks", "error", err)
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}

func (r *taskRepository) Update(ctx context.Context, task *model.Task) error {
	query := `
		UPDATE task
		SET title = $1, description = $2, status = $3, end_date = $4, updated_at = $5
		WHERE id = $6
	`

	result, err := r.db.ExecContext(ctx, query,
		task.Title, task.Description, task.Status, task.EndDate, time.Now(), task.ID,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to update task", "error", err, "task_id", task.ID)
		return fmt.Errorf("failed to update task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task not found: %s", task.ID)
	}

	r.logger.InfoContext(ctx, "task updated", "task_id", task.ID)
	return nil
}

func (r *taskRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM task WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to delete task", "error", err, "task_id", id)
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task not found: %s", id)
	}

	r.logger.InfoContext(ctx, "task deleted", "task_id", id)
	return nil
}
