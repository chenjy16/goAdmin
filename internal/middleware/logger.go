package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ZapLogger 基于zap的结构化日志中间件
func ZapLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 计算延迟
		latency := time.Since(start)

		// 获取请求信息
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()
		userAgent := c.Request.UserAgent()

		if raw != "" {
			path = path + "?" + raw
		}

		// 构建日志字段
		fields := []zap.Field{
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", userAgent),
			zap.Int("body_size", bodySize),
		}

		// 添加请求ID（如果存在）
		if requestID := c.GetHeader("X-Request-ID"); requestID != "" {
			fields = append(fields, zap.String("request_id", requestID))
		}

		// 根据状态码选择日志级别
		switch {
		case statusCode >= 500:
			logger.Error("Server error", fields...)
		case statusCode >= 400:
			logger.Warn("Client error", fields...)
		case statusCode >= 300:
			logger.Info("Redirection", fields...)
		default:
			logger.Info("Request completed", fields...)
		}
	}
}

// RequestID 请求ID中间件
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	// 简单的时间戳+随机数生成请求ID
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// StructuredLogger 结构化日志记录器
type StructuredLogger struct {
	logger *zap.Logger
}

// NewStructuredLogger 创建结构化日志记录器
func NewStructuredLogger(logger *zap.Logger) *StructuredLogger {
	return &StructuredLogger{logger: logger}
}

// LogRequest 记录请求日志
func (sl *StructuredLogger) LogRequest(c *gin.Context, fields ...zap.Field) {
	baseFields := []zap.Field{
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("client_ip", c.ClientIP()),
		zap.String("user_agent", c.Request.UserAgent()),
	}

	if requestID := c.GetString("request_id"); requestID != "" {
		baseFields = append(baseFields, zap.String("request_id", requestID))
	}

	allFields := append(baseFields, fields...)
	sl.logger.Info("Request received", allFields...)
}

// LogResponse 记录响应日志
func (sl *StructuredLogger) LogResponse(c *gin.Context, latency time.Duration, fields ...zap.Field) {
	baseFields := []zap.Field{
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.Int("status", c.Writer.Status()),
		zap.Duration("latency", latency),
		zap.String("client_ip", c.ClientIP()),
		zap.Int("body_size", c.Writer.Size()),
	}

	if requestID := c.GetString("request_id"); requestID != "" {
		baseFields = append(baseFields, zap.String("request_id", requestID))
	}

	allFields := append(baseFields, fields...)

	statusCode := c.Writer.Status()
	switch {
	case statusCode >= 500:
		sl.logger.Error("Response sent", allFields...)
	case statusCode >= 400:
		sl.logger.Warn("Response sent", allFields...)
	default:
		sl.logger.Info("Response sent", allFields...)
	}
}
