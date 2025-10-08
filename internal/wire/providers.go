package wire

import (
	"admin/internal/config"
	"admin/internal/database"
	"admin/internal/repository"
	"admin/internal/route"
	"admin/internal/services"
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
	var logger *zap.Logger
	var err error

	if cfg.Server.Mode == "release" {
		// 生产环境使用JSON格式日志
		logger, err = zap.NewProduction()
	} else {
		// 开发环境使用开发者友好的格式
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		return nil, err
	}

	// 设置全局日志器
	zap.ReplaceGlobals(logger)

	return logger, nil
}

// ProvideDatabase 提供数据库连接
func ProvideDatabase(cfg *config.Config) (*database.DB, error) {
	return database.NewConnection(cfg.Database.Driver, cfg.Database.DSN)
}

// ProvideJWTManager 提供JWT管理器
func ProvideJWTManager(cfg *config.Config) *utils.JWTManager {
	return utils.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpireTime)
}

// ProvideUserService 提供用户服务
func ProvideUserService(repoManager repository.RepositoryManager) *services.UserService {
	return services.NewUserService(repoManager.User())
}

// ProvideRouter 提供路由器
func ProvideRouter(userService *services.UserService, logger *zap.Logger) *gin.Engine {
	return route.SetupRoutes(userService, logger)
}
