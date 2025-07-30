package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/helmcode/coderun-cli/internal/client"
	"github.com/helmcode/coderun-cli/internal/utils"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all deployments",
	Long: `List all deployments in your account.

Example:
  coderun list`,
	Run: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

// calculateColumnWidths calculates the optimal width for each column based on content
func calculateColumnWidths(deployments []client.DeploymentResponse) (int, int, int, int, int, int, int) {
	// Minimum widths for headers
	idWidth := len("ID")
	appNameWidth := len("App Name")
	imageWidth := len("Image")
	replicasWidth := len("Replicas")
	statusWidth := len("Status")
	connectionWidth := len("Connection")
	createdWidth := len("Created")

	// Check content and adjust widths
	for _, deployment := range deployments {
		if len(deployment.ID) > idWidth {
			idWidth = len(deployment.ID)
		}
		if len(deployment.AppName) > appNameWidth {
			appNameWidth = len(deployment.AppName)
		}
		if len(deployment.Image) > imageWidth {
			imageWidth = len(deployment.Image)
		}

		replicasStr := fmt.Sprintf("%d", deployment.Replicas)
		if len(replicasStr) > replicasWidth {
			replicasWidth = len(replicasStr)
		}

		if len(deployment.Status) > statusWidth {
			statusWidth = len(deployment.Status)
		}

		connection := getConnectionString(&deployment)
		if len(connection) > connectionWidth {
			connectionWidth = len(connection)
		}

		createdStr := deployment.CreatedAt.Format("2006-01-02 15:04")
		if len(createdStr) > createdWidth {
			createdWidth = len(createdStr)
		}
	}

	// Add some padding
	return idWidth + 2, appNameWidth + 2, imageWidth + 2, replicasWidth + 2, statusWidth + 2, connectionWidth + 2, createdWidth + 2
}

// getConnectionString generates a user-friendly connection string based on the deployment type
func getConnectionString(deployment *client.DeploymentResponse) string {
	// HTTP deployments - use the URL field from backend if available
	if deployment.URL != nil && *deployment.URL != "" {
		return *deployment.URL
	}

	// Fallback for HTTP deployments if URL is not set yet
	if deployment.HTTPPort != nil {
		return fmt.Sprintf("https://%s.helmcode.com", deployment.AppName)
	}

	// TCP deployments - prefer TCPConnection field if available
	if deployment.TCPConnection != nil && *deployment.TCPConnection != "" {
		return *deployment.TCPConnection
	}

	// TCP with NodePort but no connection string yet
	if deployment.TCPPort != nil && deployment.TCPNodePort != nil {
		return fmt.Sprintf("67.207.79.206:%d", *deployment.TCPNodePort)
	}

	// No exposed ports
	return "Internal only"
}

func runList(cmd *cobra.Command, args []string) {
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

	// Create client and get deployments
	apiClient := client.NewClient(config.BaseURL)
	apiClient.SetToken(config.AccessToken)

	fmt.Println("Fetching deployments...")
	deploymentList, err := apiClient.ListDeployments()
	if err != nil {
		fmt.Printf("Failed to fetch deployments: %v\n", err)
		os.Exit(1)
	}

	if len(deploymentList.Deployments) == 0 {
		fmt.Println("No deployments found.")
		return
	}

	// Calculate optimal column widths
	idWidth, appNameWidth, imageWidth, replicasWidth, statusWidth, connectionWidth, createdWidth := calculateColumnWidths(deploymentList.Deployments)

	// Create format strings for headers and rows
	headerFormat := fmt.Sprintf("%%-%ds %%-%ds %%-%ds %%-%ds %%-%ds %%-%ds %%-%ds\n",
		idWidth, appNameWidth, imageWidth, replicasWidth, statusWidth, connectionWidth, createdWidth)

	separatorFormat := fmt.Sprintf("%%-%ds %%-%ds %%-%ds %%-%ds %%-%ds %%-%ds %%-%ds\n",
		idWidth, appNameWidth, imageWidth, replicasWidth, statusWidth, connectionWidth, createdWidth)

	// Display headers
	fmt.Printf(headerFormat, "ID", "App Name", "Image", "Replicas", "Status", "Connection", "Created")

	// Display separator
	fmt.Printf(separatorFormat,
		strings.Repeat("-", idWidth),
		strings.Repeat("-", appNameWidth),
		strings.Repeat("-", imageWidth),
		strings.Repeat("-", replicasWidth),
		strings.Repeat("-", statusWidth),
		strings.Repeat("-", connectionWidth),
		strings.Repeat("-", createdWidth))

	// Display rows
	for _, deployment := range deploymentList.Deployments {
		// Format created date
		createdAt := deployment.CreatedAt.Format("2006-01-02 15:04")

		// Generate connection string (no truncation)
		connection := getConnectionString(&deployment)

		fmt.Printf(headerFormat,
			deployment.ID,
			deployment.AppName,
			deployment.Image,
			fmt.Sprintf("%d", deployment.Replicas),
			deployment.Status,
			connection,
			createdAt)
	}
}
