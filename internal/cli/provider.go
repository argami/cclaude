package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/argami/cclaude-glm/internal/execution"
)

func NewProviderCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "<proveedor>",
		Short: "Ejecuta claude con el proveedor especificado",
		Long:  `Ejecuta Claude Code CLI usando el proveedor de IA especificado.

PROVEEDORES DISPONIBLES:
    mimo       Xiaomi Mimo v2 Flash
    minimax    MiniMax M2.1
    kimi       Moonshot Kimi K2
    glm        Zhipu GLM-4.7
    claude     Anthropic Claude (nativo)

EJEMPLOS:
    cclaude mimo "Explica este código"
    cclaude glm "Ayúdame con este error"
    cclaude minimax --version`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			provider := args[0]
			remainingArgs := args[1:]

			// Validate provider
			validProviders := map[string]bool{
				"mimo":    true,
				"minimax":  true,
				"kimi":     true,
				"glm":      true,
				"claude":   true,
			}

			if !validProviders[provider] {
				return fmt.Errorf("proveedor no válido: %s\nproveedores disponibles: mimo, minimax, kimi, glm, claude", provider)
			}

			// Execute with provider
			return execution.RunProvider(provider, remainingArgs)
		},
	}

	return cmd
}
