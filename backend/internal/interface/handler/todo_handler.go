package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/sikigasa/github-task-controller/backend/internal/application/usecase"
	"github.com/sikigasa/github-task-controller/backend/internal/domain/model"
)

// TodoHandler はTODOに関するHTTPリクエストを処理する
type TodoHandler struct {
	usecase *usecase.TodoUsecase
	logger  *slog.Logger
}

// NewTodoHandler は新しいTodoHandlerを作成する
func NewTodoHandler(usecase *usecase.TodoUsecase, logger *slog.Logger) *TodoHandler {
	return &TodoHandler{
		usecase: usecase,
		logger:  logger,
	}
}

// ProblemDetail はRFC 9457に準拠したエラーレスポンス
type ProblemDetail struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

// respondJSON はJSON形式でレスポンスを返す
func (h *TodoHandler) respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

// respondError はRFC 9457形式のエラーレスポンスを返す
func (h *TodoHandler) respondError(w http.ResponseWriter, r *http.Request, status int, title string, detail string) {
	problem := ProblemDetail{
		Type:     "about:blank",
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: r.URL.Path,
	}

	// ログレベルを適切に設定
	switch {
	case status >= 500:
		h.logger.Error("server error", "status", status, "title", title, "detail", detail, "path", r.URL.Path)
	case status == 401 || status == 403 || status == 409 || status == 429:
		h.logger.Warn("client error requiring attention", "status", status, "title", title, "path", r.URL.Path)
	default:
		h.logger.Info("client error", "status", status, "title", title, "path", r.URL.Path)
	}

	h.respondJSON(w, status, problem)
}

// Create はTODOを作成する
func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "Invalid Request", "リクエストボディが不正です")
		return
	}

	// バリデーション
	if req.Title == "" {
		h.respondError(w, r, http.StatusBadRequest, "Invalid Input", "タイトルは必須です")
		return
	}
	if len(req.Title) > 200 {
		h.respondError(w, r, http.StatusBadRequest, "Invalid Input", "タイトルは200文字以内にしてください")
		return
	}
	if len(req.Description) > 1000 {
		h.respondError(w, r, http.StatusBadRequest, "Invalid Input", "説明は1000文字以内にしてください")
		return
	}

	todo, err := h.usecase.Create(ctx, &req)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "Internal Server Error", "TODOの作成に失敗しました")
		return
	}

	h.respondJSON(w, http.StatusCreated, todo)
}

// Get はTODOを取得する
func (h *TodoHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	if id == "" {
		h.respondError(w, r, http.StatusBadRequest, "Invalid Request", "IDが指定されていません")
		return
	}

	todo, err := h.usecase.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			h.respondError(w, r, http.StatusNotFound, "Not Found", "指定されたTODOが見つかりません")
			return
		}
		h.respondError(w, r, http.StatusInternalServerError, "Internal Server Error", "TODOの取得に失敗しました")
		return
	}

	h.respondJSON(w, http.StatusOK, todo)
}

// List はすべてのTODOを取得する
func (h *TodoHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	todos, err := h.usecase.GetAll(ctx)
	if err != nil {
		h.respondError(w, r, http.StatusInternalServerError, "Internal Server Error", "TODOリストの取得に失敗しました")
		return
	}

	h.respondJSON(w, http.StatusOK, todos)
}

// Update はTODOを更新する
func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	if id == "" {
		h.respondError(w, r, http.StatusBadRequest, "Invalid Request", "IDが指定されていません")
		return
	}

	var req model.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "Invalid Request", "リクエストボディが不正です")
		return
	}

	// バリデーション
	if req.Title != nil && *req.Title == "" {
		h.respondError(w, r, http.StatusBadRequest, "Invalid Input", "タイトルは空にできません")
		return
	}
	if req.Title != nil && len(*req.Title) > 200 {
		h.respondError(w, r, http.StatusBadRequest, "Invalid Input", "タイトルは200文字以内にしてください")
		return
	}
	if req.Description != nil && len(*req.Description) > 1000 {
		h.respondError(w, r, http.StatusBadRequest, "Invalid Input", "説明は1000文字以内にしてください")
		return
	}

	todo, err := h.usecase.Update(ctx, id, &req)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			h.respondError(w, r, http.StatusNotFound, "Not Found", "指定されたTODOが見つかりません")
			return
		}
		h.respondError(w, r, http.StatusInternalServerError, "Internal Server Error", "TODOの更新に失敗しました")
		return
	}

	h.respondJSON(w, http.StatusOK, todo)
}

// Delete はTODOを削除する
func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")

	if id == "" {
		h.respondError(w, r, http.StatusBadRequest, "Invalid Request", "IDが指定されていません")
		return
	}

	if err := h.usecase.Delete(ctx, id); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			h.respondError(w, r, http.StatusNotFound, "Not Found", "指定されたTODOが見つかりません")
			return
		}
		h.respondError(w, r, http.StatusInternalServerError, "Internal Server Error", "TODOの削除に失敗しました")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
