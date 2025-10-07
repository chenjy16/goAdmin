//go:build wireinject
// +build wireinject

package wire

import (
	"admin/internal/config"
	"admin/internal/controllers"
	"admin/internal/database"
	"admin/internal/repository"
	"admin/internal/services"
	"admin/internal/utils"

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
		database.NewConnection,

		// JWT管理器
		ProvideJWTManager,

		// 验证器
		utils.NewCustomValidator,

		// Repository
		repository.NewRepositoryManager,

		// Services
		ProvideUserService,

		// Controllers
		controllers.NewUserController,

		// Router
		ProvideRouter,

		// App
		NewApp,
	)
	return &App{}, nil, nil
}

// App 应用程序结构
type App struct {
	Config         *config.Config
	Logger         *zap.Logger
	DB             *database.DB
	JWTManager     *utils.JWTManager
	Validator      *utils.CustomValidator
	RepoManager    repository.RepositoryManager
	UserService    *services.UserService
	UserController *controllers.UserController
	Router         *gin.Engine
}

// NewApp 创建应用程序实例
func NewApp(
	config *config.Config,
	logger *zap.Logger,
	db *database.DB,
	jwtManager *utils.JWTManager,
	validator *utils.CustomValidator,
	repoManager repository.RepositoryManager,
	userService *services.UserService,
	userController *controllers.UserController,
	router *gin.Engine,
) (*App, func()) {
	app := &App{
		Config:         config,
		Logger:         logger,
		DB:             db,
		JWTManager:     jwtManager,
		Validator:      validator,
		RepoManager:    repoManager,
		UserService:    userService,
		UserController: userController,
		Router:         router,
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
