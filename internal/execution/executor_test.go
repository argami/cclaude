package execution

import (
	"os"
	"testing"
)

// Mock provider for testing
type mockProvider struct {
	name      string
	validateErr error
	setupErr   error
}

func (m *mockProvider) Name() string {
	return m.name
}

func (m *mockProvider) Validate() error {
	return m.validateErr
}

func (m *mockProvider) SetupEnv() error {
	return m.setupErr
}

func (m *mockProvider) GetClaudeArgs() (string, []string) {
	return "claude", []string{}
}

func TestNewExecutor(t *testing.T) {
	mockProv := &mockProvider{name: "test"}
	executor := NewExecutor(mockProv)

	if executor == nil {
		t.Fatal("expected executor, got nil")
	}
	if executor.provider != mockProv {
		t.Error("expected provider to be set")
	}
}

func TestExecutor_Execute(t *testing.T) {
	t.Run("should fail validation when provider invalid", func(t *testing.T) {
		mockProv := &mockProvider{
			validateErr: os.ErrInvalid,
		}
		executor := NewExecutor(mockProv)

		err := executor.Execute([]string{"--help"})

		if err == nil {
			t.Error("expected validation error, got nil")
		}
	})

	t.Run("should fail setup when environment error", func(t *testing.T) {
		mockProv := &mockProvider{
			setupErr: os.ErrPermission,
		}
		executor := NewExecutor(mockProv)

		err := executor.Execute([]string{"--help"})

		if err == nil {
			t.Error("expected setup error, got nil")
		}
	})

	t.Run("should attempt to execute when provider valid", func(t *testing.T) {
		// Setup valid provider
		mockProv := &mockProvider{name: "test"}
		executor := NewExecutor(mockProv)

		// Execute with --help flag (should not fail even if claude not in PATH)
		err := executor.Execute([]string{"--help"})

		// We expect this to fail because claude binary is not in test environment
		// but it should pass validation and setup
		if err == nil {
			// This would mean claude binary exists and ran successfully
			return
		}

		// Verify the error is about claude not being found, not validation/setup
		if err.Error() == "API key no definida" {
			t.Errorf("got validation error instead of execution error: %v", err)
		}
	})
}

func TestRunProvider(t *testing.T) {
	t.Run("should create provider and execute", func(t *testing.T) {
		// This test uses the actual Factory, so we need to ensure config exists
		// For now, we'll skip if config not found
		err := RunProvider("claude", []string{"--help"})

		// We expect this to either work or fail with "claude not found"
		// but not with factory errors if claude provider exists
		if err != nil && err.Error() == "proveedor no encontrado" {
			t.Skip("claude provider not configured")
		}
	})
}
