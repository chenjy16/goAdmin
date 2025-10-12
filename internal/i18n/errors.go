package i18n

import (
	"go-springAi/internal/errors"
	"net/http"
)

// ErrorMessageMap 错误码到消息ID的映射
var ErrorMessageMap = map[errors.ErrorCode]string{
	// 通用错误码
	errors.ErrCodeInternal:   "error.internal",
	errors.ErrCodeConflict:   "error.conflict",
	errors.ErrCodeNotFound:   "error.not.found",
	errors.ErrCodeBadRequest: "error.bad.request",
	errors.ErrCodeTimeout:    "error.timeout",
	errors.ErrCodeRateLimit:  "error.rate.limit",

	// 认证和授权相关错误码
	errors.ErrCodeUnauthorized:     "error.unauthorized",
	errors.ErrCodeForbidden:        "error.forbidden",
	errors.ErrCodeTokenExpired:     "error.token.expired",
	errors.ErrCodeTokenInvalid:     "error.token.invalid",
	errors.ErrCodeLoginFailed:      "error.login.failed",
	errors.ErrCodePasswordWeak:     "error.password.weak",
	errors.ErrCodeAccountLocked:    "error.account.locked",
	errors.ErrCodeAccountDisabled:  "error.account.disabled",

	// 用户相关错误码
	errors.ErrCodeUserNotFound:   "error.user.not.found",
	errors.ErrCodeUserExists:     "error.user.exists",
	errors.ErrCodeUserInactive:   "error.user.inactive",
	errors.ErrCodeEmailExists:    "error.email.exists",
	errors.ErrCodeUsernameExists: "error.username.exists",

	// 数据库相关错误码
	errors.ErrCodeDatabaseQuery:       "error.database.query",
	errors.ErrCodeDatabaseConnection:  "error.database.connection",
	errors.ErrCodeDatabaseTransaction: "error.database.transaction",
	errors.ErrCodeDatabaseConstraint:  "error.database.constraint",
	errors.ErrCodeDatabaseDeadlock:    "error.database.deadlock",

	// 验证相关错误码
	errors.ErrCodeValidationFailed: "error.validation.failed",
	errors.ErrCodeInvalidFormat:    "error.invalid.format",
	errors.ErrCodeMissingField:     "error.missing.field",
	errors.ErrCodeInvalidValue:     "error.invalid.value",

	// 业务逻辑相关错误码
	errors.ErrCodeBusinessLogic:   "error.business.logic",
	errors.ErrCodeOperationFailed: "error.operation.failed",
	errors.ErrCodeResourceBusy:    "error.resource.busy",
	errors.ErrCodeQuotaExceeded:   "error.quota.exceeded",

	// 网络和外部服务相关错误码
	errors.ErrCodeNetworkError:       "error.network",
	errors.ErrCodeServiceUnavailable: "error.service.unavailable",
	errors.ErrCodeExternalService:    "error.external.service",

	// 文件和存储相关错误码
	errors.ErrCodeFileNotFound:     "error.file.not.found",
	errors.ErrCodeFileUploadFailed: "error.file.upload.failed",
	errors.ErrCodeStorageError:     "error.storage",
	errors.ErrCodeFileTooLarge:     "error.file.too.large",
}

// GetErrorMessage 获取错误的国际化消息
func (m *Manager) GetErrorMessage(lang string, appErr *errors.AppError) string {
	if messageID, exists := ErrorMessageMap[appErr.Code]; exists {
		return m.T(lang, messageID, map[string]interface{}{
			"Details": appErr.Details,
		})
	}
	
	// 如果没有找到映射，返回原始消息
	return appErr.Message
}

// NewInternalError 创建国际化的内部错误
func (m *Manager) NewInternalError(lang string, details string) *errors.AppError {
	message := m.T(lang, "error.internal", nil)
	return errors.NewAppError(errors.ErrCodeInternal, message, errors.SeverityHigh, http.StatusInternalServerError).
		WithDetails(details).
		WithStackTrace()
}

// NewValidationError 创建国际化的验证错误
func (m *Manager) NewValidationError(lang string, field string) *errors.AppError {
	message := m.T(lang, "error.validation.failed", map[string]interface{}{
		"Field": field,
	})
	return errors.NewAppError(errors.ErrCodeValidationFailed, message, errors.SeverityLow, http.StatusBadRequest)
}

// NewNotFoundError 创建国际化的未找到错误
func (m *Manager) NewNotFoundError(lang string, resource string) *errors.AppError {
	message := m.T(lang, "error.not.found", map[string]interface{}{
		"Resource": resource,
	})
	return errors.NewAppError(errors.ErrCodeNotFound, message, errors.SeverityLow, http.StatusNotFound)
}

// NewUnauthorizedError 创建国际化的未授权错误
func (m *Manager) NewUnauthorizedError(lang string) *errors.AppError {
	message := m.T(lang, "error.unauthorized", nil)
	return errors.NewAppError(errors.ErrCodeUnauthorized, message, errors.SeverityMedium, http.StatusUnauthorized)
}

// NewConflictError 创建国际化的冲突错误
func (m *Manager) NewConflictError(lang string, resource string) *errors.AppError {
	message := m.T(lang, "error.conflict", map[string]interface{}{
		"Resource": resource,
	})
	return errors.NewAppError(errors.ErrCodeConflict, message, errors.SeverityMedium, http.StatusConflict)
}

// NewTimeoutError 创建国际化的超时错误
func (m *Manager) NewTimeoutError(lang string, operation string) *errors.AppError {
	message := m.T(lang, "error.timeout", map[string]interface{}{
		"Operation": operation,
	})
	return errors.NewAppError(errors.ErrCodeTimeout, message, errors.SeverityMedium, http.StatusRequestTimeout)
}