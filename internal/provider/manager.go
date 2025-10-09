package provider

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"go-springAi/internal/logger"
)

// Manager Provider管理器
type Manager struct {
	providers map[ProviderType]Provider
	mu        sync.RWMutex
	logger    logger.Logger
}

// NewManager 创建新的Provider管理器
func NewManager(logger logger.Logger) *Manager {
	return &Manager{
		providers: make(map[ProviderType]Provider),
		logger:    logger,
	}
}

// RegisterProvider 注册Provider
func (m *Manager) RegisterProvider(provider Provider) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	providerType := provider.GetType()
	if _, exists := m.providers[providerType]; exists {
		return fmt.Errorf("provider %s already registered", providerType)
	}
	
	m.providers[providerType] = provider
	m.logger.Info("Provider registered",
		logger.String("type", string(providerType)),
		logger.String("name", provider.GetName()),
	)
	
	return nil
}

// GetProvider 获取指定类型的Provider
func (m *Manager) GetProvider(providerType ProviderType) (Provider, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	provider, exists := m.providers[providerType]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", providerType)
	}
	
	return provider, nil
}

// GetProviderByName 根据名称获取Provider
func (m *Manager) GetProviderByName(name string) (Provider, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	for _, provider := range m.providers {
		if provider.GetName() == name {
			return provider, nil
		}
	}
	
	return nil, fmt.Errorf("provider with name %s not found", name)
}

// ListProviders 列出所有注册的Provider
func (m *Manager) ListProviders() []ProviderInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	var providers []ProviderInfo
	for _, provider := range m.providers {
		// 获取模型数量
		models, err := provider.ListModels(context.Background())
		modelCount := 0
		if err == nil {
			modelCount = len(models)
		}
		
		// 检查健康状态
		healthy := provider.IsHealthy(context.Background())
		
		providers = append(providers, ProviderInfo{
			Type:        provider.GetType(),
			Name:        provider.GetName(),
			Description: fmt.Sprintf("%s AI Provider", provider.GetName()),
			Healthy:     healthy,
			ModelCount:  modelCount,
		})
	}
	
	return providers
}

// GetAvailableProviders 获取可用的Provider列表
func (m *Manager) GetAvailableProviders(ctx context.Context) []ProviderInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	var availableProviders []ProviderInfo
	for _, provider := range m.providers {
		if provider.IsHealthy(ctx) {
			// 获取模型数量
			models, err := provider.ListModels(ctx)
			modelCount := 0
			if err == nil {
				modelCount = len(models)
			}
			
			availableProviders = append(availableProviders, ProviderInfo{
				Type:        provider.GetType(),
				Name:        provider.GetName(),
				Description: fmt.Sprintf("%s AI Provider", provider.GetName()),
				Healthy:     true,
				ModelCount:  modelCount,
			})
		}
	}
	
	return availableProviders
}

// GetProviderTypes 获取所有已注册的Provider类型
func (m *Manager) GetProviderTypes() []ProviderType {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	var types []ProviderType
	for providerType := range m.providers {
		types = append(types, providerType)
	}
	
	return types
}

// IsProviderRegistered 检查Provider是否已注册
func (m *Manager) IsProviderRegistered(providerType ProviderType) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	_, exists := m.providers[providerType]
	return exists
}

// ValidateModelForProvider 验证模型是否存在于指定的提供商中
func (m *Manager) ValidateModelForProvider(ctx context.Context, providerName, modelName string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// 通过名称获取提供商
	provider, err := m.getProviderByNameUnsafe(providerName)
	if err != nil {
		return fmt.Errorf("provider %s not found: %w", providerName, err)
	}
	
	// 获取提供商的模型列表
	models, err := provider.ListModels(ctx)
	if err != nil {
		return fmt.Errorf("failed to get models for provider %s: %w", providerName, err)
	}
	
	// 检查模型是否存在
	if _, exists := models[modelName]; !exists {
		return fmt.Errorf("model %s not found in provider %s", modelName, providerName)
	}
	
	return nil
}

// GetProviderByModelWithValidation 根据模型名称获取对应的Provider，并验证模型存在
func (m *Manager) GetProviderByModelWithValidation(ctx context.Context, modelName string) (Provider, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// 遍历所有提供商，查找包含该模型的提供商
	for _, provider := range m.providers {
		models, err := provider.ListModels(ctx)
		if err != nil {
			m.logger.Warn("Failed to get models for provider",
				logger.String("provider", provider.GetName()),
				logger.ZapError(err))
			continue
		}
		
		if _, exists := models[modelName]; exists {
			return provider, nil
		}
	}
	
	return nil, fmt.Errorf("model %s not found in any registered provider", modelName)
}

// GetProviderByModel 根据模型名称获取对应的Provider（保持向后兼容）
func (m *Manager) GetProviderByModel(modelName string) (Provider, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// 根据模型名称前缀映射到提供商
	var providerType ProviderType
	
	switch {
	case strings.HasPrefix(modelName, "gpt-"):
		providerType = ProviderTypeOpenAI
	case strings.HasPrefix(modelName, "gemini-"):
		providerType = ProviderTypeGoogleAI
	case strings.HasPrefix(modelName, "claude-"):
		// 为未来的Claude支持预留
		return nil, fmt.Errorf("claude provider not implemented yet")
	default:
		// 默认使用OpenAI
		providerType = ProviderTypeOpenAI
	}
	
	provider, exists := m.providers[providerType]
	if !exists {
		return nil, fmt.Errorf("provider %s not found for model %s", providerType, modelName)
	}
	
	return provider, nil
}

// getProviderByNameUnsafe 内部方法，不加锁获取提供商（调用者需要持有锁）
func (m *Manager) getProviderByNameUnsafe(name string) (Provider, error) {
	for _, provider := range m.providers {
		if provider.GetName() == name {
			return provider, nil
		}
	}
	return nil, fmt.Errorf("provider %s not found", name)
}

// UnregisterProvider 注销Provider
func (m *Manager) UnregisterProvider(providerType ProviderType) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	provider, exists := m.providers[providerType]
	if !exists {
		return fmt.Errorf("provider %s not registered", providerType)
	}
	
	delete(m.providers, providerType)
	m.logger.Info("Provider unregistered",
		logger.String("type", string(providerType)),
		logger.String("name", provider.GetName()),
	)
	
	return nil
}

// GetHealthStatus 获取所有Provider的健康状态
func (m *Manager) GetHealthStatus(ctx context.Context) map[ProviderType]bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	status := make(map[ProviderType]bool)
	for providerType, provider := range m.providers {
		status[providerType] = provider.IsHealthy(ctx)
	}
	
	return status
}