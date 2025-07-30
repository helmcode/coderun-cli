package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/helmcode/coderun-cli/internal/client"
	"github.com/helmcode/coderun-cli/internal/utils"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy IMAGE",
	Short: "Deploy a Docker container",
	Long: `Deploy a Docker container to the CodeRun platform.

Examples:
  coderun deploy nginx:latest --name my-nginx
  coderun deploy my-app:v1.0 --name my-app --replicas 3 --cpu 500m --memory 1Gi
  coderun deploy my-app:latest --name web-app --http-port 8080 --env-file .env
  coderun deploy redis:latest --name my-redis --tcp-port 6379
  coderun deploy postgres:latest --name my-db --tcp-port 5432 --env-file database.env
  coderun deploy my-app:latest --name prod-app --replicas 2 --cpu 200m --memory 512Mi --http-port 3000 --env-file production.env`,
	Args: cobra.ExactArgs(1),
	Run:  runDeploy,
}

var (
	replicas int
	cpu      string
	memory   string
	httpPort int
	tcpPort  int
	envFile  string
	appName  string
)

func init() {
	rootCmd.AddCommand(deployCmd)

	// Add flags
	deployCmd.Flags().IntVar(&replicas, "replicas", 1, "Number of replicas")
	deployCmd.Flags().StringVar(&cpu, "cpu", "", "CPU resource limit (e.g., 100m, 0.5)")
	deployCmd.Flags().StringVar(&memory, "memory", "", "Memory resource limit (e.g., 128Mi, 1Gi)")
	deployCmd.Flags().IntVar(&httpPort, "http-port", 0, "HTTP port to expose")
	deployCmd.Flags().IntVar(&tcpPort, "tcp-port", 0, "TCP port to expose")
	deployCmd.Flags().StringVar(&envFile, "env-file", "", "Path to environment file")
	deployCmd.Flags().StringVar(&appName, "name", "", "Application name (required, 3-30 chars, lowercase letters/numbers/hyphens only)")
}

// parseValidationError tries to parse backend validation errors and return user-friendly messages
func parseValidationError(errorMsg string) string {
	// Convert to lowercase for easier matching
	lowerError := strings.ToLower(errorMsg)

	// App name validation errors
	if strings.Contains(lowerError, "app_name") {
		if strings.Contains(lowerError, "at least 3 characters") {
			return "App name must be at least 3 characters long. Use --name to specify one (e.g., --name my-app)"
		}
		if strings.Contains(lowerError, "at most 30 characters") || strings.Contains(lowerError, "no more than 30") {
			return "App name must be no more than 30 characters long"
		}
		if strings.Contains(lowerError, "lowercase") || strings.Contains(lowerError, "letters") || strings.Contains(lowerError, "hyphens") {
			return "App name must contain only lowercase letters, numbers, and hyphens"
		}
		return "Invalid app name. Use --name to specify one (3-30 chars, lowercase letters/numbers/hyphens only)"
	}

	// Port validation errors
	if strings.Contains(lowerError, "both http_port and tcp_port") || strings.Contains(lowerError, "both ports") {
		return "Cannot specify both --http-port and --tcp-port. Choose one type of port"
	}
	if strings.Contains(lowerError, "http_port") || strings.Contains(lowerError, "tcp_port") || strings.Contains(lowerError, "port") {
		return "Port must be a valid number between 1 and 65535"
	}

	// Resource validation errors
	if strings.Contains(lowerError, "cpu") && (strings.Contains(lowerError, "invalid") || strings.Contains(lowerError, "format")) {
		return "Invalid CPU value. Use format like '100m' or '0.5'"
	}
	if strings.Contains(lowerError, "memory") && (strings.Contains(lowerError, "invalid") || strings.Contains(lowerError, "format")) {
		return "Invalid memory value. Use format like '128Mi' or '1Gi'"
	}

	// Image validation errors
	if strings.Contains(lowerError, "image") && strings.Contains(lowerError, "at least 1") {
		return "Image name cannot be empty"
	}

	// Generic validation error
	if strings.Contains(lowerError, "422") || strings.Contains(lowerError, "validation") {
		return "Validation error: Please check your input parameters"
	}

	// If we can't parse it, return a cleaner version of the original error
	if strings.Contains(errorMsg, "HTTP 422:") {
		return "Validation error: Please check your input parameters and try again"
	}

	return errorMsg
}

func runDeploy(cmd *cobra.Command, args []string) {
	image := args[0]

	// Load config
	config, err := utils.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	if config.AccessToken == "" {
		fmt.Println("Please login first using 'coderun login'")
		os.Exit(1)
	}

	// Validate resource values
	if err := utils.ValidateResourceValue(cpu, "cpu"); err != nil {
		fmt.Printf("Invalid CPU value: %v\n", err)
		os.Exit(1)
	}

	if err := utils.ValidateResourceValue(memory, "memory"); err != nil {
		fmt.Printf("Invalid memory value: %v\n", err)
		os.Exit(1)
	}

	// Validate that only one of HTTP or TCP port is specified
	if httpPort > 0 && tcpPort > 0 {
		fmt.Println("Cannot specify both --http-port and --tcp-port")
		os.Exit(1)
	}

	// Validate app name if provided
	if appName != "" {
		if len(appName) < 3 {
			fmt.Println("App name must be at least 3 characters long")
			os.Exit(1)
		}
		if len(appName) > 30 {
			fmt.Println("App name must be no more than 30 characters long")
			os.Exit(1)
		}
		// Validate format using regex: only lowercase letters, numbers, and hyphens
		matched, _ := regexp.MatchString(`^[a-z0-9-]+$`, appName)
		if !matched {
			fmt.Println("App name must contain only lowercase letters, numbers, and hyphens")
			os.Exit(1)
		}
		// Cannot start or end with hyphen
		if strings.HasPrefix(appName, "-") || strings.HasSuffix(appName, "-") {
			fmt.Println("App name cannot start or end with a hyphen")
			os.Exit(1)
		}
	} else {
		fmt.Println("App name is required. Use --name to specify one (e.g., --name my-app)")
		fmt.Println("App name must be 3-30 characters long and contain only lowercase letters, numbers, and hyphens")
		os.Exit(1)
	}

	// Validate port ranges
	if httpPort > 0 && (httpPort < 1 || httpPort > 65535) {
		fmt.Println("HTTP port must be between 1 and 65535")
		os.Exit(1)
	}
	if tcpPort > 0 && (tcpPort < 1 || tcpPort > 65535) {
		fmt.Println("TCP port must be between 1 and 65535")
		os.Exit(1)
	}

	// Parse environment file if provided
	var envVars map[string]string
	if envFile != "" {
		envVars, err = utils.ParseEnvFile(envFile)
		if err != nil {
			fmt.Printf("Error parsing env file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Loaded %d environment variables from %s\n", len(envVars), envFile)
	}

	// Create deployment request
	deployReq := client.DeploymentCreate{
		AppName:         appName,
		Image:           image,
		Replicas:        replicas,
		CPULimit:        cpu,
		MemoryLimit:     memory,
		EnvironmentVars: envVars,
	}

	// Add HTTP port if specified
	if httpPort > 0 {
		deployReq.HTTPPort = &httpPort
	}

	// Add TCP port if specified
	if tcpPort > 0 {
		deployReq.TCPPort = &tcpPort
	}

	// Create client and deploy
	apiClient := client.NewClient(config.BaseURL)
	apiClient.SetToken(config.AccessToken)

	fmt.Printf("Deploying %s...\n", image)
	if httpPort > 0 {
		fmt.Println("ℹ️  Note: Deploy with HTTP port may take several minutes (waiting for TLS certificate)")
	}
	if tcpPort > 0 {
		fmt.Println("ℹ️  Note: Deploy with TCP port will be available in the NodePort range (30000-32767)")
	}
	deployment, err := apiClient.CreateDeployment(&deployReq)
	if err != nil {
		userFriendlyError := parseValidationError(err.Error())
		fmt.Printf("Deployment failed: %s\n", userFriendlyError)
		os.Exit(1)
	}

	fmt.Println("✅ Deployment created successfully!")
	fmt.Printf("Deployment ID: %s\n", deployment.ID)
	fmt.Printf("App Name: %s\n", deployment.AppName)
	fmt.Printf("Image: %s\n", deployment.Image)
	fmt.Printf("Replicas: %d\n", deployment.Replicas)

	if deployment.CPULimit != "" {
		fmt.Printf("CPU: %s\n", deployment.CPULimit)
	}
	if deployment.MemoryLimit != "" {
		fmt.Printf("Memory: %s\n", deployment.MemoryLimit)
	}
	if deployment.HTTPPort != nil {
		fmt.Printf("HTTP Port: %d\n", *deployment.HTTPPort)
	}
	if deployment.TCPPort != nil {
		fmt.Printf("TCP Port: %d\n", *deployment.TCPPort)
	}
	if deployment.TCPNodePort != nil {
		fmt.Printf("TCP NodePort: %d\n", *deployment.TCPNodePort)
	}
	if deployment.TCPConnection != nil {
		fmt.Printf("TCP Connection: %s\n", *deployment.TCPConnection)
	}
	if len(deployment.EnvironmentVars) > 0 {
		fmt.Printf("Environment Variables: %d\n", len(deployment.EnvironmentVars))
	}

	fmt.Printf("Status: %s\n", deployment.Status)
	fmt.Printf("Created: %s\n", deployment.CreatedAt.Format("2006-01-02 15:04:05"))
}
