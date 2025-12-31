package router

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/rs/cors"
	"github.com/sikigasa/github-task-controller/backend/internal/interface/handler"
	"github.com/sikigasa/github-task-controller/backend/internal/interface/middleware"
)

// Router はアプリケーションのルーティングを管理する
type Router struct {
	mux            *http.ServeMux
	todoHandler    *handler.TodoHandler
	projectHandler *handler.ProjectHandler
	taskHandler    *handler.TaskHandler
	authHandler    *handler.AuthHandler
	authMiddleware *middleware.AuthMiddleware
	logger         *slog.Logger
}

// NewRouter は新しいRouterを作成する
func NewRouter(
	todoHandler *handler.TodoHandler,
	projectHandler *handler.ProjectHandler,
	taskHandler *handler.TaskHandler,
	authHandler *handler.AuthHandler,
	authMiddleware *middleware.AuthMiddleware,
	logger *slog.Logger,
) *Router {
	return &Router{
		mux:            http.NewServeMux(),
		todoHandler:    todoHandler,
		projectHandler: projectHandler,
		taskHandler:    taskHandler,
		authHandler:    authHandler,
		authMiddleware: authMiddleware,
		logger:         logger,
	}
}

// Setup はルーティングを設定する
func (r *Router) Setup() http.Handler {
	// ヘルスチェック
	r.mux.HandleFunc("GET /health", r.healthCheck)

	// 認証エンドポイント（認証不要）
	// Google OAuth
	r.mux.HandleFunc("GET /auth/google/login", r.authHandler.Login)
	r.mux.HandleFunc("GET /auth/google/callback", r.authHandler.Callback)
	// GitHub OAuth
	r.mux.HandleFunc("GET /auth/github/login", r.authHandler.LoginGithub)
	r.mux.HandleFunc("GET /auth/github/callback", r.authHandler.CallbackGithub)
	// 共通
	r.mux.HandleFunc("POST /auth/logout", r.authHandler.Logout)
	r.mux.HandleFunc("GET /auth/me", r.authHandler.Me)

	// 認証が必要なAPIエンドポイント
	// TODOエンドポイント
	r.mux.Handle("POST /api/v1/todos", r.authMiddleware.RequireAuth(http.HandlerFunc(r.todoHandler.Create)))
	r.mux.Handle("GET /api/v1/todos", r.authMiddleware.RequireAuth(http.HandlerFunc(r.todoHandler.List)))
	r.mux.Handle("GET /api/v1/todos/{id}", r.authMiddleware.RequireAuth(http.HandlerFunc(r.todoHandler.Get)))
	r.mux.Handle("PUT /api/v1/todos/{id}", r.authMiddleware.RequireAuth(http.HandlerFunc(r.todoHandler.Update)))
	r.mux.Handle("DELETE /api/v1/todos/{id}", r.authMiddleware.RequireAuth(http.HandlerFunc(r.todoHandler.Delete)))

	// プロジェクトエンドポイント
	r.mux.Handle("POST /api/v1/projects", r.authMiddleware.RequireAuth(http.HandlerFunc(r.projectHandler.Create)))
	r.mux.Handle("GET /api/v1/projects", r.authMiddleware.RequireAuth(http.HandlerFunc(r.projectHandler.ListByUserID)))
	r.mux.Handle("GET /api/v1/projects/{id}", r.authMiddleware.RequireAuth(http.HandlerFunc(r.projectHandler.Get)))
	r.mux.Handle("PUT /api/v1/projects/{id}", r.authMiddleware.RequireAuth(http.HandlerFunc(r.projectHandler.Update)))
	r.mux.Handle("DELETE /api/v1/projects/{id}", r.authMiddleware.RequireAuth(http.HandlerFunc(r.projectHandler.Delete)))

	// タスクエンドポイント
	r.mux.Handle("POST /api/v1/tasks", r.authMiddleware.RequireAuth(http.HandlerFunc(r.taskHandler.Create)))
	r.mux.Handle("GET /api/v1/tasks", r.authMiddleware.RequireAuth(http.HandlerFunc(r.taskHandler.ListByProjectID)))
	r.mux.Handle("GET /api/v1/tasks/{id}", r.authMiddleware.RequireAuth(http.HandlerFunc(r.taskHandler.Get)))
	r.mux.Handle("PUT /api/v1/tasks/{id}", r.authMiddleware.RequireAuth(http.HandlerFunc(r.taskHandler.Update)))
	r.mux.Handle("DELETE /api/v1/tasks/{id}", r.authMiddleware.RequireAuth(http.HandlerFunc(r.taskHandler.Delete)))

	// CORS設定
	// credentials: 'include' を使用する場合、AllowedOrigins に "*" は使用不可
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Cookie"},
		ExposedHeaders:   []string{"Content-Length", "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	// ミドルウェアを適用
	var h http.Handler = r.mux
	h = r.loggingMiddleware(h)
	h = r.recoveryMiddleware(h)

	return c.Handler(h)
}

// healthCheck はヘルスチェックエンドポイント
func (r *Router) healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
		r.logger.ErrorContext(req.Context(), "failed to write health check response", "error", err)
	}
}

// loggingMiddleware はリクエストをログに記録するミドルウェア
func (r *Router) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()

		// レスポンスライターをラップして状態コードを取得
		wrapper := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapper, req)

		duration := time.Since(start)

		r.logger.InfoContext(req.Context(), "request completed",
			"method", req.Method,
			"path", req.URL.Path,
			"status", wrapper.statusCode,
			"duration_ms", duration.Milliseconds(),
			"remote_addr", req.RemoteAddr,
		)
	})
}

// recoveryMiddleware はpanicをキャッチして500エラーを返すミドルウェア
func (r *Router) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				r.logger.ErrorContext(req.Context(), "panic recovered",
					"error", err,
					"path", req.URL.Path,
				)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				if _, writeErr := w.Write([]byte(`{"type":"about:blank","title":"Internal Server Error","status":500,"detail":"予期しないエラーが発生しました"}`)); writeErr != nil {
					r.logger.ErrorContext(req.Context(), "failed to write error response", "error", writeErr)
				}
			}
		}()

		next.ServeHTTP(w, req)
	})
}

// responseWriter はhttp.ResponseWriterをラップして状態コードを記録する
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader は状態コードを記録する
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
