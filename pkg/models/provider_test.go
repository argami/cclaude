package models

import (
	"testing"
)

func TestGetProvider(t *testing.T) {
	tests := []struct {
		name          string
		providerName  string
		wantName      string
		wantBaseURL   string
		wantModel     string
		wantEnvKey    string
		wantErr       bool
	}{
		{
			name:         "mimo provider",
			providerName: "mimo",
			wantName:     "mimo",
			wantBaseURL:  "https://api.xiaomimimo.com/anthropic",
			wantModel:    "mimo-v2-flash",
			wantEnvKey:   "MIMO_API_KEY",
			wantErr:      false,
		},
		{
			name:         "minimax provider",
			providerName: "minimax",
			wantName:     "minimax",
			wantBaseURL:  "https://api.minimax.io/anthropic",
			wantModel:    "MiniMax-M2.1",
			wantEnvKey:   "MINIMAX_API_KEY",
			wantErr:      false,
		},
		{
			name:         "kimi provider",
			providerName: "kimi",
			wantName:     "kimi",
			wantBaseURL:  "https://api.kimi.com/coding/",
			wantModel:    "kimi-k2-0711-preview",
			wantEnvKey:   "KIMI_API_KEY",
			wantErr:      false,
		},
		{
			name:         "glm provider",
			providerName: "glm",
			wantName:     "glm",
			wantBaseURL:  "https://api.z.ai/api/anthropic",
			wantModel:    "glm-4.7",
			wantEnvKey:   "GLM_API_KEY",
			wantErr:      false,
		},
		{
			name:         "invalid provider",
			providerName: "invalid",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := GetProvider(tt.providerName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetProvider(%s) expected error, got nil", tt.providerName)
				}
				return
			}

			if err != nil {
				t.Errorf("GetProvider(%s) unexpected error: %v", tt.providerName, err)
				return
			}

			if provider.Name != tt.wantName {
				t.Errorf("Name = %v, want %v", provider.Name, tt.wantName)
			}
			if provider.BaseURL != tt.wantBaseURL {
				t.Errorf("BaseURL = %v, want %v", provider.BaseURL, tt.wantBaseURL)
			}
			if provider.Model != tt.wantModel {
				t.Errorf("Model = %v, want %v", provider.Model, tt.wantModel)
			}
			if provider.EnvKey != tt.wantEnvKey {
				t.Errorf("EnvKey = %v, want %v", provider.EnvKey, tt.wantEnvKey)
			}
		})
	}
}

func TestGetProviderNames(t *testing.T) {
	names := GetProviderNames()

	expected := []string{"mimo", "minimax", "kimi", "glm"}

	if len(names) != len(expected) {
		t.Errorf("GetProviderNames() returned %d names, want %d", len(names), len(expected))
	}

	for _, name := range expected {
		found := false
		for _, n := range names {
			if n == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetProviderNames() missing %s", name)
		}
	}
}
