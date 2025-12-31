package model

import "time"

// Session はセッション情報を表す
type Session struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Picture   string    `json:"picture"`
	ExpiresAt time.Time `json:"expires_at"`
}
