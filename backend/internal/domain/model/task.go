package model

import "time"

// TaskStatus はタスクのステータスを表す
type TaskStatus int

const (
	TaskStatusTodo       TaskStatus = 0
	TaskStatusInProgress TaskStatus = 1
	TaskStatusDone       TaskStatus = 2
)

// TaskPriority はタスクの優先度を表す
type TaskPriority int

const (
	TaskPriorityLow    TaskPriority = 0
	TaskPriorityMedium TaskPriority = 1
	TaskPriorityHigh   TaskPriority = 2
)

// Task はタスクを表すドメインモデル
type Task struct {
	ID                string       `json:"id"`
	ProjectID         string       `json:"project_id"`
	Title             string       `json:"title"`
	Description       string       `json:"description"`
	Status            TaskStatus   `json:"status"`
	Priority          TaskPriority `json:"priority"`
	EndDate           *time.Time   `json:"end_date,omitempty"`
	GithubItemID      *string      `json:"github_item_id,omitempty"`
	GithubIssueNumber *int         `json:"github_issue_number,omitempty"`
	GithubIssueURL    *string      `json:"github_issue_url,omitempty"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
}

// HasGithubIssue はGitHub Issueが紐づいているかを返す
func (t *Task) HasGithubIssue() bool {
	return t.GithubIssueURL != nil && *t.GithubIssueURL != ""
}
