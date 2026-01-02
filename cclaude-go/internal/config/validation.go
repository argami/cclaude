package config

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/argami/cclaude-go/pkg/types"
)

func ValidateEnvironment() error {
	// Verificar que 'claude' está disponible en el PATH
	_, err := exec.LookPath("claude")
	if err != nil {
		return fmt.Errorf("comando 'claude' no encontrado en PATH: %w", err)
	}
	return nil
}

func ValidateAPIKey(provider types.ProviderConfig) error {
	key := os.Getenv(provider.EnvVar)
	if key == "" {
		return fmt.Errorf("variable de entorno %s no configurada", provider.EnvVar)
	}
	if len(key) < 8 {
		return fmt.Errorf("API key inusualmente corta para %s", provider.Name)
	}
	return nil
}