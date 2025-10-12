package service

import (
	"context"
	"fmt"

	"go-springAi/internal/errors"
	"go-springAi/internal/logger"
)

// ProviderClient 提供商客户端接口
type ProviderClient interface {
	ValidateAPIKey(ctx context.Context) error
}

// ProviderKeyManager 提供商密钥管理器接口
type ProviderKeyManager interface {
	SetAPIKey(key string) error
	GetAPIKey() (string, error)
	ValidateKey(key string) error
}

// ProviderModelManager 提供商模型管理器接口
type ProviderModelManager interface {
	GetModel(name string) (interface{}, error)
	ListModels() map[string]interface{}
	UpdateModel(name string, config interface{}) error
	EnableModel(name string) error
	DisableModel(name string) error
}

// BaseProviderService 基础提供商服务
type BaseProviderService struct {
	providerName string
	client       ProviderClient
	keyManager   ProviderKeyManager
	modelManager ProviderModelManager
	logger       logger.Logger
}

// NewBaseProviderService 创建基础提供商服务
func NewBaseProviderService(
	providerName string,
	client ProviderClient,
	keyManager ProviderKeyManager,
	modelManager ProviderModelManager,
	log logger.Logger,
) *BaseProviderService {
	return &BaseProviderService{
		providerName: providerName,
		client:       client,
		keyManager:   keyManager,
		modelManager: modelManager,
		logger:       log,
	}
}

// ValidateAPIKey 验证API密钥的通用实现
func (s *BaseProviderService) ValidateAPIKey(ctx context.Context) error {
	s.logger.Info(fmt.Sprintf("Validating %s API key", s.providerName))
	
	err := s.client.ValidateAPIKey(ctx)
	if err != nil {
		s.logger.Error("API key validation failed", logger.ZapError(err))
		return errors.APIValidationFailed(s.providerName, err)
	}
	
	s.logger.Info("API key validation successful")
	return nil
}

// SetAPIKey 设置API密钥的通用实现
func (s *BaseProviderService) SetAPIKey(key string) error {
	s.logger.Info(fmt.Sprintf("Setting %s API key", s.providerName))
	
	err := s.keyManager.SetAPIKey(key)
	if err != nil {
		s.logger.Error("Failed to set API key", logger.ZapError(err))
		return errors.FailedToOperation("set API key", err)
	}
	
	s.logger.Info("API key set successfully")
	return nil
}

// GetAPIKey 获取API密钥的通用实现
func (s *BaseProviderService) GetAPIKey() (string, error) {
	s.logger.Debug(fmt.Sprintf("Getting %s API key", s.providerName))
	
	key, err := s.keyManager.GetAPIKey()
	if err != nil {
		s.logger.Error("Failed to get API key", logger.ZapError(err))
		return "", errors.FailedToGet("API key", err)
	}
	
	return key, nil
}

// UpdateModelConfig 更新模型配置的通用实现
func (s *BaseProviderService) UpdateModelConfig(name string, config interface{}) error {
	s.logger.Info(fmt.Sprintf("Updating %s model config", s.providerName), 
		logger.String("model", name))
	
	err := s.modelManager.UpdateModel(name, config)
	if err != nil {
		s.logger.Error("Failed to update model config", 
			logger.String("model", name), 
			logger.ZapError(err))
		return errors.FailedToUpdate("model config", err)
	}
	
	s.logger.Info("Model config updated successfully", logger.String("model", name))
	return nil
}

// EnableModel 启用模型的通用实现
func (s *BaseProviderService) EnableModel(name string) error {
	s.logger.Info(fmt.Sprintf("Enabling %s model", s.providerName), 
		logger.String("model", name))
	
	err := s.modelManager.EnableModel(name)
	if err != nil {
		s.logger.Error("Failed to enable model", 
			logger.String("model", name), 
			logger.ZapError(err))
		return errors.FailedToOperation("enable model", err)
	}
	
	s.logger.Info("Model enabled successfully", logger.String("model", name))
	return nil
}

// DisableModel 禁用模型的通用实现
func (s *BaseProviderService) DisableModel(name string) error {
	s.logger.Info(fmt.Sprintf("Disabling %s model", s.providerName), 
		logger.String("model", name))
	
	err := s.modelManager.DisableModel(name)
	if err != nil {
		s.logger.Error("Failed to disable model", 
			logger.String("model", name), 
			logger.ZapError(err))
		return errors.FailedToOperation("disable model", err)
	}
	
	s.logger.Info("Model disabled successfully", logger.String("model", name))
	return nil
}

// GetModelConfig 获取模型配置的通用实现
func (s *BaseProviderService) GetModelConfig(name string) (interface{}, error) {
	s.logger.Debug(fmt.Sprintf("Getting %s model config", s.providerName), 
		logger.String("model", name))
	
	config, err := s.modelManager.GetModel(name)
	if err != nil {
		s.logger.Error("Failed to get model config", 
			logger.String("model", name), 
			logger.ZapError(err))
		return nil, errors.FailedToGet("model config", err)
	}
	
	return config, nil
}

// ListModels 列出模型的通用实现
func (s *BaseProviderService) ListModels() map[string]interface{} {
	s.logger.Debug(fmt.Sprintf("Listing %s models", s.providerName))
	
	models := s.modelManager.ListModels()
	s.logger.Info(fmt.Sprintf("Listed %s models", s.providerName), 
		logger.Int("count", len(models)))
	
	return models
}

// LogChatCompletion 记录聊天完成请求的通用方法
func (s *BaseProviderService) LogChatCompletion(model string, messageCount int) {
	s.logger.Info(fmt.Sprintf("Processing %s chat completion", s.providerName),
		logger.String("model", model),
		logger.Int("message_count", messageCount))
}

// LogChatCompletionSuccess 记录聊天完成成功的通用方法
func (s *BaseProviderService) LogChatCompletionSuccess(model string, responseID string, usage interface{}) {
	s.logger.Info(fmt.Sprintf("%s chat completion successful", s.providerName),
		logger.String("model", model),
		logger.String("response_id", responseID))
}

// LogChatCompletionError 记录聊天完成错误的通用方法
func (s *BaseProviderService) LogChatCompletionError(model string, err error) {
	s.logger.Error(fmt.Sprintf("%s chat completion failed", s.providerName),
		logger.String("model", model),
		logger.ZapError(err))
}

// LogStreamCompletion 记录流式完成请求的通用方法
func (s *BaseProviderService) LogStreamCompletion(model string, messageCount int) {
	s.logger.Info(fmt.Sprintf("Processing %s stream completion", s.providerName),
		logger.String("model", model),
		logger.Int("message_count", messageCount))
}

// LogStreamCompletionSuccess 记录流式完成成功的通用方法
func (s *BaseProviderService) LogStreamCompletionSuccess(model string) {
	s.logger.Info(fmt.Sprintf("%s stream completion initiated", s.providerName),
		logger.String("model", model))
}

// LogStreamCompletionError 记录流式完成错误的通用方法
func (s *BaseProviderService) LogStreamCompletionError(model string, err error) {
	s.logger.Error(fmt.Sprintf("%s stream completion failed", s.providerName),
		logger.String("model", model),
		logger.ZapError(err))
}