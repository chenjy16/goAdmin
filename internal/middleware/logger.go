package middleware

import (
	"context"
	"time"

	"admin/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ZapLogger 基于zap的结构化日志中间件
func ZapLogger(zapLogger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 记录请求开始
		logRequestStart(c, start)

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

		// 构建基础字段
		baseFields := buildBaseLogFields(c, method, path, statusCode, latency, clientIP, userAgent, bodySize)

		// 添加性能相关字段
		performanceFields := buildPerformanceFields(c, latency, bodySize)
		baseFields = append(baseFields, performanceFields...)

		// 添加安全相关字段
		securityFields := buildSecurityFields(c)
		baseFields = append(baseFields, securityFields...)

		// 根据状态码选择日志级别和消息
		ctx := c.Request.Context()
		logMessage := getResponseLogMessage(statusCode, latency)

		switch {
		case statusCode >= 500:
			logger.ErrorCtx(ctx, logMessage, baseFields...)
		case statusCode >= 400:
			logger.WarnCtx(ctx, logMessage, baseFields...)
		case statusCode >= 300:
			logger.InfoCtx(ctx, logMessage, baseFields...)
		default:
			// 只有在调试模式或慢请求时记录成功请求
			if shouldLogSuccessRequest(latency) {
				logger.InfoCtx(ctx, logMessage, baseFields...)
			}
		}

		// 记录慢请求
		logSlowRequest(ctx, c, latency, baseFields)
	}
}

// logRequestStart 记录请求开始
func logRequestStart(c *gin.Context, start time.Time) {
	// 只在调试模式下记录请求开始
	if gin.Mode() == gin.DebugMode {
		ctx := c.Request.Context()
		logger.DebugCtx(ctx, "Request started",
			logger.Module(logger.ModuleMiddleware),
			logger.Component("http"),
			logger.String("method", c.Request.Method),
			logger.String("path", c.Request.URL.Path),
			logger.String("client_ip", c.ClientIP()),
			logger.RequestID(getRequestID(c)),
			logger.String("user_agent", c.Request.UserAgent()),
			logger.Time("start_time", start))
	}
}

// buildBaseLogFields 构建基础日志字段
func buildBaseLogFields(c *gin.Context, method, path string, statusCode int, latency time.Duration, clientIP, userAgent string, bodySize int) []logger.LogField {
	fields := []logger.LogField{
		logger.Module(logger.ModuleMiddleware),
		logger.Component("http"),
		logger.String("method", method),
		logger.String("path", path),
		logger.Int("status", statusCode),
		logger.Duration("latency", latency),
		logger.String("client_ip", clientIP),
		logger.String("user_agent", userAgent),
		logger.Int("body_size", bodySize),
		logger.RequestID(getRequestID(c)),
	}

	// 添加追踪ID
	if traceID := getTraceID(c); traceID != "" {
		fields = append(fields, logger.String("trace_id", traceID))
	}

	// 添加用户ID
	if userID := getUserID(c); userID != "" {
		fields = append(fields, logger.UserID(userID))
	}

	// 添加会话ID
	if sessionID := getSessionID(c); sessionID != "" {
		fields = append(fields, logger.String("session_id", sessionID))
	}

	return fields
}

// buildPerformanceFields 构建性能相关字段
func buildPerformanceFields(c *gin.Context, latency time.Duration, bodySize int) []logger.LogField {
	fields := []logger.LogField{
		logger.Float64("latency_ms", float64(latency.Nanoseconds())/1e6),
		logger.Int64("memory_usage", getMemoryUsage()),
	}

	// 添加数据库查询统计
	if dbStats := getDatabaseStats(c); dbStats != nil {
		fields = append(fields,
			logger.Int("db_queries", dbStats.QueryCount),
			logger.Duration("db_total_time", dbStats.TotalTime))
	}

	// 添加缓存统计
	if cacheStats := getCacheStats(c); cacheStats != nil {
		fields = append(fields,
			logger.Int("cache_hits", cacheStats.Hits),
			logger.Int("cache_misses", cacheStats.Misses))
	}

	return fields
}

// buildSecurityFields 构建安全相关字段
func buildSecurityFields(c *gin.Context) []logger.LogField {
	fields := []logger.LogField{}

	// 添加认证信息
	if authType := getAuthType(c); authType != "" {
		fields = append(fields, logger.String("auth_type", authType))
	}

	// 添加权限信息
	if permissions := getPermissions(c); len(permissions) > 0 {
		fields = append(fields, logger.Strings("permissions", permissions))
	}

	// 添加IP地理位置信息
	if location := getIPLocation(c); location != "" {
		fields = append(fields, logger.String("ip_location", location))
	}

	// 检查可疑活动
	if isSuspicious := checkSuspiciousActivity(c); isSuspicious {
		fields = append(fields, logger.Bool("suspicious_activity", true))
	}

	return fields
}

// getResponseLogMessage 根据状态码和延迟获取日志消息
func getResponseLogMessage(statusCode int, latency time.Duration) string {
	switch {
	case statusCode >= 500:
		return logger.MsgAPIError
	case statusCode >= 400:
		return logger.MsgAPIError
	case latency > time.Second:
		return "Slow API response"
	default:
		return logger.MsgAPIResponse
	}
}

// shouldLogSuccessRequest 判断是否应该记录成功请求
func shouldLogSuccessRequest(latency time.Duration) bool {
	// 在调试模式下记录所有请求
	if gin.Mode() == gin.DebugMode {
		return true
	}
	// 在生产模式下只记录慢请求
	return latency > 500*time.Millisecond
}

// logSlowRequest 记录慢请求
func logSlowRequest(ctx context.Context, c *gin.Context, latency time.Duration, baseFields []logger.LogField) {
	slowThreshold := time.Second
	if latency > slowThreshold {
		fields := append(baseFields,
			logger.String("alert_type", "slow_request"),
			logger.Duration("threshold", slowThreshold),
			logger.Bool("is_slow", true))

		logger.WarnCtx(ctx, "Slow request detected", fields...)
	}
}

// getSessionID 获取会话ID
func getSessionID(c *gin.Context) string {
	if sessionID, exists := c.Get("session_id"); exists {
		if id, ok := sessionID.(string); ok {
			return id
		}
	}
	return ""
}

// getMemoryUsage 获取内存使用情况（简化实现）
func getMemoryUsage() int64 {
	// 这里可以集成runtime.MemStats或其他内存监控工具
	return 0
}

// DatabaseStats 数据库统计信息
type DatabaseStats struct {
	QueryCount int
	TotalTime  time.Duration
}

// getDatabaseStats 获取数据库统计信息
func getDatabaseStats(c *gin.Context) *DatabaseStats {
	if stats, exists := c.Get("db_stats"); exists {
		if dbStats, ok := stats.(*DatabaseStats); ok {
			return dbStats
		}
	}
	return nil
}

// CacheStats 缓存统计信息
type CacheStats struct {
	Hits   int
	Misses int
}

// getCacheStats 获取缓存统计信息
func getCacheStats(c *gin.Context) *CacheStats {
	if stats, exists := c.Get("cache_stats"); exists {
		if cacheStats, ok := stats.(*CacheStats); ok {
			return cacheStats
		}
	}
	return nil
}

// getAuthType 获取认证类型
func getAuthType(c *gin.Context) string {
	if authType, exists := c.Get("auth_type"); exists {
		if auth, ok := authType.(string); ok {
			return auth
		}
	}
	return ""
}

// getPermissions 获取用户权限
func getPermissions(c *gin.Context) []string {
	if permissions, exists := c.Get("permissions"); exists {
		if perms, ok := permissions.([]string); ok {
			return perms
		}
	}
	return nil
}

// getIPLocation 获取IP地理位置
func getIPLocation(c *gin.Context) string {
	// 这里可以集成IP地理位置服务
	return ""
}

// checkSuspiciousActivity 检查可疑活动
func checkSuspiciousActivity(c *gin.Context) bool {
	// 这里可以实现可疑活动检测逻辑
	// 例如：频繁请求、异常IP、恶意User-Agent等
	return false
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
