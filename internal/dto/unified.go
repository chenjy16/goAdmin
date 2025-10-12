package dto

import "go-springAi/internal/types"

// UnifiedMessage 统一的消息结构
type UnifiedMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// SetAPIKeyRequest 设置API密钥请求
type SetAPIKeyRequest struct {
	APIKey string `json:"api_key" binding:"required"`
}

// ProvidersResponse 提供商列表响应
type ProvidersResponse struct {
	Providers []ProviderInfo `json:"providers"`
}

// 使用共享的提供商信息类型
type ProviderInfo = types.CommonProviderInfo

// ModelsResponse 模型列表响应
type ModelsResponse struct {
	Provider string                 `json:"provider"`
	Models   map[string]interface{} `json:"models"`
}

// ModelConfigResponse 模型配置响应
type ModelConfigResponse struct {
	Provider string      `json:"provider"`
	Model    string      `json:"model"`
	Config   interface{} `json:"config"`
}

// ValidationResponse API密钥验证响应
type ValidationResponse struct {
	Provider string `json:"provider"`
	Valid    bool   `json:"valid"`
	Message  string `json:"message,omitempty"`
}