package mcp

import (
	"context"
	"errors"
	"testing"

	"go-springAi/internal/googleai"

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