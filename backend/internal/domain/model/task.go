package model

import "time"

// TaskStatus はタスクのステータスを表す
type TaskStatus int

const (
	TaskStatusTodo       TaskStatus = 0
	TaskStatusInProgress TaskStatus = 1
	TaskStatusDone       TaskStatus = 2
)

// Task はタスクを表すドメインモデル
type Task struct {
	ID          string     `json:"id"`
	ProjectID   string     `json:"project_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
