package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogLevel 日志级别常量
type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	FatalLevel LogLevel = "fatal"
)

// LogField 日志字段类型
type LogField = zap.Field

// 常用字段构造函数
var (
	String    = zap.String
	Int       = zap.Int
	Int64     = zap.Int64
	Uint64    = zap.Uint64
	Uint32    = zap.Uint32
	Float64   = zap.Float64
	Bool      = zap.Bool
	Duration  = zap.Duration
	Time      = zap.Time
	ZapError  = zap.Error
	Any       = zap.Any
	Strings   = zap.Strings
	Ints      = zap.Ints
)

// 业务相关的常用字段
func UserID(id string) LogField {
	return zap.String("user_id", id)
}

func RequestID(id string) LogField {
	return zap.String("request_id", id)
}

func TraceID(id string) LogField {
	return zap.String("trace_id", id)
}

func Module(name string) LogField {
	return zap.String("module", name)
}

func Operation(op string) LogField {
	return zap.String("operation", op)
}

func Component(comp string) LogField {
	return zap.String("component", comp)
}

// Logger 统一日志接口
type Logger interface {
	Debug(msg string, fields ...LogField)
	Info(msg string, fields ...LogField)
	Warn(msg string, fields ...LogField)
	Error(msg string, fields ...LogField)
	Fatal(msg string, fields ...LogField)
	
	// 带上下文的日志方法
	DebugCtx(ctx context.Context, msg string, fields ...LogField)
	InfoCtx(ctx context.Context, msg string, fields ...LogField)
	WarnCtx(ctx context.Context, msg string, fields ...LogField)
	ErrorCtx(ctx context.Context, msg string, fields ...LogField)
	
	// 创建子日志器
	With(fields ...LogField) Logger
	WithContext(ctx context.Context) Logger
}

// zapLogger zap日志器的包装
type zapLogger struct {
	logger *zap.Logger
}

// NewLogger 创建新的日志器
func NewLogger(mode string) (Logger, error) {
	var config zap.Config
	
	if mode == "release" || mode == "production" {
		config = zap.NewProductionConfig()
		// 生产环境配置
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		config.OutputPaths = []string{"stdout"}
		config.ErrorOutputPaths = []string{"stderr"}
	} else {
		config = zap.NewDevelopmentConfig()
		// 开发环境配置
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		config.Development = true
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}
	
	// 统一的字段配置
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.NameKey = "logger"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.StacktraceKey = "stacktrace"
	
	logger, err := config.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, err
	}
	
	return &zapLogger{logger: logger}, nil
}

// NewLoggerFromZap 从现有的zap.Logger创建Logger
func NewLoggerFromZap(zapLog *zap.Logger) Logger {
	return &zapLogger{logger: zapLog}
}

func (l *zapLogger) Debug(msg string, fields ...LogField) {
	l.logger.Debug(msg, fields...)
}

func (l *zapLogger) Info(msg string, fields ...LogField) {
	l.logger.Info(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...LogField) {
	l.logger.Warn(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...LogField) {
	l.logger.Error(msg, fields...)
}

func (l *zapLogger) Fatal(msg string, fields ...LogField) {
	l.logger.Fatal(msg, fields...)
}

func (l *zapLogger) DebugCtx(ctx context.Context, msg string, fields ...LogField) {
	fields = l.addContextFields(ctx, fields)
	l.logger.Debug(msg, fields...)
}

func (l *zapLogger) InfoCtx(ctx context.Context, msg string, fields ...LogField) {
	fields = l.addContextFields(ctx, fields)
	l.logger.Info(msg, fields...)
}

func (l *zapLogger) WarnCtx(ctx context.Context, msg string, fields ...LogField) {
	fields = l.addContextFields(ctx, fields)
	l.logger.Warn(msg, fields...)
}

func (l *zapLogger) ErrorCtx(ctx context.Context, msg string, fields ...LogField) {
	fields = l.addContextFields(ctx, fields)
	l.logger.Error(msg, fields...)
}

func (l *zapLogger) With(fields ...LogField) Logger {
	return &zapLogger{logger: l.logger.With(fields...)}
}

func (l *zapLogger) WithContext(ctx context.Context) Logger {
	fields := l.extractContextFields(ctx)
	return &zapLogger{logger: l.logger.With(fields...)}
}

// addContextFields 从上下文中提取字段并添加到日志中
func (l *zapLogger) addContextFields(ctx context.Context, fields []LogField) []LogField {
	contextFields := l.extractContextFields(ctx)
	return append(fields, contextFields...)
}

// extractContextFields 从上下文中提取日志字段
func (l *zapLogger) extractContextFields(ctx context.Context) []LogField {
	var fields []LogField
	
	// 提取请求ID
	if requestID := ctx.Value("request_id"); requestID != nil {
		if id, ok := requestID.(string); ok {
			fields = append(fields, RequestID(id))
		}
	}
	
	// 提取用户ID
	if userID := ctx.Value("user_id"); userID != nil {
		if id, ok := userID.(string); ok {
			fields = append(fields, UserID(id))
		}
	}
	
	// 提取追踪ID
	if traceID := ctx.Value("trace_id"); traceID != nil {
		if id, ok := traceID.(string); ok {
			fields = append(fields, TraceID(id))
		}
	}
	
	return fields
}

// 全局日志器实例
var globalLogger Logger

// InitGlobalLogger 初始化全局日志器
func InitGlobalLogger(mode string) error {
	logger, err := NewLogger(mode)
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

// GetGlobalLogger 获取全局日志器
func GetGlobalLogger() Logger {
	if globalLogger == nil {
		// 如果全局日志器未初始化，创建一个默认的
		logger, _ := NewLogger("development")
		globalLogger = logger
	}
	return globalLogger
}

// 全局日志函数
func Debug(msg string, fields ...LogField) {
	GetGlobalLogger().Debug(msg, fields...)
}

func Info(msg string, fields ...LogField) {
	GetGlobalLogger().Info(msg, fields...)
}

func Warn(msg string, fields ...LogField) {
	GetGlobalLogger().Warn(msg, fields...)
}

func LogError(msg string, fields ...LogField) {
	GetGlobalLogger().Error(msg, fields...)
}

func Fatal(msg string, fields ...LogField) {
	GetGlobalLogger().Fatal(msg, fields...)
}

// 带上下文的全局日志函数
func DebugCtx(ctx context.Context, msg string, fields ...LogField) {
	GetGlobalLogger().DebugCtx(ctx, msg, fields...)
}

func InfoCtx(ctx context.Context, msg string, fields ...LogField) {
	GetGlobalLogger().InfoCtx(ctx, msg, fields...)
}

func WarnCtx(ctx context.Context, msg string, fields ...LogField) {
	GetGlobalLogger().WarnCtx(ctx, msg, fields...)
}

func ErrorCtx(ctx context.Context, msg string, fields ...LogField) {
	GetGlobalLogger().ErrorCtx(ctx, msg, fields...)
}

// Sync 同步日志缓冲区
func Sync() error {
	if zl, ok := globalLogger.(*zapLogger); ok {
		return zl.logger.Sync()
	}
	return nil
}