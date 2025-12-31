package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sikigasa/github-task-controller/backend/internal/application/usecase"
)

// ProjectHandler はプロジェクトのHTTPハンドラー
type ProjectHandler struct {
	usecase *usecase.ProjectUsecase
	logger  *slog.Logger
}

// NewProjectHandler は新しいProjectHandlerを作成する
func NewProjectHandler(usecase *usecase.ProjectUsecase, logger *slog.Logger) *ProjectHandler {
	return &ProjectHandler{
		usecase: usecase,
		logger:  logger,
	}
}

// CreateProjectRequest はプロジェクト作成リクエスト
type CreateProjectRequest struct {
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// UpdateProjectRequest はプロジェクト更新リクエスト
type UpdateProjectRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Create は新しいプロジェクトを作成する
func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.ErrorContext(ctx, "failed to decode request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.Title == "" {
		http.Error(w, "user_id and title are required", http.StatusBadRequest)
		return
	}

	project, err := h.usecase.CreateProject(ctx, req.UserID, req.Title, req.Description)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to create project", "error", err)
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}

// Get はIDでプロジェクトを取得する
func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	project, err := h.usecase.GetProject(ctx, id)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to get project", "error", err, "id", id)
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

// ListByUserID はユーザーIDで全プロジェクトを取得する
func (h *ProjectHandler) ListByUserID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := r.URL.Query().Get("user_id")

	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	projects, err := h.usecase.ListProjectsByUserID(ctx, userID)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to list projects", "error", err, "user_id", userID)
		http.Error(w, "Failed to list projects", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

// Update はプロジェクト情報を更新する
func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.ErrorContext(ctx, "failed to decode request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	project, err := h.usecase.UpdateProject(ctx, id, req.Title, req.Description)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to update project", "error", err, "id", id)
		http.Error(w, "Failed to update project", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

// Delete はプロジェクトを削除する
func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.usecase.DeleteProject(ctx, id); err != nil {
		h.logger.ErrorContext(ctx, "failed to delete project", "error", err, "id", id)
		http.Error(w, "Failed to delete project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
