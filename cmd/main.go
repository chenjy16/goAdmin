package main

import (
	"fmt"

	"goMcp/internal/logger"
	"goMcp/internal/wire"

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
