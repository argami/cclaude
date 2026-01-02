package provider

import (
	"testing"

	"github.com/argami/cclaude-go/pkg/types"
)

func TestGetProvider(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{"mimo", "mimo", false},
		{"minimax", "minimax", false},
		{"kimi", "kimi", false},
		{"glm", "glm", false},
		{"claude", "claude", false},
		{"invalid", "invalid", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := GetProvider(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for %s: %v", tt.input, err)
				}
				if provider == nil {
					t.Errorf("Expected provider for %s, got nil", tt.input)
				}
			}
		})
	}
}

func TestProviderDetails(t *testing.T) {
	// Test que fallará porque Providers map no existe
	provider, err := GetProvider("mimo")
	if err != nil {
		t.Fatalf("Failed to get mimo provider: %v", err)
	}

	expected := types.ProviderConfig{
		Name:        "mimo",
		BaseURL:     "https://api.xiaomimimo.com/anthropic",
		Model:       "mimo-v2-flash",
		OpusModel:   "mimo-v2-flash",
		EnvVar:      "MIMO_API_KEY",
		Description: "Xiaomi MiMo API",
	}

	if provider.Name != expected.Name {
		t.Errorf("Expected Name %s, got %s", expected.Name, provider.Name)
	}
	if provider.BaseURL != expected.BaseURL {
		t.Errorf("Expected BaseURL %s, got %s", expected.BaseURL, provider.BaseURL)
	}
	if provider.Model != expected.Model {
		t.Errorf("Expected Model %s, got %s", expected.Model, provider.Model)
	}
	if provider.EnvVar != expected.EnvVar {
		t.Errorf("Expected EnvVar %s, got %s", expected.EnvVar, provider.EnvVar)
	}
}

func TestAllProviders(t *testing.T) {
	// Verificar que todos los proveedores requeridos existen
	requiredProviders := []string{"mimo", "minimax", "kimi", "glm", "claude"}

	for _, name := range requiredProviders {
		t.Run(name, func(t *testing.T) {
			provider, err := GetProvider(name)
			if err != nil {
				t.Errorf("Provider %s should exist, got error: %v", name, err)
			}
			if provider == nil {
				t.Errorf("Provider %s should not be nil", name)
			}
		})
	}
}