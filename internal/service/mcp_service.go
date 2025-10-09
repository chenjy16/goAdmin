package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-springAi/internal/dto"
	"go-springAi/internal/googleai"
	"go-springAi/internal/logger"
	"go-springAi/internal/mcp"
	"go-springAi/internal/openai"
	"go-springAi/internal/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// OpenAIServiceAdapter OpenAI服务适配器，将OpenAIService适配为mcp.OpenAIService
type OpenAIServiceAdapter struct {
	service *OpenAIService
}

// NewOpenAIServiceAdapter 创建OpenAI服务适配器
func NewOpenAIServiceAdapter(service *OpenAIService) *OpenAIServiceAdapter {
	return &OpenAIServiceAdapter{service: service}
}

// ChatCompletion 聊天完成
func (a *OpenAIServiceAdapter) ChatCompletion(ctx context.Context, req *mcp.ChatCompletionRequest) (*mcp.ChatCompletionResponse, error) {
	// 转换请求类型
	serviceReq := &ChatCompletionRequest{
		Model:       req.Model,
		Messages:    req.Messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Stream:      req.Stream,
		Options:     req.Options,
	}

	// 调用实际服务
	serviceResp, err := a.service.ChatCompletion(ctx, serviceReq)
	if err != nil {
		return nil, err
	}

	// 转换响应类型
	mcpResp := &mcp.ChatCompletionResponse{
		ID:      serviceResp.ID,
		Object:  serviceResp.Object,
		Created: serviceResp.Created,
		Model:   serviceResp.Model,
		Choices: serviceResp.Choices,
		Usage:   serviceResp.Usage,
	}

	return mcpResp, nil
}

// ListModels 列出模型
func (a *OpenAIServiceAdapter) ListModels(ctx context.Context) (map[string]*openai.ModelConfig, error) {
	return a.service.ListModels(ctx)
}

// ValidateAPIKey 验证API密钥
func (a *OpenAIServiceAdapter) ValidateAPIKey(ctx context.Context) error {
	return a.service.ValidateAPIKey(ctx)
}

// SetAPIKey 设置API密钥
func (a *OpenAIServiceAdapter) SetAPIKey(key string) error {
	return a.service.SetAPIKey(key)
}

// GetModelConfig 获取模型配置
func (a *OpenAIServiceAdapter) GetModelConfig(name string) (*openai.ModelConfig, error) {
	return a.service.GetModelConfig(name)
}

// EnableModel 启用模型
func (a *OpenAIServiceAdapter) EnableModel(name string) error {
	return a.service.EnableModel(name)
}

// DisableModel 禁用模型
func (a *OpenAIServiceAdapter) DisableModel(name string) error {
	return a.service.DisableModel(name)
}

// GoogleAIServiceAdapter Google AI服务适配器，将GoogleAIService适配为mcp.GoogleAIService
type GoogleAIServiceAdapter struct {
	service *GoogleAIService
}

// NewGoogleAIServiceAdapter 创建Google AI服务适配器
func NewGoogleAIServiceAdapter(service *GoogleAIService) *GoogleAIServiceAdapter {
	return &GoogleAIServiceAdapter{service: service}
}

// ChatCompletion 聊天完成
func (a *GoogleAIServiceAdapter) ChatCompletion(ctx context.Context, req *mcp.GoogleAIChatCompletionRequest) (*mcp.GoogleAIChatCompletionResponse, error) {
	// 转换请求类型
	serviceReq := &GoogleAIChatCompletionRequest{
		Model:       req.Model,
		Messages:    req.Messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		TopK:        req.TopK,
		Stream:      req.Stream,
		Options:     req.Options,
	}

	// 调用实际服务
	serviceResp, err := a.service.ChatCompletion(ctx, serviceReq)
	if err != nil {
		return nil, err
	}

	// 转换响应类型
	mcpResp := &mcp.GoogleAIChatCompletionResponse{
		ID:      serviceResp.ID,
		Object:  serviceResp.Object,
		Created: serviceResp.Created,
		Model:   serviceResp.Model,
		Choices: serviceResp.Choices,
		Usage:   serviceResp.Usage,
	}

	return mcpResp, nil
}

// ListModels 列出模型
func (a *GoogleAIServiceAdapter) ListModels(ctx context.Context) (map[string]*googleai.ModelConfig, error) {
	return a.service.ListModels(ctx)
}

// ValidateAPIKey 验证API密钥
func (a *GoogleAIServiceAdapter) ValidateAPIKey(ctx context.Context) error {
	return a.service.ValidateAPIKey(ctx)
}

// SetAPIKey 设置API密钥
func (a *GoogleAIServiceAdapter) SetAPIKey(key string) error {
	return a.service.SetAPIKey(key)
}

// GetModelConfig 获取模型配置
func (a *GoogleAIServiceAdapter) GetModelConfig(name string) (*googleai.ModelConfig, error) {
	return a.service.GetModelConfig(name)
}

// EnableModel 启用模型
func (a *GoogleAIServiceAdapter) EnableModel(name string) error {
	return a.service.EnableModel(name)
}

// DisableModel 禁用模型
func (a *GoogleAIServiceAdapter) DisableModel(name string) error {
	return a.service.DisableModel(name)
}

// MCPUserService MCP用户服务接口（适配器接口）
type MCPUserService interface {
	GetUser(ctx context.Context, id int64) (*dto.UserResponse, error)
	ListUsers(ctx context.Context, page, limit int64) ([]*dto.UserResponse, error)
}

// UserServiceAdapter 用户服务适配器，将repository适配为MCPUserService
type UserServiceAdapter struct {
	userRepo repository.UserRepository
}

// NewUserServiceAdapter 创建用户服务适配器
func NewUserServiceAdapter(repoManager repository.RepositoryManager) MCPUserService {
	return &UserServiceAdapter{
		userRepo: repoManager.User(),
	}
}

// GetUser 获取用户
func (a *UserServiceAdapter) GetUser(ctx context.Context, id int64) (*dto.UserResponse, error) {
	user, err := a.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// ListUsers 列出用户
func (a *UserServiceAdapter) ListUsers(ctx context.Context, page, limit int64) ([]*dto.UserResponse, error) {
	params := repository.NewPaginationParams(page, limit)
	users, err := a.userRepo.List(ctx, params)
	if err != nil {
		return nil, err
	}
	
	return users, nil
}

// MCPService MCP服务接口
type MCPService interface {
	// Initialize 初始化MCP服务
	Initialize(ctx context.Context, req *dto.MCPInitializeRequest) (*dto.MCPInitializeResponse, error)
	// ListTools 获取工具列表
	ListTools(ctx context.Context) (*dto.MCPToolsResponse, error)
	// ExecuteTool 执行工具
	ExecuteTool(ctx context.Context, req *dto.MCPExecuteRequest) (*dto.MCPExecuteResponse, error)
	// RegisterTool 注册工具
	RegisterTool(tool mcp.Tool) error
	// GetExecutionLog 获取执行日志
	GetExecutionLog(ctx context.Context, executionID string) (*dto.MCPToolExecutionLog, error)
	// ListExecutionLogs 列出执行日志
	ListExecutionLogs(ctx context.Context, userID *string, limit int) ([]*dto.MCPToolExecutionLog, error)
}

// MCPServiceImpl MCP服务实现
type MCPServiceImpl struct {
	toolRegistry    *mcp.ToolRegistry
	userService     MCPUserService
	googleaiService *GoogleAIService
	openaiService   *OpenAIService
	executionLogs   map[string]*dto.MCPToolExecutionLog
	executionMutex  sync.RWMutex
	sseClients      map[string]chan *dto.MCPSSEEvent
	sseClientsMutex sync.RWMutex
	logger          *zap.Logger
}

// NewMCPService 创建MCP服务
func NewMCPService(userService MCPUserService, googleaiService *GoogleAIService, openaiService *OpenAIService, logger *zap.Logger) MCPService {
	service := &MCPServiceImpl{
		toolRegistry:    mcp.NewToolRegistry(),
		userService:     userService,
		googleaiService: googleaiService,
		openaiService:   openaiService,
		executionLogs:   make(map[string]*dto.MCPToolExecutionLog),
		sseClients:      make(map[string]chan *dto.MCPSSEEvent),
		logger:          logger,
	}

	// 注册默认工具
	service.registerDefaultTools()

	return service
}

// registerDefaultTools 注册默认工具
func (s *MCPServiceImpl) registerDefaultTools() {
	// 注册回显工具
	echoTool := mcp.NewEchoTool()
	s.toolRegistry.Register(echoTool)

	s.logger.Info("Default MCP tools registered",
		logger.Module(logger.ModuleService),
		logger.Component("mcp"),
		zap.Strings("tools", s.toolRegistry.GetToolNames()))
}

// Initialize 初始化MCP服务
func (s *MCPServiceImpl) Initialize(ctx context.Context, req *dto.MCPInitializeRequest) (*dto.MCPInitializeResponse, error) {
	s.logger.Info("MCP service initialization requested",
		zap.String("protocolVersion", req.ProtocolVersion),
		zap.String("clientName", req.ClientInfo.Name),
		zap.String("clientVersion", req.ClientInfo.Version))

	// 验证协议版本
	if req.ProtocolVersion != "2024-11-05" {
		return nil, fmt.Errorf("unsupported protocol version: %s", req.ProtocolVersion)
	}

	response := &dto.MCPInitializeResponse{
		ProtocolVersion: "2024-11-05",
		Capabilities: dto.MCPCapabilities{
			Tools: &dto.MCPToolsCapability{
				ListChanged: true,
			},
			Logging: &dto.MCPLoggingCapability{},
		},
		ServerInfo: dto.MCPServerInfo{
			Name:    "Admin MCP Server",
			Version: "1.0.0",
		},
		Instructions: "This is an Admin MCP Server that provides tools for user management and system operations.",
	}

	s.logger.Info("MCP service initialized successfully",
		zap.String("serverName", response.ServerInfo.Name),
		zap.String("serverVersion", response.ServerInfo.Version))

	return response, nil
}

// ListTools 获取工具列表
func (s *MCPServiceImpl) ListTools(ctx context.Context) (*dto.MCPToolsResponse, error) {
	s.logger.Info("Listing MCP tools")

	tools := s.toolRegistry.ListTools()

	s.logger.Info("MCP tools listed successfully",
		zap.Int("toolCount", len(tools)),
		zap.Strings("toolNames", s.toolRegistry.GetToolNames()))

	return &dto.MCPToolsResponse{
		Tools: tools,
	}, nil
}

// ExecuteTool 执行工具
func (s *MCPServiceImpl) ExecuteTool(ctx context.Context, req *dto.MCPExecuteRequest) (*dto.MCPExecuteResponse, error) {
	executionID := uuid.New().String()
	startTime := time.Now()

	s.logger.Info("Executing MCP tool",
		zap.String("executionId", executionID),
		zap.String("toolName", req.Name),
		zap.Any("arguments", req.Arguments))

	// 创建执行日志
	executionLog := &dto.MCPToolExecutionLog{
		ID:        executionID,
		ToolName:  req.Name,
		Arguments: req.Arguments,
		StartTime: startTime,
		RequestID: getRequestIDFromContext(ctx),
	}

	// 从上下文获取用户ID（如果有）
	if userID := getUserIDFromContext(ctx); userID != "" {
		executionLog.UserID = &userID
	}

	// 保存执行日志
	s.executionMutex.Lock()
	s.executionLogs[executionID] = executionLog
	s.executionMutex.Unlock()

	// 获取工具
	tool, exists := s.toolRegistry.GetTool(req.Name)
	if !exists {
		err := fmt.Errorf("tool not found: %s", req.Name)
		s.updateExecutionLog(executionID, nil, &dto.MCPError{
			Code:    -32601,
			Message: err.Error(),
		})
		return nil, err
	}

	// 验证参数
	if err := tool.Validate(req.Arguments); err != nil {
		s.updateExecutionLog(executionID, nil, &dto.MCPError{
			Code:    -32602,
			Message: fmt.Sprintf("Invalid parameters: %v", err),
		})
		return nil, fmt.Errorf("invalid parameters: %v", err)
	}

	// 执行工具
	result, err := tool.Execute(ctx, req.Arguments)
	endTime := time.Now()
	duration := endTime.Sub(startTime)

	if err != nil {
		s.updateExecutionLog(executionID, nil, &dto.MCPError{
			Code:    -32603,
			Message: err.Error(),
		})
		s.logger.Error("MCP tool execution failed",
			zap.String("executionId", executionID),
			zap.String("toolName", req.Name),
			zap.Error(err),
			zap.Duration("duration", duration))
		return nil, err
	}

	// 更新执行日志
	s.updateExecutionLog(executionID, result, nil)

	s.logger.Info("MCP tool executed successfully",
		zap.String("executionId", executionID),
		zap.String("toolName", req.Name),
		zap.Duration("duration", duration),
		zap.Bool("isError", result.IsError))

	// 发送SSE事件
	s.broadcastSSEEvent(&dto.MCPSSEEvent{
		ID:    executionID,
		Event: "tool_execution",
		Data:  fmt.Sprintf(`{"toolName":"%s","executionId":"%s","status":"completed"}`, req.Name, executionID),
	})

	return result, nil
}

// RegisterTool 注册工具
func (s *MCPServiceImpl) RegisterTool(tool mcp.Tool) error {
	definition := tool.GetDefinition()
	s.toolRegistry.Register(tool)

	s.logger.Info("MCP tool registered",
		zap.String("toolName", definition.Name),
		zap.String("description", definition.Description))

	// 发送工具列表变更事件
	s.broadcastSSEEvent(&dto.MCPSSEEvent{
		Event: "tools_list_changed",
		Data:  fmt.Sprintf(`{"action":"added","toolName":"%s"}`, definition.Name),
	})

	return nil
}

// GetExecutionLog 获取执行日志
func (s *MCPServiceImpl) GetExecutionLog(ctx context.Context, executionID string) (*dto.MCPToolExecutionLog, error) {
	s.executionMutex.RLock()
	defer s.executionMutex.RUnlock()

	log, exists := s.executionLogs[executionID]
	if !exists {
		return nil, fmt.Errorf("execution log not found: %s", executionID)
	}

	return log, nil
}

// ListExecutionLogs 列出执行日志
func (s *MCPServiceImpl) ListExecutionLogs(ctx context.Context, userID *string, limit int) ([]*dto.MCPToolExecutionLog, error) {
	s.executionMutex.RLock()
	defer s.executionMutex.RUnlock()

	var logs []*dto.MCPToolExecutionLog
	count := 0

	for _, log := range s.executionLogs {
		if limit > 0 && count >= limit {
			break
		}

		// 如果指定了用户ID，只返回该用户的日志
		if userID != nil && (log.UserID == nil || *log.UserID != *userID) {
			continue
		}

		logs = append(logs, log)
		count++
	}

	return logs, nil
}

// updateExecutionLog 更新执行日志
func (s *MCPServiceImpl) updateExecutionLog(executionID string, result *dto.MCPExecuteResponse, mcpError *dto.MCPError) {
	s.executionMutex.Lock()
	defer s.executionMutex.Unlock()

	if log, exists := s.executionLogs[executionID]; exists {
		endTime := time.Now()
		duration := endTime.Sub(log.StartTime)

		log.EndTime = &endTime
		log.Duration = &duration
		log.Result = result
		log.Error = mcpError
	}
}

// AddSSEClient 添加SSE客户端
func (s *MCPServiceImpl) AddSSEClient(clientID string) chan *dto.MCPSSEEvent {
	s.sseClientsMutex.Lock()
	defer s.sseClientsMutex.Unlock()

	eventChan := make(chan *dto.MCPSSEEvent, 100)
	s.sseClients[clientID] = eventChan

	s.logger.Info("SSE client added", zap.String("clientId", clientID))

	return eventChan
}

// RemoveSSEClient 移除SSE客户端
func (s *MCPServiceImpl) RemoveSSEClient(clientID string) {
	s.sseClientsMutex.Lock()
	defer s.sseClientsMutex.Unlock()

	if eventChan, exists := s.sseClients[clientID]; exists {
		close(eventChan)
		delete(s.sseClients, clientID)
		s.logger.Info("SSE client removed", zap.String("clientId", clientID))
	}
}

// broadcastSSEEvent 广播SSE事件
func (s *MCPServiceImpl) broadcastSSEEvent(event *dto.MCPSSEEvent) {
	s.sseClientsMutex.RLock()
	defer s.sseClientsMutex.RUnlock()

	for clientID, eventChan := range s.sseClients {
		select {
		case eventChan <- event:
			// 事件发送成功
		default:
			// 通道已满，移除客户端
			s.logger.Warn("SSE client channel full, removing client", zap.String("clientId", clientID))
			go s.RemoveSSEClient(clientID)
		}
	}
}

// getUserIDFromContext 从上下文获取用户ID
func getUserIDFromContext(ctx context.Context) string {
	if userID := ctx.Value("userID"); userID != nil {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

// getRequestIDFromContext 从上下文获取请求ID
func getRequestIDFromContext(ctx context.Context) string {
	if requestID := ctx.Value("request_id"); requestID != nil {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}