package controllers

import (
	"net/http"

	"go-springAi/internal/dto"
	"go-springAi/internal/logger"
	"go-springAi/internal/response"
	"go-springAi/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GoogleAIController Google AI 控制器
type GoogleAIController struct {
	BaseController
	googleaiService *service.GoogleAIService
	logger          *zap.Logger
}

// NewGoogleAIController 创建 Google AI 控制器
func NewGoogleAIController(googleaiService *service.GoogleAIService, logger *zap.Logger) *GoogleAIController {
	return &GoogleAIController{
		BaseController:  *NewBaseController(),
		googleaiService: googleaiService,
		logger:          logger,
	}
}

// ChatCompletion 聊天完成接口
func (gc *GoogleAIController) ChatCompletion(c *gin.Context) {
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("googleai"),
		logger.Operation("chat_completion"),
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path))

	var req dto.GoogleAIChatRequest

	// 使用基础控制器的统一绑定和验证方法
	if err := gc.BindAndValidate(c, &req); err != nil {
		logger.WarnCtx(c.Request.Context(), logger.MsgAPIValidation,
			logger.Module(logger.ModuleController),
			logger.Component("googleai"),
			logger.Operation("chat_completion"),
			logger.ZapError(err))
		return
	}

	// 转换为服务层请求
	serviceReq := &service.GoogleAIChatCompletionRequest{
		Model:       req.Model,
		Messages:    req.Messages,
		MaxTokens:   req.MaxTokens,
		Temperature: convertFloat64ToFloat32Ptr(req.Temperature),
		TopP:        convertFloat64ToFloat32Ptr(req.TopP),
		TopK:        req.TopK,
		Stream:      req.Stream,
		Options:     req.Options,
	}

	result, err := gc.googleaiService.ChatCompletion(c.Request.Context(), serviceReq)
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("googleai"),
			logger.Operation("chat_completion"),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIResponse,
		logger.Module(logger.ModuleController),
		logger.Component("googleai"),
		logger.Operation("chat_completion"),
		logger.String("response_id", result.ID),
		logger.Int("status", http.StatusOK))

	response.Success(c, http.StatusOK, "Chat completion successful", result)
}

// ListModels 列出可用模型
func (gc *GoogleAIController) ListModels(c *gin.Context) {
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("googleai"),
		logger.Operation("list_models"),
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path))

	models, err := gc.googleaiService.ListModels(c.Request.Context())
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("googleai"),
			logger.Operation("list_models"),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIResponse,
		logger.Module(logger.ModuleController),
		logger.Component("googleai"),
		logger.Operation("list_models"),
		logger.Int("model_count", len(models)),
		logger.Int("status", http.StatusOK))

	response.Success(c, http.StatusOK, "Models retrieved successfully", models)
}

// ValidateAPIKey 验证 API 密钥
func (gc *GoogleAIController) ValidateAPIKey(c *gin.Context) {
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("googleai"),
		logger.Operation("validate_api_key"),
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path))

	err := gc.googleaiService.ValidateAPIKey(c.Request.Context())
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("googleai"),
			logger.Operation("validate_api_key"),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIResponse,
		logger.Module(logger.ModuleController),
		logger.Component("googleai"),
		logger.Operation("validate_api_key"),
		logger.Int("status", http.StatusOK))

	response.Success(c, http.StatusOK, "API key is valid", nil)
}

// SetAPIKey 设置 API 密钥
func (gc *GoogleAIController) SetAPIKey(c *gin.Context) {
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("googleai"),
		logger.Operation("set_api_key"),
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path))

	var req dto.GoogleAISetAPIKeyRequest

	// 使用基础控制器的统一绑定和验证方法
	if err := gc.BindAndValidate(c, &req); err != nil {
		logger.WarnCtx(c.Request.Context(), logger.MsgAPIValidation,
			logger.Module(logger.ModuleController),
			logger.Component("googleai"),
			logger.Operation("set_api_key"),
			logger.ZapError(err))
		return
	}

	err := gc.googleaiService.SetAPIKey(req.APIKey)
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("googleai"),
			logger.Operation("set_api_key"),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIResponse,
		logger.Module(logger.ModuleController),
		logger.Component("googleai"),
		logger.Operation("set_api_key"),
		logger.Int("status", http.StatusOK))

	response.Success(c, http.StatusOK, "API key set successfully", nil)
}

// GetModelConfig 获取模型配置
func (gc *GoogleAIController) GetModelConfig(c *gin.Context) {
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("googleai"),
		logger.Operation("get_model_config"),
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path))

	modelName := c.Param("model")
	if modelName == "" {
		logger.WarnCtx(c.Request.Context(), logger.MsgAPIValidation,
			logger.Module(logger.ModuleController),
			logger.Component("googleai"),
			logger.Operation("get_model_config"),
			logger.String("error", "model name is required"))
		response.Error(c, http.StatusBadRequest, "Model name is required", "INVALID_MODEL_NAME")
		return
	}

	config, err := gc.googleaiService.GetModelConfig(modelName)
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("googleai"),
			logger.Operation("get_model_config"),
			logger.String("model", modelName),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIResponse,
		logger.Module(logger.ModuleController),
		logger.Component("googleai"),
		logger.Operation("get_model_config"),
		logger.String("model", modelName),
		logger.Int("status", http.StatusOK))

	response.Success(c, http.StatusOK, "Model configuration retrieved successfully", config)
}

// EnableModel 启用模型
func (gc *GoogleAIController) EnableModel(c *gin.Context) {
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("googleai"),
		logger.Operation("enable_model"),
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path))

	modelName := c.Param("model")
	if modelName == "" {
		logger.WarnCtx(c.Request.Context(), logger.MsgAPIValidation,
			logger.Module(logger.ModuleController),
			logger.Component("googleai"),
			logger.Operation("enable_model"),
			logger.String("error", "model name is required"))
		response.Error(c, http.StatusBadRequest, "Model name is required", "INVALID_MODEL_NAME")
		return
	}

	err := gc.googleaiService.EnableModel(modelName)
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("googleai"),
			logger.Operation("enable_model"),
			logger.String("model", modelName),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIResponse,
		logger.Module(logger.ModuleController),
		logger.Component("googleai"),
		logger.Operation("enable_model"),
		logger.String("model", modelName),
		logger.Int("status", http.StatusOK))

	response.Success(c, http.StatusOK, "Model enabled successfully", nil)
}

// DisableModel 禁用模型
func (gc *GoogleAIController) DisableModel(c *gin.Context) {
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("googleai"),
		logger.Operation("disable_model"),
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path))

	modelName := c.Param("model")
	if modelName == "" {
		logger.WarnCtx(c.Request.Context(), logger.MsgAPIValidation,
			logger.Module(logger.ModuleController),
			logger.Component("googleai"),
			logger.Operation("disable_model"),
			logger.String("error", "model name is required"))
		response.Error(c, http.StatusBadRequest, "Model name is required", "INVALID_MODEL_NAME")
		return
	}

	err := gc.googleaiService.DisableModel(modelName)
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("googleai"),
			logger.Operation("disable_model"),
			logger.String("model", modelName),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIResponse,
		logger.Module(logger.ModuleController),
		logger.Component("googleai"),
		logger.Operation("disable_model"),
		logger.String("model", modelName),
		logger.Int("status", http.StatusOK))

	response.Success(c, http.StatusOK, "Model disabled successfully", nil)
}