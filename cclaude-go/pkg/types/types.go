package types

import "time"

type ProviderConfig struct {
	Name        string
	BaseURL     string
	Model       string
	OpusModel   string
	EnvVar      string
	Description string
}

type AppConfig struct {
	Provider      *ProviderConfig
	Timeout       time.Duration
	Debug         bool
	ModelOverride string
	Args          []string
}