package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"sync"

	"go-springAi/internal/repository"
)

// DatabaseKeyManager 基于数据库的密钥管理器
type DatabaseKeyManager struct {
	mu           sync.RWMutex
	userID       int64
	providerType string
	encryptKey   []byte
	repo         repository.APIKeyRepository
}

// NewDatabaseKeyManager 创建新的数据库密钥管理器
func NewDatabaseKeyManager(userID int64, providerType string, repo repository.APIKeyRepository) *DatabaseKeyManager {
	// 使用固定的加密密钥（实际应用中应该从配置中获取）
	// 这里使用SHA256哈希生成固定的32字节密钥
	fixedSeed := "go-springAi-encryption-key-v1.0"
	hash := sha256.Sum256([]byte(fixedSeed))
	encryptKey := hash[:]
	
	return &DatabaseKeyManager{
		userID:       userID,
		providerType: providerType,
		encryptKey:   encryptKey,
		repo:         repo,
	}
}

// SetAPIKey 设置 API 密钥
func (km *DatabaseKeyManager) SetAPIKey(key string) error {
	if key == "" {
		return fmt.Errorf("API key cannot be empty")
	}
	
	km.mu.Lock()
	defer km.mu.Unlock()
	
	ctx := context.Background()
	
	// 加密密钥
	encryptedKey, err := km.EncryptKey(key)
	if err != nil {
		return fmt.Errorf("failed to encrypt key: %w", err)
	}
	
	// 生成密钥哈希
	keyHash := km.generateKeyHash(key)
	
	// 检查是否已存在
	exists, err := km.repo.CheckAPIKeyExists(ctx, km.userID, km.providerType)
	if err != nil {
		return fmt.Errorf("failed to check key existence: %w", err)
	}
	
	if exists {
		// 更新现有密钥
		_, err = km.repo.UpdateAPIKey(ctx, repository.UpdateAPIKeyParams{
			UserID:       km.userID,
			ProviderType: km.providerType,
			EncryptedKey: encryptedKey,
			KeyHash:      keyHash,
		})
		if err != nil {
			return fmt.Errorf("failed to update API key: %w", err)
		}
	} else {
		// 创建新密钥
		_, err = km.repo.CreateAPIKey(ctx, repository.CreateAPIKeyParams{
			UserID:       km.userID,
			ProviderType: km.providerType,
			EncryptedKey: encryptedKey,
			KeyHash:      keyHash,
			IsActive:     true,
		})
		if err != nil {
			return fmt.Errorf("failed to create API key: %w", err)
		}
	}
	
	return nil
}

// GetAPIKey 获取 API 密钥
func (km *DatabaseKeyManager) GetAPIKey() (string, error) {
	km.mu.RLock()
	defer km.mu.RUnlock()
	
	ctx := context.Background()
	
	apiKey, err := km.repo.GetAPIKey(ctx, km.userID, km.providerType)
	if err != nil {
		return "", fmt.Errorf("failed to get API key: %w", err)
	}
	
	if !apiKey.IsActive.Bool {
		return "", fmt.Errorf("API key is inactive")
	}
	
	// 解密密钥
	decryptedKey, err := km.DecryptKey(apiKey.EncryptedKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt key: %w", err)
	}
	
	return decryptedKey, nil
}

// ValidateKey 验证 API 密钥格式
func (km *DatabaseKeyManager) ValidateKey(key string) error {
	if key == "" {
		return fmt.Errorf("API key is empty")
	}
	
	// 根据提供商类型进行不同的验证
	switch km.providerType {
	case "openai":
		if len(key) < 10 {
			return fmt.Errorf("OpenAI API key is too short")
		}
		// OpenAI API 密钥通常以 "sk-" 开头
		if len(key) > 4 && key[:3] != "sk-" {
			return fmt.Errorf("OpenAI API key should start with 'sk-'")
		}
	case "googleai":
		if len(key) < 10 {
			return fmt.Errorf("Google AI API key is too short")
		}
		// Google AI API 密钥通常以 "AIza" 开头
		if len(key) > 4 && key[:4] != "AIza" {
			return fmt.Errorf("Google AI API key should start with 'AIza'")
		}
	default:
		// 通用验证
		if len(key) < 10 {
			return fmt.Errorf("API key is too short")
		}
	}
	
	return nil
}

// EncryptKey 加密 API 密钥
func (km *DatabaseKeyManager) EncryptKey(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("key cannot be empty")
	}
	
	block, err := aes.NewCipher(km.encryptKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}
	
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}
	
	ciphertext := gcm.Seal(nonce, nonce, []byte(key), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptKey 解密 API 密钥
func (km *DatabaseKeyManager) DecryptKey(encryptedKey string) (string, error) {
	if encryptedKey == "" {
		return "", fmt.Errorf("encrypted key cannot be empty")
	}
	
	data, err := base64.StdEncoding.DecodeString(encryptedKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}
	
	block, err := aes.NewCipher(km.encryptKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}
	
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}
	
	return string(plaintext), nil
}

// generateKeyHash 生成密钥哈希
func (km *DatabaseKeyManager) generateKeyHash(key string) string {
	hash := sha256.Sum256([]byte(key))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// IsActive 检查密钥是否激活
func (km *DatabaseKeyManager) IsActive() (bool, error) {
	km.mu.RLock()
	defer km.mu.RUnlock()
	
	ctx := context.Background()
	
	apiKey, err := km.repo.GetAPIKey(ctx, km.userID, km.providerType)
	if err != nil {
		return false, err
	}
	
	return apiKey.IsActive.Bool, nil
}

// Deactivate 停用密钥
func (km *DatabaseKeyManager) Deactivate() error {
	km.mu.Lock()
	defer km.mu.Unlock()
	
	ctx := context.Background()
	
	return km.repo.DeactivateAPIKey(ctx, km.userID, km.providerType)
}

// Delete 删除密钥
func (km *DatabaseKeyManager) Delete() error {
	km.mu.Lock()
	defer km.mu.Unlock()
	
	ctx := context.Background()
	
	return km.repo.DeleteAPIKey(ctx, km.userID, km.providerType)
}