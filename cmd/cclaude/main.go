package main

import (
	"os"
	"github.com/argami/cclaude-glm/internal/cli"
)

func main() {
	rootCmd := cli.NewRootCommand()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
