package provider

import (
	"fmt"
	"os"
)

// Provider defines the interface for AI providers
type Provider interface {
	// Name returns the provider name
	Name() string

	// Validate checks if the provider is properly configured
	// Returns error if API key is missing or invalid
	Validate() error

	// SetupEnv configures environment variables for the provider
	// Sets ANTHROPIC_BASE_URL, MAIN_MODEL, etc.
	SetupEnv() error

	// GetClaudeArgs returns the command and arguments to execute
	// Returns base command (e.g., "claude") and any additional arguments
	GetClaudeArgs() (string, []string)
}

// BaseProvider provides common functionality for all providers
type BaseProvider struct {
	name      string
	baseURL   string
	model     string
	envKey    string
	opusModel string
}

// Name returns the provider name
func (p *BaseProvider) Name() string {
	return p.name
}

// Validate checks if API key is configured
func (p *BaseProvider) Validate() error {
	if p.envKey == "" {
		return nil // No API key needed for native Claude
	}

	if os.Getenv(p.envKey) == "" {
		return fmt.Errorf("API key no definida: %s\nexport %s='your-api-key'", p.envKey, p.envKey)
	}

	return nil
}

// SetupEnv configures environment variables for Claude
func (p *BaseProvider) SetupEnv() error {
	// Set base URL if configured
	if p.baseURL != "" {
		os.Setenv("ANTHROPIC_BASE_URL", p.baseURL)
	}

	// Set model configuration
	if p.model != "" {
		os.Setenv("MAIN_MODEL", p.model)
		os.Setenv("ANTHROPICIC_MODEL", p.model)
		os.Setenv("ANTHROPIC_DEFAULT_SONNET_MODEL", p.model)
		os.Setenv("ANTHROPIC_DEFAULT_HAIKU_MODEL", p.model)
		os.Setenv("CLAUDE_CODE_SUBAGENT_MODEL", p.model)
	}

	// Set Opus model if configured
	if p.opusModel != "" {
		os.Setenv("CLAUDE_DEFAULT_SONNET_MODEL", p.opusModel)
		os.Setenv("CLAUDE_DEFAULT_OPUS_MODEL", p.opusModel)
	}

	// Clear native API key to force use of provider
	os.Setenv("ANTHROPIC_API_KEY", "")

	// Disable non-essential calls for compatibility
	os.Setenv("DISABLE_NON_ESSENTIAL_MODEL_CALLS", "1")
	os.Setenv("CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC", "1")

	// Set timeout
	os.Setenv("API_TIMEOUT_MS", "3000000")

	return nil
}

// GetClaudeArgs returns the command to execute
func (p *BaseProvider) GetClaudeArgs() (string, []string) {
	return "claude", []string{}
}
