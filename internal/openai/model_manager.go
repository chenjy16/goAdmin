package openai

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// FileModelManager 基于文件的模型管理器
type FileModelManager struct {
	configFile string
	models     map[string]*ModelConfig
	mu         sync.RWMutex
}

// NewFileModelManager 创建新的文件模型管理器
func NewFileModelManager(configFile string) *FileModelManager {
	mm := &FileModelManager{
		configFile: configFile,
		models:     make(map[string]*ModelConfig),
	}
	
	// 加载配置
	if err := mm.loadConfig(); err != nil {
		// 如果加载失败，使用默认配置
		mm.models = DefaultModels()
		mm.saveConfig() // 保存默认配置
	}
	
	return mm
}

// GetModel 获取模型配置
func (mm *FileModelManager) GetModel(name string) (*ModelConfig, error) {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	
	model, exists := mm.models[name]
	if !exists {
		return nil, fmt.Errorf("model %s not found", name)
	}
	
	// 返回副本以避免并发修改
	modelCopy := *model
	return &modelCopy, nil
}

// ListModels 列出所有模型
func (mm *FileModelManager) ListModels() map[string]*ModelConfig {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	
	// 返回副本以避免并发修改
	result := make(map[string]*ModelConfig)
	for name, model := range mm.models {
		modelCopy := *model
		result[name] = &modelCopy
	}
	
	return result
}

// UpdateModel 更新模型配置
func (mm *FileModelManager) UpdateModel(name string, config *ModelConfig) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	// 验证配置
	if err := mm.validateModelConfig(config); err != nil {
		return fmt.Errorf("invalid model config: %w", err)
	}
	
	// 更新配置
	mm.models[name] = config
	
	// 保存到文件
	return mm.saveConfig()
}

// EnableModel 启用模型
func (mm *FileModelManager) EnableModel(name string) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	model, exists := mm.models[name]
	if !exists {
		return fmt.Errorf("model %s not found", name)
	}
	
	model.Enabled = true
	return mm.saveConfig()
}

// DisableModel 禁用模型
func (mm *FileModelManager) DisableModel(name string) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	model, exists := mm.models[name]
	if !exists {
		return fmt.Errorf("model %s not found", name)
	}
	
	model.Enabled = false
	return mm.saveConfig()
}

// AddModel 添加新模型
func (mm *FileModelManager) AddModel(name string, config *ModelConfig) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	// 验证配置
	if err := mm.validateModelConfig(config); err != nil {
		return fmt.Errorf("invalid model config: %w", err)
	}
	
	// 设置名称
	config.Name = name
	
	// 添加模型
	mm.models[name] = config
	
	// 保存到文件
	return mm.saveConfig()
}

// RemoveModel 移除模型
func (mm *FileModelManager) RemoveModel(name string) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	if _, exists := mm.models[name]; !exists {
		return fmt.Errorf("model %s not found", name)
	}
	
	delete(mm.models, name)
	return mm.saveConfig()
}

// GetEnabledModels 获取启用的模型
func (mm *FileModelManager) GetEnabledModels() map[string]*ModelConfig {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	
	result := make(map[string]*ModelConfig)
	for name, model := range mm.models {
		if model.Enabled {
			modelCopy := *model
			result[name] = &modelCopy
		}
	}
	
	return result
}

// validateModelConfig 验证模型配置
func (mm *FileModelManager) validateModelConfig(config *ModelConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}
	
	if config.Name == "" {
		return fmt.Errorf("model name cannot be empty")
	}
	
	if config.MaxTokens <= 0 {
		return fmt.Errorf("max tokens must be positive")
	}
	
	if config.Temperature < 0 || config.Temperature > 2 {
		return fmt.Errorf("temperature must be between 0 and 2")
	}
	
	if config.TopP < 0 || config.TopP > 1 {
		return fmt.Errorf("top_p must be between 0 and 1")
	}
	
	if config.FrequencyPenalty < -2 || config.FrequencyPenalty > 2 {
		return fmt.Errorf("frequency_penalty must be between -2 and 2")
	}
	
	if config.PresencePenalty < -2 || config.PresencePenalty > 2 {
		return fmt.Errorf("presence_penalty must be between -2 and 2")
	}
	
	return nil
}

// loadConfig 从文件加载配置
func (mm *FileModelManager) loadConfig() error {
	// 检查文件是否存在
	if _, err := os.Stat(mm.configFile); os.IsNotExist(err) {
		return fmt.Errorf("config file not found")
	}
	
	// 读取文件
	data, err := os.ReadFile(mm.configFile)
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}
	
	// 解析 JSON
	var models map[string]*ModelConfig
	if err := json.Unmarshal(data, &models); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}
	
	mm.models = models
	return nil
}

// saveConfig 保存配置到文件
func (mm *FileModelManager) saveConfig() error {
	// 确保目录存在
	dir := filepath.Dir(mm.configFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}
	
	// 序列化配置
	data, err := json.MarshalIndent(mm.models, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	
	// 写入文件
	if err := os.WriteFile(mm.configFile, data, 0644); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}
	
	return nil
}

// MemoryModelManager 基于内存的模型管理器（用于测试）
type MemoryModelManager struct {
	models map[string]*ModelConfig
	mu     sync.RWMutex
}

// NewMemoryModelManager 创建新的内存模型管理器
func NewMemoryModelManager() *MemoryModelManager {
	return &MemoryModelManager{
		models: DefaultModels(),
	}
}

// GetModel 获取模型配置
func (mm *MemoryModelManager) GetModel(name string) (*ModelConfig, error) {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	
	model, exists := mm.models[name]
	if !exists {
		return nil, fmt.Errorf("model %s not found", name)
	}
	
	// 返回副本
	modelCopy := *model
	return &modelCopy, nil
}

// ListModels 列出所有模型
func (mm *MemoryModelManager) ListModels() map[string]*ModelConfig {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	
	// 返回副本
	result := make(map[string]*ModelConfig)
	for name, model := range mm.models {
		modelCopy := *model
		result[name] = &modelCopy
	}
	
	return result
}

// UpdateModel 更新模型配置
func (mm *MemoryModelManager) UpdateModel(name string, config *ModelConfig) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	mm.models[name] = config
	return nil
}

// EnableModel 启用模型
func (mm *MemoryModelManager) EnableModel(name string) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	model, exists := mm.models[name]
	if !exists {
		return fmt.Errorf("model %s not found", name)
	}
	
	model.Enabled = true
	return nil
}

// DisableModel 禁用模型
func (mm *MemoryModelManager) DisableModel(name string) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	model, exists := mm.models[name]
	if !exists {
		return fmt.Errorf("model %s not found", name)
	}
	
	model.Enabled = false
	return nil
}