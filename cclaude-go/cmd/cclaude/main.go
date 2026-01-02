package main

import (
	"fmt"
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

	// Modo interactivo
	if flagConfig.Interactive {
		im := utils.NewInteractiveMode()
		if err := im.Run(); err != nil {
			utils.HandleError(err, utils.ExitConfigError)
		}
		return
	}

	// Health check
	if flagConfig.HealthCheck {
		healthChecker := provider.NewHealthChecker()
		results := healthChecker.CheckAll()
		fmt.Println(healthChecker.FormatHealthResults(results, flagConfig.Debug))
		os.Exit(0)
	}

	// Diagnóstico completo
	if flagConfig.Diagnose {
		healthChecker := provider.NewHealthChecker()
		diagnostics := healthChecker.RunDiagnostics()
		fmt.Println("\n📊 Diagnóstico Completo:")
		for k, v := range diagnostics {
			fmt.Printf("  %s: %v\n", k, v)
		}
		os.Exit(0)
	}

	// Mostrar configuración actual
	if flagConfig.ShowConfig {
		utils.ShowConfig()
		utils.ShowTips()
		os.Exit(0)
	}

	// Listar perfiles
	if flagConfig.ListProfiles {
		pm, err := config.NewProfileManager()
		if err != nil {
			utils.HandleError(err, utils.ExitConfigError)
		}
		if err := pm.LoadProfiles(); err != nil {
			utils.HandleError(err, utils.ExitConfigError)
		}
		profiles := pm.ListProfiles()
		if len(profiles) == 0 {
			fmt.Println("No hay perfiles configurados.")
		} else {
			fmt.Println("📋 Perfiles disponibles:")
			for _, p := range profiles {
				fmt.Printf("  - %s\n", p)
			}
		}
		os.Exit(0)
	}

	// Crear perfiles por defecto
	if flagConfig.CreateProfiles {
		pm, err := config.NewProfileManager()
		if err != nil {
			utils.HandleError(err, utils.ExitConfigError)
		}
		if err := pm.CreateDefaultProfiles(); err != nil {
			utils.HandleError(err, utils.ExitConfigError)
		}
		fmt.Println("✅ Perfiles por defecto creados:")
		for _, p := range pm.ListProfiles() {
			fmt.Printf("  - %s\n", p)
		}
		os.Exit(0)
	}

	// Aplicar perfil si se especificó
	if flagConfig.Profile != "" {
		pm, err := config.NewProfileManager()
		if err != nil {
			utils.HandleError(err, utils.ExitConfigError)
		}
		if err := pm.LoadProfiles(); err != nil {
			utils.HandleError(err, utils.ExitConfigError)
		}
		profile, exists := pm.GetProfile(flagConfig.Profile)
		if !exists {
			utils.HandleError(fmt.Errorf("perfil '%s' no encontrado", flagConfig.Profile), utils.ExitConfigError)
		}
		if err := pm.ApplyProfile(profile); err != nil {
			utils.HandleError(err, utils.ExitConfigError)
		}
		utils.Info("Perfil '%s' aplicado", flagConfig.Profile)
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

	// Confirmación interactiva si está habilitada
	if flagConfig.Confirm && len(flagConfig.Args) > 0 {
		fmt.Printf("¿Ejecutar con %s? [s/n]: ", providerConfig.Name)
		var response string
		fmt.Scanln(&response)
		if response != "s" && response != "S" && response != "y" && response != "Y" {
			fmt.Println("❌ Operación cancelada")
			os.Exit(0)
		}
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