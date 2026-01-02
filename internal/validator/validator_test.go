package validator

import (
	"os"
	"testing"
)

func TestValidateAPIKey(t *testing.T) {
	// Ensure MINIMAX_API_KEY is not set for these tests
	os.Unsetenv("MINIMAX_API_KEY")

	tests := []struct {
		name         string
		providerName string
		setEnvKey    string
		wantErr      bool
	}{
		{
			name:         "valid provider with API key",
			providerName: "mimo",
			setEnvKey:    "MIMO_API_KEY",
			wantErr:      false,
		},
		{
			name:         "valid provider without API key",
			providerName: "minimax",
			setEnvKey:    "",
			wantErr:      true,
		},
		{
			name:         "invalid provider",
			providerName: "invalid",
			setEnvKey:    "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable if needed
			if tt.setEnvKey != "" {
				os.Setenv(tt.setEnvKey, "test-key-value")
				defer os.Unsetenv(tt.setEnvKey)
			}

			err := ValidateAPIKey(tt.providerName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateAPIKey(%s) expected error, got nil", tt.providerName)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateAPIKey(%s) unexpected error: %v", tt.providerName, err)
				}
			}
		})
	}
}

func TestValidateProvider(t *testing.T) {
	tests := []struct {
		name         string
		providerName string
		wantErr      bool
	}{
		{
			name:         "valid provider",
			providerName: "mimo",
			wantErr:      false,
		},
		{
			name:         "invalid provider",
			providerName: "nonexistent",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProvider(tt.providerName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateProvider(%s) expected error, got nil", tt.providerName)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateProvider(%s) unexpected error: %v", tt.providerName, err)
				}
			}
		})
	}
}
