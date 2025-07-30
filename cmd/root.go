package cmd

import (
	"github.com/spf13/cobra"
)

// Version variable set from main
var buildVersion = "dev"

// SetVersionInfo configures version information from main
func SetVersionInfo(version string) {
	buildVersion = version
	rootCmd.Version = version
} // rootCmd represents the base command when called without any subcommands
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
