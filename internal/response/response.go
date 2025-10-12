package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string, err string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Error:   err,
	})
}

// BadRequest 400错误
func BadRequest(c *gin.Context, message string, err string) {
	Error(c, http.StatusBadRequest, message, err)
}

// InternalServerError 500错误
func InternalServerError(c *gin.Context, message string, err string) {
	Error(c, http.StatusInternalServerError, message, err)
}

// 国际化响应函数

// I18nSuccess 国际化成功响应
func I18nSuccess(c *gin.Context, code int, messageID string, data interface{}, templateData map[string]interface{}) {
	message := getI18nMessage(c, messageID, templateData)
	Success(c, code, message, data)
}

// I18nError 国际化错误响应
func I18nError(c *gin.Context, code int, messageID string, err string, templateData map[string]interface{}) {
	message := getI18nMessage(c, messageID, templateData)
	Error(c, code, message, err)
}

// I18nBadRequest 国际化400错误
func I18nBadRequest(c *gin.Context, messageID string, err string, templateData map[string]interface{}) {
	I18nError(c, http.StatusBadRequest, messageID, err, templateData)
}

// I18nUnauthorized 国际化401错误
func I18nUnauthorized(c *gin.Context, messageID string, err string, templateData map[string]interface{}) {
	I18nError(c, http.StatusUnauthorized, messageID, err, templateData)
}

// I18nNotFound 国际化404错误
func I18nNotFound(c *gin.Context, messageID string, err string, templateData map[string]interface{}) {
	I18nError(c, http.StatusNotFound, messageID, err, templateData)
}

// I18nInternalServerError 国际化500错误
func I18nInternalServerError(c *gin.Context, messageID string, err string, templateData map[string]interface{}) {
	I18nError(c, http.StatusInternalServerError, messageID, err, templateData)
}

// getI18nMessage 获取国际化消息的辅助函数
func getI18nMessage(c *gin.Context, messageID string, templateData map[string]interface{}) string {
	// 尝试从上下文获取翻译函数
	if t, exists := c.Get("i18n_translate"); exists {
		if translateFunc, ok := t.(func(string, map[string]interface{}) string); ok {
			return translateFunc(messageID, templateData)
		}
	}
	
	// 如果没有翻译函数，返回消息ID作为默认值
	return messageID
}
