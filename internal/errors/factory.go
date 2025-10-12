package errors

import (
	"fmt"
)

// ErrorFactory 错误工厂，提供常见错误模式的统一创建方法
type ErrorFactory struct{}

// NewErrorFactory 创建错误工厂实例
func NewErrorFactory() *ErrorFactory {
	return &ErrorFactory{}
}

// FailedToOperation 创建"操作失败"类型的错误
func (f *ErrorFactory) FailedToOperation(operation string, cause error) *AppError {
	return NewOperationFailedError(operation).WithCause(cause).WithStackTrace()
}

// FailedToGet 创建"获取失败"类型的错误
func (f *ErrorFactory) FailedToGet(resource string, cause error) *AppError {
	message := fmt.Sprintf("Failed to get %s", resource)
	return NewAppError(ErrCodeOperationFailed, message, SeverityMedium, 500).
		WithCause(cause).WithStackTrace()
}

// FailedToCreate 创建"创建失败"类型的错误
func (f *ErrorFactory) FailedToCreate(resource string, cause error) *AppError {
	message := fmt.Sprintf("Failed to create %s", resource)
	return NewAppError(ErrCodeOperationFailed, message, SeverityMedium, 500).
		WithCause(cause).WithStackTrace()
}

// FailedToUpdate 创建"更新失败"类型的错误
func (f *ErrorFactory) FailedToUpdate(resource string, cause error) *AppError {
	message := fmt.Sprintf("Failed to update %s", resource)
	return NewAppError(ErrCodeOperationFailed, message, SeverityMedium, 500).
		WithCause(cause).WithStackTrace()
}

// FailedToDelete 创建"删除失败"类型的错误
func (f *ErrorFactory) FailedToDelete(resource string, cause error) *AppError {
	message := fmt.Sprintf("Failed to delete %s", resource)
	return NewAppError(ErrCodeOperationFailed, message, SeverityMedium, 500).
		WithCause(cause).WithStackTrace()
}

// FailedToValidate 创建"验证失败"类型的错误
func (f *ErrorFactory) FailedToValidate(resource string, cause error) *AppError {
	message := fmt.Sprintf("Failed to validate %s", resource)
	return NewAppError(ErrCodeValidationFailed, message, SeverityLow, 400).
		WithCause(cause).WithStackTrace()
}

// FailedToConnect 创建"连接失败"类型的错误
func (f *ErrorFactory) FailedToConnect(service string, cause error) *AppError {
	message := fmt.Sprintf("Failed to connect to %s", service)
	return NewAppError(ErrCodeNetworkError, message, SeverityHigh, 503).
		WithCause(cause).WithStackTrace()
}

// FailedToInitialize 创建"初始化失败"类型的错误
func (f *ErrorFactory) FailedToInitialize(component string, cause error) *AppError {
	message := fmt.Sprintf("Failed to initialize %s", component)
	return NewAppError(ErrCodeInternal, message, SeverityHigh, 500).
		WithCause(cause).WithStackTrace()
}

// FailedToExecute 创建"执行失败"类型的错误
func (f *ErrorFactory) FailedToExecute(operation string, cause error) *AppError {
	message := fmt.Sprintf("Failed to execute %s", operation)
	return NewAppError(ErrCodeOperationFailed, message, SeverityMedium, 500).
		WithCause(cause).WithStackTrace()
}

// FailedToEncode 创建"编码失败"类型的错误
func (f *ErrorFactory) FailedToEncode(format string, cause error) *AppError {
	message := fmt.Sprintf("Failed to encode to %s", format)
	return NewAppError(ErrCodeInternal, message, SeverityMedium, 500).
		WithCause(cause).WithStackTrace()
}

// FailedToDecode 创建"解码失败"类型的错误
func (f *ErrorFactory) FailedToDecode(format string, cause error) *AppError {
	message := fmt.Sprintf("Failed to decode from %s", format)
	return NewAppError(ErrCodeInternal, message, SeverityMedium, 500).
		WithCause(cause).WithStackTrace()
}

// FailedToEncrypt 创建"加密失败"类型的错误
func (f *ErrorFactory) FailedToEncrypt(cause error) *AppError {
	return NewAppError(ErrCodeInternal, "Failed to encrypt data", SeverityHigh, 500).
		WithCause(cause).WithStackTrace()
}

// FailedToDecrypt 创建"解密失败"类型的错误
func (f *ErrorFactory) FailedToDecrypt(cause error) *AppError {
	return NewAppError(ErrCodeInternal, "Failed to decrypt data", SeverityHigh, 500).
		WithCause(cause).WithStackTrace()
}

// APIValidationFailed 创建"API验证失败"类型的错误
func (f *ErrorFactory) APIValidationFailed(provider string, cause error) *AppError {
	message := fmt.Sprintf("%s API validation failed", provider)
	return NewAppError(ErrCodeUnauthorized, message, SeverityMedium, 401).
		WithCause(cause).WithStackTrace()
}

// ProviderChatFailed 创建"提供商聊天失败"类型的错误
func (f *ErrorFactory) ProviderChatFailed(provider string, cause error) *AppError {
	message := fmt.Sprintf("%s chat completion failed", provider)
	return NewAppError(ErrCodeExternalService, message, SeverityMedium, 502).
		WithCause(cause).WithStackTrace()
}

// ToolExecutionFailed 创建"工具执行失败"类型的错误
func (f *ErrorFactory) ToolExecutionFailed(toolName string, attempts int, cause error) *AppError {
	message := fmt.Sprintf("Tool execution failed after %d attempts", attempts)
	return NewAppError(ErrCodeOperationFailed, message, SeverityMedium, 500).
		WithDetails(fmt.Sprintf("Tool: %s", toolName)).
		WithCause(cause).WithStackTrace()
}

// DatabaseOperationFailed 创建"数据库操作失败"类型的错误
func (f *ErrorFactory) DatabaseOperationFailed(operation string, cause error) *AppError {
	return NewDatabaseError(operation, cause)
}

// 全局错误工厂实例
var Factory = NewErrorFactory()

// 便捷函数，直接使用全局工厂实例
func FailedToOperation(operation string, cause error) *AppError {
	return Factory.FailedToOperation(operation, cause)
}

func FailedToGet(resource string, cause error) *AppError {
	return Factory.FailedToGet(resource, cause)
}

func FailedToCreate(resource string, cause error) *AppError {
	return Factory.FailedToCreate(resource, cause)
}

func FailedToUpdate(resource string, cause error) *AppError {
	return Factory.FailedToUpdate(resource, cause)
}

func FailedToDelete(resource string, cause error) *AppError {
	return Factory.FailedToDelete(resource, cause)
}

func FailedToValidate(resource string, cause error) *AppError {
	return Factory.FailedToValidate(resource, cause)
}

func FailedToConnect(service string, cause error) *AppError {
	return Factory.FailedToConnect(service, cause)
}

func FailedToInitialize(component string, cause error) *AppError {
	return Factory.FailedToInitialize(component, cause)
}

func FailedToExecute(operation string, cause error) *AppError {
	return Factory.FailedToExecute(operation, cause)
}

func FailedToEncode(format string, cause error) *AppError {
	return Factory.FailedToEncode(format, cause)
}

func FailedToDecode(format string, cause error) *AppError {
	return Factory.FailedToDecode(format, cause)
}

func FailedToEncrypt(cause error) *AppError {
	return Factory.FailedToEncrypt(cause)
}

func FailedToDecrypt(cause error) *AppError {
	return Factory.FailedToDecrypt(cause)
}

func APIValidationFailed(provider string, cause error) *AppError {
	return Factory.APIValidationFailed(provider, cause)
}

func ProviderChatFailed(provider string, cause error) *AppError {
	return Factory.ProviderChatFailed(provider, cause)
}

func ToolExecutionFailed(toolName string, attempts int, cause error) *AppError {
	return Factory.ToolExecutionFailed(toolName, attempts, cause)
}

func DatabaseOperationFailed(operation string, cause error) *AppError {
	return Factory.DatabaseOperationFailed(operation, cause)
}