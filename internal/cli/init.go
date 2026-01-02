package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	config "github.com/argami/cclaude-glm/internal/config"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	Long:  `Initialize a new configuration file in the default location.

Creates a default config.yaml file in ~/.config/cclaude/ (or $XDG_CONFIG_HOME/cclaude/).
If a config file already exists, it will not be overwritten unless --force is used.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")

		// Get config path
		configPath, err := getConfigPath()
		if err != nil {
			return err
		}

		// Check if config already exists
		if _, err := os.Stat(configPath); err == nil && !force {
			return fmt.Errorf("config file already exists at %s (use --force to overwrite)", configPath)
		}

		// Create config directory if needed
		configDir := filepath.Dir(configPath)
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		// Create default config
		defaultConfig := generateDefaultConfig()

		// Write to file
		if err := writeConfig(configPath, defaultConfig); err != nil {
			return err
		}

		fmt.Printf("✅ Configuration file created at: %s\n", configPath)
		fmt.Println("\n📝 Next steps:")
		fmt.Println("1. Edit the configuration file to add your API keys")
		fmt.Println("2. Set environment variables for your providers:")
		fmt.Println("   export MIMO_API_KEY=your_key")
		fmt.Println("   export MINIMAX_API_KEY=your_key")
		fmt.Println("   export KIMI_API_KEY=your_key")
		fmt.Println("   export GLM_API_KEY=your_key")
		fmt.Println("\n🔍 Available providers:")
		for name := range defaultConfig.Providers {
			fmt.Printf("   - %s\n", name)
		}

		return nil
	},
}

func init() {
	initCmd.Flags().BoolP("force", "f", false, "Overwrite existing config file")
}

func getConfigPath() (string, error) {
	// Try XDG_CONFIG_HOME first
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "cclaude", "config.yaml"), nil
	}

	// Fallback to ~/.config
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config", "cclaude", "config.yaml"), nil
}

func generateDefaultConfig() *config.Config {
	return &config.Config{
		Providers: map[string]*config.Provider{
			"mimo": {
				Name:      "Mimo",
				BaseURL:   "https://api.xiaomimimo.com/anthropic",
				Model:     "mimo-v2-flash",
				EnvKey:    "MIMO_API_KEY",
				OpusModel: "mimo-v2-flash",
			},
			"minimax": {
				Name:      "MiniMax",
				BaseURL:   "https://api.minimax.io/anthropic",
				Model:     "MiniMax-M2.1",
				EnvKey:    "MINIMAX_API_KEY",
				OpusModel: "MiniMax-M2.1",
			},
			"kimi": {
				Name:      "Kimi",
				BaseURL:   "https://api.kimi.com/coding/",
				Model:     "kimi-k2-0711-preview",
				EnvKey:    "KIMI_API_KEY",
				OpusModel: "kimi-k2-thinking-turbo",
			},
			"glm": {
				Name:      "GLM",
				BaseURL:   "https://api.z.ai/api/anthropic",
				Model:     "glm-4.7",
				EnvKey:    "GLM_API_KEY",
				OpusModel: "glm-4.7",
			},
			"claude": {
				Name:      "Claude",
				BaseURL:   "",
				Model:     "",
				EnvKey:    "",
				OpusModel: "",
			},
		},
		Settings: config.Settings{
			TimeoutMs:            3000000,
			DisableNonEssential: true,
			LogLevel:              "info",
		},
	}
}

func writeConfig(path string, cfg *config.Config) error {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	// Set values
	providersMap := make(map[string]interface{})
	for name, p := range cfg.Providers {
		providersMap[name] = map[string]interface{}{
			"name":        p.Name,
			"base_url":    p.BaseURL,
			"model":       p.Model,
			"env_key":     p.EnvKey,
			"opus_model":  p.OpusModel,
		}
	}

	viper.Set("providers", providersMap)
	viper.Set("settings", map[string]interface{}{
		"timeout_ms":                    cfg.Settings.TimeoutMs,
		"disable_non_essential_calls": cfg.Settings.DisableNonEssential,
		"log_level":                      cfg.Settings.LogLevel,
	})

	return viper.WriteConfig()
}
