package provider

import (
	"fmt"
	"github.com/argami/cclaude-glm/internal/config"
)

// Factory creates provider instances based on provider name
func Factory(providerName string) (Provider, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("error cargando configuración: %w", err)
	}

	providerConfig, exists := cfg.Providers[providerName]
	if !exists {
		return nil, fmt.Errorf("proveedor no encontrado: %s", providerName)
	}

	return &BaseProvider{
		name:      providerConfig.Name,
		baseURL:   providerConfig.BaseURL,
		model:     providerConfig.Model,
		envKey:    providerConfig.EnvKey,
		opusModel: providerConfig.OpusModel,
	}, nil
}
