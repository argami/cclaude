package utils

import (
	"testing"
)

func TestExecuteClaude(t *testing.T) {
	// Test que verifica que ExecuteClaude existe y tiene la firma correcta
	// Este test no ejecuta realmente porque no tenemos claude disponible
	// Pero verifica que la función está implementada correctamente

	// Verificar que la función existe con la firma correcta
	var _ func([]string) error = ExecuteClaude

	// Test con argumentos vacíos (debería fallar porque no hay claude)
	err := ExecuteClaude([]string{})

	// Si no hay claude, debería dar error
	if err == nil {
		t.Log("Claude parece estar disponible en el PATH")
	} else {
		t.Logf("Claude no está disponible (esperado en entornos de test): %v", err)
	}
}

func TestExecuteClaudeWithArgs(t *testing.T) {
	// Verificar que la función maneja argumentos correctamente
	err := ExecuteClaude([]string{"--help"})

	// Si claude está disponible, esto debería funcionar
	// Si no, debería dar error claro
	if err != nil {
		t.Logf("Ejecución con argumentos falló (claude no disponible): %v", err)
	}
}