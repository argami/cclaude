package validator

import (
	"fmt"
	"os"

	"github.com/argami/cclaude/pkg/models"
)

// ValidateAPIKey checks if the API key for the given provider is set
func ValidateAPIKey(providerName string) error {
	provider, err := models.GetProvider(providerName)
	if err != nil {
		return fmt.Errorf("provider not found: %s", providerName)
	}

	apiKey := os.Getenv(provider.EnvKey)
	if apiKey == "" {
		return fmt.Errorf("missing API key for %s: set %s environment variable", providerName, provider.EnvKey)
	}

	return nil
}

// ValidateProvider checks if the provider name is valid
func ValidateProvider(providerName string) error {
	_, err := models.GetProvider(providerName)
	return err
}
