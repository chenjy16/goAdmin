package dto

// UnifiedMessage 统一的消息结构
type UnifiedMessage struct {
	Role    string `json:"role" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// UnifiedChatRequest 统一的聊天请求
type UnifiedChatRequest struct {
	Model       string                 `json:"model" binding:"required"`
	Messages    []UnifiedMessage       `json:"messages" binding:"required,min=1"`
	MaxTokens   *int                   `json:"max_tokens,omitempty"`
	Temperature *float64               `json:"temperature,omitempty"`
	TopP        *float64               `json:"top_p,omitempty"`
	TopK        *int                   `json:"top_k,omitempty"`
	Stream      bool                   `json:"stream,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
}

// SetAPIKeyRequest 设置API密钥请求
type SetAPIKeyRequest struct {
	APIKey string `json:"api_key" binding:"required"`
}

// ProvidersResponse 提供商列表响应
type ProvidersResponse struct {
	Providers []ProviderInfo `json:"providers"`
}

// ProviderInfo 提供商信息
type ProviderInfo struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Healthy     bool   `json:"healthy"`
	ModelCount  int    `json:"model_count"`
}

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