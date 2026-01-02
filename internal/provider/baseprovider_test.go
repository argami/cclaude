package provider

import (
	"os"
	"testing"
)

func TestBaseProvider_Name(t *testing.T) {
	provider := &BaseProvider{name: "test-provider"}

	if provider.Name() != "test-provider" {
		t.Errorf("expected 'test-provider', got '%s'", provider.Name())
	}
}

func TestBaseProvider_Validate(t *testing.T) {
	t.Run("should pass validation when no API key required", func(t *testing.T) {
		provider := &BaseProvider{
			name:   "claude",
			envKey: "",
		}

		err := provider.Validate()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("should pass validation when API key is set", func(t *testing.T) {
		// Setup
		os.Setenv("TEST_API_KEY", "test-key")
		defer os.Unsetenv("TEST_API_KEY")

		provider := &BaseProvider{
			name:   "test",
			envKey: "TEST_API_KEY",
		}

		err := provider.Validate()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("should fail validation when API key is missing", func(t *testing.T) {
		// Ensure env var is not set
		os.Unsetenv("TEST_API_KEY")

		provider := &BaseProvider{
			name:   "test",
			envKey: "TEST_API_KEY",
		}

		err := provider.Validate()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestBaseProvider_SetupEnv(t *testing.T) {
	t.Run("should set all environment variables", func(t *testing.T) {
		// Setup
		provider := &BaseProvider{
			name:      "test",
			baseURL:   "https://api.test.com",
			model:     "test-model",
			opusModel: "test-opus",
		}

		// Execute
		err := provider.SetupEnv()

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify environment variables
		if os.Getenv("ANTHROPIC_BASE_URL") != "https://api.test.com" {
			t.Errorf("expected ANTHROPIC_BASE_URL to be set")
		}
		if os.Getenv("MAIN_MODEL") != "test-model" {
			t.Errorf("expected MAIN_MODEL to be set")
		}
		if os.Getenv("ANTHROPIC_DEFAULT_SONNET_MODEL") != "test-model" {
			t.Errorf("expected ANTHROPIC_DEFAULT_SONNET_MODEL to be set")
		}
		if os.Getenv("CLAUDE_DEFAULT_OPUS_MODEL") != "test-opus" {
			t.Errorf("expected CLAUDE_DEFAULT_OPUS_MODEL to be set")
		}
		if os.Getenv("ANTHROPIC_API_KEY") != "" {
			t.Errorf("expected ANTHROPIC_API_KEY to be empty")
		}
	})

	t.Run("should work with minimal configuration", func(t *testing.T) {
		provider := &BaseProvider{
			name: "claude",
		}

		err := provider.SetupEnv()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}

func TestBaseProvider_GetClaudeArgs(t *testing.T) {
	provider := &BaseProvider{name: "test"}

	cmd, args := provider.GetClaudeArgs()

	if cmd != "claude" {
		t.Errorf("expected 'claude', got '%s'", cmd)
	}
	if len(args) != 0 {
		t.Errorf("expected no args, got %v", args)
	}
}
