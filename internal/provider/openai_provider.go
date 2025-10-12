package provider

import (
	"context"
	"io"

	"go-springAi/internal/openai"
	"go-springAi/internal/service"
	"go-springAi/internal/types"
)

// OpenAIProvider OpenAI提供商实现
type OpenAIProvider struct {
	service *service.OpenAIService
}

// NewOpenAIProvider 创建OpenAI Provider
func NewOpenAIProvider(service *service.OpenAIService) *OpenAIProvider {
	return &OpenAIProvider{
		service: service,
	}
}

// GetType 获取提供商类型
func (p *OpenAIProvider) GetType() ProviderType {
	return types.ProviderTypeOpenAI
}

// GetName 获取提供商名称
func (p *OpenAIProvider) GetName() string {
	return "OpenAI"
}

// ChatCompletion 聊天完成
func (p *OpenAIProvider) ChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// 转换统一请求为OpenAI特定请求
	openaiReq := &service.ChatCompletionRequest{
		Model:       req.Model,
		Messages:    convertToOpenAIMessages(req.Messages),
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Stream:      req.Stream,
		Options:     req.Options,
	}
	
	// 调用OpenAI服务
	resp, err := p.service.ChatCompletion(ctx, openaiReq)
	if err != nil {
		return nil, err
	}
	
	// 转换OpenAI响应为统一响应
	return &ChatResponse{
		ID:      resp.ID,
		Object:  resp.Object,
		Created: resp.Created,
		Model:   resp.Model,
		Choices: convertFromOpenAIChoices(resp.Choices),
		Usage: Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
	}, nil
}

// ChatCompletionStream 流式聊天完成
func (p *OpenAIProvider) ChatCompletionStream(ctx context.Context, req *ChatRequest) (io.ReadCloser, error) {
	// 转换统一请求为OpenAI特定请求
	openaiReq := &service.ChatCompletionRequest{
		Model:       req.Model,
		Messages:    convertToOpenAIMessages(req.Messages),
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Stream:      true,
		Options:     req.Options,
	}
	
	// 调用OpenAI服务
	return p.service.ChatCompletionStream(ctx, openaiReq)
}

// ListModels 列出可用模型（仅启用的）
func (p *OpenAIProvider) ListModels(ctx context.Context) (map[string]*ModelConfig, error) {
	models, err := p.service.ListModels(ctx)
	if err != nil {
		return nil, err
	}
	
	// 转换OpenAI模型配置为统一模型配置
	result := make(map[string]*ModelConfig)
	for name, config := range models {
		result[name] = &ModelConfig{
			Name:        config.Name,
			DisplayName: config.Name, // OpenAI使用Name作为显示名称
			MaxTokens:   config.MaxTokens,
			Temperature: config.Temperature,
			TopP:        config.TopP,
			Enabled:     config.Enabled,
		}
	}
	
	return result, nil
}

// ListAllModels 列出所有模型（包括禁用的）
func (p *OpenAIProvider) ListAllModels(ctx context.Context) (map[string]*ModelConfig, error) {
	models, err := p.service.ListAllModels(ctx)
	if err != nil {
		return nil, err
	}
	
	// 转换OpenAI模型配置为统一模型配置
	result := make(map[string]*ModelConfig)
	for name, config := range models {
		result[name] = &ModelConfig{
			Name:        config.Name,
			DisplayName: config.Name, // OpenAI使用Name作为显示名称
			MaxTokens:   config.MaxTokens,
			Temperature: config.Temperature,
			TopP:        config.TopP,
			Enabled:     config.Enabled,
		}
	}
	
	return result, nil
}

// GetModelConfig 获取模型配置
func (p *OpenAIProvider) GetModelConfig(name string) (*ModelConfig, error) {
	config, err := p.service.GetModelConfig(name)
	if err != nil {
		return nil, err
	}
	
	return &ModelConfig{
		Name:        config.Name,
		DisplayName: config.Name, // OpenAI使用Name作为显示名称
		MaxTokens:   config.MaxTokens,
		Temperature: config.Temperature,
		TopP:        config.TopP,
		Enabled:     config.Enabled,
	}, nil
}

// EnableModel 启用模型
func (p *OpenAIProvider) EnableModel(name string) error {
	return p.service.EnableModel(name)
}

// DisableModel 禁用模型
func (p *OpenAIProvider) DisableModel(name string) error {
	return p.service.DisableModel(name)
}

// ValidateAPIKey 验证API密钥
func (p *OpenAIProvider) ValidateAPIKey(ctx context.Context) error {
	return p.service.ValidateAPIKey(ctx)
}

// SetAPIKey 设置API密钥
func (p *OpenAIProvider) SetAPIKey(key string) error {
	return p.service.SetAPIKey(key)
}

// IsHealthy 检查提供商健康状态
func (p *OpenAIProvider) IsHealthy(ctx context.Context) bool {
	err := p.service.ValidateAPIKey(ctx)
	return err == nil
}

// 辅助函数：转换统一消息为OpenAI消息
func convertToOpenAIMessages(messages []Message) []openai.Message {
	result := make([]openai.Message, len(messages))
	for i, msg := range messages {
		result[i] = openai.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	return result
}

// 辅助函数：转换OpenAI选择为统一选择
func convertFromOpenAIChoices(choices []openai.Choice) []Choice {
	result := make([]Choice, len(choices))
	for i, choice := range choices {
		result[i] = Choice{
			Index: choice.Index,
			Message: Message{
				Role:    choice.Message.Role,
				Content: choice.Message.Content,
			},
			FinishReason: choice.FinishReason,
		}
	}
	return result
}