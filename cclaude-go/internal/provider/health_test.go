package provider

import (
	"os"
	"testing"
	"time"
)

func TestHealthChecker(t *testing.T) {
	// Configurar variables de entorno para tests
	os.Setenv("MIMO_API_KEY", "test-key-12345678")
	os.Setenv("MINIMAX_API_KEY", "test-key-87654321")
	defer os.Unsetenv("MIMO_API_KEY")
	defer os.Unsetenv("MINIMAX_API_KEY")

	hc := NewHealthChecker()

	// Test CheckAll
	results := hc.CheckAll()
	if len(results) == 0 {
		t.Error("Expected results from CheckAll")
	}

	// Verificar que todos los proveedores están presentes
	if len(results) != len(Providers) {
		t.Errorf("Expected %d results, got %d", len(Providers), len(results))
	}
}

func TestHealthCheckerCheckProvider(t *testing.T) {
	os.Setenv("MIMO_API_KEY", "test-key-12345678")
	defer os.Unsetenv("MIMO_API_KEY")

	hc := NewHealthChecker()

	// Test proveedor válido
	result := hc.CheckProvider("mimo")
	if result.Provider != "mimo" {
		t.Errorf("Expected provider 'mimo', got '%s'", result.Provider)
	}

	// Test proveedor inexistente
	result = hc.CheckProvider("invalid")
	if result.Provider != "invalid" {
		t.Error("Expected result for invalid provider")
	}
	if result.Healthy {
		t.Error("Invalid provider should not be healthy")
	}
}

func TestHealthCheckerCheckWithTimeout(t *testing.T) {
	os.Setenv("MIMO_API_KEY", "test-key-12345678")
	defer os.Unsetenv("MIMO_API_KEY")

	hc := NewHealthChecker()

	// Test con timeout personalizado
	result := hc.CheckWithTimeout("mimo", 1*time.Second)
	if result.Provider != "mimo" {
		t.Errorf("Expected provider 'mimo', got '%s'", result.Provider)
	}
}

func TestHealthCheckerVerifyAPIKey(t *testing.T) {
	hc := NewHealthChecker()

	// Test API key vacía
	valid, msg := hc.VerifyAPIKey("mimo", "")
	if valid {
		t.Error("Empty API key should be invalid")
	}
	if msg != "API key vacía" {
		t.Errorf("Expected 'API key vacía', got '%s'", msg)
	}

	// Test API key demasiado corta
	valid, msg = hc.VerifyAPIKey("mimo", "short")
	if valid {
		t.Error("Short API key should be invalid")
	}

	// Test API key válida
	valid, msg = hc.VerifyAPIKey("mimo", "test-key-12345678")
	if !valid {
		t.Error("Valid API key should be accepted")
	}

	// Test proveedor inexistente
	valid, msg = hc.VerifyAPIKey("invalid", "test-key")
	if valid {
		t.Error("Invalid provider should be rejected")
	}

	// Test Claude nativo
	valid, msg = hc.VerifyAPIKey("claude", "")
	if !valid {
		t.Errorf("Claude nativo should not require API key, got: %s", msg)
	}
}

func TestHealthCheckerGetProviderStats(t *testing.T) {
	hc := NewHealthChecker()

	stats := hc.GetProviderStats("mimo")
	if stats == nil {
		t.Error("Expected stats for valid provider")
	}

	if stats["name"] != "mimo" {
		t.Errorf("Expected name 'mimo', got '%v'", stats["name"])
	}
}

func TestHealthCheckerRunDiagnostics(t *testing.T) {
	os.Setenv("MIMO_API_KEY", "test-key-12345678")
	defer os.Unsetenv("MIMO_API_KEY")

	hc := NewHealthChecker()

	diagnostics := hc.RunDiagnostics()
	if diagnostics == nil {
		t.Error("Expected diagnostics")
	}

	// Verificar estructura de diagnósticos
	if _, ok := diagnostics["health_results"]; !ok {
		t.Error("Expected health_results in diagnostics")
	}
	if _, ok := diagnostics["summary"]; !ok {
		t.Error("Expected summary in diagnostics")
	}
	if _, ok := diagnostics["environment"]; !ok {
		t.Error("Expected environment in diagnostics")
	}
}

func TestHealthCheckerFormatHealthResults(t *testing.T) {
	hc := NewHealthChecker()

	// Test con resultados vacíos
	results := []HealthCheckResult{}
	output := hc.FormatHealthResults(results, false)
	if output == "" {
		t.Error("Expected non-empty output")
	}

	// Test con resultados
	result := HealthCheckResult{
		Provider:  "mimo",
		Healthy:   true,
		Latency:   100 * time.Millisecond,
		Timestamp: time.Now(),
	}
	results = []HealthCheckResult{result}
	output = hc.FormatHealthResults(results, true)
	if output == "" {
		t.Error("Expected non-empty output")
	}
	if !contains(output, "mimo") {
		t.Error("Output should contain provider name")
	}
}

// Función auxiliar
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}