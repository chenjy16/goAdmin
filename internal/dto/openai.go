package dto

import (
	"go-springAi/internal/openai"
)

// OpenAIChatRequest 聊天完成请求
type OpenAIChatRequest struct {
	Model       string                 `json:"model" binding:"required"`
	Messages    []openai.Message       `json:"messages" binding:"required,min=1"`
	MaxTokens   *int                   `json:"max_tokens,omitempty"`
	Temperature *float64               `json:"temperature,omitempty"`
	TopP        *float64               `json:"top_p,omitempty"`
	Stream      bool                   `json:"stream,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
}

// OpenAIChatResponse 聊天完成响应
type OpenAIChatResponse struct {
	ID      string           `json:"id"`
	Object  string           `json:"object"`
	Created int64            `json:"created"`
	Model   string           `json:"model"`
	Choices []openai.Choice  `json:"choices"`
	Usage   openai.Usage     `json:"usage"`
}

// OpenAISetAPIKeyRequest 设置 API 密钥请求
type OpenAISetAPIKeyRequest struct {
	APIKey string `json:"api_key" binding:"required"`
}

// OpenAIModelsResponse 模型列表响应
type OpenAIModelsResponse struct {
	Models map[string]*openai.ModelConfig `json:"models"`
}

// OpenAIModelConfigResponse 模型配置响应
type OpenAIModelConfigResponse struct {
	Model  string                `json:"model"`
	Config *openai.ModelConfig   `json:"config"`
}

// OpenAIValidationResponse API 密钥验证响应
type OpenAIValidationResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message,omitempty"`
}