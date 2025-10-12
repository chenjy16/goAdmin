package route

import (
	"go-springAi/internal/controllers"
	"go-springAi/internal/dto"

	"go-springAi/internal/i18n"
	"go-springAi/internal/middleware"
	"go-springAi/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupRoutes 设置路由
func SetupRoutes(logger *zap.Logger, jwtManager *utils.JWTManager, mcpController *controllers.MCPController, aiController *controllers.AIController, aiAssistantController *controllers.AIAssistantController, stockController *controllers.StockController, testI18nController *controllers.TestI18nController, i18nManager *i18n.Manager) *gin.Engine {
	// 创建Gin引擎
	r := gin.New()

	// 添加中间件
	r.Use(middleware.RequestID())          // 请求ID中间件
	r.Use(middleware.ZapLogger(logger))    // zap结构化日志中间件
	r.Use(middleware.ErrorHandler(logger)) // 错误处理中间件
	r.Use(middleware.Recovery())           // 恢复中间件
	r.Use(middleware.CORS())               // 跨域中间件
	r.Use(middleware.I18nMiddleware(i18nManager)) // 国际化中间件

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
			
			// MCP状态端点
			mcp.GET("/status", mcpController.GetStatus)
			
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
			aiGroup.GET("/:provider/models/all", aiController.ListAllModels) // 新增：获取所有模型（包括禁用的）
			aiGroup.GET("/:provider/config/:model", aiController.GetModelConfig)
			aiGroup.PUT("/:provider/models/:model/enable", aiController.EnableModel)
			aiGroup.PUT("/:provider/models/:model/disable", aiController.DisableModel)
			
			// API密钥管理端点（可选认证）
			aiGroup.POST("/:provider/api-key", middleware.OptionalAuthMiddleware(jwtManager, logger), aiController.SetAPIKey)
			aiGroup.POST("/:provider/validate", middleware.OptionalAuthMiddleware(jwtManager, logger), aiController.ValidateAPIKey)
			aiGroup.GET("/api-keys/status", middleware.OptionalAuthMiddleware(jwtManager, logger), aiController.GetAPIKeyStatus)
			aiGroup.GET("/:provider/api-key/plain", middleware.OptionalAuthMiddleware(jwtManager, logger), aiController.GetPlainAPIKey)
			
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

		// 股票分析端点
		stockGroup := v1.Group("/stock")
		{
			// 股票分析
			stockGroup.POST("/analyze", stockController.AnalyzeStock)
			
			// 股票比较
			stockGroup.POST("/compare", stockController.CompareStocks)
			
			// 股票报价
			stockGroup.GET("/quote/:symbol", stockController.GetStockQuote)
			
			// 股票历史数据
			stockGroup.GET("/history/:symbol", stockController.GetStockHistory)
			
			// 市场摘要
			stockGroup.GET("/market/summary", stockController.GetMarketSummary)
		}

		// 国际化测试端点
		testGroup := v1.Group("/test")
		{
			// 测试成功响应
			testGroup.GET("/success", testI18nController.TestSuccess)
			
			// 测试错误响应
			testGroup.GET("/error", testI18nController.TestError)
			
			// 测试翻译功能
			testGroup.GET("/translation", testI18nController.TestTranslation)
		}
	}

	return r
}
