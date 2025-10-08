package googleai

import "time"

// Config Google AI 配置
type Config struct {
	APIKey       string        `json:"api_key"`
	ProjectID    string        `json:"project_id"`
	Location     string        `json:"location"`
	Timeout      time.Duration `json:"timeout"`
	MaxRetries   int           `json:"max_retries"`
	DefaultModel string        `json:"default_model"`
}