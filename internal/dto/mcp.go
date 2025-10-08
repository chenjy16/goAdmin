package dto

import (
	"encoding/json"
	"time"
)

// MCPTool MCP工具定义
type MCPTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// MCPToolsResponse 工具列表响应
type MCPToolsResponse struct {
	Tools []MCPTool `json:"tools"`
}

// MCPExecuteRequest 工具执行请求
type MCPExecuteRequest struct {
	Name      string                 `json:"name" binding:"required"`
	Arguments map[string]interface{} `json:"arguments"`
}

// MCPExecuteResponse 工具执行响应
type MCPExecuteResponse struct {
	Content []MCPContent `json:"content"`
	IsError bool         `json:"isError,omitempty"`
}

// MCPContent MCP内容结构
type MCPContent struct {
	Type string      `json:"type"`
	Text string      `json:"text,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

// MCPMessage MCP消息结构（用于SSE）
type MCPMessage struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *MCPError       `json:"error,omitempty"`
}

// MCPError MCP错误结构
type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// MCPInitializeRequest 初始化请求
type MCPInitializeRequest struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    MCPCapabilities        `json:"capabilities"`
	ClientInfo      MCPClientInfo          `json:"clientInfo"`
	Meta            map[string]interface{} `json:"meta,omitempty"`
}

// MCPCapabilities MCP能力声明
type MCPCapabilities struct {
	Tools     *MCPToolsCapability     `json:"tools,omitempty"`
	Resources *MCPResourcesCapability `json:"resources,omitempty"`
	Prompts   *MCPPromptsCapability   `json:"prompts,omitempty"`
	Logging   *MCPLoggingCapability   `json:"logging,omitempty"`
}

// MCPToolsCapability 工具能力
type MCPToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// MCPResourcesCapability 资源能力
type MCPResourcesCapability struct {
	Subscribe   bool `json:"subscribe,omitempty"`
	ListChanged bool `json:"listChanged,omitempty"`
}

// MCPPromptsCapability 提示能力
type MCPPromptsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// MCPLoggingCapability 日志能力
type MCPLoggingCapability struct{}

// MCPClientInfo 客户端信息
type MCPClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// MCPInitializeResponse 初始化响应
type MCPInitializeResponse struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    MCPCapabilities        `json:"capabilities"`
	ServerInfo      MCPServerInfo          `json:"serverInfo"`
	Instructions    string                 `json:"instructions,omitempty"`
	Meta            map[string]interface{} `json:"meta,omitempty"`
}

// MCPServerInfo 服务器信息
type MCPServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// MCPSSEEvent SSE事件结构
type MCPSSEEvent struct {
	ID    string `json:"id,omitempty"`
	Event string `json:"event,omitempty"`
	Data  string `json:"data"`
	Retry int    `json:"retry,omitempty"`
}

// MCPToolExecutionLog 工具执行日志
type MCPToolExecutionLog struct {
	ID          string                 `json:"id"`
	ToolName    string                 `json:"toolName"`
	Arguments   map[string]interface{} `json:"arguments"`
	Result      *MCPExecuteResponse    `json:"result,omitempty"`
	Error       *MCPError              `json:"error,omitempty"`
	StartTime   time.Time              `json:"startTime"`
	EndTime     *time.Time             `json:"endTime,omitempty"`
	Duration    *time.Duration         `json:"duration,omitempty"`
	UserID      *string                `json:"userId,omitempty"`
	RequestID   string                 `json:"requestId"`
}