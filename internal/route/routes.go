package route

import (
	"admin/internal/controllers"
	"admin/internal/dto"
	"admin/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetupRoutes 设置路由
func SetupRoutes(logger *zap.Logger, mcpController *controllers.MCPController) *gin.Engine {
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
	}

	return r
}
