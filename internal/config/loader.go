package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Provider struct {
	Name      string
	BaseURL   string
	Model     string
	EnvKey    string
	OpusModel string
}

type Settings struct {
	TimeoutMs            int
	DisableNonEssential bool
	LogLevel              string
}

type Config struct {
	Providers map[string]*Provider
	Settings  Settings
	mu        sync.RWMutex // Protects Providers and Settings
}

var (
	currentConfig *Config
	watcher      *fsnotify.Watcher
	configMutex  sync.RWMutex
	onConfigChange []func(*Config)
	changeMutex  sync.RWMutex
)

// Load loads configuration from file with fallback to defaults
func Load() (*Config, error) {
	cfg := &Config{
		Providers: getDefaultProviders(),
		Settings:  getDefaultSettings(),
	}

	// Try loading from config paths
	configFile, err := findConfigFile()
	if err == nil {
		if loadErr := loadFromFile(cfg, configFile); loadErr == nil {
			cfg.mu.Lock()
			currentConfig = cfg
			cfg.mu.Unlock()
			return cfg, nil
		}
	}

	cfg.mu.Lock()
	currentConfig = cfg
	cfg.mu.Unlock()

	return cfg, nil
}

// GetConfig returns the current configuration (thread-safe)
func GetConfig() *Config {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return currentConfig
}

// GetProvider returns a specific provider (thread-safe)
func (c *Config) GetProvider(name string) (*Provider, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	provider, exists := c.Providers[name]
	return provider, exists
}

// ListProviders returns all provider names (thread-safe)
func (c *Config) ListProviders() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	names := make([]string, 0, len(c.Providers))
	for name := range c.Providers {
		names = append(names, name)
	}
	return names
}

// Reload reloads the configuration from file
func Reload() error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	configMutex.Lock()
	currentConfig = cfg
	configMutex.Unlock()

	// Notify change listeners
	changeMutex.RLock()
	listeners := make([]func(*Config), len(onConfigChange))
	copy(listeners, onConfigChange)
	changeMutex.RUnlock()

	for _, listener := range listeners {
		listener(cfg)
	}

	return nil
}

// OnConfigChange registers a callback for configuration changes
func OnConfigChange(callback func(*Config)) {
	changeMutex.Lock()
	defer changeMutex.Unlock()
	onConfigChange = append(onConfigChange, callback)
}

// Watch starts watching the config file for changes
func Watch() error {
	configFile, err := findConfigFile()
	if err != nil {
		return fmt.Errorf("no config file to watch: %w", err)
	}

	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}

	// Watch the config file directory
	configDir := filepath.Dir(configFile)
	if err := watcher.Add(configDir); err != nil {
		return fmt.Errorf("failed to watch directory: %w", err)
	}

	// Start watching in background
	go watchConfig()

	return nil
}

// watchConfig monitors config file changes
func watchConfig() {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				if err := Reload(); err != nil {
					fmt.Fprintf(os.Stderr, "Error reloading config: %v\n", err)
				} else {
					fmt.Println("Configuration reloaded successfully")
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Fprintf(os.Stderr, "Watcher error: %v\n", err)
		}
	}
}

// StopWatching stops the config file watcher
func StopWatching() error {
	if watcher != nil {
		return watcher.Close()
	}
	return nil
}

func loadFromFile(cfg *Config, configFile string) error {
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// Load providers - iterate through provider names and extract each
	providersMap := viper.GetStringMap("providers")
	for name := range providersMap {
		cfg.Providers[name] = &Provider{
			Name:      viper.GetString("providers." + name + ".name"),
			BaseURL:   viper.GetString("providers." + name + ".base_url"),
			Model:     viper.GetString("providers." + name + ".model"),
			EnvKey:    viper.GetString("providers." + name + ".env_key"),
			OpusModel: viper.GetString("providers." + name + ".opus_model"),
		}
	}

	// Load settings directly with proper type conversion
	cfg.Settings.TimeoutMs = viper.GetInt("settings.timeout_ms")
	cfg.Settings.DisableNonEssential = viper.GetBool("settings.disable_non_essential_calls")
	cfg.Settings.LogLevel = viper.GetString("settings.log_level")

	return nil
}

func findConfigFile() (string, error) {
	configPaths := getConfigPaths()
	for _, path := range configPaths {
		configFile := filepath.Join(path, "config.yaml")
		if _, err := os.Stat(configFile); err == nil {
			return configFile, nil
		}
	}
	return "", fmt.Errorf("no config file found in paths: %v", configPaths)
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

	// Current directory (for development)
	if cwd, err := os.Getwd(); err == nil {
		paths = append(paths, cwd)
	}

	// XDG Base Directory
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		paths = append(paths, filepath.Join(xdg, "cclaude"))
	}

	// Fallback to home directory
	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(home, ".config", "cclaude"))
	}

	// System-wide config
	paths = append(paths, "/etc/cclaude")

	return paths
}
