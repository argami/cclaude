package provider

import (
	"os"
	"os/exec"

	"github.com/argami/cclaude/internal/config"
	"github.com/argami/cclaude/pkg/models"
)

// Execute runs claude with the specified provider configuration
func Execute(providerName string, args []string, cfg *config.Config) error {
	provider, err := models.GetProvider(providerName)
	if err != nil {
		return err
	}

	// Set environment variables
	os.Setenv("ANTHROPIC_BASE_URL", provider.BaseURL)
	os.Setenv("ANTHROPIC_AUTH_TOKEN", os.Getenv(provider.EnvKey))

	// Set model if specified
	if cfg.GetModel() != "" {
		os.Setenv("ANTHROPIC_DEFAULT_OPUS_MODEL", cfg.GetModel())
	}

	// Apply environment overrides from config
	for key, value := range cfg.EnvOverrides {
		os.Setenv(key, value)
	}

	// Set common environment variables
	os.Setenv("ANTHROPIC_API_KEY", "")
	os.Setenv("ANTHROPIC_MODEL", provider.Model)
	os.Setenv("ANTHROPIC_DEFAULT_SONNET_MODEL", provider.Model)
	os.Setenv("ANTHROPIC_DEFAULT_HAIKU_MODEL", provider.Model)
	os.Setenv("CLAUDE_CODE_SUBAGENT_MODEL", provider.Model)

	// Execute claude binary
	cmd := exec.Command("claude", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// ExecuteDryRun shows what would be executed without actually running
func ExecuteDryRun(providerName string, args []string, cfg *config.Config) error {
	provider, err := models.GetProvider(providerName)
	if err != nil {
		return err
	}

	// Print configuration
	printConfiguration(provider, args, cfg)

	return nil
}

// printConfiguration prints the configuration that would be used
func printConfiguration(provider models.Provider, args []string, cfg *config.Config) {
	// Print provider info
}
