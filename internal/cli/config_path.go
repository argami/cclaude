package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var configPathCmd = &cobra.Command{
	Use:   "config-path",
	Short: "Show configuration file location",
	Long:  `Show the path where the configuration file is located.

Displays the configuration file path that cclaude will use.
This is useful for debugging configuration issues.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, err := getConfigPath()
		if err != nil {
			return err
		}

		fmt.Printf("Configuration file location: %s\n", configPath)

		// Check if file exists
		if _, err := os.Stat(configPath); err != nil {
			if os.IsNotExist(err) {
				fmt.Println("Status: File does not exist (run 'cclaude init' to create)")
			} else {
				fmt.Printf("Status: Error accessing file: %v\n", err)
			}
		} else {
			fmt.Println("Status: File exists and will be loaded")
		}

		return nil
	},
}

// NewConfigPathCommand returns the config-path command
func NewConfigPathCommand() *cobra.Command {
	return configPathCmd
}
