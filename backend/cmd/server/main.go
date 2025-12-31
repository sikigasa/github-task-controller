package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sikigasa/github-task-controller/backend/cmd/config"
	"github.com/sikigasa/github-task-controller/backend/internal/application/usecase"
	"github.com/sikigasa/github-task-controller/backend/internal/infrastructure/auth"
	"github.com/sikigasa/github-task-controller/backend/internal/infrastructure/github"
	"github.com/sikigasa/github-task-controller/backend/internal/infrastructure/persistence"
	"github.com/sikigasa/github-task-controller/backend/internal/infrastructure/session"
	"github.com/sikigasa/github-task-controller/backend/internal/interface/handler"
	"github.com/sikigasa/github-task-controller/backend/internal/interface/middleware"
	"github.com/sikigasa/github-task-controller/backend/internal/router"
)

func main() {
	os.Exit(run())
}

func run() int {
	// ロガーの初期化
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	ctx := context.Background()

	// 環境変数の読み込み
	if err := config.LoadEnv(); err != nil {
		logger.Warn("failed to load .env file, using environment variables", "error", err)
	}

	// 設定の検証
	if config.Config.OAuth.Google.ClientID == "" || config.Config.OAuth.Google.ClientSecret == "" {
		logger.Error("GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET must be set")
		return 1
	}

	if config.Config.OAuth.Github.ClientID == "" || config.Config.OAuth.Github.ClientSecret == "" {
		logger.Error("GITHUB_CLIENT_ID and GITHUB_CLIENT_SECRET must be set")
		return 1
	}

	// データベース設定
	dbConfig := persistence.DBConfig{
		Host:     config.Config.Database.Host,
		Port:     config.Config.Database.Port,
		User:     config.Config.Database.User,
		Password: config.Config.Database.Password,
		DBName:   config.Config.Database.Name,
		SSLMode:  config.Config.Database.SSLMode,
	}

	// セッションストアの初期化
	sessionStore := session.NewCookieStore([]byte(config.Config.Session.Secret))

	// データベース接続
	db, err := persistence.NewDB(ctx, dbConfig, logger)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		return 1
	}
	defer db.Close()

	// スキーマ初期化
	if err := persistence.InitSchema(ctx, db, logger); err != nil {
		logger.Error("failed to initialize schema", "error", err)
		return 1
	}

	// OAuth設定の初期化
	oauthConfig := auth.NewOAuthConfig(
		config.Config.OAuth.Google.ClientID,
		config.Config.OAuth.Google.ClientSecret,
		config.Config.OAuth.Google.RedirectURL,
		config.Config.OAuth.Github.ClientID,
		config.Config.OAuth.Github.ClientSecret,
		config.Config.OAuth.Github.RedirectURL,
		logger,
	)

	// 依存性の注入
	todoRepo := persistence.NewTodoRepository(db, logger)
	userRepo := persistence.NewUserRepository(db, logger)
	googleAccountRepo := persistence.NewGoogleAccountRepository(db, logger)
	githubAccountRepo := persistence.NewGithubAccountRepository(db, logger)
	projectRepo := persistence.NewProjectRepository(db, logger)
	taskRepo := persistence.NewTaskRepository(db, logger)

	todoUsecase := usecase.NewTodoUsecase(todoRepo, logger)
	authUsecase := usecase.NewAuthUsecase(userRepo, googleAccountRepo, githubAccountRepo, oauthConfig, logger)
	projectUsecase := usecase.NewProjectUsecase(projectRepo, logger)
	taskUsecase := usecase.NewTaskUsecase(taskRepo, logger)

	// GitHub連携
	githubClient := github.NewClient(logger)
	githubService := github.NewProjectService(githubClient, logger)
	githubUsecase := usecase.NewGithubUsecase(githubAccountRepo, projectRepo, taskRepo, githubService, logger)

	todoHandler := handler.NewTodoHandler(todoUsecase, logger)
	authHandler := handler.NewAuthHandler(authUsecase, sessionStore, config.Config.App.FrontendURL, logger)
	projectHandler := handler.NewProjectHandler(projectUsecase, logger)
	taskHandler := handler.NewTaskHandler(taskUsecase, logger)
	githubHandler := handler.NewGithubHandler(githubUsecase, logger)

	authMiddleware := middleware.NewAuthMiddleware(sessionStore, logger)

	// ルーターのセットアップ
	r := router.NewRouter(todoHandler, projectHandler, taskHandler, authHandler, githubHandler, authMiddleware, logger)
	httpHandler := r.Setup()

	// サーバーの設定
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Config.App.Port),
		Handler:      httpHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// サーバーの起動
	go func() {
		logger.Info("starting server", "port", config.Config.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
		}
	}()

	// シグナル待機
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	// シャットダウンのタイムアウト設定
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
		return 1
	}

	logger.Info("server exited gracefully")
	return 0
}
