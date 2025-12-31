package model

import "time"

// GithubPAT はGitHub Personal Access Tokenを表すドメインモデル
type GithubPAT struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	TokenEncrypted string    `json:"-"` // JSONには含めない
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
