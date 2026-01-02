package execution

import (
	"fmt"
	"os"
	"os/exec"
	"github.com/argami/cclaude-glm/internal/config"
)

type Executor struct {
	provider string
}

func NewExecutor(provider string) *Executor {
	return &Executor{provider: provider}
}

// RunProvider is a convenience function to create an executor and run it
func RunProvider(provider string, args []string) error {
	executor := NewExecutor(provider)
	return executor.Execute(args)
}

func (e *Executor) Execute(args []string) error {
	// Load config to get provider details
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("error cargando configuración: %w", err)
	}

	provider, exists := cfg.Providers[e.provider]
	if !exists {
		return fmt.Errorf("proveedor no encontrado: %s", e.provider)
	}

	// Validate API key if needed
	if provider.EnvKey != "" {
		if os.Getenv(provider.EnvKey) == "" {
			return fmt.Errorf("API key no definida: %s\nexport %s='your-api-key'", provider.EnvKey, provider.EnvKey)
		}
	}

	// Setup environment variables
	if provider.BaseURL != "" {
		os.Setenv("ANTHROPIC_BASE_URL", provider.BaseURL)
	}
	if provider.Model != "" {
		os.Setenv("MAIN_MODEL", provider.Model)
		os.Setenv("ANTHROPICIC_MODEL", provider.Model)
		os.Setenv("ANTHROPIC_DEFAULT_SONNET_MODEL", provider.Model)
		os.Setenv("ANTHROPIC_DEFAULT_HAIKU_MODEL", provider.Model)
		os.Setenv("CLAUDE_CODE_SUBAGENT_MODEL", provider.Model)
	}
	if provider.OpusModel != "" {
		os.Setenv("CLAUDE_DEFAULT_SONNET_MODEL", provider.OpusModel)
		os.Setenv("CLAUDE_DEFAULT_OPUS_MODEL", provider.OpusModel)
	}
	os.Setenv("ANTHROPIC_API_KEY", "")
	os.Setenv("DISABLE_NON_ESSENTIAL_MODEL_CALLS", "1")
	os.Setenv("CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC", "1")
	os.Setenv("API_TIMEOUT_MS", "3000000")

	// Prepare claude args
	claudeArgs := []string{"claude"}
	claudeArgs = append(claudeArgs, args...)

	// Execute claude
	cmd := exec.Command(claudeArgs[0], claudeArgs[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
