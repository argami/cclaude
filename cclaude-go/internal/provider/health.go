package provider

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/argami/cclaude-go/pkg/types"
)

// HealthCheckResult representa el resultado de un health check
type HealthCheckResult struct {
	Provider  string
	Healthy   bool
	Latency   time.Duration
	Error     string
	Timestamp time.Time
}

// HealthChecker verifica la salud de los proveedores
type HealthChecker struct {
	providers map[string]types.ProviderConfig
	timeout   time.Duration
}

// NewHealthChecker crea un nuevo verificador de salud
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		providers: Providers,
		timeout:   5 * time.Second,
	}
}

// CheckAll verifica todos los proveedores
func (hc *HealthChecker) CheckAll() []HealthCheckResult {
	results := make([]HealthCheckResult, 0, len(hc.providers))

	for name := range hc.providers {
		result := hc.CheckProvider(name)
		results = append(results, result)
	}

	return results
}

// CheckProvider verifica un proveedor específico
func (hc *HealthChecker) CheckProvider(providerName string) HealthCheckResult {
	provider, exists := hc.providers[providerName]
	if !exists {
		return HealthCheckResult{
			Provider:  providerName,
			Healthy:   false,
			Error:     "Proveedor no encontrado",
			Timestamp: time.Now(),
		}
	}

	// Verificar que la variable de entorno está configurada
	apiKey := provider.EnvVar
	if apiKey == "" {
		return HealthCheckResult{
			Provider:  providerName,
			Healthy:   false,
			Error:     "API key no configurada",
			Timestamp: time.Now(),
		}
	}

	// Verificar conectividad al endpoint (sin enviar datos reales)
	start := time.Now()
	result := HealthCheckResult{
		Provider:  providerName,
		Timestamp: start,
	}

	// Crear contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), hc.timeout)
	defer cancel()

	// Crear petición HEAD para verificar disponibilidad
	req, err := http.NewRequestWithContext(ctx, "HEAD", provider.BaseURL, nil)
	if err != nil {
		result.Healthy = false
		result.Error = fmt.Sprintf("Error creando petición: %v", err)
		return result
	}

	// Añadir headers básicos (simulando comportamiento de Claude)
	req.Header.Set("User-Agent", "cclaude-health-check/1.0")
	req.Header.Set("Accept", "application/json")

	// Ejecutar petición
	client := &http.Client{Timeout: hc.timeout}
	resp, err := client.Do(req)

	latency := time.Since(start)
	result.Latency = latency

	if err != nil {
		result.Healthy = false
		result.Error = fmt.Sprintf("Error de conexión: %v", err)
		return result
	}
	defer resp.Body.Close()

	// Verificar código de estado
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		result.Healthy = true
		result.Error = ""
	} else {
		result.Healthy = false
		result.Error = fmt.Sprintf("Código de estado: %d", resp.StatusCode)
	}

	return result
}

// CheckProviderWithAPIKey verifica un proveedor con API key específica
func (hc *HealthChecker) CheckProviderWithAPIKey(providerName string, apiKey string) HealthCheckResult {
	result := hc.CheckProvider(providerName)

	// Si el proveedor existe pero no hay API key, usamos la proporcionada
	if !result.Healthy && result.Error == "API key no configurada" && apiKey != "" {
		// Reintentar con la API key proporcionada
		provider := hc.providers[providerName]
		originalEnvVar := provider.EnvVar
		provider.EnvVar = apiKey

		result = hc.CheckProvider(providerName)

		// Restaurar
		provider.EnvVar = originalEnvVar
	}

	return result
}

// GetHealthSummary retorna un resumen de salud del sistema
func (hc *HealthChecker) GetHealthSummary() string {
	results := hc.CheckAll()

	healthyCount := 0
	for _, result := range results {
		if result.Healthy {
			healthyCount++
		}
	}

	return fmt.Sprintf("Health Status: %d/%d providers healthy", healthyCount, len(results))
}

// FormatHealthResults formatea los resultados para display
func (hc *HealthChecker) FormatHealthResults(results []HealthCheckResult, verbose bool) string {
	if len(results) == 0 {
		return "No hay proveedores para verificar"
	}

	output := "🔍 Health Check Results:\n\n"

	for _, result := range results {
		status := "✅"
		if !result.Healthy {
			status = "❌"
		}

		output += fmt.Sprintf("%s %s: ", status, result.Provider)

		if result.Healthy {
			output += fmt.Sprintf("Healthy (latency: %v)\n", result.Latency)
		} else {
			output += fmt.Sprintf("Unhealthy - %s\n", result.Error)
		}

		if verbose && result.Latency > 0 {
			output += fmt.Sprintf("   Latency: %v\n", result.Latency)
		}
	}

	return output
}

// CheckWithTimeout verifica un proveedor con timeout personalizado
func (hc *HealthChecker) CheckWithTimeout(providerName string, timeout time.Duration) HealthCheckResult {
	originalTimeout := hc.timeout
	hc.timeout = timeout
	defer func() { hc.timeout = originalTimeout }()

	return hc.CheckProvider(providerName)
}

// VerifyAPIKey verifica que una API key sea válida para un proveedor
func (hc *HealthChecker) VerifyAPIKey(providerName string, apiKey string) (bool, string) {
	// Verificar si el proveedor existe
	_, exists := hc.providers[providerName]
	if !exists {
		return false, "Proveedor no encontrado"
	}

	// Claude nativo no necesita API key
	if providerName == "claude" {
		return true, "Claude nativo no requiere API key"
	}

	// Verificar API key para otros proveedores
	if apiKey == "" {
		return false, "API key vacía"
	}

	if len(apiKey) < 8 {
		return false, "API key demasiado corta"
	}

	// Verificar formato básico según proveedor
	switch providerName {
	case "mimo", "minimax", "kimi", "glm":
		// Estos proveedores usan claves alfanuméricas
		if len(apiKey) < 16 {
			return false, fmt.Sprintf("API key %s demasiado corta", providerName)
		}
	}

	return true, "API key aparentemente válida"
}

// GetProviderStats retorna estadísticas de un proveedor
func (hc *HealthChecker) GetProviderStats(providerName string) map[string]interface{} {
	provider, exists := hc.providers[providerName]
	if !exists {
		return map[string]interface{}{
			"error": "Proveedor no encontrado",
		}
	}

	return map[string]interface{}{
		"name":        provider.Name,
		"base_url":    provider.BaseURL,
		"model":       provider.Model,
		"opus_model":  provider.OpusModel,
		"env_var":     provider.EnvVar,
		"description": provider.Description,
	}
}

// RunDiagnostics ejecuta diagnósticos completos del sistema
func (hc *HealthChecker) RunDiagnostics() map[string]interface{} {
	diagnostics := make(map[string]interface{})

	// Verificar todos los proveedores
	results := hc.CheckAll()
	diagnostics["health_results"] = results

	// Calcular métricas
	healthyCount := 0
	var totalLatency time.Duration
	latencyCount := 0

	for _, result := range results {
		if result.Healthy {
			healthyCount++
			if result.Latency > 0 {
				totalLatency += result.Latency
				latencyCount++
			}
		}
	}

	// Métricas agregadas
	diagnostics["summary"] = map[string]interface{}{
		"total_providers": len(results),
		"healthy":         healthyCount,
		"unhealthy":       len(results) - healthyCount,
		"health_rate":     fmt.Sprintf("%.1f%%", float64(healthyCount)/float64(len(results))*100),
	}

	if latencyCount > 0 {
		diagnostics["latency_stats"] = map[string]interface{}{
			"average": totalLatency / time.Duration(latencyCount),
			"total":   totalLatency,
		}
	}

	// Verificar variables de entorno
	envVars := []string{
		"MIMO_API_KEY", "MINIMAX_API_KEY", "KIMI_API_KEY", "GLM_API_KEY",
	}
	envStatus := make(map[string]string)
	for _, env := range envVars {
		if val := os.Getenv(env); val != "" {
			envStatus[env] = "configured"
		} else {
			envStatus[env] = "missing"
		}
	}
	diagnostics["environment"] = envStatus

	return diagnostics
}