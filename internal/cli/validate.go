package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	config "github.com/argami/cclaude-glm/internal/config"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration",
	Long:  `Validate the current configuration.

Checks if:
- Config file exists and is valid YAML
- All providers are properly configured
- Required environment variables are set
- No conflicting settings

Use this to troubleshoot configuration issues before running cclaude.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("error loading config: %w", err)
		}

		// Validate structure
		if err := config.Validate(cfg); err != nil {
			return fmt.Errorf("invalid config structure: %w", err)
		}

		fmt.Println("✅ Configuration is valid")

		// Show config location
		configPath := os.Getenv("XDG_CONFIG_HOME")
		if configPath == "" {
			home, _ := os.UserHomeDir()
			configPath = home + "/.config/cclaude"
		} else {
			configPath += "/cclaude"
		}
		fmt.Printf("Config location: %s/config.yaml\n", configPath)

		// Show providers
		fmt.Println("\nProviders:")
		for name, provider := range cfg.Providers {
			fmt.Printf("  %s\n", name)

			// Check if API key is set
			if provider.EnvKey != "" {
				if os.Getenv(provider.EnvKey) != "" {
					fmt.Printf("    ✅ API key set (%s)\n", provider.EnvKey)
				} else {
					fmt.Printf("    ⚠️  API key not set (%s)\n", provider.EnvKey)
				}
			}
		}

		// Show settings
		fmt.Println("\nSettings:")
		fmt.Printf("  Timeout: %dms\n", cfg.Settings.TimeoutMs)
		fmt.Printf("  Log level: %s\n", cfg.Settings.LogLevel)
		fmt.Printf("  Disable non-essential calls: %v\n", cfg.Settings.DisableNonEssential)

		return nil
	},
}

// NewValidateCommand returns the validate command
func NewValidateCommand() *cobra.Command {
	return validateCmd
}
