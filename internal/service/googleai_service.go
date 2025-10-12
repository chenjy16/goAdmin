package service

import (
	"context"
	"fmt"
	"io"
	"time"

	"go-springAi/internal/googleai"
	"go-springAi/internal/logger"
)

// GoogleAIService Google AI 服务
type GoogleAIService struct {
	*BaseProviderService
	client       googleai.Client
	keyManager   googleai.KeyManager
	modelManager googleai.ModelManager
}

// NewGoogleAIService 创建新的 Google AI 服务
func NewGoogleAIService(
	client googleai.Client,
	keyManager googleai.KeyManager,
	modelManager googleai.ModelManager,
	log logger.Logger,
) *GoogleAIService {
	// 创建适配器
	keyAdapter := &googleaiKeyManagerAdapter{keyManager}
	modelAdapter := &googleaiModelManagerAdapter{modelManager}
	
	baseService := NewBaseProviderService("googleai", client, keyAdapter, modelAdapter, log)
	return &GoogleAIService{
		BaseProviderService: baseService,
		client:              client,
		keyManager:          keyManager,
		modelManager:        modelManager,
	}
}

// GoogleAIChatCompletionRequest Google AI 聊天完成请求
type GoogleAIChatCompletionRequest struct {
	Model       string                     `json:"model"`
	Messages    []googleai.Message         `json:"messages"`
	MaxTokens   *int                       `json:"max_tokens,omitempty"`
	Temperature *float32                   `json:"temperature,omitempty"`
	TopP        *float32                   `json:"top_p,omitempty"`
	TopK        *int                       `json:"top_k,omitempty"`
	Stream      bool                       `json:"stream,omitempty"`
	Options     map[string]interface{}     `json:"options,omitempty"`
}

// GoogleAIChatCompletionResponse Google AI 聊天完成响应
type GoogleAIChatCompletionResponse struct {
	ID      string              `json:"id"`
	Object  string              `json:"object"`
	Created int64               `json:"created"`
	Model   string              `json:"model"`
	Choices []googleai.Choice   `json:"choices"`
	Usage   googleai.Usage      `json:"usage"`
}

// ChatCompletion 聊天完成
func (s *GoogleAIService) ChatCompletion(ctx context.Context, req *GoogleAIChatCompletionRequest) (*GoogleAIChatCompletionResponse, error) {
	startTime := time.Now()
	
	// 记录请求日志
	s.logger.Info("Google AI chat completion request",
		logger.String("model", req.Model),
		logger.Int("message_count", len(req.Messages)),
		logger.Bool("stream", req.Stream),
	)
	
	// 验证模型
	modelConfig, err := s.modelManager.GetModel(req.Model)
	if err != nil {
		s.logger.Error("Invalid model", logger.String("model", req.Model), logger.ZapError(err))
		return nil, fmt.Errorf("invalid model: %w", err)
	}
	
	if !modelConfig.Enabled {
		s.logger.Error("Model disabled", logger.String("model", req.Model))
		return nil, fmt.Errorf("model %s is disabled", req.Model)
	}
	
	// 构建 Google AI 请求
	googleaiReq := &googleai.ChatRequest{
		Model:    req.Model,
		Messages: req.Messages,
		Stream:   req.Stream,
	}
	
	// 应用模型配置
	s.applyModelConfig(googleaiReq, modelConfig, req)
	
	// 调用 Google AI API
	resp, err := s.client.ChatCompletion(ctx, googleaiReq)
	if err != nil {
		s.logger.Error("Google AI API error",
			logger.String("model", req.Model),
			logger.ZapError(err),
			logger.Duration("duration", time.Since(startTime)),
		)
		return nil, fmt.Errorf("google AI API error: %w", err)
	}
	
	// 记录成功日志
	s.logger.Info("Google AI chat completion success",
		logger.String("model", req.Model),
		logger.String("response_id", resp.ID),
		logger.Int("prompt_tokens", resp.Usage.PromptTokens),
		logger.Int("completion_tokens", resp.Usage.CompletionTokens),
		logger.Int("total_tokens", resp.Usage.TotalTokens),
		logger.Duration("duration", time.Since(startTime)),
	)
	
	return &GoogleAIChatCompletionResponse{
		ID:      resp.ID,
		Object:  resp.Object,
		Created: resp.Created,
		Model:   resp.Model,
		Choices: resp.Choices,
		Usage:   resp.Usage,
	}, nil
}

// ChatCompletionStream 流式聊天完成
func (s *GoogleAIService) ChatCompletionStream(ctx context.Context, req *GoogleAIChatCompletionRequest) (io.ReadCloser, error) {
	startTime := time.Now()
	
	// 记录请求日志
	s.logger.Info("Google AI chat completion stream request",
		logger.String("model", req.Model),
		logger.Int("message_count", len(req.Messages)),
	)
	
	// 验证模型
	modelConfig, err := s.modelManager.GetModel(req.Model)
	if err != nil {
		s.logger.Error("Invalid model", logger.String("model", req.Model), logger.ZapError(err))
		return nil, fmt.Errorf("invalid model: %w", err)
	}
	
	if !modelConfig.Enabled {
		s.logger.Error("Model disabled", logger.String("model", req.Model))
		return nil, fmt.Errorf("model %s is disabled", req.Model)
	}
	
	// 构建 Google AI 请求
	googleaiReq := &googleai.ChatRequest{
		Model:    req.Model,
		Messages: req.Messages,
		Stream:   true,
	}
	
	// 应用模型配置
	s.applyModelConfig(googleaiReq, modelConfig, req)
	
	// 调用 Google AI API
	stream, err := s.client.ChatCompletionStream(ctx, googleaiReq)
	if err != nil {
		s.logger.Error("Google AI API stream error",
			logger.String("model", req.Model),
			logger.ZapError(err),
			logger.Duration("duration", time.Since(startTime)),
		)
		return nil, fmt.Errorf("google AI API stream error: %w", err)
	}
	
	// 记录流开始日志
	s.logger.Info("Google AI chat completion stream started",
		logger.String("model", req.Model),
		logger.Duration("setup_duration", time.Since(startTime)),
	)
	
	return stream, nil
}

// ListModels 列出可用模型（仅启用的）
func (s *GoogleAIService) ListModels(ctx context.Context) (map[string]*googleai.ModelConfig, error) {
	s.logger.Info("Listing Google AI models")
	
	// 获取本地配置的模型
	models := s.modelManager.ListModels()
	
	// 过滤启用的模型
	enabledModels := make(map[string]*googleai.ModelConfig)
	for name, model := range models {
		if model.Enabled {
			enabledModels[name] = model
		}
	}
	
	s.logger.Info("Listed Google AI models", logger.Int("count", len(enabledModels)))
	return enabledModels, nil
}

// ListAllModels 列出所有模型（包括禁用的）
func (s *GoogleAIService) ListAllModels(ctx context.Context) (map[string]*googleai.ModelConfig, error) {
	s.logger.Info("Listing all Google AI models")
	
	// 获取本地配置的所有模型
	models := s.modelManager.ListModels()
	
	s.logger.Info("Listed all Google AI models", logger.Int("count", len(models)))
	return models, nil
}

// GetModelConfig 获取模型配置 (类型安全的包装方法)
func (s *GoogleAIService) GetModelConfig(name string) (*googleai.ModelConfig, error) {
	return s.modelManager.GetModel(name)
}

// UpdateModelConfig 更新模型配置 (类型安全的包装方法)
func (s *GoogleAIService) UpdateModelConfig(name string, config *googleai.ModelConfig) error {
	return s.modelManager.UpdateModel(name, config)
}

// applyModelConfig 应用模型配置到请求
func (s *GoogleAIService) applyModelConfig(googleaiReq *googleai.ChatRequest, modelConfig *googleai.ModelConfig, req *GoogleAIChatCompletionRequest) {
	// 应用最大令牌数
	if req.MaxTokens != nil {
		googleaiReq.MaxTokens = *req.MaxTokens
	} else {
		googleaiReq.MaxTokens = modelConfig.MaxTokens
	}
	
	// 应用温度
	if req.Temperature != nil {
		googleaiReq.Temperature = *req.Temperature
	} else {
		googleaiReq.Temperature = modelConfig.Temperature
	}
	
	// 应用 TopP
	if req.TopP != nil {
		googleaiReq.TopP = *req.TopP
	} else {
		googleaiReq.TopP = modelConfig.TopP
	}
	
	// 应用 TopK
	if req.TopK != nil {
		googleaiReq.TopK = *req.TopK
	} else {
		googleaiReq.TopK = modelConfig.TopK
	}
}

// googleaiKeyManagerAdapter 适配器，将 googleai.KeyManager 适配为 ProviderKeyManager
type googleaiKeyManagerAdapter struct {
	googleai.KeyManager
}

// googleaiModelManagerAdapter 适配器，将 googleai.ModelManager 适配为 ProviderModelManager
type googleaiModelManagerAdapter struct {
	googleai.ModelManager
}

// GetModel 实现 ProviderModelManager 接口
func (a *googleaiModelManagerAdapter) GetModel(name string) (interface{}, error) {
	return a.ModelManager.GetModel(name)
}

// ListModels 实现 ProviderModelManager 接口
func (a *googleaiModelManagerAdapter) ListModels() map[string]interface{} {
	models := a.ModelManager.ListModels()
	result := make(map[string]interface{})
	for k, v := range models {
		result[k] = v
	}
	return result
}

// UpdateModel 实现 ProviderModelManager 接口
func (a *googleaiModelManagerAdapter) UpdateModel(name string, config interface{}) error {
	if googleaiConfig, ok := config.(*googleai.ModelConfig); ok {
		return a.ModelManager.UpdateModel(name, googleaiConfig)
	}
	return fmt.Errorf("invalid config type for GoogleAI model")
}