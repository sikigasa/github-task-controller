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
	authHandler    *handler.AuthHandler
	authMiddleware *middleware.AuthMiddleware
	logger         *slog.Logger
}

// NewRouter は新しいRouterを作成する
func NewRouter(
	todoHandler *handler.TodoHandler,
	authHandler *handler.AuthHandler,
	authMiddleware *middleware.AuthMiddleware,
	logger *slog.Logger,
) *Router {
	return &Router{
		mux:            mux.NewRouter(),
		todoHandler:    todoHandler,
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
	auth.HandleFunc("/login", r.authHandler.Login).Methods(http.MethodGet)
	auth.HandleFunc("/callback", r.authHandler.Callback).Methods(http.MethodGet)
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

	// CORS設定
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // 本番環境では適切に設定すること
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	return c.Handler(r.mux)
}

// healthCheck はヘルスチェックエンドポイント
func (r *Router) healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
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
				w.Write([]byte(`{"type":"about:blank","title":"Internal Server Error","status":500,"detail":"予期しないエラーが発生しました"}`))
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
