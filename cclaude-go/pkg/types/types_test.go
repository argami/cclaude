package types

import (
	"testing"
	"time"
)

func TestProviderConfig(t *testing.T) {
	// Test que fallará inicialmente porque ProviderConfig no existe
	p := ProviderConfig{
		Name:        "mimo",
		BaseURL:     "https://api.xiaomimimo.com/anthropic",
		Model:       "mimo-v2-flash",
		OpusModel:   "mimo-v2-flash",
		EnvVar:      "MIMO_API_KEY",
		Description: "Xiaomi MiMo API",
	}

	if p.Name != "mimo" {
		t.Errorf("Expected Name 'mimo', got '%s'", p.Name)
	}
}

func TestAppConfig(t *testing.T) {
	// Test que fallará inicialmente porque AppConfig no existe
	provider := &ProviderConfig{
		Name:    "test",
		BaseURL: "https://test.com",
		Model:   "test-model",
	}

	config := AppConfig{
		Provider:      provider,
		Timeout:       5 * time.Minute,
		Debug:         false,
		ModelOverride: "",
		Args:          []string{"arg1", "arg2"},
	}

	if config.Timeout != 5*time.Minute {
		t.Errorf("Expected timeout 5m, got %v", config.Timeout)
	}

	if len(config.Args) != 2 {
		t.Errorf("Expected 2 args, got %d", len(config.Args))
	}
}

func TestStructFields(t *testing.T) {
	// Verificar que todos los campos existen con tipos correctos
	p := ProviderConfig{}
	_ = p.Name
	_ = p.BaseURL
	_ = p.Model
	_ = p.OpusModel
	_ = p.EnvVar
	_ = p.Description

	a := AppConfig{}
	_ = a.Provider
	_ = a.Timeout
	_ = a.Debug
	_ = a.ModelOverride
	_ = a.Args
}