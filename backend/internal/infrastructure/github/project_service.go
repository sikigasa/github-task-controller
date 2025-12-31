package github

import (
	"context"
	"fmt"
	"log/slog"
)

// ProjectItem はGitHub ProjectのItemを表す
type ProjectItem struct {
	ID          string
	Title       string
	Body        string
	Status      string
	IssueNumber *int
	IssueURL    *string
}

// Project はGitHub Projectを表す
type Project struct {
	ID     string
	Number int
	Title  string
}

// ProjectService はGitHub Projects V2のサービス
type ProjectService struct {
	client *Client
	logger *slog.Logger
}

// NewProjectService は新しいProjectServiceを作成する
func NewProjectService(client *Client, logger *slog.Logger) *ProjectService {
	return &ProjectService{
		client: client,
		logger: logger,
	}
}

// GetUserProjects はユーザーのProjectsを取得する
func (s *ProjectService) GetUserProjects(ctx context.Context, token string) ([]Project, error) {
	query := `
		query {
			viewer {
				projectsV2(first: 20) {
					nodes {
						id
						number
						title
					}
				}
			}
		}
	`

	result, err := s.client.GraphQLRequest(ctx, token, query, nil)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	viewer, ok := data["viewer"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid viewer format")
	}

	projectsV2, ok := viewer["projectsV2"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid projectsV2 format")
	}

	nodes, ok := projectsV2["nodes"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid nodes format")
	}

	var projects []Project
	for _, node := range nodes {
		n, ok := node.(map[string]interface{})
		if !ok {
			continue
		}
		projects = append(projects, Project{
			ID:     n["id"].(string),
			Number: int(n["number"].(float64)),
			Title:  n["title"].(string),
		})
	}

	return projects, nil
}

// GetProjectItems はProjectのItemsを取得する
func (s *ProjectService) GetProjectItems(ctx context.Context, token, owner string, projectNumber int) ([]ProjectItem, error) {
	query := `
		query($owner: String!, $number: Int!) {
			user(login: $owner) {
				projectV2(number: $number) {
					items(first: 100) {
						nodes {
							id
							content {
								... on Issue {
									title
									body
									number
									url
								}
								... on DraftIssue {
									title
									body
								}
							}
							fieldValueByName(name: "Status") {
								... on ProjectV2ItemFieldSingleSelectValue {
									name
								}
							}
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"owner":  owner,
		"number": projectNumber,
	}

	result, err := s.client.GraphQLRequest(ctx, token, query, variables)
	if err != nil {
		return nil, err
	}

	// レスポンスをパース
	items, err := s.parseProjectItems(result)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *ProjectService) parseProjectItems(result map[string]interface{}) ([]ProjectItem, error) {
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	user, ok := data["user"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid user format")
	}

	projectV2, ok := user["projectV2"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid projectV2 format")
	}

	itemsData, ok := projectV2["items"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid items format")
	}

	nodes, ok := itemsData["nodes"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid nodes format")
	}

	var items []ProjectItem
	for _, node := range nodes {
		n, ok := node.(map[string]interface{})
		if !ok {
			continue
		}

		item := ProjectItem{
			ID: n["id"].(string),
		}

		if content, ok := n["content"].(map[string]interface{}); ok {
			if title, ok := content["title"].(string); ok {
				item.Title = title
			}
			if body, ok := content["body"].(string); ok {
				item.Body = body
			}
			if number, ok := content["number"].(float64); ok {
				num := int(number)
				item.IssueNumber = &num
			}
			if url, ok := content["url"].(string); ok {
				item.IssueURL = &url
			}
		}

		if fieldValue, ok := n["fieldValueByName"].(map[string]interface{}); ok {
			if name, ok := fieldValue["name"].(string); ok {
				item.Status = name
			}
		}

		items = append(items, item)
	}

	return items, nil
}

// AddDraftIssueToProject はProjectにDraft Issueを追加する
func (s *ProjectService) AddDraftIssueToProject(ctx context.Context, token, projectID, title, body string) (*ProjectItem, error) {
	query := `
		mutation($projectId: ID!, $title: String!, $body: String) {
			addProjectV2DraftIssue(input: {projectId: $projectId, title: $title, body: $body}) {
				projectItem {
					id
				}
			}
		}
	`

	variables := map[string]interface{}{
		"projectId": projectID,
		"title":     title,
		"body":      body,
	}

	result, err := s.client.GraphQLRequest(ctx, token, query, variables)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	addResult, ok := data["addProjectV2DraftIssue"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid addProjectV2DraftIssue format")
	}

	projectItem, ok := addResult["projectItem"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid projectItem format")
	}

	return &ProjectItem{
		ID:    projectItem["id"].(string),
		Title: title,
		Body:  body,
	}, nil
}

// GetProjectID はowner/project_numberからProject IDを取得する
func (s *ProjectService) GetProjectID(ctx context.Context, token, owner string, projectNumber int) (string, error) {
	query := `
		query($owner: String!, $number: Int!) {
			user(login: $owner) {
				projectV2(number: $number) {
					id
				}
			}
		}
	`

	variables := map[string]interface{}{
		"owner":  owner,
		"number": projectNumber,
	}

	result, err := s.client.GraphQLRequest(ctx, token, query, variables)
	if err != nil {
		return "", err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format")
	}

	user, ok := data["user"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid user format")
	}

	projectV2, ok := user["projectV2"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("project not found")
	}

	return projectV2["id"].(string), nil
}

// DeleteProjectItem はProjectからItemを削除する
func (s *ProjectService) DeleteProjectItem(ctx context.Context, token, projectID, itemID string) error {
	query := `
		mutation($projectId: ID!, $itemId: ID!) {
			deleteProjectV2Item(input: {projectId: $projectId, itemId: $itemId}) {
				deletedItemId
			}
		}
	`

	variables := map[string]interface{}{
		"projectId": projectID,
		"itemId":    itemID,
	}

	_, err := s.client.GraphQLRequest(ctx, token, query, variables)
	return err
}
