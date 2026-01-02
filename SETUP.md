# Guía de Setup Rápido

## 🎯 Objetivo
Tener `cclaude.py` funcionando en 5 minutos con TDD completo.

---

## 📦 **Repositorio Oficial**
- **GitHub:** https://github.com/argami/cclaude
- **Rama:** `mimo`
- **Modelo:** MiMo V2 Flash
- **Entorno:** Claude Code + MCPs

---

## ⚡ Flujo de 5 Pasos

### 1. Entrar al Directorio
```bash
cd /Users/argami/Documents/workspace/AI/cclaude/mimo/
```

### 2. Configurar Entorno (1 minuto)
```bash
make setup
```

Esto hará:
- ✅ `chmod +x cclaude.py` (hacer ejecutable)
- ✅ Instalar dependencias (pytest, etc.)
- ✅ Verificar Python 3.8+

### 3. Validar Instalación (1 minuto)
```bash
./cclaude.py --test
```

Deberías ver:
```
🧪 Ejecutando tests de validación...

1. Validación de providers...
   ✅ mimo
   ✅ minimax
   ✅ kimi
   ✅ glm
   ✅ claude nativo

2. Validación variables de entorno...
   ⚠️  mimo: MIMO_API_KEY no está definida
   ⚠️  minimax: MINIMAX_API_KEY no está definida
   ...

RESUMEN DE TESTS
✅ Pasados: X
❌ Fallidos: Y
```

### 4. Configurar API Keys (2 minutos)

**Opción A: Temporal (para probar)**
```bash
export MIMO_API_KEY="tu_key_aqui"
export MINIMAX_API_KEY="tu_key_aqui"
export KIMI_API_KEY="tu_key_aqui"
export GLM_API_KEY="tu_key_aqui"
```

**Opción B: Permanente (recomendado)**
Añadir a `~/.zshrc` o `~/.bashrc`:
```bash
# Claude Code Providers
export MIMO_API_KEY="tu_key_aqui"
export MINIMAX_API_KEY="tu_key_aqui"
export KIMI_API_KEY="tu_key_aqui"
export GLM_API_KEY="tu_key_aqui"

# Alias opcional
alias cclaude="/Users/argami/Documents/workspace/AI/cclaude/mimo/cclaude.py"
```

Recargar:
```bash
source ~/.zshrc  # o ~/.bashrc
```

### 5. Probar Funcionamiento (1 minuto)

```bash
# Ver ayuda
./cclaude.py --help

# Listar providers
./cclaude.py --list-providers

# Probar con mimo (requiere API key)
./cclaude.py mimo --version

# Probar claude nativo (sin API key)
./cclaude.py claude --version
```

---

## 🧪 Ejecutar Tests

### Tests Rápidos
```bash
make test
```

### Tests con Cobertura
```bash
make test-coverage
# Abre htmlcov/index.html para ver reporte
```

### Suite Completa
```bash
make test-all
```

### Tests Manuales con Pytest
```bash
python3 -m pytest tests/test_cclaude.py -v
python3 -m pytest tests/test_cclaude.py -k "test_e2e" -v  # Solo E2E
python3 -m pytest tests/test_cclaude.py --cov=cclaude     # Con cobertura
```

---

## 📋 Verificación Pre-Commit

Antes de commitear cambios:
```bash
make check  # Lint + format check
make test   # Todos los tests pasan
```

---

## 🎯 Uso Diario

### Comandos Cortos (con alias)
```bash
# Añadir a ~/.zshrc:
alias cclaude="/Users/argami/Documents/workspace/AI/cclaude/mimo/cclaude.py"

# Uso:
cclaude mimo --verbose
cclaude minimax --help
cclaude claude --version
```

### Sin Alias
```bash
./cclaude.py mimo --verbose
```

---

## 🔍 Troubleshooting

### Problema: `command not found: claude`
**Solución**: Instalar Claude Code y añadir al PATH

### Problema: `MIMO_API_KEY no está definida`
**Solución**: Añadir variable de entorno (ver paso 4)

### Problema: Tests fallan
**Solución**:
```bash
make clean
make setup
make test
```

### Problema: No tiene permisos de ejecución
**Solución**:
```bash
chmod +x cclaude.py
```

---

## 📊 Estructura de Archivos

```
mimo/
├── cclaude.py              # ⭐ Script principal (ejecutable)
├── tests/
│   └── test_cclaude.py     # ⭐ Tests TDD (100+ casos)
├── requirements.txt         # Dependencias
├── Makefile                # 🤖 Automatización
├── README.md               # Documentación completa
├── SETUP.md                # Este archivo
└── .gitignore              # Git
```

---

## 🚀 Próximos Pasos

1. ✅ **Listo**: Setup básico funcionando
2. 📝 **Opcional**: Añadir a tu PATH global
3. 🧪 **Opcional**: Personalizar tests
4. 🔧 **Opcional**: Añadir nuevos providers

### Añadir a PATH Global
```bash
# Copiar a /usr/local/bin (requiere sudo)
sudo cp cclaude.py /usr/local/bin/cclaude
sudo chmod +x /usr/local/bin/cclaude

# Ahora puedes usar desde cualquier lugar
cclaude mimo --help
```

---

## 🎓 Comandos Make Reference

```bash
make help          # Ver todos los comandos
make install       # Instalar dependencias
make test          # Ejecutar tests
make test-coverage # Tests + cobertura
make lint          # Verificar calidad
make format        # Formatear código
make check         # Lint + format check
make setup         # Configurar todo
make clean         # Limpiar temporales
make test-all      # Suite completa
make run-mimo      # Ejemplo rápido
make run-claude    # Ejemplo nativo
```

---

**✅ Setup completado!** Ahora puedes usar `cclaude.py` con toda la seguridad y tests que tu framework requiere.