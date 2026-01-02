package config

import (
	"fmt"
	"os"
)

// Validate checks if the configuration is valid
func Validate(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	// Check if providers exist
	if len(cfg.Providers) == 0 {
		return fmt.Errorf("no providers configured")
	}

	// Validate each provider
	for name, provider := range cfg.Providers {
		if provider == nil {
			return fmt.Errorf("provider %s is nil", name)
		}

		// Check for required fields
		if provider.Name == "" {
			return fmt.Errorf("provider %s has empty name", name)
		}

		// Check if API key is set for providers that require it
		if provider.EnvKey != "" {
			if os.Getenv(provider.EnvKey) == "" {
				return fmt.Errorf("provider %s requires %s environment variable", name, provider.EnvKey)
			}
		}
	}

	return nil
}

// ValidateProvider checks if a specific provider is valid
func ValidateProvider(name string) error {
	cfg, err := Load()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	if err := Validate(cfg); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	provider, exists := cfg.Providers[name]
	if !exists {
		return fmt.Errorf("proveedor '%s' no configurado", name)
	}

	// Check API key if required
	if provider.EnvKey != "" {
		if os.Getenv(provider.EnvKey) == "" {
			return fmt.Errorf("API key no definida para %s: %s", name, provider.EnvKey)
		}
	}

	return nil
}
