package i18n

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var localeFS embed.FS

// Manager 国际化管理器
type Manager struct {
	bundle       *i18n.Bundle
	localizers   map[string]*i18n.Localizer
	defaultLang  string
	supportedLangs []string
}

// NewManager 创建国际化管理器
func NewManager(defaultLang string, supportedLangs []string) (*Manager, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	manager := &Manager{
		bundle:         bundle,
		localizers:     make(map[string]*i18n.Localizer),
		defaultLang:    defaultLang,
		supportedLangs: supportedLangs,
	}

	// 加载翻译文件
	if err := manager.loadTranslations(); err != nil {
		return nil, fmt.Errorf("failed to load translations: %w", err)
	}

	// 创建本地化器
	manager.createLocalizers()

	return manager, nil
}

// loadTranslations 加载翻译文件
func (m *Manager) loadTranslations() error {
	for _, lang := range m.supportedLangs {
		filename := fmt.Sprintf("locales/%s.json", lang)
		data, err := localeFS.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed to read translation file %s: %w", filename, err)
		}

		if _, err := m.bundle.ParseMessageFileBytes(data, filename); err != nil {
			return fmt.Errorf("failed to parse translation file %s: %w", filename, err)
		}
	}
	return nil
}

// createLocalizers 创建本地化器
func (m *Manager) createLocalizers() {
	for _, lang := range m.supportedLangs {
		m.localizers[lang] = i18n.NewLocalizer(m.bundle, lang)
	}
}

// GetLocalizer 获取指定语言的本地化器
func (m *Manager) GetLocalizer(lang string) *i18n.Localizer {
	if localizer, exists := m.localizers[lang]; exists {
		return localizer
	}
	// 返回默认语言的本地化器
	return m.localizers[m.defaultLang]
}

// T 翻译函数
func (m *Manager) T(lang, messageID string, templateData map[string]interface{}) string {
	localizer := m.GetLocalizer(lang)
	
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})
	if err != nil {
		// 如果翻译失败，返回消息ID
		return messageID
	}
	return msg
}

// TWithDefault 带默认值的翻译函数
func (m *Manager) TWithDefault(lang, messageID, defaultMessage string, templateData map[string]interface{}) string {
	localizer := m.GetLocalizer(lang)
	
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:      messageID,
		DefaultMessage: &i18n.Message{ID: messageID, Other: defaultMessage},
		TemplateData:   templateData,
	})
	if err != nil {
		return defaultMessage
	}
	return msg
}

// GetLanguageFromContext 从上下文获取语言
func (m *Manager) GetLanguageFromContext(c *gin.Context) string {
	// 1. 从查询参数获取
	if lang := c.Query("lang"); lang != "" && m.isSupportedLanguage(lang) {
		return lang
	}

	// 2. 从请求头获取
	if lang := c.GetHeader("Accept-Language"); lang != "" {
		if parsedLang := m.parseAcceptLanguage(lang); parsedLang != "" {
			return parsedLang
		}
	}

	// 3. 从Cookie获取
	if lang, err := c.Cookie("language"); err == nil && m.isSupportedLanguage(lang) {
		return lang
	}

	// 4. 返回默认语言
	return m.defaultLang
}

// GetLanguageFromRequest 从请求获取语言（非Gin上下文）
func (m *Manager) GetLanguageFromRequest(ctx context.Context, acceptLanguage string) string {
	if acceptLanguage != "" {
		if parsedLang := m.parseAcceptLanguage(acceptLanguage); parsedLang != "" {
			return parsedLang
		}
	}
	return m.defaultLang
}

// isSupportedLanguage 检查是否支持的语言
func (m *Manager) isSupportedLanguage(lang string) bool {
	for _, supported := range m.supportedLangs {
		if supported == lang {
			return true
		}
	}
	return false
}

// parseAcceptLanguage 解析Accept-Language头
func (m *Manager) parseAcceptLanguage(acceptLang string) string {
	// 简单解析Accept-Language头，取第一个支持的语言
	langs := strings.Split(acceptLang, ",")
	for _, lang := range langs {
		// 移除权重信息 (如 "en-US;q=0.9")
		lang = strings.TrimSpace(strings.Split(lang, ";")[0])
		
		// 检查完整匹配
		if m.isSupportedLanguage(lang) {
			return lang
		}
		
		// 检查语言前缀匹配 (如 "en-US" -> "en")
		if parts := strings.Split(lang, "-"); len(parts) > 1 {
			if m.isSupportedLanguage(parts[0]) {
				return parts[0]
			}
		}
	}
	return ""
}

// GetSupportedLanguages 获取支持的语言列表
func (m *Manager) GetSupportedLanguages() []string {
	return m.supportedLangs
}

// GetDefaultLanguage 获取默认语言
func (m *Manager) GetDefaultLanguage() string {
	return m.defaultLang
}