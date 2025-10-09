package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-springAi/internal/logger"
	"go-springAi/internal/provider"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupAIController() *AIController {
	gin.SetMode(gin.TestMode)
	zapLogger, _ := zap.NewDevelopment()
	
	// 创建logger实例
	appLogger := logger.NewLoggerFromZap(zapLogger)
	
	// 创建provider manager
	providerManager := provider.NewManager(appLogger)
	
	// 创建controller
	controller := NewAIController(providerManager, zapLogger)
	
	return controller
}

func TestAIController_NewAIController(t *testing.T) {
	zapLogger, _ := zap.NewDevelopment()
	appLogger := logger.NewLoggerFromZap(zapLogger)
	providerManager := provider.NewManager(appLogger)
	
	controller := NewAIController(providerManager, zapLogger)
	assert.NotNil(t, controller)
	assert.NotNil(t, controller.providerManager)
	assert.NotNil(t, controller.logger)
}



func TestAIController_ListModels_InvalidProvider(t *testing.T) {
	controller := setupAIController()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/ai/invalid/models", nil)
	c.Params = gin.Params{
		{Key: "provider", Value: "invalid"},
	}

	// 调用控制器方法
	controller.ListModels(c)

	// 验证结果 - 应该返回400错误
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(http.StatusBadRequest), response["code"])
	assert.Contains(t, response["message"], "Invalid provider")
}

func TestAIController_ListProviders(t *testing.T) {
	controller := setupAIController()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/ai/providers", nil)

	// 调用控制器方法
	controller.ListProviders(c)

	// 验证结果
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(http.StatusOK), response["code"])
	assert.Contains(t, response, "data")
	
	data := response["data"].(map[string]interface{})
	assert.Contains(t, data, "providers")
	
	// providers可能为nil或空列表
	if data["providers"] != nil {
		providers := data["providers"].([]interface{})
		assert.Len(t, providers, 0)
	}
}

func TestAIController_GetModelConfig_InvalidProvider(t *testing.T) {
	controller := setupAIController()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/ai/invalid/models/gpt-3.5-turbo", nil)
	c.Params = gin.Params{
		{Key: "provider", Value: "invalid"},
		{Key: "model", Value: "gpt-3.5-turbo"},
	}

	// 调用控制器方法
	controller.GetModelConfig(c)

	// 验证结果 - 应该返回400错误
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(http.StatusBadRequest), response["code"])
	assert.Contains(t, response["message"], "Invalid provider")
}

func TestAIController_EnableModel_InvalidProvider(t *testing.T) {
	controller := setupAIController()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/ai/invalid/models/gpt-3.5-turbo/enable", nil)
	c.Params = gin.Params{
		{Key: "provider", Value: "invalid"},
		{Key: "model", Value: "gpt-3.5-turbo"},
	}

	// 调用控制器方法
	controller.EnableModel(c)

	// 验证结果 - 应该返回400错误
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(http.StatusBadRequest), response["code"])
	assert.Contains(t, response["message"], "Invalid provider")
}

func TestAIController_DisableModel_InvalidProvider(t *testing.T) {
	controller := setupAIController()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/ai/invalid/models/gpt-3.5-turbo/disable", nil)
	c.Params = gin.Params{
		{Key: "provider", Value: "invalid"},
		{Key: "model", Value: "gpt-3.5-turbo"},
	}

	// 调用控制器方法
	controller.DisableModel(c)

	// 验证结果 - 应该返回400错误
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(http.StatusBadRequest), response["code"])
	assert.Contains(t, response["message"], "Invalid provider")
}

func TestAIController_ValidateAPIKey_InvalidProvider(t *testing.T) {
	controller := setupAIController()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/ai/invalid/validate", nil)
	c.Params = gin.Params{
		{Key: "provider", Value: "invalid"},
	}

	// 调用控制器方法
	controller.ValidateAPIKey(c)

	// 验证结果 - 应该返回400错误
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(http.StatusBadRequest), response["code"])
	assert.Contains(t, response["message"], "Invalid provider")
}

func TestAIController_SetAPIKey_InvalidProvider(t *testing.T) {
	controller := setupAIController()

	requestBody := map[string]string{
		"api_key": "test-key",
	}

	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/ai/invalid/api-key", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "provider", Value: "invalid"},
	}

	// 调用控制器方法
	controller.SetAPIKey(c)

	// 验证结果 - 应该返回400错误
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(http.StatusBadRequest), response["code"])
	assert.Contains(t, response["message"], "Invalid provider")
}