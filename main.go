package main

import (
	"os"

	"github.com/helmcode/coderun-cli/cmd"
)

// Build variables injected at compile time
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Pass version information to cmd package
	cmd.SetVersionInfo(version, commit, date)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
