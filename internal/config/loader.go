package config

import (
	"os"
	"path/filepath"
	"github.com/spf13/viper"
)

type Provider struct {
	Name      string `mapstructure:"name"`
	BaseURL   string `mapstructure:"base_url"`
	Model     string `mapstructure:"model"`
	EnvKey    string `mapstructure:"env_key"`
	OpusModel string `mapstructure:"opus_model"`
}

type Settings struct {
	TimeoutMs            int    `mapstructure:"timeout_ms"`
	DisableNonEssential bool   `mapstructure:"disable_non_essential_calls"`
	LogLevel              string `mapstructure:"log_level"`
}

type Config struct {
	Providers map[string]*Provider `mapstructure:"providers"`
	Settings  Settings           `mapstructure:"settings"`
}

func Load() (*Config, error) {
	cfg := &Config{
		Providers: getDefaultProviders(),
		Settings: getDefaultSettings(),
	}

	// Try loading from XDG config dirs
	configPaths := getConfigPaths()
	for _, path := range configPaths {
		configFile := filepath.Join(path, "config.yaml")
		if _, err := os.Stat(configFile); err == nil {
			viper.SetConfigFile(configFile)
			viper.SetConfigType("yaml")

			if err := viper.ReadInConfig(); err == nil {
				return cfg, nil
			}
		}
	}

	return cfg, nil
}

func getDefaultProviders() map[string]*Provider {
	return map[string]*Provider{
		"mimo": {
			Name:      "Mimo",
			BaseURL:   "https://api.xiaomimimo.com/anthropic",
			Model:     "mimo-v2-flash",
			EnvKey:    "MIMO_API_KEY",
			OpusModel: "mimo-v2-flash",
		},
		"minimax": {
			Name:      "MiniMax",
			BaseURL:   "https://api.minimax.io/anthropic",
			Model:     "MiniMax-M2.1",
			EnvKey:    "MINIMAX_API_KEY",
			OpusModel: "MiniMax-M2.1",
		},
		"kimi": {
			Name:      "Kimi",
			BaseURL:   "https://api.kimi.com/coding/",
			Model:     "kimi-k2-0711-preview",
			EnvKey:    "KIMI_API_KEY",
			OpusModel: "kimi-k2-thinking-turbo",
		},
		"glm": {
			Name:      "GLM",
			BaseURL:   "https://api.z.ai/api/anthropic",
			Model:     "glm-4.7",
			EnvKey:    "GLM_API_KEY",
			OpusModel: "glm-4.7",
		},
		"claude": {
			Name:      "Claude",
			BaseURL:   "",
			Model:     "",
			EnvKey:    "",
			OpusModel: "",
		},
	}
}

func getDefaultSettings() Settings {
	return Settings{
		TimeoutMs:            3000000,
		DisableNonEssential: true,
		LogLevel:              "info",
	}
}

func getConfigPaths() []string {
	var paths []string

	// XDG Base Directory
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		paths = append(paths, filepath.Join(xdg, "cclaude"))
	}

	// Fallback to home directory
	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(home, ".config", "cclaude"))
	}

	return paths
}
