package provider

import (
	"context"
	"io"
)

// ProviderType 提供商类型
type ProviderType string

const (
	ProviderTypeOpenAI   ProviderType = "openai"
	ProviderTypeGoogleAI ProviderType = "googleai"
	ProviderTypeMock     ProviderType = "mock"
)

// Message 统一的消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Choice 统一的选择结构
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage 统一的使用情况结构
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatRequest 统一的聊天请求结构
type ChatRequest struct {
	Model       string                 `json:"model"`
	Messages    []Message              `json:"messages"`
	MaxTokens   *int                   `json:"max_tokens,omitempty"`
	Temperature *float32               `json:"temperature,omitempty"`
	TopP        *float32               `json:"top_p,omitempty"`
	TopK        *int                   `json:"top_k,omitempty"` // Google AI 特有
	Stream      bool                   `json:"stream,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
}

// ChatResponse 统一的聊天响应结构
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// ModelConfig 统一的模型配置结构
type ModelConfig struct {
	Name         string   `json:"name"`
	DisplayName  string   `json:"display_name"`
	MaxTokens    int      `json:"max_tokens"`
	Temperature  float32  `json:"temperature"`
	TopP         float32  `json:"top_p"`
	TopK         int      `json:"top_k,omitempty"` // Google AI 特有
	Enabled      bool     `json:"enabled"`
}

// Provider 统一的AI提供商接口
type Provider interface {
	// GetType 获取提供商类型
	GetType() ProviderType
	
	// GetName 获取提供商名称
	GetName() string
	
	// ChatCompletion 聊天完成
	ChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
	
	// ChatCompletionStream 流式聊天完成
	ChatCompletionStream(ctx context.Context, req *ChatRequest) (io.ReadCloser, error)
	
	// ListModels 列出可用模型（仅启用的）
	ListModels(ctx context.Context) (map[string]*ModelConfig, error)
	
	// ListAllModels 列出所有模型（包括禁用的）
	ListAllModels(ctx context.Context) (map[string]*ModelConfig, error)
	
	// GetModelConfig 获取模型配置
	GetModelConfig(name string) (*ModelConfig, error)
	
	// EnableModel 启用模型
	EnableModel(name string) error
	
	// DisableModel 禁用模型
	DisableModel(name string) error
	
	// ValidateAPIKey 验证API密钥
	ValidateAPIKey(ctx context.Context) error
	
	// SetAPIKey 设置API密钥
	SetAPIKey(key string) error
	
	// IsHealthy 检查提供商健康状态
	IsHealthy(ctx context.Context) bool
}

// ProviderInfo 提供商信息
type ProviderInfo struct {
	Type        ProviderType `json:"type"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Healthy     bool         `json:"healthy"`
	ModelCount  int          `json:"model_count"`
}