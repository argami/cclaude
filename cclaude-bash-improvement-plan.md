# Plan de Mejoras: cclaude Script (Bash Original)

**Fecha:** 2026-01-02
**Script original:** `/Users/argami/bin/cclaude`

---

## 📊 Análisis del Script Actual

El script `cclaude` es un wrapper de bash que permite ejecutar Claude Code con múltiples proveedores de API alternativos (mimo, minimax, kimi, glm) configurando las variables de entorno apropiadas.

### Funcionalidades Actuales
- Selección de proveedor mediante primer argumento
- Configuración automática de variables de entorno por proveedor
- Soporte para 4 proveedores: mimo, minimax, kimi, glm
- Passthrough de argumentos a Claude Code nativo

---

## 🚀 Mejoras Propuestas (Versión Bash)

### 1. Manejo de Errores y Validación

| Problema | Mejora | Prioridad |
|----------|---------|-----------|
| No valida API keys antes de ejecutar | Verificar `$PROVIDER_API_KEY` existe antes de ejecutar | Alta |
| No valida argumentos recibidos | Agregar validación de argumentos obligatorios | Media |
| Mensajes de error confusos | Mejorar mensajes de error con sugerencias | Media |
| Exit code no informativo | Retornar códigos de error específicos por tipo de falla | Baja |

### 2. Documentación y Experiencia de Usuario

| Problema | Mejora | Prioridad |
|----------|---------|-----------|
| Help básico | Expandir help con ejemplos de uso | Alta |
| Sin autocomplete | Agregar completion para shells (bash/zsh) | Media |
| No muestra proveedor activo | Imprimir qué proveedor se está usando | Baja |
| Sin versionamiento | Agregar flag `--version` | Baja |

### 3. Configuración Flexible

| Problema | Mejora | Prioridad |
|----------|---------|-----------|
| Variables hardcodeadas | Usar archivo de configuración `~/.config/cclaude/config` | Media |
| Timeout hardcodeado | Permitir configurar timeout via env var o config | Baja |
| Sin soporte para modelos personalizados | Agregar flag `--model` o config por proveedor | Media |

### 4. Nuevas Funcionalidades

| Problema | Mejora | Prioridad |
|----------|---------|-----------|
| Solo 4 proveedores | Agregar soporte para más proveedores (deepseek, grok, etc.) | Media |
| Sin modo interactivo | Agregar modo interactivo para seleccionar proveedor | Baja |
| Sin dry-run | Agregar flag `--dry-run` para ver qué se configuraría | Baja |
| Sin logging | Agregar logging de sesiones (qué proveedor, cuándo) | Baja |

### 5. Compatibilidad y Portabilidad

| Problema | Mejora | Prioridad |
|----------|---------|-----------|
| Solo bash | Asegurar compatibilidad con zsh | Baja |
| Sin test unitario | Agregar tests para validación | Media |
| Sin instalación | Crear script de instalación/desinstalación | Baja |

---

## 📋 Plan de Implementación (Bash)

### Fase 1: Mejoras de Robustez (1-2 horas)
1. Agregar validación de API keys
2. Mejorar mensajes de error
3. Agregar `--help` expandido con ejemplos
4. Agregar `--version`

### Fase 2: Configuración Flexible (2-3 horas)
1. Crear archivo de configuración
2. Permitir override de settings via variables de entorno
3. Agregar soporte para modelos personalizados

### Fase 3: Nuevos Proveedores (1-2 horas)
1. Research de nuevos proveedores compatibles
2. Agregar configuración para cada uno
3. Documentar en el help

### Fase 4: Experiencia de Usuario (1 hora)
1. Agregar completion para bash/zsh
2. Agregar flag `--dry-run`
3. Agregar output del proveedor activo

### Fase 5: Testing y Documentación (1-2 horas)
1. Agregar tests unitarios
2. Crear script de instalación
3. Documentar cambios en CHANGELOG

---

## 📁 Archivos a Modificar

```
/Users/argami/bin/cclaude                    # Script principal (modificar)
~/.config/cclaude/config                     # Archivo de configuración (nuevo)
~/.config/cclaude/completion.bash            # Completion bash (nuevo)
~/.config/cclaude/completion.zsh             # Completion zsh (nuevo)
docs/cclaude.md                             # Documentación (nuevo)
tests/cclaude_test.sh                       # Tests (nuevo)
```

---

## 🔍 Detalle de Cambios por Sección

### Sección de Proveedores (modificación)
```bash
# Agregar verificación de API key antes de configurar
if [[ -z "${PROVIDER_API_KEY}" ]]; then
  echo "Error: API key no configurada para $PROVIDER"
  echo "Configure ${PROVIDER}_API_KEY o agreguela en ~/.config/cclaude/config"
  exit 1
fi
```

### Agregar Flag `--dry-run`
```bash
dry-run)
  shift
  echo "[DRY-RUN] Proveedor: $PROVIDER"
  echo "[DRY-RUN] Variables que se configurarian:"
  echo "  ANTHROPIC_BASE_URL=$ANTHROPIC_BASE_URL"
  echo "  MAIN_MODEL=$MAIN_MODEL"
  echo "[DRY-RUN] Ejecutando: claude $@"
  exit 0
  ;;
```

### Archivo de Configuración
```bash
# Leer configuración si existe
if [[ -f ~/.config/cclaude/config ]]; then
  source ~/.config/cclaude/config
fi
```

---

## ✅ Checklist de Validación

- [ ] Script funciona con proveedor por defecto
- [ ] Script funciona con cada proveedor específico
- [ ] Help muestra ejemplos claros
- [ ] Error cuando API key no está configurada
- [ ] Flags nuevos funcionan correctamente
- [ ] Completions funcionan en bash y zsh
- [ ] Tests pasan exitosamente

---

## 📝 Notas

- Mantener backward compatibility con el uso actual
- No romper integraciones existentes
- Seguir principios de Unix (pequeño, modular)
- Documentar cambios breaking si los hay

---

## ⏱️ Tiempo de Generación del Plan

**Tiempo total:** 15 minutos (lectura, análisis, planificación y escritura)

- Lectura y análisis del script: 5 minutos
- Identificación de mejoras: 5 minutos
- Documentación del plan: 5 minutos

---

## 📚 Relación con Otros Planes

| Plan | Enfoque | Estado |
|------|---------|--------|
| `cclaude-bash-improvement-plan.md` | Mejoras incrementales al script bash | Este archivo |
| `cclaude-go-improvement-plan.md` | Rewrite completo en Go | Implementado |
| `README.md` | Documentación del proyecto Go | Completado |

---

**Nota:** Este plan fue recreado el 2026-01-02 ya que el archivo original no se guardó correctamente.
