package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/sikigasa/github-task-controller/backend/internal/application/usecase"
	"github.com/sikigasa/github-task-controller/backend/internal/domain/model"
	"github.com/sikigasa/github-task-controller/backend/internal/infrastructure/session"
)

const (
	sessionName         = "auth-session"
	sessionKeyUserID    = "user_id"
	sessionKeyEmail     = "email"
	sessionKeyName      = "name"
	sessionKeyPicture   = "picture"
	sessionKeyExpiresAt = "expires_at"
	oauthStateKey       = "oauth_state"
	sessionMaxAge       = 60 * 60 * 24 * 7 // 7日間
)

// AuthHandler は認証に関するHTTPリクエストを処理する
type AuthHandler struct {
	authUsecase  *usecase.AuthUsecase
	sessionStore *session.CookieStore
	frontendURL  string
	logger       *slog.Logger
}

// NewAuthHandler は新しいAuthHandlerを作成する
func NewAuthHandler(
	authUsecase *usecase.AuthUsecase,
	sessionStore *session.CookieStore,
	frontendURL string,
	logger *slog.Logger,
) *AuthHandler {
	return &AuthHandler{
		authUsecase:  authUsecase,
		sessionStore: sessionStore,
		frontendURL:  frontendURL,
		logger:       logger,
	}
}

// Login はGoogle OAuth認証を開始する
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.InfoContext(ctx, "starting google oauth login")

	// 状態トークンを生成
	state, err := h.authUsecase.GenerateStateToken()
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to generate state token", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// セッションに状態を保存
	sess, _ := h.sessionStore.Get(r, sessionName)
	sess.Set(oauthStateKey, state)
	if err := h.sessionStore.Save(w, r, sessionName, sess); err != nil {
		h.logger.ErrorContext(ctx, "failed to save session", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Google認証URLにリダイレクト
	authURL := h.authUsecase.GetAuthURL("google", state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// LoginGithub はGitHub OAuth認証を開始する
func (h *AuthHandler) LoginGithub(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.InfoContext(ctx, "starting github oauth login")

	// 状態トークンを生成
	state, err := h.authUsecase.GenerateStateToken()
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to generate state token", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// セッションに状態を保存
	sess, _ := h.sessionStore.Get(r, sessionName)
	sess.Set(oauthStateKey, state)
	if err := h.sessionStore.Save(w, r, sessionName, sess); err != nil {
		h.logger.ErrorContext(ctx, "failed to save session", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// GitHub認証URLにリダイレクト
	authURL := h.authUsecase.GetAuthURL("github", state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// Callback はGoogle OAuth認証のコールバックを処理する
func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.InfoContext(ctx, "handling google oauth callback")

	// セッションから状態を取得
	sess, _ := h.sessionStore.Get(r, sessionName)
	savedState, ok := sess.GetString(oauthStateKey)
	if !ok || savedState == "" {
		h.logger.WarnContext(ctx, "state not found in session")
		http.Redirect(w, r, h.frontendURL+"?error=invalid_state", http.StatusTemporaryRedirect)
		return
	}

	// 状態を検証
	state := r.URL.Query().Get("state")
	if state != savedState {
		h.logger.WarnContext(ctx, "state mismatch", "expected", savedState, "got", state)
		http.Redirect(w, r, h.frontendURL+"?error=invalid_state", http.StatusTemporaryRedirect)
		return
	}

	// 認証コードを取得
	code := r.URL.Query().Get("code")
	if code == "" {
		h.logger.WarnContext(ctx, "code not found in query")
		http.Redirect(w, r, h.frontendURL+"?error=no_code", http.StatusTemporaryRedirect)
		return
	}

	// コールバックを処理してユーザー情報を取得
	user, _, err := h.authUsecase.HandleCallback(ctx, "google", code)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to handle callback", "error", err)
		http.Redirect(w, r, h.frontendURL+"?error=auth_failed", http.StatusTemporaryRedirect)
		return
	}

	// セッションにユーザー情報を保存
	sessionInfo := h.authUsecase.CreateSession(user, time.Duration(sessionMaxAge)*time.Second)
	sess.Set(sessionKeyUserID, sessionInfo.UserID)
	sess.Set(sessionKeyEmail, sessionInfo.Email)
	sess.Set(sessionKeyName, sessionInfo.Name)
	sess.Set(sessionKeyPicture, sessionInfo.Picture)
	sess.Set(sessionKeyExpiresAt, sessionInfo.ExpiresAt.Unix())
	sess.Delete(oauthStateKey)

	sess.Options.MaxAge = sessionMaxAge
	sess.Options.HttpOnly = true
	sess.Options.Secure = isHTTPS(r)
	if isHTTPS(r) {
		sess.Options.SameSite = http.SameSiteNoneMode
	} else {
		sess.Options.SameSite = http.SameSiteLaxMode
	}

	if err := h.sessionStore.Save(w, r, sessionName, sess); err != nil {
		h.logger.ErrorContext(ctx, "failed to save session", "error", err)
		http.Redirect(w, r, h.frontendURL+"?error=session_failed", http.StatusTemporaryRedirect)
		return
	}

	h.logger.InfoContext(ctx, "user logged in successfully", "user_id", user.ID)

	// フロントエンドにリダイレクト
	http.Redirect(w, r, h.frontendURL, http.StatusTemporaryRedirect)
}

// CallbackGithub はGitHub OAuth認証のコールバックを処理する
func (h *AuthHandler) CallbackGithub(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.InfoContext(ctx, "handling github oauth callback")

	// セッションから状態を取得
	sess, _ := h.sessionStore.Get(r, sessionName)
	savedState, ok := sess.GetString(oauthStateKey)
	if !ok || savedState == "" {
		h.logger.WarnContext(ctx, "state not found in session")
		http.Redirect(w, r, h.frontendURL+"?error=invalid_state", http.StatusTemporaryRedirect)
		return
	}

	// 状態を検証
	state := r.URL.Query().Get("state")
	if state != savedState {
		h.logger.WarnContext(ctx, "state mismatch", "expected", savedState, "got", state)
		http.Redirect(w, r, h.frontendURL+"?error=invalid_state", http.StatusTemporaryRedirect)
		return
	}

	// 認証コードを取得
	code := r.URL.Query().Get("code")
	if code == "" {
		h.logger.WarnContext(ctx, "code not found in query")
		http.Redirect(w, r, h.frontendURL+"?error=no_code", http.StatusTemporaryRedirect)
		return
	}

	// コールバックを処理してユーザー情報を取得
	user, _, err := h.authUsecase.HandleCallback(ctx, "github", code)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to handle callback", "error", err)
		http.Redirect(w, r, h.frontendURL+"?error=auth_failed", http.StatusTemporaryRedirect)
		return
	}

	// セッションにユーザー情報を保存
	sessionInfo := h.authUsecase.CreateSession(user, time.Duration(sessionMaxAge)*time.Second)
	sess.Set(sessionKeyUserID, sessionInfo.UserID)
	sess.Set(sessionKeyEmail, sessionInfo.Email)
	sess.Set(sessionKeyName, sessionInfo.Name)
	sess.Set(sessionKeyPicture, sessionInfo.Picture)
	sess.Set(sessionKeyExpiresAt, sessionInfo.ExpiresAt.Unix())
	sess.Delete(oauthStateKey)

	sess.Options.MaxAge = sessionMaxAge
	sess.Options.HttpOnly = true
	sess.Options.Secure = isHTTPS(r)
	if isHTTPS(r) {
		sess.Options.SameSite = http.SameSiteNoneMode
	} else {
		sess.Options.SameSite = http.SameSiteLaxMode
	}

	if err := h.sessionStore.Save(w, r, sessionName, sess); err != nil {
		h.logger.ErrorContext(ctx, "failed to save session", "error", err)
		http.Redirect(w, r, h.frontendURL+"?error=session_failed", http.StatusTemporaryRedirect)
		return
	}

	h.logger.InfoContext(ctx, "user logged in successfully", "user_id", user.ID)

	// フロントエンドにリダイレクト
	http.Redirect(w, r, h.frontendURL, http.StatusTemporaryRedirect)
}

// Logout はログアウト処理を行う
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.InfoContext(ctx, "logging out user")

	// セッションを削除
	h.sessionStore.Delete(w, sessionName)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "logged out successfully"}); err != nil {
		h.logger.ErrorContext(ctx, "failed to encode response", "error", err)
	}
}

// Me は現在ログイン中のユーザー情報を返す
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// セッションからユーザー情報を取得
	sess, _ := h.sessionStore.Get(r, sessionName)
	userID, ok := sess.GetString(sessionKeyUserID)
	if !ok || userID == "" {
		h.logger.InfoContext(ctx, "user not authenticated")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// セッション有効期限を確認
	if sess.IsExpired(sessionKeyExpiresAt) {
		h.logger.InfoContext(ctx, "session expired", "user_id", userID)
		h.sessionStore.Delete(w, sessionName)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ユーザー情報を取得
	user, err := h.authUsecase.GetUserByID(ctx, userID)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to get user", "user_id", userID, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// レスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      user.ID,
		"email":   user.Email,
		"name":    user.Name,
		"picture": user.ImageURL,
	}); err != nil {
		h.logger.ErrorContext(ctx, "failed to encode response", "error", err)
	}
}

// GetSessionFromRequest はリクエストからセッション情報を取得する
func (h *AuthHandler) GetSessionFromRequest(r *http.Request) (*model.Session, error) {
	sess, err := h.sessionStore.Get(r, sessionName)
	if err != nil {
		return nil, err
	}

	userID, ok := sess.GetString(sessionKeyUserID)
	if !ok || userID == "" {
		return nil, nil
	}

	if sess.IsExpired(sessionKeyExpiresAt) {
		return nil, nil
	}

	email, _ := sess.GetString(sessionKeyEmail)
	name, _ := sess.GetString(sessionKeyName)
	picture, _ := sess.GetString(sessionKeyPicture)
	expiresAt, _ := sess.GetInt64(sessionKeyExpiresAt)

	return &model.Session{
		UserID:    userID,
		Email:     email,
		Name:      name,
		Picture:   picture,
		ExpiresAt: time.Unix(expiresAt, 0),
	}, nil
}

// isHTTPS はリクエストがHTTPS経由かどうかを判定する
// プロキシ（Railway等）の場合はX-Forwarded-Protoヘッダーも確認する
func isHTTPS(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	return r.Header.Get("X-Forwarded-Proto") == "https"
}
