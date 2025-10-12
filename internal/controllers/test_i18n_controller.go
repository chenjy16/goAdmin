package controllers

import (
	"go-springAi/internal/middleware"
	"go-springAi/internal/response"

	"github.com/gin-gonic/gin"
)

// TestI18nController 测试国际化控制器
type TestI18nController struct{}

// NewTestI18nController 创建测试国际化控制器
func NewTestI18nController() *TestI18nController {
	return &TestI18nController{}
}

// TestSuccess 测试成功响应的国际化
func (c *TestI18nController) TestSuccess(ctx *gin.Context) {
	// 使用国际化成功响应
	response.I18nSuccess(ctx, 200, "success", gin.H{
		"message": middleware.T(ctx, "test_success", nil),
		"data":    "Hello World",
	}, nil)
}

// TestError 测试错误响应的国际化
func (c *TestI18nController) TestError(ctx *gin.Context) {
	// 使用国际化错误响应
	response.I18nBadRequest(ctx, "validation_error", "Bad request example", nil)
}

// TestTranslation 测试翻译功能
func (c *TestI18nController) TestTranslation(ctx *gin.Context) {
	// 翻译消息
	message := middleware.T(ctx, "welcome_message", nil)
	errorMsg := middleware.T(ctx, "error_demo", nil)
	
	response.I18nSuccess(ctx, 200, "success", gin.H{
		"welcome":     message,
		"error_demo":  errorMsg,
		"language":    middleware.GetLanguageFromContext(ctx),
	}, nil)
}