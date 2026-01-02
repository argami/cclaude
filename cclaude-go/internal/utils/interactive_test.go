package utils

import (
	"os"
	"testing"
)

func TestInteractiveMode(t *testing.T) {
	im := NewInteractiveMode()
	if im == nil {
		t.Error("Expected InteractiveMode to be created")
	}
	if im.reader == nil {
		t.Error("Expected reader to be initialized")
	}
	if len(im.providers) == 0 {
		t.Error("Expected providers to be initialized")
	}
}

func TestShowConfig(t *testing.T) {
	// Capturar stdout
	old := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	ShowConfig()

	w.Close()
	os.Stdout = old

	// Verificar que no hay errores
	// ShowConfig debería ejecutarse sin panic
}

func TestShowTips(t *testing.T) {
	// Capturar stdout
	old := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	ShowTips()

	w.Close()
	os.Stdout = old

	// Verificar que no hay errores
	// ShowTips debería ejecutarse sin panic
}

func TestNewInteractiveMode(t *testing.T) {
	im := NewInteractiveMode()
	if im == nil {
		t.Error("Expected non-nil InteractiveMode")
	}
	if im.reader == nil {
		t.Error("Expected reader to be initialized")
	}
}