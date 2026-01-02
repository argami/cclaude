package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/argami/cclaude/internal/config"
	"github.com/argami/cclaude/internal/provider"
	"github.com/argami/cclaude/internal/validator"
	"github.com/argami/cclaude/pkg/models"
)

// Version is set at build time via LDFLAGS
var version = "dev"

// Flags holds all CLI flags and arguments
type Flags struct {
	Provider string
	Model    string
	DryRun   bool
	Verbose  bool
	Version  bool
	Help     bool
	Args     []string
}

func parseFlags() *Flags {
	flags := &Flags{}

	flag.StringVar(&flags.Provider, "provider", "", "Provider to use (mimo, minimax, kimi, glm)")
	flag.StringVar(&flags.Model, "model", "", "Model to use for the provider")
	flag.BoolVar(&flags.DryRun, "dry-run", false, "Show configuration without executing")
	flag.BoolVar(&flags.DryRun, "n", false, "Show configuration without executing (shorthand)")
	flag.BoolVar(&flags.Verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&flags.Verbose, "v", false, "Enable verbose output (shorthand)")
	flag.BoolVar(&flags.Version, "version", false, "Show version information")
	flag.BoolVar(&flags.Help, "help", false, "Show help information")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: cclaude [FLAGS] [ARGUMENTS]\n\n")
		fmt.Fprintf(os.Stderr, "A wrapper for Claude Code with multiple provider support.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  cclaude --provider minimax \"Analyze this code\"\n")
		fmt.Fprintf(os.Stderr, "  cclaude --dry-run --provider mimo \"Hello\"\n")
		fmt.Fprintf(os.Stderr, "  cclaude --verbose \"Complex task\"\n")
	}

	flag.Parse()

	// Collect positional arguments
	flags.Args = flag.Args()

	return flags
}

func main() {
	// Parse flags
	flags := parseFlags()

	// Handle special flags
	if flags.Help {
		flag.Usage()
		return
	}

	if flags.Version {
		fmt.Printf("cclaude version %s\n", version)
		return
	}

	// Load configuration
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "cclaude", "cclaude.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = "./config/cclaude.yaml"
	}
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load config: %v\n", err)
		cfg, _ = config.LoadConfig("") // Fall back to defaults
	}

	// Determine provider
	providerName := flags.Provider
	if providerName == "" {
		providerName = cfg.GetProvider()
	}

	// If no provider specified, run claude directly
	if providerName == "" {
		runClaude(flags.Args)
		return
	}

	// Validate provider
	if err := validator.ValidateProvider(providerName); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Available providers: %v\n", models.GetProviderNames())
		os.Exit(1)
	}

	// Validate API key
	if err := validator.ValidateAPIKey(providerName); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Handle dry-run
	if flags.DryRun {
		if err := provider.ExecuteDryRun(providerName, flags.Args, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Execute with provider
	if err := provider.Execute(providerName, flags.Args, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// runClaude executes claude directly without any provider wrapper
func runClaude(args []string) {
	cmd := exec.Command("claude", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
}
