package wire

import (
	"admin/internal/config"
	"admin/internal/controllers"
	"admin/internal/database"
	"admin/internal/logger"
	"admin/internal/repository"
	"admin/internal/route"
	"admin/internal/service"
	"admin/internal/utils"

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

// ProvideRouter 提供路由器
func ProvideRouter(mcpController *controllers.MCPController, logger *zap.Logger) *gin.Engine {
	return route.SetupRoutes(logger, mcpController)
}
