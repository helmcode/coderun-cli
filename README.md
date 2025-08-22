# CodeRun CLI

CodeRun CLI is a command-line tool that allows you to deploy and manage Docker containers on a Helmcode Kubernetes platform easily.

## 🚀 Installation

### Quick Installation (Recommended)

#### Linux (AMD64)
```bash
curl -L https://github.com/helmcode/coderun-cli/releases/latest/download/coderun-linux-amd64 -o coderun
chmod +x coderun
sudo mv coderun /usr/local/bin/
```

#### Linux (ARM64)
```bash
curl -L https://github.com/helmcode/coderun-cli/releases/latest/download/coderun-linux-arm64 -o coderun
chmod +x coderun
sudo mv coderun /usr/local/bin/
```

#### macOS (Intel)
```bash
curl -L https://github.com/helmcode/coderun-cli/releases/latest/download/coderun-darwin-amd64 -o coderun
chmod +x coderun
sudo mv coderun /usr/local/bin/
```

#### macOS (Apple Silicon)
```bash
curl -L https://github.com/helmcode/coderun-cli/releases/latest/download/coderun-darwin-arm64 -o coderun
chmod +x coderun
sudo mv coderun /usr/local/bin/
```

#### Windows
1. Download the appropriate file from [Releases](https://github.com/helmcode/coderun-cli/releases/latest)
2. Rename it to `coderun.exe`
3. Place it in your PATH

### Build from Source

```bash
git clone https://github.com/helmcode/coderun-cli.git
cd coderun-cli
go build -o coderun .
```

## 📋 Installation Verification

```bash
coderun --version
```

## 🔧 Basic Usage

### 1. Authentication
```bash
coderun login
```

### 2. Deploy an Application

#### Web Applications (HTTP)
```bash
# Basic deployment
coderun deploy nginx:latest --name my-web-app --http-port 80

# With custom resources
coderun deploy my-app:v1.0 --name web-app --http-port 8080 --replicas 3 --cpu 500m --memory 1Gi

# With environment variables
coderun deploy my-app:latest --name prod-app --http-port 3000 --env-file .env
```

#### TCP Applications (Databases, etc.)
```bash
# Redis
coderun deploy redis:latest --name my-redis --tcp-port 6379

# PostgreSQL
coderun deploy postgres:latest --name my-db --tcp-port 5432 --env-file database.env

# Custom TCP application
coderun deploy my-tcp-app:latest --name tcp-service --tcp-port 9000
```

### 3. Deployment Management

#### List deployments
```bash
coderun list
```

#### View detailed status
```bash
coderun status <DEPLOYMENT_ID>
```

#### Delete deployment
```bash
coderun delete <DEPLOYMENT_ID>
```

## 📖 Available Commands

| Command | Description |
|---------|-------------|
| `login` | Authenticate with the platform |
| `deploy` | Deploy an application |
| `list` | List all deployments |
| `status` | View detailed deployment status |
| `delete` | Delete a deployment |

## 🔗 Connection Types

### HTTP/HTTPS
- Web applications are automatically exposed with HTTPS
- URL format: `https://app-name-id.helmcode.com`
- Automatic TLS certificates

### TCP
- TCP applications are exposed on the LoadBalancer
- Format: `app-name-id.helmcode.com:port`
- Ideal for databases, TCP APIs, etc.

## ⚙️ Configuration Options

### Deploy Command Flags

| Flag | Description | Example |
|------|-------------|---------|
| `--name` | Application name (required) | `--name my-app` |
| `--replicas` | Number of replicas | `--replicas 3` |
| `--cpu` | CPU limit | `--cpu 500m` |
| `--memory` | Memory limit | `--memory 1Gi` |
| `--http-port` | HTTP port to expose | `--http-port 8080` |
| `--tcp-port` | TCP port to expose | `--tcp-port 5432` |
| `--env-file` | Environment variables file | `--env-file .env` |

### Environment File Format (.env)
```bash
DATABASE_URL=postgres://user:pass@host:5432/db
API_KEY=your-secret-key
DEBUG=true
```

## 🔍 Practical Examples

### Deploy WordPress
```bash
coderun deploy wordpress:latest --name my-blog --http-port 80 --replicas 2
```

### Deploy Node.js API
```bash
coderun deploy my-api:v2.1 --name api-service --http-port 3000 --cpu 200m --memory 512Mi --env-file api.env
```

### Deploy Redis for Caching
```bash
coderun deploy redis:alpine --name cache --tcp-port 6379
```

### Deploy MongoDB
```bash
coderun deploy mongo:latest --name database --tcp-port 27017 --env-file mongo.env
```

## 🔒 Validations

The CLI includes automatic validations for:
- ✅ Application names (3-30 characters, lowercase, letters/numbers/hyphens)
- ✅ Ports in valid range (1-65535)
- ✅ HTTP/TCP mutual exclusion (only one allowed)
- ✅ Resource format (CPU/memory)
- ✅ Authentication verification

## 🚦 Deployment States

| State | Description |
|-------|-------------|
| `pending` | Deployment being created |
| `running` | Application running correctly |
| `failed` | Deployment error |
| `stopped` | Application stopped |

## 🐛 Troubleshooting

### Error: "Please login first"
```bash
coderun login
```

### Error: "App name is required"
```bash
# Add the --name flag
coderun deploy nginx:latest --name my-application
```

### Error: "Cannot specify both --http-port and --tcp-port"
```bash
# Use only one of them
coderun deploy my-app:latest --name app --http-port 8080
# Or
coderun deploy my-app:latest --name app --tcp-port 9000
```

## 📦 Releases

Releases are automatically generated when a tag is created in the repository:

```bash
git tag v0.0.1
git push origin v0.0.1
```

This automatically triggers compilation for all platforms and creates a release on GitHub.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📁 Complete documentation
1. [Intro](https://gist.github.com/sre-helmcode/8b67790107fa890ce078696edb5967e0#file-intro-md)
2. [Basic Usage](https://gist.github.com/sre-helmcode/8b67790107fa890ce078696edb5967e0#file-basic_usage-md)
3. [Basic Deploy](https://gist.github.com/sre-helmcode/8b67790107fa890ce078696edb5967e0#file-basic_deploy-md)
4. [Advanced Deploy](https://gist.github.com/sre-helmcode/8b67790107fa890ce078696edb5967e0#file-advanced_deploy-md)
5. [Debugging](https://gist.github.com/sre-helmcode/8b67790107fa890ce078696edb5967e0#file-debug-md)
6. [Build and Deploy from a Dockerfile](https://gist.github.com/sre-helmcode/8b67790107fa890ce078696edb5967e0#file-build_and_deploy-md)

## 📄 License

[Apache License 2.0](https://github.com/helmcode/coderun-cli/blob/main/LICENSE)

## 🔗 Links

- [Report bugs](https://github.com/helmcode/coderun-cli/issues)
- [Request features](https://github.com/helmcode/coderun-cli/issues)
