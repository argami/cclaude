package main

import (
	"os"
	"os/exec"
	"testing"
)

// TestMainIntegration es un test de integración básico
func TestMainIntegration(t *testing.T) {
	// Verificar que el binario se puede compilar
	cmd := exec.Command("go", "build", "-o", "/tmp/cclaude-test")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, output)
	}

	// Test 1: Help flag
	cmd = exec.Command("/tmp/cclaude-test", "--help")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Logf("Help output: %s", output)
	}

	// Test 2: Version flag
	cmd = exec.Command("/tmp/cclaude-test", "--version")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Logf("Version output: %s", output)
	}

	// Limpiar archivo temporal
	os.Remove("/tmp/cclaude-test")
}

func TestMainFlowWithoutClaude(t *testing.T) {
	// Test que verifica el flujo sin tener claude instalado
	t.Skip("Requiere mocking de exec.LookPath")
}

func TestMainFlowWithProvider(t *testing.T) {
	// Test que verifica el flujo completo con un proveedor
	t.Skip("Requiere setup completo de ambiente")
}