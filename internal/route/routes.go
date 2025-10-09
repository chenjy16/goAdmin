package route

import (
	"go-springAi/internal/controllers"
	"go-springAi/internal/dto"
	"go-springAi/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupRoutes 设置路由
func SetupRoutes(logger *zap.Logger, mcpController *controllers.MCPController, aiController *controllers.AIController, aiAssistantController *controllers.AIAssistantController) *gin.Engine {
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

		// MCP相关路由
		mcp := v1.Group("/mcp")
		{
			// MCP初始化端点
			mcp.POST("/initialize", middleware.ValidateJSONFactory(&dto.MCPInitializeRequest{}), mcpController.Initialize)
			
			// 工具管理端点
			mcp.GET("/tools", mcpController.ListTools)
			mcp.POST("/execute", middleware.ValidateJSONFactory(&dto.MCPExecuteRequest{}), mcpController.ExecuteTool)
			
			// SSE流式端点
			mcp.GET("/sse", mcpController.StreamSSE)
			
			// 执行日志端点
			mcp.GET("/logs", mcpController.ListExecutionLogs)
			mcp.GET("/logs/:id", mcpController.GetExecutionLog)
		}



		// 统一AI API端点
		aiGroup := v1.Group("/ai")
		{
			// 模型管理端点
			aiGroup.GET("/:provider/models", aiController.ListModels)
			aiGroup.GET("/:provider/config/:model", aiController.GetModelConfig)
			aiGroup.PUT("/:provider/models/:model/enable", aiController.EnableModel)
			aiGroup.PUT("/:provider/models/:model/disable", aiController.DisableModel)
			
			// API密钥管理端点
			aiGroup.POST("/:provider/api-key", aiController.SetAPIKey)
			aiGroup.POST("/:provider/validate", aiController.ValidateAPIKey)
			
			// 提供商管理端点
			aiGroup.GET("/providers", aiController.ListProviders)
		}

		// AI助手端点
		assistantGroup := v1.Group("/assistant")
		{
			// 初始化AI助手
			assistantGroup.POST("/initialize", aiAssistantController.Initialize)
			
			// AI助手聊天端点
			assistantGroup.POST("/chat", aiAssistantController.Chat)
		}
	}

	return r
}
