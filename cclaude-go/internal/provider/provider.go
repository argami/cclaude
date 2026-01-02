package provider

import (
	"fmt"

	"github.com/argami/cclaude-go/pkg/types"
)

var Providers = map[string]types.ProviderConfig{
	"mimo": {
		Name:        "mimo",
		BaseURL:     "https://api.xiaomimimo.com/anthropic",
		Model:       "mimo-v2-flash",
		OpusModel:   "mimo-v2-flash",
		EnvVar:      "MIMO_API_KEY",
		Description: "Xiaomi MiMo API",
	},
	"minimax": {
		Name:        "minimax",
		BaseURL:     "https://api.minimax.io/anthropic",
		Model:       "MiniMax-M2.1",
		OpusModel:   "MiniMax-M2.1",
		EnvVar:      "MINIMAX_API_KEY",
		Description: "MiniMax API",
	},
	"kimi": {
		Name:        "kimi",
		BaseURL:     "https://api.kimi.com/coding/",
		Model:       "kimi-k2-0711-preview",
		OpusModel:   "kimi-k2-thinking-turbo",
		EnvVar:      "KIMI_API_KEY",
		Description: "Kimi API",
	},
	"glm": {
		Name:        "glm",
		BaseURL:     "https://api.z.ai/api/anthropic",
		Model:       "glm-4.7",
		OpusModel:   "glm-4.7",
		EnvVar:      "GLM_API_KEY",
		Description: "Zhipu AI API",
	},
	"claude": {
		Name:        "claude",
		BaseURL:     "",
		Model:       "",
		OpusModel:   "",
		EnvVar:      "",
		Description: "Claude Native",
	},
}

func GetProvider(name string) (*types.ProviderConfig, error) {
	if provider, exists := Providers[name]; exists {
		return &provider, nil
	}
	return nil, fmt.Errorf("proveedor no encontrado: %s", name)
}