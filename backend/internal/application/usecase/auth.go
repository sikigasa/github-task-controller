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
	"github.com/sikigasa/github-task-controller/backend/internal/domain/model"
	"github.com/sikigasa/github-task-controller/backend/internal/domain/repository"
	"github.com/sikigasa/github-task-controller/backend/internal/infrastructure/auth"
	"golang.org/x/oauth2"
)

// AuthUsecase は認証に関するビジネスロジックを実装する
type AuthUsecase struct {
	userRepo          repository.UserRepository
	googleAccountRepo repository.GoogleAccountRepository
	githubAccountRepo repository.GithubAccountRepository
	oauthConfig       *auth.OAuthConfig
	logger            *slog.Logger
}

// NewAuthUsecase は新しいAuthUsecaseを作成する
func NewAuthUsecase(
	userRepo repository.UserRepository,
	googleAccountRepo repository.GoogleAccountRepository,
	githubAccountRepo repository.GithubAccountRepository,
	oauthConfig *auth.OAuthConfig,
	logger *slog.Logger,
) *AuthUsecase {
	return &AuthUsecase{
		userRepo:          userRepo,
		googleAccountRepo: googleAccountRepo,
		githubAccountRepo: githubAccountRepo,
		oauthConfig:       oauthConfig,
		logger:            logger,
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
func (u *AuthUsecase) GetAuthURL(provider string, state string) string {
	var providerType auth.ProviderType
	switch provider {
	case "google":
		providerType = auth.ProviderGoogle
	case "github":
		providerType = auth.ProviderGithub
	default:
		providerType = auth.ProviderGoogle
	}
	return u.oauthConfig.GetAuthURL(providerType, state)
}

// HandleCallback はOAuthコールバックを処理する
func (u *AuthUsecase) HandleCallback(ctx context.Context, provider string, code string) (*model.User, *oauth2.Token, error) {
	u.logger.InfoContext(ctx, "handling oauth callback", "provider", provider)

	var providerType auth.ProviderType
	switch provider {
	case "google":
		providerType = auth.ProviderGoogle
	case "github":
		providerType = auth.ProviderGithub
	default:
		return nil, nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	// トークンを取得
	token, err := u.oauthConfig.Exchange(ctx, providerType, code)
	if err != nil {
		u.logger.ErrorContext(ctx, "failed to exchange token", "provider", provider, "error", err)
		return nil, nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	switch providerType {
	case auth.ProviderGoogle:
		return u.handleGoogleCallback(ctx, token)
	case auth.ProviderGithub:
		return u.handleGithubCallback(ctx, token)
	default:
		return nil, nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// handleGoogleCallback はGoogleのOAuthコールバックを処理する
func (u *AuthUsecase) handleGoogleCallback(ctx context.Context, token *oauth2.Token) (*model.User, *oauth2.Token, error) {
	// ユーザー情報を取得
	googleUserInfo, err := u.oauthConfig.GetGoogleUserInfo(ctx, token)
	if err != nil {
		u.logger.ErrorContext(ctx, "failed to get google user info", "error", err)
		return nil, nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// メールが確認されていない場合はエラー
	if !googleUserInfo.VerifiedEmail {
		u.logger.WarnContext(ctx, "email not verified", "email", googleUserInfo.Email)
		return nil, nil, errors.New("email is not verified")
	}

	// 既存のGoogleアカウントを検索
	googleAccount, err := u.googleAccountRepo.FindByProviderAccountID(ctx, "google", googleUserInfo.ID)
	if err != nil && err.Error() != fmt.Sprintf("google account not found: %s", googleUserInfo.ID) {
		u.logger.ErrorContext(ctx, "failed to find google account", "error", err)
		return nil, nil, fmt.Errorf("failed to find google account: %w", err)
	}

	now := time.Now()
	var domainUser *model.User

	if googleAccount != nil {
		// 既存のユーザーを取得
		domainUser, err = u.userRepo.FindByID(ctx, googleAccount.UserID)
		if err != nil {
			u.logger.ErrorContext(ctx, "failed to find user", "user_id", googleAccount.UserID, "error", err)
			return nil, nil, fmt.Errorf("failed to find user: %w", err)
		}

		// ユーザー情報を更新
		domainUser.Name = googleUserInfo.Name
		domainUser.ImageURL = googleUserInfo.Picture
		domainUser.UpdatedAt = now

		if err := u.userRepo.Update(ctx, domainUser); err != nil {
			u.logger.ErrorContext(ctx, "failed to update user", "error", err)
			return nil, nil, fmt.Errorf("failed to update user: %w", err)
		}

		// Googleアカウント情報を更新
		googleAccount.AccessToken = token.AccessToken
		if token.RefreshToken != "" {
			googleAccount.RefreshToken = token.RefreshToken
		}
		if !token.Expiry.IsZero() {
			googleAccount.ExpiresAt = &token.Expiry
		}
		googleAccount.UpdatedAt = now

		if err := u.googleAccountRepo.Update(ctx, googleAccount); err != nil {
			u.logger.ErrorContext(ctx, "failed to update google account", "error", err)
			return nil, nil, fmt.Errorf("failed to update google account: %w", err)
		}
	} else {
		// 新規ユーザーの場合、メールで既存ユーザーを検索
		domainUser, err = u.userRepo.FindByEmail(ctx, googleUserInfo.Email)
		if err != nil && err.Error() != fmt.Sprintf("user not found: %s", googleUserInfo.Email) {
			u.logger.ErrorContext(ctx, "failed to find user by email", "error", err)
			return nil, nil, fmt.Errorf("failed to find user: %w", err)
		}

		if domainUser == nil {
			// 新規ユーザーを作成
			domainUser = &model.User{
				ID:        uuid.New().String(),
				Email:     googleUserInfo.Email,
				Name:      googleUserInfo.Name,
				ImageURL:  googleUserInfo.Picture,
				CreatedAt: now,
				UpdatedAt: now,
			}

			if err := u.userRepo.Create(ctx, domainUser); err != nil {
				u.logger.ErrorContext(ctx, "failed to create user", "error", err)
				return nil, nil, fmt.Errorf("failed to create user: %w", err)
			}

			u.logger.InfoContext(ctx, "user created successfully", "user_id", domainUser.ID)
		}

		// Googleアカウントを作成
		googleAccount = &model.GoogleAccount{
			ID:                uuid.New().String(),
			UserID:            domainUser.ID,
			Provider:          "google",
			ProviderAccountID: googleUserInfo.ID,
			AccessToken:       token.AccessToken,
			RefreshToken:      token.RefreshToken,
			CreatedAt:         now,
			UpdatedAt:         now,
		}
		if !token.Expiry.IsZero() {
			googleAccount.ExpiresAt = &token.Expiry
		}

		if err := u.googleAccountRepo.Create(ctx, googleAccount); err != nil {
			u.logger.ErrorContext(ctx, "failed to create google account", "error", err)
			return nil, nil, fmt.Errorf("failed to create google account: %w", err)
		}

		u.logger.InfoContext(ctx, "google account created successfully", "account_id", googleAccount.ID)
	}

	return domainUser, token, nil
}

// handleGithubCallback はGitHubのOAuthコールバックを処理する
func (u *AuthUsecase) handleGithubCallback(ctx context.Context, token *oauth2.Token) (*model.User, *oauth2.Token, error) {
	// ユーザー情報を取得
	githubUserInfo, err := u.oauthConfig.GetGithubUserInfo(ctx, token)
	if err != nil {
		u.logger.ErrorContext(ctx, "failed to get github user info", "error", err)
		return nil, nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// メールアドレスが取得できない場合はエラー
	if githubUserInfo.Email == "" {
		u.logger.WarnContext(ctx, "email not found", "login", githubUserInfo.Login)
		return nil, nil, errors.New("email is not available")
	}

	// 既存のGitHubアカウントを検索
	githubAccount, err := u.githubAccountRepo.FindByProviderAccountID(ctx, "github", fmt.Sprintf("%d", githubUserInfo.ID))
	if err != nil && err.Error() != fmt.Sprintf("github account not found: %d", githubUserInfo.ID) {
		u.logger.ErrorContext(ctx, "failed to find github account", "error", err)
		return nil, nil, fmt.Errorf("failed to find github account: %w", err)
	}

	now := time.Now()
	var domainUser *model.User

	if githubAccount != nil {
		// 既存のユーザーを取得
		domainUser, err = u.userRepo.FindByID(ctx, githubAccount.UserID)
		if err != nil {
			u.logger.ErrorContext(ctx, "failed to find user", "user_id", githubAccount.UserID, "error", err)
			return nil, nil, fmt.Errorf("failed to find user: %w", err)
		}

		// ユーザー情報を更新
		domainUser.Name = githubUserInfo.Name
		if domainUser.Name == "" {
			domainUser.Name = githubUserInfo.Login
		}
		domainUser.ImageURL = githubUserInfo.AvatarURL
		domainUser.UpdatedAt = now

		if err := u.userRepo.Update(ctx, domainUser); err != nil {
			u.logger.ErrorContext(ctx, "failed to update user", "error", err)
			return nil, nil, fmt.Errorf("failed to update user: %w", err)
		}

		// GitHubアカウント情報を更新
		githubAccount.AccessToken = token.AccessToken
		if token.RefreshToken != "" {
			githubAccount.RefreshToken = token.RefreshToken
		}
		if !token.Expiry.IsZero() {
			githubAccount.ExpiresAt = &token.Expiry
		}
		githubAccount.UpdatedAt = now

		if err := u.githubAccountRepo.Update(ctx, githubAccount); err != nil {
			u.logger.ErrorContext(ctx, "failed to update github account", "error", err)
			return nil, nil, fmt.Errorf("failed to update github account: %w", err)
		}
	} else {
		// 新規ユーザーの場合、メールで既存ユーザーを検索
		domainUser, err = u.userRepo.FindByEmail(ctx, githubUserInfo.Email)
		if err != nil && err.Error() != fmt.Sprintf("user not found: %s", githubUserInfo.Email) {
			u.logger.ErrorContext(ctx, "failed to find user by email", "error", err)
			return nil, nil, fmt.Errorf("failed to find user: %w", err)
		}

		if domainUser == nil {
			// 新規ユーザーを作成
			userName := githubUserInfo.Name
			if userName == "" {
				userName = githubUserInfo.Login
			}

			domainUser = &model.User{
				ID:        uuid.New().String(),
				Email:     githubUserInfo.Email,
				Name:      userName,
				ImageURL:  githubUserInfo.AvatarURL,
				CreatedAt: now,
				UpdatedAt: now,
			}

			if err := u.userRepo.Create(ctx, domainUser); err != nil {
				u.logger.ErrorContext(ctx, "failed to create user", "error", err)
				return nil, nil, fmt.Errorf("failed to create user: %w", err)
			}

			u.logger.InfoContext(ctx, "user created successfully", "user_id", domainUser.ID)
		}

		// GitHubアカウントを作成
		githubAccount = &model.GithubAccount{
			ID:                uuid.New().String(),
			UserID:            domainUser.ID,
			Provider:          "github",
			ProviderAccountID: fmt.Sprintf("%d", githubUserInfo.ID),
			AccessToken:       token.AccessToken,
			RefreshToken:      token.RefreshToken,
			CreatedAt:         now,
			UpdatedAt:         now,
		}
		if !token.Expiry.IsZero() {
			githubAccount.ExpiresAt = &token.Expiry
		}

		if err := u.githubAccountRepo.Create(ctx, githubAccount); err != nil {
			u.logger.ErrorContext(ctx, "failed to create github account", "error", err)
			return nil, nil, fmt.Errorf("failed to create github account: %w", err)
		}

		u.logger.InfoContext(ctx, "github account created successfully", "account_id", githubAccount.ID)
	}

	return domainUser, token, nil
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
		Picture:   user.ImageURL,
		ExpiresAt: time.Now().Add(expiresIn),
	}
}
