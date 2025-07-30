package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/helmcode/coderun-cli/internal/client"
	"github.com/helmcode/coderun-cli/internal/utils"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status DEPLOYMENT_ID",
	Short: "Get deployment status",
	Long: `Get the status of a specific deployment by deployment ID.

Examples:
  coderun status abc12345-def6-7890-ghij-klmnopqrstuv
  coderun status abc12345def67890`,
	Args: cobra.ExactArgs(1),
	Run:  runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) {
	deploymentID := args[0]

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

	// Create client and get status
	apiClient := client.NewClient(config.BaseURL)
	apiClient.SetToken(config.AccessToken)

	fmt.Printf("Getting status for deployment '%s'...\n", deploymentID)
	status, err := apiClient.GetDeploymentStatus(deploymentID)
	if err != nil {
		fmt.Printf("Failed to get deployment status: %v\n", err)
		os.Exit(1)
	}

	// Display status information
	fmt.Printf("\n📊 Deployment Status for '%s'\n", status.AppName)
	fmt.Printf("─────────────────────────────────\n")
	fmt.Printf("Status: %s\n", status.Status)
	fmt.Printf("Replicas Ready: %d/%d\n", status.ReplicasReady, status.ReplicasDesired)

	if status.URL != nil {
		fmt.Printf("URL: %s\n", *status.URL)

		// Show TLS certificate information if available
		if status.TLSCertificate != nil {
			if status.TLSCertificate.Ready {
				fmt.Printf("🔐 TLS Certificate: ✅ Ready\n")
			} else {
				fmt.Printf("🔐 TLS Certificate: ⏳ %s\n", status.TLSCertificate.Status)
				if status.TLSCertificate.Message != "" {
					fmt.Printf("    📝 %s\n", status.TLSCertificate.Message)
				}
			}
		}

		// Show note about URL accessibility if it exists
		if status.URLNote != nil {
			fmt.Printf("    %s\n", *status.URLNote)
		}
	}

	if status.TCPConnection != nil {
		fmt.Printf("TCP Connection: %s\n", *status.TCPConnection)
	}

	// Show persistent storage information if configured
	if status.PersistentVolumeSize != "" && status.PersistentVolumeMountPath != "" {
		fmt.Printf("💾 Persistent Storage:\n")
		fmt.Printf("    Size: %s\n", status.PersistentVolumeSize)
		fmt.Printf("    Mount Path: %s\n", status.PersistentVolumeMountPath)
	}

	if len(status.Pods) > 0 {
		fmt.Printf("\n� Pods:\n")
		for i, pod := range status.Pods {
			fmt.Printf("  Pod %d:\n", i+1)
			for key, value := range pod {
				fmt.Printf("    %s: %s\n", key, value)
			}
			fmt.Println()
		}
	}
}
