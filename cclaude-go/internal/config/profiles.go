package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Profile representa un perfil de configuración
type Profile struct {
	Name        string
	Provider    string
	Model       string
	Timeout     string
	Environment map[string]string
}

// ProfileManager gestiona perfiles de configuración
type ProfileManager struct {
	configDir string
	profiles  map[string]Profile
}

// NewProfileManager crea un nuevo gestor de perfiles
func NewProfileManager() (*ProfileManager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener directorio home: %w", err)
	}

	configDir := filepath.Join(home, ".config", "cclaude")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("no se pudo crear directorio config: %w", err)
	}

	return &ProfileManager{
		configDir: configDir,
		profiles:  make(map[string]Profile),
	}, nil
}

// LoadProfiles carga todos los perfiles desde el directorio de configuración
func (pm *ProfileManager) LoadProfiles() error {
	profileDir := filepath.Join(pm.configDir, "profiles")
	if _, err := os.Stat(profileDir); os.IsNotExist(err) {
		return nil // No hay perfiles, no es un error
	}

	entries, err := os.ReadDir(profileDir)
	if err != nil {
		return fmt.Errorf("no se pudo leer directorio perfiles: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".conf") {
			continue
		}

		profileName := strings.TrimSuffix(entry.Name(), ".conf")
		profile, err := pm.loadProfileFile(filepath.Join(profileDir, entry.Name()))
		if err != nil {
			return fmt.Errorf("no se pudo cargar perfil %s: %w", profileName, err)
		}

		pm.profiles[profileName] = profile
	}

	return nil
}

// loadProfileFile lee y parsea un archivo de perfil
func (pm *ProfileManager) loadProfileFile(path string) (Profile, error) {
	file, err := os.Open(path)
	if err != nil {
		return Profile{}, err
	}
	defer file.Close()

	profile := Profile{
		Environment: make(map[string]string),
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "provider":
			profile.Provider = value
		case "model":
			profile.Model = value
		case "timeout":
			profile.Timeout = value
		default:
			if strings.HasPrefix(key, "ENV_") {
				envKey := strings.TrimPrefix(key, "ENV_")
				profile.Environment[envKey] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return Profile{}, err
	}

	return profile, nil
}

// SaveProfile guarda un perfil en disco
func (pm *ProfileManager) SaveProfile(name string, profile Profile) error {
	profileDir := filepath.Join(pm.configDir, "profiles")
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		return fmt.Errorf("no se pudo crear directorio perfiles: %w", err)
	}

	profilePath := filepath.Join(profileDir, name+".conf")
	file, err := os.Create(profilePath)
	if err != nil {
		return fmt.Errorf("no se pudo crear archivo perfil: %w", err)
	}
	defer file.Close()

	// Escribir configuración básica
	if profile.Provider != "" {
		fmt.Fprintf(file, "provider=%s\n", profile.Provider)
	}
	if profile.Model != "" {
		fmt.Fprintf(file, "model=%s\n", profile.Model)
	}
	if profile.Timeout != "" {
		fmt.Fprintf(file, "timeout=%s\n", profile.Timeout)
	}

	// Escribir variables de entorno
	for key, value := range profile.Environment {
		fmt.Fprintf(file, "ENV_%s=%s\n", key, value)
	}

	return nil
}

// GetProfile obtiene un perfil por nombre
func (pm *ProfileManager) GetProfile(name string) (Profile, bool) {
	profile, exists := pm.profiles[name]
	return profile, exists
}

// ListProfiles retorna todos los perfiles disponibles
func (pm *ProfileManager) ListProfiles() []string {
	var names []string
	for name := range pm.profiles {
		names = append(names, name)
	}
	return names
}

// ApplyProfile aplica un perfil a la configuración actual
func (pm *ProfileManager) ApplyProfile(profile Profile) error {
	// Aplicar variables de entorno del perfil
	for key, value := range profile.Environment {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("no se pudo establecer variable %s: %w", key, err)
		}
	}

	// Aplicar configuración específica
	if profile.Model != "" {
		os.Setenv("MAIN_MODEL", profile.Model)
		os.Setenv("ANTHROPIC_MODEL", profile.Model)
	}

	return nil
}

// CreateDefaultProfiles crea perfiles predefinidos
func (pm *ProfileManager) CreateDefaultProfiles() error {
	defaults := map[string]Profile{
		"dev": {
			Provider: "mimo",
			Model:    "mimo-v2-flash",
			Timeout:  "5m",
			Environment: map[string]string{
				"DEBUG": "true",
			},
		},
		"prod": {
			Provider: "mimo",
			Model:    "mimo-v2-flash",
			Timeout:  "10m",
			Environment: map[string]string{
				"DEBUG": "false",
			},
		},
		"test": {
			Provider: "claude",
			Model:    "claude-3-5-sonnet",
			Timeout:  "2m",
			Environment: map[string]string{
				"DEBUG": "true",
			},
		},
	}

	for name, profile := range defaults {
		if err := pm.SaveProfile(name, profile); err != nil {
			return fmt.Errorf("no se pudo crear perfil %s: %w", name, err)
		}
		pm.profiles[name] = profile
	}

	return nil
}

// DeleteProfile elimina un perfil
func (pm *ProfileManager) DeleteProfile(name string) error {
	profileDir := filepath.Join(pm.configDir, "profiles")
	profilePath := filepath.Join(profileDir, name+".conf")

	if err := os.Remove(profilePath); err != nil {
		return fmt.Errorf("no se pudo eliminar perfil %s: %w", name, err)
	}

	delete(pm.profiles, name)
	return nil
}

// ProfileExists verifica si un perfil existe
func (pm *ProfileManager) ProfileExists(name string) bool {
	_, exists := pm.profiles[name]
	return exists
}