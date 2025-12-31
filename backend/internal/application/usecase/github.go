package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sikigasa/github-task-controller/backend/internal/domain/repository"
	"github.com/sikigasa/github-task-controller/backend/internal/infrastructure/github"
)

// GithubUsecase はGitHub連携のユースケース
type GithubUsecase struct {
	githubAccountRepo repository.GithubAccountRepository
	projectRepo       repository.ProjectRepository
	taskRepo          repository.TaskRepository
	githubService     *github.ProjectService
	logger            *slog.Logger
}

// NewGithubUsecase は新しいGithubUsecaseを作成する
func NewGithubUsecase(
	githubAccountRepo repository.GithubAccountRepository,
	projectRepo repository.ProjectRepository,
	taskRepo repository.TaskRepository,
	githubService *github.ProjectService,
	logger *slog.Logger,
) *GithubUsecase {
	return &GithubUsecase{
		githubAccountRepo: githubAccountRepo,
		projectRepo:       projectRepo,
		taskRepo:          taskRepo,
		githubService:     githubService,
		logger:            logger,
	}
}

// GithubConnectionStatus はGitHub連携状態を表す
type GithubConnectionStatus struct {
	IsConnected bool   `json:"is_connected"`
	HasPAT      bool   `json:"has_pat"`
	Username    string `json:"username,omitempty"`
}

// GetConnectionStatus はユーザーのGitHub連携状態を取得する
func (u *GithubUsecase) GetConnectionStatus(ctx context.Context, userID string) (*GithubConnectionStatus, error) {
	account, err := u.githubAccountRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find github account: %w", err)
	}

	if account == nil {
		return &GithubConnectionStatus{
			IsConnected: false,
			HasPAT:      false,
		}, nil
	}

	return &GithubConnectionStatus{
		IsConnected: true,
		HasPAT:      account.HasPAT(),
		Username:    account.ProviderAccountID,
	}, nil
}

// SavePAT はPATを保存する（簡易実装：本番では暗号化必須）
func (u *GithubUsecase) SavePAT(ctx context.Context, userID, pat string) error {
	account, err := u.githubAccountRepo.FindByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find github account: %w", err)
	}

	if account == nil {
		return fmt.Errorf("github account not found, please login with GitHub first")
	}

	// TODO: 本番環境では暗号化する
	account.PATEncrypted = &pat

	if err := u.githubAccountRepo.Update(ctx, account); err != nil {
		return fmt.Errorf("failed to update github account: %w", err)
	}

	u.logger.InfoContext(ctx, "PAT saved", "user_id", userID)
	return nil
}

// DeletePAT はPATを削除する
func (u *GithubUsecase) DeletePAT(ctx context.Context, userID string) error {
	account, err := u.githubAccountRepo.FindByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find github account: %w", err)
	}

	if account == nil {
		return fmt.Errorf("github account not found")
	}

	account.PATEncrypted = nil

	if err := u.githubAccountRepo.Update(ctx, account); err != nil {
		return fmt.Errorf("failed to update github account: %w", err)
	}

	u.logger.InfoContext(ctx, "PAT deleted", "user_id", userID)
	return nil
}

// GetToken はユーザーのGitHubトークンを取得する（PAT優先、なければOAuthトークン）
func (u *GithubUsecase) GetToken(ctx context.Context, userID string) (string, error) {
	account, err := u.githubAccountRepo.FindByUserID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to find github account: %w", err)
	}

	if account == nil {
		return "", fmt.Errorf("github account not found")
	}

	// PAT優先
	if account.HasPAT() {
		return *account.PATEncrypted, nil
	}

	// OAuthトークン
	if account.AccessToken != "" {
		return account.AccessToken, nil
	}

	return "", fmt.Errorf("no valid token found")
}

// ListGithubProjects はユーザーのGitHub Projectsを取得する
func (u *GithubUsecase) ListGithubProjects(ctx context.Context, userID string) ([]github.Project, error) {
	token, err := u.GetToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	projects, err := u.githubService.GetUserProjects(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get github projects: %w", err)
	}

	return projects, nil
}

// LinkProjectToGithub はプロジェクトをGitHub Projectに連携する
func (u *GithubUsecase) LinkProjectToGithub(ctx context.Context, userID, projectID, githubOwner, githubRepo string, githubProjectNumber int) error {
	project, err := u.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to find project: %w", err)
	}

	if project.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	project.GithubOwner = &githubOwner
	project.GithubRepo = &githubRepo
	project.GithubProjectNumber = &githubProjectNumber

	if err := u.projectRepo.Update(ctx, project); err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	u.logger.InfoContext(ctx, "project linked to github", "project_id", projectID, "github_project", githubProjectNumber)
	return nil
}

// UnlinkProjectFromGithub はプロジェクトのGitHub連携を解除する
func (u *GithubUsecase) UnlinkProjectFromGithub(ctx context.Context, userID, projectID string) error {
	project, err := u.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("failed to find project: %w", err)
	}

	if project.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	project.GithubOwner = nil
	project.GithubRepo = nil
	project.GithubProjectNumber = nil

	if err := u.projectRepo.Update(ctx, project); err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	u.logger.InfoContext(ctx, "project unlinked from github", "project_id", projectID)
	return nil
}

// SyncTaskToGithub はタスクをGitHub Projectに同期する
func (u *GithubUsecase) SyncTaskToGithub(ctx context.Context, userID, taskID string) error {
	task, err := u.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to find task: %w", err)
	}

	project, err := u.projectRepo.FindByID(ctx, task.ProjectID)
	if err != nil {
		return fmt.Errorf("failed to find project: %w", err)
	}

	if project.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	if !project.IsGithubLinked() {
		return fmt.Errorf("project is not linked to github")
	}

	token, err := u.GetToken(ctx, userID)
	if err != nil {
		return err
	}

	// GitHub Project IDを取得
	projectGithubID, err := u.githubService.GetProjectID(ctx, token, *project.GithubOwner, *project.GithubProjectNumber)
	if err != nil {
		return fmt.Errorf("failed to get github project id: %w", err)
	}

	// Draft Issueとして追加
	item, err := u.githubService.AddDraftIssueToProject(ctx, token, projectGithubID, task.Title, task.Description)
	if err != nil {
		return fmt.Errorf("failed to add task to github: %w", err)
	}

	// タスクにGitHub Item IDを保存
	task.GithubItemID = &item.ID
	if err := u.taskRepo.Update(ctx, task); err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	u.logger.InfoContext(ctx, "task synced to github", "task_id", taskID, "github_item_id", item.ID)
	return nil
}
