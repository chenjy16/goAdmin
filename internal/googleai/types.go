package googleai

import (
	"context"
	"io"
)

// Message 聊天消息
type Message struct {
	Role    string `json:"role"`    // user, model
	Content string `json:"content"`
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Model            string    `json:"model"`
	Messages         []Message `json:"messages"`
	MaxTokens        int       `json:"max_tokens,omitempty"`
	Temperature      float32   `json:"temperature,omitempty"`
	TopP             float32   `json:"top_p,omitempty"`
	TopK             int       `json:"top_k,omitempty"`
	Stream           bool      `json:"stream,omitempty"`
}

// Choice 响应选择
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage 使用统计
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// StreamChoice 流式响应选择
type StreamChoice struct {
	Index int `json:"index"`
	Delta struct {
		Role    string `json:"role,omitempty"`
		Content string `json:"content,omitempty"`
	} `json:"delta"`
	FinishReason *string `json:"finish_reason"`
}

// StreamResponse 流式响应
type StreamResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []StreamChoice `json:"choices"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

// ModelConfig 模型配置
type ModelConfig struct {
	Name        string  `json:"name"`
	DisplayName string  `json:"display_name"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float32 `json:"temperature"`
	TopP        float32 `json:"top_p"`
	TopK        int     `json:"top_k"`
	Enabled     bool    `json:"enabled"`
}

// Client Google AI 客户端接口
type Client interface {
	// ChatCompletion 聊天完成
	ChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error)

	// ChatCompletionStream 流式聊天完成
	ChatCompletionStream(ctx context.Context, req *ChatRequest) (io.ReadCloser, error)

	// ListModels 列出可用模型
	ListModels(ctx context.Context) ([]string, error)

	// ValidateAPIKey 验证API密钥
	ValidateAPIKey(ctx context.Context) error

	// ResetClient 重置客户端，强制重新初始化
	ResetClient()
}

// ModelManager 模型管理器接口
type ModelManager interface {
	// GetModel 获取模型配置
	GetModel(name string) (*ModelConfig, error)

	// ListModels 列出所有模型
	ListModels() map[string]*ModelConfig

	// UpdateModel 更新模型配置
	UpdateModel(name string, config *ModelConfig) error

	// EnableModel 启用模型
	EnableModel(name string) error

	// DisableModel 禁用模型
	DisableModel(name string) error
}

// KeyManager API密钥管理器接口
type KeyManager interface {
	// SetAPIKey 设置API密钥
	SetAPIKey(key string) error

	// GetAPIKey 获取API密钥
	GetAPIKey() (string, error)

	// ValidateKey 验证密钥
	ValidateKey(key string) error

	// EncryptKey 加密密钥
	EncryptKey(key string) (string, error)

	// DecryptKey 解密密钥
	DecryptKey(encryptedKey string) (string, error)
}