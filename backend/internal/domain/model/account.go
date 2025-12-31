package model

import "time"

// GithubAccount はGitHubアカウント認証情報を表すドメインモデル
type GithubAccount struct {
	ID                string     `json:"id"`
	UserID            string     `json:"user_id"`
	Provider          string     `json:"provider"`
	ProviderAccountID string     `json:"provider_account_id"`
	AccessToken       string     `json:"access_token,omitempty"`
	RefreshToken      string     `json:"refresh_token,omitempty"`
	ExpiresAt         *time.Time `json:"expires_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// GoogleAccount はGoogleアカウント認証情報を表すドメインモデル
type GoogleAccount struct {
	ID                string     `json:"id"`
	UserID            string     `json:"user_id"`
	Provider          string     `json:"provider"`
	ProviderAccountID string     `json:"provider_account_id"`
	AccessToken       string     `json:"access_token,omitempty"`
	RefreshToken      string     `json:"refresh_token,omitempty"`
	ExpiresAt         *time.Time `json:"expires_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
