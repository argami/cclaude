# cclaude - Claude Code Multi-Provider Wrapper

**Language:** Go | **Portability:** Binario único, zero dependencies

---

## Tabla de Contenidos

- [Descripción](#descripción)
- [Instalación](#instalación)
- [Uso](#uso)
- [Proveedores](#proveedores)
- [Configuración](#configuración)
- [Desarrollo](#desarrollo)
- [Distribución](#distribución)

---

## Descripción

`cclaude` es un wrapper de Go para ejecutar Claude Code con múltiples proveedores de API alternativos (mimo, minimax, kimi, glm).

### Características

- ✅ **Portabilidad extrema**: Binario único, zero dependencies
- ✅ **Type safety**: Go typed para menos bugs en runtime
- ✅ **Testing robusto**: Unit e integration tests incluidos
- ✅ **Cross-platform**: Linux y macOS (amd64 + arm64)
- ✅ **Configuración flexible**: YAML + environment variables + CLI flags

### Arquitectura

```
cclaude/
├── cmd/
│   └── cclaude/main.go        # Entry point
├── internal/
│   ├── config/config.go       # Configuración
│   ├── provider/provider.go   # Lógica de proveedores
│   └── validator/validator.go # Validación de API keys
├── pkg/
│   ├── models/provider.go     # Structs de proveedores
│   └── cli/flags.go           # Parsing de flags
├── config/cclaude.yaml        # Ejemplo de configuración
├── Makefile                   # Build y distribución
└── README.md                  # Este archivo
```

---

## Instalación

### Opción 1: Desde el código fuente

```bash
git clone https://github.com/argami/cclaude.git
cd cclaude
make install
```

### Opción 2: Binario pre-compilado

```bash
# Linux amd64
curl -L https://github.com/argami/cclaude/releases/download/v1.0.0/cclaude-linux-amd64 -o cclaude
chmod +x cclaude
sudo mv cclaude /usr/local/bin/

# macOS arm64
curl -L https://github.com/argami/cclaude/releases/download/v1.0.0/cclaude-darwin-arm64 -o cclaude
chmod +x cclaude
sudo mv cclaude /usr/local/bin/
```

### Opción 3: Build desde Makefile

```bash
make build      # Compila el binario ./cclaude
make install    # Instala en /usr/local/bin/cclaude
```

---

## Uso

### Sintaxis básica

```bash
cclaude [FLAGS] [ARGUMENTOS PARA CLAUDE]
```

### Flags disponibles

| Flag | Descripción |
|------|-------------|
| `--provider, -p` | Proveedor a usar (mimo, minimax, kimi, glm) |
| `--model, -m` | Modelo específico para el proveedor |
| `--dry-run, -n` | Muestra configuración sin ejecutar |
| `--verbose, -v` | Output detallado de debugging |
| `--version` | Muestra versión |
| `--help` | Muestra ayuda |

### Ejemplos de uso

#### Uso básico

```bash
# Proveedor por defecto (Claude nativo)
cclaude "Write a function to calculate fibonacci"

# Con proveedor específico
cclaude minimax "Analyze this code"
cclaude mimo --help
```

#### Con flags

```bash
# Dry-run: ver configuración sin ejecutar
cclaude --dry-run --provider minimax "Hello"

# Verbose: output detallado
cclaude --verbose --provider kimi "Complex task"

# Modelo específico
cclaude --provider minimax --model MiniMax-M2.1 "Task"
```

#### Con argumentos para Claude

```bash
# Arguments passed through to claude
cclaude mimo --print --dangerously-skip-permissions-check "Help me refactor this code"

# Multi-word prompts
cclaude glm "Create a REST API with Go and Gin framework"
```

---

## Proveedores

| Proveedor | API Key Env | URL | Modelo Default |
|-----------|-------------|-----|----------------|
| **Claude (native)** | `ANTHROPIC_API_KEY` | Native | - |
| **mimo** | `MIMO_API_KEY` | `https://api.xiaomimimo.com/anthropic` | `mimo-v2-flash` |
| **minimax** | `MINIMAX_API_KEY` | `https://api.minimax.io/anthropic` | `MiniMax-M2.1` |
| **kimi** | `KIMI_API_KEY` | `https://api.kimi.com/coding/` | `kimi-k2-0711-preview` |
| **glm** | `GLM_API_KEY` | `https://api.z.ai/api/anthropic` | `glm-4.7` |

### Configuración de API Keys

```bash
# Opción 1: Variables de entorno
export MIMO_API_KEY="tu-api-key-aqui"
export MINIMAX_API_KEY="tu-api-key-aqui"
export KIMI_API_KEY="tu-api-key-aqui"
export GLM_API_KEY="tu-api-key-aqui"

# Opción 2: En config.yaml
# Ver sección de configuración

# Opción 3: En línea (no recomendado)
cclaude MINIMAX_API_KEY=xxx minimax "Task"
```

---

## Configuración

### Archivo de configuración (config/cclaude.yaml)

```yaml
# cclaude Configuration File
# Located at: ~/.config/cclaude/cclaude.yaml or ./config/cclaude.yaml

# Proveedor por defecto
provider: minimax

# Modelo por defecto (opcional, usa el del proveedor si no se especifica)
model: MiniMax-M2.1

# Timeout en milisegundos (3000000 = 50 minutos)
timeout_ms: 3000000

# Variables de entorno personalizadas para cada proveedor
env_overrides:
  minimax:
    DISABLE_NON_ESSENTIAL_MODEL_CALLS: "1"
  mimo:
    CUSTOM_VAR: "custom_value"

# Configuración global
global:
  verbose: false
  dry_run: false
```

### Precedencia de configuración

```
1. CLI flags (máxima prioridad)
2. Environment variables
3. Archivo de configuración
4. Valores por defecto
```

### Ejemplo de configuración completa

```yaml
provider: minimax
model: MiniMax-M2.1
timeout_ms: 3000000

env_overrides:
  all:
    DISABLE_NON_ESSENTIAL_MODEL_CALLS: "1"
    CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC: "1"
  minimax:
    ANTHROPIC_DEFAULT_OPUS_MODEL: "MiniMax-M2.1"
  mimo:
    ANTHROPIC_DEFAULT_OPUS_MODEL: "mimo-v2-flash"
```

---

## Desarrollo

### Requisitos

- Go 1.21+
- make

### Setup del proyecto

```bash
# Clonar el repositorio
git clone https://github.com/argami/cclaude.git
cd cclaude

# Inicializar módulo (si es necesario)
go mod init github.com/argami/cclaude

# Descargar dependencias
go mod tidy

# Compilar
make build

# Ejecutar tests
make test
```

### Comandos de desarrollo

```bash
make build        # Compila el binario
make test         # Ejecuta todos los tests
make clean        # Limpia binarios
make release      # Compila para todas las plataformas
make install      # Instala en /usr/local/bin
```

### Estructura del proyecto

```
cclaude/
├── cmd/
│   └── cclaude/
│       └── main.go              # Entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Carga y merge de configuración
│   ├── provider/
│   │   └── provider.go          # Lógica de ejecución
│   └── validator/
│       └── validator.go         # Validación de API keys
├── pkg/
│   ├── cli/
│   │   └── flags.go             # Parsing de flags
│   └── models/
│       └── provider.go          # Structs de providers
├── config/
│   └── cclaude.yaml             # Ejemplo de configuración
├── Makefile                     # Build automation
├── go.mod                       # Dependencias
└── README.md                    # Este archivo
```

### Agregar un nuevo proveedor

1. Editar `pkg/models/provider.go`:

```go
var providers = map[string]Provider{
    // Proveedores existentes...
    "nuevo-proveedor": {
        Name:    "nuevo-proveedor",
        BaseURL: "https://api.nuevo.com/anthropic",
        Model:   "nuevo-modelo",
        EnvKey:  "NUEVO_PROVEEDOR_API_KEY",
    },
}
```

2. Agregar tests en `pkg/models/provider_test.go`

3. Documentar en este README

### Testing

```bash
# Todos los tests con coverage
go test ./... -v -cover

# Tests específicos
go test ./pkg/models/ -v
go test ./internal/validator/ -v

# Verificar coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

---

## Distribución

### Build para todas las plataformas

```bash
make release
```

Esto genera:
- `cclaude-linux-amd64`
- `cclaude-darwin-amd64`
- `cclaude-darwin-arm64`

### Instalación manual

```bash
# Compilar
make build

# Instalar
sudo cp cclaude /usr/local/bin/
sudo chmod +x /usr/local/bin/cclaude

# Verificar
cclaude --version
```

### Crear release en GitHub

```bash
# Tag de versión
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# Crear binarios
make release

# Subir binarios a GitHub Releases
```

---

## Troubleshooting

### API Key no encontrada

```bash
Error: Missing API key for minimax: set MINIMAX_API_KEY

Solución:
export MINIMAX_API_KEY="tu-api-key"
```

### Error al ejecutar claude

```bash
Error: exec: "claude": not found

Solución: Asegúrate de que Claude Code está instalado y en PATH
which claude
```

### Permiso denegado

```bash
Permission denied: /usr/local/bin/cclaude

Solución:
sudo cp cclaude /usr/local/bin/
```

---

## Changelog

### v1.0.0 (2026-01-02)

- ✅ Versión inicial
- ✅ Proveedores: mimo, minimax, kimi, glm
- ✅ Flags: --provider, --model, --dry-run, --verbose, --help, --version
- ✅ Configuración YAML
- ✅ Tests unitarios
- ✅ Cross-compilation para Linux y macOS

---

## Licencia

MIT License - Ver archivo LICENSE para detalles.

---

## Contribuciones

1. Fork el repositorio
2. Crea una rama (`git checkout -b feature/nueva-funcionalidad`)
3. Commit tus cambios (`git commit -am 'Agrega nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Crea un Pull Request

---

**Documentación generada:** 25 minutos
