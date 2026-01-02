package flags

import (
	"flag"
	"os"
)

type FlagConfig struct {
	Provider        string
	Timeout         string
	Debug           bool
	Help            bool
	Version         bool
	ModelOverride   string
	ConfigFile      string
	Interactive     bool
	HealthCheck     bool
	Diagnose        bool
	ShowConfig      bool
	Confirm         bool
	Profile         string
	ListProfiles    bool
	CreateProfiles  bool
	Args            []string
}

func Parse() (*FlagConfig, error) {
	return ParseWithArgs(os.Args[1:])
}

func ParseWithArgs(args []string) (*FlagConfig, error) {
	var flags FlagConfig

	// Crear un nuevo FlagSet para evitar conflictos entre tests
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// Flags básicos
	fs.StringVar(&flags.Provider, "provider", "", "Proveedor de API (mimo, minimax, kimi, glm, claude)")
	fs.StringVar(&flags.Provider, "p", "", "Abreviatura para --provider")
	fs.StringVar(&flags.Timeout, "timeout", "5m", "Timeout para la ejecución")
	fs.BoolVar(&flags.Debug, "debug", false, "Modo debug")
	fs.BoolVar(&flags.Help, "help", false, "Mostrar ayuda")
	fs.BoolVar(&flags.Version, "version", false, "Mostrar versión")
	fs.StringVar(&flags.ModelOverride, "model", "", "Sobrescribir modelo por defecto")
	fs.StringVar(&flags.ConfigFile, "config", "", "Archivo de configuración personalizado")

	// Nuevos flags para funcionalidades extra
	fs.BoolVar(&flags.Interactive, "interactive", false, "Modo interactivo guiado")
	fs.BoolVar(&flags.Interactive, "i", false, "Abreviatura para --interactive")
	fs.BoolVar(&flags.HealthCheck, "health-check", false, "Verificar salud de proveedores")
	fs.BoolVar(&flags.HealthCheck, "hc", false, "Abreviatura para --health-check")
	fs.BoolVar(&flags.Diagnose, "diagnose", false, "Diagnóstico completo del sistema")
	fs.BoolVar(&flags.Diagnose, "d", false, "Abreviatura para --diagnose")
	fs.BoolVar(&flags.ShowConfig, "show-config", false, "Mostrar configuración actual")
	fs.BoolVar(&flags.ShowConfig, "sc", false, "Abreviatura para --show-config")
	fs.BoolVar(&flags.Confirm, "confirm", false, "Solicitar confirmación antes de ejecutar")
	fs.BoolVar(&flags.Confirm, "c", false, "Abreviatura para --confirm")
	fs.StringVar(&flags.Profile, "profile", "", "Usar perfil de configuración")
	fs.StringVar(&flags.Profile, "pr", "", "Abreviatura para --profile")
	fs.BoolVar(&flags.ListProfiles, "list-profiles", false, "Listar perfiles disponibles")
	fs.BoolVar(&flags.ListProfiles, "lp", false, "Abreviatura para --list-profiles")
	fs.BoolVar(&flags.CreateProfiles, "create-profiles", false, "Crear perfiles por defecto")
	fs.BoolVar(&flags.CreateProfiles, "cp", false, "Abreviatura para --create-profiles")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	// Si no hay proveedor y no son flags de ayuda, usar el primer argumento
	if flags.Provider == "" && fs.NArg() > 0 {
		flags.Provider = fs.Arg(0)
		// Capturar argumentos después del proveedor
		if fs.NArg() > 1 {
			flags.Args = fs.Args()[1:]
		}
	} else if flags.Provider != "" && fs.NArg() > 0 {
		// Si el proveedor fue especificado por flag, capturar todos los args
		flags.Args = fs.Args()
	} else {
		flags.Args = fs.Args()
	}

	return &flags, nil
}