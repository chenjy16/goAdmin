package repository

import (
	"context"

	"go-springAi/internal/database/generated/api_keys"
)

// APIKeyRepository API密钥数据访问层接口
type APIKeyRepository interface {
	// CreateAPIKey 创建API密钥
	CreateAPIKey(ctx context.Context, params CreateAPIKeyParams) (*api_keys.ApiKey, error)
	
	// GetAPIKey 获取指定用户和提供商的API密钥
	GetAPIKey(ctx context.Context, userID int64, providerType string) (*api_keys.ApiKey, error)
	
	// GetAPIKeyByID 根据ID获取API密钥
	GetAPIKeyByID(ctx context.Context, id int64) (*api_keys.ApiKey, error)
	
	// ListAPIKeysByUser 获取指定用户的所有API密钥
	ListAPIKeysByUser(ctx context.Context, userID int64) ([]api_keys.ApiKey, error)
	
	// ListAPIKeysByProvider 获取指定提供商的所有API密钥
	ListAPIKeysByProvider(ctx context.Context, providerType string) ([]api_keys.ApiKey, error)
	
	// UpdateAPIKey 更新API密钥
	UpdateAPIKey(ctx context.Context, params UpdateAPIKeyParams) (*api_keys.ApiKey, error)
	
	// DeactivateAPIKey 停用API密钥
	DeactivateAPIKey(ctx context.Context, userID int64, providerType string) error
	
	// DeleteAPIKey 删除API密钥
	DeleteAPIKey(ctx context.Context, userID int64, providerType string) error
	
	// CheckAPIKeyExists 检查API密钥是否存在
	CheckAPIKeyExists(ctx context.Context, userID int64, providerType string) (bool, error)
	
	// CountAPIKeysByUser 统计用户的API密钥数量
	CountAPIKeysByUser(ctx context.Context, userID int64) (int64, error)
	
	// CountAPIKeysByProvider 统计提供商的API密钥数量
	CountAPIKeysByProvider(ctx context.Context, providerType string) (int64, error)
}

// CreateAPIKeyParams 创建API密钥参数
type CreateAPIKeyParams struct {
	UserID       int64  `json:"user_id"`
	ProviderType string `json:"provider_type"`
	EncryptedKey string `json:"encrypted_key"`
	KeyHash      string `json:"key_hash"`
	IsActive     bool   `json:"is_active"`
}

// UpdateAPIKeyParams 更新API密钥参数
type UpdateAPIKeyParams struct {
	UserID       int64  `json:"user_id"`
	ProviderType string `json:"provider_type"`
	EncryptedKey string `json:"encrypted_key"`
	KeyHash      string `json:"key_hash"`
}