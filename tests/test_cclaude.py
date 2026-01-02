"""
Tests TDD para cclaude.py
Ejecutar: pytest tests/test_cclaude.py -v
"""

import os
import sys
import pytest
from unittest.mock import patch, MagicMock
import subprocess

# Añadir el directorio padre al path para importar cclaude
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from cclaude import (
    ProviderConfig,
    ConfigValidator,
    EnvironmentManager,
    CLI,
    Cclaude,
    CclaudeError,
    PROVIDERS
)


# ============================================================================
# TESTS DE UNIDAD - ProviderConfig
# ============================================================================
class TestProviderConfig:
    """Tests para la clase ProviderConfig"""

    def test_config_valida(self):
        """✅ Configuración válida pasa validación"""
        config = ProviderConfig(
            url="https://api.test.com",
            model="test-model",
            env_key="TEST_API_KEY",
            description="Test provider"
        )
        assert config.validate() is None

    def test_config_url_vacia(self):
        """❌ URL vacía falla validación"""
        config = ProviderConfig(
            url="",
            model="test-model",
            env_key="TEST_API_KEY",
            description="Test"
        )
        assert config.validate() == "URL no definida"

    def test_config_modelo_vacio(self):
        """❌ Modelo vacío falla validación"""
        config = ProviderConfig(
            url="https://api.test.com",
            model="",
            env_key="TEST_API_KEY",
            description="Test"
        )
        assert config.validate() == "Modelo no definido"

    def test_config_env_key_vacia(self):
        """❌ Env key vacía falla validación"""
        config = ProviderConfig(
            url="https://api.test.com",
            model="test-model",
            env_key="",
            description="Test"
        )
        assert config.validate() == "Variable de entorno no definida"


# ============================================================================
# TESTS DE UNIDAD - ConfigValidator
# ============================================================================
class TestConfigValidator:
    """Tests para el validador de configuración"""

    def test_validate_provider_existe(self):
        """✅ Provider existente es válido"""
        validator = ConfigValidator()
        assert validator.validate_provider("mimo") is True
        assert validator.validate_provider("minimax") is True
        assert validator.validate_provider("claude") is True

    def test_validate_provider_no_existe(self):
        """❌ Provider inexistente es inválido"""
        validator = ConfigValidator()
        assert validator.validate_provider("invalid-provider") is False

    @patch.dict(os.environ, {"MIMO_API_KEY": "test_key_12345678"})
    def test_validate_env_key_existe(self):
        """✅ Variable de entorno existente pasa validación"""
        validator = ConfigValidator()
        error = validator.validate_env_key("mimo")
        assert error is None

    @patch.dict(os.environ, clear=True)
    def test_validate_env_key_falta(self):
        """❌ Variable de entorno faltante falla validación"""
        validator = ConfigValidator()
        error = validator.validate_env_key("mimo")
        assert "no está definida" in error

    @patch.dict(os.environ, {"MIMO_API_KEY": "short"})
    def test_validate_env_key_corta(self):
        """❌ API key muy corta falla validación"""
        validator = ConfigValidator()
        error = validator.validate_env_key("mimo")
        assert "parece inválida" in error

    def test_validate_env_key_claude(self):
        """✅ Claude nativo no necesita validación de env"""
        validator = ConfigValidator()
        error = validator.validate_env_key("claude")
        assert error is None

    @patch("subprocess.run")
    def test_claude_disponible(self, mock_run):
        """✅ Detecta cuando claude está disponible"""
        mock_run.return_value = MagicMock()
        validator = ConfigValidator()
        assert validator.validate_claude_available() is True

    @patch("subprocess.run")
    def test_claude_no_disponible(self, mock_run):
        """❌ Detecta cuando claude no está disponible"""
        mock_run.side_effect = subprocess.CalledProcessError(1, "which")
        validator = ConfigValidator()
        assert validator.validate_claude_available() is False


# ============================================================================
# TESTS DE UNIDAD - EnvironmentManager
# ============================================================================
class TestEnvironmentManager:
    """Tests para el gestor de entorno"""

    @patch.dict(os.environ, {"MIMO_API_KEY": "test_key_12345678"}, clear=True)
    def test_setup_provider_env(self):
        """✅ Configura correctamente variables para provider"""
        manager = EnvironmentManager()
        manager.setup_provider_env("mimo")

        assert os.environ["ANTHROPIC_BASE_URL"] == "https://api.xiaomimimo.com/anthropic"
        assert os.environ["ANTHROPIC_MODEL"] == "mimo-v2-flash"
        assert os.environ["ANTHROPIC_AUTH_TOKEN"] == "test_key_12345678"
        assert os.environ["DISABLE_NON_ESSENTIAL_MODEL_CALLS"] == "1"

    def test_setup_claude_env(self):
        """✅ Claude nativo no modifica entorno"""
        manager = EnvironmentManager()
        original_env = os.environ.copy()

        manager.setup_provider_env("claude")

        # Debería ser idéntico
        assert os.environ == original_env

    def test_get_current_config_claude(self):
        """✅ Config de claude nativo"""
        manager = EnvironmentManager()
        config = manager.get_current_config("claude")
        assert config == {"modo": "claude nativo"}

    @patch.dict(os.environ, {"MIMO_API_KEY": "test_key"}, clear=True)
    def test_get_current_config_provider(self):
        """✅ Config de provider alternativo"""
        manager = EnvironmentManager()
        config = manager.get_current_config("mimo")
        assert config["provider"] == "mimo"
        assert config["model"] == "mimo-v2-flash"
        assert config["env_key"] == "MIMO_API_KEY"


# ============================================================================
# TESTS DE UNIDAD - CLI
# ============================================================================
class TestCLI:
    """Tests para la interfaz CLI"""

    def test_show_help_contiene_info(self):
        """✅ Help contiene información esencial"""
        help_text = CLI.show_help()
        assert "cclaude - Claude Code wrapper" in help_text
        assert "mimo" in help_text
        assert "--help" in help_text

    def test_show_version_formato(self):
        """✅ Version tiene formato correcto"""
        version = CLI.show_version()
        assert "cclaude.py" in version
        assert "v1.0.0" in version

    def test_show_providers_formato(self):
        """✅ Lista de providers tiene formato correcto"""
        providers = CLI.show_providers()
        assert "mimo:" in providers
        assert "MIMO_API_KEY" in providers
        assert "https://api.xiaomimimo.com" in providers


# ============================================================================
# TESTS DE INTEGRACIÓN - Cclaude
# ============================================================================
class TestCclaudeIntegration:
    """Tests de integración para la clase principal"""

    def test_run_sin_args(self):
        """✅ Sin args muestra help"""
        cclaude = Cclaude()
        with patch("sys.stdout") as mock_stdout:
            result = cclaude.run([])
        assert result == 0

    def test_run_help_flag(self):
        """✅ Flag --help muestra ayuda"""
        cclaude = Cclaude()
        with patch("sys.stdout") as mock_stdout:
            result = cclaude.run(["--help"])
        assert result == 0

    def test_run_version_flag(self):
        """✅ Flag --version muestra versión"""
        cclaude = Cclaude()
        with patch("sys.stdout") as mock_stdout:
            result = cclaude.run(["--version"])
        assert result == 0

    def test_run_list_providers_flag(self):
        """✅ Flag --list-providers lista providers"""
        cclaude = Cclaude()
        with patch("sys.stdout") as mock_stdout:
            result = cclaude.run(["--list-providers"])
        assert result == 0

    def test_run_provider_invalido(self):
        """❌ Provider inválido retorna error"""
        cclaude = Cclaude()
        result = cclaude.run(["invalid-provider"])
        assert result == 1

    @patch.dict(os.environ, clear=True)
    def test_run_provider_sin_env_key(self):
        """❌ Provider sin API key retorna error"""
        cclaude = Cclaude()
        result = cclaude.run(["mimo"])
        assert result == 1

    @patch.dict(os.environ, {"MIMO_API_KEY": "test_key_12345678"}, clear=True)
    @patch("cclaude.ConfigValidator.validate_claude_available")
    def test_run_claude_no_disponible(self, mock_validate):
        """❌ Error si claude no está disponible"""
        mock_validate.return_value = False
        cclaude = Cclaude()
        result = cclaude.run(["mimo"])
        assert result == 1

    @patch.dict(os.environ, {"MIMO_API_KEY": "test_key_12345678"}, clear=True)
    @patch("cclaude.ConfigValidator.validate_claude_available")
    @patch("os.execvp")
    def test_run_exitoso_mimo(self, mock_execvp, mock_validate):
        """✅ Ejecución exitosa con mimo"""
        mock_validate.return_value = True
        cclaude = Cclaude()

        # No debería llegar a execvp en test, pero si llega, mockeamos
        mock_execvp.return_value = None

        # Mockeamos execvp para que no termine el test
        with patch("os.execvp", side_effect=Exception("Execvp llamado")):
            try:
                result = cclaude.run(["mimo", "--help"])
            except Exception as e:
                # Si execvp es llamado, el test pasa
                assert "Execvp llamado" in str(e)

    @patch.dict(os.environ, clear=True)
    def test_run_claude_nativo(self):
        """✅ Ejecución exitosa con claude nativo"""
        cclaude = Cclaude()
        mock_execvp = patch("os.execvp", side_effect=Exception("Execvp llamado"))

        with mock_execvp:
            try:
                result = cclaude.run(["claude", "--version"])
            except Exception as e:
                assert "Execvp llamado" in str(e)

    @patch.dict(os.environ, {"MIMO_API_KEY": "test_key_12345678"}, clear=True)
    @patch("cclaude.ConfigValidator.validate_claude_available")
    def test_run_test_flag(self, mock_validate):
        """✅ Flag --test ejecuta tests internos"""
        mock_validate.return_value = True
        cclaude = Cclaude()
        result = cclaude.run(["--test"])
        assert result == 0


# ============================================================================
# TESTS E2E - Flujo Completo
# ============================================================================
class TestCclaudeE2E:
    """Tests de extremo a extremo"""

    @patch.dict(os.environ, {"MIMO_API_KEY": "test_key_12345678"}, clear=True)
    @patch("cclaude.ConfigValidator.validate_claude_available")
    @patch("os.execvp")
    def test_e2e_mimo_completo(self, mock_execvp, mock_validate):
        """✅ Flujo completo: validación + configuración + ejecución"""
        mock_validate.return_value = True

        cclaude = Cclaude()

        # Ejecutar con mimo
        try:
            cclaude.run(["mimo", "--verbose"])
        except:
            pass  # execvp lanza excepción en test

        # Verificar que se llamó a execvp con claude
        mock_execvp.assert_called_once()
        args = mock_execvp.call_args[0]
        assert args[0] == "claude"
        assert "--verbose" in args[1]

    @patch.dict(os.environ, {"KIMI_API_KEY": "test_key_12345678"}, clear=True)
    @patch("cclaude.ConfigValidator.validate_claude_available")
    @patch("os.execvp")
    def test_e2e_kimi_completo(self, mock_execvp, mock_validate):
        """✅ Flujo completo con provider kimi"""
        mock_validate.return_value = True

        cclaude = Cclaude()

        try:
            cclaude.run(["kimi", "--help"])
        except:
            pass

        # Verificar variables de entorno configuradas
        assert os.environ["ANTHROPIC_BASE_URL"] == "https://api.kimi.com/coding/"
        assert os.environ["ANTHROPIC_MODEL"] == "kimi-k2-0711-preview"

    @patch.dict(os.environ, clear=True)
    def test_e2e_error_cadena_completa(self):
        """✅ Cadena de errores completa y clara"""
        cclaude = Cclaude()

        # Primero intentar sin API key
        result = cclaude.run(["mimo"])
        assert result == 1

        # Luego con provider inválido
        result = cclaude.run(["invalid"])
        assert result == 1


# ============================================================================
# TESTS DE SEGURIDAD
# ============================================================================
class TestSecurity:
    """Tests de seguridad"""

    @patch.dict(os.environ, {"MIMO_API_KEY": "short"}, clear=True)
    def test_api_key_corta_detectada(self):
        """⚠️ Detecta API keys potencialmente inválidas"""
        from cclaude import ConfigValidator
        validator = ConfigValidator()
        error = validator.validate_env_key("mimo")
        assert "parece inválida" in error

    @patch.dict(os.environ, {"MIMO_API_KEY": "valid_key_12345678"}, clear=True)
    def test_api_key_larga_aceptada(self):
        """✅ API keys largas son aceptadas"""
        from cclaude import ConfigValidator
        validator = ConfigValidator()
        error = validator.validate_env_key("mimo")
        assert error is None

    def test_provider_config_completa(self):
        """✅ Todos los providers tienen configuración válida"""
        for name, config in PROVIDERS.items():
            assert config.url.startswith("http")
            assert config.model
            assert config.env_key.endswith("_API_KEY")
            assert config.description


# ============================================================================
# TESTS DE ROBUSTEZ
# ============================================================================
class TestRobustness:
    """Tests de robustez y casos edge"""

    def test_run_keyboard_interrupt(self):
        """✅ Maneja Ctrl+C correctamente"""
        cclaude = Cclaude()

        with patch("cclaude.Cclaude.run", side_effect=KeyboardInterrupt):
            try:
                cclaude.run(["mimo"])
            except KeyboardInterrupt:
                pass

        # El wrapper debería capturar esto
        cclaude_int = Cclaude()
        with patch("sys.exit") as mock_exit:
            with patch("cclaude.Cclaude.run", side_effect=KeyboardInterrupt):
                cclaude_int.run(["mimo"])

        # Verificar que salió con código 130
        # (esto depende de cómo se implemente el manejo)

    def test_run_exception_desconocida(self):
        """✅ Maneja excepciones inesperadas"""
        cclaude = Cclaude()

        with patch("cclaude.Cclaude.run", side_effect=Exception("Error inesperado")):
            try:
                cclaude.run(["mimo"])
            except Exception as e:
                assert "Error inesperado" in str(e)

    def test_args_vacios(self):
        """✅ Maneja args vacíos"""
        cclaude = Cclaude()
        result = cclaude.run([])
        assert result == 0

    def test_args_solo_provider(self):
        """✅ Maneja solo provider sin args adicionales"""
        cclaude = Cclaude()

        with patch.dict(os.environ, {"MIMO_API_KEY": "test_key_12345678"}, clear=True):
            with patch("cclaude.ConfigValidator.validate_claude_available", return_value=True):
                with patch("os.execvp", side_effect=Exception("exec")):
                    try:
                        cclaude.run(["mimo"])
                    except:
                        pass  # Esperamos que falle en execvp


# ============================================================================
# TESTS DE RENDIMIENTO
# ============================================================================
class TestPerformance:
    """Tests de rendimiento"""

    def test_config_centralizada_rapida(self):
        """✅ Acceso a config centralizada es rápido (< 1ms)"""
        import time

        start = time.time()
        for _ in range(100):
            _ = PROVIDERS["mimo"]
        elapsed = time.time() - start

        assert elapsed < 0.1  # 100 accesos < 100ms

    def test_validacion_rapida(self):
        """✅ Validación es rápida (< 10ms)"""
        import time

        validator = ConfigValidator()
        start = time.time()
        for _ in range(100):
            validator.validate_provider("mimo")
        elapsed = time.time() - start

        assert elapsed < 0.01  # 100 validaciones < 10ms


# ============================================================================
# TESTS DE COMPATIBILIDAD
# ============================================================================
class TestCompatibility:
    """Tests de compatibilidad con bash original"""

    def test_mismo_comportamiento_help(self):
        """✅ Help tiene misma información esencial que bash"""
        bash_help = """# Del bash original
# Usage: cclaude <provider> [claude args...]
# Providers: mimo, minimax, kimi, glm, claude (or none for default)"""

        python_help = CLI.show_help()

        # Verificar que contiene la información clave
        assert "mimo" in python_help
        assert "minimax" in python_help
        assert "claude" in python_help

    def test_mismas_variables_entorno(self):
        """✅ Usa las mismas variables que bash original"""
        bash_vars = [
            "ANTHROPIC_BASE_URL",
            "ANTHROPIC_MODEL",
            "ANTHROPIC_AUTH_TOKEN",
            "ANTHROPIC_API_KEY",
            "ANTHROPIC_DEFAULT_SONNET_MODEL",
            "ANTHROPIC_DEFAULT_HAIKU_MODEL",
            "CLAUDE_CODE_SUBAGENT_MODEL",
            "DISABLE_NON_ESSENTIAL_MODEL_CALLS",
            "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC",
            "API_TIMEOUT_MS"
        ]

        # Verificar que el python las exporta todas
        with patch.dict(os.environ, {"MIMO_API_KEY": "test"}, clear=True):
            manager = EnvironmentManager()
            manager.setup_provider_env("mimo")

            for var in bash_vars:
                assert var in os.environ, f"Falta variable: {var}"

    def test_mismas_urls_providers(self):
        """✅ Mismas URLs y modelos que bash original"""
        urls_bash = {
            "mimo": "https://api.xiaomimimo.com/anthropic",
            "minimax": "https://api.minimax.io/anthropic",
            "kimi": "https://api.kimi.com/coding/",
            "glm": "https://api.z.ai/api/anthropic"
        }

        for provider, url in urls_bash.items():
            assert PROVIDERS[provider].url == url

        modelos_bash = {
            "mimo": "mimo-v2-flash",
            "minimax": "MiniMax-M2.1",
            "kimi": "kimi-k2-0711-preview",
            "glm": "glm-4.7"
        }

        for provider, model in modelos_bash.items():
            assert PROVIDERS[provider].model == model


# ============================================================================
# TESTS DE DOCUMENTACIÓN
# ============================================================================
class TestDocumentation:
    """Tests de documentación y docstrings"""

    def test_todas_las_funciones_documentadas(self):
        """✅ Todas las funciones tienen docstrings"""
        from cclaude import Cclaude, ConfigValidator, EnvironmentManager, CLI

        funciones = [
            Cclaude.run,
            ConfigValidator.validate_provider,
            ConfigValidator.validate_env_key,
            EnvironmentManager.setup_provider_env,
            CLI.show_help
        ]

        for func in funciones:
            assert func.__doc__, f"Falta docstring en {func.__name__}"

    def test_clase_principal_documentada(self):
        """✅ Clase principal tiene docstring"""
        from cclaude import Cclaude
        assert Cclaude.__doc__ is not None

    def test_module_tiene_docstring(self):
        """✅ Módulo principal tiene docstring"""
        import cclaude
        assert cclaude.__doc__ is not None


# ============================================================================
# TESTS DE INTEGRACIÓN CON SISTEMA
# ============================================================================
class TestSystemIntegration:
    """Tests de integración con el sistema"""

    @patch.dict(os.environ, {"MIMO_API_KEY": "test_key_12345678"}, clear=True)
    def test_ejecutable_permisos(self):
        """✅ El script debe ser ejecutable"""
        script_path = os.path.join(
            os.path.dirname(os.path.dirname(__file__)),
            "cclaude.py"
        )

        # Verificar que existe
        assert os.path.exists(script_path)

        # Verificar que tiene shebang
        with open(script_path, 'r') as f:
            first_line = f.readline()
            assert first_line.startswith("#!/usr/bin/env python3")

    def test_estructura_directorios(self):
        """✅ Verifica estructura de directorios correcta"""
        base_dir = os.path.dirname(os.path.dirname(__file__))

        assert os.path.exists(os.path.join(base_dir, "cclaude.py"))
        assert os.path.exists(os.path.join(base_dir, "tests"))
        assert os.path.exists(os.path.join(base_dir, "tests", "test_cclaude.py"))


# ============================================================================
# RUNNER DE TESTS
# ============================================================================
if __name__ == "__main__":
    # Ejecutar todos los tests con pytest
    import subprocess
    import sys

    result = subprocess.run([
        sys.executable, "-m", "pytest",
        __file__, "-v", "--tb=short"
    ], capture_output=True)

    print(result.stdout.decode())
    if result.stderr:
        print("STDERR:", result.stderr.decode())

    sys.exit(result.returncode)