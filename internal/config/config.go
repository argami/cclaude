package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Provider     string            `yaml:"provider"`
	Model        string            `yaml:"model"`
	TimeoutMS    int               `yaml:"timeout_ms"`
	EnvOverrides map[string]string `yaml:"env_overrides"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	config := &Config{}

	// Try to load from file
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return defaults if file doesn't exist
			return getDefaults(), nil
		}
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply environment variable overrides
	applyEnvOverrides(config)

	return config, nil
}

// applyEnvOverrides applies environment variable overrides to config
func applyEnvOverrides(config *Config) {
	// Override provider from environment
	if v := os.Getenv("CCLAUDE_PROVIDER"); v != "" && config.Provider == "" {
		config.Provider = v
	}

	// Override model from environment
	if v := os.Getenv("CCLAUDE_MODEL"); v != "" && config.Model == "" {
		config.Model = v
	}
}

// getDefaults returns the default configuration
func getDefaults() *Config {
	return &Config{
		Provider:     "",
		Model:        "",
		TimeoutMS:    3000000, // 50 minutes default
		EnvOverrides: make(map[string]string),
	}
}

// GetProvider returns the configured provider or empty string for default
func (c *Config) GetProvider() string {
	return c.Provider
}

// GetModel returns the configured model or empty string for provider default
func (c *Config) GetModel() string {
	return c.Model
}

// GetTimeoutMS returns the timeout in milliseconds
func (c *Config) GetTimeoutMS() int {
	if c.TimeoutMS == 0 {
		return 3000000
	}
	return c.TimeoutMS
}
