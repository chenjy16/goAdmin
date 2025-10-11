package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

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
	Messages     []openai.Message `json:"messages"`
	Model        string           `json:"model,omitempty"`
	MaxTokens    *int             `json:"max_tokens,omitempty"`
	Temperature  *float32         `json:"temperature,omitempty"`
	UseTools     bool             `json:"use_tools,omitempty"`
	Provider     string           `json:"provider,omitempty"`     // 指定提供商
	SelectedTool string           `json:"selected_tool,omitempty"` // 指定要使用的工具
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
		zap.String("selected_tool", req.SelectedTool))

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
				s.logger.Warn("Failed to get mock provider, falling back to mock-gpt-3.5-turbo", zap.Error(err))
				provider, err = s.providerManager.GetProviderByModel("mock-gpt-3.5-turbo") // 回退到免费的mock模型
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
	if req.UseTools || req.SelectedTool != "" {
		toolsResp, err := s.mcpClient.ListTools(ctx)
		if err != nil {
			s.logger.Error("Failed to get available tools", zap.Error(err))
			return nil, fmt.Errorf("failed to get available tools: %w", err)
		}
		
		// 根据SelectedTool过滤工具
		if req.SelectedTool != "" {
			availableTools = s.filterTool(toolsResp.Tools, req.SelectedTool)
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

// filterTool 根据单个选定工具过滤工具列表
func (s *AIAssistantService) filterTool(allTools []dto.MCPTool, selectedTool string) []dto.MCPTool {
	if selectedTool == "" {
		return allTools
	}
	
	var filtered []dto.MCPTool
	for _, tool := range allTools {
		if tool.Name == selectedTool {
			filtered = append(filtered, tool)
			break // 只需要找到一个匹配的工具
		}
	}
	
	s.logger.Info("Filtered tool", 
		zap.Int("total_tools", len(allTools)),
		zap.Int("selected_tools", len(filtered)),
		zap.String("tool_name", selectedTool))
	
	return filtered
}

// chatWithOpenAI 回退到原有的OpenAI实现（向后兼容）
func (s *AIAssistantService) chatWithOpenAI(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	s.logger.Info("Falling back to OpenAI implementation")
	
	// 如果启用工具或指定了工具，先获取可用工具列表
	var availableTools []dto.MCPTool
	if req.UseTools || req.SelectedTool != "" {
		toolsResp, err := s.mcpClient.ListTools(ctx)
		if err != nil {
			s.logger.Error("Failed to get available tools", zap.Error(err))
			return nil, fmt.Errorf("failed to get available tools: %w", err)
		}
		
		// 根据SelectedTool过滤工具
		if req.SelectedTool != "" {
			availableTools = s.filterTool(toolsResp.Tools, req.SelectedTool)
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
	builder.WriteString("# Financial AI Assistant - Professional Stock Analysis Expert\n\n")
	
	builder.WriteString("## Your Role & Expertise\n")
	builder.WriteString("You are a senior financial analyst and investment advisor with deep expertise in:\n")
	builder.WriteString("- **Stock Market Analysis**: Technical and fundamental analysis, market trends, sector analysis\n")
	builder.WriteString("- **Investment Strategy**: Portfolio optimization, risk assessment, valuation models\n")
	builder.WriteString("- **Financial Data Interpretation**: Reading financial statements, ratio analysis, performance metrics\n")
	builder.WriteString("- **Market Intelligence**: Economic indicators, industry trends, competitive analysis\n\n")
	
	builder.WriteString("## Core Responsibilities\n")
	builder.WriteString("1. **Data-Driven Analysis**: Always use tools to gather real-time, accurate financial data\n")
	builder.WriteString("2. **Professional Insights**: Provide expert-level analysis suitable for serious investors\n")
	builder.WriteString("3. **Risk Awareness**: Highlight potential risks and market uncertainties\n")
	builder.WriteString("4. **Actionable Recommendations**: Offer practical, implementable investment guidance\n")
	builder.WriteString("5. **Educational Value**: Explain complex financial concepts clearly\n\n")
	
	builder.WriteString("## Tool Usage Instructions\n")
	
	builder.WriteString("### When to Use Tools (Decision Matrix)\n")
	builder.WriteString("**ALWAYS use tools when users ask about:**\n")
	builder.WriteString("- Current stock prices, market data, or real-time information\n")
	builder.WriteString("- Specific company financial performance or metrics\n")
	builder.WriteString("- Stock comparisons or relative analysis\n")
	builder.WriteString("- Historical price movements or trends\n")
	builder.WriteString("- Portfolio analysis or investment recommendations\n\n")
	
	builder.WriteString("**DO NOT use tools for:**\n")
	builder.WriteString("- General financial education or concept explanations\n")
	builder.WriteString("- Hypothetical scenarios or theoretical discussions\n")
	builder.WriteString("- Market news interpretation (unless specific data is needed)\n")
	builder.WriteString("- Basic investment advice that doesn't require current data\n\n")
	
	builder.WriteString("### Tool Call Format\n")
	builder.WriteString("When you need to use a tool, respond with a JSON object in this exact format:\n")
	builder.WriteString("```json\n")
	builder.WriteString(`{"tool_call": {"name": "tool_name", "arguments": {...}}}`)
	builder.WriteString("\n```\n\n")
	
	builder.WriteString("### Critical Guidelines\n")
	builder.WriteString("- **One tool per response**: Never call multiple tools simultaneously\n")
	builder.WriteString("- **Single line JSON**: Provide the tool_call JSON in exactly one line\n")
	builder.WriteString("- **Complete arguments**: Include all required parameters with valid values\n")
	builder.WriteString("- **Immediate execution**: Call tools as soon as you identify the need\n")
	builder.WriteString("- **Clear intent**: Briefly explain what you're analyzing before the tool call\n\n")
	
	builder.WriteString("## Error Recovery Strategy\n")
	builder.WriteString("If a tool call fails or returns an error:\n")
	builder.WriteString("1. **Acknowledge the limitation**: Clearly state what data is unavailable\n")
	builder.WriteString("2. **Provide alternative analysis**: Use available information or general market knowledge\n")
	builder.WriteString("3. **Suggest manual verification**: Recommend users verify critical information independently\n")
	builder.WriteString("4. **Maintain professionalism**: Continue providing valuable insights despite data limitations\n")
	builder.WriteString("5. **Be transparent**: Explain how the missing data affects your analysis\n\n")
	
	builder.WriteString("## Complete Analysis Examples\n")
	builder.WriteString("### Example 1: Single Stock Analysis\n")
	builder.WriteString("**User Question**: \"How has Apple stock performed this year?\"\n")
	builder.WriteString("**Your Response**: \"I'll analyze Apple's stock performance for you.\"\n")
	builder.WriteString("**Tool Call**: ")
	builder.WriteString(`{"tool_call": {"name": "stock_analysis", "arguments": {"symbol": "AAPL", "period": "1y"}}}`)
	builder.WriteString("\n**Follow-up Analysis**: Provide comprehensive analysis of the results including price trends, volume patterns, key events, and investment implications.\n\n")
	
	builder.WriteString("### Example 2: Comparative Analysis\n")
	builder.WriteString("**User Question**: \"Should I invest in Apple or Google?\"\n")
	builder.WriteString("**Your Response**: \"Let me compare these two tech giants for you.\"\n")
	builder.WriteString("**Tool Call**: ")
	builder.WriteString(`{"tool_call": {"name": "stock_comparison", "arguments": {"symbols": ["AAPL", "GOOGL"], "metrics": ["price", "volume", "market_cap", "pe_ratio"]}}}`)
	builder.WriteString("\n**Follow-up Analysis**: Compare financial metrics, growth prospects, risk factors, and provide investment recommendation based on data.\n\n")
	
	builder.WriteString("### Example 3: Error Handling\n")
	builder.WriteString("**Scenario**: Tool call fails or returns incomplete data\n")
	builder.WriteString("**Your Response**: \"I apologize, but I'm currently unable to access real-time data for [specific stock]. However, based on recent market trends and available information, I can provide the following analysis... I recommend verifying current prices through your broker or financial platform.\"\n\n")
	
	builder.WriteString("Available tools:\n")
	
	// 工具已经在调用方过滤过了，这里直接使用
	for _, tool := range tools {
		builder.WriteString(fmt.Sprintf("### %s\n", tool.Name))
		builder.WriteString(fmt.Sprintf("Description: %s\n", tool.Description))
		if schemaBytes, err := json.Marshal(tool.InputSchema); err == nil {
			builder.WriteString(fmt.Sprintf("Schema: %s\n\n", string(schemaBytes)))
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
	
	// 清理输入内容，移除多余的空白字符
	content = strings.TrimSpace(content)
	if content == "" {
		s.logger.Warn("Empty content provided for tool call parsing")
		return toolCalls
	}
	
	// 支持多种JSON格式的解析策略
	strategies := []func(string) []ToolCall{
		s.parseDirectJSON,
		s.parseWrappedToolCall,
		s.parseCodeBlockJSON,
		s.parseMultipleToolCalls,
	}
	
	for i, strategy := range strategies {
		if parsedCalls := strategy(content); len(parsedCalls) > 0 {
			s.logger.Info("Tool calls parsed successfully", 
				zap.Int("strategy", i+1), 
				zap.Int("count", len(parsedCalls)))
			return parsedCalls
		}
	}
	
	s.logger.Warn("No tool calls found in content", zap.String("content_preview", s.truncateString(content, 100)))
	return toolCalls
}

// parseDirectJSON 尝试直接解析整个内容作为JSON
func (s *AIAssistantService) parseDirectJSON(content string) []ToolCall {
	var toolCalls []ToolCall
	
	// 尝试解析为单个工具调用
	var singleCall ToolCall
	if err := json.Unmarshal([]byte(content), &singleCall); err == nil && singleCall.Name != "" {
		if singleCall.Arguments == nil {
			singleCall.Arguments = make(map[string]interface{})
		}
		toolCalls = append(toolCalls, singleCall)
		return toolCalls
	}
	
	// 尝试解析为工具调用数组
	var multipleCalls []ToolCall
	if err := json.Unmarshal([]byte(content), &multipleCalls); err == nil && len(multipleCalls) > 0 {
		for _, call := range multipleCalls {
			if call.Name != "" {
				if call.Arguments == nil {
					call.Arguments = make(map[string]interface{})
				}
				toolCalls = append(toolCalls, call)
			}
		}
		return toolCalls
	}
	
	return toolCalls
}

// parseWrappedToolCall 解析包装在tool_call字段中的JSON
func (s *AIAssistantService) parseWrappedToolCall(content string) []ToolCall {
	var toolCalls []ToolCall
	
	var wrapper map[string]interface{}
	if err := json.Unmarshal([]byte(content), &wrapper); err != nil {
		return toolCalls
	}
	
	// 检查tool_call字段
	if toolCallData, ok := wrapper["tool_call"]; ok {
		if call := s.extractToolCallFromInterface(toolCallData); call != nil {
			toolCalls = append(toolCalls, *call)
		}
	}
	
	// 检查tool_calls字段（数组）
	if toolCallsData, ok := wrapper["tool_calls"]; ok {
		if callsArray, ok := toolCallsData.([]interface{}); ok {
			for _, callData := range callsArray {
				if call := s.extractToolCallFromInterface(callData); call != nil {
					toolCalls = append(toolCalls, *call)
				}
			}
		}
	}
	
	return toolCalls
}

// parseCodeBlockJSON 从代码块中提取JSON
func (s *AIAssistantService) parseCodeBlockJSON(content string) []ToolCall {
	var toolCalls []ToolCall
	
	// 使用正则表达式查找JSON代码块
	jsonBlockRegex := regexp.MustCompile("```(?:json)?\n?([^`]+)\n?```")
	matches := jsonBlockRegex.FindAllStringSubmatch(content, -1)
	
	for _, match := range matches {
		if len(match) > 1 {
			jsonContent := strings.TrimSpace(match[1])
			if calls := s.parseDirectJSON(jsonContent); len(calls) > 0 {
				toolCalls = append(toolCalls, calls...)
			}
		}
	}
	
	return toolCalls
}

// parseMultipleToolCalls 使用改进的括号匹配算法查找多个工具调用
func (s *AIAssistantService) parseMultipleToolCalls(content string) []ToolCall {
	var toolCalls []ToolCall
	
	// 查找所有可能的JSON对象起始位置
	patterns := []string{`{"tool_call"`, `{"name"`, `[{"name"`}
	
	for _, pattern := range patterns {
		startIndex := 0
		for {
			index := strings.Index(content[startIndex:], pattern)
			if index == -1 {
				break
			}
			
			actualIndex := startIndex + index
			if jsonStr := s.extractJSONObject(content, actualIndex); jsonStr != "" {
				// 尝试解析提取的JSON
				if calls := s.parseDirectJSON(jsonStr); len(calls) > 0 {
					toolCalls = append(toolCalls, calls...)
				} else if calls := s.parseWrappedToolCall(jsonStr); len(calls) > 0 {
					toolCalls = append(toolCalls, calls...)
				}
			}
			
			startIndex = actualIndex + 1
		}
	}
	
	return s.deduplicateToolCalls(toolCalls)
}

// extractJSONObject 从指定位置提取完整的JSON对象
func (s *AIAssistantService) extractJSONObject(content string, startIndex int) string {
	if startIndex >= len(content) {
		return ""
	}
	
	remaining := content[startIndex:]
	braceCount := 0
	inString := false
	escaped := false
	
	for i, char := range remaining {
		if escaped {
			escaped = false
			continue
		}
		
		if char == '\\' {
			escaped = true
			continue
		}
		
		if char == '"' {
			inString = !inString
			continue
		}
		
		if !inString {
			if char == '{' {
				braceCount++
			} else if char == '}' {
				braceCount--
				if braceCount == 0 {
					return remaining[:i+1]
				}
			}
		}
	}
	
	return ""
}

// extractToolCallFromInterface 从interface{}中提取ToolCall
func (s *AIAssistantService) extractToolCallFromInterface(data interface{}) *ToolCall {
	callMap, ok := data.(map[string]interface{})
	if !ok {
		return nil
	}
	
	name, ok := callMap["name"].(string)
	if !ok || name == "" {
		return nil
	}
	
	toolCall := &ToolCall{
		Name:      name,
		Arguments: make(map[string]interface{}),
	}
	
	if args, ok := callMap["arguments"].(map[string]interface{}); ok {
		toolCall.Arguments = args
	}
	
	return toolCall
}

// deduplicateToolCalls 去除重复的工具调用
func (s *AIAssistantService) deduplicateToolCalls(toolCalls []ToolCall) []ToolCall {
	seen := make(map[string]bool)
	var unique []ToolCall
	
	for _, call := range toolCalls {
		// 创建唯一标识符
		key := call.Name
		if argsBytes, err := json.Marshal(call.Arguments); err == nil {
			key += string(argsBytes)
		}
		
		if !seen[key] {
			seen[key] = true
			unique = append(unique, call)
		}
	}
	
	return unique
}

// truncateString 截断字符串用于日志记录
func (s *AIAssistantService) truncateString(str string, maxLen int) string {
	if len(str) <= maxLen {
		return str
	}
	return str[:maxLen] + "..."
}

// executeToolCall 执行工具调用
func (s *AIAssistantService) executeToolCall(ctx context.Context, toolCall ToolCall) ToolCallExecution {
	execution := ToolCallExecution{
		ToolName:  toolCall.Name,
		Arguments: toolCall.Arguments,
	}
	
	// 执行MCP工具，带有超时控制和重试机制
	mcpReq := &dto.MCPExecuteRequest{
		Name:      toolCall.Name,
		Arguments: toolCall.Arguments,
	}
	
	result, err := s.executeToolWithRetry(ctx, mcpReq, toolCall.Name)
	if err != nil {
		execution.Error = err.Error()
		s.logger.Error("Tool execution failed after retries",
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

// executeToolWithRetry 执行工具调用，带有超时控制和重试机制
func (s *AIAssistantService) executeToolWithRetry(ctx context.Context, req *dto.MCPExecuteRequest, toolName string) (*dto.MCPExecuteResponse, error) {
	const (
		maxRetries = 3
		baseDelay  = 1 * time.Second
		maxDelay   = 10 * time.Second
		timeout    = 30 * time.Second
	)
	
	var lastErr error
	
	for attempt := 0; attempt < maxRetries; attempt++ {
		// 为每次尝试创建带超时的上下文
		timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
		
		s.logger.Info("Executing tool",
			zap.String("tool", toolName),
			zap.Int("attempt", attempt+1),
			zap.Int("max_attempts", maxRetries))
		
		result, err := s.mcpClient.ExecuteTool(timeoutCtx, req)
		cancel() // 立即释放资源
		
		if err == nil {
			if result != nil && !result.IsError {
				// 成功执行
				if attempt > 0 {
					s.logger.Info("Tool execution succeeded after retry",
						zap.String("tool", toolName),
						zap.Int("attempt", attempt+1))
				}
				return result, nil
			} else if result != nil && result.IsError {
				// 工具返回了错误结果，但这不是网络或系统错误
				errorContent := ""
				if len(result.Content) > 0 {
					if contentBytes, err := json.Marshal(result.Content); err == nil {
						errorContent = string(contentBytes)
					}
				}
				s.logger.Warn("Tool returned error result",
					zap.String("tool", toolName),
					zap.String("error", errorContent))
				return result, nil
			}
		}
		
		lastErr = err
		
		// 检查是否应该重试
		if !s.shouldRetryError(err) {
			s.logger.Warn("Error is not retryable, stopping attempts",
				zap.String("tool", toolName),
				zap.Error(err))
			break
		}
		
		// 如果不是最后一次尝试，等待后重试
		if attempt < maxRetries-1 {
			delay := s.calculateBackoffDelay(attempt, baseDelay, maxDelay)
			s.logger.Info("Tool execution failed, retrying",
				zap.String("tool", toolName),
				zap.Int("attempt", attempt+1),
				zap.Duration("retry_delay", delay),
				zap.Error(err))
			
			select {
			case <-time.After(delay):
				// 继续重试
			case <-ctx.Done():
				// 上下文被取消
				return nil, ctx.Err()
			}
		}
	}
	
	return nil, fmt.Errorf("tool execution failed after %d attempts: %w", maxRetries, lastErr)
}

// shouldRetryError 判断错误是否应该重试
func (s *AIAssistantService) shouldRetryError(err error) bool {
	if err == nil {
		return false
	}
	
	errStr := err.Error()
	
	// 网络相关错误通常可以重试
	retryableErrors := []string{
		"timeout",
		"connection refused",
		"connection reset",
		"temporary failure",
		"network is unreachable",
		"no such host",
		"context deadline exceeded",
		"i/o timeout",
		"EOF",
	}
	
	for _, retryableErr := range retryableErrors {
		if strings.Contains(strings.ToLower(errStr), retryableErr) {
			return true
		}
	}
	
	// 检查是否是上下文超时
	if err == context.DeadlineExceeded || err == context.Canceled {
		return true
	}
	
	return false
}

// calculateBackoffDelay 计算指数退避延迟
func (s *AIAssistantService) calculateBackoffDelay(attempt int, baseDelay, maxDelay time.Duration) time.Duration {
	// 指数退避：baseDelay * 2^attempt
	delay := baseDelay * time.Duration(1<<uint(attempt))
	
	// 添加一些随机性以避免雷群效应
	jitter := time.Duration(float64(delay) * 0.1 * (0.5 - float64(attempt%2)))
	delay += jitter
	
	// 确保不超过最大延迟
	if delay > maxDelay {
		delay = maxDelay
	}
	
	return delay
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
	resultsBuilder.WriteString("## Tool Execution Results\n\n")
	
	successCount := 0
	errorCount := 0
	
	for i, exec := range executions {
		resultsBuilder.WriteString(fmt.Sprintf("### Tool %d: %s\n", i+1, exec.ToolName))
		
		// 添加工具参数信息
		if len(exec.Arguments) > 0 {
			if argsBytes, err := json.Marshal(exec.Arguments); err == nil {
				resultsBuilder.WriteString(fmt.Sprintf("**Parameters:** %s\n", string(argsBytes)))
			}
		}
		
		if exec.Error != "" {
			resultsBuilder.WriteString(fmt.Sprintf("**Status:** ❌ Error\n"))
			resultsBuilder.WriteString(fmt.Sprintf("**Error Details:** %s\n", exec.Error))
			errorCount++
		} else if exec.Result != nil {
			if exec.Result.IsError {
				resultsBuilder.WriteString(fmt.Sprintf("**Status:** ⚠️ Tool Error\n"))
				errorCount++
			} else {
				resultsBuilder.WriteString(fmt.Sprintf("**Status:** ✅ Success\n"))
				successCount++
			}
			
			resultsBuilder.WriteString("**Results:**\n")
			for _, content := range exec.Result.Content {
				resultsBuilder.WriteString(fmt.Sprintf("- %s\n", content.Text))
			}
		}
		resultsBuilder.WriteString("\n")
	}
	
	// 构建提供商请求的消息格式
	providerMessages := make([]ProviderMessage, 0, len(originalReq.Messages)+3)
	
	// 添加系统消息，定义分析师角色
	systemPrompt := s.buildAnalysisSystemPrompt(successCount, errorCount)
	providerMessages = append(providerMessages, ProviderMessage{
		Role:    "system",
		Content: systemPrompt,
	})
	
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
	
	// 添加生成最终回复的详细指令
	analysisPrompt := s.buildAnalysisPrompt(executions)
	providerMessages = append(providerMessages, ProviderMessage{
		Role:    "user",
		Content: analysisPrompt,
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

// buildAnalysisSystemPrompt 构建分析系统提示
func (s *AIAssistantService) buildAnalysisSystemPrompt(successCount, errorCount int) string {
	var builder strings.Builder
	builder.WriteString("You are a professional financial analyst with expertise in stock market analysis, investment strategies, and financial data interpretation. ")
	builder.WriteString("Your role is to provide comprehensive, data-driven financial analysis based on the tool execution results.\n\n")
	
	builder.WriteString("## Analysis Guidelines:\n")
	builder.WriteString("1. **Data Interpretation**: Analyze the numerical data, trends, and patterns from the tool results\n")
	builder.WriteString("2. **Context Integration**: Consider market conditions, company fundamentals, and industry trends\n")
	builder.WriteString("3. **Risk Assessment**: Identify potential risks and opportunities\n")
	builder.WriteString("4. **Professional Tone**: Use clear, professional language suitable for investors\n")
	builder.WriteString("5. **Actionable Insights**: Provide practical recommendations when appropriate\n\n")
	
	if errorCount > 0 {
		builder.WriteString("⚠️ **Note**: Some tools encountered errors. Acknowledge these limitations in your analysis and work with available data.\n\n")
	}
	
	return builder.String()
}

// buildAnalysisPrompt 构建分析提示
func (s *AIAssistantService) buildAnalysisPrompt(executions []ToolCallExecution) string {
	var builder strings.Builder
	
	builder.WriteString("Based on the tool execution results above, please provide a comprehensive financial analysis report with the following structure:\n\n")
	
	builder.WriteString("## 📊 Executive Summary\n")
	builder.WriteString("Provide a concise overview of the key findings and main insights.\n\n")
	
	builder.WriteString("## 📈 Data Analysis\n")
	builder.WriteString("Analyze the specific data points, metrics, and trends from the tool results. Include:\n")
	builder.WriteString("- Key financial metrics and their implications\n")
	builder.WriteString("- Trend analysis and patterns\n")
	builder.WriteString("- Comparative analysis (if applicable)\n\n")
	
	builder.WriteString("## 🎯 Investment Insights\n")
	builder.WriteString("Provide investment-focused analysis including:\n")
	builder.WriteString("- Market position and competitive advantages\n")
	builder.WriteString("- Growth prospects and potential catalysts\n")
	builder.WriteString("- Valuation considerations\n\n")
	
	builder.WriteString("## ⚠️ Risk Factors\n")
	builder.WriteString("Identify and explain potential risks and challenges.\n\n")
	
	// 根据工具类型添加特定指导
	toolTypes := make(map[string]bool)
	for _, exec := range executions {
		toolTypes[exec.ToolName] = true
	}
	
	if toolTypes["stock_comparison"] {
		builder.WriteString("## 🔄 Comparative Analysis\n")
		builder.WriteString("Provide detailed comparison between the analyzed stocks, highlighting relative strengths and weaknesses.\n\n")
	}
	
	if toolTypes["yahoo_finance"] || toolTypes["stock_analysis"] {
		builder.WriteString("## 📊 Technical & Fundamental Analysis\n")
		builder.WriteString("Combine technical indicators with fundamental analysis for a comprehensive view.\n\n")
	}
	
	builder.WriteString("## 💡 Recommendations\n")
	builder.WriteString("Provide clear, actionable recommendations based on your analysis. Include:\n")
	builder.WriteString("- Investment thesis (if applicable)\n")
	builder.WriteString("- Suggested actions or considerations\n")
	builder.WriteString("- Timeline and monitoring points\n\n")
	
	builder.WriteString("**Important**: Ensure your analysis is objective, data-driven, and acknowledges any limitations from tool errors or missing data. ")
	builder.WriteString("Use professional financial terminology and provide context for technical concepts when necessary.")
	
	return builder.String()
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