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

	"github.com/sikigasa/github-task-controller/backend/internal/application/usecase"
	"github.com/sikigasa/github-task-controller/backend/internal/infrastructure/persistence"
	"github.com/sikigasa/github-task-controller/backend/internal/interface/handler"
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
	dbConfig := persistence.DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "todoapp"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

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

	// 依存性の注入
	todoRepo := persistence.NewTodoRepository(db, logger)
	todoUsecase := usecase.NewTodoUsecase(todoRepo, logger)
	todoHandler := handler.NewTodoHandler(todoUsecase, logger)

	// ルーターのセットアップ
	r := router.NewRouter(todoHandler, logger)
	handler := r.Setup()

	// サーバーの設定
	port := getEnv("PORT", "8080")
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// サーバーの起動
	go func() {
		logger.Info("starting server", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
		}
	}()

	// グレースフルシャットダウンの設定
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

// getEnv は環境変数を取得し、存在しない場合はデフォルト値を返す
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
