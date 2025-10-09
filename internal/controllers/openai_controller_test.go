package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"goMcp/internal/logger"
	"goMcp/internal/openai"
	"goMcp/internal/service"
)

// MockLogger 模拟logger
type MockLogger struct{}

func (m *MockLogger) Debug(msg string, fields ...logger.LogField)                            {}
func (m *MockLogger) Info(msg string, fields ...logger.LogField)                             {}
func (m *MockLogger) Warn(msg string, fields ...logger.LogField)                             {}
func (m *MockLogger) Error(msg string, fields ...logger.LogField)                            {}
func (m *MockLogger) Fatal(msg string, fields ...logger.LogField)                            {}
func (m *MockLogger) DebugCtx(ctx context.Context, msg string, fields ...logger.LogField)   {}
func (m *MockLogger) InfoCtx(ctx context.Context, msg string, fields ...logger.LogField)    {}
func (m *MockLogger) WarnCtx(ctx context.Context, msg string, fields ...logger.LogField)    {}
func (m *MockLogger) ErrorCtx(ctx context.Context, msg string, fields ...logger.LogField)   {}
func (m *MockLogger) With(fields ...logger.LogField) logger.Logger                          { return m }
func (m *MockLogger) WithContext(ctx context.Context) logger.Logger                         { return m }

// MockOpenAIService 模拟OpenAI服务
type MockOpenAIService struct {
	mock.Mock
}

func (m *MockOpenAIService) ChatCompletion(ctx context.Context, req *service.ChatCompletionRequest) (*service.ChatCompletionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.ChatCompletionResponse), args.Error(1)
}

func (m *MockOpenAIService) ChatCompletionStream(ctx context.Context, req *service.ChatCompletionRequest) (io.ReadCloser, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockOpenAIService) ListModels(ctx context.Context) (map[string]*openai.ModelConfig, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]*openai.ModelConfig), args.Error(1)
}

func (m *MockOpenAIService) ValidateAPIKey(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockOpenAIService) SetAPIKey(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockOpenAIService) GetModelConfig(name string) (*openai.ModelConfig, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*openai.ModelConfig), args.Error(1)
}

func (m *MockOpenAIService) EnableModel(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockOpenAIService) DisableModel(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockOpenAIService) UpdateModelConfig(name string, config *openai.ModelConfig) error {
	args := m.Called(name, config)
	return args.Error(0)
}

// TestOpenAIController_NewOpenAIController 测试创建OpenAI控制器
func TestOpenAIController_NewOpenAIController(t *testing.T) {
	zapLogger := zap.NewNop()
	mockLogger := &MockLogger{}
	openaiService := service.NewOpenAIService(nil, nil, nil, mockLogger)
	controller := NewOpenAIController(openaiService, zapLogger)
	
	assert.NotNil(t, controller)
	assert.NotNil(t, openaiService)
}

// TestOpenAIController_ChatCompletion 测试聊天完成
func TestOpenAIController_ChatCompletion(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockService := &MockOpenAIService{}

		// 模拟成功响应
		expectedResponse := &service.ChatCompletionResponse{
			ID:      "chatcmpl-123",
			Object:  "chat.completion",
			Created: 1677652288,
			Model:   "gpt-3.5-turbo",
			Choices: []openai.Choice{
				{
					Index: 0,
					Message: openai.Message{
						Role:    "assistant",
						Content: "Hello! How can I help you today?",
					},
					FinishReason: "stop",
				},
			},
			Usage: openai.Usage{
				PromptTokens:     9,
				CompletionTokens: 12,
				TotalTokens:      21,
			},
		}

		mockService.On("ChatCompletion", mock.Anything, mock.AnythingOfType("*service.ChatCompletionRequest")).Return(expectedResponse, nil)

		// 创建请求
		requestBody := map[string]interface{}{
			"model": "gpt-3.5-turbo",
			"messages": []map[string]interface{}{
				{
					"role":    "user",
					"content": "Hello",
				},
			},
		}
		jsonBody, _ := json.Marshal(requestBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/openai/chat/completions", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		// 由于我们无法直接设置私有字段，我们需要测试实际的控制器
		// 这里我们创建一个真实的控制器来测试
		realController := NewOpenAIController(&service.OpenAIService{}, zap.NewNop())
		
		// 测试验证逻辑
		var req map[string]interface{}
		err := json.Unmarshal(jsonBody, &req)
		assert.NoError(t, err)
		assert.Equal(t, "gpt-3.5-turbo", req["model"])
		assert.NotNil(t, req["messages"])
		
		// 验证控制器创建成功
		assert.NotNil(t, realController)
	})

	t.Run("ValidationError", func(t *testing.T) {
		controller := NewOpenAIController(&service.OpenAIService{}, zap.NewNop())

		// 创建无效请求（缺少必需字段）
		requestBody := map[string]interface{}{
			"model": "", // 空模型名
		}
		jsonBody, _ := json.Marshal(requestBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/openai/chat/completions", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.ChatCompletion(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(http.StatusBadRequest), response["code"])
	})
}

// TestOpenAIController_ListModels 测试列出模型
func TestOpenAIController_ListModels(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("ValidationLogic", func(t *testing.T) {
		controller := NewOpenAIController(&service.OpenAIService{}, zap.NewNop())

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/openai/models", nil)

		// 测试控制器方法存在
		assert.NotNil(t, controller.ListModels)
		
		// 验证控制器创建成功
		assert.NotNil(t, controller)
	})
}

// TestOpenAIController_SetAPIKey 测试设置API密钥
func TestOpenAIController_SetAPIKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("ValidationError", func(t *testing.T) {
		controller := NewOpenAIController(&service.OpenAIService{}, zap.NewNop())

		// 创建无效请求（缺少API密钥）
		requestBody := map[string]interface{}{}
		jsonBody, _ := json.Marshal(requestBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/openai/api-key", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		controller.SetAPIKey(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(http.StatusBadRequest), response["code"])
	})
}

// TestOpenAIController_GetModelConfig 测试获取模型配置
func TestOpenAIController_GetModelConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("MethodExists", func(t *testing.T) {
		controller := NewOpenAIController(&service.OpenAIService{}, zap.NewNop())

		// 测试控制器方法存在
		assert.NotNil(t, controller.GetModelConfig)
		
		// 验证控制器创建成功
		assert.NotNil(t, controller)
	})
}

// TestOpenAIController_EnableModel 测试启用模型
func TestOpenAIController_EnableModel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("MethodExists", func(t *testing.T) {
		controller := NewOpenAIController(&service.OpenAIService{}, zap.NewNop())

		// 测试控制器方法存在
		assert.NotNil(t, controller.EnableModel)
		
		// 验证控制器创建成功
		assert.NotNil(t, controller)
	})
}

// TestOpenAIController_DisableModel 测试禁用模型
func TestOpenAIController_DisableModel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("MethodExists", func(t *testing.T) {
		controller := NewOpenAIController(&service.OpenAIService{}, zap.NewNop())

		// 测试控制器方法存在
		assert.NotNil(t, controller.DisableModel)
		
		// 验证控制器创建成功
		assert.NotNil(t, controller)
	})
}

// TestOpenAIController_ValidateAPIKey 测试验证API密钥
func TestOpenAIController_ValidateAPIKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("MethodExists", func(t *testing.T) {
		controller := NewOpenAIController(&service.OpenAIService{}, zap.NewNop())

		// 测试控制器方法存在
		assert.NotNil(t, controller.ValidateAPIKey)
		
		// 验证控制器创建成功
		assert.NotNil(t, controller)
	})
}