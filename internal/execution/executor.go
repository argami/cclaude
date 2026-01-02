package execution

import (
	"fmt"
	"os"
	"os/exec"
	"github.com/argami/cclaude-glm/internal/provider"
)

type Executor struct {
	provider provider.Provider
}

func NewExecutor(p provider.Provider) *Executor {
	return &Executor{provider: p}
}

// RunProvider is a convenience function that creates a provider and executes
func RunProvider(providerName string, args []string) error {
	// Create provider using factory
	p, err := provider.Factory(providerName)
	if err != nil {
		return err
	}

	executor := NewExecutor(p)
	return executor.Execute(args)
}

func (e *Executor) Execute(args []string) error {
	// Validate provider configuration
	if err := e.provider.Validate(); err != nil {
		return err
	}

	// Setup environment variables
	if err := e.provider.SetupEnv(); err != nil {
		return err
	}

	// Get claude command and args
	cmdName, cmdArgs := e.provider.GetClaudeArgs()

	// Prepare claude args - prepend user args
	claudeArgs := append([]string{cmdName}, cmdArgs...)
	claudeArgs = append(claudeArgs, args...)

	// Check if claude binary exists
	cmdPath, err := exec.LookPath(claudeArgs[0])
	if err != nil {
		return fmt.Errorf("claude CLI no encontrado en PATH: %w", err)
	}

	if cmdPath == "" {
		return fmt.Errorf("claude CLI no encontrado en PATH")
	}

	// Execute claude
	cmd := exec.Command(cmdPath, claudeArgs[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
