package errors

import (
	"fmt"
	"net/http"
	"runtime"
	"time"
)

// ErrorCode 错误码类型
type ErrorCode string

// ErrorSeverity 错误严重级别
type ErrorSeverity string

const (
	SeverityLow      ErrorSeverity = "LOW"      // 低级别：用户输入错误等
	SeverityMedium   ErrorSeverity = "MEDIUM"   // 中级别：业务逻辑错误等
	SeverityHigh     ErrorSeverity = "HIGH"     // 高级别：系统错误等
	SeverityCritical ErrorSeverity = "CRITICAL" // 严重级别：数据库连接失败等
)

// 定义错误码常量
const (
	// 通用错误码
	ErrCodeInternal ErrorCode = "INTERNAL_ERROR"
	ErrCodeConflict ErrorCode = "CONFLICT"

	// 用户相关错误码
	ErrCodeUserNotFound ErrorCode = "USER_NOT_FOUND"

	// 数据库相关错误码
	ErrCodeDatabaseQuery ErrorCode = "DATABASE_QUERY_ERROR"

	// 验证相关错误码
	ErrCodeValidationFailed ErrorCode = "VALIDATION_FAILED"
)

// AppError 应用程序自定义错误
type AppError struct {
	Code       ErrorCode     `json:"code"`
	Message    string        `json:"message"`
	Details    string        `json:"details,omitempty"`
	Severity   ErrorSeverity `json:"severity"`
	HTTPStatus int           `json:"-"`
	Timestamp  time.Time     `json:"timestamp"`
	StackTrace []string      `json:"stack_trace,omitempty"`
	Cause      error         `json:"-"`
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 实现错误链
func (e *AppError) Unwrap() error {
	return e.Cause
}

// WithCause 添加原因错误
func (e *AppError) WithCause(cause error) *AppError {
	e.Cause = cause
	return e
}

// WithDetails 添加详细信息
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// WithStackTrace 添加堆栈跟踪
func (e *AppError) WithStackTrace() *AppError {
	e.StackTrace = getStackTrace()
	return e
}

// NewAppError 创建新的应用程序错误
func NewAppError(code ErrorCode, message string, severity ErrorSeverity, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Severity:   severity,
		HTTPStatus: httpStatus,
		Timestamp:  time.Now(),
	}
}

// 预定义的错误创建函数

// NewInternalError 创建内部错误
func NewInternalError(message string) *AppError {
	return NewAppError(ErrCodeInternal, message, SeverityHigh, http.StatusInternalServerError).WithStackTrace()
}

// NewValidationError 创建验证错误
func NewValidationError(message string) *AppError {
	return NewAppError(ErrCodeValidationFailed, message, SeverityLow, http.StatusBadRequest)
}



// NewConflictError 创建冲突错误
func NewConflictError(message string) *AppError {
	return NewAppError(ErrCodeConflict, message, SeverityMedium, http.StatusConflict)
}

// 用户相关错误

// NewUserNotFoundError 创建用户未找到错误
func NewUserNotFoundError() *AppError {
	return NewAppError(ErrCodeUserNotFound, "User not found", SeverityLow, http.StatusNotFound)
}



// 数据库相关错误

// NewDatabaseError 创建数据库错误
func NewDatabaseError(operation string, cause error) *AppError {
	return NewAppError(ErrCodeDatabaseQuery,
		fmt.Sprintf("Database %s failed", operation),
		SeverityHigh, http.StatusInternalServerError).WithCause(cause).WithStackTrace()
}



// 工具函数

// getStackTrace 获取堆栈跟踪
func getStackTrace() []string {
	var stackTrace []string
	// 获取更多的堆栈信息，但限制在合理范围内
	for i := 2; i < 15; i++ { // 跳过当前函数和调用者
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			// 只保留文件名，不包含完整路径，减少日志大小
			shortFile := file
			if lastSlash := len(file) - 1; lastSlash >= 0 {
				for i := lastSlash; i >= 0; i-- {
					if file[i] == '/' {
						shortFile = file[i+1:]
						break
					}
				}
			}
			stackTrace = append(stackTrace, fmt.Sprintf("%s:%d %s", shortFile, line, fn.Name()))
		}
	}
	return stackTrace
}

// IsAppError 检查是否为应用程序错误
func IsAppError(err error) (*AppError, bool) {
	if appErr, ok := err.(*AppError); ok {
		return appErr, true
	}
	return nil, false
}

// WrapError 包装普通错误为应用程序错误
func WrapError(err error, code ErrorCode, message string, severity ErrorSeverity, httpStatus int) *AppError {
	if err == nil {
		return nil
	}

	// 如果已经是 AppError，直接返回
	if appErr, ok := IsAppError(err); ok {
		return appErr
	}

	return NewAppError(code, message, severity, httpStatus).WithCause(err).WithStackTrace()
}
