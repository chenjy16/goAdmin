package mcp

import (
	"context"
	"testing"

	"go-springAi/internal/dto"

	"github.com/stretchr/testify/assert"
)

// TestTool 测试工具（用于测试目的）
type TestTool struct {
	*BaseTool
}

// NewTestTool 创建测试工具
func NewTestTool() *TestTool {
	return &TestTool{
		BaseTool: &BaseTool{
			Name:        "test_tool",
			Description: "A simple test tool",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"input": map[string]interface{}{
						"type":        "string",
						"description": "Test input",
					},
				},
				"required": []string{"input"},
			},
		},
	}
}

// Execute 执行测试工具
func (tt *TestTool) Execute(ctx context.Context, args map[string]interface{}) (*dto.MCPExecuteResponse, error) {
	return &dto.MCPExecuteResponse{
		Content: []dto.MCPContent{
			{
				Type: "text",
				Text: "Test result",
			},
		},
		IsError: false,
	}, nil
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
		testTool := NewTestTool()
		
		registry.Register(testTool)
		
		tool, exists := registry.GetTool("test_tool")
		assert.True(t, exists)
		assert.Equal(t, testTool, tool)
		
		_, exists = registry.GetTool("nonexistent")
		assert.False(t, exists)
	})

	t.Run("ListTools", func(t *testing.T) {
		registry := NewToolRegistry()
		testTool := NewTestTool()
		
		registry.Register(testTool)
		
		tools := registry.ListTools()
		assert.Equal(t, 1, len(tools))
		assert.Equal(t, "test_tool", tools[0].Name)
	})

	t.Run("GetToolNames", func(t *testing.T) {
		registry := NewToolRegistry()
		testTool := NewTestTool()
		
		registry.Register(testTool)
		
		names := registry.GetToolNames()
		assert.Equal(t, 1, len(names))
		assert.Contains(t, names, "test_tool")
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