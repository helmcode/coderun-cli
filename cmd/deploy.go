package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/helmcode/coderun-cli/internal/client"
	"github.com/helmcode/coderun-cli/internal/utils"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy [IMAGE]",
	Short: "Deploy a Docker container or build from source",
	Long: `Deploy a Docker container to the CodeRun platform or build from source.

Deploy from existing image:
  coderun deploy nginx:latest --name my-nginx
  coderun deploy redis:latest --name my-redis --tcp-port 6379
  coderun deploy my-app:latest --name prod-app --replicas 2 --cpu 200m --memory 512Mi --http-port 3000 --env-file production.env

Build from source:
  coderun deploy --build . --name my-app
  coderun deploy --build ./my-app --name my-app --dockerfile Dockerfile.prod
  coderun deploy --build . --name web-app --http-port 8080 --env-file .env

With persistent storage (automatically forces replicas to 1):
  coderun deploy postgres:15 --name my-postgres --tcp-port 5432 --storage-size 5Gi --storage-path /var/lib/postgresql/data
  coderun deploy mysql:8 --name my-mysql --tcp-port 3306 --storage-size 10Gi --storage-path /var/lib/mysql
  coderun deploy nginx:latest --name web-server --http-port 80 --storage-size 1Gi --storage-path /usr/share/nginx/html`,
	Args: func(cmd *cobra.Command, args []string) error {
		// If --build is specified, IMAGE argument is optional
		if buildContext != "" {
			return cobra.MaximumNArgs(0)(cmd, args)
		}
		// Otherwise, IMAGE argument is required
		return cobra.ExactArgs(1)(cmd, args)
	},
	Run: runDeploy,
}

var (
	replicas                  int
	cpu                       string
	memory                    string
	httpPort                  int
	tcpPort                   int
	envFile                   string
	appName                   string
	persistentVolumeSize      string
	persistentVolumeMountPath string
	// Build flags
	buildContext   string
	dockerfilePath string
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

	// Persistent storage flags
	deployCmd.Flags().StringVar(&persistentVolumeSize, "storage-size", "", "Size of persistent volume (e.g., '1Gi', '500Mi', '10Gi')")
	deployCmd.Flags().StringVar(&persistentVolumeMountPath, "storage-path", "", "Path where to mount the volume (e.g., '/data', '/var/lib/mysql')")

	// Build flags
	deployCmd.Flags().StringVar(&buildContext, "build", "", "Build from source. Specify the build context directory (e.g., './my-app' or '.')")
	deployCmd.Flags().StringVar(&dockerfilePath, "dockerfile", "Dockerfile", "Path to Dockerfile relative to build context (default: 'Dockerfile')")
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

	// Persistent storage validation errors
	if strings.Contains(lowerError, "persistent_volume_size") || strings.Contains(lowerError, "storage") {
		if strings.Contains(lowerError, "together") || strings.Contains(lowerError, "provided together") {
			return "When using persistent storage, both --storage-size and --storage-path are required"
		}
		if strings.Contains(lowerError, "format") || strings.Contains(lowerError, "10gi") || strings.Contains(lowerError, "500mi") {
			return "Storage size must be in format like '1Gi', '500Mi', '10Gi'"
		}
		return "Invalid storage configuration. Use --storage-size and --storage-path together"
	}
	if strings.Contains(lowerError, "persistent_volume_mount_path") || strings.Contains(lowerError, "mount") {
		if strings.Contains(lowerError, "absolute") || strings.Contains(lowerError, "starting with") {
			return "Storage path must be an absolute path starting with '/' (e.g., '/data', '/var/lib/mysql')"
		}
		return "Invalid storage path. Must be absolute path like '/data' or '/var/lib/mysql'"
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
	var image string

	// Determine if we're building from source or deploying an existing image
	isBuild := buildContext != ""

	if !isBuild {
		if len(args) == 0 {
			fmt.Println("Either specify an IMAGE to deploy or use --build to build from source")
			os.Exit(1)
		}
		image = args[0]
	}

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

	// Validate persistent storage flags
	if persistentVolumeSize != "" || persistentVolumeMountPath != "" {
		// Both flags must be provided together
		if persistentVolumeSize == "" {
			fmt.Println("When using persistent storage, both --storage-size and --storage-path are required")
			os.Exit(1)
		}
		if persistentVolumeMountPath == "" {
			fmt.Println("When using persistent storage, both --storage-size and --storage-path are required")
			os.Exit(1)
		}

		// Validate storage size format
		matched, _ := regexp.MatchString(`^\d+[MGT]i$`, persistentVolumeSize)
		if !matched {
			fmt.Println("Storage size must be in format like '1Gi', '500Mi', '10Gi'")
			os.Exit(1)
		}

		// Validate mount path format (must be absolute path)
		if !strings.HasPrefix(persistentVolumeMountPath, "/") {
			fmt.Println("Storage path must be an absolute path starting with '/' (e.g., '/data', '/var/lib/mysql')")
			os.Exit(1)
		}

		// Force replicas to 1 when using persistent storage
		if replicas > 1 {
			fmt.Printf("Warning: Persistent storage requested, forcing replicas to 1 (was %d)\n", replicas)
			replicas = 1
		}
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

	// Create client
	apiClient := client.NewClient(config.BaseURL)
	apiClient.SetToken(config.AccessToken)

	// Handle build from source
	if isBuild {
		fmt.Printf("Building from source in %s...\n", buildContext)

		// Validate build context
		if _, err := os.Stat(buildContext); os.IsNotExist(err) {
			fmt.Printf("Build context directory does not exist: %s\n", buildContext)
			os.Exit(1)
		}

		// Validate Dockerfile
		if err := utils.ValidateDockerfile(buildContext, dockerfilePath); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Create build context archive
		contextArchivePath := utils.GenerateBuildContextPath(appName)
		defer os.Remove(contextArchivePath) // Clean up

		fmt.Printf("Creating build context archive...\n")
		if err := utils.CreateBuildContext(buildContext, contextArchivePath); err != nil {
			fmt.Printf("Error creating build context: %v\n", err)
			os.Exit(1)
		}

		// Upload and start build
		fmt.Printf("Uploading build context and starting build...\n")
		buildResp, err := apiClient.CreateBuild(contextArchivePath, appName, dockerfilePath)
		if err != nil {
			userFriendlyError := parseValidationError(err.Error())
			fmt.Printf("Build failed: %s\n", userFriendlyError)
			os.Exit(1)
		}

		fmt.Printf("✅ Build started successfully!\n")
		fmt.Printf("Build ID: %s\n", buildResp.ID)
		fmt.Printf("Status: %s\n", buildResp.Status)
		fmt.Printf("Image URI: %s\n", buildResp.ImageURI)

		// Wait for build to complete
		fmt.Printf("Waiting for build to complete...\n")
		for {
			time.Sleep(5 * time.Second)

			status, err := apiClient.GetBuildStatus(buildResp.ID)
			if err != nil {
				fmt.Printf("Error checking build status: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Build status: %s\n", status.Status)

			if status.Status == "completed" {
				fmt.Printf("✅ Build completed successfully!\n")

				// Show build logs for successful builds too
				fmt.Println("\n📋 Build logs:")
				fmt.Println("================")
				logs, err := apiClient.GetBuildLogs(status.ID)
				if err != nil {
					fmt.Printf("❌ Could not retrieve build logs: %v\n", err)
				} else if logs == "" {
					fmt.Println("No logs available")
				} else {
					fmt.Println(logs)
				}
				fmt.Println("================\n")

				image = status.ImageURI
				break
			} else if status.Status == "failed" {
				fmt.Printf("❌ Build failed!\n")

				// Try to get build logs to show the error
				fmt.Println("\n📋 Build logs:")
				fmt.Println("================")
				logs, err := apiClient.GetBuildLogs(status.ID)
				if err != nil {
					fmt.Printf("❌ Could not retrieve build logs: %v\n", err)
				} else if logs == "" {
					fmt.Println("No logs available")
				} else {
					fmt.Println(logs)
				}
				fmt.Println("================")

				os.Exit(1)
			}
		}
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

	// Add persistent storage if specified
	if persistentVolumeSize != "" && persistentVolumeMountPath != "" {
		deployReq.PersistentVolumeSize = persistentVolumeSize
		deployReq.PersistentVolumeMountPath = persistentVolumeMountPath
	}

	// Add HTTP port if specified
	if httpPort > 0 {
		deployReq.HTTPPort = &httpPort
	}

	// Add TCP port if specified
	if tcpPort > 0 {
		deployReq.TCPPort = &tcpPort
	}

	// Deploy the application
	if isBuild {
		fmt.Printf("Deploying built image %s...\n", image)
	} else {
		fmt.Printf("Deploying %s...\n", image)
	}

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
	if deployment.URL != nil {
		fmt.Printf("HTTP URL: %s\n", *deployment.URL)
	}
	if len(deployment.EnvironmentVars) > 0 {
		fmt.Printf("Environment Variables: %d\n", len(deployment.EnvironmentVars))
	}

	fmt.Printf("Status: %s\n", deployment.Status)
	fmt.Printf("Created: %s\n", deployment.CreatedAt.Format("2006-01-02 15:04:05"))

	if isBuild {
		fmt.Println("\n🚀 Successfully built and deployed from source!")
	}
}
