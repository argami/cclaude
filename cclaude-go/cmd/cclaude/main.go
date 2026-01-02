package main

import (
	"os"
	"time"

	"github.com/argami/cclaude-go/internal/config"
	"github.com/argami/cclaude-go/internal/flags"
	"github.com/argami/cclaude-go/internal/provider"
	"github.com/argami/cclaude-go/internal/utils"
)

func main() {
	// Parsear flags
	flagConfig, err := flags.Parse()
	if err != nil {
		utils.HandleError(err, utils.ExitConfigError)
	}

	// Manejar flags de ayuda y versión
	if flagConfig.Help {
		utils.ShowHelp()
		os.Exit(0)
	}

	if flagConfig.Version {
		utils.ShowVersion()
		os.Exit(0)
	}

	// Validar ambiente
	if err := config.ValidateEnvironment(); err != nil {
		utils.HandleError(err, utils.ExitValidationError)
	}

	// Obtener proveedor
	providerConfig, err := provider.GetProvider(flagConfig.Provider)
	if err != nil {
		// Si no hay proveedor o es "claude", ejecutar nativo
		if flagConfig.Provider == "" || flagConfig.Provider == "claude" {
			utils.Info("Ejecutando Claude nativo")
			args := flagConfig.Args
			if len(args) == 0 {
				args = flagConfig.Args
			}
			if err := utils.ExecuteClaude(args); err != nil {
				utils.HandleError(err, utils.ExitClaudeNotFound)
			}
			return
		}
		utils.HandleError(err, utils.ExitProviderNotFound)
	}

	// Validar API key
	if err := config.ValidateAPIKey(*providerConfig); err != nil {
		utils.HandleError(err, utils.ExitAPIKeyMissing)
	}

	// Configurar timeout
	timeout, err := time.ParseDuration(flagConfig.Timeout)
	if err != nil {
		timeout = 5 * time.Minute
	}

	// Configurar variables de entorno
	authToken := os.Getenv(providerConfig.EnvVar)
	if err := utils.SetupEnvironment(*providerConfig, authToken, flagConfig.ModelOverride); err != nil {
		utils.HandleError(err, utils.ExitConfigError)
	}

	// Logging de configuración
	utils.Info("Proveedor: %s", providerConfig.Name)
	utils.Info("Modelo: %s", providerConfig.Model)
	utils.Info("Timeout: %s", timeout)

	if flagConfig.Debug {
		utils.SetLogLevel(utils.LevelDebug)
		utils.Debug("Modo debug habilitado")
		utils.Debug("Base URL: %s", providerConfig.BaseURL)
	}

	// Ejecutar claude con argumentos restantes
	claudeArgs := flagConfig.Args
	if len(claudeArgs) == 0 {
		claudeArgs = flagConfig.Args
	}

	if err := utils.ExecuteClaude(claudeArgs); err != nil {
		utils.HandleError(err, utils.ExitClaudeNotFound)
	}
}