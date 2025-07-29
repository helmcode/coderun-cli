package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/helmcode/coderun-cli/internal/client"
	"github.com/helmcode/coderun-cli/internal/utils"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy IMAGE",
	Short: "Deploy a Docker container",
	Long: `Deploy a Docker container to the CodeRun platform.

Examples:
  coderun deploy nginx:latest
  coderun deploy my-app:v1.0 --replicas=3 --cpu=500m --memory=1Gi
  coderun deploy my-app:latest --http-port=8080 --env-file=.env
  coderun deploy my-app:latest --replicas=2 --cpu=200m --memory=512Mi --http-port=3000 --env-file=production.env`,
	Args: cobra.ExactArgs(1),
	Run:  runDeploy,
}

var (
	replicas int
	cpu      string
	memory   string
	httpPort int
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
	deployCmd.Flags().StringVar(&envFile, "env-file", "", "Path to environment file")
	deployCmd.Flags().StringVar(&appName, "name", "", "Application name (optional, auto-generated if not provided)")
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

	// Create client and deploy
	apiClient := client.NewClient(config.BaseURL)
	apiClient.SetToken(config.AccessToken)

	fmt.Printf("Deploying %s...\n", image)
	if httpPort > 0 {
		fmt.Println("ℹ️  Note: Deploy with HTTP port may take several minutes (waiting for TLS certificate)")
	}
	deployment, err := apiClient.CreateDeployment(&deployReq)
	if err != nil {
		fmt.Printf("Deployment failed: %v\n", err)
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
	if len(deployment.EnvironmentVars) > 0 {
		fmt.Printf("Environment Variables: %d\n", len(deployment.EnvironmentVars))
	}

	fmt.Printf("Status: %s\n", deployment.Status)
	fmt.Printf("Created: %s\n", deployment.CreatedAt.Format("2006-01-02 15:04:05"))
}
