package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

const (
	graphQLEndpoint = "https://api.github.com/graphql"
	restAPIBase     = "https://api.github.com"
)

// Client はGitHub APIクライアント
type Client struct {
	httpClient *http.Client
	logger     *slog.Logger
}

// NewClient は新しいGitHub APIクライアントを作成する
func NewClient(logger *slog.Logger) *Client {
	return &Client{
		httpClient: &http.Client{},
		logger:     logger,
	}
}

// GraphQLRequest はGraphQLリクエストを実行する
func (c *Client) GraphQLRequest(ctx context.Context, token, query string, variables map[string]interface{}) (map[string]interface{}, error) {
	body := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", graphQLEndpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.ErrorContext(ctx, "GitHub API error", "status", resp.StatusCode, "body", string(respBody))
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if errors, ok := result["errors"]; ok {
		c.logger.ErrorContext(ctx, "GraphQL errors", "errors", errors)
		return nil, fmt.Errorf("GraphQL errors: %v", errors)
	}

	return result, nil
}

// RESTRequest はREST APIリクエストを実行する
func (c *Client) RESTRequest(ctx context.Context, token, method, path string, body interface{}) (map[string]interface{}, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, restAPIBase+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		c.logger.ErrorContext(ctx, "GitHub REST API error", "status", resp.StatusCode, "body", string(respBody))
		return nil, fmt.Errorf("GitHub REST API error: %s", resp.Status)
	}

	if len(respBody) == 0 {
		return nil, nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}
