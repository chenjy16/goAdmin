package main

import (
	"fmt"
	"log"

	"admin/internal/wire"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 使用wire初始化应用
	app, cleanup, err := wire.InitializeApp(".")
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer cleanup()

	// 设置Gin模式
	gin.SetMode(app.Config.Server.Mode)

	// 设置路由
	router := app.Router

	// 启动服务器
	addr := fmt.Sprintf("%s:%s", app.Config.Server.Host, app.Config.Server.Port)
	app.Logger.Info("Server starting",
		zap.String("address", addr),
		zap.String("mode", app.Config.Server.Mode))

	if err := router.Run(addr); err != nil {
		app.Logger.Fatal("Failed to start server", zap.Error(err))
	}
}
