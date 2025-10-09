package middleware

import (
	"go-springAi/internal/errors"
	"go-springAi/internal/logger"
	"go-springAi/internal/response"
	"context"
	"strings"
	"time"

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
func logError(zapLogger *zap.Logger, c *gin.Context, err error) {
	ctx := c.Request.Context()
	
	// 基础字段
	baseFields := []zap.Field{
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("client_ip", c.ClientIP()),
		zap.String("user_agent", c.Request.UserAgent()),
		zap.String("request_id", getRequestID(c)),
		zap.String("trace_id", getTraceID(c)),
		zap.Time("error_time", time.Now()),
	}

	// 添加用户信息（如果存在）
	if userID := getUserID(c); userID != "" {
		baseFields = append(baseFields, zap.String("user_id", userID))
	}

	// 添加查询参数和路径参数
	if c.Request.URL.RawQuery != "" {
		baseFields = append(baseFields, zap.String("query_params", c.Request.URL.RawQuery))
	}
	
	// 添加路径参数
	if len(c.Params) > 0 {
		params := make(map[string]string)
		for _, param := range c.Params {
			params[param.Key] = param.Value
		}
		baseFields = append(baseFields, zap.Any("path_params", params))
	}

	if appErr, ok := errors.IsAppError(err); ok {
		// 应用程序错误
		fields := append(baseFields,
			zap.String("error_code", string(appErr.Code)),
			zap.String("error_message", appErr.Message),
			zap.String("error_severity", string(appErr.Severity)),
			zap.Int("http_status", appErr.HTTPStatus),
			zap.Time("error_timestamp", appErr.Timestamp),
		)

		if appErr.Details != "" {
			fields = append(fields, zap.String("error_details", appErr.Details))
		}

		if appErr.Cause != nil {
			fields = append(fields, zap.Error(appErr.Cause))
		}

		if len(appErr.StackTrace) > 0 {
			fields = append(fields, zap.Strings("stack_trace", appErr.StackTrace))
		}

		// 根据错误严重级别和HTTP状态码记录日志
		logMessage := getLogMessage(appErr)
		
		switch appErr.Severity {
		case errors.SeverityCritical:
			logger.ErrorCtx(ctx, logMessage, fields...)
		case errors.SeverityHigh:
			if appErr.HTTPStatus >= 500 {
				logger.ErrorCtx(ctx, logMessage, fields...)
			} else {
				logger.WarnCtx(ctx, logMessage, fields...)
			}
		case errors.SeverityMedium:
			if appErr.HTTPStatus >= 500 {
				logger.ErrorCtx(ctx, logMessage, fields...)
			} else if appErr.HTTPStatus >= 400 {
				logger.WarnCtx(ctx, logMessage, fields...)
			} else {
				logger.InfoCtx(ctx, logMessage, fields...)
			}
		case errors.SeverityLow:
			if appErr.HTTPStatus >= 500 {
				logger.WarnCtx(ctx, logMessage, fields...)
			} else {
				logger.InfoCtx(ctx, logMessage, fields...)
			}
		default:
			logger.ErrorCtx(ctx, logMessage, fields...)
		}

		// 记录安全相关的错误
		logSecurityEvent(ctx, appErr, c)
		
	} else {
		// 普通错误，包装为内部错误
		fields := append(baseFields,
			zap.Error(err),
			zap.String("error_type", "unexpected_error"),
		)
		
		logger.ErrorCtx(ctx, "Unexpected error occurred", fields...)
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

// 辅助函数

// getRequestID 获取请求ID
func getRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return c.GetHeader("X-Request-ID")
}

// getTraceID 获取追踪ID
func getTraceID(c *gin.Context) string {
	if traceID, exists := c.Get("trace_id"); exists {
		if id, ok := traceID.(string); ok {
			return id
		}
	}
	return c.GetHeader("X-Trace-ID")
}

// getUserID 获取用户ID
func getUserID(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

// getLogMessage 根据错误类型生成日志消息
func getLogMessage(appErr *errors.AppError) string {
	switch {
	case strings.Contains(string(appErr.Code), "AUTH"):
		return logger.MsgAuthLoginFailed
	case strings.Contains(string(appErr.Code), "USER"):
		return logger.MsgUserValidation
	case strings.Contains(string(appErr.Code), "DATABASE"):
		return logger.MsgDBError
	case strings.Contains(string(appErr.Code), "VALIDATION"):
		return logger.MsgAPIValidation
	case strings.Contains(string(appErr.Code), "MCP"):
		return logger.MsgAPIError
	case appErr.HTTPStatus >= 500:
		return logger.MsgAPIError
	case appErr.HTTPStatus >= 400:
		return logger.MsgAPIError
	default:
		return logger.MsgAPIResponse
	}
}

// logSecurityEvent 记录安全相关事件
func logSecurityEvent(ctx context.Context, appErr *errors.AppError, c *gin.Context) {
	// 只记录安全相关的错误
	securityCodes := []errors.ErrorCode{
		errors.ErrCodeUnauthorized,
		errors.ErrCodeForbidden,
		errors.ErrCodeTokenExpired,
		errors.ErrCodeTokenInvalid,
		errors.ErrCodeLoginFailed,
		errors.ErrCodeAccountLocked,
		errors.ErrCodeAccountDisabled,
	}

	isSecurityEvent := false
	for _, code := range securityCodes {
		if appErr.Code == code {
			isSecurityEvent = true
			break
		}
	}

	if !isSecurityEvent {
		return
	}

	// 记录安全事件
	fields := []logger.LogField{
		logger.String("security_event", string(appErr.Code)),
		logger.String("client_ip", c.ClientIP()),
		logger.String("user_agent", c.Request.UserAgent()),
		logger.String("path", c.Request.URL.Path),
		logger.String("method", c.Request.Method),
		logger.String("request_id", getRequestID(c)),
		logger.Time("event_time", time.Now()),
	}

	if userID := getUserID(c); userID != "" {
		fields = append(fields, logger.String("user_id", userID))
	}

	// 添加额外的安全上下文
	if referer := c.GetHeader("Referer"); referer != "" {
		fields = append(fields, logger.String("referer", referer))
	}

	if origin := c.GetHeader("Origin"); origin != "" {
		fields = append(fields, logger.String("origin", origin))
	}

	logger.WarnCtx(ctx, "Security event detected", fields...)
}
