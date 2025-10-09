package service

import (
	"context"
	"fmt"
	"io"
	"time"

	"goMcp/internal/logger"
	"goMcp/internal/openai"
)

// OpenAIService OpenAI 服务
type OpenAIService struct {
	client       openai.Client
	keyManager   openai.KeyManager
	modelManager openai.ModelManager
	logger       logger.Logger
}

// NewOpenAIService 创建新的 OpenAI 服务
func NewOpenAIService(
	client openai.Client,
	keyManager openai.KeyManager,
	modelManager openai.ModelManager,
	log logger.Logger,
) *OpenAIService {
	return &OpenAIService{
		client:       client,
		keyManager:   keyManager,
		modelManager: modelManager,
		logger:       log,
	}
}

// ChatCompletionRequest 聊天完成请求
type ChatCompletionRequest struct {
	Model       string                `json:"model"`
	Messages    []openai.Message      `json:"messages"`
	MaxTokens   *int                  `json:"max_tokens,omitempty"`
	Temperature *float32              `json:"temperature,omitempty"`
	TopP        *float32              `json:"top_p,omitempty"`
	Stream      bool                  `json:"stream,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
}

// ChatCompletionResponse 聊天完成响应
type ChatCompletionResponse struct {
	ID      string           `json:"id"`
	Object  string           `json:"object"`
	Created int64            `json:"created"`
	Model   string           `json:"model"`
	Choices []openai.Choice  `json:"choices"`
	Usage   openai.Usage     `json:"usage"`
}

// ChatCompletion 聊天完成
func (s *OpenAIService) ChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	startTime := time.Now()
	
	// 记录请求日志
	s.logger.Info("OpenAI chat completion request",
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
	
	// 构建 OpenAI 请求
	openaiReq := &openai.ChatRequest{
		Model:    req.Model,
		Messages: req.Messages,
		Stream:   req.Stream,
	}
	
	// 应用模型配置
	s.applyModelConfig(openaiReq, modelConfig, req)
	
	// 调用 OpenAI API
	resp, err := s.client.ChatCompletion(ctx, openaiReq)
	if err != nil {
		s.logger.Error("OpenAI API error",
			logger.String("model", req.Model),
			logger.ZapError(err),
			logger.Duration("duration", time.Since(startTime)),
		)
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}
	
	// 记录成功日志
	s.logger.Info("OpenAI chat completion success",
		logger.String("model", req.Model),
		logger.String("response_id", resp.ID),
		logger.Int("prompt_tokens", resp.Usage.PromptTokens),
		logger.Int("completion_tokens", resp.Usage.CompletionTokens),
		logger.Int("total_tokens", resp.Usage.TotalTokens),
		logger.Duration("duration", time.Since(startTime)),
	)
	
	return &ChatCompletionResponse{
		ID:      resp.ID,
		Object:  resp.Object,
		Created: resp.Created,
		Model:   resp.Model,
		Choices: resp.Choices,
		Usage:   resp.Usage,
	}, nil
}

// ChatCompletionStream 流式聊天完成
func (s *OpenAIService) ChatCompletionStream(ctx context.Context, req *ChatCompletionRequest) (io.ReadCloser, error) {
	startTime := time.Now()
	
	// 记录请求日志
	s.logger.Info("OpenAI chat completion stream request",
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
	
	// 构建 OpenAI 请求
	openaiReq := &openai.ChatRequest{
		Model:    req.Model,
		Messages: req.Messages,
		Stream:   true,
	}
	
	// 应用模型配置
	s.applyModelConfig(openaiReq, modelConfig, req)
	
	// 调用 OpenAI API
	stream, err := s.client.ChatCompletionStream(ctx, openaiReq)
	if err != nil {
		s.logger.Error("OpenAI API stream error",
			logger.String("model", req.Model),
			logger.ZapError(err),
			logger.Duration("duration", time.Since(startTime)),
		)
		return nil, fmt.Errorf("OpenAI API stream error: %w", err)
	}
	
	// 记录流开始日志
	s.logger.Info("OpenAI chat completion stream started",
		logger.String("model", req.Model),
		logger.Duration("setup_duration", time.Since(startTime)),
	)
	
	return stream, nil
}

// ListModels 列出可用模型
func (s *OpenAIService) ListModels(ctx context.Context) (map[string]*openai.ModelConfig, error) {
	s.logger.Info("Listing OpenAI models")
	
	// 获取本地配置的模型
	models := s.modelManager.ListModels()
	
	// 过滤启用的模型
	enabledModels := make(map[string]*openai.ModelConfig)
	for name, model := range models {
		if model.Enabled {
			enabledModels[name] = model
		}
	}
	
	s.logger.Info("Listed OpenAI models", logger.Int("count", len(enabledModels)))
	return enabledModels, nil
}

// ValidateAPIKey 验证 API 密钥
func (s *OpenAIService) ValidateAPIKey(ctx context.Context) error {
	s.logger.Info("Validating OpenAI API key")
	
	err := s.client.ValidateAPIKey(ctx)
	if err != nil {
		s.logger.Error("API key validation failed", logger.ZapError(err))
		return fmt.Errorf("API key validation failed: %w", err)
	}
	
	s.logger.Info("API key validation successful")
	return nil
}

// SetAPIKey 设置 API 密钥
func (s *OpenAIService) SetAPIKey(key string) error {
	s.logger.Info("Setting OpenAI API key")
	
	err := s.keyManager.SetAPIKey(key)
	if err != nil {
		s.logger.Error("Failed to set API key", logger.ZapError(err))
		return fmt.Errorf("failed to set API key: %w", err)
	}
	
	s.logger.Info("API key set successfully")
	return nil
}

// GetModelConfig 获取模型配置
func (s *OpenAIService) GetModelConfig(name string) (*openai.ModelConfig, error) {
	return s.modelManager.GetModel(name)
}

// UpdateModelConfig 更新模型配置
func (s *OpenAIService) UpdateModelConfig(name string, config *openai.ModelConfig) error {
	s.logger.Info("Updating model config", logger.String("model", name))
	
	err := s.modelManager.UpdateModel(name, config)
	if err != nil {
		s.logger.Error("Failed to update model config",
			logger.String("model", name),
			logger.ZapError(err),
		)
		return fmt.Errorf("failed to update model config: %w", err)
	}
	
	s.logger.Info("Model config updated successfully", logger.String("model", name))
	return nil
}

// EnableModel 启用模型
func (s *OpenAIService) EnableModel(name string) error {
	s.logger.Info("Enabling model", logger.String("model", name))
	
	err := s.modelManager.EnableModel(name)
	if err != nil {
		s.logger.Error("Failed to enable model", logger.String("model", name), logger.ZapError(err))
		return fmt.Errorf("failed to enable model: %w", err)
	}
	
	s.logger.Info("Model enabled successfully", logger.String("model", name))
	return nil
}

// DisableModel 禁用模型
func (s *OpenAIService) DisableModel(name string) error {
	s.logger.Info("Disabling model", logger.String("model", name))
	
	err := s.modelManager.DisableModel(name)
	if err != nil {
		s.logger.Error("Failed to disable model", logger.String("model", name), logger.ZapError(err))
		return fmt.Errorf("failed to disable model: %w", err)
	}
	
	s.logger.Info("Model disabled successfully", logger.String("model", name))
	return nil
}

// applyModelConfig 应用模型配置到请求
func (s *OpenAIService) applyModelConfig(openaiReq *openai.ChatRequest, modelConfig *openai.ModelConfig, req *ChatCompletionRequest) {
	// 应用最大令牌数
	if req.MaxTokens != nil {
		openaiReq.MaxTokens = *req.MaxTokens
	} else {
		openaiReq.MaxTokens = modelConfig.MaxTokens
	}
	
	// 应用温度
	if req.Temperature != nil {
		openaiReq.Temperature = *req.Temperature
	} else {
		openaiReq.Temperature = modelConfig.Temperature
	}
	
	// 应用 TopP
	if req.TopP != nil {
		openaiReq.TopP = *req.TopP
	} else {
		openaiReq.TopP = modelConfig.TopP
	}
	
	// 应用频率惩罚
	openaiReq.FrequencyPenalty = modelConfig.FrequencyPenalty
	
	// 应用存在惩罚
	openaiReq.PresencePenalty = modelConfig.PresencePenalty
}