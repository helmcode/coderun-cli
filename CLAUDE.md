# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Build
```bash
go build -o coderun .
```

### Run Tests
```bash
go test ./...
```

### Install Binary (Optional)
```bash
sudo mv coderun /usr/local/bin/
```

## Project Architecture

This is a Go CLI application built with Cobra that provides a Container-as-a-Service interface for deploying Docker containers to Kubernetes. The application consists of:

### Core Structure
- **`main.go`**: Entry point that calls `cmd.Execute()`
- **`cmd/`**: Contains all CLI commands using Cobra framework
  - `root.go`: Base command definition and configuration
  - `login.go`, `deploy.go`, `list.go`, `status.go`, `delete.go`, `logs.go`: Individual command implementations
- **`internal/client/`**: HTTP API client implementation
  - `client.go`: Core HTTP client with authentication
  - `auth.go`, `deployments.go`: API endpoint implementations
  - `types.go`: Request/response type definitions
- **`internal/utils/`**: Utility functions
  - `config.go`: Configuration file management (~/.coderun/config.json)
  - `env.go`: Environment file parsing and resource validation

### Key Design Patterns
- Uses Cobra for CLI structure with subcommands
- HTTP client with bearer token authentication
- Configuration stored in `~/.coderun/config.json` with access tokens
- Environment variables loaded from files using `KEY=VALUE` format
- Resource validation for CPU (e.g., `100m`, `0.5`) and memory (e.g., `128Mi`, `1Gi`)
- Long timeout (10 minutes) for deployment operations

### API Integration
- Communicates with CodeRun API (default: `http://localhost:8000`)
- Supports custom API URL via `--api-url` flag or `CODERUN_API_URL` environment variable
- Handles deployment creation, listing, status checking, and deletion
- Supports container resource limits, replicas, HTTP port exposure, and environment variables

### Commands Available
- `coderun login`: Authenticate and store access token
- `coderun deploy IMAGE`: Deploy containers with optional flags for replicas, resources, ports, env files
- `coderun list`: List all deployments in table format
- `coderun status APP_NAME`: Get detailed deployment status
- `coderun delete DEPLOYMENT_ID`: Remove deployments
- `coderun logs DEPLOYMENT_ID`: View container logs

The CLI is designed to be the primary interface for the CodeRun platform, abstracting Kubernetes complexity while providing essential container deployment and management capabilities.