package flags

import (
	"testing"
)

func TestParseProviderPositional(t *testing.T) {
	config, err := ParseWithArgs([]string{"mimo"})
	if err != nil {
		t.Fatalf("ParseWithArgs() error = %v", err)
	}

	if config.Provider != "mimo" {
		t.Errorf("Expected provider 'mimo', got '%s'", config.Provider)
	}
}

func TestParseProviderFlag(t *testing.T) {
	config, err := ParseWithArgs([]string{"--provider", "minimax"})
	if err != nil {
		t.Fatalf("ParseWithArgs() error = %v", err)
	}

	if config.Provider != "minimax" {
		t.Errorf("Expected provider 'minimax', got '%s'", config.Provider)
	}
}

func TestParseShortFlag(t *testing.T) {
	config, err := ParseWithArgs([]string{"-p", "kimi"})
	if err != nil {
		t.Fatalf("ParseWithArgs() error = %v", err)
	}

	if config.Provider != "kimi" {
		t.Errorf("Expected provider 'kimi', got '%s'", config.Provider)
	}
}

func TestParseAllFlags(t *testing.T) {
	config, err := ParseWithArgs([]string{
		"--provider", "glm",
		"--timeout", "10m",
		"--debug",
		"--model", "custom-model",
		"--config", "/path/to/config",
		"arg1", "arg2",
	})
	if err != nil {
		t.Fatalf("ParseWithArgs() error = %v", err)
	}

	if config.Provider != "glm" {
		t.Errorf("Expected provider 'glm', got '%s'", config.Provider)
	}
	if config.Timeout != "10m" {
		t.Errorf("Expected timeout '10m', got '%s'", config.Timeout)
	}
	if !config.Debug {
		t.Error("Expected Debug=true")
	}
	if config.ModelOverride != "custom-model" {
		t.Errorf("Expected model 'custom-model', got '%s'", config.ModelOverride)
	}
	if config.ConfigFile != "/path/to/config" {
		t.Errorf("Expected config '/path/to/config', got '%s'", config.ConfigFile)
	}
	if len(config.Args) != 2 {
		t.Errorf("Expected 2 args, got %d", len(config.Args))
	}
}

func TestParseHelp(t *testing.T) {
	config, err := ParseWithArgs([]string{"--help"})
	if err != nil {
		t.Fatalf("ParseWithArgs() error = %v", err)
	}

	if !config.Help {
		t.Error("Expected Help=true")
	}
}

func TestParseVersion(t *testing.T) {
	config, err := ParseWithArgs([]string{"--version"})
	if err != nil {
		t.Fatalf("ParseWithArgs() error = %v", err)
	}

	if !config.Version {
		t.Error("Expected Version=true")
	}
}

func TestParseClaudeArgs(t *testing.T) {
	config, err := ParseWithArgs([]string{"mimo", "--help", "some-arg"})
	if err != nil {
		t.Fatalf("ParseWithArgs() error = %v", err)
	}

	if config.Provider != "mimo" {
		t.Errorf("Expected provider 'mimo', got '%s'", config.Provider)
	}
	if len(config.Args) != 2 {
		t.Errorf("Expected 2 args, got %d: %v", len(config.Args), config.Args)
	}
}

func TestParseEmpty(t *testing.T) {
	config, err := ParseWithArgs([]string{})
	if err != nil {
		t.Fatalf("ParseWithArgs() error = %v", err)
	}

	if config.Provider != "" {
		t.Errorf("Expected empty provider, got '%s'", config.Provider)
	}
}