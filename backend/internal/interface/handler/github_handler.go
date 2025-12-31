package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sikigasa/github-task-controller/backend/internal/application/usecase"
	"github.com/sikigasa/github-task-controller/backend/internal/interface/middleware"
)

// GithubHandler はGitHub連携のHTTPハンドラー
type GithubHandler struct {
	usecase *usecase.GithubUsecase
	logger  *slog.Logger
}

// NewGithubHandler は新しいGithubHandlerを作成する
func NewGithubHandler(usecase *usecase.GithubUsecase, logger *slog.Logger) *GithubHandler {
	return &GithubHandler{
		usecase: usecase,
		logger:  logger,
	}
}

// GetConnectionStatus はGitHub連携状態を取得する
func (h *GithubHandler) GetConnectionStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := middleware.GetUserIDFromContext(ctx)

	status, err := h.usecase.GetConnectionStatus(ctx, userID)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to get connection status", "error", err)
		http.Error(w, "Failed to get connection status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		h.logger.ErrorContext(ctx, "failed to encode response", "error", err)
	}
}

// SavePATRequest はPAT保存リクエスト
type SavePATRequest struct {
	PAT string `json:"pat"`
}

// SavePAT はPATを保存する
func (h *GithubHandler) SavePAT(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := middleware.GetUserIDFromContext(ctx)

	var req SavePATRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.ErrorContext(ctx, "failed to decode request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.PAT == "" {
		http.Error(w, "PAT is required", http.StatusBadRequest)
		return
	}

	if err := h.usecase.SavePAT(ctx, userID, req.PAT); err != nil {
		h.logger.ErrorContext(ctx, "failed to save PAT", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeletePAT はPATを削除する
func (h *GithubHandler) DeletePAT(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := middleware.GetUserIDFromContext(ctx)

	if err := h.usecase.DeletePAT(ctx, userID); err != nil {
		h.logger.ErrorContext(ctx, "failed to delete PAT", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListGithubProjects はユーザーのGitHub Projectsを取得する
func (h *GithubHandler) ListGithubProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := middleware.GetUserIDFromContext(ctx)

	projects, err := h.usecase.ListGithubProjects(ctx, userID)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to list github projects", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(projects); err != nil {
		h.logger.ErrorContext(ctx, "failed to encode response", "error", err)
	}
}

// LinkProjectRequest はプロジェクト連携リクエスト
type LinkProjectRequest struct {
	GithubOwner         string `json:"github_owner"`
	GithubRepo          string `json:"github_repo"`
	GithubProjectNumber int    `json:"github_project_number"`
}

// LinkProject はプロジェクトをGitHub Projectに連携する
func (h *GithubHandler) LinkProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := middleware.GetUserIDFromContext(ctx)
	projectID := r.PathValue("id")

	var req LinkProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.ErrorContext(ctx, "failed to decode request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.GithubOwner == "" || req.GithubProjectNumber == 0 {
		http.Error(w, "github_owner and github_project_number are required", http.StatusBadRequest)
		return
	}

	if err := h.usecase.LinkProjectToGithub(ctx, userID, projectID, req.GithubOwner, req.GithubRepo, req.GithubProjectNumber); err != nil {
		h.logger.ErrorContext(ctx, "failed to link project", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UnlinkProject はプロジェクトのGitHub連携を解除する
func (h *GithubHandler) UnlinkProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := middleware.GetUserIDFromContext(ctx)
	projectID := r.PathValue("id")

	if err := h.usecase.UnlinkProjectFromGithub(ctx, userID, projectID); err != nil {
		h.logger.ErrorContext(ctx, "failed to unlink project", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SyncTaskToGithub はタスクをGitHub Projectに同期する
func (h *GithubHandler) SyncTaskToGithub(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := middleware.GetUserIDFromContext(ctx)
	taskID := r.PathValue("id")

	if err := h.usecase.SyncTaskToGithub(ctx, userID, taskID); err != nil {
		h.logger.ErrorContext(ctx, "failed to sync task", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
