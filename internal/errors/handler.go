package errors

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ErrorHandler 统一的错误处理器
type ErrorHandler struct {
	i18nManager I18nManager
}

// I18nManager 国际化管理器接口
type I18nManager interface {
	GetErrorMessage(lang string, appErr *AppError) string
	T(lang string, key string, params map[string]interface{}) string
}

// NewErrorHandler 创建错误处理器
func NewErrorHandler(i18nManager I18nManager) *ErrorHandler {
	return &ErrorHandler{
		i18nManager: i18nManager,
	}
}

// HandleError 统一的错误处理方法
func (h *ErrorHandler) HandleError(c *gin.Context, err error) {
	// 获取语言设置
	lang := h.getLanguage(c)

	// 检查是否为应用程序错误
	if appErr, ok := IsAppError(err); ok {
		h.handleAppError(c, appErr, lang)
		return
	}

	// 检查是否为验证错误
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		h.handleValidationErrors(c, validationErrors, lang)
		return
	}

	// 处理其他类型的错误
	h.handleGenericError(c, err, lang)
}

// handleAppError 处理应用程序错误
func (h *ErrorHandler) handleAppError(c *gin.Context, appErr *AppError, lang string) {
	// 获取国际化消息
	message := appErr.Message
	if h.i18nManager != nil {
		message = h.i18nManager.GetErrorMessage(lang, appErr)
	}

	// 构建响应
	response := gin.H{
		"error": gin.H{
			"code":      appErr.Code,
			"message":   message,
			"timestamp": appErr.Timestamp,
		},
	}

	// 在开发环境下添加详细信息
	if gin.Mode() == gin.DebugMode {
		if appErr.Details != "" {
			response["error"].(gin.H)["details"] = appErr.Details
		}
		if len(appErr.StackTrace) > 0 {
			response["error"].(gin.H)["stack_trace"] = appErr.StackTrace
		}
	}

	c.JSON(appErr.HTTPStatus, response)
}

// handleValidationErrors 处理验证错误
func (h *ErrorHandler) handleValidationErrors(c *gin.Context, validationErrors validator.ValidationErrors, lang string) {
	var errorMessages []string
	
	for _, e := range validationErrors {
		message := h.getValidationErrorMessage(e, lang)
		errorMessages = append(errorMessages, message)
	}

	appErr := NewValidationError("Validation failed")
	appErr.Details = fmt.Sprintf("Fields: %v", errorMessages)

	h.handleAppError(c, appErr, lang)
}

// handleGenericError 处理通用错误
func (h *ErrorHandler) handleGenericError(c *gin.Context, err error, lang string) {
	appErr := NewInternalError(err.Error()).WithStackTrace()
	h.handleAppError(c, appErr, lang)
}

// getLanguage 获取请求的语言设置
func (h *ErrorHandler) getLanguage(c *gin.Context) string {
	// 优先从查询参数获取
	if lang := c.Query("lang"); lang != "" {
		return lang
	}

	// 从Header获取
	if lang := c.GetHeader("Accept-Language"); lang != "" {
		// 简单解析，取第一个语言
		if len(lang) >= 2 {
			return lang[:2]
		}
	}

	// 默认语言
	return "zh"
}

// getValidationErrorMessage 获取验证错误的消息
func (h *ErrorHandler) getValidationErrorMessage(e validator.FieldError, lang string) string {
	field := e.Field()
	tag := e.Tag()

	// 如果有国际化管理器，使用国际化消息
	if h.i18nManager != nil {
		key := fmt.Sprintf("validation.%s", tag)
		message := h.i18nManager.T(lang, key, map[string]interface{}{
			"Field": field,
			"Value": e.Value(),
			"Param": e.Param(),
		})
		if message != key { // 如果找到了翻译
			return message
		}
	}

	// 默认英文消息
	switch tag {
	case "required":
		return fmt.Sprintf("Field '%s' is required", field)
	case "email":
		return fmt.Sprintf("Field '%s' must be a valid email", field)
	case "min":
		return fmt.Sprintf("Field '%s' must be at least %s characters", field, e.Param())
	case "max":
		return fmt.Sprintf("Field '%s' must be at most %s characters", field, e.Param())
	case "len":
		return fmt.Sprintf("Field '%s' must be exactly %s characters", field, e.Param())
	case "numeric":
		return fmt.Sprintf("Field '%s' must be numeric", field)
	case "alpha":
		return fmt.Sprintf("Field '%s' must contain only letters", field)
	case "alphanum":
		return fmt.Sprintf("Field '%s' must contain only letters and numbers", field)
	default:
		return fmt.Sprintf("Field '%s' is invalid", field)
	}
}

// HandlePanic 处理panic恢复
func (h *ErrorHandler) HandlePanic(c *gin.Context, recovered interface{}) {
	lang := h.getLanguage(c)
	
	appErr := NewInternalError(fmt.Sprintf("Panic recovered: %v", recovered)).
		WithStackTrace()
	
	h.handleAppError(c, appErr, lang)
}

// ErrorMiddleware 错误处理中间件
func (h *ErrorHandler) ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 处理panic
		defer func() {
			if recovered := recover(); recovered != nil {
				h.HandlePanic(c, recovered)
				c.Abort()
			}
		}()

		c.Next()

		// 处理在处理过程中产生的错误
		if len(c.Errors) > 0 {
			// 只处理最后一个错误
			err := c.Errors.Last().Err
			h.HandleError(c, err)
		}
	}
}