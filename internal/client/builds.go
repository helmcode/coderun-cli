package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// CreateBuild uploads a build context and creates a build
func (c *Client) CreateBuild(contextPath, appName, dockerfilePath string) (*BuildResponse, error) {
	// Open the context file
	file, err := os.Open(contextPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open context file: %w", err)
	}
	defer file.Close()

	// Create a buffer for the multipart form
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add the context file
	part, err := writer.CreateFormFile("context_file", filepath.Base(contextPath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// Add other form fields
	if err := writer.WriteField("app_name", appName); err != nil {
		return nil, fmt.Errorf("failed to write app_name field: %w", err)
	}

	if err := writer.WriteField("dockerfile_path", dockerfilePath); err != nil {
		return nil, fmt.Errorf("failed to write dockerfile_path field: %w", err)
	}

	// Close the writer to finalize the form
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", c.BaseURL+"/api/v1/builds/upload", &body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	// Make the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, handleAPIError(resp)
	}

	var buildResp BuildResponse
	if err := json.NewDecoder(resp.Body).Decode(&buildResp); err != nil {
		return nil, fmt.Errorf("failed to decode build response: %w", err)
	}

	return &buildResp, nil
}

// GetBuildStatus gets build status by ID
func (c *Client) GetBuildStatus(buildID string) (*BuildResponse, error) {
	resp, err := c.makeRequest("GET", "/api/v1/builds/"+buildID, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, handleAPIError(resp)
	}

	var buildResp BuildResponse
	if err := json.NewDecoder(resp.Body).Decode(&buildResp); err != nil {
		return nil, fmt.Errorf("failed to decode build response: %w", err)
	}

	return &buildResp, nil
}

// GetBuildLogs gets build logs by ID
func (c *Client) GetBuildLogs(buildID string) (string, error) {
	resp, err := c.makeRequest("GET", "/api/v1/builds/"+buildID+"/logs", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", handleAPIError(resp)
	}

	var logResp struct {
		Logs string `json:"logs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&logResp); err != nil {
		return "", fmt.Errorf("failed to decode build logs: %w", err)
	}

	return logResp.Logs, nil
}
