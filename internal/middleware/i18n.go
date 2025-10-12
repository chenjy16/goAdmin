package middleware

import (
	"go-springAi/internal/i18n"

	"github.com/gin-gonic/gin"
)

const (
	// LanguageContextKey 语言上下文键
	LanguageContextKey = "language"
	// I18nManagerContextKey 国际化管理器上下文键
	I18nManagerContextKey = "i18n_manager"
)

// I18nMiddleware 国际化中间件
func I18nMiddleware(manager *i18n.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求语言
		lang := manager.GetLanguageFromContext(c)
		
		// 设置语言到上下文
		c.Set(LanguageContextKey, lang)
		c.Set(I18nManagerContextKey, manager)
		
		// 创建翻译函数并添加到上下文
		translateFunc := func(messageID string, templateData map[string]interface{}) string {
			return manager.T(lang, messageID, templateData)
		}
		c.Set("i18n_translate", translateFunc)
		
		// 设置响应头
		c.Header("Content-Language", lang)
		
		c.Next()
	}
}

// GetLanguageFromContext 从上下文获取语言
func GetLanguageFromContext(c *gin.Context) string {
	if lang, exists := c.Get(LanguageContextKey); exists {
		if langStr, ok := lang.(string); ok {
			return langStr
		}
	}
	return "en" // 默认语言
}

// GetI18nManagerFromContext 从上下文获取国际化管理器
func GetI18nManagerFromContext(c *gin.Context) *i18n.Manager {
	if manager, exists := c.Get(I18nManagerContextKey); exists {
		if i18nManager, ok := manager.(*i18n.Manager); ok {
			return i18nManager
		}
	}
	return nil
}

// T 翻译函数（从上下文获取语言和管理器）
func T(c *gin.Context, messageID string, templateData map[string]interface{}) string {
	manager := GetI18nManagerFromContext(c)
	if manager == nil {
		return messageID
	}
	
	lang := GetLanguageFromContext(c)
	return manager.T(lang, messageID, templateData)
}

// TWithDefault 带默认值的翻译函数
func TWithDefault(c *gin.Context, messageID, defaultMessage string, templateData map[string]interface{}) string {
	manager := GetI18nManagerFromContext(c)
	if manager == nil {
		return defaultMessage
	}
	
	lang := GetLanguageFromContext(c)
	return manager.TWithDefault(lang, messageID, defaultMessage, templateData)
}