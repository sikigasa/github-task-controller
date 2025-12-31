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

// ProjectUsecase はプロジェクトに関するユースケース
type ProjectUsecase struct {
	projectRepo repository.ProjectRepository
	logger      *slog.Logger
}

// NewProjectUsecase は新しいProjectUsecaseを作成する
func NewProjectUsecase(projectRepo repository.ProjectRepository, logger *slog.Logger) *ProjectUsecase {
	return &ProjectUsecase{
		projectRepo: projectRepo,
		logger:      logger,
	}
}

// CreateProject は新しいプロジェクトを作成する
func (u *ProjectUsecase) CreateProject(ctx context.Context, userID, title, description string) (*model.Project, error) {
	now := time.Now()
	project := &model.Project{
		ID:          uuid.New().String(),
		UserID:      userID,
		Title:       title,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := u.projectRepo.Create(ctx, project); err != nil {
		u.logger.ErrorContext(ctx, "failed to create project", "error", err)
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	u.logger.InfoContext(ctx, "project created", "project_id", project.ID, "user_id", userID)
	return project, nil
}

// GetProject はIDでプロジェクトを取得する
func (u *ProjectUsecase) GetProject(ctx context.Context, id string) (*model.Project, error) {
	project, err := u.projectRepo.FindByID(ctx, id)
	if err != nil {
		u.logger.ErrorContext(ctx, "failed to get project", "error", err, "project_id", id)
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return project, nil
}

// ListProjectsByUserID はユーザーIDで全プロジェクトを取得する
func (u *ProjectUsecase) ListProjectsByUserID(ctx context.Context, userID string) ([]*model.Project, error) {
	projects, err := u.projectRepo.FindByUserID(ctx, userID)
	if err != nil {
		u.logger.ErrorContext(ctx, "failed to list projects", "error", err, "user_id", userID)
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}

	return projects, nil
}

// UpdateProject はプロジェクト情報を更新する
func (u *ProjectUsecase) UpdateProject(ctx context.Context, id, title, description string) (*model.Project, error) {
	project, err := u.projectRepo.FindByID(ctx, id)
	if err != nil {
		u.logger.ErrorContext(ctx, "failed to find project", "error", err, "project_id", id)
		return nil, fmt.Errorf("failed to find project: %w", err)
	}

	project.Title = title
	project.Description = description
	project.UpdatedAt = time.Now()

	if err := u.projectRepo.Update(ctx, project); err != nil {
		u.logger.ErrorContext(ctx, "failed to update project", "error", err, "project_id", id)
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	u.logger.InfoContext(ctx, "project updated", "project_id", id)
	return project, nil
}

// DeleteProject はプロジェクトを削除する
func (u *ProjectUsecase) DeleteProject(ctx context.Context, id string) error {
	if err := u.projectRepo.Delete(ctx, id); err != nil {
		u.logger.ErrorContext(ctx, "failed to delete project", "error", err, "project_id", id)
		return fmt.Errorf("failed to delete project: %w", err)
	}

	u.logger.InfoContext(ctx, "project deleted", "project_id", id)
	return nil
}
