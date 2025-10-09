//go:build wireinject
// +build wireinject

package wire

import (
	"goMcp/internal/config"
	"goMcp/internal/controllers"
	"goMcp/internal/database"
	"goMcp/internal/provider"
	"goMcp/internal/repository"
	"goMcp/internal/service"
	"goMcp/internal/utils"

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

		// 验证器
		utils.NewCustomValidator,

		// Repository
		repository.NewRepositoryManager,

		// Services
		ProvideMCPService,
		ProvideOpenAIService,
		ProvideGoogleAIService,
		ProvideAIAssistantService,

		// Controllers
		ProvideMCPController,
		ProvideOpenAIController,
		ProvideGoogleAIController,
		ProvideAIAssistantController,

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
	Validator              *utils.CustomValidator
	RepoManager            repository.RepositoryManager
	MCPService             service.MCPService
	OpenAIService          *service.OpenAIService
	GoogleAIService        *service.GoogleAIService
	AIAssistantService     *service.AIAssistantService
	MCPController          *controllers.MCPController
	OpenAIController       *controllers.OpenAIController
	GoogleAIController     *controllers.GoogleAIController
	AIAssistantController  *controllers.AIAssistantController
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
	validator *utils.CustomValidator,
	repoManager repository.RepositoryManager,
	mcpService service.MCPService,
	openaiService *service.OpenAIService,
	googleaiService *service.GoogleAIService,
	aiAssistantService *service.AIAssistantService,
	mcpController *controllers.MCPController,
	openaiController *controllers.OpenAIController,
	googleaiController *controllers.GoogleAIController,
	aiAssistantController *controllers.AIAssistantController,
	providerManager *provider.Manager,
	aiController *controllers.AIController,
	router *gin.Engine,
) (*App, func()) {
	app := &App{
		Config:                config,
		Logger:                logger,
		DB:                    db,
		JWTManager:            jwtManager,
		Validator:             validator,
		RepoManager:           repoManager,
		MCPService:            mcpService,
		OpenAIService:         openaiService,
		GoogleAIService:       googleaiService,
		AIAssistantService:    aiAssistantService,
		MCPController:         mcpController,
		OpenAIController:      openaiController,
		GoogleAIController:    googleaiController,
		AIAssistantController: aiAssistantController,
		ProviderManager:       providerManager,
		AIController:          aiController,
		Router:                router,
	}

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
