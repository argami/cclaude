# Contributing a cclaude-glm

Gracias por tu interés en contribuir a cclaude-glm! Este documento proporciona directrices y procedimientos para contribuir al proyecto.

## 📋 Índice

- [Código de Conducta](#código-de-conducta)
- [Cómo Contribuir](#cómo-contribuir)
- [Proceso de Desarrollo](#proceso-de-desarrollo)
- [Estándares de Código](#estándares-de-código)
- [Commit Messages](#commit-messages)
- [Testing](#testing)
- [Pull Requests](#pull-requests)

## 🤝 Código de Conducta

Al participar en este proyecto, te comprometes a mantener un ambiente inclusivo y respetuoso. Por favor:

- Ser respetuoso con otros contribuidores
- Usar lenguaje inclusivo
- Aceptar críticas constructivas
- Enfocarse en lo que es mejor para la comunidad

## 🚀 Cómo Contribuir

### Reporting Bugs

Antes de crear un issue, busca si ya existe uno similar. Si encuentras un bug:

1. Usa un título claro y descriptivo
2. Incluye pasos para reproducir el problema
3. Proporciona información del entorno (OS, Go version, etc.)
4. Adjunta logs o screenshots si son relevantes

### Sugerencias de Features

1. Busca issues existentes primero
2. Explica claramente el caso de uso
3. Describe el comportamiento esperado
4. Considera si es alineado con los objetivos del proyecto

## 🔧 Proceso de Desarrollo

### Setup del Entorno

```bash
# Fork y clona tu repositorio
git clone https://github.com/tu-usuario/cclaude-glm.git
cd cclaude-glm

# Agrega el remoto original
git remote add upstream https://github.com/argami/cclaude-glm.git

# Instala dependencias
go mod download

# Instala pre-commit hooks
pre-commit install
```

### Creando una Rama

```bash
# Actualiza tu rama main
git checkout main
git pull upstream main

# Crea una rama para tu feature
git checkout -b feature/nombre-de-tu-feature
# o para un bugfix
git checkout -b fix/nombre-del-bug
```

### Flujo de Trabajo

1. **TDD Primero**: Escribe tests ANTES de escribir código
2. **Codificación**: Implementa la funcionalidad
3. **Testing**: Ejecuta todos los tests
4. **Linting**: Asegúrate de que pase el linting
5. **Documentación**: Actualiza la documentación si es necesario
6. **Commit**: Commitea tus cambios con un mensaje claro
7. **Push**: Envía tus cambios a tu fork
8. **Pull Request**: Crea un PR

## 📝 Estándares de Código

### Guía de Estilo Go

Seguimos las [Effective Go guidelines](https://go.dev/doc/effective_go) y [Uber Go Style Guide](https://github.com/uber-go/guide).

#### Nombres

- **Packages**: `lowercase`, sin guiones bajos
- **Constants**: `MixedCase` o `UPPER_CASE` para exportadas
- **Variables**: `mixedCase` o `camelCase`
- **Interfaces**: Generalmente `-er` sufijo (ej: `Provider`, `Executor`)
- **Tests**: `Test<FunctionName>`

#### Ejemplo

```go
// ✅ Bueno
type Provider interface {
    Name() string
    Validate() error
}

type BaseProvider struct {
    name   string
    envKey string
}

func (p *BaseProvider) Name() string {
    return p.name
}

// ❌ Malo
type provider interface {  // Should be exported
    name string              // Should be exported if used externally
}
```

### Formateo

```bash
# Formatear código
gofmt -w .

# Formatear y simplificar
goimports -w .
```

El proyecto usa pre-commit hooks que ejecutan `gofmt` automáticamente.

## 💬 Commit Messages

Seguimos [Conventional Commits](https://www.conventionalcommits.org/):

### Formato

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Tipos

- `feat`: Nueva feature
- `fix`: Bug fix
- `docs`: Cambios en documentación
- `style**: Cambios de formato (sin lógica)
- `refactor`: Refactorización de código
- `test`: Agregar o actualizar tests
- `chore`: Cambios en build/process/herramientas

### Ejemplos

```
# ✅ Bueno
feat(provider): add support for new XYZ LLM provider

Implements Provider interface for XYZ service with:
- API key validation
- Environment variable setup
- Integration tests

Closes #123

# ✅ Bueno
fix(cli): handle missing config file gracefully

Returns error message instead of panicking when
~/.config/cclaude/config.yaml is not found.

Fixes #456

# ❌ Malo
update stuff
fix bug
add tests
```

### Proceso de Commit

Los pre-commit hooks validarán:
- Formato del código (gofmt)
- Linting (ruff)
- Tests (go test)
- Formato del commit message

Si algo falla, el commit será rechazado.

## ✅ Testing

### Tipos de Tests

1. **Unit Tests**: Prueban funciones individuales
2. **Integration Tests**: Prueban interacción entre componentes
3. **E2E Tests**: Prueban flujos completos de la CLI

### Escribiendo Tests

```go
// ✅ Bueno - Tablas de prueba y subtests
func TestFactory(t *testing.T) {
    tests := []struct {
        name        string
        provider    string
        expectError bool
    }{
        {"valid mimo provider", "mimo", false},
        {"invalid provider", "nonexistent", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            provider, err := Factory(tt.provider)
            if tt.expectError {
                if err == nil {
                    t.Error("expected error, got nil")
                }
            } else {
                if err != nil {
                    t.Errorf("unexpected error: %v", err)
                }
                if provider == nil {
                    t.Error("expected provider, got nil")
                }
            }
        })
    }
}

// ❌ Malo - Sin estructura
func TestFactory(t *testing.T) {
    p, _ := Factory("mimo")
    if p == nil {
        t.Fatal("failed")
    }
}
```

### Coverage

Buscamos mantener un coverage alto:
- Objetivo: >70% en código de negocio
- Ejecutar: `go test -cover ./...`

## 🔄 Pull Requests

### Antes de Crear un PR

1. **Tests**: Todos los tests deben pasar
2. **Linting**: Sin errores de linting
3. **Build**: El proyecto debe compilar sin errores
4. **Docs**: Actualiza la documentación si es necesario

### Creando un PR

1. Título claro y descriptivo
2. Describe los cambios en la descripción
3. Referencia issues relacionados (ej: `Closes #123`)
4. Agrega screenshots si es aplicable

### Plantilla de PR

```markdown
## Descripción
Breve descripción de los cambios.

## Tipo de Cambio
- [ ] Bug fix (non-breaking change)
- [ ] New feature (non-breaking change)
- [ ] Breaking change (fix/feature que rompe compatibilidad)

## Testing
- [ ] Tests unitarios incluidos/pasan
- [ ] Tests de integración incluidos/pasan
- [ ] Tests E2E incluidos/pasan

## Checklist
- [ ] Mi código sigue los estándares de estilo
- [ ] Me he documentado los cambios en docs/
- [ ] He actualizado README.md si es necesario
- [ ] Todos los tests pasan
- [ ] Sin errores de linting
```

### Revisión de PR

Los mantenedores revisarán el PR y pueden:
- Solicitar cambios
- Hacer preguntas
- Proponer mejoras
- Aprobar o rechazar el PR

## 🏗️ Arquitectura

### Estructura de Directorios

```
cclaude-glm/
├── cmd/                    # Aplicaciones CLI
│   └── cclaude/
├── internal/               # Código privado
│   ├── cli/               # CLI logic
│   ├── config/            # Configuración
│   ├── execution/         # Ejecución
│   └── provider/          # Providers
├── tests/                 # Tests adicionales
├── .pre-commit-config.yaml
├── go.mod
├── go.sum
└── README.md
```

### Patrones

- **Factory Pattern**: Creación de providers
- **Strategy Pattern**: Diferentes ejecutores
- **Interface Segregation**: Interfaces limpias y enfocadas

## 📚 Recursos de Aprendizaje

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Conventional Commits](https://www.conventionalcommits.org/)

## ❓ Preguntas Frecuentes

### ¿Necesito permiso para contribuir?

No! Solo sigue el proceso descrito arriba.

### ¿Puedo trabajar en cualquier issue?

Sí, pero es mejor comentar primero para evitar trabajo duplicado.

### ¿Qué hago si mi PR es rechazado?

No te preocupes, lee los comentarios y haz los cambios solicitados. Estamos aquí para ayudar!

### ¿Cómo puedo contactar a los mantenedores?

Abre un issue con la etiqueta `question`.

## 🙏 Gracias

Gracias por tu tiempo y esfuerzo en mejorar cclaude-glm!
