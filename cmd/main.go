package main

import (
	"context"
	"fmt"

	"go-springAi/internal/dto"
	"go-springAi/internal/logger"
	"go-springAi/internal/wire"

	"github.com/gin-gonic/gin"
)

func main() {
	// 使用wire初始化应用
	app, cleanup, err := wire.InitializeApp(".")
	if err != nil {
		logger.Fatal(logger.MsgServerError,
			logger.ZapError(err),
			logger.Module(logger.ModuleServer),
			logger.Operation(logger.OpStart))
	}
	defer cleanup()

	// 自动初始化MCP系统
	initializeMCPSystem(app)

	// 设置Gin模式
	gin.SetMode(app.Config.Server.Mode)

	// 设置路由
	router := app.Router

	// 启动服务器
	addr := fmt.Sprintf("%s:%s", app.Config.Server.Host, app.Config.Server.Port)
	
	// 使用统一日志记录服务器启动
	logger.Info(logger.MsgServerStarting,
		logger.String("address", addr),
		logger.String("mode", app.Config.Server.Mode),
		logger.Module(logger.ModuleServer),
		logger.Operation(logger.OpStart))

	if err := router.Run(addr); err != nil {
		logger.Fatal(logger.MsgServerError,
			logger.ZapError(err),
			logger.String("address", addr),
			logger.Module(logger.ModuleServer),
			logger.Operation(logger.OpStart))
	}
}

// initializeMCPSystem 自动初始化MCP系统
func initializeMCPSystem(app *wire.App) {
	if app.MCPService == nil {
		logger.Warn(logger.MsgServerError,
			logger.String("message", "MCP service is not available, skipping auto-initialization"),
			logger.Module(logger.ModuleServer),
			logger.Operation(logger.OpStart))
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
		logger.Warn(logger.MsgServerError,
			logger.ZapError(err),
			logger.String("message", "Failed to auto-initialize MCP system"),
			logger.Module(logger.ModuleServer),
			logger.Operation(logger.OpStart))
		return
	}

	logger.Info(logger.MsgServerStarting,
		logger.String("protocolVersion", response.ProtocolVersion),
		logger.String("serverName", response.ServerInfo.Name),
		logger.String("serverVersion", response.ServerInfo.Version),
		logger.String("message", "MCP system auto-initialized successfully"),
		logger.Module(logger.ModuleServer),
		logger.Operation(logger.OpStart))
}
