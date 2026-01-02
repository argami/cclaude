package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "cclaude.yaml")

	// Test loading from non-existent file (should return defaults)
	t.Run("non-existent file", func(t *testing.T) {
		cfg, err := LoadConfig("/nonexistent/path/cclaude.yaml")
		if err != nil {
			t.Errorf("LoadConfig() unexpected error: %v", err)
		}
		if cfg == nil {
			t.Error("LoadConfig() returned nil config")
		}
		if cfg.GetTimeoutMS() != 3000000 {
			t.Errorf("GetTimeoutMS() = %d, want 3000000", cfg.GetTimeoutMS())
		}
	})

	// Test loading from valid file
	t.Run("valid config file", func(t *testing.T) {
		content := `
provider: minimax
model: MiniMax-M2.1
timeout_ms: 600000
`
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}

		cfg, err := LoadConfig(configPath)
		if err != nil {
			t.Errorf("LoadConfig() unexpected error: %v", err)
		}

		if cfg.Provider != "minimax" {
			t.Errorf("Provider = %s, want minimax", cfg.Provider)
		}
		if cfg.Model != "MiniMax-M2.1" {
			t.Errorf("Model = %s, want MiniMax-M2.1", cfg.Model)
		}
		if cfg.GetTimeoutMS() != 600000 {
			t.Errorf("GetTimeoutMS() = %d, want 600000", cfg.GetTimeoutMS())
		}
	})

	// Test environment variable override (only when config value is empty)
	t.Run("env override for empty config", func(t *testing.T) {
		content := `
provider:
model:
`
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}

		os.Setenv("CCLAUDE_PROVIDER", "mimo")
		os.Setenv("CCLAUDE_MODEL", "kimi-k2")
		defer os.Unsetenv("CCLAUDE_PROVIDER")
		defer os.Unsetenv("CCLAUDE_MODEL")

		cfg, err := LoadConfig(configPath)
		if err != nil {
			t.Errorf("LoadConfig() unexpected error: %v", err)
		}

		if cfg.Provider != "mimo" {
			t.Errorf("Provider = %s, want mimo (from env)", cfg.Provider)
		}
		if cfg.Model != "kimi-k2" {
			t.Errorf("Model = %s, want kimi-k2 (from env)", cfg.Model)
		}
	})
}

func TestGetDefaults(t *testing.T) {
	cfg := getDefaults()

	if cfg.Provider != "" {
		t.Errorf("Provider = %s, want empty", cfg.Provider)
	}
	if cfg.Model != "" {
		t.Errorf("Model = %s, want empty", cfg.Model)
	}
	if cfg.GetTimeoutMS() != 3000000 {
		t.Errorf("GetTimeoutMS() = %d, want 3000000", cfg.GetTimeoutMS())
	}
	if cfg.EnvOverrides == nil {
		t.Error("EnvOverrides should not be nil")
	}
}
