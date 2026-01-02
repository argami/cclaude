# Plan de Implementación: cclaude en Go

**Fecha de creación:** 2026-01-02
**Autor:** Análisis automatizado
**Tiempo de análisis y diseño:** 25 minutos
**Tiempo estimado de implementación:** 2-3 semanas
**Lenguaje objetivo:** Go (Golang)

---

## 📋 Resumen Ejecutivo

Migración del wrapper bash `cclaude` a una aplicación Go nativa, manteniendo la funcionalidad actual pero añadiendo robustez, seguridad y portabilidad multiplataforma.

**Estado actual:** ✅ Bash básico funcional
**Objetivo:** 🔧 Go nativo con binarios auto-contenidos
**Ventaja principal:** Distribución simple + portabilidad total (incluye Windows)

---

## 🎯 Arquitectura Go Propuesta

### Estructura del Proyecto
```
cclaude-go/
├── cmd/
│   └── cclaude/
│       └── main.go              # Punto de entrada
├── internal/
│   ├── provider/                # Lógica de proveedores
│   │   ├── provider.go
│   │   ├── mimo.go
│   │   ├── minimax.go
│   │   ├── kimi.go
│   │   ├── glm.go
│   │   └── claude.go
│   ├── config/                  # Gestión de configuración
│   │   ├── config.go
│   │   └── validation.go
│   ├── flags/                   # Parsing de argumentos
│   │   └── flags.go
│   └── utils/                   # Utilidades
│       ├── env.go
│       ├── errors.go
│       └── logging.go
├── pkg/                         # Paquetes reutilizables
│   └── types/
│       └── types.go
├── go.mod
├── go.sum
├── README.md
├── Makefile
└── .goreleaser.yml             # Para builds automatizados
```

### Componentes Principales

#### 1. **Estructura de Datos**
```go
type Provider struct {
    Name         string
    BaseURL      string
    Model        string
    DefaultOpus  string
    EnvKey       string
    AuthToken    string
}

type Config struct {
    Provider      Provider
    Timeout       time.Duration
    Debug         bool
    ConfigFile    string
    ModelOverride string
}
```

#### 2. **Flujo de Ejecución**
```
1. Parsear argumentos → flags.Parse()
2. Validar ambiente → config.Validate()
3. Cargar configuración → config.Load()
4. Seleccionar proveedor → provider.Get()
5. Configurar variables → env.Setup()
6. Ejecutar claude → exec.Command()
7. Manejar errores → errors.Handle()
```

---

## 🚀 Plan de Implementación Detallado

### Fase 1: Fundamentos (Día 1-3) - CRÍTICO

#### 1.1 Estructura Base y Módulos
```bash
# Inicializar proyecto Go
mkdir cclaude-go && cd cclaude-go
go mod init github.com/argami/cclaude-go

# Estructura de directorios
mkdir -p cmd/cclaude internal/{provider,config,flags,utils} pkg/types
```

#### 1.2 Tipos y Estructuras de Datos
```go
// internal/types/types.go
package types

import "time"

type ProviderConfig struct {
    Name        string
    BaseURL     string
    Model       string
    OpusModel   string
    EnvVar      string
    Description string
}

type AppConfig struct {
    Provider      *ProviderConfig
    Timeout       time.Duration
    Debug         bool
    ModelOverride string
    Args          []string
}
```

#### 1.3 Sistema de Proveedores
```go
// internal/provider/provider.go
package provider

var Providers = map[string]ProviderConfig{
    "mimo": {
        Name:       "mimo",
        BaseURL:    "https://api.xiaomimimo.com/anthropic",
        Model:      "mimo-v2-flash",
        OpusModel:  "mimo-v2-flash",
        EnvVar:     "MIMO_API_KEY",
        Description: "Xiaomi MiMo API",
    },
    // ... otros proveedores
}

func GetProvider(name string) (*ProviderConfig, error) {
    if provider, exists := Providers[name]; exists {
        return &provider, nil
    }
    return nil, fmt.Errorf("proveedor no encontrado: %s", name)
}
```

### Fase 2: Lógica de Configuración (Día 4-6)

#### 2.1 Manejo de Configuración
```go
// internal/config/config.go
package config

import (
    "os"
    "path/filepath"
    "time"
)

const (
    DefaultTimeout = 5 * time.Minute
    ConfigFileName = ".cclaude-config"
)

type ConfigLoader struct {
    ConfigPath string
}

func (cl *ConfigLoader) Load() (*AppConfig, error) {
    // Cargar desde archivo de configuración
    // Cargar desde variables de entorno
    // Combinar con flags de CLI
}
```

#### 2.2 Validación de Ambiente
```go
// internal/config/validation.go
package config

import (
    "os/exec"
    "strings"
)

func ValidateEnvironment() error {
    // Verificar que 'claude' está disponible
    if _, err := exec.LookPath("claude"); err != nil {
        return fmt.Errorf("comando 'claude' no encontrado en PATH")
    }
    return nil
}

func ValidateAPIKey(provider ProviderConfig) error {
    key := os.Getenv(provider.EnvVar)
    if key == "" {
        return fmt.Errorf("variable de entorno %s no configurada", provider.EnvVar)
    }
    if len(key) < 8 {
        return fmt.Errorf("API key inusualmente corta para %s", provider.Name)
    }
    return nil
}
```

#### 2.3 Parsing de Argumentos
```go
// internal/flags/flags.go
package flags

import (
    "flag"
    "fmt"
    "os"
)

type FlagConfig struct {
    Provider      string
    Timeout       string
    Debug         bool
    Help          bool
    Version       bool
    ModelOverride string
    ConfigFile    string
}

func Parse() (*FlagConfig, error) {
    var flags FlagConfig

    flag.StringVar(&flags.Provider, "provider", "", "Proveedor de API (mimo, minimax, kimi, glm, claude)")
    flag.StringVar(&flags.Provider, "p", "", "Abreviatura para --provider")
    flag.StringVar(&flags.Timeout, "timeout", "5m", "Timeout para la ejecución")
    flag.BoolVar(&flags.Debug, "debug", false, "Modo debug")
    flag.BoolVar(&flags.Help, "help", false, "Mostrar ayuda")
    flag.BoolVar(&flags.Version, "version", false, "Mostrar versión")
    flag.StringVar(&flags.ModelOverride, "model", "", "Sobrescribir modelo por defecto")
    flag.StringVar(&flags.ConfigFile, "config", "", "Archivo de configuración personalizado")

    flag.Parse()

    // Si no hay proveedor y no son flags de ayuda, usar el primer argumento
    if flags.Provider == "" && flag.NArg() > 0 {
        flags.Provider = flag.Arg(0)
    }

    return &flags, nil
}
```

### Fase 3: Ejecución y Manejo de Errores (Día 7-10)

#### 3.1 Configuración de Variables de Entorno
```go
// internal/utils/env.go
package utils

import (
    "os"
    "fmt"
)

func SetupEnvironment(provider ProviderConfig, authToken string, modelOverride string) error {
    // Limpiar variables anteriores
    os.Unsetenv("ANTHROPIC_BASE_URL")
    os.Unsetenv("MAIN_MODEL")
    os.Unsetenv("ANTHROPIC_AUTH_TOKEN")

    // Configurar nuevas variables
    if err := os.Setenv("ANTHROPIC_BASE_URL", provider.BaseURL); err != nil {
        return fmt.Errorf("error configurando ANTHROPIC_BASE_URL: %w", err)
    }

    model := provider.Model
    if modelOverride != "" {
        model = modelOverride
    }

    os.Setenv("MAIN_MODEL", model)
    os.Setenv("ANTHROPIC_AUTH_TOKEN", authToken)
    os.Setenv("ANTHROPIC_DEFAULT_OPUS_MODEL", provider.OpusModel)
    os.Setenv("ANTHROPIC_MODEL", model)
    os.Setenv("ANTHROPIC_DEFAULT_SONNET_MODEL", model)
    os.Setenv("ANTHROPIC_DEFAULT_HAIKU_MODEL", model)
    os.Setenv("CLAUDE_CODE_SUBAGENT_MODEL", model)
    os.Setenv("DISABLE_NON_ESSENTIAL_MODEL_CALLS", "1")
    os.Setenv("CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC", "1")
    os.Setenv("API_TIMEOUT_MS", "3000000")

    return nil
}
```

#### 3.2 Ejecución de Claude
```go
// internal/utils/exec.go
package utils

import (
    "os"
    "os/exec"
    "syscall"
)

func ExecuteClaude(args []string) error {
    cmd := exec.Command("claude", args...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    // Propagar señales (Ctrl+C)
    cmd.SysProcAttr = &syscall.SysProcAttr{
        // Configuración específica por plataforma
    }

    return cmd.Run()
}
```

#### 3.3 Sistema de Errores
```go
// internal/utils/errors.go
package utils

import (
    "fmt"
    "os"
)

type ExitCode int

const (
    ExitSuccess ExitCode = iota
    ExitProviderNotFound
    ExitAPIKeyMissing
    ExitClaudeNotFound
    ExitConfigError
    ExitValidationError
)

func HandleError(err error, code ExitCode) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
        os.Exit(int(code))
    }
}

func HandleErrorf(format string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, "❌ Error: %s\n", fmt.Sprintf(format, args...))
    os.Exit(1)
}
```

### Fase 4: Logging y Ayuda (Día 11-13)

#### 4.1 Sistema de Logging
```go
// internal/utils/logging.go
package utils

import (
    "fmt"
    "os"
    "time"
)

type LogLevel int

const (
    LevelSilent LogLevel = iota
    LevelError
    LevelWarn
    LevelInfo
    LevelDebug
)

var currentLevel = LevelInfo

func SetLogLevel(level LogLevel) {
    currentLevel = level
}

func log(level LogLevel, format string, args ...interface{}) {
    if level > currentLevel {
        return
    }

    timestamp := time.Now().Format("2006-01-02 15:04:05")
    prefix := ""

    switch level {
    case LevelError:
        prefix = "❌"
    case LevelWarn:
        prefix = "⚠️"
    case LevelInfo:
        prefix = "ℹ️"
    case LevelDebug:
        prefix = "🔍"
    }

    message := fmt.Sprintf(format, args...)
    fmt.Printf("[%s] %s %s\n", timestamp, prefix, message)
}

func Info(format string, args ...interface{})  { log(LevelInfo, format, args...) }
func Warn(format string, args ...interface{})  { log(LevelWarn, format, args...) }
func Error(format string, args ...interface{}) { log(LevelError, format, args...) }
func Debug(format string, args ...interface{}) { log(LevelDebug, format, args...) }
```

#### 4.2 Sistema de Ayuda
```go
// internal/utils/help.go
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
```

### Fase 5: Main y Orquestación (Día 14-16)

#### 5.1 Punto de Entrada
```go
// cmd/cclaude/main.go
package main

import (
    "os"
    "time"

    "github.com/argami/cclaude-go/internal/config"
    "github.com/argami/cclaude-go/internal/flags"
    "github.com/argami/cclaude-go/internal/provider"
    "github.com/argami/cclaude-go/internal/utils"
)

func main() {
    // Parsear flags
    flagConfig, err := flags.Parse()
    if err != nil {
        utils.HandleError(err, utils.ExitConfigError)
    }

    // Manejar flags de ayuda y versión
    if flagConfig.Help {
        utils.ShowHelp()
        os.Exit(0)
    }

    if flagConfig.Version {
        utils.ShowVersion()
        os.Exit(0)
    }

    // Validar ambiente
    if err := config.ValidateEnvironment(); err != nil {
        utils.HandleError(err, utils.ExitValidationError)
    }

    // Obtener proveedor
    providerConfig, err := provider.GetProvider(flagConfig.Provider)
    if err != nil {
        // Si no hay proveedor o es "claude", ejecutar nativo
        if flagConfig.Provider == "" || flagConfig.Provider == "claude" {
            utils.Info("Ejecutando Claude nativo")
            args := flagConfig.Args
            if len(args) == 0 {
                args = os.Args[1:]
            }
            if err := utils.ExecuteClaude(args); err != nil {
                utils.HandleError(err, utils.ExitClaudeNotFound)
            }
            return
        }
        utils.HandleError(err, utils.ExitProviderNotFound)
    }

    // Validar API key
    if err := config.ValidateAPIKey(*providerConfig); err != nil {
        utils.HandleError(err, utils.ExitAPIKeyMissing)
    }

    // Configurar timeout
    timeout, err := time.ParseDuration(flagConfig.Timeout)
    if err != nil {
        timeout = 5 * time.Minute
    }

    // Configurar variables de entorno
    authToken := os.Getenv(providerConfig.EnvVar)
    if err := utils.SetupEnvironment(*providerConfig, authToken, flagConfig.ModelOverride); err != nil {
        utils.HandleError(err, utils.ExitConfigError)
    }

    // Logging de configuración
    utils.Info("Proveedor: %s", providerConfig.Name)
    utils.Info("Modelo: %s", providerConfig.Model)
    utils.Info("Timeout: %s", timeout)

    if flagConfig.Debug {
        utils.SetLogLevel(utils.LevelDebug)
        utils.Debug("Modo debug habilitado")
        utils.Debug("Base URL: %s", providerConfig.BaseURL)
    }

    // Ejecutar claude con argumentos restantes
    claudeArgs := flagConfig.Args
    if len(claudeArgs) == 0 {
        claudeArgs = flagConfig.Args
    }

    if err := utils.ExecuteClaude(claudeArgs); err != nil {
        utils.HandleError(err, utils.ExitClaudeNotFound)
    }
}
```

### Fase 6: Build y Distribución (Día 17-18)

#### 6.1 Makefile
```makefile
# Makefile
BINARY_NAME=cclaude
VERSION=1.0.0
BUILD_DIR=build

.PHONY: build build-linux build-macos build-windows clean test install

build: build-linux build-macos build-windows

build-linux:
	@echo "Building Linux amd64..."
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/cclaude
	@chmod +x $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64

build-macos:
	@echo "Building macOS amd64..."
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-macos-amd64 ./cmd/cclaude
	@chmod +x $(BUILD_DIR)/$(BINARY_NAME)-macos-amd64
	@echo "Building macOS arm64..."
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-macos-arm64 ./cmd/cclaude
	@chmod +x $(BUILD_DIR)/$(BINARY_NAME)-macos-arm64

build-windows:
	@echo "Building Windows amd64..."
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/cclaude

clean:
	@echo "Cleaning builds..."
	@rm -rf $(BUILD_DIR)

test:
	@echo "Running tests..."
	@go test ./...

install:
	@echo "Installing to /usr/local/bin..."
	@go build -o /usr/local/bin/$(BINARY_NAME) ./cmd/cclaude

# Cross-compile all platforms
cross-compile:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/cclaude
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-macos-amd64 ./cmd/cclaude
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-macos-arm64 ./cmd/cclaude
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/cclaude
```

#### 6.2 GoReleaser Config
```yaml
# .goreleaser.yml
project_name: cclaude

builds:
  - id: cclaude
    binary: cclaude
    main: ./cmd/cclaude
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    env:
      - CGO_ENABLED=0
    flags:
      - -trimpath
    ldflags:
      - -s -w -X main.version={{.Version}} -X main.commit={{.Commit}} -X main.date={{.Date}}

archives:
  - id: cclaude
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    format_overrides:
      - goos: windows
        format: zip
    files:
      - LICENSE
      - README.md

checksum:
  name_template: 'checksums.txt'

changelog:
  sort: asc
  filters:
    exclude:
      - '^docs:'
      - '^test:'
```

### Fase 7: Tests y Validación (Día 19-21)

#### 7.1 Tests Unitarios
```go
// internal/provider/provider_test.go
package provider

import "testing"

func TestGetProvider(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        expectError bool
    }{
        {"mimo", "mimo", false},
        {"minimax", "minimax", false},
        {"invalid", "invalid", true},
        {"empty", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            provider, err := GetProvider(tt.input)
            if tt.expectError {
                if err == nil {
                    t.Errorf("Expected error for %s, got nil", tt.input)
                }
            } else {
                if err != nil {
                    t.Errorf("Unexpected error for %s: %v", tt.input, err)
                }
                if provider == nil {
                    t.Errorf("Expected provider for %s, got nil", tt.input)
                }
            }
        })
    }
}
```

#### 7.2 Tests de Integración
```go
// internal/config/validation_test.go
package config

import (
    "os"
    "testing"
)

func TestValidateAPIKey(t *testing.T) {
    // Test con key válida
    os.Setenv("MIMO_API_KEY", "test-key-12345678")
    provider := Providers["mimo"]
    err := ValidateAPIKey(provider)
    if err != nil {
        t.Errorf("Expected no error with valid key, got: %v", err)
    }

    // Test con key faltante
    os.Unsetenv("MIMO_API_KEY")
    err = ValidateAPIKey(provider)
    if err == nil {
        t.Error("Expected error with missing key")
    }

    // Test con key corta
    os.Setenv("MIMO_API_KEY", "short")
    err = ValidateAPIKey(provider)
    if err == nil {
        t.Error("Expected error with short key")
    }
}
```

### Fase 8: Documentación y Ejemplos (Día 22-23)

#### 8.1 README.md
```markdown
# cclaude-go

Wrapper multi-proveedor para Claude Code escrito en Go, con portabilidad nativa y robustez mejorada.

## Características

- ✅ **Multi-plataforma**: Binarios para Linux, macOS (Intel/Apple Silicon), Windows
- ✅ **Sin dependencias**: Un solo archivo binario auto-contenido
- ✅ **Validación robusta**: Chequeos de ambiente y API keys
- ✅ **Configuración flexible**: Archivos de config + variables de entorno + flags
- ✅ **Logging estructurado**: Niveles de debug, info, warning, error
- ✅ **Manejo de errores**: Códigos de salida específicos
- ✅ **Timeout configurable**: Prevención de ejecuciones colgadas

## Instalación

### Desde binarios pre-compilados
```bash
# Linux
curl -L https://github.com/argami/cclaude-go/releases/latest/download/cclaude-linux-amd64 -o cclaude
chmod +x cclaude
sudo mv cclaude /usr/local/bin/

# macOS
curl -L https://github.com/argami/cclaude-go/releases/latest/download/cclaude-macos-arm64 -o cclaude
chmod +x cclaude
sudo mv cclaude /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/argami/cclaude-go/releases/latest/download/cclaude-windows-amd64.exe" -OutFile "cclaude.exe"
Move-Item -Path "cclaude.exe" -Destination "$env:USERPROFILE\AppData\Local\Microsoft\WindowsApps\cclaude.exe"
```

### Compilando desde fuente
```bash
git clone https://github.com/argami/cclaude-go.git
cd cclaude-go
make build
sudo make install
```

## Uso

```bash
# Proveedores alternativos
cclaude mimo "tu pregunta aquí"
cclaude minimax --help
cclaude kimi --model "kimi-k2-thinking-turbo" "analiza esto"

# Claude nativo
cclaude claude --version
cclaude --help

# Modo debug
cclaude mimo --debug "test query"

# Timeout personalizado
cclaude minimax --timeout 10m "tarea larga"
```

## Configuración

### Variables de entorno
```bash
export MIMO_API_KEY="tu-key-aqui"
export MINIMAX_API_KEY="tu-key-aqui"
export KIMI_API_KEY="tu-key-aqui"
export GLM_API_KEY="tu-key-aqui"
```

### Archivo de configuración (~/.cclaude-config)
```
MIMO_API_KEY=your-mimo-key
MINIMAX_API_KEY=your-minimax-key
KIMI_API_KEY=your-kimi-key
GLM_API_KEY=your-glm-key
TIMEOUT=10m
DEBUG=false
```

## Desarrollo

```bash
# Estructura del proyecto
cclaude-go/
├── cmd/cclaude/          # Punto de entrada
├── internal/             # Paquetes internos
│   ├── provider/         # Lógica de proveedores
│   ├── config/           # Configuración y validación
│   ├── flags/            # Parsing de argumentos
│   └── utils/            # Utilidades
├── pkg/                  # Paquetes públicos
├── go.mod
└── Makefile

# Build y test
make build        # Compilar para todas las plataformas
make test         # Ejecutar tests
make install      # Instalar localmente
```

## Migración desde bash

Si tienes la versión bash instalada:
```bash
# Backup del original
sudo cp /usr/local/bin/cclaude /usr/local/bin/cclaude-bash

# Instalar versión Go
sudo make install

# Verificar
cclaude --version
```

## Licencia

MIT
```

### Fase 9: GitHub Actions y CI/CD (Día 24-25)

#### 9.1 GitHub Actions Workflow
```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Run tests
        run: make test

      - name: Build binaries
        run: make cross-compile

      - name: Create Release
        uses: goreleaser/goreleaser-action@v5
        with:
          distribution: goreleaser
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

## 📊 Métricas de Éxito

### Rendimiento
- **Tiempo de inicio**: < 50ms (vs 100ms+ bash)
- **Uso de memoria**: < 10MB adicional
- **Overhead**: < 5% vs bash original

### Portabilidad
- **Binarios**: 4 plataformas × 2 arquitecturas = 8 builds
- **Tamaño binario**: ~5-10MB por build
- **Sin dependencias externas**: 100% auto-contenido

### Robustez
- **Validación**: 100% de inputs validados
- **Manejo de errores**: Códigos específicos por tipo
- **Logging**: 4 niveles de detalle
- **Tests**: >80% cobertura

### Seguridad
- **Exposición de keys**: 0 (nunca en logs)
- **Sanitización de inputs**: 100%
- **Timeout configurable**: Prevención de bloqueos

---

## 🎯 Roadmap de Implementación

### Semana 1: Fundamentos
- [ ] Estructura de proyecto y go.mod
- [ ] Tipos y estructuras de datos
- [ ] Sistema de proveedores
- [ ] Parsing básico de flags

### Semana 2: Lógica Principal
- [ ] Validación de ambiente y API keys
- [ ] Configuración de variables de entorno
- [ ] Ejecución de Claude con timeout
- [ ] Manejo de errores y logging

### Semana 3: Tests y Distribución
- [ ] Tests unitarios y de integración
- [ ] Makefile y builds multi-plataforma
- [ ] Documentación completa
- [ ] GitHub Actions para CI/CD

### Semana 4: Polish y Release
- [ ] Tests de integración en plataformas reales
- [ ] Documentación de migración
- [ ] Release de binarios
- [ ] Validación final de portabilidad

---

## 📝 Decisiones de Diseño Clave

### 1. **Por qué Go sobre otros lenguajes**
- **Portabilidad nativa**: Compilación cruzada sin toolchains externas
- **Performance**: Arranque rápido, bajo consumo de recursos
- **Distribución**: Un solo binario sin dependencias
- **Seguridad**: Tipos fuertes, manejo explícito de errores
- **Ecosistema**: Herramientas maduras (go mod, testing, goreleaser)

### 2. **Arquitectura por capas**
- **internal/**: Paquetes privados, no reutilizables externamente
- **pkg/**: Paquetes públicos (si se necesitan en el futuro)
- **cmd/****: Puntos de entrada claros
- **Separación de responsabilidades**: Config, provider, execution, utils

### 3. **Manejo de errores**
- **Códigos de salida específicos**: Facilita debugging y scripting
- **Errores envueltos**: Contexto completo con `fmt.Errorf("...: %w", err)`
- **Fail fast**: Validación temprana, sin ejecuciones parciales

### 4. **Configuración flexible**
- **Jerarquía**: Flags > Env Vars > Config File > Defaults
- **Backward compatibility**: Mantiene misma interfaz que bash original
- **Extensibilidad**: Fácil añadir nuevos proveedores o configuraciones

---

## ⏱️ Tiempo de Planificación

**Análisis de requerimientos:** 10 minutos
**Diseño de arquitectura:** 8 minutos
**Planificación detallada:** 7 minutos
**Total:** 25 minutos

**Fecha de finalización:** 2026-01-02 15:00 UTC

---

## 📚 Recursos y Referencias

### Documentación Go
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Standard Library](https://pkg.go.dev/std)

### Herramientas
- **GoReleaser**: Distribución automatizada
- **Make**: Build automation
- **GitHub Actions**: CI/CD
- **Go Test**: Testing framework

### Patrones de Diseño
- **Clean Architecture**: Separación de responsabilidades
- **Dependency Injection**: Facilita testing
- **Error Wrapping**: Contexto en errores
- **Context Pattern**: Timeout y cancelación

---

## ✅ Próximos Pasos

1. **Aprobación del plan**: Revisar y confirmar arquitectura
2. **Setup inicial**: `go mod init` y estructura de directorios
3. **Implementación incremental**: Fase por fase según roadmap
4. **Testing temprano**: Tests unitarios en cada fase
5. **Iteración**: Feedback loops con builds funcionales

---

**Estado del Plan:** ✅ **COMPLETO** - Listo para implementación

**Decisión final:** Go es la elección óptima para portabilidad + robustez + distribución simple.