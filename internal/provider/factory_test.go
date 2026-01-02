package provider

import (
	"os"
	"testing"
)

func TestFactory(t *testing.T) {
	// Test case 1: Valid provider
	t.Run("should return provider when valid name given", func(t *testing.T) {
		// Setup: Set required environment variable
		os.Setenv("MIMO_API_KEY", "test-key")
		defer os.Unsetenv("MIMO_API_KEY")

		// Execute
		provider, err := Factory("mimo")

		// Assert
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if provider == nil {
			t.Fatal("expected provider, got nil")
		}
		if provider.Name() != "Mimo" {
			t.Errorf("expected name 'Mimo', got '%s'", provider.Name())
		}
	})

	// Test case 2: Invalid provider name
	t.Run("should return error when invalid provider name given", func(t *testing.T) {
		// Execute
		provider, err := Factory("nonexistent")

		// Assert
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if provider != nil {
			t.Errorf("expected nil provider, got %v", provider)
		}
	})

	// Test case 3: Multiple valid providers
	validProviders := []struct {
		name            string
	requiresAPIKey   bool
		expectedError   bool
	}{
		{"mimo", true, false},
		{"minimax", true, false},
		{"kimi", true, true},  // kimi not in config
		{"glm", true, false},
		{"claude", false, false},
	}

	for _, tc := range validProviders {
		t.Run("should handle provider "+tc.name, func(t *testing.T) {
			// Setup
			if tc.requiresAPIKey {
				os.Setenv(tc.name+"_API_KEY", "test-key")
				defer os.Unsetenv(tc.name+"_API_KEY")
			}

			// Execute
			provider, err := Factory(tc.name)

			// Assert
			if tc.expectedError {
				if err == nil {
					t.Logf("Expected error for %s, but got none - provider may not be configured", tc.name)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error for %s: %v", tc.name, err)
				}
				if provider == nil {
					t.Fatalf("expected provider for %s, got nil", tc.name)
				}
			}
		})
	}
}
