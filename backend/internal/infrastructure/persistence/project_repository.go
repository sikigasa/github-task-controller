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

type projectRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewProjectRepository は新しいProjectRepositoryを作成する
func NewProjectRepository(db *sql.DB, logger *slog.Logger) repository.ProjectRepository {
	return &projectRepository{
		db:     db,
		logger: logger,
	}
}

func (r *projectRepository) Create(ctx context.Context, project *model.Project) error {
	query := `
		INSERT INTO project (id, user_id, title, description, github_owner, github_repo, github_project_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		project.ID, project.UserID, project.Title, project.Description,
		project.GithubOwner, project.GithubRepo, project.GithubProjectNumber,
		project.CreatedAt, project.UpdatedAt,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to create project", "error", err)
		return fmt.Errorf("failed to create project: %w", err)
	}

	r.logger.InfoContext(ctx, "project created", "project_id", project.ID)
	return nil
}

func (r *projectRepository) FindByID(ctx context.Context, id string) (*model.Project, error) {
	query := `
		SELECT id, user_id, title, description, github_owner, github_repo, github_project_number, created_at, updated_at
		FROM project
		WHERE id = $1
	`

	var project model.Project
	var githubOwner, githubRepo sql.NullString
	var githubProjectNumber sql.NullInt32
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&project.ID, &project.UserID, &project.Title, &project.Description,
		&githubOwner, &githubRepo, &githubProjectNumber,
		&project.CreatedAt, &project.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("project not found: %s", id)
	}
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to find project by id", "error", err, "id", id)
		return nil, fmt.Errorf("failed to find project by id: %w", err)
	}

	if githubOwner.Valid {
		project.GithubOwner = &githubOwner.String
	}
	if githubRepo.Valid {
		project.GithubRepo = &githubRepo.String
	}
	if githubProjectNumber.Valid {
		num := int(githubProjectNumber.Int32)
		project.GithubProjectNumber = &num
	}

	return &project, nil
}

func (r *projectRepository) FindByUserID(ctx context.Context, userID string) ([]*model.Project, error) {
	query := `
		SELECT id, user_id, title, description, github_owner, github_repo, github_project_number, created_at, updated_at
		FROM project
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to find projects by user_id", "error", err, "user_id", userID)
		return nil, fmt.Errorf("failed to find projects by user_id: %w", err)
	}
	defer rows.Close()

	var projects []*model.Project
	for rows.Next() {
		var project model.Project
		var githubOwner, githubRepo sql.NullString
		var githubProjectNumber sql.NullInt32
		err := rows.Scan(
			&project.ID, &project.UserID, &project.Title, &project.Description,
			&githubOwner, &githubRepo, &githubProjectNumber,
			&project.CreatedAt, &project.UpdatedAt,
		)
		if err != nil {
			r.logger.ErrorContext(ctx, "failed to scan project", "error", err)
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		if githubOwner.Valid {
			project.GithubOwner = &githubOwner.String
		}
		if githubRepo.Valid {
			project.GithubRepo = &githubRepo.String
		}
		if githubProjectNumber.Valid {
			num := int(githubProjectNumber.Int32)
			project.GithubProjectNumber = &num
		}
		projects = append(projects, &project)
	}

	if err = rows.Err(); err != nil {
		r.logger.ErrorContext(ctx, "error iterating projects", "error", err)
		return nil, fmt.Errorf("error iterating projects: %w", err)
	}

	return projects, nil
}

func (r *projectRepository) Update(ctx context.Context, project *model.Project) error {
	query := `
		UPDATE project
		SET title = $1, description = $2, github_owner = $3, github_repo = $4, github_project_number = $5, updated_at = $6
		WHERE id = $7
	`

	result, err := r.db.ExecContext(ctx, query,
		project.Title, project.Description,
		project.GithubOwner, project.GithubRepo, project.GithubProjectNumber,
		time.Now(), project.ID,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to update project", "error", err, "project_id", project.ID)
		return fmt.Errorf("failed to update project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("project not found: %s", project.ID)
	}

	r.logger.InfoContext(ctx, "project updated", "project_id", project.ID)
	return nil
}

func (r *projectRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM project WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to delete project", "error", err, "project_id", id)
		return fmt.Errorf("failed to delete project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("project not found: %s", id)
	}

	r.logger.InfoContext(ctx, "project deleted", "project_id", id)
	return nil
}
