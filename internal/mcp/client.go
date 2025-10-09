package mcp

import (
	"context"
	"fmt"

	"admin/internal/dto"
)

// MCPServiceInterface MCP服务接口（避免循环导入）
type MCPServiceInterface interface {
	Initialize(ctx context.Context, req *dto.MCPInitializeRequest) (*dto.MCPInitializeResponse, error)
	ListTools(ctx context.Context) (*dto.MCPToolsResponse, error)
	ExecuteTool(ctx context.Context, req *dto.MCPExecuteRequest) (*dto.MCPExecuteResponse, error)
	GetExecutionLog(ctx context.Context, executionID string) (*dto.MCPToolExecutionLog, error)
	ListExecutionLogs(ctx context.Context, userID *string, limit int) ([]*dto.MCPToolExecutionLog, error)
}

// InternalMCPClient 内部MCP客户端接口
type InternalMCPClient interface {
	// Initialize 初始化MCP连接
	Initialize(ctx context.Context, req *dto.MCPInitializeRequest) (*dto.MCPInitializeResponse, error)
	// ListTools 获取可用工具列表
	ListTools(ctx context.Context) (*dto.MCPToolsResponse, error)
	// ExecuteTool 执行工具
	ExecuteTool(ctx context.Context, req *dto.MCPExecuteRequest) (*dto.MCPExecuteResponse, error)
	// GetExecutionLog 获取执行日志
	GetExecutionLog(ctx context.Context, executionID string) (*dto.MCPToolExecutionLog, error)
	// ListExecutionLogs 列出执行日志
	ListExecutionLogs(ctx context.Context, userID *string, limit int) ([]*dto.MCPToolExecutionLog, error)
}

// InternalMCPClientImpl 内部MCP客户端实现
type InternalMCPClientImpl struct {
	mcpService  MCPServiceInterface
	clientInfo  dto.MCPClientInfo
	initialized bool
}

// NewInternalMCPClient 创建内部MCP客户端
func NewInternalMCPClient(mcpService MCPServiceInterface, clientInfo dto.MCPClientInfo) InternalMCPClient {
	return &InternalMCPClientImpl{
		mcpService:  mcpService,
		clientInfo:  clientInfo,
		initialized: false,
	}
}

// Initialize 初始化MCP连接
func (c *InternalMCPClientImpl) Initialize(ctx context.Context, req *dto.MCPInitializeRequest) (*dto.MCPInitializeResponse, error) {
	if req == nil {
		req = &dto.MCPInitializeRequest{
			ProtocolVersion: "2024-11-05",
			Capabilities: dto.MCPCapabilities{
				Tools: &dto.MCPToolsCapability{
					ListChanged: true,
				},
			},
			ClientInfo: c.clientInfo,
		}
	}

	resp, err := c.mcpService.Initialize(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MCP client: %w", err)
	}

	c.initialized = true
	return resp, nil
}

// ListTools 获取可用工具列表
func (c *InternalMCPClientImpl) ListTools(ctx context.Context) (*dto.MCPToolsResponse, error) {
	if !c.initialized {
		return nil, fmt.Errorf("MCP client not initialized")
	}

	return c.mcpService.ListTools(ctx)
}

// ExecuteTool 执行工具
func (c *InternalMCPClientImpl) ExecuteTool(ctx context.Context, req *dto.MCPExecuteRequest) (*dto.MCPExecuteResponse, error) {
	if !c.initialized {
		return nil, fmt.Errorf("MCP client not initialized")
	}

	if req == nil {
		return nil, fmt.Errorf("execute request cannot be nil")
	}

	return c.mcpService.ExecuteTool(ctx, req)
}

// GetExecutionLog 获取执行日志
func (c *InternalMCPClientImpl) GetExecutionLog(ctx context.Context, executionID string) (*dto.MCPToolExecutionLog, error) {
	if !c.initialized {
		return nil, fmt.Errorf("MCP client not initialized")
	}

	return c.mcpService.GetExecutionLog(ctx, executionID)
}

// ListExecutionLogs 列出执行日志
func (c *InternalMCPClientImpl) ListExecutionLogs(ctx context.Context, userID *string, limit int) ([]*dto.MCPToolExecutionLog, error) {
	if !c.initialized {
		return nil, fmt.Errorf("MCP client not initialized")
	}

	return c.mcpService.ListExecutionLogs(ctx, userID, limit)
}

// MCPClientManager MCP客户端管理器
type MCPClientManager struct {
	clients map[string]InternalMCPClient
}

// NewMCPClientManager 创建MCP客户端管理器
func NewMCPClientManager() *MCPClientManager {
	return &MCPClientManager{
		clients: make(map[string]InternalMCPClient),
	}
}

// RegisterClient 注册客户端
func (m *MCPClientManager) RegisterClient(name string, client InternalMCPClient) {
	m.clients[name] = client
}

// GetClient 获取客户端
func (m *MCPClientManager) GetClient(name string) (InternalMCPClient, bool) {
	client, exists := m.clients[name]
	return client, exists
}

// ListClients 列出所有客户端
func (m *MCPClientManager) ListClients() []string {
	names := make([]string, 0, len(m.clients))
	for name := range m.clients {
		names = append(names, name)
	}
	return names
}