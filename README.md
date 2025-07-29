# CodeRun CLI

CodeRun CLI is a command-line tool for deploying and managing Docker containers on the CodeRun platform.

## Installation

### Build from source code

```bash
git clone <repository-url>
cd cli
go build -o coderun .
```

### Move binary to your PATH (optional)

```bash
sudo mv coderun /usr/local/bin/
```

## Configuration

### 1. Login

First you need to authenticate with your CodeRun account:

```bash
coderun login
```

You will be prompted for your email and password. The access token will be saved in `~/.coderun/config.json`.

### 2. Configure API URL (optional)

By default, the CLI uses `https://api.coderun.dev`. To use a local instance:

```bash
coderun --api-url http://localhost:8000 login
```

Or you can configure the environment variable:

```bash
export CODERUN_API_URL="http://localhost:8000"
```

## Available Commands

### `coderun login`

Authenticate with the CodeRun platform.

```bash
coderun login
```

### `coderun deploy`

Deploy a Docker container.

#### Basic options:

```bash
# Simple deployment
coderun deploy nginx:latest

# With specific replicas
coderun deploy nginx:latest --replicas=3

# With resource limits
coderun deploy myapp:v1.0 --cpu=500m --memory=1Gi

# With exposed HTTP port
coderun deploy myapp:v1.0 --http-port=8080
```

#### Using environment files:

```bash
# Load environment variables from file
coderun deploy myapp:v1.0 --env-file=production.env

# Complete example
coderun deploy myapp:v1.0 \
  --replicas=2 \
  --cpu=200m \
  --memory=512Mi \
  --http-port=3000 \
  --env-file=.env \
  --name=my-production-app
```

#### Parameters:

- `--replicas`: Number of replicas (default: 1)
- `--cpu`: CPU limit (example: `100m`, `0.5`, `1`)
- `--memory`: Memory limit (example: `128Mi`, `1Gi`, `512Ki`) 
- `--http-port`: HTTP port to expose
- `--env-file`: Path to environment variables file
- `--name`: Application name (optional, auto-generated)

### `coderun list`

List all deployments.

```bash
coderun list
```

Example output:
```
Found 3 deployment(s):

ID         App Name             Image                          Replicas Status     Created
----------  -------------------- ------------------------------ -------- ---------- ----------------
abc12345..  my-nginx             nginx:latest                   2        Running    2024-01-15 14:30
def67890..  my-api               mycompany/api:v2.1             1        Running    2024-01-15 13:45
ghi24680..  my-worker            mycompany/worker:latest        3        Pending    2024-01-15 14:25
```

### `coderun status`

Get detailed status of a deployment by application ID.

```bash
# Get status by ID
coderun status abc12345def67890
```

### `coderun delete`

Delete a deployment by ID.

```bash
# Delete by specific ID
coderun delete abc12345def67890

# Using command combination
coderun delete $(coderun list | grep my-old-app | awk '{print $1}')
```

## Environment file format

The environment file must follow the `KEY=VALUE` format:

```bash
# example.env
DATABASE_URL=postgresql://user:pass@localhost:5432/db
API_KEY="secret-key-with-special-chars"
DEBUG=true
PORT=8080
```

### Rules:

- One variable per line
- Format: `KEY=VALUE`
- Comments start with `#`
- Empty lines are ignored
- Values can be quoted (`"` or `'`)
- No spaces around the `=`

## Usage examples

### Deploy a simple web application

```bash
# 1. Login
coderun login

# 2. Deploy nginx
coderun deploy nginx:latest --replicas=2 --http-port=80

# 3. Check status
coderun list
coderun status nginx
```

### Deploy an application with database

```bash
# 1. Create environment file
cat > production.env << EOF
DATABASE_URL=postgresql://user:pass@db:5432/myapp
REDIS_URL=redis://redis:6379
API_KEY=your-secret-key
APP_ENV=production
EOF

# 2. Deploy the application
coderun deploy mycompany/webapp:v2.0 \
  --replicas=3 \
  --cpu=500m \
  --memory=1Gi \
  --http-port=8080 \
  --env-file=production.env \
  --name=webapp-prod

# 3. Verify deployment
coderun status abc12345def67890
```

### Deployment management

```bash
# List all deployments
coderun list

# View details of specific deployment
coderun status abc12345def67890

# Delete old deployment
coderun delete abc12345def67890
```

## Troubleshooting

### Authentication error

```bash
Error: Authentication failed
```

**Solution**: Run `coderun login` to authenticate again.

### Connection error

```bash
Error: Failed to connect to API
```

**Solution**: Verify API URL with `--api-url` or the `CODERUN_API_URL` environment variable.

### Resource format error

```bash
Error: Invalid CPU format '1core' (examples: 100m, 0.5, 1)
```

**Solution**: Use valid formats:
- CPU: `100m`, `0.5`, `1`, `2`
- Memory: `128Mi`, `1Gi`, `512Ki`

### Environment file error

```bash
Error: Invalid format at line 5: INVALID LINE (expected KEY=VALUE)
```

**Solution**: Verify each line follows the `KEY=VALUE` format.

## Configuration files

- **Config**: `~/.coderun/config.json`
- **Token**: Stored in configuration file

## Development

### Build

```bash
go build -o coderun .
```

### Run tests

```bash
go test ./...
```

### Add new commands

1. Create file in `cmd/`
2. Implement command with Cobra
3. Register in `cmd/root.go`
