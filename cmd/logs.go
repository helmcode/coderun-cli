package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/helmcode/coderun-cli/internal/client"
	"github.com/helmcode/coderun-cli/internal/utils"
	"github.com/spf13/cobra"
)

var (
	logLines int
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs [DEPLOYMENT_ID]",
	Short: "Get logs from a deployment",
	Long: `Get logs from all pods in a deployment. Similar to 'docker logs' but for Kubernetes deployments.

Shows the most recent logs from all pods in the deployment.

Examples:
  coderun logs abc12345-6789-def0-1234-567890abcdef
  coderun logs abc12345-6789-def0-1234-567890abcdef --lines=200`,
	Args: cobra.ExactArgs(1),
	Run:  runLogs,
}

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.Flags().IntVarP(&logLines, "lines", "n", 100, "Number of lines to show from the end of the logs")
}

func runLogs(cmd *cobra.Command, args []string) {
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

	// Create client and get logs
	apiClient := client.NewClient(config.BaseURL)
	apiClient.SetToken(config.AccessToken)

	fmt.Printf("Fetching logs for deployment %s...\n", deploymentID)
	logsResponse, err := apiClient.GetDeploymentLogs(deploymentID, logLines)
	if err != nil {
		fmt.Printf("Failed to get logs: %v\n", err)
		os.Exit(1)
	}

	// Display deployment info
	fmt.Printf("📦 Deployment: %s (%s)\n", logsResponse.AppName, logsResponse.DeploymentID)
	fmt.Printf("🐳 Image: %s\n", logsResponse.Image)
	fmt.Printf("📊 Status: %s\n", logsResponse.Status)
	fmt.Printf("🔢 Total Pods: %d\n", logsResponse.TotalPods)
	fmt.Printf("📝 Lines: %d\n\n", logLines)

	if len(logsResponse.Logs) == 0 {
		fmt.Println("ℹ️  No pods found or no logs available")
		return
	}

	// Display logs for each pod
	for podName, podInfo := range logsResponse.Logs {
		fmt.Printf("🚀 Pod: %s (Status: %s, Restarts: %d)\n", podName, podInfo.Status, podInfo.RestartCount)
		fmt.Println("📋 Logs:")
		fmt.Println(strings.Repeat("-", 80))

		if strings.TrimSpace(podInfo.Logs) == "" {
			fmt.Println("(No logs available)")
		} else {
			// Add some basic formatting
			lines := strings.Split(podInfo.Logs, "\n")
			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					fmt.Printf("   %s\n", line)
				}
			}
		}

		fmt.Println(strings.Repeat("-", 80))
		fmt.Println()
	}
}
