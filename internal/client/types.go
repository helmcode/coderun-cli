package client

import (
	"time"
)

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// DeploymentCreate represents a deployment creation request
type DeploymentCreate struct {
	AppName                   string            `json:"app_name"`
	Image                     string            `json:"image"`
	Replicas                  int               `json:"replicas"`
	CPULimit                  string            `json:"cpu_limit,omitempty"`
	MemoryLimit               string            `json:"memory_limit,omitempty"`
	CPURequest                string            `json:"cpu_request,omitempty"`
	MemoryRequest             string            `json:"memory_request,omitempty"`
	HTTPPort                  *int              `json:"http_port,omitempty"`
	TCPPort                   *int              `json:"tcp_port,omitempty"`
	EnvironmentVars           map[string]string `json:"environment_vars,omitempty"`
	PersistentVolumeSize      string            `json:"persistent_volume_size,omitempty"`
	PersistentVolumeMountPath string            `json:"persistent_volume_mount_path,omitempty"`
}

// DeploymentResponse represents a deployment response
type DeploymentResponse struct {
	ID                        string            `json:"id"`
	ClientID                  string            `json:"client_id"`
	AppName                   string            `json:"app_name"`
	Image                     string            `json:"image"`
	Replicas                  int               `json:"replicas"`
	CPULimit                  string            `json:"cpu_limit"`
	MemoryLimit               string            `json:"memory_limit"`
	CPURequest                string            `json:"cpu_request"`
	MemoryRequest             string            `json:"memory_request"`
	HTTPPort                  *int              `json:"http_port"`
	TCPPort                   *int              `json:"tcp_port"`
	TCPNodePort               *int              `json:"tcp_node_port"`
	EnvironmentVars           map[string]string `json:"environment_vars"`
	PersistentVolumeSize      string            `json:"persistent_volume_size,omitempty"`
	PersistentVolumeMountPath string            `json:"persistent_volume_mount_path,omitempty"`
	Status                    string            `json:"status"`
	CreatedAt                 time.Time         `json:"created_at"`
	UpdatedAt                 *time.Time        `json:"updated_at"`
	URL                       *string           `json:"url"`
	TCPConnection             *string           `json:"tcp_connection"`
}

// DeploymentList represents a list of deployments
type DeploymentList struct {
	Deployments []DeploymentResponse `json:"deployments"`
	Total       int                  `json:"total"`
}

// DeploymentStatus represents deployment status details
type DeploymentStatus struct {
	AppName                   string              `json:"app_name"`
	Status                    string              `json:"status"`
	ReplicasReady             int                 `json:"replicas_ready"`
	ReplicasDesired           int                 `json:"replicas_desired"`
	Pods                      []map[string]string `json:"pods"`
	URL                       *string             `json:"url"`
	TCPConnection             *string             `json:"tcp_connection"`
	URLNote                   *string             `json:"url_note,omitempty"`
	TLSCertificate            *TLSCertificateInfo `json:"tls_certificate,omitempty"`
	PersistentVolumeSize      string              `json:"persistent_volume_size,omitempty"`
	PersistentVolumeMountPath string              `json:"persistent_volume_mount_path,omitempty"`
}

// TLSCertificateInfo represents TLS certificate status
type TLSCertificateInfo struct {
	Ready   bool   `json:"ready"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// UserInfo represents user information
type UserInfo struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Namespace string    `json:"namespace"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// PodLogs represents logs from a single pod
type PodLogs struct {
	Status       string `json:"status"`
	Logs         string `json:"logs"`
	RestartCount int    `json:"restart_count"`
}

// DeploymentLogsResponse represents logs response from a deployment
type DeploymentLogsResponse struct {
	DeploymentID string             `json:"deployment_id"`
	AppName      string             `json:"app_name"`
	Image        string             `json:"image"`
	Status       string             `json:"status"`
	Deployment   string             `json:"deployment"`
	TotalPods    int                `json:"total_pods"`
	Error        string             `json:"error,omitempty"`
	Logs         map[string]PodLogs `json:"logs"`
}

// BuildRequest represents a build request (used for uploading context)
type BuildRequest struct {
	AppName        string `json:"app_name"`
	DockerfilePath string `json:"dockerfile_path"`
}

// BuildResponse represents a build response
type BuildResponse struct {
	ID             string     `json:"id"`
	ClientID       string     `json:"client_id"`
	AppName        string     `json:"app_name"`
	Tag            string     `json:"tag"`
	Status         string     `json:"status"`
	ImageURI       string     `json:"image_uri"`
	DockerfilePath string     `json:"dockerfile_path"`
	K8sJobName     string     `json:"k8s_job_name"`
	CreatedAt      time.Time  `json:"created_at"`
	StartedAt      *time.Time `json:"started_at"`
	FinishedAt     *time.Time `json:"finished_at"`
}
