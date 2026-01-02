package models

import "fmt"

// Provider represents a Claude API provider configuration
type Provider struct {
	Name    string
	BaseURL string
	Model   string
	EnvKey  string
}

// providers map contains all supported providers
var providers = map[string]Provider{
	"mimo": {
		Name:    "mimo",
		BaseURL: "https://api.xiaomimimo.com/anthropic",
		Model:   "mimo-v2-flash",
		EnvKey:  "MIMO_API_KEY",
	},
	"minimax": {
		Name:    "minimax",
		BaseURL: "https://api.minimax.io/anthropic",
		Model:   "MiniMax-M2.1",
		EnvKey:  "MINIMAX_API_KEY",
	},
	"kimi": {
		Name:    "kimi",
		BaseURL: "https://api.kimi.com/coding/",
		Model:   "kimi-k2-0711-preview",
		EnvKey:  "KIMI_API_KEY",
	},
	"glm": {
		Name:    "glm",
		BaseURL: "https://api.z.ai/api/anthropic",
		Model:   "glm-4.7",
		EnvKey:  "GLM_API_KEY",
	},
}

// GetProvider returns the provider configuration for the given name
func GetProvider(name string) (Provider, error) {
	p, ok := providers[name]
	if !ok {
		return Provider{}, fmt.Errorf("provider not found: %s", name)
	}
	return p, nil
}

// GetProviderNames returns a slice of all available provider names
func GetProviderNames() []string {
	names := make([]string, 0, len(providers))
	for name := range providers {
		names = append(names, name)
	}
	return names
}
