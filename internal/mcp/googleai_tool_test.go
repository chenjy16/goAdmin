package mcp

import (
	"context"
	"errors"
	"testing"

	"goMcp/internal/googleai"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestGoogleAIChatTool 测试Google AI聊天工具
func TestGoogleAIChatTool(t *testing.T) {
	t.Run("NewGoogleAIChatTool", func(t *testing.T) {
		mockService := &MockGoogleAIService{}
		tool := NewGoogleAIChatTool(mockService)
		
		assert.NotNil(t, tool)
		assert.Equal(t, "googleai_chat", tool.Name)
		assert.Equal(t, "Chat with Google AI models (Gemini Pro, Gemini Pro Vision, etc.)", tool.Description)
		assert.Equal(t, mockService, tool.googleaiService)
	})

	t.Run("Validate_Success", func(t *testing.T) {
		mockService := &MockGoogleAIService{}
		tool := NewGoogleAIChatTool(mockService)
		
		args := map[string]interface{}{
			"model": "gemini-pro",
			"messages": []interface{}{
				map[string]interface{}{
					"role":    "user",
					"content": "Hello",
				},
			},
		}
		
		err := tool.Validate(args)
		assert.NoError(t, err)
	})

	t.Run("Validate_MissingMessages", func(t *testing.T) {
		mockService := &MockGoogleAIService{}
		tool := NewGoogleAIChatTool(mockService)
		
		args := map[string]interface{}{
			"model": "gemini-pro",
		}
		
		err := tool.Validate(args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "messages parameter is required")
	})

	t.Run("Validate_EmptyMessages", func(t *testing.T) {
		mockService := &MockGoogleAIService{}
		tool := NewGoogleAIChatTool(mockService)
		
		args := map[string]interface{}{
			"model":    "gemini-pro",
			"messages": []interface{}{},
		}
		
		err := tool.Validate(args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "messages array cannot be empty")
	})

	t.Run("Validate_InvalidMessageFormat", func(t *testing.T) {
		mockService := &MockGoogleAIService{}
		tool := NewGoogleAIChatTool(mockService)
		
		args := map[string]interface{}{
			"model": "gemini-pro",
			"messages": []interface{}{
				"invalid message format",
			},
		}
		
		err := tool.Validate(args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "message at index 0 must be an object")
	})

	t.Run("Validate_InvalidRole", func(t *testing.T) {
		mockService := &MockGoogleAIService{}
		tool := NewGoogleAIChatTool(mockService)
		
		args := map[string]interface{}{
			"model": "gemini-pro",
			"messages": []interface{}{
				map[string]interface{}{
					"role":    "invalid_role",
					"content": "Hello",
				},
			},
		}
		
		err := tool.Validate(args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid role")
	})

	t.Run("Validate_InvalidTemperature", func(t *testing.T) {
		mockService := &MockGoogleAIService{}
		tool := NewGoogleAIChatTool(mockService)
		
		args := map[string]interface{}{
			"model": "gemini-pro",
			"messages": []interface{}{
				map[string]interface{}{
					"role":    "user",
					"content": "Hello",
				},
			},
			"temperature": 3.0, // Invalid: > 2
		}
		
		err := tool.Validate(args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "temperature must be between 0 and 2")
	})

	t.Run("Execute_Success", func(t *testing.T) {
		mockService := &MockGoogleAIService{}
		tool := NewGoogleAIChatTool(mockService)
		ctx := context.Background()
		
		expectedResponse := &GoogleAIChatCompletionResponse{
			ID:      "test-id",
			Object:  "chat.completion",
			Created: 1234567890,
			Model:   "gemini-pro",
			Choices: []googleai.Choice{
				{
					Index: 0,
					Message: googleai.Message{
						Role:    "model",
						Content: "Hello! How can I help you?",
					},
					FinishReason: "stop",
				},
			},
			Usage: googleai.Usage{
				PromptTokens:     10,
				CompletionTokens: 8,
				TotalTokens:      18,
			},
		}
		
		mockService.On("ChatCompletion", ctx, mock.AnythingOfType("*mcp.GoogleAIChatCompletionRequest")).Return(expectedResponse, nil)
		
		args := map[string]interface{}{
			"model": "gemini-pro",
			"messages": []interface{}{
				map[string]interface{}{
					"role":    "user",
					"content": "Hello",
				},
			},
		}
		
		response, err := tool.Execute(ctx, args)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.False(t, response.IsError)
		assert.Equal(t, 1, len(response.Content))
		assert.Equal(t, "text", response.Content[0].Type)
		assert.NotNil(t, response.Content[0].Data)
		
		mockService.AssertExpectations(t)
	})

	t.Run("Execute_ServiceError", func(t *testing.T) {
		mockService := &MockGoogleAIService{}
		tool := NewGoogleAIChatTool(mockService)
		ctx := context.Background()
		
		mockService.On("ChatCompletion", ctx, mock.AnythingOfType("*mcp.GoogleAIChatCompletionRequest")).Return(nil, errors.New("API error"))
		
		args := map[string]interface{}{
			"model": "gemini-pro",
			"messages": []interface{}{
				map[string]interface{}{
					"role":    "user",
					"content": "Hello",
				},
			},
		}
		
		response, err := tool.Execute(ctx, args)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.IsError)
		assert.Contains(t, response.Content[0].Text, "Google AI API error")
		
		mockService.AssertExpectations(t)
	})
}

// TestGoogleAIModelsTool 测试Google AI模型工具
func TestGoogleAIModelsTool(t *testing.T) {
	t.Run("NewGoogleAIModelsTool", func(t *testing.T) {
		mockService := &MockGoogleAIService{}
		tool := NewGoogleAIModelsTool(mockService)
		
		assert.NotNil(t, tool)
		assert.Equal(t, "googleai_models", tool.Name)
		assert.Equal(t, mockService, tool.googleaiService)
	})

	t.Run("Execute_Success", func(t *testing.T) {
		mockService := &MockGoogleAIService{}
		tool := NewGoogleAIModelsTool(mockService)
		ctx := context.Background()
		
		expectedModels := map[string]*googleai.ModelConfig{
			"gemini-pro": {
				Name:        "gemini-pro",
				DisplayName: "Gemini Pro",
				Enabled:     true,
				MaxTokens:   8192,
			},
			"gemini-pro-vision": {
				Name:        "gemini-pro-vision",
				DisplayName: "Gemini Pro Vision",
				Enabled:     false,
				MaxTokens:   4096,
			},
		}
		
		mockService.On("ListModels", ctx).Return(expectedModels, nil)
		
		args := map[string]interface{}{}
		
		response, err := tool.Execute(ctx, args)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.False(t, response.IsError)
		assert.Equal(t, 1, len(response.Content))
		assert.NotNil(t, response.Content[0].Data)
		
		mockService.AssertExpectations(t)
	})

	t.Run("Execute_ServiceError", func(t *testing.T) {
		mockService := &MockGoogleAIService{}
		tool := NewGoogleAIModelsTool(mockService)
		ctx := context.Background()
		
		mockService.On("ListModels", ctx).Return(nil, errors.New("failed to list models"))
		
		args := map[string]interface{}{}
		
		response, err := tool.Execute(ctx, args)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.IsError)
		assert.Contains(t, response.Content[0].Text, "Failed to list models")
		
		mockService.AssertExpectations(t)
	})
}

// TestGoogleAIConfigTool 测试Google AI配置工具
func TestGoogleAIConfigTool(t *testing.T) {
	t.Run("NewGoogleAIConfigTool", func(t *testing.T) {
		mockService := &MockGoogleAIService{}
		tool := NewGoogleAIConfigTool(mockService)
		
		assert.NotNil(t, tool)
		assert.Equal(t, "googleai_config", tool.Name)
		assert.Equal(t, mockService, tool.googleaiService)
	})

	t.Run("Validate_Success", func(t *testing.T) {
		mockService := &MockGoogleAIService{}
		tool := NewGoogleAIConfigTool(mockService)
		
		args := map[string]interface{}{
			"action": "set_api_key",
			"api_key": "test-api-key",
		}
		
		err := tool.Validate(args)
		assert.NoError(t, err)
	})

	t.Run("Validate_MissingAction", func(t *testing.T) {
		mockService := &MockGoogleAIService{}
		tool := NewGoogleAIConfigTool(mockService)
		
		args := map[string]interface{}{}
		
		err := tool.Validate(args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "action parameter is required")
	})

	t.Run("Validate_InvalidAction", func(t *testing.T) {
		mockService := &MockGoogleAIService{}
		tool := NewGoogleAIConfigTool(mockService)
		
		args := map[string]interface{}{
			"action": "invalid_action",
		}
		
		err := tool.Validate(args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid action")
	})

	t.Run("Execute_SetAPIKey_Success", func(t *testing.T) {
		mockService := &MockGoogleAIService{}
		tool := NewGoogleAIConfigTool(mockService)
		ctx := context.Background()
		
		mockService.On("SetAPIKey", "test-api-key").Return(nil)
		
		args := map[string]interface{}{
			"action":  "set_api_key",
			"api_key": "test-api-key",
		}
		
		response, err := tool.Execute(ctx, args)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.False(t, response.IsError)
		assert.NotNil(t, response.Content[0].Data)
		
		mockService.AssertExpectations(t)
	})

	t.Run("Execute_ValidateAPIKey_Success", func(t *testing.T) {
		mockService := &MockGoogleAIService{}
		tool := NewGoogleAIConfigTool(mockService)
		ctx := context.Background()
		
		mockService.On("ValidateAPIKey", ctx).Return(nil)
		
		args := map[string]interface{}{
			"action": "validate_api_key",
		}
		
		response, err := tool.Execute(ctx, args)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.False(t, response.IsError)
		assert.NotNil(t, response.Content[0].Data)
		
		mockService.AssertExpectations(t)
	})

	t.Run("Execute_EnableModel_Success", func(t *testing.T) {
		mockService := &MockGoogleAIService{}
		tool := NewGoogleAIConfigTool(mockService)
		ctx := context.Background()
		
		mockService.On("EnableModel", "gemini-pro").Return(nil)
		
		args := map[string]interface{}{
			"action": "enable_model",
			"model":  "gemini-pro",
		}
		
		response, err := tool.Execute(ctx, args)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.False(t, response.IsError)
		assert.NotNil(t, response.Content[0].Data)
		
		mockService.AssertExpectations(t)
	})

	t.Run("Execute_DisableModel_Success", func(t *testing.T) {
		mockService := &MockGoogleAIService{}
		tool := NewGoogleAIConfigTool(mockService)
		ctx := context.Background()
		
		mockService.On("DisableModel", "gemini-pro").Return(nil)
		
		args := map[string]interface{}{
			"action": "disable_model",
			"model":  "gemini-pro",
		}
		
		response, err := tool.Execute(ctx, args)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.False(t, response.IsError)
		assert.NotNil(t, response.Content[0].Data)
		
		mockService.AssertExpectations(t)
	})

	t.Run("Execute_ServiceError", func(t *testing.T) {
		mockService := &MockGoogleAIService{}
		tool := NewGoogleAIConfigTool(mockService)
		ctx := context.Background()
		
		mockService.On("SetAPIKey", "invalid-key").Return(errors.New("invalid API key"))
		
		args := map[string]interface{}{
			"action":  "set_api_key",
			"api_key": "invalid-key",
		}
		
		response, err := tool.Execute(ctx, args)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.IsError)
		assert.Contains(t, response.Content[0].Text, "Failed to set API key")
		
		mockService.AssertExpectations(t)
	})
}