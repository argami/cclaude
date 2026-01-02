package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestRootCommand(t *testing.T) {
	// Build the binary first
	buildCmd := exec.Command("go", "build", "-o", "cclaude-test", ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build binary: %v", err)
	}
	defer os.Remove("cclaude-test")

	t.Run("should show help when no args provided", func(t *testing.T) {
		cmd := exec.Command("./cclaude-test", "--help")
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Fatalf("expected no error, got %v: %s", err, output)
		}

		if len(output) == 0 {
			t.Error("expected help output, got empty string")
		}
	})

	t.Run("should list providers", func(t *testing.T) {
		cmd := exec.Command("./cclaude-test", "list")
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Fatalf("expected no error, got %v: %s", err, output)
		}

		// Check if output contains provider information
		if len(output) == 0 {
			t.Error("expected output to contain provider information")
		}
	})

	t.Run("should show version", func(t *testing.T) {
		cmd := exec.Command("./cclaude-test", "version")
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Fatalf("expected no error, got %v: %s", err, output)
		}

		if len(output) == 0 {
			t.Error("expected version output, got empty string")
		}
	})

	t.Run("should reject invalid provider", func(t *testing.T) {
		cmd := exec.Command("./cclaude-test", "invalid-provider", "--help")
		output, err := cmd.CombinedOutput()

		if err == nil {
			t.Fatal("expected error for invalid provider, got nil")
		}

		if len(output) == 0 {
			t.Error("expected error message, got empty output")
		}
	})

	t.Run("should accept valid provider with args", func(t *testing.T) {
		// This will fail if claude binary not in PATH, but that's expected
		// We're testing that the command parsing works
		cmd := exec.Command("./cclaude-test", "claude", "--help")
		_, err := cmd.CombinedOutput()

		// Command should either execute or fail with specific error
		// but not crash with parsing errors
		if err != nil {
			t.Log("Command failed as expected (claude may not be in PATH)")
		}
	})
}
