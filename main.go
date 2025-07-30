package main

import (
	"os"

	"github.com/helmcode/coderun-cli/cmd"
)

// Build version injected at compile time
var version = "dev"

func main() {
	// Pass version information to cmd package
	cmd.SetVersionInfo(version)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
