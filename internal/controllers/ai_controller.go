package controllers

import (
	"net/http"
	"strconv"

	"go-springAi/internal/dto"
	"go-springAi/internal/logger"
	"go-springAi/internal/middleware"
	"go-springAi/internal/provider"
	"go-springAi/internal/response"
	"go-springAi/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AIController 统一AI控制器
type AIController struct {
	BaseController
	providerManager *provider.Manager
	apiKeyService   service.APIKeyService
	logger          *zap.Logger
}

// NewAIController 创建统一AI控制器
func NewAIController(providerManager *provider.Manager, apiKeyService service.APIKeyService, logger *zap.Logger) *AIController {
	return &AIController{
		BaseController:  *NewBaseController(),
		providerManager: providerManager,
		apiKeyService:   apiKeyService,
		logger:          logger,
	}
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

// ListAllModels 列出指定提供商的所有模型（包括禁用的，用于模型管理）
func (ac *AIController) ListAllModels(c *gin.Context) {
	providerType := c.Param("provider")

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("ai"),
		logger.Operation("list_all_models"),
		logger.String("provider", providerType))

	// 获取Provider
	prov, err := ac.providerManager.GetProvider(provider.ProviderType(providerType))
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("list_all_models"),
			logger.String("provider", providerType),
			logger.ZapError(err))
		response.Error(c, http.StatusBadRequest, "Invalid provider", err.Error())
		return
	}

	models, err := prov.ListAllModels(c.Request.Context())
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("list_all_models"),
			logger.String("provider", providerType),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "All models retrieved successfully", gin.H{
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

	// 获取用户ID，如果没有认证则使用默认用户ID
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		// 在没有认证的情况下，使用默认用户ID 1
		userID = 1
		logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("set_api_key"),
			logger.String("provider", providerType),
			logger.String("user_id", "1"),
			logger.String("auth_status", "unauthenticated"))
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("ai"),
		logger.Operation("set_api_key"),
		logger.String("provider", providerType),
		logger.String("user_id", strconv.FormatInt(userID, 10)))

	var req dto.SetAPIKeyRequest
	if err := ac.BindAndValidate(c, &req); err != nil {
		logger.WarnCtx(c.Request.Context(), logger.MsgAPIValidation,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("set_api_key"),
			logger.String("provider", providerType),
			logger.String("user_id", strconv.FormatInt(userID, 10)),
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
			logger.String("user_id", strconv.FormatInt(userID, 10)),
			logger.ZapError(err))
		response.Error(c, http.StatusBadRequest, "Invalid provider", err.Error())
		return
	}

	// 保存API密钥到数据库
	err = ac.apiKeyService.SetAPIKey(c.Request.Context(), userID, providerType, req.APIKey)
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("set_api_key"),
			logger.String("provider", providerType),
			logger.String("user_id", strconv.FormatInt(userID, 10)),
			logger.ZapError(err))
		response.Error(c, http.StatusInternalServerError, "Failed to save API key", err.Error())
		return
	}

	// 设置Provider的API密钥（用于当前会话）
	err = prov.SetAPIKey(req.APIKey)
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("set_api_key"),
			logger.String("provider", providerType),
			logger.String("user_id", strconv.FormatInt(userID, 10)),
			logger.ZapError(err))
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "API key set successfully", gin.H{
		"provider": providerType,
	})
}

// APIKeyInfo API密钥信息结构体
type APIKeyInfo struct {
	HasKey    bool   `json:"has_key"`
	MaskedKey string `json:"masked_key,omitempty"`
}

// GetAPIKeyStatus 获取用户的API密钥状态
func (ac *AIController) GetAPIKeyStatus(c *gin.Context) {
	// 获取用户ID，如果没有认证则使用默认用户ID
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		// 在没有认证的情况下，使用默认用户ID 1
		userID = 1
		logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("get_api_key_status"),
			logger.String("user_id", "1"),
			logger.String("auth_status", "unauthenticated"))
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("ai"),
		logger.Operation("get_api_key_status"),
		logger.String("user_id", strconv.FormatInt(userID, 10)))

	// 获取所有提供商的API密钥状态
	apiKeyStatus := make(map[string]APIKeyInfo)
	
	// 获取所有支持的提供商类型
	supportedProviders := []string{"openai", "googleai", "mock"}
	
	for _, providerType := range supportedProviders {
		hasKey, err := ac.apiKeyService.CheckAPIKeyExists(c.Request.Context(), userID, providerType)
		if err != nil {
			logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
				logger.Module(logger.ModuleController),
				logger.Component("ai"),
				logger.Operation("get_api_key_status"),
				logger.String("user_id", strconv.FormatInt(userID, 10)),
				logger.String("provider", providerType),
				logger.ZapError(err))
			// 如果检查失败，默认为false
			apiKeyStatus[providerType] = APIKeyInfo{HasKey: false}
			continue
		}
		
		keyInfo := APIKeyInfo{HasKey: hasKey}
		
		// 如果有密钥，获取脱敏的密钥信息
		if hasKey {
			maskedKey, err := ac.apiKeyService.GetMaskedAPIKey(c.Request.Context(), userID, providerType)
			if err != nil {
				logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
					logger.Module(logger.ModuleController),
					logger.Component("ai"),
					logger.Operation("get_masked_api_key"),
					logger.String("user_id", strconv.FormatInt(userID, 10)),
					logger.String("provider", providerType),
					logger.ZapError(err))
				// 如果获取脱敏密钥失败，仍然标记为有密钥，但不显示脱敏信息
				keyInfo.MaskedKey = "****"
			} else {
				keyInfo.MaskedKey = maskedKey
			}
		}
		
		apiKeyStatus[providerType] = keyInfo
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("ai"),
		logger.Operation("get_api_key_status"),
		logger.String("user_id", strconv.FormatInt(userID, 10)))

	response.Success(c, http.StatusOK, "API key status retrieved successfully", apiKeyStatus)
}

// GetPlainAPIKey 获取明文API密钥
func (ac *AIController) GetPlainAPIKey(c *gin.Context) {
	providerType := c.Param("provider")
	
	// 获取用户ID，如果没有认证则使用默认用户ID
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		// 在没有认证的情况下，使用默认用户ID 1
		userID = 1
		logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("get_plain_api_key"),
			logger.String("provider", providerType),
			logger.String("user_id", "1"),
			logger.String("auth_status", "unauthenticated"))
	}

	logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
		logger.Module(logger.ModuleController),
		logger.Component("ai"),
		logger.Operation("get_plain_api_key"),
		logger.String("provider", providerType),
		logger.String("user_id", strconv.FormatInt(userID, 10)))

	// 验证提供商类型
	if !ac.isValidProviderType(providerType) {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("get_plain_api_key"),
			logger.String("provider", providerType),
			logger.String("user_id", strconv.FormatInt(userID, 10)),
			logger.String("error", "invalid provider type"))
		response.Error(c, http.StatusBadRequest, "Invalid provider type", "Unsupported provider: "+providerType)
		return
	}

	// 检查API密钥是否存在
	hasKey, err := ac.apiKeyService.CheckAPIKeyExists(c.Request.Context(), userID, providerType)
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("get_plain_api_key"),
			logger.String("provider", providerType),
			logger.String("user_id", strconv.FormatInt(userID, 10)),
			logger.ZapError(err))
		response.Error(c, http.StatusInternalServerError, "Failed to check API key", err.Error())
		return
	}

	if !hasKey {
		logger.InfoCtx(c.Request.Context(), logger.MsgAPIRequest,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("get_plain_api_key"),
			logger.String("provider", providerType),
			logger.String("user_id", strconv.FormatInt(userID, 10)),
			logger.String("result", "no_key_found"))
		response.Error(c, http.StatusNotFound, "API key not found", "No API key configured for this provider")
		return
	}

	// 获取明文API密钥
	plainKey, err := ac.apiKeyService.GetAPIKey(c.Request.Context(), userID, providerType)
	if err != nil {
		logger.ErrorCtx(c.Request.Context(), logger.MsgAPIError,
			logger.Module(logger.ModuleController),
			logger.Component("ai"),
			logger.Operation("get_plain_api_key"),
			logger.String("provider", providerType),
			logger.String("user_id", strconv.FormatInt(userID, 10)),
			logger.ZapError(err))
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve API key", err.Error())
		return
	}

	// 记录明文密钥访问日志（安全审计）
	logger.InfoCtx(c.Request.Context(), "API key accessed in plain text",
		logger.Module(logger.ModuleController),
		logger.Component("ai"),
		logger.Operation("get_plain_api_key"),
		logger.String("provider", providerType),
		logger.String("user_id", strconv.FormatInt(userID, 10)),
		logger.String("access_type", "plain_text"))

	response.Success(c, http.StatusOK, "Plain API key retrieved successfully", gin.H{
		"provider": providerType,
		"api_key":  plainKey,
	})
}

// isValidProviderType 验证提供商类型是否有效
func (ac *AIController) isValidProviderType(providerType string) bool {
	validProviders := []string{"openai", "googleai", "mock"}
	for _, valid := range validProviders {
		if providerType == valid {
			return true
		}
	}
	return false
}
