package dto

import (
	"go-springAi/internal/googleai"
)

// GoogleAIChatRequest Google AI 聊天完成请求
type GoogleAIChatRequest struct {
	Model       string                     `json:"model" binding:"required"`
	Messages    []googleai.Message         `json:"messages" binding:"required,min=1"`
	MaxTokens   *int                       `json:"max_tokens,omitempty"`
	Temperature *float64                   `json:"temperature,omitempty"`
	TopP        *float64                   `json:"top_p,omitempty"`
	TopK        *int                       `json:"top_k,omitempty"`
	Stream      bool                       `json:"stream,omitempty"`
	Options     map[string]interface{}     `json:"options,omitempty"`
}

// GoogleAIChatResponse Google AI 聊天完成响应
type GoogleAIChatResponse struct {
	ID      string              `json:"id"`
	Object  string              `json:"object"`
	Created int64               `json:"created"`
	Model   string              `json:"model"`
	Choices []googleai.Choice   `json:"choices"`
	Usage   googleai.Usage      `json:"usage"`
}

// GoogleAISetAPIKeyRequest 设置 API 密钥请求
type GoogleAISetAPIKeyRequest struct {
	APIKey string `json:"api_key" binding:"required"`
}

// GoogleAIModelsResponse 模型列表响应
type GoogleAIModelsResponse struct {
	Models map[string]*googleai.ModelConfig `json:"models"`
}

// GoogleAIModelConfigResponse 模型配置响应
type GoogleAIModelConfigResponse struct {
	Model  string                    `json:"model"`
	Config *googleai.ModelConfig     `json:"config"`
}

// GoogleAIValidationResponse API 密钥验证响应
type GoogleAIValidationResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message,omitempty"`
}