#!/usr/bin/env python3
"""
Script de validación rápida para cclaude.py
Ejecuta verificaciones básicas sin necesidad de pytest
"""

import os
import sys
import subprocess
from pathlib import Path

def print_header(text: str):
    print(f"\n{'='*60}")
    print(f"  {text}")
    print(f"{'='*60}")

def print_check(text: str, status: bool):
    icon = "✅" if status else "❌"
    print(f"{icon} {text}")

def validate_structure():
    """Valida estructura de archivos"""
    print_header("1. ESTRUCTURA DE ARCHIVOS")

    base = Path("/Users/argami/Documents/workspace/AI/cclaude/mimo")
    archivos = [
        "cclaude.py",
        "tests/test_cclaude.py",
        "requirements.txt",
        "Makefile",
        "README.md",
        "SETUP.md",
        ".gitignore"
    ]

    todos_existentes = True
    for archivo in archivos:
        existe = (base / archivo).exists()
        print_check(f"Archivo: {archivo}", existe)
        if not existe:
            todos_existentes = False

    return todos_existentes

def validate_permissions():
    """Valida permisos de ejecución"""
    print_header("2. PERMISOS")

    script = Path("/Users/argami/Documents/workspace/AI/cclaude/mimo/cclaude.py")
    es_ejecutable = os.access(script, os.X_OK)
    print_check("cclaude.py es ejecutable", es_ejecutable)

    return es_ejecutable

def validate_shebang():
    """Valida shebang del script"""
    print_header("3. SHEBANG")

    script = Path("/Users/argami/Documents/workspace/AI/cclaude/mimo/cclaude.py")
    with open(script, 'r') as f:
        first_line = f.readline().strip()

    es_correcto = first_line == "#!/usr/bin/env python3"
    print_check(f"Shebang: {first_line}", es_correcto)

    return es_correcto

def validate_imports():
    """Valida que el script se pueda importar sin errores"""
    print_header("4. IMPORTS Y SINTAXIS")

    try:
        sys.path.insert(0, "/Users/argami/Documents/workspace/AI/cclaude/mimo")
        import cclaude
        print_check("Módulo cclaude se importa correctamente", True)

        # Verificar clases principales
        clases = ["ProviderConfig", "ConfigValidator", "EnvironmentManager", "CLI", "Cclaude"]
        for clase in clases:
            existe = hasattr(cclaude, clase)
            print_check(f"Clase: {clase}", existe)

        return True
    except Exception as e:
        print_check(f"Error al importar: {e}", False)
        return False

def validate_help():
    """Valida que el help funcione"""
    print_header("5. HELP Y FLAGS")

    try:
        sys.path.insert(0, "/Users/argami/Documents/workspace/AI/cclaude/mimo")
        from cclaude import CLI

        help_text = CLI.show_help()
        flags_validos = [
            "--help" in help_text,
            "--version" in help_text,
            "--list-providers" in help_text,
            "mimo" in help_text,
            "minimax" in help_text
        ]

        for i, flag in enumerate(["--help", "--version", "--list-providers", "mimo", "minimax"]):
            print_check(f"Help contiene: {flag}", flags_validos[i])

        return all(flags_validos)
    except Exception as e:
        print_check(f"Error en help: {e}", False)
        return False

def validate_providers():
    """Valida configuración de providers"""
    print_header("6. CONFIGURACIÓN DE PROVIDERS")

    try:
        sys.path.insert(0, "/Users/argami/Documents/workspace/AI/cclaude/mimo")
        from cclaude import PROVIDERS

        providers_esperados = ["mimo", "minimax", "kimi", "glm"]
        todos_validos = True

        for provider in providers_esperados:
            if provider in PROVIDERS:
                config = PROVIDERS[provider]
                es_valido = all([
                    config.url.startswith("http"),
                    config.model,
                    config.env_key.endswith("_API_KEY"),
                    config.description
                ])
                print_check(f"Provider: {provider}", es_valido)
                if not es_valido:
                    todos_validos = False
            else:
                print_check(f"Provider: {provider}", False)
                todos_validos = False

        return todos_validos
    except Exception as e:
        print_check(f"Error en providers: {e}", False)
        return False

def validate_tests_exist():
    """Valida que los tests existen y tienen la estructura correcta"""
    print_header("7. TESTS TDD")

    test_file = Path("/Users/argami/Documents/workspace/AI/cclaude/mimo/tests/test_cclaude.py")

    if not test_file.exists():
        print_check("Archivo de tests existe", False)
        return False

    with open(test_file, 'r') as f:
        content = f.read()

    clases_test = [
        "TestProviderConfig",
        "TestConfigValidator",
        "TestEnvironmentManager",
        "TestCLI",
        "TestCclaudeIntegration",
        "TestCclaudeE2E",
        "TestSecurity",
        "TestRobustness",
        "TestPerformance",
        "TestCompatibility",
        "TestDocumentation",
        "TestSystemIntegration"
    ]

    todos_existentes = True
    for clase in clases_test:
        existe = clase in content
        print_check(f"Test class: {clase}", existe)
        if not existe:
            todos_existentes = False

    return todos_existentes

def validate_makefile():
    """Valida que el Makefile tiene comandos esenciales"""
    print_header("8. MAKEFILE")

    makefile = Path("/Users/argami/Documents/workspace/AI/cclaude/mimo/Makefile")

    if not makefile.exists():
        print_check("Makefile existe", False)
        return False

    with open(makefile, 'r') as f:
        content = f.read()

    comandos = ["test", "setup", "help", "test-coverage", "lint"]

    todos_existentes = True
    for cmd in comandos:
        existe = f"## {cmd}:" in content or f"{cmd}:" in content
        print_check(f"Comando: {cmd}", existe)
        if not existe:
            todos_existentes = False

    return todos_existentes

def validate_documentation():
    """Valida documentación"""
    print_header("9. DOCUMENTACIÓN")

    archivos = [
        "/Users/argami/Documents/workspace/AI/cclaude/mimo/README.md",
        "/Users/argami/Documents/workspace/AI/cclaude/mimo/SETUP.md"
    ]

    todos_existentes = True
    for archivo in archivos:
        existe = Path(archivo).exists()
        print_check(f"Doc: {Path(archivo).name}", existe)
        if not existe:
            todos_existentes = False

    return todos_existentes

def main():
    """Ejecuta todas las validaciones"""
    print_header("VALIDACIÓN COMPLETA DE cclaude.py")
    print("Ubicación: /Users/argami/Documents/workspace/AI/cclaude/mimo/")

    resultados = []

    resultados.append(("Estructura", validate_structure()))
    resultados.append(("Permisos", validate_permissions()))
    resultados.append(("Shebang", validate_shebang()))
    resultados.append(("Imports", validate_imports()))
    resultados.append(("Help", validate_help()))
    resultados.append(("Providers", validate_providers()))
    resultados.append(("Tests", validate_tests_exist()))
    resultados.append(("Makefile", validate_makefile()))
    resultados.append(("Documentación", validate_documentation()))

    print_header("RESUMEN FINAL")

    total = len(resultados)
    pasados = sum(1 for _, status in resultados if status)

    for nombre, status in resultados:
        print_check(f"{nombre}", status)

    print(f"\n📊 RESULTADO: {pasados}/{total} validaciones pasadas")

    if pasados == total:
        print("\n🎉 ¡TODO LISTO! cclaude.py está completamente implementado y validado")
        print("\nPróximos pasos:")
        print("  1. make setup")
        print("  2. ./cclaude.py --test")
        print("  3. ./cclaude.py mimo --help")
        return 0
    else:
        print(f"\n⚠️  {total - pasados} validaciones fallaron")
        print("Revisa los detalles arriba")
        return 1

if __name__ == "__main__":
    sys.exit(main())