package provider

import (
	"context"
	"io"

	"goMcp/internal/googleai"
	"goMcp/internal/service"
)

// GoogleAIProvider GoogleAI提供商实现
type GoogleAIProvider struct {
	service *service.GoogleAIService
}

// NewGoogleAIProvider 创建GoogleAI Provider
func NewGoogleAIProvider(service *service.GoogleAIService) *GoogleAIProvider {
	return &GoogleAIProvider{
		service: service,
	}
}

// GetType 获取提供商类型
func (p *GoogleAIProvider) GetType() ProviderType {
	return ProviderTypeGoogleAI
}

// GetName 获取提供商名称
func (p *GoogleAIProvider) GetName() string {
	return "Google AI"
}

// ChatCompletion 聊天完成
func (p *GoogleAIProvider) ChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// 转换统一请求为GoogleAI特定请求
	googleaiReq := &service.GoogleAIChatCompletionRequest{
		Model:       req.Model,
		Messages:    convertToGoogleAIMessages(req.Messages),
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		TopK:        req.TopK,
		Stream:      req.Stream,
		Options:     req.Options,
	}
	
	// 调用GoogleAI服务
	resp, err := p.service.ChatCompletion(ctx, googleaiReq)
	if err != nil {
		return nil, err
	}
	
	// 转换GoogleAI响应为统一响应
	return &ChatResponse{
		ID:      resp.ID,
		Object:  resp.Object,
		Created: resp.Created,
		Model:   resp.Model,
		Choices: convertFromGoogleAIChoices(resp.Choices),
		Usage: Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
	}, nil
}

// ChatCompletionStream 流式聊天完成
func (p *GoogleAIProvider) ChatCompletionStream(ctx context.Context, req *ChatRequest) (io.ReadCloser, error) {
	// 转换统一请求为GoogleAI特定请求
	googleaiReq := &service.GoogleAIChatCompletionRequest{
		Model:       req.Model,
		Messages:    convertToGoogleAIMessages(req.Messages),
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		TopK:        req.TopK,
		Stream:      true,
		Options:     req.Options,
	}
	
	// 调用GoogleAI服务
	return p.service.ChatCompletionStream(ctx, googleaiReq)
}

// ListModels 列出可用模型
func (p *GoogleAIProvider) ListModels(ctx context.Context) (map[string]*ModelConfig, error) {
	models, err := p.service.ListModels(ctx)
	if err != nil {
		return nil, err
	}
	
	// 转换GoogleAI模型配置为统一模型配置
	result := make(map[string]*ModelConfig)
	for name, config := range models {
		result[name] = &ModelConfig{
			Name:        config.Name,
			DisplayName: config.DisplayName,
			MaxTokens:   config.MaxTokens,
			Temperature: config.Temperature,
			TopP:        config.TopP,
			TopK:        config.TopK,
			Enabled:     config.Enabled,
		}
	}
	
	return result, nil
}

// GetModelConfig 获取模型配置
func (p *GoogleAIProvider) GetModelConfig(name string) (*ModelConfig, error) {
	config, err := p.service.GetModelConfig(name)
	if err != nil {
		return nil, err
	}
	
	return &ModelConfig{
		Name:        config.Name,
		DisplayName: config.DisplayName,
		MaxTokens:   config.MaxTokens,
		Temperature: config.Temperature,
		TopP:        config.TopP,
		TopK:        config.TopK,
		Enabled:     config.Enabled,
	}, nil
}

// EnableModel 启用模型
func (p *GoogleAIProvider) EnableModel(name string) error {
	return p.service.EnableModel(name)
}

// DisableModel 禁用模型
func (p *GoogleAIProvider) DisableModel(name string) error {
	return p.service.DisableModel(name)
}

// ValidateAPIKey 验证API密钥
func (p *GoogleAIProvider) ValidateAPIKey(ctx context.Context) error {
	return p.service.ValidateAPIKey(ctx)
}

// SetAPIKey 设置API密钥
func (p *GoogleAIProvider) SetAPIKey(key string) error {
	return p.service.SetAPIKey(key)
}

// IsHealthy 检查提供商健康状态
func (p *GoogleAIProvider) IsHealthy(ctx context.Context) bool {
	err := p.service.ValidateAPIKey(ctx)
	return err == nil
}

// 辅助函数：转换统一消息为GoogleAI消息
func convertToGoogleAIMessages(messages []Message) []googleai.Message {
	result := make([]googleai.Message, len(messages))
	for i, msg := range messages {
		result[i] = googleai.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	return result
}

// 辅助函数：转换GoogleAI选择为统一选择
func convertFromGoogleAIChoices(choices []googleai.Choice) []Choice {
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