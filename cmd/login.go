package cmd

import (
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/helmcode/coderun-cli/internal/client"
	"github.com/helmcode/coderun-cli/internal/utils"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to CodeRun platform",
	Long: `Login to the CodeRun platform using your email and password.
This will store an authentication token for subsequent commands.

Example:
  coderun login`,
	Run: runLogin,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func runLogin(cmd *cobra.Command, args []string) {
	// Load current config
	config, err := utils.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Prompt for email
	fmt.Print("Email: ")
	var email string
	if _, err := fmt.Scanln(&email); err != nil {
		fmt.Printf("Error reading email: %v\n", err)
		os.Exit(1)
	}

	// Prompt for password (hidden input)
	fmt.Print("Password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("\nError reading password: %v\n", err)
		os.Exit(1)
	}
	password := string(passwordBytes)
	fmt.Println() // Add newline after password input

	// Create client and attempt login
	apiClient := client.NewClient(config.BaseURL)

	fmt.Println("Logging in...")
	loginResp, err := apiClient.Login(email, password)
	if err != nil {
		fmt.Printf("Login failed: %v\n", err)
		os.Exit(1)
	}

	// Save token to config
	config.AccessToken = loginResp.AccessToken
	if err := utils.SaveConfig(config); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Successfully logged in!")
	fmt.Printf("Token saved to config file\n")
}
