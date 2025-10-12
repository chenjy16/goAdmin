package provider

import (
	"context"
	"io"
	"go-springAi/internal/types"
)

// 使用共享的通用类型定义
type Message = types.CommonMessage
type Choice = types.CommonChoice
type Usage = types.CommonUsage
type ChatRequest = types.CommonChatRequest
type ChatResponse = types.CommonChatResponse

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

// 使用共享的通用提供商信息类型
type ProviderInfo = types.CommonProviderInfo
type ProviderType = types.ProviderType