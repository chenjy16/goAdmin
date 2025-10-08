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
	ErrCodeInternal   ErrorCode = "INTERNAL_ERROR"
	ErrCodeConflict   ErrorCode = "CONFLICT"
	ErrCodeNotFound   ErrorCode = "NOT_FOUND"
	ErrCodeBadRequest ErrorCode = "BAD_REQUEST"
	ErrCodeTimeout    ErrorCode = "TIMEOUT"
	ErrCodeRateLimit  ErrorCode = "RATE_LIMIT_EXCEEDED"

	// 认证和授权相关错误码
	ErrCodeUnauthorized     ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden        ErrorCode = "FORBIDDEN"
	ErrCodeTokenExpired     ErrorCode = "TOKEN_EXPIRED"
	ErrCodeTokenInvalid     ErrorCode = "TOKEN_INVALID"
	ErrCodeLoginFailed      ErrorCode = "LOGIN_FAILED"
	ErrCodePasswordWeak     ErrorCode = "PASSWORD_WEAK"
	ErrCodeAccountLocked    ErrorCode = "ACCOUNT_LOCKED"
	ErrCodeAccountDisabled  ErrorCode = "ACCOUNT_DISABLED"

	// 用户相关错误码
	ErrCodeUserNotFound     ErrorCode = "USER_NOT_FOUND"
	ErrCodeUserExists       ErrorCode = "USER_ALREADY_EXISTS"
	ErrCodeUserInactive     ErrorCode = "USER_INACTIVE"
	ErrCodeEmailExists      ErrorCode = "EMAIL_ALREADY_EXISTS"
	ErrCodeUsernameExists   ErrorCode = "USERNAME_ALREADY_EXISTS"

	// 数据库相关错误码
	ErrCodeDatabaseQuery      ErrorCode = "DATABASE_QUERY_ERROR"
	ErrCodeDatabaseConnection ErrorCode = "DATABASE_CONNECTION_ERROR"
	ErrCodeDatabaseTransaction ErrorCode = "DATABASE_TRANSACTION_ERROR"
	ErrCodeDatabaseConstraint  ErrorCode = "DATABASE_CONSTRAINT_ERROR"
	ErrCodeDatabaseDeadlock    ErrorCode = "DATABASE_DEADLOCK_ERROR"

	// 验证相关错误码
	ErrCodeValidationFailed ErrorCode = "VALIDATION_FAILED"
	ErrCodeInvalidFormat    ErrorCode = "INVALID_FORMAT"
	ErrCodeMissingField     ErrorCode = "MISSING_REQUIRED_FIELD"
	ErrCodeInvalidValue     ErrorCode = "INVALID_VALUE"

	// 业务逻辑相关错误码
	ErrCodeBusinessLogic    ErrorCode = "BUSINESS_LOGIC_ERROR"
	ErrCodeOperationFailed  ErrorCode = "OPERATION_FAILED"
	ErrCodeResourceBusy     ErrorCode = "RESOURCE_BUSY"
	ErrCodeQuotaExceeded    ErrorCode = "QUOTA_EXCEEDED"

	// 网络和外部服务相关错误码
	ErrCodeNetworkError     ErrorCode = "NETWORK_ERROR"
	ErrCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	ErrCodeExternalService  ErrorCode = "EXTERNAL_SERVICE_ERROR"

	// 文件和存储相关错误码
	ErrCodeFileNotFound     ErrorCode = "FILE_NOT_FOUND"
	ErrCodeFileUploadFailed ErrorCode = "FILE_UPLOAD_FAILED"
	ErrCodeStorageError     ErrorCode = "STORAGE_ERROR"
	ErrCodeFileTooLarge     ErrorCode = "FILE_TOO_LARGE"

	// MCP相关错误码
	ErrCodeMCPInitFailed    ErrorCode = "MCP_INIT_FAILED"
	ErrCodeMCPToolNotFound  ErrorCode = "MCP_TOOL_NOT_FOUND"
	ErrCodeMCPExecuteFailed ErrorCode = "MCP_EXECUTE_FAILED"
	ErrCodeMCPInvalidParams ErrorCode = "MCP_INVALID_PARAMS"
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

// NewBadRequestError 创建请求错误
func NewBadRequestError(message string) *AppError {
	return NewAppError(ErrCodeBadRequest, message, SeverityLow, http.StatusBadRequest)
}

// NewNotFoundError 创建资源未找到错误
func NewNotFoundError(resource string) *AppError {
	return NewAppError(ErrCodeNotFound, fmt.Sprintf("%s not found", resource), SeverityLow, http.StatusNotFound)
}

// NewConflictError 创建冲突错误
func NewConflictError(message string) *AppError {
	return NewAppError(ErrCodeConflict, message, SeverityMedium, http.StatusConflict)
}

// NewTimeoutError 创建超时错误
func NewTimeoutError(operation string) *AppError {
	return NewAppError(ErrCodeTimeout, fmt.Sprintf("%s timeout", operation), SeverityMedium, http.StatusRequestTimeout)
}

// NewRateLimitError 创建限流错误
func NewRateLimitError() *AppError {
	return NewAppError(ErrCodeRateLimit, "Rate limit exceeded", SeverityMedium, http.StatusTooManyRequests)
}

// 认证和授权相关错误

// NewUnauthorizedError 创建未授权错误
func NewUnauthorizedError(message string) *AppError {
	if message == "" {
		message = "Unauthorized access"
	}
	return NewAppError(ErrCodeUnauthorized, message, SeverityMedium, http.StatusUnauthorized)
}

// NewForbiddenError 创建禁止访问错误
func NewForbiddenError(message string) *AppError {
	if message == "" {
		message = "Access forbidden"
	}
	return NewAppError(ErrCodeForbidden, message, SeverityMedium, http.StatusForbidden)
}

// NewTokenExpiredError 创建令牌过期错误
func NewTokenExpiredError() *AppError {
	return NewAppError(ErrCodeTokenExpired, "Token has expired", SeverityLow, http.StatusUnauthorized)
}

// NewTokenInvalidError 创建令牌无效错误
func NewTokenInvalidError() *AppError {
	return NewAppError(ErrCodeTokenInvalid, "Invalid token", SeverityLow, http.StatusUnauthorized)
}

// NewLoginFailedError 创建登录失败错误
func NewLoginFailedError() *AppError {
	return NewAppError(ErrCodeLoginFailed, "Invalid username or password", SeverityLow, http.StatusUnauthorized)
}

// NewAccountLockedError 创建账户锁定错误
func NewAccountLockedError() *AppError {
	return NewAppError(ErrCodeAccountLocked, "Account is locked", SeverityMedium, http.StatusForbidden)
}

// NewAccountDisabledError 创建账户禁用错误
func NewAccountDisabledError() *AppError {
	return NewAppError(ErrCodeAccountDisabled, "Account is disabled", SeverityMedium, http.StatusForbidden)
}

// 用户相关错误

// NewUserNotFoundError 创建用户未找到错误
func NewUserNotFoundError() *AppError {
	return NewAppError(ErrCodeUserNotFound, "User not found", SeverityLow, http.StatusNotFound)
}

// NewUserExistsError 创建用户已存在错误
func NewUserExistsError() *AppError {
	return NewAppError(ErrCodeUserExists, "User already exists", SeverityLow, http.StatusConflict)
}

// NewEmailExistsError 创建邮箱已存在错误
func NewEmailExistsError() *AppError {
	return NewAppError(ErrCodeEmailExists, "Email already exists", SeverityLow, http.StatusConflict)
}

// NewUsernameExistsError 创建用户名已存在错误
func NewUsernameExistsError() *AppError {
	return NewAppError(ErrCodeUsernameExists, "Username already exists", SeverityLow, http.StatusConflict)
}

// NewUserInactiveError 创建用户未激活错误
func NewUserInactiveError() *AppError {
	return NewAppError(ErrCodeUserInactive, "User account is inactive", SeverityLow, http.StatusForbidden)
}

// 数据库相关错误

// NewDatabaseError 创建数据库错误
func NewDatabaseError(operation string, cause error) *AppError {
	return NewAppError(ErrCodeDatabaseQuery,
		fmt.Sprintf("Database %s failed", operation),
		SeverityHigh, http.StatusInternalServerError).WithCause(cause).WithStackTrace()
}

// NewDatabaseConnectionError 创建数据库连接错误
func NewDatabaseConnectionError(cause error) *AppError {
	return NewAppError(ErrCodeDatabaseConnection,
		"Database connection failed",
		SeverityCritical, http.StatusInternalServerError).WithCause(cause).WithStackTrace()
}

// NewDatabaseTransactionError 创建数据库事务错误
func NewDatabaseTransactionError(operation string, cause error) *AppError {
	return NewAppError(ErrCodeDatabaseTransaction,
		fmt.Sprintf("Database transaction %s failed", operation),
		SeverityHigh, http.StatusInternalServerError).WithCause(cause).WithStackTrace()
}

// NewDatabaseConstraintError 创建数据库约束错误
func NewDatabaseConstraintError(constraint string) *AppError {
	return NewAppError(ErrCodeDatabaseConstraint,
		fmt.Sprintf("Database constraint violation: %s", constraint),
		SeverityMedium, http.StatusConflict)
}

// 业务逻辑相关错误

// NewBusinessLogicError 创建业务逻辑错误
func NewBusinessLogicError(message string) *AppError {
	return NewAppError(ErrCodeBusinessLogic, message, SeverityMedium, http.StatusBadRequest)
}

// NewOperationFailedError 创建操作失败错误
func NewOperationFailedError(operation string) *AppError {
	return NewAppError(ErrCodeOperationFailed,
		fmt.Sprintf("Operation %s failed", operation),
		SeverityMedium, http.StatusInternalServerError)
}

// NewResourceBusyError 创建资源忙错误
func NewResourceBusyError(resource string) *AppError {
	return NewAppError(ErrCodeResourceBusy,
		fmt.Sprintf("Resource %s is busy", resource),
		SeverityMedium, http.StatusConflict)
}

// 网络和外部服务相关错误

// NewNetworkError 创建网络错误
func NewNetworkError(operation string, cause error) *AppError {
	return NewAppError(ErrCodeNetworkError,
		fmt.Sprintf("Network %s failed", operation),
		SeverityHigh, http.StatusBadGateway).WithCause(cause)
}

// NewServiceUnavailableError 创建服务不可用错误
func NewServiceUnavailableError(service string) *AppError {
	return NewAppError(ErrCodeServiceUnavailable,
		fmt.Sprintf("Service %s is unavailable", service),
		SeverityHigh, http.StatusServiceUnavailable)
}

// 文件和存储相关错误

// NewFileNotFoundError 创建文件未找到错误
func NewFileNotFoundError(filename string) *AppError {
	return NewAppError(ErrCodeFileNotFound,
		fmt.Sprintf("File %s not found", filename),
		SeverityLow, http.StatusNotFound)
}

// NewFileUploadFailedError 创建文件上传失败错误
func NewFileUploadFailedError(reason string) *AppError {
	return NewAppError(ErrCodeFileUploadFailed,
		fmt.Sprintf("File upload failed: %s", reason),
		SeverityMedium, http.StatusBadRequest)
}

// NewFileTooLargeError 创建文件过大错误
func NewFileTooLargeError(maxSize string) *AppError {
	return NewAppError(ErrCodeFileTooLarge,
		fmt.Sprintf("File size exceeds maximum allowed size of %s", maxSize),
		SeverityLow, http.StatusRequestEntityTooLarge)
}

// MCP相关错误

// NewMCPInitFailedError 创建MCP初始化失败错误
func NewMCPInitFailedError(reason string) *AppError {
	return NewAppError(ErrCodeMCPInitFailed,
		fmt.Sprintf("MCP initialization failed: %s", reason),
		SeverityHigh, http.StatusInternalServerError)
}

// NewMCPToolNotFoundError 创建MCP工具未找到错误
func NewMCPToolNotFoundError(toolName string) *AppError {
	return NewAppError(ErrCodeMCPToolNotFound,
		fmt.Sprintf("MCP tool '%s' not found", toolName),
		SeverityLow, http.StatusNotFound)
}

// NewMCPExecuteFailedError 创建MCP执行失败错误
func NewMCPExecuteFailedError(toolName string, reason string) *AppError {
	return NewAppError(ErrCodeMCPExecuteFailed,
		fmt.Sprintf("MCP tool '%s' execution failed: %s", toolName, reason),
		SeverityMedium, http.StatusInternalServerError)
}

// NewMCPInvalidParamsError 创建MCP参数无效错误
func NewMCPInvalidParamsError(toolName string, reason string) *AppError {
	return NewAppError(ErrCodeMCPInvalidParams,
		fmt.Sprintf("Invalid parameters for MCP tool '%s': %s", toolName, reason),
		SeverityLow, http.StatusBadRequest)
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
