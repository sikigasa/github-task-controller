package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// OAuthConfig はOAuth認証の設定を保持する
type OAuthConfig struct {
	Config *oauth2.Config
	Logger *slog.Logger
}

// NewOAuthConfig は新しいOAuthConfigを作成する
func NewOAuthConfig(clientID, clientSecret, redirectURL string, logger *slog.Logger) *OAuthConfig {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &OAuthConfig{
		Config: config,
		Logger: logger,
	}
}

// GetAuthURL は認証URLを生成する
func (o *OAuthConfig) GetAuthURL(state string) string {
	return o.Config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

// Exchange は認証コードをトークンに交換する
func (o *OAuthConfig) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := o.Config.Exchange(ctx, code)
	if err != nil {
		o.Logger.ErrorContext(ctx, "failed to exchange token", "error", err)
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

// GetUserInfo はアクセストークンを使用してGoogleからユーザー情報を取得する
func (o *OAuthConfig) GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	client := o.Config.Client(ctx, token)

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
