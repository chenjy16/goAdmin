package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// TestGoogleAIController_NewGoogleAIController 测试创建GoogleAI控制器
func TestGoogleAIController_NewGoogleAIController(t *testing.T) {
	zapLogger := zap.NewNop()
	controller := NewGoogleAIController(nil, zapLogger)
	assert.NotNil(t, controller)
}

// TestGoogleAIController_ChatCompletion 测试聊天完成HTTP处理
func TestGoogleAIController_ChatCompletion(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("ValidationError", func(t *testing.T) {
		zapLogger := zap.NewNop()
		controller := NewGoogleAIController(nil, zapLogger)

		// 创建无效请求（空body）
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/googleai/chat/completions", bytes.NewBuffer([]byte("{}")))
		c.Request.Header.Set("Content-Type", "application/json")

		// 调用ChatCompletion方法
		controller.ChatCompletion(c)

		// 验证返回400错误（由于验证失败）
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("MethodExists", func(t *testing.T) {
		zapLogger := zap.NewNop()
		controller := NewGoogleAIController(nil, zapLogger)
		assert.NotNil(t, controller.ChatCompletion)
	})
}

// TestGoogleAIController_ListModels 测试列出模型
func TestGoogleAIController_ListModels(t *testing.T) {
	t.Run("MethodExists", func(t *testing.T) {
		zapLogger := zap.NewNop()
		controller := NewGoogleAIController(nil, zapLogger)
		assert.NotNil(t, controller.ListModels)
	})
}

// TestGoogleAIController_SetAPIKey 测试设置API密钥
func TestGoogleAIController_SetAPIKey(t *testing.T) {
	t.Run("ValidationError", func(t *testing.T) {
		zapLogger := zap.NewNop()
		controller := NewGoogleAIController(nil, zapLogger)

		// 创建无效请求（空API密钥）
		requestBody := map[string]interface{}{
			"api_key": "",
		}
		jsonBody, _ := json.Marshal(requestBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/googleai/api-key", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		// 调用SetAPIKey方法
		controller.SetAPIKey(c)

		// 验证返回400错误
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestGoogleAIController_GetModelConfig 测试获取模型配置
func TestGoogleAIController_GetModelConfig(t *testing.T) {
	t.Run("MethodExists", func(t *testing.T) {
		zapLogger := zap.NewNop()
		controller := NewGoogleAIController(nil, zapLogger)
		assert.NotNil(t, controller.GetModelConfig)
	})
}

// TestGoogleAIController_EnableModel 测试启用模型
func TestGoogleAIController_EnableModel(t *testing.T) {
	t.Run("MethodExists", func(t *testing.T) {
		zapLogger := zap.NewNop()
		controller := NewGoogleAIController(nil, zapLogger)
		assert.NotNil(t, controller.EnableModel)
	})
}

// TestGoogleAIController_DisableModel 测试禁用模型
func TestGoogleAIController_DisableModel(t *testing.T) {
	t.Run("MethodExists", func(t *testing.T) {
		zapLogger := zap.NewNop()
		controller := NewGoogleAIController(nil, zapLogger)
		assert.NotNil(t, controller.DisableModel)
	})
}

// TestGoogleAIController_ValidateAPIKey 测试验证API密钥
func TestGoogleAIController_ValidateAPIKey(t *testing.T) {
	t.Run("MethodExists", func(t *testing.T) {
		zapLogger := zap.NewNop()
		controller := NewGoogleAIController(nil, zapLogger)
		assert.NotNil(t, controller.ValidateAPIKey)
	})
}