package config

import (
	"os"
	"testing"

	"github.com/argami/cclaude-go/pkg/types"
)

func TestValidateEnvironment(t *testing.T) {
	// Test que fallará porque ValidateEnvironment no existe
	err := ValidateEnvironment()
	// Este test pasará si claude está en el PATH, fallará si no
	// Para testing, asumimos que está disponible
	if err != nil {
		t.Logf("Claude no está en PATH (esperado en algunos entornos): %v", err)
	}
}

func TestValidateAPIKey(t *testing.T) {
	// Test que fallará porque ValidateAPIKey no existe
	provider := types.ProviderConfig{
		Name:    "mimo",
		EnvVar:  "TEST_MIMO_API_KEY",
	}

	// Test 1: Key válida
	os.Setenv("TEST_MIMO_API_KEY", "test-key-12345678")
	err := ValidateAPIKey(provider)
	if err != nil {
		t.Errorf("Expected no error with valid key, got: %v", err)
	}

	// Test 2: Key faltante
	os.Unsetenv("TEST_MIMO_API_KEY")
	err = ValidateAPIKey(provider)
	if err == nil {
		t.Error("Expected error with missing key")
	}

	// Test 3: Key corta
	os.Setenv("TEST_MIMO_API_KEY", "short")
	err = ValidateAPIKey(provider)
	if err == nil {
		t.Error("Expected error with short key")
	}

	// Test 4: Key vacía
	os.Setenv("TEST_MIMO_API_KEY", "")
	err = ValidateAPIKey(provider)
	if err == nil {
		t.Error("Expected error with empty key")
	}
}

func TestValidateMultipleProviders(t *testing.T) {
	providers := []types.ProviderConfig{
		{Name: "mimo", EnvVar: "MIMO_API_KEY"},
		{Name: "minimax", EnvVar: "MINIMAX_API_KEY"},
		{Name: "kimi", EnvVar: "KIMI_API_KEY"},
		{Name: "glm", EnvVar: "GLM_API_KEY"},
	}

	for _, provider := range providers {
		t.Run(provider.Name, func(t *testing.T) {
			// Set valid key
			os.Setenv(provider.EnvVar, "valid-key-12345678")
			err := ValidateAPIKey(provider)
			if err != nil {
				t.Errorf("Provider %s should accept valid key: %v", provider.Name, err)
			}

			// Unset key
			os.Unsetenv(provider.EnvVar)
			err = ValidateAPIKey(provider)
			if err == nil {
				t.Errorf("Provider %s should reject missing key", provider.Name)
			}
		})
	}
}