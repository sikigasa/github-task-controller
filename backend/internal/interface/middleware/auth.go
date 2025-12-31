package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
)

// ContextKey はコンテキストキーの型
type ContextKey string

const (
	// UserIDKey はコンテキストからユーザーIDを取得するためのキー
	UserIDKey ContextKey = "user_id"
	// SessionKey はコンテキストからセッション情報を取得するためのキー
	SessionKey ContextKey = "session"
)

const (
	sessionName         = "auth-session"
	sessionKeyUserID    = "user_id"
	sessionKeyExpiresAt = "expires_at"
)

// AuthMiddleware は認証ミドルウェア
type AuthMiddleware struct {
	sessionStore sessions.Store
	logger       *slog.Logger
}

// NewAuthMiddleware は新しいAuthMiddlewareを作成する
func NewAuthMiddleware(sessionStore sessions.Store, logger *slog.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		sessionStore: sessionStore,
		logger:       logger,
	}
}

// RequireAuth は認証が必要なエンドポイント用のミドルウェア
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// セッションからユーザー情報を取得
		session, err := m.sessionStore.Get(r, sessionName)
		if err != nil {
			m.logger.ErrorContext(ctx, "failed to get session", "error", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, ok := session.Values[sessionKeyUserID].(string)
		if !ok || userID == "" {
			m.logger.InfoContext(ctx, "user not authenticated")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// セッション有効期限を確認
		expiresAt, ok := session.Values[sessionKeyExpiresAt].(int64)
		if !ok || time.Now().Unix() > expiresAt {
			m.logger.InfoContext(ctx, "session expired", "user_id", userID)
			session.Options.MaxAge = -1
			session.Save(r, w)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// コンテキストにユーザーIDを追加
		ctx = context.WithValue(ctx, UserIDKey, userID)
		ctx = context.WithValue(ctx, SessionKey, session.Values)

		m.logger.InfoContext(ctx, "user authenticated", "user_id", userID)

		// 次のハンドラーを実行
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth は認証がオプションのエンドポイント用のミドルウェア
// 認証情報があればコンテキストに追加するが、なくてもエラーにしない
func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// セッションからユーザー情報を取得
		session, err := m.sessionStore.Get(r, sessionName)
		if err != nil {
			// エラーがあっても続行
			next.ServeHTTP(w, r)
			return
		}

		userID, ok := session.Values[sessionKeyUserID].(string)
		if !ok || userID == "" {
			// ユーザーIDがなくても続行
			next.ServeHTTP(w, r)
			return
		}

		// セッション有効期限を確認
		expiresAt, ok := session.Values[sessionKeyExpiresAt].(int64)
		if !ok || time.Now().Unix() > expiresAt {
			// 期限切れでも続行
			next.ServeHTTP(w, r)
			return
		}

		// コンテキストにユーザーIDを追加
		ctx = context.WithValue(ctx, UserIDKey, userID)
		ctx = context.WithValue(ctx, SessionKey, session.Values)

		// 次のハンドラーを実行
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext はコンテキストからユーザーIDを取得する
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}
