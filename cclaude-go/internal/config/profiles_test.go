package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProfileManager(t *testing.T) {
	// Crear un directorio temporal para tests
	tempDir := t.TempDir()

	// Mock el directorio de configuración
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	pm, err := NewProfileManager()
	if err != nil {
		t.Fatalf("NewProfileManager() error = %v", err)
	}

	// Test crear perfiles por defecto
	err = pm.CreateDefaultProfiles()
	if err != nil {
		t.Fatalf("CreateDefaultProfiles() error = %v", err)
	}

	// Test listar perfiles
	profiles := pm.ListProfiles()
	if len(profiles) == 0 {
		t.Error("Expected profiles to be created")
	}

	// Test obtener perfil
	profile, exists := pm.GetProfile("dev")
	if !exists {
		t.Error("Expected dev profile to exist")
	}
	if profile.Provider != "mimo" {
		t.Errorf("Expected provider 'mimo', got '%s'", profile.Provider)
	}

	// Test guardar perfil personalizado
	customProfile := Profile{
		Provider: "kimi",
		Model:    "kimi-custom",
		Timeout:  "10m",
		Environment: map[string]string{
			"DEBUG": "true",
		},
	}
	err = pm.SaveProfile("custom", customProfile)
	if err != nil {
		t.Fatalf("SaveProfile() error = %v", err)
	}

	// Test aplicar perfil
	err = pm.ApplyProfile(customProfile)
	if err != nil {
		t.Fatalf("ApplyProfile() error = %v", err)
	}

	// Verificar variables de entorno aplicadas
	if os.Getenv("MAIN_MODEL") != "kimi-custom" {
		t.Error("MAIN_MODEL not applied correctly")
	}

	// Test eliminar perfil
	err = pm.DeleteProfile("custom")
	if err != nil {
		t.Fatalf("DeleteProfile() error = %v", err)
	}

	// Verificar eliminación
	if pm.ProfileExists("custom") {
		t.Error("Profile should be deleted")
	}
}

func TestProfileManagerLoadProfiles(t *testing.T) {
	// Crear directorio temporal
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Crear archivo de perfil manualmente
	profileDir := filepath.Join(tempDir, ".config", "cclaude", "profiles")
	os.MkdirAll(profileDir, 0755)

	profileContent := `provider=test
model=test-model
timeout=5m
ENV_DEBUG=true
`
	err := os.WriteFile(filepath.Join(profileDir, "test.conf"), []byte(profileContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test profile: %v", err)
	}

	pm, err := NewProfileManager()
	if err != nil {
		t.Fatalf("NewProfileManager() error = %v", err)
	}

	err = pm.LoadProfiles()
	if err != nil {
		t.Fatalf("LoadProfiles() error = %v", err)
	}

	profile, exists := pm.GetProfile("test")
	if !exists {
		t.Error("Expected test profile to be loaded")
	}
	if profile.Model != "test-model" {
		t.Errorf("Expected model 'test-model', got '%s'", profile.Model)
	}
}

func TestProfileManagerErrors(t *testing.T) {
	// Test con directorio home inaccesible
	os.Setenv("HOME", "/nonexistent")
	defer os.Unsetenv("HOME")

	_, err := NewProfileManager()
	if err == nil {
		t.Error("Expected error with invalid home directory")
	}
}