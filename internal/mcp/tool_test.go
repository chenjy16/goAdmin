package mcp

import (
	"context"
	"testing"

	"go-springAi/internal/googleai"
	"go-springAi/internal/openai"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)



// MockGoogleAIService 模拟Google AI服务
type MockGoogleAIService struct {
	mock.Mock
}

func (m *MockGoogleAIService) ChatCompletion(ctx context.Context, req *GoogleAIChatCompletionRequest) (*GoogleAIChatCompletionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GoogleAIChatCompletionResponse), args.Error(1)
}

func (m *MockGoogleAIService) ListModels(ctx context.Context) (map[string]*googleai.ModelConfig, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]*googleai.ModelConfig), args.Error(1)
}

func (m *MockGoogleAIService) ValidateAPIKey(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockGoogleAIService) SetAPIKey(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockGoogleAIService) GetModelConfig(name string) (*googleai.ModelConfig, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*googleai.ModelConfig), args.Error(1)
}

func (m *MockGoogleAIService) EnableModel(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockGoogleAIService) DisableModel(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

// MockOpenAIService 模拟OpenAI服务
type MockOpenAIService struct {
	mock.Mock
}

func (m *MockOpenAIService) ChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ChatCompletionResponse), args.Error(1)
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

// TestToolRegistry 测试工具注册表
func TestToolRegistry(t *testing.T) {
	t.Run("NewToolRegistry", func(t *testing.T) {
		registry := NewToolRegistry()
		assert.NotNil(t, registry)
		assert.NotNil(t, registry.tools)
		assert.Equal(t, 0, len(registry.tools))
	})

	t.Run("Register and GetTool", func(t *testing.T) {
		registry := NewToolRegistry()
		echoTool := NewEchoTool()
		
		registry.Register(echoTool)
		
		tool, exists := registry.GetTool("echo")
		assert.True(t, exists)
		assert.Equal(t, echoTool, tool)
		
		_, exists = registry.GetTool("nonexistent")
		assert.False(t, exists)
	})

	t.Run("ListTools", func(t *testing.T) {
		registry := NewToolRegistry()
		echoTool := NewEchoTool()
		
		registry.Register(echoTool)
		
		tools := registry.ListTools()
		assert.Equal(t, 1, len(tools))
		assert.Equal(t, "echo", tools[0].Name)
	})

	t.Run("GetToolNames", func(t *testing.T) {
		registry := NewToolRegistry()
		echoTool := NewEchoTool()
		
		registry.Register(echoTool)
		
		names := registry.GetToolNames()
		assert.Equal(t, 1, len(names))
		assert.Contains(t, names, "echo")
	})
}

// TestBaseTool 测试基础工具
func TestBaseTool(t *testing.T) {
	t.Run("GetDefinition", func(t *testing.T) {
		baseTool := &BaseTool{
			Name:        "test_tool",
			Description: "Test tool description",
			InputSchema: map[string]interface{}{
				"type": "object",
			},
		}
		
		definition := baseTool.GetDefinition()
		assert.Equal(t, "test_tool", definition.Name)
		assert.Equal(t, "Test tool description", definition.Description)
		assert.Equal(t, map[string]interface{}{"type": "object"}, definition.InputSchema)
	})

	t.Run("Validate", func(t *testing.T) {
		baseTool := &BaseTool{}
		err := baseTool.Validate(map[string]interface{}{})
		assert.NoError(t, err)
	})
}

// TestEchoTool 测试回显工具
func TestEchoTool(t *testing.T) {
	t.Run("NewEchoTool", func(t *testing.T) {
		tool := NewEchoTool()
		assert.NotNil(t, tool)
		assert.Equal(t, "echo", tool.Name)
		assert.Equal(t, "Echo the input message back to the user", tool.Description)
	})

	t.Run("Execute_Success", func(t *testing.T) {
		tool := NewEchoTool()
		ctx := context.Background()
		args := map[string]interface{}{
			"message": "Hello, World!",
		}
		
		response, err := tool.Execute(ctx, args)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.False(t, response.IsError)
		assert.Equal(t, 1, len(response.Content))
		assert.Equal(t, "text", response.Content[0].Type)
		assert.Equal(t, "Echo: Hello, World!", response.Content[0].Text)
	})

	t.Run("Execute_MissingMessage", func(t *testing.T) {
		tool := NewEchoTool()
		ctx := context.Background()
		args := map[string]interface{}{}
		
		response, err := tool.Execute(ctx, args)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "message parameter is required")
	})

	t.Run("Execute_InvalidMessageType", func(t *testing.T) {
		tool := NewEchoTool()
		ctx := context.Background()
		args := map[string]interface{}{
			"message": 123,
		}
		
		response, err := tool.Execute(ctx, args)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "message parameter is required and must be a string")
	})
}