package model

import "time"

// User はユーザー情報を表すドメインモデル
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Picture はImageURLのエイリアス（後方互換性のため）
func (u *User) Picture() string {
	return u.ImageURL
}
