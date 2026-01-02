package flags

import (
	"flag"
	"os"
)

type FlagConfig struct {
	Provider      string
	Timeout       string
	Debug         bool
	Help          bool
	Version       bool
	ModelOverride string
	ConfigFile    string
	Args          []string
}

func Parse() (*FlagConfig, error) {
	return ParseWithArgs(os.Args[1:])
}

func ParseWithArgs(args []string) (*FlagConfig, error) {
	var flags FlagConfig

	// Crear un nuevo FlagSet para evitar conflictos entre tests
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	fs.StringVar(&flags.Provider, "provider", "", "Proveedor de API (mimo, minimax, kimi, glm, claude)")
	fs.StringVar(&flags.Provider, "p", "", "Abreviatura para --provider")
	fs.StringVar(&flags.Timeout, "timeout", "5m", "Timeout para la ejecución")
	fs.BoolVar(&flags.Debug, "debug", false, "Modo debug")
	fs.BoolVar(&flags.Help, "help", false, "Mostrar ayuda")
	fs.BoolVar(&flags.Version, "version", false, "Mostrar versión")
	fs.StringVar(&flags.ModelOverride, "model", "", "Sobrescribir modelo por defecto")
	fs.StringVar(&flags.ConfigFile, "config", "", "Archivo de configuración personalizado")

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