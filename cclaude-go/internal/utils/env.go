package utils

import (
	"fmt"
	"os"

	"github.com/argami/cclaude-go/pkg/types"
)

func SetupEnvironment(provider types.ProviderConfig, authToken string, modelOverride string) error {
	// Limpiar variables anteriores
	os.Unsetenv("ANTHROPIC_BASE_URL")
	os.Unsetenv("MAIN_MODEL")
	os.Unsetenv("ANTHROPIC_AUTH_TOKEN")

	// Configurar nuevas variables
	if err := os.Setenv("ANTHROPIC_BASE_URL", provider.BaseURL); err != nil {
		return fmt.Errorf("error configurando ANTHROPIC_BASE_URL: %w", err)
	}

	model := provider.Model
	if modelOverride != "" {
		model = modelOverride
	}

	os.Setenv("MAIN_MODEL", model)
	os.Setenv("ANTHROPIC_AUTH_TOKEN", authToken)
	os.Setenv("ANTHROPIC_DEFAULT_OPUS_MODEL", provider.OpusModel)
	os.Setenv("ANTHROPIC_MODEL", model)
	os.Setenv("ANTHROPIC_DEFAULT_SONNET_MODEL", model)
	os.Setenv("ANTHROPIC_DEFAULT_HAIKU_MODEL", model)
	os.Setenv("CLAUDE_CODE_SUBAGENT_MODEL", model)
	os.Setenv("DISABLE_NON_ESSENTIAL_MODEL_CALLS", "1")
	os.Setenv("CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC", "1")
	os.Setenv("API_TIMEOUT_MS", "3000000")

	return nil
}