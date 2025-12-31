package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// ProviderType はOAuthプロバイダーの種類
type ProviderType string

const (
	ProviderGoogle ProviderType = "google"
	ProviderGithub ProviderType = "github"
)

// OAuthConfig はOAuth認証の設定を保持する
type OAuthConfig struct {
	GoogleConfig *oauth2.Config
	GithubConfig *oauth2.Config
	Logger       *slog.Logger
}

// NewOAuthConfig は新しいOAuthConfigを作成する
func NewOAuthConfig(
	googleClientID, googleClientSecret, googleRedirectURL string,
	githubClientID, githubClientSecret, githubRedirectURL string,
	logger *slog.Logger,
) *OAuthConfig {
	googleConfig := &oauth2.Config{
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		RedirectURL:  googleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	githubConfig := &oauth2.Config{
		ClientID:     githubClientID,
		ClientSecret: githubClientSecret,
		RedirectURL:  githubRedirectURL,
		Scopes: []string{
			"user:email",
			"read:user",
		},
		Endpoint: github.Endpoint,
	}

	return &OAuthConfig{
		GoogleConfig: googleConfig,
		GithubConfig: githubConfig,
		Logger:       logger,
	}
}

// GetAuthURL は認証URLを生成する
func (o *OAuthConfig) GetAuthURL(provider ProviderType, state string) string {
	switch provider {
	case ProviderGoogle:
		return o.GoogleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	case ProviderGithub:
		return o.GithubConfig.AuthCodeURL(state)
	default:
		return ""
	}
}

// Exchange は認証コードをトークンに交換する
func (o *OAuthConfig) Exchange(ctx context.Context, provider ProviderType, code string) (*oauth2.Token, error) {
	var token *oauth2.Token
	var err error

	switch provider {
	case ProviderGoogle:
		token, err = o.GoogleConfig.Exchange(ctx, code)
	case ProviderGithub:
		token, err = o.GithubConfig.Exchange(ctx, code)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	if err != nil {
		o.Logger.ErrorContext(ctx, "failed to exchange token", "provider", provider, "error", err)
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}
	return token, nil
}

// GoogleUserInfo はGoogleから取得したユーザー情報
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// GithubUserInfo はGitHubから取得したユーザー情報
type GithubUserInfo struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Bio       string `json:"bio"`
}

// GithubEmail はGitHubのメールアドレス情報
type GithubEmail struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}

// GetGoogleUserInfo はアクセストークンを使用してGoogleからユーザー情報を取得する
func (o *OAuthConfig) GetGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	client := o.GoogleConfig.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		o.Logger.ErrorContext(ctx, "failed to get user info", "error", err)
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		o.Logger.ErrorContext(ctx, "google api returned non-200 status",
			"status", resp.StatusCode,
			"body", string(body))
		return nil, fmt.Errorf("google api returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		o.Logger.ErrorContext(ctx, "failed to read response body", "error", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		o.Logger.ErrorContext(ctx, "failed to unmarshal user info", "error", err)
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	return &userInfo, nil
}

// GetGithubUserInfo はアクセストークンを使用してGitHubからユーザー情報を取得する
func (o *OAuthConfig) GetGithubUserInfo(ctx context.Context, token *oauth2.Token) (*GithubUserInfo, error) {
	client := o.GithubConfig.Client(ctx, token)

	// ユーザー情報を取得
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		o.Logger.ErrorContext(ctx, "failed to get user info", "error", err)
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		o.Logger.ErrorContext(ctx, "github api returned non-200 status",
			"status", resp.StatusCode,
			"body", string(body))
		return nil, fmt.Errorf("github api returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		o.Logger.ErrorContext(ctx, "failed to read response body", "error", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo GithubUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		o.Logger.ErrorContext(ctx, "failed to unmarshal user info", "error", err)
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	// メールアドレスがない場合、別途取得
	if userInfo.Email == "" {
		email, err := o.getGithubPrimaryEmail(ctx, client)
		if err != nil {
			o.Logger.WarnContext(ctx, "failed to get github email", "error", err)
		} else {
			userInfo.Email = email
		}
	}

	return &userInfo, nil
}

// getGithubPrimaryEmail はGitHubから主要なメールアドレスを取得する
func (o *OAuthConfig) getGithubPrimaryEmail(ctx context.Context, client *http.Client) (string, error) {
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return "", fmt.Errorf("failed to get emails: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github api returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var emails []GithubEmail
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", fmt.Errorf("failed to unmarshal emails: %w", err)
	}

	// 認証済みかつプライマリなメールアドレスを探す
	for _, email := range emails {
		if email.Primary && email.Verified {
			return email.Email, nil
		}
	}

	// 認証済みのメールアドレスを探す
	for _, email := range emails {
		if email.Verified {
			return email.Email, nil
		}
	}

	return "", fmt.Errorf("no verified email found")
}
