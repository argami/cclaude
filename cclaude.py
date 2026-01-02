#!/usr/bin/env python3
"""
cclaude - Claude Code wrapper multi-provider con validación robusta
Uso: cclaude <provider> [claude args...]
Providers: mimo, minimax, kimi, glm, claude

Ejemplos:
    cclaude mimo --verbose
    cclaude minimax --help
    cclaude claude --version
"""

import os
import sys
from dataclasses import dataclass
from typing import Optional, Dict, List
import subprocess
from pathlib import Path


@dataclass
class ProviderConfig:
    """Configuración de provider con validación de seguridad"""
    url: str
    model: str
    env_key: str
    description: str

    def validate(self) -> Optional[str]:
        """Valida la configuración del provider"""
        if not self.url:
            return "URL no definida"
        if not self.model:
            return "Modelo no definido"
        if not self.env_key:
            return "Variable de entorno no definida"
        return None


# Configuración centralizada de providers
PROVIDERS: Dict[str, ProviderConfig] = {
    "mimo": ProviderConfig(
        url="https://api.xiaomimimo.com/anthropic",
        model="mimo-v2-flash",
        env_key="MIMO_API_KEY",
        description="Xiaomi MiMo V2 Flash"
    ),
    "minimax": ProviderConfig(
        url="https://api.minimax.io/anthropic",
        model="MiniMax-M2.1",
        env_key="MINIMAX_API_KEY",
        description="MiniMax M2.1"
    ),
    "kimi": ProviderConfig(
        url="https://api.kimi.com/coding/",
        model="kimi-k2-0711-preview",
        env_key="KIMI_API_KEY",
        description="Kimi K2 0711 Preview"
    ),
    "glm": ProviderConfig(
        url="https://api.z.ai/api/anthropic",
        model="glm-4.7",
        env_key="GLM_API_KEY",
        description="GLM 4.7"
    ),
}


class CclaudeError(Exception):
    """Excepción específica para errores de cclaude"""
    pass


class ConfigValidator:
    """Validador de configuración y entorno"""

    @staticmethod
    def validate_provider(provider: str) -> bool:
        """Valida que el provider exista"""
        return provider in PROVIDERS or provider == "claude"

    @staticmethod
    def validate_env_key(provider: str) -> Optional[str]:
        """Valida que la variable de entorno exista"""
        if provider == "claude":
            return None

        config = PROVIDERS[provider]
        api_key = os.getenv(config.env_key)

        if not api_key:
            return f"❌ Variable {config.env_key} no está definida"

        if len(api_key) < 10:
            return f"❌ {config.env_key} parece inválida (muy corta)"

        return None

    @staticmethod
    def validate_claude_available() -> bool:
        """Verifica que el comando 'claude' esté disponible"""
        try:
            subprocess.run(
                ["which", "claude"],
                capture_output=True,
                check=True,
                timeout=5
            )
            return True
        except (subprocess.CalledProcessError, FileNotFoundError):
            return False


class EnvironmentManager:
    """Gestiona las variables de entorno para cada provider"""

    @staticmethod
    def setup_provider_env(provider: str) -> None:
        """Configura el entorno para un provider específico"""
        if provider == "claude":
            return  # No necesita configuración especial

        config = PROVIDERS[provider]
        api_key = os.getenv(config.env_key)

        # Variables base
        os.environ["ANTHROPIC_BASE_URL"] = config.url
        os.environ["ANTHROPIC_MODEL"] = config.model
        os.environ["ANTHROPIC_AUTH_TOKEN"] = api_key
        os.environ["ANTHROPIC_API_KEY"] = ""  # Limpiar por seguridad

        # Variables para Claude Code
        os.environ["ANTHROPIC_DEFAULT_SONNET_MODEL"] = config.model
        os.environ["ANTHROPIC_DEFAULT_HAIKU_MODEL"] = config.model
        os.environ["CLAUDE_CODE_SUBAGENT_MODEL"] = config.model

        # Optimización de tráfico
        os.environ["DISABLE_NON_ESSENTIAL_MODEL_CALLS"] = "1"
        os.environ["CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC"] = "1"

        # Timeout extendido
        os.environ["API_TIMEOUT_MS"] = "3000000"

    @staticmethod
    def get_current_config(provider: str) -> Dict[str, str]:
        """Obtiene la configuración actual para logging"""
        if provider == "claude":
            return {"modo": "claude nativo"}

        config = PROVIDERS[provider]
        return {
            "provider": provider,
            "model": config.model,
            "url": config.url,
            "env_key": config.env_key,
            "description": config.description
        }


class CLI:
    """Maneja la interfaz de línea de comandos"""

    @staticmethod
    def show_help() -> str:
        """Muestra ayuda detallada"""
        providers_list = ", ".join(PROVIDERS.keys())
        return f"""
cclaude - Claude Code wrapper multi-provider

Uso:
    cclaude <provider> [args...]

Providers disponibles:
    {providers_list}
    claude (modo nativo sin configuración)

Flags:
    --help, -h          Mostrar esta ayuda
    --version, -v       Mostrar versión
    --list-providers    Listar providers con detalles
    --test              Ejecutar tests de validación

Ejemplos:
    cclaude mimo --verbose
    cclaude minimax --help
    cclaude claude --version

Variables de entorno requeridas:
    MIMO_API_KEY, MINIMAX_API_KEY, KIMI_API_KEY, GLM_API_KEY
"""

    @staticmethod
    def show_version() -> str:
        """Muestra versión del script"""
        return "cclaude.py v1.0.0 - Python wrapper para Claude Code"

    @staticmethod
    def show_providers() -> str:
        """Muestra detalles de todos los providers"""
        output = ["Providers disponibles:\n"]
        for name, config in PROVIDERS.items():
            output.append(f"  {name}:")
            output.append(f"    Descripción: {config.description}")
            output.append(f"    Modelo: {config.model}")
            output.append(f"    URL: {config.url}")
            output.append(f"    Variable: {config.env_key}")
            output.append("")
        return "\n".join(output)

    @staticmethod
    def print_config(config: Dict[str, str]) -> None:
        """Imprime configuración actual con formato bonito"""
        print("\n" + "="*50)
        print("🎯 CONFIGURACIÓN ACTIVADA")
        print("="*50)
        for key, value in config.items():
            print(f"  {key:<15}: {value}")
        print("="*50 + "\n")


class Cclaude:
    """Clase principal de cclaude"""

    def __init__(self):
        self.validator = ConfigValidator()
        self.env_manager = EnvironmentManager()
        self.cli = CLI()

    def run(self, args: List[str]) -> int:
        """Ejecuta cclaude con los argumentos proporcionados"""
        try:
            # Parsear flags sin provider
            if not args:
                print(self.cli.show_help())
                return 0

            if args[0] in ["--help", "-h"]:
                print(self.cli.show_help())
                return 0

            if args[0] in ["--version", "-v"]:
                print(self.cli.show_version())
                return 0

            if args[0] in ["--list-providers", "-l"]:
                print(self.cli.show_providers())
                return 0

            if args[0] == "--test":
                return self.run_tests()

            # Extraer provider
            provider = args[0]
            claude_args = args[1:]

            # Validar provider
            if not self.validator.validate_provider(provider):
                print(f"❌ Provider inválido: {provider}")
                print(f"✅ Providers disponibles: {', '.join(PROVIDERS.keys())}, claude")
                return 1

            # Validar entorno (solo para providers alternativos)
            env_error = self.validator.validate_env_key(provider)
            if env_error:
                print(env_error)
                return 1

            # Verificar claude disponible
            if not self.validator.validate_claude_available():
                print("❌ Comando 'claude' no encontrado en PATH")
                print("   Asegúrate de tener Claude Code instalado")
                return 1

            # Configurar entorno
            if provider != "claude":
                self.env_manager.setup_provider_env(provider)
                config = self.env_manager.get_current_config(provider)
                self.cli.print_config(config)

            # Ejecutar claude
            print(f"🚀 Ejecutando claude con provider: {provider}")
            print(f"   Args: {' '.join(claude_args) if claude_args else '(ninguno)'}")
            print()

            os.execvp("claude", ["claude"] + claude_args)

        except KeyboardInterrupt:
            print("\n⚠️ Operación cancelada por usuario")
            return 130
        except Exception as e:
            print(f"❌ Error inesperado: {e}")
            return 1

    def run_tests(self) -> int:
        """Ejecuta tests de validación del sistema"""
        print("🧪 Ejecutando tests de validación...\n")

        tests_passed = 0
        tests_failed = 0

        # Test 1: Validación de providers
        print("1. Validación de providers...")
        for provider in PROVIDERS.keys():
            if self.validator.validate_provider(provider):
                print(f"   ✅ {provider}")
                tests_passed += 1
            else:
                print(f"   ❌ {provider}")
                tests_failed += 1

        # Test 2: Validación de claude nativo
        print("\n2. Validación claude nativo...")
        if self.validator.validate_provider("claude"):
            print("   ✅ claude nativo")
            tests_passed += 1
        else:
            print("   ❌ claude nativo")
            tests_failed += 1

        # Test 3: Validación de variables de entorno
        print("\n3. Validación variables de entorno...")
        for provider in PROVIDERS.keys():
            error = self.validator.validate_env_key(provider)
            if error:
                print(f"   ⚠️  {provider}: {error}")
            else:
                print(f"   ✅ {provider}: API key presente")
                tests_passed += 1

        # Test 4: Verificar claude disponible
        print("\n4. Verificar claude en PATH...")
        if self.validator.validate_claude_available():
            print("   ✅ claude disponible")
            tests_passed += 1
        else:
            print("   ⚠️  claude no encontrado (no es crítico para tests)")

        # Resumen
        print("\n" + "="*40)
        print("RESUMEN DE TESTS")
        print("="*40)
        print(f"✅ Pasados: {tests_passed}")
        print(f"❌ Fallidos: {tests_failed}")
        print(f"📊 Total: {tests_passed + tests_failed}")

        return 0 if tests_failed == 0 else 1


def main():
    """Punto de entrada principal"""
    cclaude = Cclaude()
    exit_code = cclaude.run(sys.argv[1:])
    sys.exit(exit_code)


if __name__ == "__main__":
    main()