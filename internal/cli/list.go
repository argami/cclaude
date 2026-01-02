package cli

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/argami/cclaude-glm/internal/config"
)

func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Aliases: []string{"ls"},
		Short: "Lista todos los proveedores disponibles",
		Long:  `Lista todos los proveedores de IA soportados por cclaude.

Muestra:
    - Nombre del proveedor
    - URL base
    - Modelo por defecto
    - Estado de configuración de API key`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("error cargando configuración: %w", err)
			}

			fmt.Println("Proveedores disponibles:")
			fmt.Println()

			for name, provider := range cfg.Providers {
				status := "✅ Configurado"

				// Check if API key is set
				if provider.EnvKey != "" {
					if os.Getenv(provider.EnvKey) == "" {
						status = "❌ Falta API key: " + provider.EnvKey
					}
				}

				fmt.Printf("  %-10s %s\n", name, status)

				if provider.BaseURL != "" {
					fmt.Printf("             URL: %s\n", provider.BaseURL)
				}
				if provider.Model != "" {
					fmt.Printf("             Modelo: %s\n", provider.Model)
				}
				fmt.Println()
			}

			return nil
		},
	}

	return cmd
}
