package model

import "time"

// Todo はTODOアイテムを表すドメインモデル
type Todo struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTodoRequest はTODO作成リクエストを表す
type CreateTodoRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=200"`
	Description string `json:"description" validate:"max=1000"`
}

// UpdateTodoRequest はTODO更新リクエストを表す
type UpdateTodoRequest struct {
	Title       *string `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	Completed   *bool   `json:"completed,omitempty"`
}
