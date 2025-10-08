package openai

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// FileKeyManager 基于文件的密钥管理器
type FileKeyManager struct {
	keyFile    string
	passphrase string
	mu         sync.RWMutex
}

// NewFileKeyManager 创建新的文件密钥管理器
func NewFileKeyManager(keyFile, passphrase string) *FileKeyManager {
	return &FileKeyManager{
		keyFile:    keyFile,
		passphrase: passphrase,
	}
}

// SetAPIKey 设置 API 密钥
func (km *FileKeyManager) SetAPIKey(key string) error {
	km.mu.Lock()
	defer km.mu.Unlock()
	
	// 验证密钥格式
	if err := km.ValidateKey(key); err != nil {
		return fmt.Errorf("invalid API key: %w", err)
	}
	
	// 加密密钥
	encryptedKey, err := km.EncryptKey(key)
	if err != nil {
		return fmt.Errorf("encrypt key: %w", err)
	}
	
	// 确保目录存在
	dir := filepath.Dir(km.keyFile)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}
	
	// 写入文件
	if err := os.WriteFile(km.keyFile, []byte(encryptedKey), 0600); err != nil {
		return fmt.Errorf("write key file: %w", err)
	}
	
	return nil
}

// GetAPIKey 获取 API 密钥
func (km *FileKeyManager) GetAPIKey() (string, error) {
	km.mu.RLock()
	defer km.mu.RUnlock()
	
	// 检查文件是否存在
	if _, err := os.Stat(km.keyFile); os.IsNotExist(err) {
		return "", fmt.Errorf("API key not found")
	}
	
	// 读取加密的密钥
	encryptedData, err := os.ReadFile(km.keyFile)
	if err != nil {
		return "", fmt.Errorf("read key file: %w", err)
	}
	
	// 解密密钥
	key, err := km.DecryptKey(string(encryptedData))
	if err != nil {
		return "", fmt.Errorf("decrypt key: %w", err)
	}
	
	return key, nil
}

// ValidateKey 验证密钥格式
func (km *FileKeyManager) ValidateKey(key string) error {
	if key == "" {
		return fmt.Errorf("API key cannot be empty")
	}
	
	// OpenAI API 密钥格式验证
	// sk-开头，后跟字母数字字符
	pattern := `^sk-[a-zA-Z0-9]{48}$`
	matched, err := regexp.MatchString(pattern, key)
	if err != nil {
		return fmt.Errorf("regex error: %w", err)
	}
	
	if !matched {
		return fmt.Errorf("invalid OpenAI API key format")
	}
	
	return nil
}

// EncryptKey 加密密钥
func (km *FileKeyManager) EncryptKey(key string) (string, error) {
	// 生成密钥
	hash := sha256.Sum256([]byte(km.passphrase))
	
	// 创建 AES 加密器
	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return "", fmt.Errorf("create cipher: %w", err)
	}
	
	// 创建 GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create GCM: %w", err)
	}
	
	// 生成随机 nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("generate nonce: %w", err)
	}
	
	// 加密
	ciphertext := gcm.Seal(nonce, nonce, []byte(key), nil)
	
	// 编码为 base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptKey 解密密钥
func (km *FileKeyManager) DecryptKey(encryptedKey string) (string, error) {
	// 解码 base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedKey)
	if err != nil {
		return "", fmt.Errorf("decode base64: %w", err)
	}
	
	// 生成密钥
	hash := sha256.Sum256([]byte(km.passphrase))
	
	// 创建 AES 解密器
	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return "", fmt.Errorf("create cipher: %w", err)
	}
	
	// 创建 GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create GCM: %w", err)
	}
	
	// 检查密文长度
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	
	// 提取 nonce 和密文
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	
	// 解密
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt: %w", err)
	}
	
	return string(plaintext), nil
}

// MemoryKeyManager 基于内存的密钥管理器（用于测试）
type MemoryKeyManager struct {
	key string
	mu  sync.RWMutex
}

// NewMemoryKeyManager 创建新的内存密钥管理器
func NewMemoryKeyManager() *MemoryKeyManager {
	return &MemoryKeyManager{}
}

// SetAPIKey 设置 API 密钥
func (km *MemoryKeyManager) SetAPIKey(key string) error {
	km.mu.Lock()
	defer km.mu.Unlock()
	
	// 验证密钥格式
	if err := km.ValidateKey(key); err != nil {
		return fmt.Errorf("invalid API key: %w", err)
	}
	
	km.key = key
	return nil
}

// GetAPIKey 获取 API 密钥
func (km *MemoryKeyManager) GetAPIKey() (string, error) {
	km.mu.RLock()
	defer km.mu.RUnlock()
	
	if km.key == "" {
		return "", fmt.Errorf("API key not found")
	}
	
	return km.key, nil
}

// ValidateKey 验证密钥格式
func (km *MemoryKeyManager) ValidateKey(key string) error {
	if key == "" {
		return fmt.Errorf("API key cannot be empty")
	}
	
	// 简单的格式验证
	if !strings.HasPrefix(key, "sk-") {
		return fmt.Errorf("invalid OpenAI API key format")
	}
	
	return nil
}

// EncryptKey 加密密钥（内存管理器不需要加密）
func (km *MemoryKeyManager) EncryptKey(key string) (string, error) {
	return key, nil
}

// DecryptKey 解密密钥（内存管理器不需要解密）
func (km *MemoryKeyManager) DecryptKey(encryptedKey string) (string, error) {
	return encryptedKey, nil
}