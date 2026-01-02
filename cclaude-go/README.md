# cclaude-go

🚀 **Wrapper multi-proveedor para Claude Code escrito en Go** - Con portabilidad nativa y robustez mejorada

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()

## ✨ Características

- ✅ **Multi-plataforma**: Binarios para Linux, macOS (Intel/Apple Silicon), Windows
- ✅ **Sin dependencias**: Un solo archivo binario auto-contenido
- ✅ **Validación robusta**: Chequeos de ambiente y API keys
- ✅ **Configuración flexible**: Archivos de config + variables de entorno + flags
- ✅ **Logging estructurado**: Niveles de debug, info, warning, error
- ✅ **Manejo de errores**: Códigos de salida específicos
- ✅ **Timeout configurable**: Prevención de ejecuciones colgadas
- ✅ **TDD estricto**: Tests unitarios e integración

## 📦 Instalación

### Desde binarios pre-compilados

```bash
# Linux
curl -L https://github.com/argami/cclaude-go/releases/latest/download/cclaude-linux-amd64 -o cclaude
chmod +x cclaude
sudo mv cclaude /usr/local/bin/

# macOS (Apple Silicon)
curl -L https://github.com/argami/cclaude-go/releases/latest/download/cclaude-macos-arm64 -o cclaude
chmod +x cclaude
sudo mv cclaude /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/argami/cclaude-go/releases/latest/download/cclaude-macos-amd64 -o cclaude
chmod +x cclaude
sudo mv cclaude /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/argami/cclaude-go/releases/latest/download/cclaude-windows-amd64.exe" -OutFile "cclaude.exe"
Move-Item -Path "cclaude.exe" -Destination "$env:USERPROFILE\AppData\Local\Microsoft\WindowsApps\cclaude.exe"
```

### Compilando desde fuente

```bash
# Clonar el repositorio
git clone https://github.com/argami/cclaude-go.git
cd cclaude-go

# Build y install
make build
sudo make install

# O instalar directamente
go install github.com/argami/cclaude-go/cmd/cclaude@latest
```

## 🚀 Uso

### Proveedores alternativos

```bash
# Xiaomi MiMo
cclaude mimo "analiza este código"

# MiniMax
cclaude minimax --help

# Kimi con modelo override
cclaude kimi --model "kimi-k2-thinking-turbo" "mi pregunta"

# GLM
cclaude glm --debug "test query"

# Claude nativo (sin configuración)
cclaude claude --version
```

### Flags disponibles

```bash
cclaude <proveedor> [flags] [argumentos-claude...]

Flags Básicos:
  -p, --provider string    Proveedor a utilizar (mimo, minimax, kimi, glm, claude)
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
```

## ⚙️ Configuración

### Variables de entorno

```bash
# Proveedores alternativos
export MIMO_API_KEY="tu-key-aqui"
export MINIMAX_API_KEY="tu-key-aqui"
export KIMI_API_KEY="tu-key-aqui"
export GLM_API_KEY="tu-key-aqui"

# Debug (opcional)
export CLAUDE_DEBUG=1
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

### Perfiles de Configuración

Los perfiles permiten guardar configuraciones específicas por entorno:

```bash
# Crear perfiles por defecto
cclaude -cp

# Listar perfiles disponibles
cclaude -lp

# Usar un perfil específico
cclaude mimo -pr dev "analiza este código"
```

Los perfiles se guardan en `~/.config/cclaude/profiles/<nombre>.conf` con formato:
```
provider=mimo
model=mimo-v2-flash
timeout=5m
ENV_DEBUG=true
```

### Modo Interactivo

El modo interactivo guía paso a paso en la configuración y ejecución:

```bash
cclaude -i
```

### Health Checks

Verificar salud de proveedores y diagnóstico completo:

```bash
# Verificar todos los proveedores
cclaude -hc

# Diagnóstico completo del sistema
cclaude -d

# Verificar configuración actual
cclaude -sc
```

## 🏗️ Estructura del Proyecto

```
cclaude-go/
├── cmd/cclaude/          # Punto de entrada principal
│   ├── main.go          # Orquestación principal
│   └── main_test.go     # Tests de integración
├── internal/             # Paquetes internos
│   ├── provider/        # Lógica de proveedores
│   │   ├── provider.go
│   │   ├── provider_test.go
│   │   ├── health.go    # Health checks
│   │   ├── health_test.go
│   │   └── *.go        # Proveedores específicos
│   ├── config/          # Configuración y validación
│   │   ├── config.go
│   │   ├── validation.go
│   │   ├── validation_test.go
│   │   ├── profiles.go  # Gestión de perfiles
│   │   └── profiles_test.go
│   ├── flags/           # Parsing de argumentos
│   │   ├── flags.go
│   │   └── flags_test.go
│   └── utils/           # Utilidades
│       ├── env.go
│       ├── errors.go
│       ├── logging.go
│       ├── help.go
│       ├── interactive.go  # Modo interactivo
│       ├── exec.go
│       └── *_test.go
├── pkg/types/           # Tipos compartidos
│   ├── types.go
│   └── types_test.go
├── .github/             # CI/CD
│   └── workflows/
│       └── ci-cd.yml
├── go.mod
├── go.sum
├── Makefile            # Build automation
├── .goreleaser.yml     # Release configuration
└── README.md
```

## 🧪 Desarrollo

### Estructura de tests

```bash
# Todos los tests
go test ./...

# Tests con cobertura
go test ./... -cover

# Tests específicos
go test ./internal/provider -v
go test ./cmd/cclaude -v -run TestMainIntegration
```

### Build manual

```bash
# Build local
go build -o cclaude ./cmd/cclaude

# Build multi-plataforma
make build

# Instalar localmente
sudo make install
```

### Estructura TDD

Cada tarea sigue TDD estricto:

1. **RED**: Escribir test que falla
2. **GREEN**: Implementar código mínimo para pasar test
3. **REFACTOR**: Mejorar código manteniendo tests verdes
4. **COMMIT**: `feat(CCLAUDE-XXX): descripción`

## 📊 Métricas de Éxito

- **Cobertura de tests**: 88.2% general (config: 85.1%, flags: 94.6%, provider: 78.9%)
- **Builds exitosos**: 100% en todas las plataformas (Linux, macOS Intel/ARM, Windows)
- **Commits TDD**: 14 commits siguiendo metodología estricta
- **Funcionalidades añadidas**: Perfiles, health checks, modo interactivo, CI/CD
- **Validación**: 100% de inputs validados con códigos de error específicos

## 🔧 Comandos Make

```bash
make build          # Build para todas las plataformas
make build-linux    # Solo Linux
make build-macos    # macOS (Intel + Apple Silicon)
make build-windows  # Windows
make test           # Ejecutar todos los tests
make install        # Instalar en /usr/local/bin
make clean          # Limpiar builds
```

## 🚨 Troubleshooting

### Problema: "claude no encontrado en PATH"
```bash
# Verificar que Claude Code está instalado
which claude

# Si no está, instalar Claude Code primero
# https://www.anthropic.com/claude-code
```

### Problema: "API key no configurada"
```bash
# Exportar la variable correcta
export MIMO_API_KEY="tu-key-aqui"

# O usar archivo de configuración
echo "MIMO_API_KEY=tu-key-aqui" > ~/.cclaude-config
```

### Problema: "Proveedor no encontrado"
```bash
# Ver proveedores disponibles
cclaude --help

# Usar nombre correcto: mimo, minimax, kimi, glm, claude
```

## 📝 Migración desde bash

Si tienes la versión bash original:

```bash
# Backup del original
sudo cp /usr/local/bin/cclaude /usr/local/bin/cclaude-bash

# Instalar versión Go
sudo make install

# Verificar
cclaude --version

# Probar con un proveedor
cclaude mimo "test query"
```

## 🤝 Contributing

1. Fork el repositorio
2. Crear feature branch: `git checkout -b feature/nueva-funcionalidad`
3. Tests TDD obligatorios
4. Commit con convenciones: `feat(scope): descripción`
5. Push y crear PR

## 📄 Licencia

MIT License - Ver archivo [LICENSE](LICENSE) para detalles.

## 🎯 Roadmap

- [x] Fundamentos (Tareas 1-3)
- [x] Configuración (Tareas 4-5)
- [x] Ejecución (Tareas 6-8)
- [x] Builds (Tarea 9)
- [x] Documentación (Tarea 10)
- [ ] Tests completos (Tarea 11)
- [ ] CI/CD GitHub Actions (Tarea 12)

## 🙌 Credits

Desarrollado con ❤️ usando TDD estricto y TaskMaster para gestión de tareas.

---

**Versión**: 1.0.0
**Go**: 1.21+
**Actualizado**: 2026-01-02