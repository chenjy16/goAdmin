package googleai

import (
	"fmt"
	"sync"
)

// modelManager Google AI 模型管理器
type modelManager struct {
	mu     sync.RWMutex
	models map[string]*ModelConfig
}

// NewModelManager 创建新的模型管理器
func NewModelManager() ModelManager {
	mm := &modelManager{
		models: make(map[string]*ModelConfig),
	}
	
	// 初始化默认模型
	mm.initDefaultModels()
	
	return mm
}

// initDefaultModels 初始化默认模型配置
func (mm *modelManager) initDefaultModels() {
	defaultModels := []*ModelConfig{
		{
			Name:        "gemini-1.5-flash",
			DisplayName: "Gemini 1.5 Flash",
			MaxTokens:   8192,
			Temperature: 0.7,
			TopP:        0.9,
			TopK:        40,
			Enabled:     true,
		},
		{
			Name:        "gemini-1.5-pro",
			DisplayName: "Gemini 1.5 Pro",
			MaxTokens:   8192,
			Temperature: 0.7,
			TopP:        0.9,
			TopK:        40,
			Enabled:     true,
		},
		{
			Name:        "gemini-2.0-flash-exp",
			DisplayName: "Gemini 2.0 Flash (Experimental)",
			MaxTokens:   8192,
			Temperature: 0.7,
			TopP:        0.9,
			TopK:        40,
			Enabled:     true,
		},
		{
			Name:        "gemini-exp-1206",
			DisplayName: "Gemini Experimental 1206",
			MaxTokens:   8192,
			Temperature: 0.7,
			TopP:        0.9,
			TopK:        40,
			Enabled:     false, // 实验性模型默认禁用
		},
	}
	
	for _, model := range defaultModels {
		mm.models[model.Name] = model
	}
}

// GetModel 获取模型配置
func (mm *modelManager) GetModel(name string) (*ModelConfig, error) {
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
func (mm *modelManager) ListModels() map[string]*ModelConfig {
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
func (mm *modelManager) UpdateModel(name string, config *ModelConfig) error {
	if config == nil {
		return fmt.Errorf("model config cannot be nil")
	}
	
	if config.Name != name {
		return fmt.Errorf("model name mismatch: expected %s, got %s", name, config.Name)
	}
	
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	// 验证配置
	if err := mm.validateModelConfig(config); err != nil {
		return fmt.Errorf("invalid model config: %w", err)
	}
	
	// 创建副本并存储
	configCopy := *config
	mm.models[name] = &configCopy
	
	return nil
}

// EnableModel 启用模型
func (mm *modelManager) EnableModel(name string) error {
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
func (mm *modelManager) DisableModel(name string) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	model, exists := mm.models[name]
	if !exists {
		return fmt.Errorf("model %s not found", name)
	}
	
	model.Enabled = false
	return nil
}

// validateModelConfig 验证模型配置
func (mm *modelManager) validateModelConfig(config *ModelConfig) error {
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
	
	if config.TopK < 0 {
		return fmt.Errorf("top_k must be non-negative")
	}
	
	return nil
}