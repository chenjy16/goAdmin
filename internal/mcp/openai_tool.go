package mcp

import (
	"context"
	"fmt"

	"go-springAi/internal/dto"
	"go-springAi/internal/openai"
)

// OpenAIService 接口，避免导入循环
type OpenAIService interface {
	ChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error)
	ListModels(ctx context.Context) (map[string]*openai.ModelConfig, error)
	ValidateAPIKey(ctx context.Context) error
	SetAPIKey(key string) error
	GetModelConfig(name string) (*openai.ModelConfig, error)
	EnableModel(name string) error
	DisableModel(name string) error
}

// ChatCompletionRequest 聊天完成请求
type ChatCompletionRequest struct {
	Model       string                `json:"model"`
	Messages    []openai.Message      `json:"messages"`
	MaxTokens   *int                  `json:"max_tokens,omitempty"`
	Temperature *float32              `json:"temperature,omitempty"`
	TopP        *float32              `json:"top_p,omitempty"`
	Stream      bool                  `json:"stream,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
}

// ChatCompletionResponse 聊天完成响应
type ChatCompletionResponse struct {
	ID      string           `json:"id"`
	Object  string           `json:"object"`
	Created int64            `json:"created"`
	Model   string           `json:"model"`
	Choices []openai.Choice  `json:"choices"`
	Usage   openai.Usage     `json:"usage"`
}

// OpenAIChatTool OpenAI 聊天工具
type OpenAIChatTool struct {
	*BaseTool
	openaiService OpenAIService
}

// NewOpenAIChatTool 创建 OpenAI 聊天工具
func NewOpenAIChatTool(openaiService OpenAIService) *OpenAIChatTool {
	baseTool := &BaseTool{
		Name:        "openai_chat",
		Description: "Chat with OpenAI models (GPT-4, GPT-3.5-turbo, etc.)",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"model": map[string]interface{}{
					"type":        "string",
					"description": "The OpenAI model to use (e.g., gpt-4, gpt-3.5-turbo)",
					"default":     "gpt-3.5-turbo",
				},
				"messages": map[string]interface{}{
					"type":        "array",
					"description": "Array of message objects",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"role": map[string]interface{}{
								"type":        "string",
								"description": "Message role (system, user, assistant)",
								"enum":        []string{"system", "user", "assistant"},
							},
							"content": map[string]interface{}{
								"type":        "string",
								"description": "Message content",
							},
						},
						"required": []string{"role", "content"},
					},
				},
				"max_tokens": map[string]interface{}{
					"type":        "number",
					"description": "Maximum number of tokens to generate (optional)",
					"minimum":     1,
					"maximum":     4096,
				},
				"temperature": map[string]interface{}{
					"type":        "number",
					"description": "Sampling temperature (0-2, optional)",
					"minimum":     0,
					"maximum":     2,
				},
				"top_p": map[string]interface{}{
					"type":        "number",
					"description": "Top-p sampling (0-1, optional)",
					"minimum":     0,
					"maximum":     1,
				},
				"stream": map[string]interface{}{
					"type":        "boolean",
					"description": "Whether to stream the response (optional)",
					"default":     false,
				},
			},
			"required": []string{"messages"},
		},
	}
	
	return &OpenAIChatTool{
		BaseTool:      baseTool,
		openaiService: openaiService,
	}
}

// Validate 验证参数
func (oct *OpenAIChatTool) Validate(args map[string]interface{}) error {
	// 验证 messages 参数
	messages, ok := args["messages"]
	if !ok {
		return fmt.Errorf("messages parameter is required")
	}

	messagesArray, ok := messages.([]interface{})
	if !ok {
		return fmt.Errorf("messages must be an array")
	}

	if len(messagesArray) == 0 {
		return fmt.Errorf("messages array cannot be empty")
	}

	// 验证每个消息对象
	for i, msg := range messagesArray {
		msgObj, ok := msg.(map[string]interface{})
		if !ok {
			return fmt.Errorf("message at index %d must be an object", i)
		}

		role, ok := msgObj["role"].(string)
		if !ok {
			return fmt.Errorf("message at index %d must have a role field", i)
		}

		if role != "system" && role != "user" && role != "assistant" {
			return fmt.Errorf("message at index %d has invalid role: %s", i, role)
		}

		content, ok := msgObj["content"].(string)
		if !ok {
			return fmt.Errorf("message at index %d must have a content field", i)
		}

		if content == "" {
			return fmt.Errorf("message at index %d cannot have empty content", i)
		}
	}

	// 验证可选参数
	if maxTokens, ok := args["max_tokens"]; ok {
		if _, ok := maxTokens.(float64); !ok {
			return fmt.Errorf("max_tokens must be a number")
		}
	}

	if temperature, ok := args["temperature"]; ok {
		if temp, ok := temperature.(float64); !ok || temp < 0 || temp > 2 {
			return fmt.Errorf("temperature must be a number between 0 and 2")
		}
	}

	if topP, ok := args["top_p"]; ok {
		if tp, ok := topP.(float64); !ok || tp < 0 || tp > 1 {
			return fmt.Errorf("top_p must be a number between 0 and 1")
		}
	}

	return nil
}

// Execute 执行 OpenAI 聊天工具
func (oct *OpenAIChatTool) Execute(ctx context.Context, args map[string]interface{}) (*dto.MCPExecuteResponse, error) {
	// 构建请求
	req := &ChatCompletionRequest{
		Model:    "gpt-3.5-turbo", // 默认模型
		Messages: []openai.Message{},
		Stream:   false,
	}

	// 设置模型
	if model, ok := args["model"].(string); ok {
		req.Model = model
	}

	// 设置消息
	messagesArray := args["messages"].([]interface{})
	for _, msg := range messagesArray {
		msgObj := msg.(map[string]interface{})
		req.Messages = append(req.Messages, openai.Message{
			Role:    msgObj["role"].(string),
			Content: msgObj["content"].(string),
		})
	}

	// 设置可选参数
	if maxTokens, ok := args["max_tokens"].(float64); ok {
		maxTokensInt := int(maxTokens)
		req.MaxTokens = &maxTokensInt
	}

	if temperature, ok := args["temperature"].(float64); ok {
		tempFloat := float32(temperature)
		req.Temperature = &tempFloat
	}

	if topP, ok := args["top_p"].(float64); ok {
		topPFloat := float32(topP)
		req.TopP = &topPFloat
	}

	if stream, ok := args["stream"].(bool); ok {
		req.Stream = stream
	}

	// 如果是流式请求，返回特殊响应
	if req.Stream {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: "Streaming is not supported in MCP tool execution. Use the API endpoint for streaming.",
				},
			},
			IsError: true,
		}, nil
	}

	// 执行聊天完成
	response, err := oct.openaiService.ChatCompletion(ctx, req)
	if err != nil {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("OpenAI API error: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// 构建响应内容
	var responseText string
	if len(response.Choices) > 0 {
		responseText = response.Choices[0].Message.Content
	} else {
		responseText = "No response generated"
	}

	return &dto.MCPExecuteResponse{
		Content: []dto.MCPContent{
			{
				Type: "text",
				Text: responseText,
				Data: map[string]interface{}{
					"model":             response.Model,
					"id":                response.ID,
					"usage":             response.Usage,
					"finish_reason":     response.Choices[0].FinishReason,
					"prompt_tokens":     response.Usage.PromptTokens,
					"completion_tokens": response.Usage.CompletionTokens,
					"total_tokens":      response.Usage.TotalTokens,
				},
			},
		},
		IsError: false,
	}, nil
}

// OpenAIModelsTool OpenAI 模型列表工具
type OpenAIModelsTool struct {
	*BaseTool
	openaiService OpenAIService
}

// NewOpenAIModelsTool 创建 OpenAI 模型列表工具
func NewOpenAIModelsTool(openaiService OpenAIService) *OpenAIModelsTool {
	baseTool := &BaseTool{
		Name:        "openai_models",
		Description: "List available OpenAI models and their configurations",
		InputSchema: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
	}
	
	return &OpenAIModelsTool{
		BaseTool:      baseTool,
		openaiService: openaiService,
	}
}

// Execute 执行 OpenAI 模型列表工具
func (omt *OpenAIModelsTool) Execute(ctx context.Context, args map[string]interface{}) (*dto.MCPExecuteResponse, error) {
	models, err := omt.openaiService.ListModels(ctx)
	if err != nil {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Error listing models: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// 构建模型信息文本
	var modelInfo string
	modelCount := 0
	for name, config := range models {
		if config.Enabled {
			modelInfo += fmt.Sprintf("- %s (max_tokens: %d, temperature: %.2f)\n",
				name, config.MaxTokens, config.Temperature)
			modelCount++
		}
	}

	if modelCount == 0 {
		modelInfo = "No enabled models found"
	}

	return &dto.MCPExecuteResponse{
		Content: []dto.MCPContent{
			{
				Type: "text",
				Text: fmt.Sprintf("Available OpenAI models (%d):\n%s", modelCount, modelInfo),
				Data: models,
			},
		},
		IsError: false,
	}, nil
}

// OpenAIConfigTool OpenAI 配置工具
type OpenAIConfigTool struct {
	*BaseTool
	openaiService OpenAIService
}

// NewOpenAIConfigTool 创建 OpenAI 配置工具
func NewOpenAIConfigTool(openaiService OpenAIService) *OpenAIConfigTool {
	baseTool := &BaseTool{
		Name:        "openai_config",
		Description: "Manage OpenAI configuration (set API key, enable/disable models)",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"action": map[string]interface{}{
					"type":        "string",
					"description": "Action to perform",
					"enum":        []string{"set_api_key", "validate_key", "enable_model", "disable_model", "get_model_config"},
				},
				"api_key": map[string]interface{}{
					"type":        "string",
					"description": "OpenAI API key (for set_api_key action)",
				},
				"model": map[string]interface{}{
					"type":        "string",
					"description": "Model name (for model-related actions)",
				},
			},
			"required": []string{"action"},
		},
	}
	
	return &OpenAIConfigTool{
		BaseTool:      baseTool,
		openaiService: openaiService,
	}
}

// Validate 验证参数
func (oct *OpenAIConfigTool) Validate(args map[string]interface{}) error {
	action, ok := args["action"].(string)
	if !ok {
		return fmt.Errorf("action parameter is required")
	}

	switch action {
	case "set_api_key":
		if _, ok := args["api_key"].(string); !ok {
			return fmt.Errorf("api_key parameter is required for set_api_key action")
		}
	case "enable_model", "disable_model", "get_model_config":
		if _, ok := args["model"].(string); !ok {
			return fmt.Errorf("model parameter is required for %s action", action)
		}
	case "validate_key":
		// 不需要额外参数
	default:
		return fmt.Errorf("invalid action: %s", action)
	}

	return nil
}

// Execute 执行 OpenAI 配置工具
func (oct *OpenAIConfigTool) Execute(ctx context.Context, args map[string]interface{}) (*dto.MCPExecuteResponse, error) {
	action := args["action"].(string)

	switch action {
	case "set_api_key":
		apiKey := args["api_key"].(string)
		err := oct.openaiService.SetAPIKey(apiKey)
		if err != nil {
			return &dto.MCPExecuteResponse{
				Content: []dto.MCPContent{
					{
						Type: "text",
						Text: fmt.Sprintf("Failed to set API key: %v", err),
					},
				},
				IsError: true,
			}, nil
		}

		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: "API key set successfully",
				},
			},
			IsError: false,
		}, nil

	case "validate_key":
		err := oct.openaiService.ValidateAPIKey(ctx)
		if err != nil {
			return &dto.MCPExecuteResponse{
				Content: []dto.MCPContent{
					{
						Type: "text",
						Text: fmt.Sprintf("API key validation failed: %v", err),
					},
				},
				IsError: true,
			}, nil
		}

		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: "API key is valid",
				},
			},
			IsError: false,
		}, nil

	case "enable_model":
		model := args["model"].(string)
		err := oct.openaiService.EnableModel(model)
		if err != nil {
			return &dto.MCPExecuteResponse{
				Content: []dto.MCPContent{
					{
						Type: "text",
						Text: fmt.Sprintf("Failed to enable model %s: %v", model, err),
					},
				},
				IsError: true,
			}, nil
		}

		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Model %s enabled successfully", model),
				},
			},
			IsError: false,
		}, nil

	case "disable_model":
		model := args["model"].(string)
		err := oct.openaiService.DisableModel(model)
		if err != nil {
			return &dto.MCPExecuteResponse{
				Content: []dto.MCPContent{
					{
						Type: "text",
						Text: fmt.Sprintf("Failed to disable model %s: %v", model, err),
					},
				},
				IsError: true,
			}, nil
		}

		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Model %s disabled successfully", model),
				},
			},
			IsError: false,
		}, nil

	case "get_model_config":
		model := args["model"].(string)
		config, err := oct.openaiService.GetModelConfig(model)
		if err != nil {
			return &dto.MCPExecuteResponse{
				Content: []dto.MCPContent{
					{
						Type: "text",
						Text: fmt.Sprintf("Failed to get model config for %s: %v", model, err),
					},
				},
				IsError: true,
			}, nil
		}

		configText := fmt.Sprintf("Model: %s\nEnabled: %t\nMax Tokens: %d\nTemperature: %.2f\nTop P: %.2f\nFrequency Penalty: %.2f\nPresence Penalty: %.2f",
			config.Name, config.Enabled, config.MaxTokens, config.Temperature, config.TopP, config.FrequencyPenalty, config.PresencePenalty)

		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: configText,
					Data: config,
				},
			},
			IsError: false,
		}, nil

	default:
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Unknown action: %s", action),
				},
			},
			IsError: true,
		}, nil
	}
}