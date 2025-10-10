package wire

import (
	"context"
	"time"

	"go-springAi/internal/config"
	"go-springAi/internal/controllers"
	"go-springAi/internal/database"
	"go-springAi/internal/googleai"
	"go-springAi/internal/logger"
	"go-springAi/internal/openai"
	"go-springAi/internal/provider"
	"go-springAi/internal/repository"
	"go-springAi/internal/route"
	"go-springAi/internal/service"
	"go-springAi/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ProvideConfig 提供配置
func ProvideConfig(configPath string) (*config.Config, error) {
	return config.LoadConfig(configPath)
}

// ProvideLogger 提供日志器
func ProvideLogger(cfg *config.Config) (*zap.Logger, error) {
	// 初始化全局日志器
	if err := logger.InitGlobalLogger(cfg.Server.Mode); err != nil {
		return nil, err
	}

	// 记录配置加载成功
	logger.Info(logger.MsgConfigLoaded,
		logger.Module(logger.ModuleConfig),
		logger.String("mode", cfg.Server.Mode))

	// 为了兼容现有代码，仍然返回zap.Logger
	var zapLogger *zap.Logger
	var err error
	if cfg.Server.Mode == "release" {
		zapLogger, err = zap.NewProduction()
	} else {
		zapLogger, err = zap.NewDevelopment()
	}

	if err != nil {
		return nil, err
	}

	// 设置全局日志器
	zap.ReplaceGlobals(zapLogger)

	return zapLogger, nil
}

// ProvideDatabase 提供数据库连接
func ProvideDatabase(cfg *config.Config) (*database.DB, error) {
	return database.NewConnection(cfg.Database.Driver, cfg.Database.DSN)
}

// ProvideJWTManager 提供JWT管理器
func ProvideJWTManager(cfg *config.Config) *utils.JWTManager {
	return utils.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpireTime)
}

// ProvideMCPService 提供MCP服务
func ProvideMCPService(repoManager repository.RepositoryManager, logger *zap.Logger) service.MCPService {
	userService := service.NewUserServiceAdapter(repoManager)
	return service.NewMCPService(userService, logger)
}

// ProvideMCPController 提供MCP控制器
func ProvideMCPController(mcpService service.MCPService, logger *zap.Logger) *controllers.MCPController {
	return controllers.NewMCPController(mcpService, logger)
}

// ProvideOpenAIService 提供OpenAI服务
func ProvideOpenAIService(cfg *config.Config, zapLogger *zap.Logger) *service.OpenAIService {
	// 创建OpenAI配置
	openaiConfig := &openai.Config{
		APIKey:       cfg.OpenAI.APIKey,
		BaseURL:      cfg.OpenAI.BaseURL,
		Timeout:      time.Duration(cfg.OpenAI.Timeout) * time.Second,
		MaxRetries:   cfg.OpenAI.MaxRetries,
		DefaultModel: cfg.OpenAI.DefaultModel,
	}

	// 创建内存管理器
	keyManager := openai.NewMemoryKeyManager()
	modelManager := openai.NewMemoryModelManager()

	// 将配置中的API密钥设置到密钥管理器中
	if cfg.OpenAI.APIKey != "" {
		if err := keyManager.SetAPIKey(cfg.OpenAI.APIKey); err != nil {
			logger.LogError("Failed to set OpenAI API key",
				logger.Module(logger.ModuleService),
				logger.String("error", err.Error()))
		}
	}

	// 创建HTTP客户端，传入密钥管理器
	httpClient := openai.NewHTTPClient(openaiConfig, keyManager)

	// 使用全局日志器
	globalLogger := logger.GetGlobalLogger()

	return service.NewOpenAIService(httpClient, keyManager, modelManager, globalLogger)
}



// ProvideGoogleAIService 提供Google AI服务
func ProvideGoogleAIService(cfg *config.Config, zapLogger *zap.Logger) (*service.GoogleAIService, error) {
	// 创建Google AI配置
	googleaiConfig := &googleai.Config{
		APIKey:       cfg.GoogleAI.APIKey,
		ProjectID:    cfg.GoogleAI.ProjectID,
		Location:     cfg.GoogleAI.Location,
		Timeout:      time.Duration(cfg.GoogleAI.Timeout) * time.Second,
		MaxRetries:   cfg.GoogleAI.MaxRetries,
		DefaultModel: cfg.GoogleAI.DefaultModel,
	}

	// 创建内存管理器
	keyManager := googleai.NewKeyManager(cfg.GoogleAI.APIKey)
	modelManager := googleai.NewModelManager()

	// 创建HTTP客户端，传递keyManager
	httpClient, err := googleai.NewHTTPClient(googleaiConfig, keyManager)
	if err != nil {
		return nil, err
	}

	// 使用全局日志器
	globalLogger := logger.GetGlobalLogger()

	return service.NewGoogleAIService(httpClient, keyManager, modelManager, globalLogger), nil
}



// ProvideProviderManager 提供Provider管理器
func ProvideProviderManager(openaiService *service.OpenAIService, googleaiService *service.GoogleAIService, zapLogger *zap.Logger) *provider.Manager {
	// 使用全局日志器
	globalLogger := logger.GetGlobalLogger()
	manager := provider.NewManager(globalLogger)
	
	// 创建并注册OpenAI Provider
	openaiProvider := provider.NewOpenAIProvider(openaiService)
	manager.RegisterProvider(openaiProvider)
	
	// 创建并注册Google AI Provider
	googleaiProvider := provider.NewGoogleAIProvider(googleaiService)
	manager.RegisterProvider(googleaiProvider)
	
	// 创建并注册Mock Provider（用于测试）
	mockProvider := provider.NewMockProvider("mock", provider.ProviderTypeMock)
	manager.RegisterProvider(mockProvider)
	
	return manager
}

// ProvideAPIKeyService 提供API密钥服务
func ProvideAPIKeyService(repoManager repository.RepositoryManager) service.APIKeyService {
	return service.NewAPIKeyService(repoManager.APIKey())
}

// ProvideAIController 提供AI控制器
func ProvideAIController(providerManager *provider.Manager, apiKeyService service.APIKeyService, logger *zap.Logger) *controllers.AIController {
	return controllers.NewAIController(providerManager, apiKeyService, logger)
}

// ProvideAIAssistantService 提供AI助手服务
func ProvideAIAssistantService(mcpService service.MCPService, openaiService *service.OpenAIService, providerManager *provider.Manager, logger *zap.Logger) *service.AIAssistantService {
	// 创建适配器来实现接口
	adapter := &ProviderManagerAdapter{manager: providerManager}
	return service.NewAIAssistantService(mcpService, openaiService, adapter, logger)
}

// ProviderManagerAdapter 适配器，将provider.Manager适配为service.ProviderManager接口
type ProviderManagerAdapter struct {
	manager *provider.Manager
}

func (a *ProviderManagerAdapter) GetProviderByModel(modelName string) (service.ProviderInterface, error) {
	provider, err := a.manager.GetProviderByModel(modelName)
	if err != nil {
		return nil, err
	}
	
	// 创建Provider适配器
	return &ProviderAdapter{provider: provider}, nil
}

func (a *ProviderManagerAdapter) GetProviderByName(name string) (service.ProviderInterface, error) {
	provider, err := a.manager.GetProviderByName(name)
	if err != nil {
		return nil, err
	}
	
	// 创建Provider适配器
	return &ProviderAdapter{provider: provider}, nil
}

func (a *ProviderManagerAdapter) ValidateModelForProvider(ctx context.Context, providerName, modelName string) error {
	return a.manager.ValidateModelForProvider(ctx, providerName, modelName)
}

func (a *ProviderManagerAdapter) GetProviderByModelWithValidation(ctx context.Context, modelName string) (service.ProviderInterface, error) {
	provider, err := a.manager.GetProviderByModelWithValidation(ctx, modelName)
	if err != nil {
		return nil, err
	}
	
	// 创建Provider适配器
	return &ProviderAdapter{provider: provider}, nil
}

// ProviderAdapter 适配器，将provider.Provider适配为service.ProviderInterface接口
type ProviderAdapter struct {
	provider provider.Provider
}

func (a *ProviderAdapter) GetType() string {
	return string(a.provider.GetType())
}

func (a *ProviderAdapter) GetName() string {
	return a.provider.GetName()
}

func (a *ProviderAdapter) ChatCompletion(ctx context.Context, request *service.ProviderChatRequest) (*service.ProviderChatResponse, error) {
	// 转换请求格式
	providerMessages := make([]provider.Message, len(request.Messages))
	for i, msg := range request.Messages {
		providerMessages[i] = provider.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	
	providerReq := &provider.ChatRequest{
		Model:       request.Model,
		Messages:    providerMessages,
		MaxTokens:   request.MaxTokens,
		Temperature: request.Temperature,
		TopP:        request.TopP,
		TopK:        request.TopK,
		Stream:      request.Stream,
		Options:     request.Options,
	}
	
	// 调用实际的provider
	resp, err := a.provider.ChatCompletion(ctx, providerReq)
	if err != nil {
		return nil, err
	}
	
	// 转换响应格式
	serviceChoices := make([]service.ProviderChoice, len(resp.Choices))
	for i, choice := range resp.Choices {
		serviceChoices[i] = service.ProviderChoice{
			Index: choice.Index,
			Message: service.ProviderMessage{
				Role:    choice.Message.Role,
				Content: choice.Message.Content,
			},
			FinishReason: choice.FinishReason,
		}
	}
	
	return &service.ProviderChatResponse{
		ID:      resp.ID,
		Object:  resp.Object,
		Created: resp.Created,
		Model:   resp.Model,
		Choices: serviceChoices,
		Usage: service.ProviderUsage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
	}, nil
}

// ProvideAIAssistantController 提供AI助手控制器
func ProvideAIAssistantController(aiAssistantService *service.AIAssistantService, logger *zap.Logger) *controllers.AIAssistantController {
	return controllers.NewAIAssistantController(aiAssistantService, logger)
}

// ProvideRouter 提供路由器
func ProvideRouter(jwtManager *utils.JWTManager, mcpController *controllers.MCPController, aiController *controllers.AIController, aiAssistantController *controllers.AIAssistantController, logger *zap.Logger) *gin.Engine {
	return route.SetupRoutes(logger, jwtManager, mcpController, aiController, aiAssistantController)
}
