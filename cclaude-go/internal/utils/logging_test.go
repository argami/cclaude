package utils

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestLogLevels(t *testing.T) {
	tests := []struct {
		name     string
		logLevel LogLevel
		function func(string, ...interface{})
		should   string
	}{
		{"Info at Info level", LevelInfo, Info, "show"},
		{"Warn at Info level", LevelInfo, Warn, "show"},
		{"Error at Info level", LevelInfo, Error, "show"},
		{"Debug at Info level", LevelInfo, Debug, "hide"},

		{"Info at Error level", LevelError, Info, "hide"},
		{"Warn at Error level", LevelError, Warn, "hide"},
		{"Error at Error level", LevelError, Error, "show"},
		{"Debug at Error level", LevelError, Debug, "hide"},

		{"All at Debug level", LevelDebug, Info, "show"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capturar stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			SetLogLevel(tt.logLevel)
			tt.function("test message")

			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			contains := strings.Contains(output, "test message")
			if tt.should == "show" && !contains {
				t.Errorf("Expected log to show, but got: %s", output)
			}
			if tt.should == "hide" && contains {
				t.Errorf("Expected log to hide, but got: %s", output)
			}
		})
	}
}

func TestLogFormat(t *testing.T) {
	// Capturar stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	SetLogLevel(LevelInfo)
	Info("test message")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verificar formato: [timestamp] ℹ️  message
	if !strings.Contains(output, "ℹ️") {
		t.Errorf("Expected info emoji in output, got: %s", output)
	}
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected message in output, got: %s", output)
	}
	if !strings.Contains(output, "[") || !strings.Contains(output, "]") {
		t.Errorf("Expected timestamp brackets in output, got: %s", output)
	}
}

func TestLogPrefixes(t *testing.T) {
	tests := []struct {
		name     string
		function func(string, ...interface{})
		prefix   string
	}{
		{"Info", Info, "ℹ️"},
		{"Warn", Warn, "⚠️"},
		{"Error", Error, "❌"},
		{"Debug", Debug, "🔍"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capturar stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			SetLogLevel(LevelDebug)
			tt.function("test")

			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			if !strings.Contains(output, tt.prefix) {
				t.Errorf("Expected prefix %s in output, got: %s", tt.prefix, output)
			}
		})
	}
}