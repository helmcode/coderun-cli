package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/helmcode/coderun-cli/internal/client"
	"github.com/helmcode/coderun-cli/internal/utils"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete DEPLOYMENT_ID_OR_APP_NAME",
	Short: "Delete a deployment",
	Long: `Delete a deployment by its ID or app name.

Examples:
  coderun delete abc12345-def6-7890-ghij-klmnopqrstuv  # Delete by full ID
  coderun delete abc12345                              # Delete by partial ID
  coderun delete my-app --by-name                      # Delete by app name`,
	Args: cobra.ExactArgs(1),
	Run:  runDelete,
}

var byAppName bool

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&byAppName, "by-name", false, "Delete by app name instead of ID")
}

func runDelete(cmd *cobra.Command, args []string) {
	identifier := args[0]

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

	// Create client
	apiClient := client.NewClient(config.BaseURL)
	apiClient.SetToken(config.AccessToken)

	var deploymentID string

	if byAppName {
		// Look up deployment ID by app name
		fmt.Printf("Looking up deployment for app '%s'...\n", identifier)
		deploymentList, err := apiClient.ListDeployments()
		if err != nil {
			fmt.Printf("Failed to fetch deployments: %v\n", err)
			os.Exit(1)
		}

		found := false
		for _, deployment := range deploymentList.Deployments {
			if deployment.AppName == identifier {
				deploymentID = deployment.ID
				found = true
				break
			}
		}

		if !found {
			fmt.Printf("No deployment found with app name: %s\n", identifier)
			os.Exit(1)
		}
	} else {
		// Use identifier as deployment ID (could be full or partial)
		deploymentID = identifier
	}

	fmt.Printf("Deleting deployment %s...\n", deploymentID)
	err = apiClient.DeleteDeployment(deploymentID)
	if err != nil {
		fmt.Printf("Failed to delete deployment: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Deployment %s deleted successfully!\n", deploymentID)
}
