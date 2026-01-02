package utils

import (
	"os"
	"testing"

	"github.com/argami/cclaude-go/pkg/types"
)

func TestSetupEnvironment(t *testing.T) {
	// Test que fallará porque SetupEnvironment no existe
	provider := types.ProviderConfig{
		Name:       "mimo",
		BaseURL:    "https://api.xiaomimimo.com/anthropic",
		Model:      "mimo-v2-flash",
		OpusModel:  "mimo-v2-flash",
		EnvVar:     "TEST_MIMO_API_KEY",
	}

	authToken := "test-auth-token-12345678"
	modelOverride := ""

	// Limpiar variables anteriores
	os.Unsetenv("ANTHROPIC_BASE_URL")
	os.Unsetenv("MAIN_MODEL")
	os.Unsetenv("ANTHROPIC_AUTH_TOKEN")

	err := SetupEnvironment(provider, authToken, modelOverride)
	if err != nil {
		t.Fatalf("SetupEnvironment() error = %v", err)
	}

	// Verificar que todas las variables se establecieron correctamente
	tests := []struct {
		name     string
		envVar   string
		expected string
	}{
		{"Base URL", "ANTHROPIC_BASE_URL", provider.BaseURL},
		{"Main Model", "MAIN_MODEL", provider.Model},
		{"Auth Token", "ANTHROPIC_AUTH_TOKEN", authToken},
		{"Opus Model", "ANTHROPIC_DEFAULT_OPUS_MODEL", provider.OpusModel},
		{"Model", "ANTHROPIC_MODEL", provider.Model},
		{"Sonnet Model", "ANTHROPIC_DEFAULT_SONNET_MODEL", provider.Model},
		{"Haiku Model", "ANTHROPIC_DEFAULT_HAIKU_MODEL", provider.Model},
		{"Subagent Model", "CLAUDE_CODE_SUBAGENT_MODEL", provider.Model},
		{"Disable Essential", "DISABLE_NON_ESSENTIAL_MODEL_CALLS", "1"},
		{"Disable Traffic", "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC", "1"},
		{"Timeout", "API_TIMEOUT_MS", "3000000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := os.Getenv(tt.envVar)
			if value != tt.expected {
				t.Errorf("Expected %s=%s, got %s=%s", tt.envVar, tt.expected, tt.envVar, value)
			}
		})
	}
}

func TestSetupEnvironmentWithOverride(t *testing.T) {
	provider := types.ProviderConfig{
		Name:       "kimi",
		BaseURL:    "https://api.kimi.com/coding/",
		Model:      "kimi-k2-0711-preview",
		OpusModel:  "kimi-k2-thinking-turbo",
		EnvVar:     "TEST_KIMI_API_KEY",
	}

	authToken := "test-kimi-token"
	modelOverride := "kimi-k2-thinking-turbo"

	err := SetupEnvironment(provider, authToken, modelOverride)
	if err != nil {
		t.Fatalf("SetupEnvironment() error = %v", err)
	}

	// Verificar que el override funciona
	expectedModel := modelOverride
	actualModel := os.Getenv("MAIN_MODEL")
	if actualModel != expectedModel {
		t.Errorf("Expected MAIN_MODEL=%s, got %s", expectedModel, actualModel)
	}
}

func TestEnvironmentCleanup(t *testing.T) {
	// Establecer variables previas
	os.Setenv("ANTHROPIC_BASE_URL", "old-url")
	os.Setenv("ANTHROPIC_AUTH_TOKEN", "old-token")

	provider := types.ProviderConfig{
		Name:       "glm",
		BaseURL:    "https://api.z.ai/api/anthropic",
		Model:      "glm-4.7",
		OpusModel:  "glm-4.7",
		EnvVar:     "TEST_GLM_API_KEY",
	}

	err := SetupEnvironment(provider, "new-token", "")
	if err != nil {
		t.Fatalf("SetupEnvironment() error = %v", err)
	}

	// Verificar que las variables antiguas fueron limpiadas
	oldURL := os.Getenv("ANTHROPIC_BASE_URL")
	if oldURL == "old-url" {
		t.Error("ANTHROPIC_BASE_URL should have been cleaned")
	}

	oldToken := os.Getenv("ANTHROPIC_AUTH_TOKEN")
	if oldToken == "old-token" {
		t.Error("ANTHROPIC_AUTH_TOKEN should have been cleaned")
	}
}

func TestSetupEnvironmentErrors(t *testing.T) {
	// Test que verifica que la función maneja errores correctamente
	provider := types.ProviderConfig{
		Name:    "test",
		BaseURL: "https://test.com",
		Model:   "test-model",
	}

	// Esto debería funcionar sin errores
	err := SetupEnvironment(provider, "test-token", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}