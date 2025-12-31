package router

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sikigasa/github-task-controller/backend/internal/interface/handler"
	"github.com/sikigasa/github-task-controller/backend/internal/interface/middleware"
)

// Router はアプリケーションのルーティングを管理する
type Router struct {
	mux            *mux.Router
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
		mux:            mux.NewRouter(),
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
	// ミドルウェアの設定
	r.mux.Use(r.loggingMiddleware)
	r.mux.Use(r.recoveryMiddleware)

	// ヘルスチェック
	r.mux.HandleFunc("/health", r.healthCheck).Methods(http.MethodGet)

	// 認証エンドポイント（認証不要）
	auth := r.mux.PathPrefix("/auth").Subrouter()
	// Google OAuth
	auth.HandleFunc("/google/login", r.authHandler.Login).Methods(http.MethodGet)
	auth.HandleFunc("/google/callback", r.authHandler.Callback).Methods(http.MethodGet)
	// GitHub OAuth
	auth.HandleFunc("/github/login", r.authHandler.LoginGithub).Methods(http.MethodGet)
	auth.HandleFunc("/github/callback", r.authHandler.CallbackGithub).Methods(http.MethodGet)
	// 共通
	auth.HandleFunc("/logout", r.authHandler.Logout).Methods(http.MethodPost)
	auth.HandleFunc("/me", r.authHandler.Me).Methods(http.MethodGet)

	// APIルーティング
	api := r.mux.PathPrefix("/api/v1").Subrouter()

	// 認証が必要なTODOエンドポイント
	protectedAPI := api.PathPrefix("").Subrouter()
	protectedAPI.Use(r.authMiddleware.RequireAuth)
	protectedAPI.HandleFunc("/todos", r.todoHandler.Create).Methods(http.MethodPost)
	protectedAPI.HandleFunc("/todos", r.todoHandler.List).Methods(http.MethodGet)
	protectedAPI.HandleFunc("/todos/{id}", r.todoHandler.Get).Methods(http.MethodGet)
	protectedAPI.HandleFunc("/todos/{id}", r.todoHandler.Update).Methods(http.MethodPut)
	protectedAPI.HandleFunc("/todos/{id}", r.todoHandler.Delete).Methods(http.MethodDelete)

	// プロジェクトエンドポイント
	protectedAPI.HandleFunc("/projects", r.projectHandler.Create).Methods(http.MethodPost)
	protectedAPI.HandleFunc("/projects", r.projectHandler.ListByUserID).Methods(http.MethodGet)
	protectedAPI.HandleFunc("/projects/{id}", r.projectHandler.Get).Methods(http.MethodGet)
	protectedAPI.HandleFunc("/projects/{id}", r.projectHandler.Update).Methods(http.MethodPut)
	protectedAPI.HandleFunc("/projects/{id}", r.projectHandler.Delete).Methods(http.MethodDelete)

	// タスクエンドポイント
	protectedAPI.HandleFunc("/tasks", r.taskHandler.Create).Methods(http.MethodPost)
	protectedAPI.HandleFunc("/tasks", r.taskHandler.ListByProjectID).Methods(http.MethodGet)
	protectedAPI.HandleFunc("/tasks/{id}", r.taskHandler.Get).Methods(http.MethodGet)
	protectedAPI.HandleFunc("/tasks/{id}", r.taskHandler.Update).Methods(http.MethodPut)
	protectedAPI.HandleFunc("/tasks/{id}", r.taskHandler.Delete).Methods(http.MethodDelete)

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

	return c.Handler(r.mux)
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
