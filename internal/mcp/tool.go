package mcp

import (
	"context"

	"go-springAi/internal/dto"
)

// Tool MCP工具接口
type Tool interface {
	// GetDefinition 获取工具定义
	GetDefinition() dto.MCPTool
	// Execute 执行工具
	Execute(ctx context.Context, args map[string]interface{}) (*dto.MCPExecuteResponse, error)
	// Validate 验证参数
	Validate(args map[string]interface{}) error
}

// ToolRegistry 工具注册表
type ToolRegistry struct {
	tools map[string]Tool
}

// NewToolRegistry 创建工具注册表
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]Tool),
	}
}

// Register 注册工具
func (tr *ToolRegistry) Register(tool Tool) {
	definition := tool.GetDefinition()
	tr.tools[definition.Name] = tool
}

// GetTool 获取工具
func (tr *ToolRegistry) GetTool(name string) (Tool, bool) {
	tool, exists := tr.tools[name]
	return tool, exists
}

// ListTools 列出所有工具
func (tr *ToolRegistry) ListTools() []dto.MCPTool {
	tools := make([]dto.MCPTool, 0, len(tr.tools))
	for _, tool := range tr.tools {
		tools = append(tools, tool.GetDefinition())
	}
	return tools
}

// GetToolNames 获取所有工具名称
func (tr *ToolRegistry) GetToolNames() []string {
	names := make([]string, 0, len(tr.tools))
	for name := range tr.tools {
		names = append(names, name)
	}
	return names
}

// BaseTool 基础工具结构
type BaseTool struct {
	Name        string
	Description string
	InputSchema map[string]interface{}
}

// GetDefinition 实现Tool接口
func (bt *BaseTool) GetDefinition() dto.MCPTool {
	return dto.MCPTool{
		Name:        bt.Name,
		Description: bt.Description,
		InputSchema: bt.InputSchema,
	}
}

// Validate 基础验证实现
func (bt *BaseTool) Validate(args map[string]interface{}) error {
	// 基础验证逻辑，可以在具体工具中重写
	return nil
}