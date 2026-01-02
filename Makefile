# Makefile para cclaude.py
# Cumple con tu regla de "herramientas de automatización en scripts/"

.PHONY: help install test test-coverage lint format check setup clean

# Variables
PYTHON := python3
PYTEST := $(PYTHON) -m pytest
BLACK := $(PYTHON) -m black
RUFF := $(PYTHON) -m ruff

## help: Muestra esta ayuda
help:
	@echo "Uso: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

## install: Instala dependencias de desarrollo
install:
	@echo "📦 Instalando dependencias..."
	$(PYTHON) -m pip install --upgrade pip
	$(PYTHON) -m pip install -r requirements.txt

## test: Ejecuta todos los tests
test:
	@echo "🧪 Ejecutando tests..."
	$(PYTEST) tests/test_cclaude.py -v

## test-coverage: Ejecuta tests con cobertura
test-coverage:
	@echo "📊 Ejecutando tests con cobertura..."
	$(PYTEST) tests/test_cclaude.py --cov=cclaude --cov-report=html --cov-report=term-missing

## lint: Verifica calidad del código
lint:
	@echo "🔍 Verificando calidad del código..."
	$(RUFF) check cclaude.py tests/test_cclaude.py

## format: Formatea el código
format:
	@echo "✨ Formateando código..."
	$(BLACK) cclaude.py tests/test_cclaude.py

## check: Lint + format check
check:
	@echo "✅ Verificando todo..."
	$(RUFF) check cclaude.py tests/test_cclaude.py
	$(BLACK) --check cclaude.py tests/test_cclaude.py

## setup: Configura todo para desarrollo
setup:
	@echo "🚀 Configurando entorno de desarrollo..."
	chmod +x cclaude.py
	$(PYTHON) -m pip install --upgrade pip
	$(PYTHON) -m pip install -r requirements.txt
	@echo "✅ Listo! Puedes usar: ./cclaude.py mimo --help"

## clean: Limpia archivos temporales
clean:
	@echo "🧹 Limpiando..."
	rm -rf .pytest_cache/ __pycache__/ tests/__pycache__/ htmlcov/ .coverage
	rm -rf *.pyc tests/*.pyc

## run-mimo: Ejemplo rápido con mimo
run-mimo:
	@echo "🚀 Ejecutando con mimo..."
	./cclaude.py mimo --help

## run-claude: Ejemplo rápido con claude nativo
run-claude:
	@echo "🚀 Ejecutando con claude nativo..."
	./cclaude.py claude --version

## test-all: Suite completa de validación
test-all: check test test-coverage
	@echo "✅ Suite de validación completada!"

## docs: Genera documentación (si se añade sphinx)
docs:
	@echo "📚 Generando documentación..."
	@echo "Nota: Ejecutar manualmente si se añade sphinx"

# Mensaje de bienvenida
welcome:
	@echo "🌟 Bienvenido a cclaude.py"
	@echo ""
	@echo "Primeros pasos:"
	@echo "  1. make setup          # Configurar todo"
	@echo "  2. make test-all       # Ejecutar suite completa"
	@echo "  3. ./cclaude.py --help # Probar el CLI"
	@echo ""
	@echo "Uso diario:"
	@echo "  ./cclaude.py mimo --verbose"
	@echo "  ./cclaude.py minimax --help"
	@echo "  ./cclaude.py claude --version"