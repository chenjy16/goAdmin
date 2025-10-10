package service

import (
	"context"
	"fmt"

	"go-springAi/internal/repository"
	"go-springAi/internal/database/generated/api_keys"
)

// APIKeyService API密钥服务接口
type APIKeyService interface {
	// SetAPIKey 设置用户的API密钥
	SetAPIKey(ctx context.Context, userID int64, providerType, apiKey string) error
	
	// GetAPIKey 获取用户的API密钥
	GetAPIKey(ctx context.Context, userID int64, providerType string) (string, error)
	
	// ValidateAPIKey 验证API密钥格式
	ValidateAPIKey(providerType, apiKey string) error
	
	// ListUserAPIKeys 获取用户的所有API密钥
	ListUserAPIKeys(ctx context.Context, userID int64) ([]api_keys.ApiKey, error)
	
	// DeactivateAPIKey 停用API密钥
	DeactivateAPIKey(ctx context.Context, userID int64, providerType string) error
	
	// DeleteAPIKey 删除API密钥
	DeleteAPIKey(ctx context.Context, userID int64, providerType string) error
	
	// CheckAPIKeyExists 检查API密钥是否存在
	CheckAPIKeyExists(ctx context.Context, userID int64, providerType string) (bool, error)
	
	// GetMaskedAPIKey 获取脱敏的API密钥
	GetMaskedAPIKey(ctx context.Context, userID int64, providerType string) (string, error)
	
	// GetKeyManager 获取密钥管理器
	GetKeyManager(userID int64, providerType string) *DatabaseKeyManager
}

// apiKeyService API密钥服务实现
type apiKeyService struct {
	repo repository.APIKeyRepository
}

// NewAPIKeyService 创建新的API密钥服务
func NewAPIKeyService(repo repository.APIKeyRepository) APIKeyService {
	return &apiKeyService{
		repo: repo,
	}
}

// SetAPIKey 设置用户的API密钥
func (s *apiKeyService) SetAPIKey(ctx context.Context, userID int64, providerType, apiKey string) error {
	// 验证API密钥格式
	if err := s.ValidateAPIKey(providerType, apiKey); err != nil {
		return fmt.Errorf("invalid API key: %w", err)
	}
	
	// 创建密钥管理器
	keyManager := NewDatabaseKeyManager(userID, providerType, s.repo)
	
	// 设置API密钥
	return keyManager.SetAPIKey(apiKey)
}

// GetAPIKey 获取用户的API密钥
func (s *apiKeyService) GetAPIKey(ctx context.Context, userID int64, providerType string) (string, error) {
	// 创建密钥管理器
	keyManager := NewDatabaseKeyManager(userID, providerType, s.repo)
	
	// 获取API密钥
	return keyManager.GetAPIKey()
}

// ValidateAPIKey 验证API密钥格式
func (s *apiKeyService) ValidateAPIKey(providerType, apiKey string) error {
	if apiKey == "" {
		return fmt.Errorf("API key is empty")
	}
	
	// 根据提供商类型进行不同的验证
	switch providerType {
	case "openai":
		if len(apiKey) < 10 {
			return fmt.Errorf("OpenAI API key is too short")
		}
		// OpenAI API 密钥通常以 "sk-" 开头
		if len(apiKey) > 4 && apiKey[:3] != "sk-" {
			return fmt.Errorf("OpenAI API key should start with 'sk-'")
		}
	case "googleai":
		if len(apiKey) < 10 {
			return fmt.Errorf("Google AI API key is too short")
		}
		// Google AI API 密钥通常以 "AIza" 开头
		if len(apiKey) > 4 && apiKey[:4] != "AIza" {
			return fmt.Errorf("Google AI API key should start with 'AIza'")
		}
	case "mock":
		// Mock provider 允许任何格式的密钥
		if len(apiKey) < 1 {
			return fmt.Errorf("Mock API key cannot be empty")
		}
	default:
		// 通用验证
		if len(apiKey) < 10 {
			return fmt.Errorf("API key is too short")
		}
	}
	
	return nil
}

// ListUserAPIKeys 获取用户的所有API密钥
func (s *apiKeyService) ListUserAPIKeys(ctx context.Context, userID int64) ([]api_keys.ApiKey, error) {
	return s.repo.ListAPIKeysByUser(ctx, userID)
}

// DeactivateAPIKey 停用API密钥
func (s *apiKeyService) DeactivateAPIKey(ctx context.Context, userID int64, providerType string) error {
	return s.repo.DeactivateAPIKey(ctx, userID, providerType)
}

// DeleteAPIKey 删除API密钥
func (s *apiKeyService) DeleteAPIKey(ctx context.Context, userID int64, providerType string) error {
	return s.repo.DeleteAPIKey(ctx, userID, providerType)
}

// CheckAPIKeyExists 检查API密钥是否存在
func (s *apiKeyService) CheckAPIKeyExists(ctx context.Context, userID int64, providerType string) (bool, error) {
	return s.repo.CheckAPIKeyExists(ctx, userID, providerType)
}

// GetMaskedAPIKey 获取脱敏的API密钥
func (s *apiKeyService) GetMaskedAPIKey(ctx context.Context, userID int64, providerType string) (string, error) {
	// 获取完整的API密钥
	apiKey, err := s.GetAPIKey(ctx, userID, providerType)
	if err != nil {
		return "", err
	}
	
	// 对API密钥进行脱敏处理
	return maskAPIKey(apiKey), nil
}

// maskAPIKey 对API密钥进行脱敏处理
func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 8 {
		// 如果密钥太短，只显示前2位和后2位
		if len(apiKey) <= 4 {
			return "****"
		}
		return apiKey[:2] + "****" + apiKey[len(apiKey)-2:]
	}
	
	// 显示前4位和后4位，中间用星号替代
	return apiKey[:4] + "****" + apiKey[len(apiKey)-4:]
}

// GetKeyManager 获取密钥管理器
func (s *apiKeyService) GetKeyManager(userID int64, providerType string) *DatabaseKeyManager {
	return NewDatabaseKeyManager(userID, providerType, s.repo)
}