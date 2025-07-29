package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// CreateDeployment creates a new deployment
func (c *Client) CreateDeployment(deployment *DeploymentCreate) (*DeploymentResponse, error) {
	resp, err := c.makeRequest("POST", "/api/v1/deploy", deployment)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, handleAPIError(resp)
	}

	var deploymentResp DeploymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&deploymentResp); err != nil {
		return nil, fmt.Errorf("failed to decode deployment response: %w", err)
	}

	return &deploymentResp, nil
}

// ListDeployments lists all deployments
func (c *Client) ListDeployments() (*DeploymentList, error) {
	resp, err := c.makeRequest("GET", "/api/v1/deployments", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, handleAPIError(resp)
	}

	var deploymentList DeploymentList
	if err := json.NewDecoder(resp.Body).Decode(&deploymentList); err != nil {
		return nil, fmt.Errorf("failed to decode deployment list: %w", err)
	}

	return &deploymentList, nil
}

// GetDeploymentStatus gets deployment status by ID
func (c *Client) GetDeploymentStatus(deploymentID string) (*DeploymentStatus, error) {
	endpoint := fmt.Sprintf("/api/v1/deployments/%s/status", deploymentID)

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, handleAPIError(resp)
	}

	var status DeploymentStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode deployment status: %w", err)
	}

	return &status, nil
}

// DeleteDeployment deletes a deployment by ID
func (c *Client) DeleteDeployment(deploymentID string) error {
	endpoint := fmt.Sprintf("/api/v1/deployments/%s", deploymentID)

	resp, err := c.makeRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return handleAPIError(resp)
	}

	return nil
}

// GetDeploymentLogs gets logs from a deployment by ID
func (c *Client) GetDeploymentLogs(deploymentID string, lines int) (*DeploymentLogsResponse, error) {
	endpoint := fmt.Sprintf("/api/v1/deployments/%s/logs?lines=%d", deploymentID, lines)

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, handleAPIError(resp)
	}

	var logs DeploymentLogsResponse
	if err := json.NewDecoder(resp.Body).Decode(&logs); err != nil {
		return nil, fmt.Errorf("failed to decode deployment logs: %w", err)
	}

	return &logs, nil
}
