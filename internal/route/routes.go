package route

import (
	"admin/internal/controllers"
	"admin/internal/dto"
	"admin/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupRoutes 设置路由
func SetupRoutes(logger *zap.Logger, mcpController *controllers.MCPController, openaiController *controllers.OpenAIController, googleaiController *controllers.GoogleAIController, aiController *controllers.AIController) *gin.Engine {
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

		// OpenAI 相关路由
		openaiGroup := v1.Group("/openai")
		{
			openaiGroup.POST("/chat/completions", openaiController.ChatCompletion)
			openaiGroup.GET("/models", openaiController.ListModels)
			openaiGroup.POST("/validate", openaiController.ValidateAPIKey)
			openaiGroup.POST("/api-key", openaiController.SetAPIKey)
			openaiGroup.GET("/config/:model", openaiController.GetModelConfig)
			openaiGroup.PUT("/models/:model/enable", openaiController.EnableModel)
			openaiGroup.PUT("/models/:model/disable", openaiController.DisableModel)
		}

		// Google AI 相关路由
		googleaiGroup := v1.Group("/googleai")
		{
			googleaiGroup.POST("/chat/completions", googleaiController.ChatCompletion)
			googleaiGroup.GET("/models", googleaiController.ListModels)
			googleaiGroup.POST("/validate", googleaiController.ValidateAPIKey)
			googleaiGroup.POST("/api-key", googleaiController.SetAPIKey)
			googleaiGroup.GET("/config/:model", googleaiController.GetModelConfig)
			googleaiGroup.PUT("/models/:model/enable", googleaiController.EnableModel)
			googleaiGroup.PUT("/models/:model/disable", googleaiController.DisableModel)
		}

		// 统一AI API端点
		aiGroup := v1.Group("/ai")
		{
			// 聊天完成端点 - 支持所有提供商
			aiGroup.POST("/:provider/chat/completions", aiController.ChatCompletion)
			
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
	}

	return r
}
