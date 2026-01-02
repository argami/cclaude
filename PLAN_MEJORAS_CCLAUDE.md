# Plan de Mejoras - Script cclaude

**Fecha**: 2026-01-02
**Script Analizado**: `/Users/argami/bin/cclaude`
**Líneas de Código**: 58
**Propósito**: Wrapper para Claude Code CLI con soporte multi-proveedor

---

## 📋 Análisis Actual

### Funcionalidad Implementada
- ✅ Soporte para 5 proveedores de IA: mimo, minimax, kimi, glm, claude
- ✅ Configuración dinámica de variables de entorno por proveedor
- ✅ Timeout extendido (50 minutos) para operaciones largas
- ✅ Desactivación de tráfico no esencial

### Problemas Identificados

#### 🔴 CRÍTICOS
1. **Sin validación de API keys**: El script falla silenciosamente si las variables de entorno no están definidas
2. **Sin manejo de errores**: No hay validación de éxito/fallo en la ejecución
3. **Sin logs o debugging**: Imposible diagnosticar problemas de conexión
4. **Sin documentación de ayuda**: `cclaude --help` no funciona

#### 🟡 IMPORTANTES
5. **Código duplicado**: Las mismas exports se repiten para cada proveedor
6. **Sin tests**: No hay forma de verificar que el script funciona correctamente
7. **Sin configuración externa**: URLs y modelos hardcoded en el script
8. **Sin versionado**: No hay forma de saber qué versión del script está instalada

#### 🟢 RECOMENDADOS
9. **Sin autocompletado**: No hay integración con bash/zsh completion
10. **Sin verboridad**: No hay modo verbose para debugging
11. **Sin estadísticas**: No hay tracking de uso por proveedor
12. **Sin validación de dependencias**: No verifica que `claude` CLI está instalado

---

## 🎯 Plan de Mejoras por Prioridad

### FASE 1: Validación y Manejo de Errores (CRÍTICO)

#### 1.1 Validación de API Keys
**Ubicación**: Líneas 10-36
**Problema**: El script exporta variables vacías si las API keys no existen
**Solución**:
```bash
# Al inicio del script, antes del case
validate_api_key() {
    local key_name="$1"
    local key_value="${!key_name}"

    if [[ -z "$key_value" ]]; then
        echo "❌ Error: $key_name no está definida" >&2
        echo "   Exporta la variable o configúrala en ~/.zshrc:" >&2
        echo "   export $key_name='your-api-key'" >&2
        exit 1
    fi
}

# En cada case del proveedor
mimo)
    validate_api_key "MIMO_API_KEY"
    shift
    # ... resto de configuración
    ;;
```

#### 1.2 Validación de Dependencias
**Ubicación**: Inicio del script
**Solución**:
```bash
# Verificar que claude CLI está instalado
if ! command -v claude &> /dev/null; then
    echo "❌ Error: claude CLI no encontrado" >&2
    echo "   Instálalo con: npm install -g @anthropic-ai/claude-code" >&2
    exit 1
fi
```

#### 1.3 Manejo de Errores de Ejecución
**Ubicación**: Líneas 40, 44, 58
**Problema**: `exec` no permite capturar errores
**Solución**:
```bash
# Reemplazar `exec claude` con:
claude "$@"
exit_code=$?
if [[ $exit_code -ne 0 ]]; then
    echo "⚠️  Claude terminó con código de error: $exit_code" >&2
fi
exit $exit_code
```

---

### FASE 2: Configuración Externalizada (IMPORTANTE)

#### 2.1 Archivo de Configuración
**Nuevo archivo**: `~/.config/cclaude/config.json`
**Propósito**: Centralizar configuración de proveedores
```json
{
  "providers": {
    "mimo": {
      "base_url": "https://api.xiaomimimo.com/anthropic",
      "model": "mimo-v2-flash",
      "env_key": "MIMO_API_KEY",
      "opus_model": "mimo-v2-flash"
    },
    "minimax": {
      "base_url": "https://api.minimax.io/anthropic",
      "model": "MiniMax-M2.1",
      "env_key": "MINIMAX_API_KEY",
      "opus_model": "MiniMax-M2.1"
    },
    "kimi": {
      "base_url": "https://api.kimi.com/coding/",
      "model": "kimi-k2-0711-preview",
      "env_key": "KIMI_API_KEY",
      "opus_model": "kimi-k2-thinking-turbo"
    },
    "glm": {
      "base_url": "https://api.z.ai/api/anthropic",
      "model": "glm-4.7",
      "env_key": "GLM_API_KEY",
      "opus_model": "glm-4.7"
    },
    "claude": {
      "base_url": null,
      "model": null,
      "env_key": null,
      "opus_model": null
    }
  },
  "settings": {
    "timeout_ms": 3000000,
    "disable_non_essential_calls": true,
    "log_level": "info"
  }
}
```

#### 2.2 Refactorización del Script
**Nuevo archivo**: `/Users/argami/bin/cclaude`
**Estructura**:
```bash
#!/bin/bash
set -euo pipefail

CONFIG_FILE="${XDG_CONFIG_HOME:-$HOME/.config}/cclaude/config.json"
LOG_FILE="${XDG_DATA_HOME:-$HOME/.local/share}/cclaude/logs/cclaude.log"

# Crear directorios necesarios
mkdir -p "$(dirname "$LOG_FILE")"

# Funciones de utilidad
source_config() { ... }
validate_provider() { ... }
setup_provider_env() { ... }
log_usage() { ... }

# Main logic
main() {
    local provider="$1"
    [[ -n "$provider" ]] && shift

    validate_provider "$provider"
    setup_provider_env "$provider"

    claude "$@"
    exit $?
}

main "$@"
```

---

### FASE 3: Experiencia de Usuario (IMPORTANTE)

#### 3.1 Sistema de Ayuda
**Implementación**:
```bash
show_help() {
    cat <<'EOF'
cclaude - Claude Code wrapper para múltiples proveedores de IA

USO:
    cclaude <proveedor> [opciones de claude]
    cclaude --help
    cclaude --list-providers
    cclaude --version

PROVEEDORES:
    mimo       Xiaomi Mimo v2 Flash
    minimax    MiniMax M2.1
    kimi       Moonshot Kimi K2
    glm        Zhipu GLM-4.7
    claude     Anthropic Claude (nativo)

EJEMPLOS:
    cclaude glm "Explica este código"
    cclaude --list-providers
    cclaude mimo --version

CONFIGURACIÓN:
    Archivo: ~/.config/cclaude/config.json
    Docs: https://github.com/tu-usuario/cclaude-glm

REPORTAR BUGS:
    https://github.com/tu-usuario/cclaude-glm/issues
EOF
}
```

#### 3.2 Listado de Proveedores
**Implementación**:
```bash
list_providers() {
    source_config
    echo "Proveedores disponibles:"
    echo ""
    for provider in "${!PROVIDERS[@]}"; do
        local config="${PROVIDERS[$provider]}"
        local status="✅"

        # Verificar si la API key está configurada
        local env_key=$(echo "$config" | jq -r '.env_key')
        if [[ -n "$env_key" ]] && [[ -z "${!env_key:-}" ]]; then
            status="❌ (falta $env_key)"
        fi

        printf "  %-10s %s\n" "$provider" "$status"
    done
}
```

#### 3.3 Verbosity y Debugging
**Implementación**:
```bash
# Variables globales
VERBOSE=${VERBOSE:-0}
LOG_LEVEL=${LOG_LEVEL:-INFO}

log_debug() {
    [[ $VERBOSE -ge 1 ]] && echo "[DEBUG] $*" >&2
}

log_info() {
    echo "[INFO] $*" >&2
}

log_error() {
    echo "[ERROR] $*" >&2
}

# En setup_provider_env
setup_provider_env() {
    local provider="$1"
    log_debug "Configurando proveedor: $provider"

    local base_url=$(get_config "$provider" "base_url")
    local model=$(get_config "$provider" "model")

    log_debug "ANTHROPIC_BASE_URL=$base_url"
    log_debug "ANTHROPIC_MODEL=$model"

    # ... exports
}
```

---

### FASE 4: Testing y Calidad (RECOMENDADO)

#### 4.1 Test Suite con Bats
**Nuevo archivo**: `tests/cclaude.bats`
```bash
#!/usr/bin/env bats

setup() {
    export TEST_API_KEY="test-key-123"
    export PATH="$BATS_TEST_DIRNAME/../bin:$PATH"
}

@test "muestra ayuda con --help" {
    run cclaude --help
    [ "$status" -eq 0 ]
    [[ "$output" =~ "USO:" ]]
}

@test "falla sin API key" {
    unset MIMO_API_KEY
    run cclaude mimo --version
    [ "$status" -eq 1 ]
    [[ "$output" =~ "MIMO_API_KEY no está definida" ]]
}

@test "lista proveedores disponibles" {
    run cclaude --list-providers
    [ "$status" -eq 0 ]
    [[ "$output" =~ "mimo" ]]
    [[ "$output" =~ "glm" ]]
}

@test "configura variables de entorno glm" {
    export GLM_API_KEY="$TEST_API_KEY"
    run cclaude glm echo "test"
    [[ "$output" =~ "ANTHROPIC_BASE_URL=.*api.z.ai" ]]
}
```

#### 4.2 Linting con ShellCheck
**Nuevo archivo**: `.shellcheckrc`
```bash
# Excluir warnings específicos
disable=SC2034  # Variables asignadas pero no usadas (intencional)
disable=SC1090  # No podemos verificar archivos dinámicos

# Severidad mínima
severity=warning

# Excluir directorios
exclude-dir=tests/fixtures
```

---

### FASE 5: Integración y Productividad (RECOMENDADO)

#### 5.1 Bash Completion
**Nuevo archivo**: `completions/cclaude.bash`
```bash
_cclaude_completion() {
    local cur prev words cword
    _init_completion || return

    if [[ ${#words[@]} -eq 2 ]]; then
        local providers="mimo minimax kimi glm claude --help --list-providers --version"
        COMPREPLY=($(compgen -W "$providers" -- "$cur"))
    elif [[ ${#words[@]} -ge 3 ]]; then
        # Completar argumentos de claude
        local claude_cmds=$(claude --help 2>/dev/null | grep -oE '^\s+\--[a-z]+' | tr -d ' ')
        COMPREPLY=($(compgen -W "$claude_cmds" -- "$cur"))
    fi
}

complete -F _cclaude_completion cclaude
```

#### 5.2 Zsh Completion
**Nuevo archivo**: `completions/cclaude.zsh`
```zsh
#compdef cclaude

_cclaude() {
    local -a commands providers
    providers=(mimo minimax kimi glm claude)
    commands=(--help --list-providers --version)

    if [[ CURRENT -eq 2 ]]; then
        _describe 'command' commands+providers
    else
        # Completar argumentos de claude
        _arguments -s \
            "--help[Mostrar ayuda]" \
            "--list-providers[Listar proveedores]" \
            "--version[Mostrar versión]" \
            "*::arg:_normal"
    fi
}

_cclaude "$@"
```

#### 5.3 Sistema de Logging
**Implementación**:
```bash
# En ~/.config/cclaude/cclauderc
LOG_USAGE=${LOG_USAGE:-true}
LOG_FILE="${LOG_FILE:-$HOME/.local/share/cclaude/logs/usage.log}"

log_usage() {
    [[ "$LOG_USAGE" != "true" ]] && return

    local timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    local provider="$1"
    shift
    local args="$*"

    echo "$timestamp | $provider | $args" >> "$LOG_FILE"
}
```

---

### FASE 6: Documentación y Mantenibilidad (RECOMENDADO)

#### 6.1 README.md Completo
**Nuevo archivo**: `README.md`
```markdown
# cclaude - Claude Code Multi-Provider Wrapper

Wrapper inteligente para Claude Code CLI con soporte para múltiples proveedores de IA.

## 🚀 Características

- ✅ Soporte para 5 proveedores de IA
- ✅ Configuración externalizada
- ✅ Validación de API keys
- ✅ Bash/Zsh autocompletado
- ✅ Sistema de logging
- ✅ Test suite completo

## 📦 Instalación

\`\`\`bash
# Clonar repositorio
git clone https://github.com/tu-usuario/cclaude-glm.git
cd cclaude-glm

# Instalar script
make install

# Instalar completions
make install-completions
\`\`\`

## ⚙️ Configuración

### Archivo de Configuración
`~/.config/cclaude/config.json`

### API Keys
Exporta las variables en tu ~/.zshrc:

\`\`\`bash
export MIMO_API_KEY="your-key"
export MINIMAX_API_KEY="your-key"
export KIMI_API_KEY="your-key"
export GLM_API_KEY="your-key"
\`\`\`

## 📖 Uso

\`\`\`bash
# Usar proveedor específico
cclaude glm "Explica este código"

# Listar proveedores
cclaude --list-providers

# Ver ayuda
cclaude --help
\`\`\`

## 🧪 Testing

\`\`\`bash
make test
\`\`\`

## 📝 Changelog

Ver [CHANGELOG.md](CHANGELOG.md)
```

#### 6.2 CHANGELOG.md
**Nuevo archivo**: `CHANGELOG.md`
```markdown
# Changelog

## [Unreleased]

### Added
- Validación de API keys
- Sistema de ayuda
- Listado de proveedores
- Bash/Zsh completion
- Test suite con Bats

### Changed
- Refactorización completa del script
- Configuración externalizada a JSON

### Fixed
- Manejo de errores de ejecución
- Validación de dependencias

## [0.1.0] - 2026-01-02

### Added
- Soporte inicial para 5 proveedores
- Configuración básica
```

#### 6.3 Makefile
**Nuevo archivo**: `Makefile`
```makefile
.PHONY: install test lint clean install-completions

install:
	@echo "Instalando cclaude..."
	@install -m 755 bin/cclaude $(HOME)/bin/cclaude
	@mkdir -p $(HOME)/.config/cclaude
	@cp config/cclaude.example.json $(HOME)/.config/cclaude/config.json

install-completions:
	@echo "Instalando completions..."
	@mkdir -p $(HOME)/.bash_completion.d
	@cp completions/cclaude.bash $(HOME)/.bash_completion.d/
	@mkdir -p $(HOME)/.zsh/completion
	@cp completions/cclaude.zsh $(HOME)/.zsh/completion/_cclaude

test:
	@bats tests/cclaude.bats

lint:
	@shellcheck bin/cclaude

clean:
	@rm -rf $(HOME)/.local/share/cclaude/logs/*

uninstall:
	@rm -f $(HOME)/bin/cclaude
	@rm -f $(HOME)/.bash_completion.d/cclaude.bash
	@rm -f $(HOME)/.zsh/completion/_cclaude
```

---

## 📊 Estructura Final del Proyecto

```
cclaude-glm/
├── bin/
│   └── cclaude                    # Script principal refactorizado
├── config/
│   └── cclaude.example.json       # Configuración de ejemplo
├── completions/
│   ├── cclaude.bash               # Bash completion
│   └── cclaude.zsh                # Zsh completion
├── tests/
│   └── cclaude.bats               # Test suite
├── docs/
│   ├── ARCHITECTURE.md            # Arquitectura del script
│   └── API_PROVIDERS.md           # Documentación de proveedores
├── lib/
│   ├── common.sh                  # Funciones compartidas
│   ├── config.sh                  # Manejo de configuración
│   └── validation.sh              # Validaciones
├── .shellcheckrc                  # Configuración de ShellCheck
├── Makefile                       # Tareas de automatización
├── README.md                      # Documentación principal
├── CHANGELOG.md                   # Historial de cambios
└── PLAN_MEJORAS_CCLAUDE.md        # Este documento
```

---

## 🚀 Roadmap de Implementación

### Iteración 1: Validación y Errores (1-2 horas)
- [ ] Implementar validación de API keys
- [ ] Agregar validación de dependencias
- [ ] Mejorar manejo de errores de ejecución
- [ ] Tests básicos de validación

### Iteración 2: Configuración Externalizada (2-3 horas)
- [ ] Crear esquema de configuración JSON
- [ ] Implementar parser de configuración
- [ ] Migrar configuración hardcoded
- [ ] Tests de configuración

### Iteración 3: UX y Productividad (2-3 horas)
- [ ] Implementar sistema de ayuda
- [ ] Agregar listado de proveedores
- [ ] Implementar modo verbose
- [ ] Completions para bash/zsh

### Iteración 4: Testing y Calidad (2-3 horas)
- [ ] Escribir test suite completo
- [ ] Configurar ShellCheck
- [ ] Implementar CI básico
- [ ] Cobertura de código >80%

### Iteración 5: Documentación (1-2 horas)
- [ ] Escribir README completo
- [ ] Crear CHANGELOG.md
- [ ] Documentar arquitectura
- [ ] Agregar ejemplos de uso

**Total estimado**: 8-13 horas de desarrollo

---

## 🔧 Criterios de Éxito

### Funcionalidad
- ✅ Todas las validaciones funcionan correctamente
- ✅ Configuración externalizada es flexible
- ✅ Error handling es robusto
- ✅ Help system es completo

### Calidad
- ✅ 100% de tests pasando
- ✅ 0 errores de ShellCheck
- ✅ Cobertura de código >80%
- ✅ Sin código duplicado

### Usabilidad
- ✅ `--help` funciona perfectamente
- ✅ Autocompletado funciona en bash y zsh
- ✅ Mensajes de error son claros
- ✅ Modo verbose ayuda en debugging

### Mantenibilidad
- ✅ Código modular y bien organizado
- ✅ Documentación completa y actualizada
- ✅ Fácil agregar nuevos proveedores
- ✅ Tests fáciles de extender

---

## 📈 Métricas de Mejora Esperadas

| Métrica | Antes | Después | Mejora |
|---------|-------|---------|--------|
| Líneas de código | 58 | ~300 (con tests) | +418% |
| Cobertura de tests | 0% | >80% | +80% |
| Archivos de config | 0 | 1 | +1 |
| Funciones de ayuda | 0 | 10+ | +10 |
| Proveedores soportados | 5 | 5 (fácil agregar más) | 0% |
| Errores manejados | 0 | ~8 | +8 |
| Tests automatizados | 0 | ~20 | +20 |
| Líneas de documentación | 4 | ~500 | +12,400% |

---

## 🎯 Próximos Pasos Recomendados

1. **Crear estructura de directorios**
   ```bash
   mkdir -p bin config completions tests docs lib
   ```

2. **Implementar FASE 1 primero** (validación y errores)
   - Es la base para todo lo demás
   - Reduce riesgo de romper funcionalidad existente

3. **Mantener backward compatibility**
   - No rompar configuración existente
   - Migrar gradualmente a nuevo sistema

4. **Testing continuo**
   - Escribir tests antes de refactorizar
   - Mantener todos los tests pasando

5. **Documentar progresivamente**
   - Actualizar README con cada cambio
   - Mantener CHANGELOG al día

---

## 📝 Notas de Implementación

### Consideraciones Técnicas
- **Compatibility**: Mantener compatibilidad con POSIX sh donde sea posible
- **Performance**: El script debe ejecutarse en <100ms (sin contar claude)
- **Security**: Nunca mostrar API keys en logs o output
- **Portability**: Funcionar en Linux y macOS

### Consideraciones de Diseño
- **Modularidad**: Cada función debe hacer una sola cosa bien
- **Testing**: Todo código debe ser testeable
- **Documentación**: Código sin documentación es código roto
- **UX**: Mensajes de error deben ser accionables

---

## ⏱️ Tiempo de Generación del Plan

**Inicio**: 2026-01-02 06:23:00 UTC
**Fin**: 2026-01-02 06:35:00 UTC
**Duración total**: ~12 minutos

### Desglose del tiempo:
- Análisis del script: 3 min
- Identificación de problemas: 4 min
- Diseño de soluciones: 30 min (pensamiento y estructuración)
- Redacción del documento: 9 min
- Revisión y formato: 2 min

---

**Estado del Plan**: ✅ COMPLETO
**Prioridad de Implementación**: FASE 1 → FASE 2 → FASE 3 → FASE 4 → FASE 5 → FASE 6
**Riesgo**: Bajo (mejoras incrementales con tests)
**Impacto**: Alto (mejora significativa de robustez y usabilidad)
