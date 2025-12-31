package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/sikigasa/github-task-controller/backend/internal/domain/repository"
	"github.com/sikigasa/github-task-controller/backend/internal/infrastructure/auth"
	"github.com/sikigasa/github-task-controller/backend/internal/model"
	"golang.org/x/oauth2"
)

// AuthUsecase は認証に関するビジネスロジックを実装する
type AuthUsecase struct {
	userRepo    repository.UserRepository
	oauthConfig *auth.OAuthConfig
	logger      *slog.Logger
}

// NewAuthUsecase は新しいAuthUsecaseを作成する
func NewAuthUsecase(
	userRepo repository.UserRepository,
	oauthConfig *auth.OAuthConfig,
	logger *slog.Logger,
) *AuthUsecase {
	return &AuthUsecase{
		userRepo:    userRepo,
		oauthConfig: oauthConfig,
		logger:      logger,
	}
}

// GenerateStateToken はCSRF対策用のランダムな状態トークンを生成する
func (u *AuthUsecase) GenerateStateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		u.logger.Error("failed to generate state token", "error", err)
		return "", fmt.Errorf("failed to generate state token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GetAuthURL は認証URLを取得する
func (u *AuthUsecase) GetAuthURL(state string) string {
	return u.oauthConfig.GetAuthURL(state)
}

// HandleCallback はOAuthコールバックを処理する
func (u *AuthUsecase) HandleCallback(ctx context.Context, code string) (*model.User, *oauth2.Token, error) {
	u.logger.InfoContext(ctx, "handling oauth callback")

	// トークンを取得
	token, err := u.oauthConfig.Exchange(ctx, code)
	if err != nil {
		u.logger.ErrorContext(ctx, "failed to exchange token", "error", err)
		return nil, nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	// ユーザー情報を取得
	googleUserInfo, err := u.oauthConfig.GetUserInfo(ctx, token)
	if err != nil {
		u.logger.ErrorContext(ctx, "failed to get user info", "error", err)
		return nil, nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// メールが確認されていない場合はエラー
	if !googleUserInfo.VerifiedEmail {
		u.logger.WarnContext(ctx, "email not verified", "email", googleUserInfo.Email)
		return nil, nil, errors.New("email is not verified")
	}

	// 既存ユーザーを検索
	user, err := u.userRepo.FindByGoogleID(ctx, googleUserInfo.ID)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		u.logger.ErrorContext(ctx, "failed to find user by google_id", "error", err)
		return nil, nil, fmt.Errorf("failed to find user: %w", err)
	}

	now := time.Now()

	// 新規ユーザーの場合は作成
	if user == nil {
		u.logger.InfoContext(ctx, "creating new user", "email", googleUserInfo.Email)

		user = &model.User{
			ID:           uuid.New().String(),
			Email:        googleUserInfo.Email,
			Name:         googleUserInfo.Name,
			Picture:      googleUserInfo.Picture,
			GoogleID:     googleUserInfo.ID,
			RefreshToken: token.RefreshToken,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		if err := u.userRepo.Create(ctx, user); err != nil {
			u.logger.ErrorContext(ctx, "failed to create user", "error", err)
			return nil, nil, fmt.Errorf("failed to create user: %w", err)
		}

		u.logger.InfoContext(ctx, "user created successfully", "user_id", user.ID)
	} else {
		// 既存ユーザーの情報を更新
		u.logger.InfoContext(ctx, "updating existing user", "user_id", user.ID)

		user.Name = googleUserInfo.Name
		user.Picture = googleUserInfo.Picture
		if token.RefreshToken != "" {
			user.RefreshToken = token.RefreshToken
		}
		user.UpdatedAt = now

		if err := u.userRepo.Update(ctx, user); err != nil {
			u.logger.ErrorContext(ctx, "failed to update user", "error", err)
			return nil, nil, fmt.Errorf("failed to update user: %w", err)
		}

		u.logger.InfoContext(ctx, "user updated successfully", "user_id", user.ID)
	}

	return user, token, nil
}

// GetUserByID はIDでユーザーを取得する
func (u *AuthUsecase) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	u.logger.InfoContext(ctx, "getting user by id", "id", id)

	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		u.logger.ErrorContext(ctx, "failed to get user", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// CreateSession はセッション情報を作成する
func (u *AuthUsecase) CreateSession(user *model.User, expiresIn time.Duration) *model.Session {
	return &model.Session{
		UserID:    user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Picture:   user.Picture,
		ExpiresAt: time.Now().Add(expiresIn),
	}
}
