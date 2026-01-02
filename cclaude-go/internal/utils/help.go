package utils

import "fmt"

func ShowHelp() {
	helpText := `cclaude - Wrapper multi-proveedor para Claude Code (Go Edition)

Uso: cclaude <proveedor> [flags] [argumentos-claude...]

Proveedores:
  mimo      - Xiaomi MiMo API (requiere MIMO_API_KEY)
  minimax   - MiniMax API (requiere MINIMAX_API_KEY)
  kimi      - Kimi API (requiere KIMI_API_KEY)
  glm       - Zhipu AI API (requiere GLM_API_KEY)
  claude    - Claude nativo (sin configuración adicional)

Flags Básicos:
  -p, --provider string    Proveedor a utilizar
  --timeout string         Timeout de ejecución (default: "5m")
  --debug                  Habilitar modo debug
  --model string           Sobrescribir modelo por defecto
  --config string          Archivo de configuración personalizado
  --help                   Mostrar esta ayuda
  --version                Mostrar versión

Flags de Gestión:
  -i, --interactive        Modo interactivo guiado
  -hc, --health-check      Verificar salud de proveedores
  -d, --diagnose           Diagnóstico completo del sistema
  -sc, --show-config       Mostrar configuración actual
  -c, --confirm            Solicitar confirmación antes de ejecutar
  -pr, --profile string    Usar perfil de configuración
  -lp, --list-profiles     Listar perfiles disponibles
  -cp, --create-profiles   Crear perfiles por defecto

Perfiles de Configuración:
  Los perfiles permiten guardar configuraciones específicas por entorno.
  Se guardan en ~/.config/cclaude/profiles/<nombre>.conf

Ejemplos Básicos:
  cclaude mimo --help
  cclaude minimax "analiza este código"
  cclaude claude --version
  cclaude kimi --model "kimi-k2-thinking-turbo" "mi pregunta"

Ejemplos con Nuevas Funcionalidades:
  cclaude -i                    # Modo interactivo
  cclaude -hc                   # Verificar salud de proveedores
  cclaude -d                    # Diagnóstico completo
  cclaude -sc                   # Ver configuración actual
  cclaude -cp                   # Crear perfiles por defecto
  cclaude -lp                   # Listar perfiles
  cclaude mimo -pr dev          # Usar perfil dev
  cclaude minimax -c "test"     # Confirmar antes de ejecutar

Variables de Entorno:
  MIMO_API_KEY, MINIMAX_API_KEY, KIMI_API_KEY, GLM_API_KEY
  CLAUDE_DEBUG, MAIN_MODEL, ANTHROPIC_BASE_URL

Configuración:
  ~/.cclaude-config para ajustes personalizados
  Formato: KEY=VALUE por línea

Documentación:
  https://github.com/argami/cclaude-go
`
	fmt.Println(helpText)
}

func ShowVersion() {
	fmt.Println("cclaude-go v1.0.0")
	fmt.Println("Wrapper multi-proveedor para Claude Code")
	fmt.Println("Compilado con Go 1.21+")
	fmt.Println("Build: TDD-88.2% coverage")
}