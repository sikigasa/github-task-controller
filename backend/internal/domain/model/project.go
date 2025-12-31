package model

import "time"

// Project はプロジェクトを表すドメインモデル
type Project struct {
	ID                  string    `json:"id"`
	UserID              string    `json:"user_id"`
	Title               string    `json:"title"`
	Description         string    `json:"description"`
	GithubOwner         *string   `json:"github_owner,omitempty"`
	GithubRepo          *string   `json:"github_repo,omitempty"`
	GithubProjectNumber *int      `json:"github_project_number,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// IsGithubLinked はGitHub連携が設定されているかを返す
func (p *Project) IsGithubLinked() bool {
	return p.GithubOwner != nil && p.GithubRepo != nil && p.GithubProjectNumber != nil
}
