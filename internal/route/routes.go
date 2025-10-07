package route

import (
	"admin/internal/controllers"
	"admin/internal/dto"
	"admin/internal/middleware"
	"admin/internal/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupRoutes 设置路由
func SetupRoutes(userService *services.UserService, logger *zap.Logger) *gin.Engine {
	// 创建Gin引擎
	r := gin.New()

	// 添加中间件
	r.Use(middleware.RequestID())          // 请求ID中间件
	r.Use(middleware.ZapLogger(logger))    // zap结构化日志中间件
	r.Use(middleware.ErrorHandler(logger)) // 错误处理中间件
	r.Use(middleware.Recovery())           // 恢复中间件
	r.Use(middleware.CORS())               // 跨域中间件

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// API版本分组
	v1 := r.Group("/api/v1")
	{
		// 用户相关路由
		userController := controllers.NewUserController(userService)
		users := v1.Group("/users")
		{
			users.POST("", middleware.ValidateJSONFactory(&dto.CreateUserRequest{}), userController.CreateUser)
			users.GET("", userController.ListUsers)
			users.GET("/:id", userController.GetUser)
			users.PUT("/:id", middleware.ValidateJSONFactory(&dto.UpdateUserRequest{}), userController.UpdateUser)
			users.DELETE("/:id", userController.DeleteUser)
		}
	}

	return r
}
