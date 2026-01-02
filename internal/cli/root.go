package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/argami/cclaude-glm/internal/execution"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "cclaude",
		Short: "Claude Code wrapper para múltiples proveedores de IA",
		Long: `cclaude es un wrapper para Claude Code CLI que permite usar
diferentes proveedores de IA (Mimo, MiniMax, Kimi, GLM, Claude nativo).

USO:
    cclaude <proveedor> [opciones de claude]

PROVEEDORES:
    mimo       Xiaomi Mimo v2 Flash
    minimax    MiniMax M2.1
    kimi       Moonshot Kimi K2
    glm        Zhipu GLM-4.7
    claude     Anthropic Claude (nativo)

EJEMPLOS:
    cclaude glm "Explica este código"
    cclaude --list-providers
    cclaude -v --version`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If no provider specified and no args, show help
			if len(args) == 0 {
				return cmd.Help()
			}

			// First arg is the provider
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

			// Execute with provider using execution package
			return execution.RunProvider(provider, remainingArgs)
		},
	}

	rootCmd.AddCommand(NewListCommand())
	rootCmd.AddCommand(NewVersionCommand())
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(NewConfigPathCommand())
	rootCmd.AddCommand(NewValidateCommand())

	return rootCmd
}
