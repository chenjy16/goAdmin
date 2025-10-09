package mcp

import (
	"context"
	"errors"
	"testing"

	"go-springAi/internal/openai"

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