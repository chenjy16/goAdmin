package controllers

import (
	"net/http"

	"admin/internal/dto"
	"admin/internal/logger"
	"admin/internal/response"
	"admin/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OpenAIController OpenAI 控制器
type OpenAIController struct {
	BaseController
	openaiService *service.OpenAIService
	logger        *zap.Logger
}

// NewOpenAIController 创建 OpenAI 控制器
func NewOpenAIController(openaiService *service.OpenAIService, logger *zap.Logger) *OpenAIController {
	return &OpenAIController{
		BaseController: *NewBaseController(),
		openaiService:  openaiService,
		logger:         logger,
	}
}

// ChatCompletion 聊天完成接口
func (oc *OpenAIController) ChatCompletion(c *gin.Context) {
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("openai"),
		logger.Operation("chat_completion"),
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path))

	var req dto.OpenAIChatRequest

	// 使用基础控制器的统一绑定和验证方法
	if err := oc.BindAndValidate(c, &req); err != nil {
		logger.WarnCtx(c.Request.Context(), logger.MsgAPIValidation,
			logger.Module(logger.ModuleController),
			logger.Component("openai"),
			logger.Operation("chat_completion"),
			logger.ZapError(err))
		return
	}

	// 转换为服务层请求
	serviceReq := &service.ChatCompletionRequest{
		Model:       req.Model,
		Messages:    req.Messages,
		MaxTokens:   req.MaxTokens,
		Temperature: convertFloat64ToFloat32Ptr(req.Temperature),
		TopP:        convertFloat64ToFloat32Ptr(req.TopP),
		Stream:      req.Stream,
		Options:     req.Options,
	}

	result, err := oc.openaiService.ChatCompletion(c.Request.Context(), serviceReq)
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("openai"),
			logger.Operation("chat_completion"),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIResponse,
		logger.Module(logger.ModuleController),
		logger.Component("openai"),
		logger.Operation("chat_completion"),
		logger.String("response_id", result.ID),
		logger.Int("status", http.StatusOK))

	response.Success(c, http.StatusOK, "Chat completion successful", result)
}

// ListModels 列出可用模型
func (oc *OpenAIController) ListModels(c *gin.Context) {
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("openai"),
		logger.Operation("list_models"),
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path))

	models, err := oc.openaiService.ListModels(c.Request.Context())
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("openai"),
			logger.Operation("list_models"),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIResponse,
		logger.Module(logger.ModuleController),
		logger.Component("openai"),
		logger.Operation("list_models"),
		logger.Int("model_count", len(models)),
		logger.Int("status", http.StatusOK))

	response.Success(c, http.StatusOK, "Models retrieved successfully", models)
}

// ValidateAPIKey 验证 API 密钥
func (oc *OpenAIController) ValidateAPIKey(c *gin.Context) {
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("openai"),
		logger.Operation("validate_api_key"),
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path))

	err := oc.openaiService.ValidateAPIKey(c.Request.Context())
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("openai"),
			logger.Operation("validate_api_key"),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIResponse,
		logger.Module(logger.ModuleController),
		logger.Component("openai"),
		logger.Operation("validate_api_key"),
		logger.Int("status", http.StatusOK))

	response.Success(c, http.StatusOK, "API key is valid", nil)
}

// SetAPIKey 设置 API 密钥
func (oc *OpenAIController) SetAPIKey(c *gin.Context) {
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("openai"),
		logger.Operation("set_api_key"),
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path))

	var req dto.OpenAISetAPIKeyRequest

	// 使用基础控制器的统一绑定和验证方法
	if err := oc.BindAndValidate(c, &req); err != nil {
		logger.WarnCtx(c.Request.Context(), logger.MsgAPIValidation,
			logger.Module(logger.ModuleController),
			logger.Component("openai"),
			logger.Operation("set_api_key"),
			logger.ZapError(err))
		return
	}

	err := oc.openaiService.SetAPIKey(req.APIKey)
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("openai"),
			logger.Operation("set_api_key"),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIResponse,
		logger.Module(logger.ModuleController),
		logger.Component("openai"),
		logger.Operation("set_api_key"),
		logger.Int("status", http.StatusOK))

	response.Success(c, http.StatusOK, "API key set successfully", nil)
}

// GetModelConfig 获取模型配置
func (oc *OpenAIController) GetModelConfig(c *gin.Context) {
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("openai"),
		logger.Operation("get_model_config"),
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path))

	modelName := c.Param("model")
	if modelName == "" {
		logger.WarnCtx(c.Request.Context(), logger.MsgAPIValidation,
			logger.Module(logger.ModuleController),
			logger.Component("openai"),
			logger.Operation("get_model_config"),
			logger.String("error", "model name is required"))
		response.Error(c, http.StatusBadRequest, "Model name is required", "INVALID_MODEL_NAME")
		return
	}

	config, err := oc.openaiService.GetModelConfig(modelName)
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("openai"),
			logger.Operation("get_model_config"),
			logger.String("model", modelName),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIResponse,
		logger.Module(logger.ModuleController),
		logger.Component("openai"),
		logger.Operation("get_model_config"),
		logger.String("model", modelName),
		logger.Int("status", http.StatusOK))

	response.Success(c, http.StatusOK, "Model configuration retrieved successfully", config)
}

// EnableModel 启用模型
func (oc *OpenAIController) EnableModel(c *gin.Context) {
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("openai"),
		logger.Operation("enable_model"),
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path))

	modelName := c.Param("model")
	if modelName == "" {
		logger.WarnCtx(c.Request.Context(), logger.MsgAPIValidation,
			logger.Module(logger.ModuleController),
			logger.Component("openai"),
			logger.Operation("enable_model"),
			logger.String("error", "model name is required"))
		response.Error(c, http.StatusBadRequest, "Model name is required", "INVALID_MODEL_NAME")
		return
	}

	err := oc.openaiService.EnableModel(modelName)
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("openai"),
			logger.Operation("enable_model"),
			logger.String("model", modelName),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIResponse,
		logger.Module(logger.ModuleController),
		logger.Component("openai"),
		logger.Operation("enable_model"),
		logger.String("model", modelName),
		logger.Int("status", http.StatusOK))

	response.Success(c, http.StatusOK, "Model enabled successfully", nil)
}

// DisableModel 禁用模型
func (oc *OpenAIController) DisableModel(c *gin.Context) {
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("openai"),
		logger.Operation("disable_model"),
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path))

	modelName := c.Param("model")
	if modelName == "" {
		logger.WarnCtx(c.Request.Context(), logger.MsgAPIValidation,
			logger.Module(logger.ModuleController),
			logger.Component("openai"),
			logger.Operation("disable_model"),
			logger.String("error", "model name is required"))
		response.Error(c, http.StatusBadRequest, "Model name is required", "INVALID_MODEL_NAME")
		return
	}

	err := oc.openaiService.DisableModel(modelName)
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("openai"),
			logger.Operation("disable_model"),
			logger.String("model", modelName),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIResponse,
		logger.Module(logger.ModuleController),
		logger.Component("openai"),
		logger.Operation("disable_model"),
		logger.String("model", modelName),
		logger.Int("status", http.StatusOK))

	response.Success(c, http.StatusOK, "Model disabled successfully", nil)
}

// convertFloat64ToFloat32Ptr 将 *float64 转换为 *float32
func convertFloat64ToFloat32Ptr(f64 *float64) *float32 {
	if f64 == nil {
		return nil
	}
	f32 := float32(*f64)
	return &f32
}