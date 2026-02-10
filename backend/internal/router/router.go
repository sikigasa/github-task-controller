package router

import (
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	githubHandler  *handler.GithubHandler
	authMiddleware *middleware.AuthMiddleware
	logger         *slog.Logger
	staticDir      string
	frontendURL    string
}

// NewRouter は新しいRouterを作成する
func NewRouter(
	todoHandler *handler.TodoHandler,
	projectHandler *handler.ProjectHandler,
	taskHandler *handler.TaskHandler,
	authHandler *handler.AuthHandler,
	githubHandler *handler.GithubHandler,
	authMiddleware *middleware.AuthMiddleware,
	frontendURL string,
	logger *slog.Logger,
) *Router {
	// 静的ファイルディレクトリ（環境変数で設定可能）
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = "../frontend/dist"
	}

	return &Router{
		mux:            http.NewServeMux(),
		todoHandler:    todoHandler,
		projectHandler: projectHandler,
		taskHandler:    taskHandler,
		authHandler:    authHandler,
		githubHandler:  githubHandler,
		authMiddleware: authMiddleware,
		logger:         logger,
		staticDir:      staticDir,
		frontendURL:    frontendURL,
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

	// GitHub連携エンドポイント
	r.mux.Handle("GET /api/v1/github/status", r.authMiddleware.RequireAuth(http.HandlerFunc(r.githubHandler.GetConnectionStatus)))
	r.mux.Handle("POST /api/v1/github/pat", r.authMiddleware.RequireAuth(http.HandlerFunc(r.githubHandler.SavePAT)))
	r.mux.Handle("DELETE /api/v1/github/pat", r.authMiddleware.RequireAuth(http.HandlerFunc(r.githubHandler.DeletePAT)))
	r.mux.Handle("GET /api/v1/github/projects", r.authMiddleware.RequireAuth(http.HandlerFunc(r.githubHandler.ListGithubProjects)))
	r.mux.Handle("POST /api/v1/projects/{id}/github/link", r.authMiddleware.RequireAuth(http.HandlerFunc(r.githubHandler.LinkProject)))
	r.mux.Handle("DELETE /api/v1/projects/{id}/github/link", r.authMiddleware.RequireAuth(http.HandlerFunc(r.githubHandler.UnlinkProject)))
	r.mux.Handle("POST /api/v1/tasks/{id}/github/sync", r.authMiddleware.RequireAuth(http.HandlerFunc(r.githubHandler.SyncTaskToGithub)))

	// SPA静的ファイル配信（本番環境用）
	r.mux.HandleFunc("/", r.spaHandler)

	// CORS設定
	// credentials: 'include' を使用する場合、AllowedOrigins に "*" は使用不可
	allowedOrigins := []string{"http://localhost:5173", "http://127.0.0.1:5173"}
	if r.frontendURL != "" && r.frontendURL != "http://localhost:5173" {
		allowedOrigins = append(allowedOrigins, r.frontendURL)
	}
	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
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

// spaHandler はSPA用の静的ファイル配信とfallbackを処理する
func (r *Router) spaHandler(w http.ResponseWriter, req *http.Request) {
	// 静的ファイルディレクトリが存在しない場合は404
	if _, err := os.Stat(r.staticDir); os.IsNotExist(err) {
		http.NotFound(w, req)
		return
	}

	path := req.URL.Path
	if path == "/" {
		path = "/index.html"
	}

	// ファイルパスを構築
	filePath := filepath.Join(r.staticDir, filepath.Clean(path))

	// セキュリティチェック: staticDir外へのアクセスを防止
	absStaticDir, _ := filepath.Abs(r.staticDir)
	absFilePath, _ := filepath.Abs(filePath)
	if !strings.HasPrefix(absFilePath, absStaticDir) {
		http.NotFound(w, req)
		return
	}

	// ファイルが存在するか確認
	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() {
		// ファイルが存在しない場合はindex.htmlを返す（SPA fallback）
		indexPath := filepath.Join(r.staticDir, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			http.ServeFile(w, req, indexPath)
			return
		}
		http.NotFound(w, req)
		return
	}

	// 静的ファイルを配信
	http.ServeFile(w, req, filePath)
}

// spaFileServer はSPA用のファイルサーバーを作成する（未使用だが参考用）
func (r *Router) spaFileServer(fsys fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(fsys))
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// APIパスは除外
		if strings.HasPrefix(req.URL.Path, "/api/") || strings.HasPrefix(req.URL.Path, "/auth/") {
			http.NotFound(w, req)
			return
		}
		fileServer.ServeHTTP(w, req)
	})
}
