package middleware

import (
	"admin/internal/errors"
	"admin/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorHandler 错误处理中间件
func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// 记录错误日志
			logError(logger, c, err)

			// 处理错误响应
			handleErrorResponse(c, err)
		}
	}
}

// logError 记录错误日志
func logError(logger *zap.Logger, c *gin.Context, err error) {
	if appErr, ok := errors.IsAppError(err); ok {
		// 应用程序错误
		fields := []zap.Field{
			zap.String("error_code", string(appErr.Code)),
			zap.String("message", appErr.Message),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Time("timestamp", appErr.Timestamp),
		}

		if appErr.Details != "" {
			fields = append(fields, zap.String("details", appErr.Details))
		}

		if appErr.Cause != nil {
			fields = append(fields, zap.Error(appErr.Cause))
		}

		if len(appErr.StackTrace) > 0 {
			fields = append(fields, zap.Strings("stack_trace", appErr.StackTrace))
		}

		// 根据错误级别记录日志
		if appErr.HTTPStatus >= 500 {
			logger.Error("Internal server error", fields...)
		} else if appErr.HTTPStatus >= 400 {
			logger.Warn("Client error", fields...)
		} else {
			logger.Info("Request error", fields...)
		}
	} else {
		// 普通错误
		logger.Error("Unexpected error",
			zap.Error(err),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)
	}
}

// handleErrorResponse 处理错误响应
func handleErrorResponse(c *gin.Context, err error) {
	if appErr, ok := errors.IsAppError(err); ok {
		// 应用程序错误
		response.Error(c, appErr.HTTPStatus, appErr.Message, string(appErr.Code))
	} else {
		// 普通错误，返回通用内部服务器错误
		response.InternalServerError(c, "Internal Server Error", "INTERNAL_ERROR")
	}
}

// AbortWithError 中止请求并设置错误
func AbortWithError(c *gin.Context, err error) {
	c.Error(err)
	c.Abort()
}

// AbortWithAppError 中止请求并设置应用程序错误
func AbortWithAppError(c *gin.Context, appErr *errors.AppError) {
	c.Error(appErr)
	c.Abort()
}
