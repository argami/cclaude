package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/argami/cclaude-go/internal/provider"
)

// InteractiveMode representa el modo interactivo
type InteractiveMode struct {
	reader    *bufio.Reader
	providers map[string]string
}

// NewInteractiveMode crea un nuevo modo interactivo
func NewInteractiveMode() *InteractiveMode {
	return &InteractiveMode{
		reader: bufio.NewReader(os.Stdin),
		providers: map[string]string{
			"mimo":    "Xiaomi MiMo API",
			"minimax": "MiniMax API",
			"kimi":    "Kimi API",
			"glm":     "Zhipu AI API",
			"claude":  "Claude nativo",
		},
	}
}

// Run ejecuta el modo interactivo
func (im *InteractiveMode) Run() error {
	fmt.Println("🚀 Modo Interactivo cclaude-go")
	fmt.Println("==============================")
	fmt.Println("Este modo te guiará a través de la configuración y ejecución.")
	fmt.Println("Escribe 'salir' o 'exit' en cualquier momento para terminar.")
	fmt.Println()

	for {
		provider, err := im.selectProvider()
		if err != nil {
			return err
		}

		if provider == "" {
			continue // El usuario seleccionó salir
		}

		args, err := im.getArguments()
		if err != nil {
			return err
		}

		if len(args) == 0 && provider != "claude" {
			fmt.Println("⚠️  Advertencia: No se proporcionaron argumentos.")
			if !im.confirm("¿Deseas continuar de todos modos?") {
				continue
			}
		}

		// Confirmar ejecución
		if im.confirm(fmt.Sprintf("¿Ejecutar cclaude con proveedor '%s'?", provider)) {
			return im.execute(provider, args)
		}

		fmt.Println()
	}
}

// selectProvider muestra el selector de proveedores
func (im *InteractiveMode) selectProvider() (string, error) {
	fmt.Println("📋 Proveedores disponibles:")
	i := 1
	for key, desc := range im.providers {
		fmt.Printf("  %d. %-10s - %s\n", i, key, desc)
		i++
	}
	fmt.Println("  0. Salir")

	for {
		fmt.Print("\nSelecciona un proveedor (número o nombre): ")
		input, err := im.readLine()
		if err != nil {
			return "", err
		}

		input = strings.TrimSpace(input)
		if input == "0" || input == "salir" || input == "exit" {
			return "", nil
		}

		// Verificar si es número
		if num, err := im.parseProviderNumber(input); err == nil && num > 0 && num <= len(im.providers) {
			i := 1
			for key := range im.providers {
				if i == num {
					return key, nil
				}
				i++
			}
		}

		// Verificar si es nombre
		if _, exists := im.providers[input]; exists {
			return input, nil
		}

		fmt.Println("❌ Opción inválida. Intenta de nuevo.")
	}
}

// getArguments solicita los argumentos para Claude
func (im *InteractiveMode) getArguments() ([]string, error) {
	fmt.Println("\n📝 Argumentos para Claude:")
	fmt.Println("   Ejemplos: --help, --version, 'analiza este código', etc.")
	fmt.Println("   Deja vacío para ejecución interactiva")
	fmt.Print("Argumentos (separados por espacios): ")

	input, err := im.readLine()
	if err != nil {
		return nil, err
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return []string{}, nil
	}

	// Parsear argumentos
	return strings.Fields(input), nil
}

// confirm solicita confirmación del usuario
func (im *InteractiveMode) confirm(message string) bool {
	for {
		fmt.Printf("%s [s/n]: ", message)
		input, err := im.readLine()
		if err != nil {
			return false
		}

		input = strings.ToLower(strings.TrimSpace(input))
		if input == "s" || input == "si" || input == "y" || input == "yes" {
			return true
		}
		if input == "n" || input == "no" {
			return false
		}

		fmt.Println("❌ Responde 's' o 'n'")
	}
}

// execute ejecuta cclaude con los parámetros seleccionados
func (im *InteractiveMode) execute(provider string, args []string) error {
	fmt.Println("\n🔄 Ejecutando...")
	fmt.Printf("   Proveedor: %s\n", provider)
	if len(args) > 0 {
		fmt.Printf("   Argumentos: %v\n", args)
	} else {
		fmt.Println("   Argumentos: (interactiva)")
	}
	fmt.Println()

	// Preparar argumentos para ExecuteClaude
	fullArgs := []string{provider}
	fullArgs = append(fullArgs, args...)

	// Ejecutar
	if err := ExecuteClaude(fullArgs); err != nil {
		fmt.Printf("❌ Error durante la ejecución: %v\n", err)
		if im.confirm("¿Deseas intentar de nuevo?") {
			return nil // Continuar el bucle
		}
		return err
	}

	fmt.Println("✅ Ejecución completada exitosamente")
	return nil
}

// readLine lee una línea de entrada
func (im *InteractiveMode) readLine() (string, error) {
	input, err := im.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// parseProviderNumber convierte un string a número de proveedor
func (im *InteractiveMode) parseProviderNumber(input string) (int, error) {
	var num int
	_, err := fmt.Sscanf(input, "%d", &num)
	if err != nil {
		return 0, err
	}
	return num, nil
}

// ShowConfig muestra la configuración actual
func ShowConfig() {
	// Variables de entorno
	envVars := []string{
		"MIMO_API_KEY", "MINIMAX_API_KEY", "KIMI_API_KEY", "GLM_API_KEY",
		"CLAUDE_DEBUG", "MAIN_MODEL", "ANTHROPIC_BASE_URL",
	}

	for _, env := range envVars {
		val := os.Getenv(env)
		if val != "" {
			if strings.Contains(env, "KEY") {
				// Ocultar API keys
				fmt.Printf("  ✅ %s: %s... (configurado)\n", env, val[:8])
			} else {
				fmt.Printf("  ✅ %s: %s\n", env, val)
			}
		} else {
			fmt.Printf("  ❌ %s: (no configurado)\n", env)
		}
	}

	// Verificar binario claude
	fmt.Println("\n🔍 Verificación de sistema:")
	if path, err := exec.LookPath("claude"); err == nil {
		fmt.Printf("  ✅ claude encontrado: %s\n", path)
	} else {
		fmt.Println("  ❌ claude no encontrado en PATH")
	}
}

// ShowTips muestra consejos de uso
func ShowTips() {
	tips := []string{
		"💡 Consejos:",
		"  • Usa 'cclaude --help' para ver todas las opciones",
		"  • Configura variables de entorno en ~/.cclaude-config",
		"  • Usa perfiles para diferentes entornos (dev/prod/test)",
		"  • Habilita modo debug con --debug para troubleshooting",
		"  • Verifica salud de proveedores con --health-check",
		"  • Usa modo interactivo para exploración guiada",
		"",
		"🔑 Variables de entorno esenciales:",
		"  • MIMO_API_KEY, MINIMAX_API_KEY, KIMI_API_KEY, GLM_API_KEY",
		"",
		"🚀 Ejemplos rápidos:",
		"  • cclaude mimo --help",
		"  • cclaude minimax 'analiza este código'",
		"  • cclaude claude --version",
		"  • cclaude kimi --model 'kimi-k2-thinking-turbo' 'mi pregunta'",
	}

	for _, tip := range tips {
		fmt.Println(tip)
	}
}

// InteractiveHealthCheck ejecuta health checks interactivos
func InteractiveHealthCheck() error {
	fmt.Println("🔍 Health Check Interactivo")
	fmt.Println("===========================")

	im := NewInteractiveMode()
	healthChecker := provider.NewHealthChecker()

	// Mostrar opciones
	fmt.Println("\nOpciones:")
	fmt.Println("  1. Verificar todos los proveedores")
	fmt.Println("  2. Verificar proveedor específico")
	fmt.Println("  3. Verificar API key")
	fmt.Println("  4. Diagnóstico completo")
	fmt.Println("  0. Salir")

	for {
		fmt.Print("\nSelecciona opción: ")
		input, err := im.readLine()
		if err != nil {
			return err
		}

		input = strings.TrimSpace(input)

		switch input {
		case "0", "salir", "exit":
			return nil

		case "1":
			results := healthChecker.CheckAll()
			fmt.Println(healthChecker.FormatHealthResults(results, true))

		case "2":
			providerName, err := im.selectProvider()
			if err != nil {
				return err
			}
			if providerName == "" {
				continue
			}
			result := healthChecker.CheckProvider(providerName)
			fmt.Println(healthChecker.FormatHealthResults([]provider.HealthCheckResult{result}, true))

		case "3":
			provider, err := im.selectProvider()
			if err != nil {
				return err
			}
			if provider == "" {
				continue
			}
			fmt.Print("Introduce tu API key: ")
			apiKey, _ := im.readLine()
			valid, msg := healthChecker.VerifyAPIKey(provider, apiKey)
			if valid {
				fmt.Printf("✅ %s\n", msg)
			} else {
				fmt.Printf("❌ %s\n", msg)
			}

		case "4":
			diagnostics := healthChecker.RunDiagnostics()
			fmt.Println("\n📊 Diagnóstico Completo:")
			for k, v := range diagnostics {
				fmt.Printf("  %s: %v\n", k, v)
			}

		default:
			fmt.Println("❌ Opción inválida")
		}
	}
}