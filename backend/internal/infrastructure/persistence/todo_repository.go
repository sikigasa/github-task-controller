package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/sikigasa/github-task-controller/backend/internal/domain/model"
	"github.com/sikigasa/github-task-controller/backend/internal/domain/repository"
)

// TodoRepositoryImpl はTodoRepositoryの実装
type TodoRepositoryImpl struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewTodoRepository は新しいTodoRepositoryImplを作成する
func NewTodoRepository(db *sql.DB, logger *slog.Logger) repository.TodoRepository {
	return &TodoRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

// Create は新しいTODOを作成する
func (r *TodoRepositoryImpl) Create(ctx context.Context, todo *model.Todo) error {
	query := `
		INSERT INTO todos (id, title, description, completed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		todo.ID,
		todo.Title,
		todo.Description,
		todo.Completed,
		todo.CreatedAt,
		todo.UpdatedAt,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to insert todo", "error", err)
		return fmt.Errorf("failed to insert todo: %w", err)
	}

	return nil
}

// FindByID はIDでTODOを取得する
func (r *TodoRepositoryImpl) FindByID(ctx context.Context, id string) (*model.Todo, error) {
	query := `
		SELECT id, title, description, completed, created_at, updated_at
		FROM todos
		WHERE id = $1
	`

	var todo model.Todo
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		r.logger.ErrorContext(ctx, "failed to query todo", "id", id, "error", err)
		return nil, fmt.Errorf("failed to query todo: %w", err)
	}

	return &todo, nil
}

// FindAll はすべてのTODOを取得する
func (r *TodoRepositoryImpl) FindAll(ctx context.Context) ([]*model.Todo, error) {
	query := `
		SELECT id, title, description, completed, created_at, updated_at
		FROM todos
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to query todos", "error", err)
		return nil, fmt.Errorf("failed to query todos: %w", err)
	}
	defer rows.Close()

	var todos []*model.Todo
	for rows.Next() {
		var todo model.Todo
		if err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		); err != nil {
			r.logger.ErrorContext(ctx, "failed to scan todo", "error", err)
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}
		todos = append(todos, &todo)
	}

	if err := rows.Err(); err != nil {
		r.logger.ErrorContext(ctx, "rows error", "error", err)
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return todos, nil
}

// Update はTODOを更新する
func (r *TodoRepositoryImpl) Update(ctx context.Context, todo *model.Todo) error {
	query := `
		UPDATE todos
		SET title = $2, description = $3, completed = $4, updated_at = $5
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		todo.ID,
		todo.Title,
		todo.Description,
		todo.Completed,
		todo.UpdatedAt,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to update todo", "id", todo.ID, "error", err)
		return fmt.Errorf("failed to update todo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to get rows affected", "error", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return model.ErrNotFound
	}

	return nil
}

// Delete はTODOを削除する
func (r *TodoRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM todos WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to delete todo", "id", id, "error", err)
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to get rows affected", "error", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return model.ErrNotFound
	}

	return nil
}
