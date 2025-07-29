package cmd

import (
	"fmt"
	"os"

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

	// Display deployments with full IDs
	fmt.Printf("%-36s %-20s %-30s %-8s %-10s %-16s\n",
		"ID", "App Name", "Image", "Replicas", "Status", "Created")
	fmt.Printf("%-36s %-20s %-30s %-8s %-10s %-16s\n",
		"------------------------------------", "--------------------", "------------------------------",
		"--------", "----------", "----------------")

	// Add rows with full IDs
	for _, deployment := range deploymentList.Deployments {
		// Truncate app name if too long
		appName := deployment.AppName
		if len(appName) > 20 {
			appName = appName[:18] + ".."
		}

		// Truncate image if too long
		image := deployment.Image
		if len(image) > 30 {
			image = image[:28] + ".."
		}

		// Format created date
		createdAt := deployment.CreatedAt.Format("2006-01-02 15:04")

		fmt.Printf("%-36s %-20s %-30s %-8d %-10s %-16s\n",
			deployment.ID,
			appName,
			image,
			deployment.Replicas,
			deployment.Status,
			createdAt)
	}
}
