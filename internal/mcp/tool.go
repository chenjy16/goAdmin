package mcp

import (
	"context"
	"fmt"

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

// EchoTool 回显工具（示例工具）
type EchoTool struct {
	*BaseTool
}

// NewEchoTool 创建回显工具
func NewEchoTool() *EchoTool {
	return &EchoTool{
		BaseTool: &BaseTool{
			Name:        "echo",
			Description: "Echo the input message back to the user",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"message": map[string]interface{}{
						"type":        "string",
						"description": "The message to echo back",
					},
				},
				"required": []string{"message"},
			},
		},
	}
}

// Execute 执行回显工具
func (et *EchoTool) Execute(ctx context.Context, args map[string]interface{}) (*dto.MCPExecuteResponse, error) {
	message, ok := args["message"].(string)
	if !ok {
		return nil, fmt.Errorf("message parameter is required and must be a string")
	}

	return &dto.MCPExecuteResponse{
		Content: []dto.MCPContent{
			{
				Type: "text",
				Text: fmt.Sprintf("Echo: %s", message),
			},
		},
		IsError: false,
	}, nil
}

// UserInfoTool 用户信息工具
type UserInfoTool struct {
	*BaseTool
	userService UserService // 注入用户服务
}

// UserService 用户服务接口（简化版）
type UserService interface {
	GetUser(ctx context.Context, id int64) (*dto.UserResponse, error)
	ListUsers(ctx context.Context, page, limit int64) ([]*dto.UserResponse, error)
}

// NewUserInfoTool 创建用户信息工具
func NewUserInfoTool(userService UserService) *UserInfoTool {
	return &UserInfoTool{
		BaseTool: &BaseTool{
			Name:        "get_user_info",
			Description: "Get user information by user ID or list all users",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"user_id": map[string]interface{}{
						"type":        "number",
						"description": "The user ID to get information for (optional)",
					},
					"list_all": map[string]interface{}{
						"type":        "boolean",
						"description": "Whether to list all users (optional)",
					},
					"page": map[string]interface{}{
						"type":        "number",
						"description": "Page number for listing users (default: 1)",
					},
					"limit": map[string]interface{}{
						"type":        "number",
						"description": "Number of users per page (default: 10)",
					},
				},
			},
		},
		userService: userService,
	}
}

// Execute 执行用户信息工具
func (uit *UserInfoTool) Execute(ctx context.Context, args map[string]interface{}) (*dto.MCPExecuteResponse, error) {
	// 检查是否要列出所有用户
	if listAll, ok := args["list_all"].(bool); ok && listAll {
		page := int64(1)
		limit := int64(10)
		
		if p, ok := args["page"].(float64); ok {
			page = int64(p)
		}
		if l, ok := args["limit"].(float64); ok {
			limit = int64(l)
		}

		users, err := uit.userService.ListUsers(ctx, page, limit)
		if err != nil {
			return &dto.MCPExecuteResponse{
				Content: []dto.MCPContent{
					{
						Type: "text",
						Text: fmt.Sprintf("Error listing users: %v", err),
					},
				},
				IsError: true,
			}, nil
		}

		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Found %d users", len(users)),
					Data: users,
				},
			},
			IsError: false,
		}, nil
	}

	// 获取特定用户信息
	userID, ok := args["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("user_id parameter is required when list_all is not true")
	}

	user, err := uit.userService.GetUser(ctx, int64(userID))
	if err != nil {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Error getting user: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	return &dto.MCPExecuteResponse{
		Content: []dto.MCPContent{
			{
				Type: "text",
				Text: fmt.Sprintf("User information for ID %d", int64(userID)),
				Data: user,
			},
		},
		IsError: false,
	}, nil
}