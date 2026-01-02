package cli

import (
	"flag"
	"fmt"
	"os"
)

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

// ParseFlags parses command line flags and returns the Flags struct
func ParseFlags() (*Flags, error) {
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

	return flags, nil
}

// MustParseFlags parses flags and exits on error
func MustParseFlags() *Flags {
	flags, err := ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}
	return flags
}
