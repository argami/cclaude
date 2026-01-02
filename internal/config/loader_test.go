package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Run("should return default config when no file exists", func(t *testing.T) {
		// Save current directory
	origWd, _ := os.Getwd()
		defer os.Chdir(origWd)

		// Change to temp directory
		tmpDir := t.TempDir()
		os.Chdir(tmpDir)

		cfg, err := Load()

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if cfg == nil {
			t.Fatal("expected config, got nil")
		}

		if len(cfg.Providers) == 0 {
			t.Error("expected default providers, got none")
		}

		if cfg.Settings.TimeoutMs == 0 {
			t.Error("expected default timeout, got 0")
		}
	})

	t.Run("should load config from file when file exists", func(t *testing.T) {
		origWd, _ := os.Getwd()
		defer os.Chdir(origWd)

		// Create temp config
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "config.yaml")

		// Write test config
		configContent := `
providers:
  test:
    name: Test Provider
    base_url: https://test.com
    model: test-model
    env_key: TEST_API_KEY
    opus_model: test-opus
settings:
  timeout_ms: 5000
  disable_non_essential_calls: false
  log_level: debug
`
		if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
			t.Fatalf("failed to write test config: %v", err)
		}

		// Change to temp dir
		os.Chdir(tmpDir)

		// Debug: check if config file exists
		if _, err := os.Stat(configFile); err != nil {
			t.Fatalf("config file doesn't exist: %v", err)
		}

		cfg, err := Load()

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if cfg == nil {
			t.Fatal("expected config, got nil")
		}

		// Check loaded values
		if cfg.Settings.TimeoutMs != 5000 {
			t.Errorf("expected timeout 5000, got %d", cfg.Settings.TimeoutMs)
		}

		if cfg.Settings.DisableNonEssential {
			t.Errorf("expected disable_non_essential_calls to be false, got %v", cfg.Settings.DisableNonEssential)
		}

		if cfg.Settings.LogLevel != "debug" {
			t.Errorf("expected log_level 'debug', got '%s'", cfg.Settings.LogLevel)
		}
	})
}

func TestGetConfig(t *testing.T) {
	// Initialize config
	cfg, _ := Load()
	currentConfig = cfg

	t.Run("should return config without error", func(t *testing.T) {
		cfg := GetConfig()

		if cfg == nil {
			t.Fatal("expected config, got nil")
		}

		if len(cfg.Providers) == 0 {
			t.Error("expected providers, got none")
		}
	})
}

func TestGetProvider(t *testing.T) {
	cfg, _ := Load()
	currentConfig = cfg

	t.Run("should return existing provider", func(t *testing.T) {
		provider, exists := cfg.GetProvider("mimo")

		if !exists {
			t.Error("expected provider to exist")
		}

		if provider == nil {
			t.Fatal("expected provider, got nil")
		}

		if provider.Name != "Mimo" {
			t.Errorf("expected name 'Mimo', got '%s'", provider.Name)
		}
	})

	t.Run("should return false for non-existent provider", func(t *testing.T) {
		_, exists := cfg.GetProvider("nonexistent")

		if exists {
			t.Error("expected false for non-existent provider")
		}
	})
}

func TestListProviders(t *testing.T) {
	cfg, _ := Load()
	currentConfig = cfg

	t.Run("should return all provider names", func(t *testing.T) {
		providers := cfg.ListProviders()

		if len(providers) == 0 {
			t.Fatal("expected providers, got empty list")
		}

		expectedProviders := []string{"claude", "glm", "kimi", "mimo", "minimax"}

		if len(providers) != len(expectedProviders) {
			t.Errorf("expected %d providers, got %d", len(expectedProviders), len(providers))
		}

		for _, expected := range expectedProviders {
			found := false
			for _, p := range providers {
				if p == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected provider '%s' in list", expected)
			}
		}
	})
}

func TestFindConfigFile(t *testing.T) {
	t.Run("should find config in XDG_CONFIG_HOME", func(t *testing.T) {
		origXdg := os.Getenv("XDG_CONFIG_HOME")
		defer os.Setenv("XDG_CONFIG_HOME", origXdg)

		tmpDir := t.TempDir()
		os.Setenv("XDG_CONFIG_HOME", tmpDir)

		// Create config directory
		configDir := filepath.Join(tmpDir, "cclaude")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			t.Fatalf("failed to create test config dir: %v", err)
		}

		// Create config file
		configFile := filepath.Join(configDir, "config.yaml")
		if err := os.WriteFile(configFile, []byte("providers: {}"), 0644); err != nil {
			t.Fatalf("failed to create test config file: %v", err)
		}

		path, err := findConfigFile()

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if path != configFile {
			t.Errorf("expected path '%s', got '%s'", configFile, path)
		}
	})

	t.Run("should_find config in home directory", func(t *testing.T) {
		origXdg := os.Getenv("XDG_CONFIG_HOME")
		defer os.Setenv("XDG_CONFIG_HOME", origXdg)

		origHome := os.Getenv("HOME")
		defer os.Setenv("HOME", origHome)

		tmpDir := t.TempDir()
		os.Setenv("XDG_CONFIG_HOME", "/nonexistent")
		os.Setenv("HOME", tmpDir)

		configDir := filepath.Join(tmpDir, ".config", "cclaude")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			t.Fatalf("failed to create test config dir: %v", err)
		}

		configFile := filepath.Join(configDir, "config.yaml")
		if err := os.WriteFile(configFile, []byte("providers: {}"), 0644); err != nil {
			t.Fatalf("failed to create test config file: %v", err)
		}

		path, err := findConfigFile()

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if path != configFile {
			t.Errorf("expected path '%s', got '%s'", configFile, path)
		}
	})
		}
	})

	t.Run("should return error when no config exists", func(t *testing.T) {
		origXdg := os.Getenv("XDG_CONFIG_HOME")
		defer os.Setenv("XDG_CONFIG_HOME", origXdg)

		origHome, _ := os.UserHomeDir()
		defer os.Setenv("HOME", origHome)

		os.Setenv("XDG_CONFIG_HOME", "/nonexistent")
		os.Setenv("HOME", "/nonexistent")

		_, err := findConfigFile()

		if err == nil {
			t.Error("expected error when no config file exists")
		}
	})
}

func TestOnConfigChange(t *testing.T) {
	t.Run("should register and notify callback", func(t *testing.T) {
		called := false
		callback := func(cfg *Config) {
			called = true
		}

		OnConfigChange(callback)

		// Trigger a reload by directly calling the unexported changeMutex
		// In real usage, this would be triggered by Watch()
		Reload()

		if !called {
			t.Error("expected callback to be called")
		}
	})
}

func TestConfigReload(t *testing.T) {
	t.Run("should reload config when file changes", func(t *testing.T) {
		origWd, _ := os.Getwd()
		defer os.Chdir(origWd)

		tmpDir := t.TempDir()
		os.Chdir(tmpDir)

		// Create initial config
		configFile := filepath.Join(tmpDir, "config.yaml")
		initialConfig := `
providers:
  test:
    name: Initial
    base_url: https://test.com
    model: v1
settings:
  timeout_ms: 1000
`
		os.WriteFile(configFile, []byte(initialConfig), 0644)

		// Load initial config
		Load()

		// Modify config
		modifiedConfig := `
providers:
  test:
    name: Modified
    base_url: https://modified.com
    model: v2
settings:
  timeout_ms: 2000
`
		os.WriteFile(configFile, []byte(modifiedConfig), 0644)

		// Reload
		err := Reload()

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify new values
		cfg := GetConfig()
	provider, _ := cfg.GetProvider("test")

		if provider.Name != "Modified" {
			t.Errorf("expected name 'Modified', got '%s'", provider.Name)
		}

		if cfg.Settings.TimeoutMs != 2000 {
			t.Errorf("expected timeout 2000, got %d", cfg.Settings.TimeoutMs)
		}
	})
}

func TestGetConfigPaths(t *testing.T) {
	t.Run("should include current directory in paths", func(t *testing.T) {
		origWd, _ := os.Getwd()
		defer os.Chdir(origWd)

		paths := getConfigPaths()

		if len(paths) == 0 {
			t.Fatal("expected at least one path")
		}

		// Current directory should be first
		if paths[0] != origWd {
			t.Errorf("expected first path to be current directory, got '%s'", paths[0])
		}
	})

	t.Run("should include XDG_CONFIG_HOME when set", func(t *testing.T) {
		origXdg := os.Getenv("XDG_CONFIG_HOME")
		defer os.Setenv("XDG_CONFIG_HOME", origXdg)

		expectedPath := "/test/xdg/cclaude"
		os.Setenv("XDG_CONFIG_HOME", "/test/xdg")

		paths := getConfigPaths()

		found := false
		for _, path := range paths {
			if path == expectedPath {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("expected path '%s' in paths", expectedPath)
		}
	})
}
