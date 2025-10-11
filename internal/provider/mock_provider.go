package provider

import (
	"context"
	"fmt"
	"io"
	"strings"
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
		"mock-gpt-3.5-turbo": {
			Name:        "mock-gpt-3.5-turbo",
			DisplayName: "Mock GPT-3.5 Turbo",
			MaxTokens:   4096,
			Temperature: 0.7,
			TopP:        0.9,
			Enabled:     true,
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
	// 检查是否有系统消息包含工具信息
	hasToolInfo := false
	userMessage := ""
	
	for _, msg := range req.Messages {
		if msg.Role == "system" {
			if strings.Contains(msg.Content, "stock_analysis") || 
				strings.Contains(msg.Content, "stock_compare") || 
				strings.Contains(msg.Content, "stock_advice") {
				hasToolInfo = true
			}
		}
		if msg.Role == "user" {
			userMessage = msg.Content
		}
	}
	
	var responseContent string
	
	// 如果有工具信息且用户询问股票相关问题，返回相应的工具调用
	if hasToolInfo {
		lowerMessage := strings.ToLower(userMessage)
		
		// 检测股票比较请求（优先级最高）
		if strings.Contains(lowerMessage, "比较") || strings.Contains(lowerMessage, "compare") ||
		   strings.Contains(lowerMessage, "vs") || strings.Contains(lowerMessage, "对比") {
			
			responseContent = `我来为您比较这两只股票。

<tool_call>
{
  "name": "stock_compare",
  "arguments": {
    "symbol1": "AAPL",
    "symbol2": "TSLA",
    "comparison_type": "comprehensive"
  }
}
</tool_call>`
		
		// 检测股票分析请求
		} else if strings.Contains(lowerMessage, "分析") || strings.Contains(lowerMessage, "analysis") ||
		   strings.Contains(lowerMessage, "aapl") || strings.Contains(lowerMessage, "苹果") ||
		   strings.Contains(lowerMessage, "tsla") || strings.Contains(lowerMessage, "特斯拉") {
			
			symbol := "AAPL"
			if strings.Contains(lowerMessage, "tsla") || strings.Contains(lowerMessage, "特斯拉") {
				symbol = "TSLA"
			}
			
			responseContent = fmt.Sprintf(`我来为您分析%s的股票。

<tool_call>
{
  "name": "stock_analysis",
  "arguments": {
    "symbol": "%s",
    "analysis_type": "comprehensive"
  }
}
</tool_call>`, symbol, symbol)
		
		// 检测投资建议请求
		} else if strings.Contains(lowerMessage, "建议") || strings.Contains(lowerMessage, "advice") ||
				  strings.Contains(lowerMessage, "推荐") || strings.Contains(lowerMessage, "投资") {
			
			symbol := "AAPL"
			if strings.Contains(lowerMessage, "tsla") || strings.Contains(lowerMessage, "特斯拉") {
				symbol = "TSLA"
			}
			
			responseContent = fmt.Sprintf(`我来为您提供%s的投资建议。

<tool_call>
{
  "name": "stock_advice",
  "arguments": {
    "symbol": "%s",
    "risk_tolerance": "moderate",
    "investment_horizon": "long_term"
  }
}
</tool_call>`, symbol, symbol)
		
		// 通用股票查询
		} else if strings.Contains(lowerMessage, "股票") || strings.Contains(lowerMessage, "stock") {
			responseContent = `我来为您分析股票。

<tool_call>
{
  "name": "stock_analysis",
  "arguments": {
    "symbol": "AAPL",
    "analysis_type": "comprehensive"
  }
}
</tool_call>`
		} else {
			// 普通响应
			responseContent = fmt.Sprintf("这是来自 %s 提供商的模拟响应，当前使用的模型是: %s。您的消息是: %s", p.name, req.Model, userMessage)
		}
	} else {
		// 普通响应
		responseContent = fmt.Sprintf("这是来自 %s 提供商的模拟响应，当前使用的模型是: %s。您的消息是: %s", p.name, req.Model, userMessage)
	}
	
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
					Content: responseContent,
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