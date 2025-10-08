package provider

import (
	"context"
	"fmt"
	"sync"

	"admin/internal/logger"
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