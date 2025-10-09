package controllers

import (
	"net/http"

	"go-springAi/internal/logger"
	"go-springAi/internal/response"
	"go-springAi/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AIAssistantController AI助手控制器
type AIAssistantController struct {
	*BaseController
	aiAssistantService *service.AIAssistantService
	logger             *zap.Logger
}

// NewAIAssistantController 创建AI助手控制器
func NewAIAssistantController(aiAssistantService *service.AIAssistantService, logger *zap.Logger) *AIAssistantController {
	return &AIAssistantController{
		BaseController:     NewBaseController(),
		aiAssistantService: aiAssistantService,
		logger:             logger,
	}
}

// Chat AI助手聊天接口
func (ac *AIAssistantController) Chat(c *gin.Context) {
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("ai_assistant"),
		logger.Operation("chat"),
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path))

	var req service.ChatRequest

	// 使用基础控制器的统一绑定和验证方法
	if err := ac.BindAndValidate(c, &req); err != nil {
		logger.WarnCtx(c.Request.Context(), logger.MsgAPIValidation,
			logger.Module(logger.ModuleController),
			logger.Component("ai_assistant"),
			logger.Operation("chat"),
			logger.ZapError(err))
		return
	}

	// 设置默认模型
	if req.Model == "" {
		req.Model = "gpt-3.5-turbo"
	}

	result, err := ac.aiAssistantService.Chat(c.Request.Context(), &req)
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("ai_assistant"),
			logger.Operation("chat"),
			logger.ZapError(err),
			logger.String("error_message", err.Error()),
			logger.String("model", req.Model))
		c.Error(err)
		return
	}

	var finishReason string
	var toolCallsCount int
	if len(result.Choices) > 0 {
		finishReason = result.Choices[0].FinishReason
		toolCallsCount = len(result.Choices[0].ToolCalls)
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIResponse,
		logger.Module(logger.ModuleController),
		logger.Component("ai_assistant"),
		logger.Operation("chat"),
		logger.String("finish_reason", finishReason),
		logger.Int("tool_calls", toolCallsCount),
		logger.Int("status", http.StatusOK))

	response.Success(c, http.StatusOK, "Chat completed successfully", result)
}

// Initialize 初始化AI助手
func (ac *AIAssistantController) Initialize(c *gin.Context) {
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("ai_assistant"),
		logger.Operation("initialize"),
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path))

	err := ac.aiAssistantService.Initialize(c.Request.Context())
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("ai_assistant"),
			logger.Operation("initialize"),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIResponse,
		logger.Module(logger.ModuleController),
		logger.Component("ai_assistant"),
		logger.Operation("initialize"),
		logger.Int("status", http.StatusOK))

	response.Success(c, http.StatusOK, "AI assistant initialized successfully", gin.H{
		"status": "initialized",
	})
}