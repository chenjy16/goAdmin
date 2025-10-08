package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"admin/internal/dto"
	"admin/internal/logger"
	"admin/internal/response"
	"admin/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// MCPController MCP控制器
type MCPController struct {
	*BaseController
	mcpService service.MCPService
	logger     *zap.Logger
}

// NewMCPController 创建MCP控制器
func NewMCPController(mcpService service.MCPService, logger *zap.Logger) *MCPController {
	return &MCPController{
		BaseController: NewBaseController(),
		mcpService:     mcpService,
		logger:         logger,
	}
}

// Initialize 初始化MCP服务
func (mc *MCPController) Initialize(c *gin.Context) {
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("mcp"),
		logger.Operation("initialize"),
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path))

	var req dto.MCPInitializeRequest

	// 使用基础控制器的统一绑定和验证方法
	if err := mc.BindAndValidate(c, &req); err != nil {
		logger.WarnCtx(c.Request.Context(), logger.MsgAPIValidation,
			logger.Module(logger.ModuleController),
			logger.Component("mcp"),
			logger.Operation("initialize"),
			logger.ZapError(err))
		return
	}

	result, err := mc.mcpService.Initialize(c.Request.Context(), &req)
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("mcp"),
			logger.Operation("initialize"),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIResponse,
		logger.Module(logger.ModuleController),
		logger.Component("mcp"),
		logger.Operation("initialize"),
		logger.String("protocolVersion", result.ProtocolVersion),
		logger.Int("status", http.StatusOK))

	response.Success(c, http.StatusOK, "MCP service initialized successfully", result)
}

// ListTools 获取工具列表
func (mc *MCPController) ListTools(c *gin.Context) {
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("mcp"),
		logger.Operation("list_tools"),
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path))

	result, err := mc.mcpService.ListTools(c.Request.Context())
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("mcp"),
			logger.Operation("list_tools"),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIResponse,
		logger.Module(logger.ModuleController),
		logger.Component("mcp"),
		logger.Operation("list_tools"),
		logger.Int("toolCount", len(result.Tools)),
		logger.Int("status", http.StatusOK))

	response.Success(c, http.StatusOK, "Tools retrieved successfully", result)
}

// ExecuteTool 执行工具
func (mc *MCPController) ExecuteTool(c *gin.Context) {
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("mcp"),
		logger.Operation("execute_tool"),
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path))

	var req dto.MCPExecuteRequest

	// 使用基础控制器的统一绑定和验证方法
	if err := mc.BindAndValidate(c, &req); err != nil {
		logger.WarnCtx(c.Request.Context(), logger.MsgAPIValidation,
			logger.Module(logger.ModuleController),
			logger.Component("mcp"),
			logger.Operation("execute_tool"),
			logger.String("toolName", req.Name),
			logger.ZapError(err))
		return
	}

	result, err := mc.mcpService.ExecuteTool(c.Request.Context(), &req)
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("mcp"),
			logger.Operation("execute_tool"),
			logger.String("toolName", req.Name),
			logger.ZapError(err))
		mc.HandleError(c, err)
		return
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIResponse,
		logger.Module(logger.ModuleController),
		logger.Component("mcp"),
		logger.Operation("execute_tool"),
		logger.String("toolName", req.Name),
		logger.Bool("isError", result.IsError),
		logger.Int("status", http.StatusOK))

	response.Success(c, http.StatusOK, "Tool executed successfully", result)
}

// StreamSSE SSE流式端点
func (mc *MCPController) StreamSSE(c *gin.Context) {
	clientID := uuid.New().String()

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("mcp"),
		logger.Operation("stream_sse"),
		logger.String("clientId", clientID),
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path))

	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Cache-Control")

	// 添加SSE客户端
	eventChan := mc.mcpService.(*service.MCPServiceImpl).AddSSEClient(clientID)
	defer mc.mcpService.(*service.MCPServiceImpl).RemoveSSEClient(clientID)

	// 发送初始连接事件
	initialEvent := &dto.MCPSSEEvent{
		ID:    uuid.New().String(),
		Event: "connected",
		Data:  fmt.Sprintf(`{"clientId":"%s","timestamp":"%s"}`, clientID, time.Now().Format(time.RFC3339)),
	}

	mc.writeSSEEvent(c, initialEvent)

	// 创建上下文用于取消
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	// 监听客户端断开连接
	go func() {
		<-c.Request.Context().Done()
		cancel()
	}()

	// 定期发送心跳
	heartbeatTicker := time.NewTicker(30 * time.Second)
	defer heartbeatTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.InfoCtx(c.Request.Context(), "SSE client disconnected",
				logger.Module(logger.ModuleController),
				logger.Component("mcp"),
				logger.String("clientId", clientID))
			return

		case event := <-eventChan:
			if err := mc.writeSSEEvent(c, event); err != nil {
				logger.ErrorCtx(c.Request.Context(), "Failed to write SSE event",
					logger.Module(logger.ModuleController),
					logger.Component("mcp"),
					logger.String("clientId", clientID),
					logger.ZapError(err))
				return
			}

		case <-heartbeatTicker.C:
			heartbeatEvent := &dto.MCPSSEEvent{
				ID:    uuid.New().String(),
				Event: "heartbeat",
				Data:  fmt.Sprintf(`{"timestamp":"%s"}`, time.Now().Format(time.RFC3339)),
			}
			if err := mc.writeSSEEvent(c, heartbeatEvent); err != nil {
				logger.ErrorCtx(c.Request.Context(), "Failed to write heartbeat event",
					logger.Module(logger.ModuleController),
					logger.Component("mcp"),
					logger.String("clientId", clientID),
					logger.ZapError(err))
				return
			}
		}
	}
}

// writeSSEEvent 写入SSE事件
func (mc *MCPController) writeSSEEvent(c *gin.Context, event *dto.MCPSSEEvent) error {
	writer := c.Writer

	if event.ID != "" {
		if _, err := fmt.Fprintf(writer, "id: %s\n", event.ID); err != nil {
			return err
		}
	}

	if event.Event != "" {
		if _, err := fmt.Fprintf(writer, "event: %s\n", event.Event); err != nil {
			return err
		}
	}

	if event.Retry > 0 {
		if _, err := fmt.Fprintf(writer, "retry: %d\n", event.Retry); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(writer, "data: %s\n\n", event.Data); err != nil {
		return err
	}

	// 刷新缓冲区
	if flusher, ok := writer.(http.Flusher); ok {
		flusher.Flush()
	}

	return nil
}

// GetExecutionLog 获取执行日志
func (mc *MCPController) GetExecutionLog(c *gin.Context) {
	executionID := c.Param("id")
	if executionID == "" {
		response.Error(c, http.StatusBadRequest, "Execution ID is required", "INVALID_EXECUTION_ID")
		return
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("mcp"),
		logger.Operation("get_execution_log"),
		logger.String("executionId", executionID),
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path))

	result, err := mc.mcpService.GetExecutionLog(c.Request.Context(), executionID)
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("mcp"),
			logger.Operation("get_execution_log"),
			logger.String("executionId", executionID),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIResponse,
		logger.Module(logger.ModuleController),
		logger.Component("mcp"),
		logger.Operation("get_execution_log"),
		logger.String("executionId", executionID),
		logger.Int("status", http.StatusOK))

	response.Success(c, http.StatusOK, "Execution log retrieved successfully", result)
}

// ListExecutionLogs 列出执行日志
func (mc *MCPController) ListExecutionLogs(c *gin.Context) {
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("mcp"),
		logger.Operation("list_execution_logs"),
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path))

	// 解析查询参数
	var userID *string
	if uid := c.Query("user_id"); uid != "" {
		userID = &uid
	}

	limit := 50 // 默认限制
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	result, err := mc.mcpService.ListExecutionLogs(c.Request.Context(), userID, limit)
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("mcp"),
			logger.Operation("list_execution_logs"),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIResponse,
		logger.Module(logger.ModuleController),
		logger.Component("mcp"),
		logger.Operation("list_execution_logs"),
		logger.Int("logCount", len(result)),
		logger.Int("status", http.StatusOK))

	response.Success(c, http.StatusOK, "Execution logs retrieved successfully", map[string]interface{}{
		"logs":  result,
		"count": len(result),
		"limit": limit,
	})
}