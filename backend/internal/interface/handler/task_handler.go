package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/sikigasa/github-task-controller/backend/internal/application/usecase"
	"github.com/sikigasa/github-task-controller/backend/internal/domain/model"
)

// TaskHandler はタスクのHTTPハンドラー
type TaskHandler struct {
	usecase *usecase.TaskUsecase
	logger  *slog.Logger
}

// NewTaskHandler は新しいTaskHandlerを作成する
func NewTaskHandler(usecase *usecase.TaskUsecase, logger *slog.Logger) *TaskHandler {
	return &TaskHandler{
		usecase: usecase,
		logger:  logger,
	}
}

// CreateTaskRequest はタスク作成リクエスト
type CreateTaskRequest struct {
	ProjectID   string     `json:"project_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      int        `json:"status"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

// UpdateTaskRequest はタスク更新リクエスト
type UpdateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      int        `json:"status"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

// Create は新しいタスクを作成する
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.ErrorContext(ctx, "failed to decode request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ProjectID == "" || req.Title == "" {
		http.Error(w, "project_id and title are required", http.StatusBadRequest)
		return
	}

	task, err := h.usecase.CreateTask(ctx, req.ProjectID, req.Title, req.Description, model.TaskStatus(req.Status), req.EndDate)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to create task", "error", err)
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		h.logger.ErrorContext(ctx, "failed to encode response", "error", err)
	}
}

// Get はIDでタスクを取得する
func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	task, err := h.usecase.GetTask(ctx, id)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to get task", "error", err, "id", id)
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		h.logger.ErrorContext(ctx, "failed to encode response", "error", err)
	}
}

// ListByProjectID はプロジェクトIDで全タスクを取得する
func (h *TaskHandler) ListByProjectID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectID := r.URL.Query().Get("project_id")

	if projectID == "" {
		http.Error(w, "project_id is required", http.StatusBadRequest)
		return
	}

	tasks, err := h.usecase.ListTasksByProjectID(ctx, projectID)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to list tasks", "error", err, "project_id", projectID)
		http.Error(w, "Failed to list tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		h.logger.ErrorContext(ctx, "failed to encode response", "error", err)
	}
}

// Update はタスク情報を更新する
func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.ErrorContext(ctx, "failed to decode request", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	task, err := h.usecase.UpdateTask(ctx, id, req.Title, req.Description, model.TaskStatus(req.Status), req.EndDate)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to update task", "error", err, "id", id)
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		h.logger.ErrorContext(ctx, "failed to encode response", "error", err)
	}
}

// Delete はタスクを削除する
func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	if err := h.usecase.DeleteTask(ctx, id); err != nil {
		h.logger.ErrorContext(ctx, "failed to delete task", "error", err, "id", id)
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
