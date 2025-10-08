package openai

import (
	"time"
)

// Config OpenAI 配置
type Config struct {
	APIKey      string        `json:"api_key" yaml:"api_key"`
	BaseURL     string        `json:"base_url" yaml:"base_url"`
	Timeout     time.Duration `json:"timeout" yaml:"timeout"`
	MaxRetries  int           `json:"max_retries" yaml:"max_retries"`
	DefaultModel string       `json:"default_model" yaml:"default_model"`
}

// ModelConfig 模型配置
type ModelConfig struct {
	Name            string  `json:"name"`
	MaxTokens       int     `json:"max_tokens"`
	Temperature     float32 `json:"temperature"`
	TopP            float32 `json:"top_p"`
	FrequencyPenalty float32 `json:"frequency_penalty"`
	PresencePenalty float32 `json:"presence_penalty"`
	Enabled         bool    `json:"enabled"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		BaseURL:      "https://api.openai.com/v1",
		Timeout:      30 * time.Second,
		MaxRetries:   3,
		DefaultModel: "gpt-3.5-turbo",
	}
}

// DefaultModels 返回默认支持的模型配置
func DefaultModels() map[string]*ModelConfig {
	return map[string]*ModelConfig{
		"gpt-4": {
			Name:             "gpt-4",
			MaxTokens:        8192,
			Temperature:      0.7,
			TopP:             1.0,
			FrequencyPenalty: 0.0,
			PresencePenalty:  0.0,
			Enabled:          true,
		},
		"gpt-4-turbo": {
			Name:             "gpt-4-turbo",
			MaxTokens:        128000,
			Temperature:      0.7,
			TopP:             1.0,
			FrequencyPenalty: 0.0,
			PresencePenalty:  0.0,
			Enabled:          true,
		},
		"gpt-3.5-turbo": {
			Name:             "gpt-3.5-turbo",
			MaxTokens:        4096,
			Temperature:      0.7,
			TopP:             1.0,
			FrequencyPenalty: 0.0,
			PresencePenalty:  0.0,
			Enabled:          true,
		},
		"gpt-3.5-turbo-16k": {
			Name:             "gpt-3.5-turbo-16k",
			MaxTokens:        16384,
			Temperature:      0.7,
			TopP:             1.0,
			FrequencyPenalty: 0.0,
			PresencePenalty:  0.0,
			Enabled:          true,
		},
	}
}