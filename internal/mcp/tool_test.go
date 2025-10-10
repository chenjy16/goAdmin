package mcp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

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