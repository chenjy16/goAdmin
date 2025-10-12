//go:build wireinject
// +build wireinject

package wire

import (
	"context"
	
	"go-springAi/internal/config"
	"go-springAi/internal/controllers"
	"go-springAi/internal/database"
	"go-springAi/internal/dto"
	"go-springAi/internal/errors"

	"go-springAi/internal/i18n"
	"go-springAi/internal/provider"
	"go-springAi/internal/repository"
	"go-springAi/internal/service"
	"go-springAi/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.uber.org/zap"
)

// InitializeApp 初始化应用程序
func InitializeApp(configPath string) (*App, func(), error) {
	wire.Build(
		// 配置
		ProvideConfig,

		// 日志
		ProvideLogger,

		// 数据库
		ProvideDatabase,

		// JWT管理器
		ProvideJWTManager,

		// 国际化管理器
		ProvideI18nManager,

		// 错误处理器
		ProvideErrorHandler,

		// 验证器
		utils.NewCustomValidator,

		// Repository
		repository.NewRepositoryManager,

		// Services
		ProvideMCPService,
		ProvideInternalMCPClient,
		ProvideOpenAIService,
		ProvideGoogleAIService,
		ProvideAPIKeyService,
		ProvideStockAnalysisService,
		ProvideAIAssistantService,

		// Controllers
		ProvideMCPController,
		ProvideAIAssistantController,
		ProvideTestI18nController,
		ProvideStockController,

		// Provider Manager
		ProvideProviderManager,

		// AI Controller
		ProvideAIController,

		// Router
		ProvideRouter,

		// App
		NewApp,
	)
	return &App{}, nil, nil
}

// App 应用程序结构
type App struct {
	Config                 *config.Config
	Logger                 *zap.Logger
	DB                     *database.DB
	JWTManager             *utils.JWTManager
	I18nManager            *i18n.Manager
	ErrorHandler           *errors.ErrorHandler
	Validator              *utils.CustomValidator
	RepoManager            repository.RepositoryManager
	MCPService             service.MCPService
	OpenAIService          *service.OpenAIService
	GoogleAIService        *service.GoogleAIService
	APIKeyService          service.APIKeyService
	StockAnalysisService   *service.StockAnalysisService
	AIAssistantService     *service.AIAssistantService
	MCPController          *controllers.MCPController
	AIAssistantController  *controllers.AIAssistantController
	TestI18nController     *controllers.TestI18nController
	StockController        *controllers.StockController
	ProviderManager        *provider.Manager
	AIController           *controllers.AIController
	Router                 *gin.Engine
}

// NewApp 创建应用程序实例
func NewApp(
	config *config.Config,
	logger *zap.Logger,
	db *database.DB,
	jwtManager *utils.JWTManager,
	i18nManager *i18n.Manager,
	errorHandler *errors.ErrorHandler,
	validator *utils.CustomValidator,
	repoManager repository.RepositoryManager,
	mcpService service.MCPService,
	openaiService *service.OpenAIService,
	googleaiService *service.GoogleAIService,
	apiKeyService service.APIKeyService,
	stockAnalysisService *service.StockAnalysisService,
	aiAssistantService *service.AIAssistantService,
	mcpController *controllers.MCPController,
	aiAssistantController *controllers.AIAssistantController,
	testI18nController *controllers.TestI18nController,
	stockController *controllers.StockController,
	providerManager *provider.Manager,
	aiController *controllers.AIController,
	router *gin.Engine,
) (*App, func()) {
	app := &App{
		Config:                config,
		Logger:                logger,
		DB:                    db,
		JWTManager:            jwtManager,
		I18nManager:           i18nManager,
		ErrorHandler:          errorHandler,
		Validator:             validator,
		RepoManager:           repoManager,
		MCPService:            mcpService,
		OpenAIService:         openaiService,
		GoogleAIService:       googleaiService,
		APIKeyService:         apiKeyService,
		StockAnalysisService:  stockAnalysisService,
		AIAssistantService:    aiAssistantService,
		MCPController:         mcpController,
		AIAssistantController: aiAssistantController,
		TestI18nController:    testI18nController,
		StockController:       stockController,
		ProviderManager:       providerManager,
		AIController:          aiController,
		Router:                router,
	}

	// 自动初始化MCP系统
	app.initializeMCPSystem()

	// 清理函数
	cleanup := func() {
		if app.DB != nil {
			app.DB.Close()
		}
		if app.Logger != nil {
			app.Logger.Sync()
		}
	}

	return app, cleanup
}

// initializeMCPSystem 自动初始化MCP系统
func (app *App) initializeMCPSystem() {
	if app.MCPService == nil {
		app.Logger.Warn("MCP service is not available, skipping auto-initialization")
		return
	}

	// 创建初始化请求
	initReq := &dto.MCPInitializeRequest{
		ProtocolVersion: "2024-11-05",
		Capabilities: dto.MCPCapabilities{
			Tools: &dto.MCPToolsCapability{
				ListChanged: true,
			},
			Logging: &dto.MCPLoggingCapability{},
		},
		ClientInfo: dto.MCPClientInfo{
			Name:    "Auto-initialized MCP Server",
			Version: "1.0.0",
		},
	}

	// 使用context.Background()进行初始化
	ctx := context.Background()
	
	// 执行初始化
	response, err := app.MCPService.Initialize(ctx, initReq)
	if err != nil {
		app.Logger.Error("Failed to auto-initialize MCP system",
			zap.Error(err),
			zap.String("module", "startup"),
			zap.String("operation", "mcp_auto_init"))
		return
	}

	app.Logger.Info("MCP system auto-initialized successfully",
		zap.String("protocolVersion", response.ProtocolVersion),
		zap.String("serverName", response.ServerInfo.Name),
		zap.String("serverVersion", response.ServerInfo.Version),
		zap.String("module", "startup"),
		zap.String("operation", "mcp_auto_init"))
}
