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

// ProviderManager 提供商管理器接口
type ProviderManager interface {
	GetProviderByModel(modelName string) (ProviderInterface, error)
	GetProviderByName(name string) (ProviderInterface, error)
	ValidateModelForProvider(ctx context.Context, providerName, modelName string) error
	GetProviderByModelWithValidation(ctx context.Context, modelName string) (ProviderInterface, error)
}

// ProviderInterface 定义Provider接口，避免循环导入
type ProviderInterface interface {
	GetType() string
	GetName() string
	ChatCompletion(ctx context.Context, request *ProviderChatRequest) (*ProviderChatResponse, error)
}

// ProviderChatRequest 提供商聊天请求结构
type ProviderChatRequest struct {
	Model       string                 `json:"model"`
	Messages    []ProviderMessage      `json:"messages"`
	MaxTokens   *int                   `json:"max_tokens,omitempty"`
	Temperature *float32               `json:"temperature,omitempty"`
	TopP        *float32               `json:"top_p,omitempty"`
	TopK        *int                   `json:"top_k,omitempty"`
	Stream      bool                   `json:"stream,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
}

// ProviderChatResponse 提供商聊天响应结构
type ProviderChatResponse struct {
	ID      string                `json:"id"`
	Object  string                `json:"object"`
	Created int64                 `json:"created"`
	Model   string                `json:"model"`
	Choices []ProviderChoice      `json:"choices"`
	Usage   ProviderUsage         `json:"usage"`
}

// ProviderMessage 提供商消息结构
type ProviderMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ProviderChoice 提供商选择结构
type ProviderChoice struct {
	Index        int             `json:"index"`
	Message      ProviderMessage `json:"message"`
	FinishReason string          `json:"finish_reason"`
}

// ProviderUsage 提供商使用情况结构
type ProviderUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// AIAssistantService AI助手服务，集成MCP客户端和Provider管理器
type AIAssistantService struct {
	mcpClient       mcp.InternalMCPClient
	openaiService   *OpenAIService
	providerManager ProviderManager
	logger          *zap.Logger
}

// NewAIAssistantService 创建AI助手服务
func NewAIAssistantService(
	mcpClient mcp.InternalMCPClient,
	openaiService *OpenAIService,
	providerManager ProviderManager,
	logger *zap.Logger,
) *AIAssistantService {
	return &AIAssistantService{
		mcpClient:       mcpClient,
		openaiService:   openaiService,
		providerManager: providerManager,
		logger:          logger,
	}
}

// ChatRequest AI助手聊天请求
type ChatRequest struct {
	Messages      []openai.Message `json:"messages"`
	Model         string           `json:"model,omitempty"`
	MaxTokens     *int             `json:"max_tokens,omitempty"`
	Temperature   *float32         `json:"temperature,omitempty"`
	UseTools      bool             `json:"use_tools,omitempty"`
	Provider      string           `json:"provider,omitempty"`      // 指定提供商
	SelectedTools []string         `json:"selected_tools,omitempty"` // 指定要使用的工具列表
}

// ChatResponse AI助手聊天响应
type ChatResponse struct {
	ID      string                `json:"id"`
	Object  string                `json:"object"`
	Created int64                 `json:"created"`
	Model   string                `json:"model"`
	Choices []ChatChoice          `json:"choices"`
	Usage   openai.Usage          `json:"usage"`
}

// ChatChoice 聊天选择
type ChatChoice struct {
	Index        int                  `json:"index"`
	Message      openai.Message       `json:"message"`
	FinishReason string               `json:"finish_reason"`
	ToolCalls    []ToolCallExecution  `json:"tool_calls,omitempty"`
}

// ToolCallExecution 工具调用执行结果
type ToolCallExecution struct {
	ToolName    string                 `json:"tool_name"`
	Arguments   map[string]interface{} `json:"arguments"`
	Result      *dto.MCPExecuteResponse `json:"result"`
	Error       string                 `json:"error,omitempty"`
	ExecutionID string                 `json:"execution_id,omitempty"`
}

// Chat 进行AI对话，支持动态提供商选择和工具调用
func (s *AIAssistantService) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	s.logger.Info("AI assistant chat request",
		zap.String("model", req.Model),
		zap.String("provider", req.Provider),
		zap.Int("message_count", len(req.Messages)),
		zap.Bool("use_tools", req.UseTools),
		zap.Strings("selected_tools", req.SelectedTools))

	// 1. 动态提供商选择和模型验证
	var provider ProviderInterface
	var err error
	
	if req.Provider != "" {
		// 如果明确指定了提供商，尝试通过提供商名称获取
		s.logger.Info("Using explicitly specified provider", zap.String("provider", req.Provider))
		provider, err = s.providerManager.GetProviderByName(req.Provider)
		if err != nil {
			s.logger.Error("Failed to get provider by name", 
				zap.String("provider", req.Provider), zap.Error(err))
			return nil, fmt.Errorf("provider %s not found", req.Provider)
		}
		
		// 验证模型是否存在于指定的提供商中
		if req.Model != "" {
			if validateErr := s.providerManager.ValidateModelForProvider(ctx, req.Provider, req.Model); validateErr != nil {
				s.logger.Error("Model validation failed", 
					zap.String("provider", req.Provider),
					zap.String("model", req.Model),
					zap.Error(validateErr))
				return nil, fmt.Errorf("model %s not supported by provider %s", req.Model, req.Provider)
			}
		}
	} else {
		// 根据模型名称自动选择提供商（使用验证版本）
		if req.Model != "" {
			provider, err = s.providerManager.GetProviderByModelWithValidation(ctx, req.Model)
			if err != nil {
				s.logger.Warn("Failed to find provider with model validation, falling back to prefix matching", 
					zap.String("model", req.Model), zap.Error(err))
				// 回退到原有的前缀匹配方式
				provider, err = s.providerManager.GetProviderByModel(req.Model)
			}
		} else {
			// 如果没有指定模型，使用Mock提供商作为默认提供商
			s.logger.Info("No model specified, using default mock provider")
			provider, err = s.providerManager.GetProviderByName("mock")
			if err != nil {
				s.logger.Warn("Failed to get mock provider, falling back to gpt-3.5-turbo", zap.Error(err))
				provider, err = s.providerManager.GetProviderByModel("gpt-3.5-turbo") // 回退到默认模型
			} else {
				// 为Mock提供商设置默认模型
				if req.Model == "" {
					req.Model = "mock-gpt-3.5-turbo"
				}
			}
		}
	}
	
	if err != nil {
		s.logger.Error("Failed to get provider", zap.Error(err))
		// 回退到原有的OpenAI实现
		return s.chatWithOpenAI(ctx, req)
	}

	// 2. 工具过滤和获取
	var availableTools []dto.MCPTool
	if req.UseTools || len(req.SelectedTools) > 0 {
		toolsResp, err := s.mcpClient.ListTools(ctx)
		if err != nil {
			s.logger.Error("Failed to get available tools", zap.Error(err))
			return nil, fmt.Errorf("failed to get available tools: %w", err)
		}
		
		// 根据SelectedTools过滤工具
		if len(req.SelectedTools) > 0 {
			availableTools = s.filterTools(toolsResp.Tools, req.SelectedTools)
		} else {
			availableTools = toolsResp.Tools
		}
	}

	// 3. 使用动态选择的提供商进行聊天
	s.logger.Info("Using provider for chat", 
		zap.String("provider_type", provider.GetType()),
		zap.String("provider_name", provider.GetName()),
		zap.Int("available_tools", len(availableTools)))

	// 构建提供商聊天请求
	providerMessages := make([]ProviderMessage, len(req.Messages))
	for i, msg := range req.Messages {
		providerMessages[i] = ProviderMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// 检查是否需要添加工具信息到系统消息
	if len(availableTools) > 0 {
		toolsInfo := s.buildToolsSystemMessage(availableTools)
		systemMsg := ProviderMessage{
			Role:    "system",
			Content: toolsInfo,
		}
		
		// 如果第一条消息已经是系统消息，则替换；否则添加到开头
		if len(providerMessages) > 0 && providerMessages[0].Role == "system" {
			providerMessages[0] = systemMsg
		} else {
			providerMessages = append([]ProviderMessage{systemMsg}, providerMessages...)
		}
	}

	providerReq := &ProviderChatRequest{
		Model:       req.Model,
		Messages:    providerMessages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
	}

	// 调用提供商
	providerResp, err := provider.ChatCompletion(ctx, providerReq)
	if err != nil {
		s.logger.Error("Provider chat failed", zap.Error(err))
		return nil, fmt.Errorf("provider chat failed: %w", err)
	}

	// 转换响应格式
	if len(providerResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from provider")
	}

	choice := providerResp.Choices[0]
	response := &ChatResponse{
		ID:      providerResp.ID,
		Object:  providerResp.Object,
		Created: providerResp.Created,
		Model:   providerResp.Model,
		Choices: []ChatChoice{
			{
				Index: 0,
				Message: openai.Message{
					Role:    choice.Message.Role,
					Content: choice.Message.Content,
				},
				FinishReason: choice.FinishReason,
			},
		},
		Usage: openai.Usage{
			PromptTokens:     providerResp.Usage.PromptTokens,
			CompletionTokens: providerResp.Usage.CompletionTokens,
			TotalTokens:      providerResp.Usage.TotalTokens,
		},
	}

	// 4. 处理工具调用（如果需要）
	// 检查是否有可用工具
	if len(availableTools) > 0 && len(response.Choices) > 0 && response.Choices[0].Message.Content != "" {
		toolCalls := s.parseToolCalls(response.Choices[0].Message.Content)
		if len(toolCalls) > 0 {
			s.logger.Info("Executing tool calls", zap.Int("count", len(toolCalls)))
			
			executions := make([]ToolCallExecution, 0, len(toolCalls))
			for _, toolCall := range toolCalls {
				execution := s.executeToolCall(ctx, toolCall)
				executions = append(executions, execution)
			}
			
			response.Choices[0].ToolCalls = executions
			
			// 如果有工具调用结果，可以选择再次调用提供商生成最终回复
			if s.shouldGenerateFinalResponse(executions) {
				finalResp, err := s.generateFinalResponse(ctx, provider, req, executions)
				if err != nil {
					s.logger.Warn("Failed to generate final response", zap.Error(err))
				} else {
					response.Choices[0].Message = finalResp
				}
			}
		}
	}

	return response, nil
}

// filterTools 根据选择的工具名称过滤工具列表
func (s *AIAssistantService) filterTools(allTools []dto.MCPTool, selectedTools []string) []dto.MCPTool {
	if len(selectedTools) == 0 {
		return allTools
	}
	
	selectedSet := make(map[string]bool)
	for _, toolName := range selectedTools {
		selectedSet[toolName] = true
	}
	
	var filtered []dto.MCPTool
	for _, tool := range allTools {
		if selectedSet[tool.Name] {
			filtered = append(filtered, tool)
		}
	}
	
	s.logger.Info("Filtered tools", 
		zap.Int("total_tools", len(allTools)),
		zap.Int("selected_tools", len(filtered)),
		zap.Strings("tool_names", selectedTools))
	
	return filtered
}

// chatWithOpenAI 回退到原有的OpenAI实现（向后兼容）
func (s *AIAssistantService) chatWithOpenAI(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	s.logger.Info("Falling back to OpenAI implementation")
	
	// 如果启用工具或指定了工具，先获取可用工具列表
	var availableTools []dto.MCPTool
	if req.UseTools || len(req.SelectedTools) > 0 {
		toolsResp, err := s.mcpClient.ListTools(ctx)
		if err != nil {
			s.logger.Error("Failed to get available tools", zap.Error(err))
			return nil, fmt.Errorf("failed to get available tools: %w", err)
		}
		
		// 根据SelectedTools过滤工具
		if len(req.SelectedTools) > 0 {
			availableTools = s.filterTools(toolsResp.Tools, req.SelectedTools)
		} else {
			availableTools = toolsResp.Tools
		}
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
		ID:      openaiResp.ID,
		Object:  openaiResp.Object,
		Created: openaiResp.Created,
		Model:   openaiResp.Model,
		Choices: []ChatChoice{
			{
				Index:        choice.Index,
				Message:      choice.Message,
				FinishReason: choice.FinishReason,
			},
		},
		Usage: openaiResp.Usage,
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
			
			response.Choices[0].ToolCalls = executions
			
			// 如果有工具调用结果，可以选择再次调用OpenAI生成最终回复
			if s.shouldGenerateFinalResponse(executions) {
				// 获取OpenAI提供商用于生成最终回复
				openaiProvider, err := s.providerManager.GetProviderByName("OpenAI")
				if err != nil {
					s.logger.Warn("Failed to get OpenAI provider for final response", zap.Error(err))
				} else {
					finalResp, err := s.generateFinalResponse(ctx, openaiProvider, req, executions)
					if err != nil {
						s.logger.Warn("Failed to generate final response", zap.Error(err))
					} else {
						response.Choices[0].Message = finalResp
					}
				}
			}
		}
	}

	return response, nil
}

// buildToolsSystemMessage 构建工具系统消息
func (s *AIAssistantService) buildToolsSystemMessage(tools []dto.MCPTool) string {
	var builder strings.Builder
	builder.WriteString("You are a professional financial AI assistant with access to comprehensive stock analysis tools. ")
	builder.WriteString("When users ask about stocks, investments, or financial analysis, you should use the appropriate tools to provide accurate, data-driven insights. ")
	builder.WriteString("When you need to use a tool, respond with a JSON object in this format: ")
	builder.WriteString(`{"tool_call": {"name": "tool_name", "arguments": {...}}}`)
	builder.WriteString("\n\nAvailable tools:\n")
	
	// 工具已经在调用方过滤过了，这里直接使用
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
	
	s.logger.Info("Parsing tool calls", zap.String("content", content))
	
	// 解析XML风格的tool_call标签
	startTag := "<tool_call>"
	endTag := "</tool_call>"
	
	startIndex := strings.Index(content, startTag)
	s.logger.Info("Tool call search", zap.Int("startIndex", startIndex))
	for startIndex != -1 {
		endIndex := strings.Index(content[startIndex:], endTag)
		if endIndex == -1 {
			break
		}
		
		// 提取tool_call内容
		toolCallContent := content[startIndex+len(startTag) : startIndex+endIndex]
		toolCallContent = strings.TrimSpace(toolCallContent)
		
		// 解析JSON内容
		var toolCallData map[string]interface{}
		if err := json.Unmarshal([]byte(toolCallContent), &toolCallData); err == nil {
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
		
		// 继续查找下一个tool_call
		nextSearchStart := startIndex + endIndex + len(endTag)
		nextIndex := strings.Index(content[nextSearchStart:], startTag)
		if nextIndex != -1 {
			startIndex = nextSearchStart + nextIndex
		} else {
			startIndex = -1
		}
	}
	
	s.logger.Info("Tool calls parsed", zap.Int("count", len(toolCalls)))
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
func (s *AIAssistantService) generateFinalResponse(ctx context.Context, provider ProviderInterface, originalReq *ChatRequest, executions []ToolCallExecution) (openai.Message, error) {
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
	
	// 构建提供商请求的消息格式
	providerMessages := make([]ProviderMessage, 0, len(originalReq.Messages)+2)
	
	// 转换原始消息
	for _, msg := range originalReq.Messages {
		providerMessages = append(providerMessages, ProviderMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}
	
	// 添加工具执行结果
	providerMessages = append(providerMessages, ProviderMessage{
		Role:    "assistant",
		Content: resultsBuilder.String(),
	})
	
	// 添加生成最终回复的指令
	providerMessages = append(providerMessages, ProviderMessage{
		Role:    "user",
		Content: "Please provide a natural language summary of the tool execution results above.",
	})
	
	// 使用动态选择的提供商生成最终回复
	finalReq := &ProviderChatRequest{
		Model:       originalReq.Model,
		Messages:    providerMessages,
		MaxTokens:   originalReq.MaxTokens,
		Temperature: originalReq.Temperature,
	}
	
	resp, err := provider.ChatCompletion(ctx, finalReq)
	if err != nil {
		return openai.Message{}, fmt.Errorf("failed to generate final response with provider %s: %w", provider.GetName(), err)
	}
	
	if len(resp.Choices) == 0 {
		return openai.Message{}, fmt.Errorf("no response from provider %s", provider.GetName())
	}
	
	// 转换回 openai.Message 格式
	return openai.Message{
		Role:    resp.Choices[0].Message.Role,
		Content: resp.Choices[0].Message.Content,
	}, nil
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