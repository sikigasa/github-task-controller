package model

import "time"

// User はユーザー情報を表すドメインモデル
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Picture      string    `json:"picture"`
	GoogleID     string    `json:"google_id"`
	RefreshToken string    `json:"-"` // JSONには含めない
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// GoogleUserInfo はGoogle OAuthから取得するユーザー情報
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

// Session はセッション情報を表す
type Session struct {
	UserID    string
	Email     string
	Name      string
	Picture   string
	ExpiresAt time.Time
}
