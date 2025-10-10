package provider

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"
)

// MockProvider 模拟提供商实现，用于测试
type MockProvider struct {
	name string
	providerType ProviderType
	models map[string]*ModelConfig
	mu sync.RWMutex
}

// NewMockProvider 创建模拟提供商
func NewMockProvider(name string, providerType ProviderType) *MockProvider {
	p := &MockProvider{
		name: name,
		providerType: providerType,
		models: make(map[string]*ModelConfig),
	}
	
	// 初始化默认模型配置
	p.initDefaultModels()
	
	return p
}

// initDefaultModels 初始化默认模型配置
func (p *MockProvider) initDefaultModels() {
	p.models = map[string]*ModelConfig{
		"mock-model-1": {
			Name:        "mock-model-1",
			DisplayName: "Mock Model 1",
			MaxTokens:   4096,
			Temperature: 0.7,
			TopP:        0.9,
			Enabled:     true,
		},
		"mock-model-2": {
			Name:        "mock-model-2", 
			DisplayName: "Mock Model 2",
			MaxTokens:   8192,
			Temperature: 0.7,
			TopP:        0.9,
			Enabled:     true,
		},
		"mock-model-3": {
			Name:        "mock-model-3", 
			DisplayName: "Mock Model 3 (Disabled)",
			MaxTokens:   2048,
			Temperature: 0.5,
			TopP:        0.8,
			Enabled:     false,
		},
	}
}

// GetType 获取提供商类型
func (p *MockProvider) GetType() ProviderType {
	return p.providerType
}

// GetName 获取提供商名称
func (p *MockProvider) GetName() string {
	return p.name
}

// ChatCompletion 模拟聊天完成
func (p *MockProvider) ChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// 模拟响应
	response := &ChatResponse{
		ID:      fmt.Sprintf("mock-%d", time.Now().Unix()),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []Choice{
			{
				Index: 0,
				Message: Message{
					Role:    "assistant",
					Content: fmt.Sprintf("这是来自 %s 提供商的模拟响应，当前使用的模型是: %s。您的消息是: %s", p.name, req.Model, req.Messages[len(req.Messages)-1].Content),
				},
				FinishReason: "stop",
			},
		},
		Usage: Usage{
			PromptTokens:     50,
			CompletionTokens: 20,
			TotalTokens:      70,
		},
	}

	return response, nil
}

// ChatCompletionStream 模拟流式聊天完成
func (p *MockProvider) ChatCompletionStream(ctx context.Context, req *ChatRequest) (io.ReadCloser, error) {
	return nil, fmt.Errorf("stream not implemented for mock provider")
}

// ListModels 列出模型（仅启用的）
func (p *MockProvider) ListModels(ctx context.Context) (map[string]*ModelConfig, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	enabledModels := make(map[string]*ModelConfig)
	for name, model := range p.models {
		if model.Enabled {
			// 创建副本以避免并发修改
			enabledModels[name] = &ModelConfig{
				Name:        model.Name,
				DisplayName: model.DisplayName,
				MaxTokens:   model.MaxTokens,
				Temperature: model.Temperature,
				TopP:        model.TopP,
				Enabled:     model.Enabled,
			}
		}
	}
	return enabledModels, nil
}

// ListAllModels 列出所有模型（包括禁用的）
func (p *MockProvider) ListAllModels(ctx context.Context) (map[string]*ModelConfig, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	allModels := make(map[string]*ModelConfig)
	for name, model := range p.models {
		// 创建副本以避免并发修改
		allModels[name] = &ModelConfig{
			Name:        model.Name,
			DisplayName: model.DisplayName,
			MaxTokens:   model.MaxTokens,
			Temperature: model.Temperature,
			TopP:        model.TopP,
			Enabled:     model.Enabled,
		}
	}
	return allModels, nil
}

// GetModelConfig 获取模型配置
func (p *MockProvider) GetModelConfig(name string) (*ModelConfig, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	if model, exists := p.models[name]; exists {
		// 返回副本以避免并发修改
		return &ModelConfig{
			Name:        model.Name,
			DisplayName: model.DisplayName,
			MaxTokens:   model.MaxTokens,
			Temperature: model.Temperature,
			TopP:        model.TopP,
			Enabled:     model.Enabled,
		}, nil
	}
	return nil, fmt.Errorf("model %s not found", name)
}

// EnableModel 启用模型
func (p *MockProvider) EnableModel(name string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if model, exists := p.models[name]; exists {
		model.Enabled = true
		return nil
	}
	return fmt.Errorf("model %s not found", name)
}

// DisableModel 禁用模型
func (p *MockProvider) DisableModel(name string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	if model, exists := p.models[name]; exists {
		model.Enabled = false
		return nil
	}
	return fmt.Errorf("model %s not found", name)
}

// ValidateAPIKey 验证API密钥
func (p *MockProvider) ValidateAPIKey(ctx context.Context) error {
	return nil // 模拟验证成功
}

// SetAPIKey 设置API密钥
func (p *MockProvider) SetAPIKey(key string) error {
	return nil // 模拟设置成功
}

// IsHealthy 检查健康状态
func (p *MockProvider) IsHealthy(ctx context.Context) bool {
	return true // 模拟健康
}