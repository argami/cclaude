package utils

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestShowHelp(t *testing.T) {
	// Capturar stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	ShowHelp()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verificar que contiene información clave
	checks := []string{
		"cclaude - Wrapper multi-proveedor",
		"Uso: cclaude <proveedor>",
		"Proveedores:",
		"mimo",
		"minimax",
		"kimi",
		"glm",
		"claude",
		"Flags:",
		"--provider",
		"--help",
		"--version",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("Expected help output to contain '%s'", check)
		}
	}
}

func TestShowVersion(t *testing.T) {
	// Capturar stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	ShowVersion()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verificar formato de versión
	if !strings.Contains(output, "cclaude-go") {
		t.Errorf("Expected version output to contain 'cclaude-go', got: %s", output)
	}
	if !strings.Contains(output, "v1.0.0") {
		t.Errorf("Expected version output to contain 'v1.0.0', got: %s", output)
	}
}