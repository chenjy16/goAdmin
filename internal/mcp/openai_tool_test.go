package mcp

import (
	"context"
	"errors"
	"testing"

	"admin/internal/openai"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestOpenAIChatTool 测试OpenAI聊天工具
func TestOpenAIChatTool(t *testing.T) {
	t.Run("NewOpenAIChatTool", func(t *testing.T) {
		mockService := &MockOpenAIService{}
		tool := NewOpenAIChatTool(mockService)
		
		assert.NotNil(t, tool)
		assert.Equal(t, "openai_chat", tool.Name)
		assert.Equal(t, "Chat with OpenAI models (GPT-4, GPT-3.5-turbo, etc.)", tool.Description)
		assert.Equal(t, mockService, tool.openaiService)
	})

	t.Run("Validate_Success", func(t *testing.T) {
		mockService := &MockOpenAIService{}
		tool := NewOpenAIChatTool(mockService)
		
		args := map[string]interface{}{
			"model": "gpt-4",
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
		mockService := &MockOpenAIService{}
		tool := NewOpenAIChatTool(mockService)
		
		args := map[string]interface{}{
			"model": "gpt-4",
		}
		
		err := tool.Validate(args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "messages parameter is required")
	})

	t.Run("Validate_EmptyMessages", func(t *testing.T) {
		mockService := &MockOpenAIService{}
		tool := NewOpenAIChatTool(mockService)
		
		args := map[string]interface{}{
			"model":    "gpt-4",
			"messages": []interface{}{},
		}
		
		err := tool.Validate(args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "messages array cannot be empty")
	})

	t.Run("Validate_InvalidMessageFormat", func(t *testing.T) {
		mockService := &MockOpenAIService{}
		tool := NewOpenAIChatTool(mockService)
		
		args := map[string]interface{}{
			"model": "gpt-4",
			"messages": []interface{}{
				"invalid message format",
			},
		}
		
		err := tool.Validate(args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "message at index 0 must be an object")
	})

	t.Run("Validate_InvalidRole", func(t *testing.T) {
		mockService := &MockOpenAIService{}
		tool := NewOpenAIChatTool(mockService)
		
		args := map[string]interface{}{
			"model": "gpt-4",
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
		mockService := &MockOpenAIService{}
		tool := NewOpenAIChatTool(mockService)
		
		args := map[string]interface{}{
			"model": "gpt-4",
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
		assert.Contains(t, err.Error(), "temperature must be a number between 0 and 2")
	})

	t.Run("Execute_Success", func(t *testing.T) {
		mockService := &MockOpenAIService{}
		tool := NewOpenAIChatTool(mockService)
		ctx := context.Background()
		
		expectedResponse := &ChatCompletionResponse{
			ID:      "test-id",
			Object:  "chat.completion",
			Created: 1234567890,
			Model:   "gpt-4",
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
				PromptTokens:     10,
				CompletionTokens: 9,
				TotalTokens:      19,
			},
		}
		
		mockService.On("ChatCompletion", ctx, mock.AnythingOfType("*mcp.ChatCompletionRequest")).Return(expectedResponse, nil)
		
		args := map[string]interface{}{
			"model": "gpt-4",
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
		assert.Contains(t, response.Content[0].Text, "Hello! How can I help you today?")
		
		mockService.AssertExpectations(t)
	})

	t.Run("Execute_ServiceError", func(t *testing.T) {
		mockService := &MockOpenAIService{}
		tool := NewOpenAIChatTool(mockService)
		ctx := context.Background()
		
		mockService.On("ChatCompletion", ctx, mock.AnythingOfType("*mcp.ChatCompletionRequest")).Return(nil, errors.New("API error"))
		
		args := map[string]interface{}{
			"model": "gpt-4",
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
		assert.Contains(t, response.Content[0].Text, "OpenAI API error")
		
		mockService.AssertExpectations(t)
	})
}

// TestOpenAIModelsTool 测试OpenAI模型工具
func TestOpenAIModelsTool(t *testing.T) {
	t.Run("NewOpenAIModelsTool", func(t *testing.T) {
		mockService := &MockOpenAIService{}
		tool := NewOpenAIModelsTool(mockService)
		
		assert.NotNil(t, tool)
		assert.Equal(t, "openai_models", tool.Name)
		assert.Equal(t, mockService, tool.openaiService)
	})

	t.Run("Execute_Success", func(t *testing.T) {
		mockService := &MockOpenAIService{}
		tool := NewOpenAIModelsTool(mockService)
		ctx := context.Background()
		
		expectedModels := map[string]*openai.ModelConfig{
			"gpt-4": {
				Name:             "gpt-4",
				Enabled:          true,
				MaxTokens:        8192,
				Temperature:      0.7,
				TopP:             1.0,
				FrequencyPenalty: 0.0,
				PresencePenalty:  0.0,
			},
			"gpt-3.5-turbo": {
				Name:             "gpt-3.5-turbo",
				Enabled:          false,
				MaxTokens:        4096,
				Temperature:      0.7,
				TopP:             1.0,
				FrequencyPenalty: 0.0,
				PresencePenalty:  0.0,
			},
		}
		
		mockService.On("ListModels", ctx).Return(expectedModels, nil)
		
		args := map[string]interface{}{}
		
		response, err := tool.Execute(ctx, args)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.False(t, response.IsError)
		assert.Equal(t, 1, len(response.Content))
		assert.Contains(t, response.Content[0].Text, "Available OpenAI models")
		
		mockService.AssertExpectations(t)
	})

	t.Run("Execute_ServiceError", func(t *testing.T) {
		mockService := &MockOpenAIService{}
		tool := NewOpenAIModelsTool(mockService)
		ctx := context.Background()
		
		mockService.On("ListModels", ctx).Return(nil, errors.New("failed to list models"))
		
		args := map[string]interface{}{}
		
		response, err := tool.Execute(ctx, args)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.True(t, response.IsError)
		assert.Contains(t, response.Content[0].Text, "Error listing models")
		
		mockService.AssertExpectations(t)
	})
}

// TestOpenAIConfigTool 测试OpenAI配置工具
func TestOpenAIConfigTool(t *testing.T) {
	t.Run("NewOpenAIConfigTool", func(t *testing.T) {
		mockService := &MockOpenAIService{}
		tool := NewOpenAIConfigTool(mockService)
		
		assert.NotNil(t, tool)
		assert.Equal(t, "openai_config", tool.Name)
		assert.Equal(t, mockService, tool.openaiService)
	})

	t.Run("Validate_Success", func(t *testing.T) {
		mockService := &MockOpenAIService{}
		tool := NewOpenAIConfigTool(mockService)
		
		args := map[string]interface{}{
			"action": "set_api_key",
			"api_key": "test-api-key",
		}
		
		err := tool.Validate(args)
		assert.NoError(t, err)
	})

	t.Run("Validate_MissingAction", func(t *testing.T) {
		mockService := &MockOpenAIService{}
		tool := NewOpenAIConfigTool(mockService)
		
		args := map[string]interface{}{}
		
		err := tool.Validate(args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "action parameter is required")
	})

	t.Run("Validate_InvalidAction", func(t *testing.T) {
		mockService := &MockOpenAIService{}
		tool := NewOpenAIConfigTool(mockService)
		
		args := map[string]interface{}{
			"action": "invalid_action",
		}
		
		err := tool.Validate(args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid action")
	})

	t.Run("Execute_SetAPIKey_Success", func(t *testing.T) {
		mockService := &MockOpenAIService{}
		tool := NewOpenAIConfigTool(mockService)
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
		assert.Contains(t, response.Content[0].Text, "API key set successfully")
		
		mockService.AssertExpectations(t)
	})

	t.Run("Execute_ValidateAPIKey_Success", func(t *testing.T) {
		mockService := &MockOpenAIService{}
		tool := NewOpenAIConfigTool(mockService)
		ctx := context.Background()
		
		mockService.On("ValidateAPIKey", ctx).Return(nil)
		
		args := map[string]interface{}{
			"action": "validate_key",
		}
		
		response, err := tool.Execute(ctx, args)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.False(t, response.IsError)
		assert.Contains(t, response.Content[0].Text, "API key is valid")
		
		mockService.AssertExpectations(t)
	})

	t.Run("Execute_EnableModel_Success", func(t *testing.T) {
		mockService := &MockOpenAIService{}
		tool := NewOpenAIConfigTool(mockService)
		ctx := context.Background()
		
		mockService.On("EnableModel", "gpt-4").Return(nil)
		
		args := map[string]interface{}{
			"action": "enable_model",
			"model":  "gpt-4",
		}
		
		response, err := tool.Execute(ctx, args)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.False(t, response.IsError)
		assert.Contains(t, response.Content[0].Text, "Model gpt-4 enabled successfully")
		
		mockService.AssertExpectations(t)
	})

	t.Run("Execute_DisableModel_Success", func(t *testing.T) {
		mockService := &MockOpenAIService{}
		tool := NewOpenAIConfigTool(mockService)
		ctx := context.Background()
		
		mockService.On("DisableModel", "gpt-4").Return(nil)
		
		args := map[string]interface{}{
			"action": "disable_model",
			"model":  "gpt-4",
		}
		
		response, err := tool.Execute(ctx, args)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.False(t, response.IsError)
		assert.Contains(t, response.Content[0].Text, "Model gpt-4 disabled successfully")
		
		mockService.AssertExpectations(t)
	})

	t.Run("Execute_ServiceError", func(t *testing.T) {
		mockService := &MockOpenAIService{}
		tool := NewOpenAIConfigTool(mockService)
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