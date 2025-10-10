package repository

import (
	"context"
	"database/sql"
	"fmt"

	"go-springAi/internal/database"
	"go-springAi/internal/database/generated/api_keys"
	"go-springAi/internal/errors"
)

// apiKeyRepository API密钥数据访问层实现
type apiKeyRepository struct {
	db *database.DB
}

// NewAPIKeyRepository 创建API密钥数据访问层
func NewAPIKeyRepository(db *database.DB) APIKeyRepository {
	return &apiKeyRepository{
		db: db,
	}
}

// CreateAPIKey 创建API密钥
func (r *apiKeyRepository) CreateAPIKey(ctx context.Context, params CreateAPIKeyParams) (*api_keys.ApiKey, error) {
	apiKey, err := r.db.APIKeys.CreateAPIKey(ctx, api_keys.CreateAPIKeyParams{
		UserID:       params.UserID,
		ProviderType: params.ProviderType,
		EncryptedKey: params.EncryptedKey,
		KeyHash:      params.KeyHash,
		IsActive:     sql.NullBool{Bool: params.IsActive, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create API key: %w", err)
	}
	return &apiKey, nil
}

// GetAPIKey 获取指定用户和提供商的API密钥
func (r *apiKeyRepository) GetAPIKey(ctx context.Context, userID int64, providerType string) (*api_keys.ApiKey, error) {
	apiKey, err := r.db.APIKeys.GetAPIKey(ctx, api_keys.GetAPIKeyParams{
		UserID:       userID,
		ProviderType: providerType,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("API key")
		}
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}
	return &apiKey, nil
}

// GetAPIKeyByID 根据ID获取API密钥
func (r *apiKeyRepository) GetAPIKeyByID(ctx context.Context, id int64) (*api_keys.ApiKey, error) {
	apiKey, err := r.db.APIKeys.GetAPIKeyByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("API key")
		}
		return nil, fmt.Errorf("failed to get API key by ID: %w", err)
	}
	return &apiKey, nil
}

// ListAPIKeysByUser 获取指定用户的所有API密钥
func (r *apiKeyRepository) ListAPIKeysByUser(ctx context.Context, userID int64) ([]api_keys.ApiKey, error) {
	apiKeys, err := r.db.APIKeys.ListAPIKeysByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list API keys by user: %w", err)
	}
	return apiKeys, nil
}

// ListAPIKeysByProvider 获取指定提供商的所有API密钥
func (r *apiKeyRepository) ListAPIKeysByProvider(ctx context.Context, providerType string) ([]api_keys.ApiKey, error) {
	apiKeys, err := r.db.APIKeys.ListAPIKeysByProvider(ctx, providerType)
	if err != nil {
		return nil, fmt.Errorf("failed to list API keys by provider: %w", err)
	}
	return apiKeys, nil
}

// UpdateAPIKey 更新API密钥
func (r *apiKeyRepository) UpdateAPIKey(ctx context.Context, params UpdateAPIKeyParams) (*api_keys.ApiKey, error) {
	apiKey, err := r.db.APIKeys.UpdateAPIKey(ctx, api_keys.UpdateAPIKeyParams{
		UserID:       params.UserID,
		ProviderType: params.ProviderType,
		EncryptedKey: params.EncryptedKey,
		KeyHash:      params.KeyHash,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("API key")
		}
		return nil, fmt.Errorf("failed to update API key: %w", err)
	}
	return &apiKey, nil
}

// DeactivateAPIKey 停用API密钥
func (r *apiKeyRepository) DeactivateAPIKey(ctx context.Context, userID int64, providerType string) error {
	err := r.db.APIKeys.DeactivateAPIKey(ctx, api_keys.DeactivateAPIKeyParams{
		UserID:       userID,
		ProviderType: providerType,
	})
	if err != nil {
		return fmt.Errorf("failed to deactivate API key: %w", err)
	}
	return nil
}

// DeleteAPIKey 删除API密钥
func (r *apiKeyRepository) DeleteAPIKey(ctx context.Context, userID int64, providerType string) error {
	err := r.db.APIKeys.DeleteAPIKey(ctx, api_keys.DeleteAPIKeyParams{
		UserID:       userID,
		ProviderType: providerType,
	})
	if err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}
	return nil
}

// CheckAPIKeyExists 检查API密钥是否存在
func (r *apiKeyRepository) CheckAPIKeyExists(ctx context.Context, userID int64, providerType string) (bool, error) {
	count, err := r.db.APIKeys.CheckAPIKeyExists(ctx, api_keys.CheckAPIKeyExistsParams{
		UserID:       userID,
		ProviderType: providerType,
	})
	if err != nil {
		return false, fmt.Errorf("failed to check API key exists: %w", err)
	}
	return count > 0, nil
}

// CountAPIKeysByUser 统计用户的API密钥数量
func (r *apiKeyRepository) CountAPIKeysByUser(ctx context.Context, userID int64) (int64, error) {
	count, err := r.db.APIKeys.CountAPIKeysByUser(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to count API keys by user: %w", err)
	}
	return count, nil
}

// CountAPIKeysByProvider 统计提供商的API密钥数量
func (r *apiKeyRepository) CountAPIKeysByProvider(ctx context.Context, providerType string) (int64, error) {
	count, err := r.db.APIKeys.CountAPIKeysByProvider(ctx, providerType)
	if err != nil {
		return 0, fmt.Errorf("failed to count API keys by provider: %w", err)
	}
	return count, nil
}