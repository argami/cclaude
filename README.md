# cclaude-glm

**Claude Code wrapper para múltiples proveedores de IA**

cclaude es un wrapper para Claude Code CLI que permite usar diferentes proveedores de IA (Mimo, MiniMax, Kimi, GLM, Claude nativo) con un solo binario compilado y sin dependencias de runtime.

## 🚀 Características

- ✅ **Single binary** - Un solo ejecutable para todas las plataformas
- ✅ **Zero runtime dependencies** - Solo el binario compilado
- ✅ **Multi-plataforma** - Linux, macOS, Windows (amd64/arm64)
- ✅ **Configuración externalizada** - YAML en `~/.config/cclaude/`
- ✅ **CLI completo** - Ayuda integrada con Cobra
- ✅ **Fácil instalación** - `curl + chmod` y listo

## 📦 Instalación

### Requisitos Previos

- Claude Code CLI instalado y disponible en el PATH
- API keys de los proveedores configuradas como variables de entorno

### Binarios Precompilados

Descarga el binario para tu plataforma:

```bash
# Linux amd64
curl -L https://github.com/argami/cclaude-glm/releases/latest/download/cclaude-linux-amd64 -o cclaude
chmod +x cclaude

# macOS amd64 (Intel)
curl -L https://github.com/argami/cclaude-glm/releases/latest/download/cclaude-darwin-amd64 -o cclaude
chmod +x cclaude

# macOS arm64 (Apple Silicon)
curl -L https://github.com/argami/cclaude-glm/releases/latest/download/cclaude-darwin-arm64 -o cclaude
chmod +x cclaude

# Windows amd64
curl -L https://github.com/argami/cclaude-glm/releases/latest/download/cclaude-windows-amd64.exe -o cclaude.exe
```

### Desde Fuente

```bash
# Clonar repositorio
git clone https://github.com/argami/cclaude-glm.git
cd cclaude-glm

# Compilar
go build -o cclaude ./cmd/cclaude

# Instalar (opcional)
sudo mv cclaude /usr/local/bin/
```

## ⚙️ Configuración

### Variables de Entorno

Configura las API keys de los proveedores:

```bash
export MIMO_API_KEY="tu-api-key-mimo"
export MINIMAX_API_KEY="tu-api-key-minimax"
export KIMI_API_KEY="tu-api-key-kimi"
export GLM_API_KEY="tu-api-key-glm"
```

### Archivo de Configuración (Opcional)

Crea `~/.config/cclaude/config.yaml` para customizar proveedores:

```yaml
providers:
  mimo:
    name: Mimo
    base_url: https://api.xiaomimimo.com/anthropic
    model: mimo-v2-flash
    env_key: MIMO_API_KEY
    opus_model: mimo-v2-flash

  minimax:
    name: MiniMax
    base_url: https://api.minimax.io/anthropic
    model: MiniMax-M2.1
    env_key: MINIMAX_API_KEY
    opus_model: MiniMax-M2.1

  kimi:
    name: Kimi
    base_url: https://api.kimi.com/coding/
    model: kimi-k2-0711-preview
    env_key: KIMI_API_KEY
    opus_model: kimi-k2-thinking-turbo

  glm:
    name: GLM
    base_url: https://api.z.ai/api/anthropic
    model: glm-4.7
    env_key: GLM_API_KEY
    opus_model: glm-4.7

  claude:
    name: Claude
    base_url: ""
    model: ""
    env_key: ""
    opus_model: ""

settings:
  timeout_ms: 3000000
  disable_non_essential_calls: true
  log_level: info
```

### Ubicaciones del Archivo de Configuración

El archivo `config.yaml` se busca en el siguiente orden:

1. **Directorio actual** - `./config.yaml` (para desarrollo)
2. **XDG_CONFIG_HOME** - `$XDG_CONFIG_HOME/cclaude/config.yaml`
3. **Home directory** - `~/.config/cclaude/config.yaml`
4. **System-wide** - `/etc/cclaude/config.yaml`

### Hot-Reload de Configuración

La configuración se recarga automáticamente cuando detecta cambios en el archivo:

```bash
# Iniciar el modo de observación (automático al cargar config)
cclaude --watch

# Simplemente edita el archivo config.yaml
vim ~/.config/cclaude/config.yaml

# Los cambios se aplican automáticamente sin reiniciar
```

**Nota**: El hot-reload usa `fsnotify` para detectar cambios en tiempo real.

### Comandos de Configuración

```bash
# Inicializar configuración por defecto
cclaude init

# Crear config con --force (sobrescribe existente)
cclaude init --force

# Mostrar dónde está el archivo de configuración
cclaude config-path
# Output: Config file: /home/user/.config/cclaude/config.yaml

# Validar configuración
cclaude validate
# Muestra estado de cada provider y validación de API keys
```


## 📖 Uso

### Comandos Básicos

```bash
# Mostrar ayuda
cclaude --help
cclaude -h

# Listar proveedores disponibles
cclaude list
cclaude ls

# Mostrar versión
cclaude version
cclaude v

# Inicializar configuración por defecto
cclaude init

# Mostrar ubicación del archivo de configuración
cclaude config-path

# Validar configuración y estado de providers
cclaude validate
```

### Usar con un Proveedor

```bash
# Usar proveedor GLM
cclaude glm "Explica este código"

# Usar proveedor Mimo
cclaude mimo "Ayúdame con este error"

# Usar proveedor Kimi
cclaude kimi "Optimiza este rendimiento"

# Pasar argumentos adicionales a Claude
cclaude glm --version
cclaude minimax --help
```

### Proveedores Disponibles

| Proveedor | Descripción | Modelo |
|-----------|-------------|--------|
| `mimo` | Xiaomi Mimo v2 Flash | mimo-v2-flash |
| `minimax` | MiniMax M2.1 | MiniMax-M2.1 |
| `kimi` | Moonshot Kimi K2 | kimi-k2-0711-preview |
| `glm` | Zhipu GLM-4.7 | glm-4.7 |
| `claude` | Anthropic Claude (nativo) | Default |

## 🛠️ Desarrollo

### Estructura del Proyecto

```
cclaude-glm/
├── cmd/
│   └── cclaude/
│       └── main.go                 # Entry point
├── internal/
│   ├── cli/
│   │   ├── root.go               # Comando raíz
│   │   ├── provider.go           # Comando provider
│   │   ├── list.go              # Comando list
│   │   ├── version.go            # Comando version
│   │   ├── init.go               # Comando init (FASE-2)
│   │   ├── config_path.go        # Comando config-path (FASE-2)
│   │   └── validate.go           # Comando validate (FASE-2)
│   ├── config/
│   │   ├── loader.go             # Carga configuración con hot-reload (FASE-2)
│   │   └── loader_test.go        # Tests de configuración (FASE-2)
│   └── execution/
│       └── executor.go            # Ejecuta claude CLI
├── go.mod                          # Go modules
├── go.sum                          # Go module checksums
└── README.md
```

### Compilar desde Fuente

```bash
# Instalar dependencias
go mod tidy

# Compilar
go build -o cclaude ./cmd/cclaude

# Ejecutar
./cclaude --help
```

### Compilar para Múltiples Plataformas

```bash
# Linux amd64
GOOS=linux GOARCH=amd64 go build -o cclaude-linux-amd64 ./cmd/cclaude

# Linux arm64
GOOS=linux GOARCH=arm64 go build -o cclaude-linux-arm64 ./cmd/cclaude

# macOS amd64 (Intel)
GOOS=darwin GOARCH=amd64 go build -o cclaude-darwin-amd64 ./cmd/cclaude

# macOS arm64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o cclaude-darwin-arm64 ./cmd/cclaude

# Windows amd64
GOOS=windows GOARCH=amd64 go build -o cclaude-windows-amd64.exe ./cmd/cclaude
```

## 📋 Plan de Implementación

Este proyecto sigue un plan de 6 fases:

- ✅ **FASE 1**: Estructura Base y CLI - **Completado**
  - CLI básica con Cobra
  - Sistema de providers con interfaz limpia
  - Executor pattern
  - Tests unitarios, integración y E2E
  - Pre-commit hooks (fmt, lint, test, commit-msg)
  - Validación de configuración
- ✅ **FASE 2**: Sistema de Configuración - **Completado**
  - File watching con fsnotify para hot-reload
  - Config reloading sin reiniciar aplicación
  - Múltiples rutas de configuración (XDG, ~/.config, actual, /etc)
  - Comandos CLI: init, config-path, validate
  - Thread-safe config access con mutexes
  - Observer pattern para notificaciones de cambios
- ⏳ **FASE 3**: Sistema de Providers Avanzado
- ⏳ **FASE 4**: Testing Extensivo
- ⏳ **FASE 5**: Multi-Platform Builds
- ⏳ **FASE 6**: Completions y Features Avanzadas

### FASE-1 Detalles Completados

**Testing**:
- ✅ Unit tests para provider factory
- ✅ Unit tests para BaseProvider methods
- ✅ Integration tests para executor
- ✅ E2E tests para CLI commands

**Calidad**:
- ✅ Pre-commit hooks con go-fmt, ruff, go-test
- ✅ Commit message validation (Conventional Commits)
- ✅ Build check automático
- ✅ Config validation module

## 🤝 Contribuyendo

Contribuciones son bienvenidas! Por favor:

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📝 Licencia

Este proyecto está licenciado bajo la MIT License.

## 🔗 Comparación con Script Original

### Ventajas sobre Bash

| Característica | Bash Script | Go Binary |
|----------------|-------------|------------|
| Portabilidad | Requiere Bash | Single binary multi-plataforma |
| Dependencias | Bash + utilidades | Zero runtime dependencies |
| Performance | Interpretado | Compilado (más rápido) |
| Distribución | Script + perms | Solo binario |
| Type Safety | Dinámico | Estático |
| Testing | Difícil | Nativo (go test) |

### Migración desde Script Bash

Si vienes del script Bash original:

```bash
# Antes (Bash)
cclaude glm "algún texto"
```

```bash
# Ahora (Go)
cclaude glm "algún texto"
```

La sintaxis es casi idéntica, pero con un binario compilado en lugar de un script.

## 📚 Referencias

- [Claude Code Documentation](https://docs.anthropic.com/)
- [Cobra Documentation](https://github.com/spf13/cobra)
- [Viper Documentation](https://github.com/spf13/viper)
- [fsnotify Documentation](https://github.com/fsnotify/fsnotify)
- [Go Modules](https://go.dev/doc/modules/create)
- [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html)
