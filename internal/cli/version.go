package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewVersionCommand() *cobra.Command {
	version := "v0.1.0"

	cmd := &cobra.Command{
		Use:   "version",
		Aliases: []string{"v"},
		Short: "Muestra la versión de cclaude",
		Long:  `Muestra información de la versión de cclaude incluyendo número de versión y fecha de compilación.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("cclaude %s\n", version)
			fmt.Println("Claude Code wrapper para múltiples proveedores de IA")
			return nil
		},
	}

	return cmd
}
