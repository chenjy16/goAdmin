package controllers

import (
	"io"
	"net/http"

	"go-springAi/internal/dto"
	"go-springAi/internal/logger"
	"go-springAi/internal/provider"
	"go-springAi/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AIController 统一AI控制器
type AIController struct {
	BaseController
	providerManager *provider.Manager
	logger          *zap.Logger
}

// NewAIController 创建统一AI控制器
func NewAIController(providerManager *provider.Manager, logger *zap.Logger) *AIController {
	return &AIController{
		BaseController:  *NewBaseController(),
		providerManager: providerManager,
		logger:          logger,
	}
}

// ChatCompletion 统一聊天完成接口
func (ac *AIController) ChatCompletion(c *gin.Context) {
	providerType := c.Param("provider")
	
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("ai"),
		logger.Operation("chat_completion"),
		logger.String("provider", providerType),
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path))

	var req dto.UnifiedChatRequest

	// 使用基础控制器的统一绑定和验证方法
	if err := ac.BindAndValidate(c, &req); err != nil {
		logger.WarnCtx(c.Request.Context(), logger.MsgAPIValidation,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("chat_completion"),
			logger.String("provider", providerType),
			logger.ZapError(err))
		return
	}

	// 获取Provider
	prov, err := ac.providerManager.GetProvider(provider.ProviderType(providerType))
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("chat_completion"),
			logger.String("provider", providerType),
			logger.ZapError(err))
		response.Error(c, http.StatusBadRequest, "Invalid provider", err.Error())
		return
	}

	// 转换为Provider请求
	providerReq := &provider.ChatRequest{
		Model:       req.Model,
		Messages:    convertToProviderMessages(req.Messages),
		MaxTokens:   req.MaxTokens,
		Temperature: convertFloat64ToFloat32Ptr(req.Temperature),
		TopP:        convertFloat64ToFloat32Ptr(req.TopP),
		TopK:        req.TopK,
		Stream:      req.Stream,
		Options:     req.Options,
	}

	// 处理流式响应
	if req.Stream {
		stream, err := prov.ChatCompletionStream(c.Request.Context(), providerReq)
		if err != nil {
			logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
				logger.Module(logger.ModuleController),
				logger.Component("ai"),
				logger.Operation("chat_completion_stream"),
				logger.String("provider", providerType),
				logger.ZapError(err))
			c.Error(err)
			return
		}
		defer stream.Close()

		// 设置SSE头
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("Access-Control-Allow-Origin", "*")

		// 流式传输数据
		c.Stream(func(w io.Writer) bool {
			buffer := make([]byte, 1024)
			n, err := stream.Read(buffer)
			if err != nil {
				return false
			}
			w.Write(buffer[:n])
			return true
		})
		return
	}

	// 处理普通响应
	result, err := prov.ChatCompletion(c.Request.Context(), providerReq)
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("chat_completion"),
			logger.String("provider", providerType),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIResponse,
		logger.Module(logger.ModuleController),
		logger.Component("ai"),
		logger.Operation("chat_completion"),
		logger.String("provider", providerType),
		logger.String("response_id", result.ID),
		logger.Int("status", http.StatusOK))

	response.Success(c, http.StatusOK, "Chat completion successful", result)
}

// ListModels 列出指定提供商的模型
func (ac *AIController) ListModels(c *gin.Context) {
	providerType := c.Param("provider")
	
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("ai"),
		logger.Operation("list_models"),
		logger.String("provider", providerType))

	// 获取Provider
	prov, err := ac.providerManager.GetProvider(provider.ProviderType(providerType))
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("list_models"),
			logger.String("provider", providerType),
			logger.ZapError(err))
		response.Error(c, http.StatusBadRequest, "Invalid provider", err.Error())
		return
	}

	models, err := prov.ListModels(c.Request.Context())
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("list_models"),
			logger.String("provider", providerType),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Models retrieved successfully", gin.H{
		"provider": providerType,
		"models":   models,
	})
}

// ListProviders 列出所有可用的提供商
func (ac *AIController) ListProviders(c *gin.Context) {
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("ai"),
		logger.Operation("list_providers"))

	providers := ac.providerManager.ListProviders()

	response.Success(c, http.StatusOK, "Providers retrieved successfully", gin.H{
		"providers": providers,
	})
}

// GetModelConfig 获取指定提供商的模型配置
func (ac *AIController) GetModelConfig(c *gin.Context) {
	providerType := c.Param("provider")
	modelName := c.Param("model")
	
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("ai"),
		logger.Operation("get_model_config"),
		logger.String("provider", providerType),
		logger.String("model", modelName))

	// 获取Provider
	prov, err := ac.providerManager.GetProvider(provider.ProviderType(providerType))
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("get_model_config"),
			logger.String("provider", providerType),
			logger.ZapError(err))
		response.Error(c, http.StatusBadRequest, "Invalid provider", err.Error())
		return
	}

	config, err := prov.GetModelConfig(modelName)
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("get_model_config"),
			logger.String("provider", providerType),
			logger.String("model", modelName),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Model config retrieved successfully", gin.H{
		"provider": providerType,
		"model":    modelName,
		"config":   config,
	})
}

// EnableModel 启用指定提供商的模型
func (ac *AIController) EnableModel(c *gin.Context) {
	providerType := c.Param("provider")
	modelName := c.Param("model")
	
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("ai"),
		logger.Operation("enable_model"),
		logger.String("provider", providerType),
		logger.String("model", modelName))

	// 获取Provider
	prov, err := ac.providerManager.GetProvider(provider.ProviderType(providerType))
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("enable_model"),
			logger.String("provider", providerType),
			logger.ZapError(err))
		response.Error(c, http.StatusBadRequest, "Invalid provider", err.Error())
		return
	}

	err = prov.EnableModel(modelName)
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("enable_model"),
			logger.String("provider", providerType),
			logger.String("model", modelName),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Model enabled successfully", gin.H{
		"provider": providerType,
		"model":    modelName,
	})
}

// DisableModel 禁用指定提供商的模型
func (ac *AIController) DisableModel(c *gin.Context) {
	providerType := c.Param("provider")
	modelName := c.Param("model")
	
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("ai"),
		logger.Operation("disable_model"),
		logger.String("provider", providerType),
		logger.String("model", modelName))

	// 获取Provider
	prov, err := ac.providerManager.GetProvider(provider.ProviderType(providerType))
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("disable_model"),
			logger.String("provider", providerType),
			logger.ZapError(err))
		response.Error(c, http.StatusBadRequest, "Invalid provider", err.Error())
		return
	}

	err = prov.DisableModel(modelName)
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("disable_model"),
			logger.String("provider", providerType),
			logger.String("model", modelName),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Model disabled successfully", gin.H{
		"provider": providerType,
		"model":    modelName,
	})
}

// ValidateAPIKey 验证指定提供商的API密钥
func (ac *AIController) ValidateAPIKey(c *gin.Context) {
	providerType := c.Param("provider")
	
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("ai"),
		logger.Operation("validate_api_key"),
		logger.String("provider", providerType))

	// 获取Provider
	prov, err := ac.providerManager.GetProvider(provider.ProviderType(providerType))
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("validate_api_key"),
			logger.String("provider", providerType),
			logger.ZapError(err))
		response.Error(c, http.StatusBadRequest, "Invalid provider", err.Error())
		return
	}

	err = prov.ValidateAPIKey(c.Request.Context())
	if err != nil {
		response.Success(c, http.StatusOK, "API key validation failed", gin.H{
			"provider": providerType,
			"valid":    false,
			"message":  err.Error(),
		})
		return
	}

	response.Success(c, http.StatusOK, "API key is valid", gin.H{
		"provider": providerType,
		"valid":    true,
	})
}

// SetAPIKey 设置指定提供商的API密钥
func (ac *AIController) SetAPIKey(c *gin.Context) {
	providerType := c.Param("provider")
	
	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("ai"),
		logger.Operation("set_api_key"),
		logger.String("provider", providerType))

	var req dto.SetAPIKeyRequest
	if err := ac.BindAndValidate(c, &req); err != nil {
		logger.WarnCtx(c.Request.Context(), logger.MsgAPIValidation,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("set_api_key"),
			logger.String("provider", providerType),
			logger.ZapError(err))
		return
	}

	// 获取Provider
	prov, err := ac.providerManager.GetProvider(provider.ProviderType(providerType))
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("set_api_key"),
			logger.String("provider", providerType),
			logger.ZapError(err))
		response.Error(c, http.StatusBadRequest, "Invalid provider", err.Error())
		return
	}

	err = prov.SetAPIKey(req.APIKey)
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("set_api_key"),
			logger.String("provider", providerType),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "API key set successfully", gin.H{
		"provider": providerType,
	})
}

// 辅助函数：转换DTO消息为Provider消息
func convertToProviderMessages(messages []dto.UnifiedMessage) []provider.Message {
	result := make([]provider.Message, len(messages))
	for i, msg := range messages {
		result[i] = provider.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	return result
}

// 辅助函数：转换float64指针为float32指针
func convertFloat64ToFloat32Ptr(f64Ptr *float64) *float32 {
	if f64Ptr == nil {
		return nil
	}
	f32 := float32(*f64Ptr)
	return &f32
}