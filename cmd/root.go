package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version variables set from main
var (
	buildVersion = "dev"
	buildCommit  = "none"
	buildDate    = "unknown"
)

// SetVersionInfo configures version information from main
func SetVersionInfo(version, commit, date string) {
	buildVersion = version
	buildCommit = commit
	buildDate = date

	// Update version in rootCmd
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "coderun",
	Short: "CodeRun Container-as-a-Service CLI",
	Long: `CodeRun CLI allows you to deploy and manage Docker containers 
in a Kubernetes cluster with ease. Deploy your applications with simple commands.

Examples:
  coderun login                                    # Authenticate with your account
  coderun deploy nginx:latest --replicas 2        # Deploy nginx with 2 replicas
  coderun list                                     # List all your deployments
  coderun status <DEPLOYMENT_ID>                  # Check status by deployment ID`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Disable completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
