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
  help      - Mostrar esta ayuda

Flags:
  -p, --provider string    Proveedor a utilizar
  --timeout string         Timeout de ejecución (default: "5m")
  --debug                  Habilitar modo debug
  --model string           Sobrescribir modelo por defecto
  --config string          Archivo de configuración personalizado
  --help                   Mostrar esta ayuda
  --version                Mostrar versión

Ejemplos:
  cclaude mimo --help
  cclaude minimax "analiza este código"
  cclaude claude --version
  cclaude kimi --model "kimi-k2-thinking-turbo" "mi pregunta"

Variables de Entorno:
  MIMO_API_KEY, MINIMAX_API_KEY, KIMI_API_KEY, GLM_API_KEY

Configuración:
  ~/.cclaude-config para ajustes personalizados
  Formato: KEY=VALUE por línea
`
	fmt.Println(helpText)
}

func ShowVersion() {
	fmt.Println("cclaude-go v1.0.0")
	fmt.Println("Wrapper multi-proveedor para Claude Code")
	fmt.Println("Compilado con Go 1.21+")
}