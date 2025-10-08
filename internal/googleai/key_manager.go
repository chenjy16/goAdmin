package googleai

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"sync"
)

// keyManager Google AI 密钥管理器
type keyManager struct {
	mu         sync.RWMutex
	apiKey     string
	encryptKey []byte // 用于加密的密钥
}

// NewKeyManager 创建新的密钥管理器
func NewKeyManager(apiKey string) KeyManager {
	// 生成一个简单的加密密钥（实际应用中应该从配置中获取）
	encryptKey := make([]byte, 32)
	rand.Read(encryptKey)
	
	return &keyManager{
		apiKey:     apiKey,
		encryptKey: encryptKey,
	}
}

// GetAPIKey 获取 API 密钥
func (km *keyManager) GetAPIKey() (string, error) {
	km.mu.RLock()
	defer km.mu.RUnlock()
	
	if km.apiKey == "" {
		return "", fmt.Errorf("API key is not set")
	}
	
	return km.apiKey, nil
}

// SetAPIKey 设置 API 密钥
func (km *keyManager) SetAPIKey(key string) error {
	if key == "" {
		return fmt.Errorf("API key cannot be empty")
	}
	
	km.mu.Lock()
	defer km.mu.Unlock()
	km.apiKey = key
	return nil
}

// ValidateKey 验证 API 密钥格式
func (km *keyManager) ValidateKey(key string) error {
	if key == "" {
		return fmt.Errorf("API key is empty")
	}
	
	// Google AI API 密钥通常以 "AIza" 开头
	if len(key) < 10 {
		return fmt.Errorf("API key is too short")
	}
	
	return nil
}

// EncryptKey 加密 API 密钥
func (km *keyManager) EncryptKey(key string) (string, error) {
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
func (km *keyManager) DecryptKey(encryptedKey string) (string, error) {
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