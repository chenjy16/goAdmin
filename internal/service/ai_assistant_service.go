package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go-springAi/internal/dto"
	"go-springAi/internal/mcp"
	"go-springAi/internal/openai"

	"go.uber.org/zap"
)

// AIAssistantService AI助手服务，集成MCP客户端和OpenAI
type AIAssistantService struct {
	mcpClient     mcp.InternalMCPClient
	openaiService *OpenAIService
	logger        *zap.Logger
}

// NewAIAssistantService 创建AI助手服务
func NewAIAssistantService(
	mcpClient mcp.InternalMCPClient,
	openaiService *OpenAIService,
	logger *zap.Logger,
) *AIAssistantService {
	return &AIAssistantService{
		mcpClient:     mcpClient,
		openaiService: openaiService,
		logger:        logger,
	}
}

// ChatRequest AI助手聊天请求
type ChatRequest struct {
	Messages    []openai.Message `json:"messages"`
	Model       string           `json:"model,omitempty"`
	MaxTokens   *int             `json:"max_tokens,omitempty"`
	Temperature *float32         `json:"temperature,omitempty"`
	UseTools    bool             `json:"use_tools,omitempty"`
}

// ChatResponse AI助手聊天响应
type ChatResponse struct {
	Message       openai.Message         `json:"message"`
	ToolCalls     []ToolCallExecution    `json:"tool_calls,omitempty"`
	Usage         openai.Usage           `json:"usage"`
	FinishReason  string                 `json:"finish_reason"`
}

// ToolCallExecution 工具调用执行结果
type ToolCallExecution struct {
	ToolName    string                 `json:"tool_name"`
	Arguments   map[string]interface{} `json:"arguments"`
	Result      *dto.MCPExecuteResponse `json:"result"`
	Error       string                 `json:"error,omitempty"`
	ExecutionID string                 `json:"execution_id,omitempty"`
}

// Chat 进行AI对话，支持工具调用
func (s *AIAssistantService) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	s.logger.Info("AI assistant chat request",
		zap.String("model", req.Model),
		zap.Int("message_count", len(req.Messages)),
		zap.Bool("use_tools", req.UseTools))

	// 如果启用工具，先获取可用工具列表
	var availableTools []dto.MCPTool
	if req.UseTools {
		toolsResp, err := s.mcpClient.ListTools(ctx)
		if err != nil {
			s.logger.Error("Failed to get available tools", zap.Error(err))
			return nil, fmt.Errorf("failed to get available tools: %w", err)
		}
		availableTools = toolsResp.Tools
	}

	// 构建OpenAI请求
	openaiReq := &ChatCompletionRequest{
		Model:       req.Model,
		Messages:    req.Messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
	}

	// 如果有可用工具，添加工具信息到系统消息
	if len(availableTools) > 0 {
		toolsInfo := s.buildToolsSystemMessage(availableTools)
		openaiReq.Messages = s.addSystemMessage(openaiReq.Messages, toolsInfo)
	}

	// 调用OpenAI
	openaiResp, err := s.openaiService.ChatCompletion(ctx, openaiReq)
	if err != nil {
		s.logger.Error("OpenAI chat completion failed", zap.Error(err))
		return nil, fmt.Errorf("OpenAI chat completion failed: %w", err)
	}

	if len(openaiResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	choice := openaiResp.Choices[0]
	response := &ChatResponse{
		Message:      choice.Message,
		Usage:        openaiResp.Usage,
		FinishReason: choice.FinishReason,
	}

	// 检查是否需要执行工具调用
	if req.UseTools && len(availableTools) > 0 {
		toolCalls := s.parseToolCalls(choice.Message.Content)
		if len(toolCalls) > 0 {
			s.logger.Info("Executing tool calls", zap.Int("count", len(toolCalls)))
			
			executions := make([]ToolCallExecution, 0, len(toolCalls))
			for _, toolCall := range toolCalls {
				execution := s.executeToolCall(ctx, toolCall)
				executions = append(executions, execution)
			}
			
			response.ToolCalls = executions
			
			// 如果有工具调用结果，可以选择再次调用OpenAI生成最终回复
			if s.shouldGenerateFinalResponse(executions) {
				finalResp, err := s.generateFinalResponse(ctx, req, executions)
				if err != nil {
					s.logger.Warn("Failed to generate final response", zap.Error(err))
				} else {
					response.Message = finalResp
				}
			}
		}
	}

	return response, nil
}

// buildToolsSystemMessage 构建工具系统消息
func (s *AIAssistantService) buildToolsSystemMessage(tools []dto.MCPTool) string {
	var builder strings.Builder
	builder.WriteString("You have access to the following tools. ")
	builder.WriteString("When you need to use a tool, respond with a JSON object in this format: ")
	builder.WriteString(`{"tool_call": {"name": "tool_name", "arguments": {...}}}`)
	builder.WriteString("\n\nAvailable tools:\n")
	
	for _, tool := range tools {
		builder.WriteString(fmt.Sprintf("- %s: %s\n", tool.Name, tool.Description))
		if schemaBytes, err := json.Marshal(tool.InputSchema); err == nil {
			builder.WriteString(fmt.Sprintf("  Schema: %s\n", string(schemaBytes)))
		}
	}
	
	return builder.String()
}

// addSystemMessage 添加系统消息
func (s *AIAssistantService) addSystemMessage(messages []openai.Message, systemContent string) []openai.Message {
	systemMsg := openai.Message{
		Role:    "system",
		Content: systemContent,
	}
	
	// 如果第一条消息已经是系统消息，则替换；否则添加到开头
	if len(messages) > 0 && messages[0].Role == "system" {
		messages[0] = systemMsg
		return messages
	}
	
	return append([]openai.Message{systemMsg}, messages...)
}

// ToolCall 工具调用结构
type ToolCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// parseToolCalls 解析工具调用
func (s *AIAssistantService) parseToolCalls(content string) []ToolCall {
	var toolCalls []ToolCall
	
	// 简单的JSON解析，寻找tool_call对象
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "tool_call") {
			var parsed map[string]interface{}
			if err := json.Unmarshal([]byte(line), &parsed); err == nil {
				if toolCallData, ok := parsed["tool_call"].(map[string]interface{}); ok {
					if name, ok := toolCallData["name"].(string); ok {
						toolCall := ToolCall{
							Name: name,
						}
						if args, ok := toolCallData["arguments"].(map[string]interface{}); ok {
							toolCall.Arguments = args
						}
						toolCalls = append(toolCalls, toolCall)
					}
				}
			}
		}
	}
	
	return toolCalls
}

// executeToolCall 执行工具调用
func (s *AIAssistantService) executeToolCall(ctx context.Context, toolCall ToolCall) ToolCallExecution {
	execution := ToolCallExecution{
		ToolName:  toolCall.Name,
		Arguments: toolCall.Arguments,
	}
	
	// 执行MCP工具
	mcpReq := &dto.MCPExecuteRequest{
		Name:      toolCall.Name,
		Arguments: toolCall.Arguments,
	}
	
	result, err := s.mcpClient.ExecuteTool(ctx, mcpReq)
	if err != nil {
		execution.Error = err.Error()
		s.logger.Error("Tool execution failed",
			zap.String("tool", toolCall.Name),
			zap.Error(err))
	} else {
		execution.Result = result
		s.logger.Info("Tool executed successfully",
			zap.String("tool", toolCall.Name),
			zap.Bool("is_error", result.IsError))
	}
	
	return execution
}

// shouldGenerateFinalResponse 判断是否需要生成最终回复
func (s *AIAssistantService) shouldGenerateFinalResponse(executions []ToolCallExecution) bool {
	// 如果有任何工具执行成功，则生成最终回复
	for _, exec := range executions {
		if exec.Error == "" && exec.Result != nil && !exec.Result.IsError {
			return true
		}
	}
	return false
}

// generateFinalResponse 生成最终回复
func (s *AIAssistantService) generateFinalResponse(ctx context.Context, originalReq *ChatRequest, executions []ToolCallExecution) (openai.Message, error) {
	// 构建包含工具执行结果的消息
	var resultsBuilder strings.Builder
	resultsBuilder.WriteString("Tool execution results:\n")
	
	for _, exec := range executions {
		resultsBuilder.WriteString(fmt.Sprintf("Tool: %s\n", exec.ToolName))
		if exec.Error != "" {
			resultsBuilder.WriteString(fmt.Sprintf("Error: %s\n", exec.Error))
		} else if exec.Result != nil {
			for _, content := range exec.Result.Content {
				resultsBuilder.WriteString(fmt.Sprintf("Result: %s\n", content.Text))
			}
		}
		resultsBuilder.WriteString("\n")
	}
	
	// 添加工具结果到消息历史
	messages := append(originalReq.Messages, openai.Message{
		Role:    "assistant",
		Content: resultsBuilder.String(),
	})
	
	messages = append(messages, openai.Message{
		Role:    "user",
		Content: "Please provide a natural language summary of the tool execution results above.",
	})
	
	// 调用OpenAI生成最终回复
	finalReq := &ChatCompletionRequest{
		Model:       originalReq.Model,
		Messages:    messages,
		MaxTokens:   originalReq.MaxTokens,
		Temperature: originalReq.Temperature,
	}
	
	resp, err := s.openaiService.ChatCompletion(ctx, finalReq)
	if err != nil {
		return openai.Message{}, err
	}
	
	if len(resp.Choices) == 0 {
		return openai.Message{}, fmt.Errorf("no response from OpenAI")
	}
	
	return resp.Choices[0].Message, nil
}

// Initialize 初始化AI助手服务
func (s *AIAssistantService) Initialize(ctx context.Context) error {
	// 初始化MCP客户端
	initReq := &dto.MCPInitializeRequest{
		ProtocolVersion: "2024-11-05",
		Capabilities: dto.MCPCapabilities{
			Tools: &dto.MCPToolsCapability{
				ListChanged: true,
			},
		},
		ClientInfo: dto.MCPClientInfo{
			Name:    "AI Assistant",
			Version: "1.0.0",
		},
	}
	
	_, err := s.mcpClient.Initialize(ctx, initReq)
	if err != nil {
		return fmt.Errorf("failed to initialize MCP client: %w", err)
	}
	
	s.logger.Info("AI assistant service initialized successfully")
	return nil
}