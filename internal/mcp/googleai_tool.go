package mcp

import (
	"context"
	"fmt"

	"admin/internal/dto"
	"admin/internal/googleai"
)

// GoogleAIService 接口，避免导入循环
type GoogleAIService interface {
	ChatCompletion(ctx context.Context, req *GoogleAIChatCompletionRequest) (*GoogleAIChatCompletionResponse, error)
	ListModels(ctx context.Context) (map[string]*googleai.ModelConfig, error)
	ValidateAPIKey(ctx context.Context) error
	SetAPIKey(key string) error
	GetModelConfig(name string) (*googleai.ModelConfig, error)
	EnableModel(name string) error
	DisableModel(name string) error
}

// GoogleAIChatCompletionRequest Google AI 聊天完成请求
type GoogleAIChatCompletionRequest struct {
	Model       string                `json:"model"`
	Messages    []googleai.Message    `json:"messages"`
	MaxTokens   *int                  `json:"max_tokens,omitempty"`
	Temperature *float32              `json:"temperature,omitempty"`
	TopP        *float32              `json:"top_p,omitempty"`
	TopK        *int                  `json:"top_k,omitempty"`
	Stream      bool                  `json:"stream,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
}

// GoogleAIChatCompletionResponse Google AI 聊天完成响应
type GoogleAIChatCompletionResponse struct {
	ID      string            `json:"id"`
	Object  string            `json:"object"`
	Created int64             `json:"created"`
	Model   string            `json:"model"`
	Choices []googleai.Choice `json:"choices"`
	Usage   googleai.Usage    `json:"usage"`
}

// GoogleAIChatTool Google AI 聊天工具
type GoogleAIChatTool struct {
	*BaseTool
	googleaiService GoogleAIService
}

// NewGoogleAIChatTool 创建 Google AI 聊天工具
func NewGoogleAIChatTool(googleaiService GoogleAIService) *GoogleAIChatTool {
	baseTool := &BaseTool{
		Name:        "googleai_chat",
		Description: "Chat with Google AI models (Gemini Pro, Gemini Pro Vision, etc.)",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"model": map[string]interface{}{
					"type":        "string",
					"description": "The Google AI model to use (e.g., gemini-pro, gemini-pro-vision)",
					"default":     "gemini-pro",
				},
				"messages": map[string]interface{}{
					"type":        "array",
					"description": "Array of message objects",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"role": map[string]interface{}{
								"type":        "string",
								"description": "Message role (user, model)",
								"enum":        []string{"user", "model"},
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
					"maximum":     8192,
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
				"top_k": map[string]interface{}{
					"type":        "number",
					"description": "Top-k sampling (1-40, optional)",
					"minimum":     1,
					"maximum":     40,
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
	
	return &GoogleAIChatTool{
		BaseTool:        baseTool,
		googleaiService: googleaiService,
	}
}

// Validate 验证参数
func (gct *GoogleAIChatTool) Validate(args map[string]interface{}) error {
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

		if role != "user" && role != "model" {
			return fmt.Errorf("message at index %d has invalid role: %s", i, role)
		}

		content, ok := msgObj["content"].(string)
		if !ok {
			return fmt.Errorf("message at index %d must have a content field", i)
		}

		if content == "" {
			return fmt.Errorf("message at index %d content cannot be empty", i)
		}
	}

	// 验证可选参数
	if maxTokens, ok := args["max_tokens"]; ok {
		if maxTokensFloat, ok := maxTokens.(float64); ok {
			if maxTokensFloat < 1 || maxTokensFloat > 8192 {
				return fmt.Errorf("max_tokens must be between 1 and 8192")
			}
		} else {
			return fmt.Errorf("max_tokens must be a number")
		}
	}

	if temperature, ok := args["temperature"]; ok {
		if tempFloat, ok := temperature.(float64); ok {
			if tempFloat < 0 || tempFloat > 2 {
				return fmt.Errorf("temperature must be between 0 and 2")
			}
		} else {
			return fmt.Errorf("temperature must be a number")
		}
	}

	if topP, ok := args["top_p"]; ok {
		if topPFloat, ok := topP.(float64); ok {
			if topPFloat < 0 || topPFloat > 1 {
				return fmt.Errorf("top_p must be between 0 and 1")
			}
		} else {
			return fmt.Errorf("top_p must be a number")
		}
	}

	if topK, ok := args["top_k"]; ok {
		if topKFloat, ok := topK.(float64); ok {
			if topKFloat < 1 || topKFloat > 40 {
				return fmt.Errorf("top_k must be between 1 and 40")
			}
		} else {
			return fmt.Errorf("top_k must be a number")
		}
	}

	return nil
}

// Execute 执行聊天
func (gct *GoogleAIChatTool) Execute(ctx context.Context, args map[string]interface{}) (*dto.MCPExecuteResponse, error) {
	// 验证参数
	if err := gct.Validate(args); err != nil {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{{
				Type: "text",
				Text: err.Error(),
			}},
			IsError: true,
		}, nil
	}

	// 构建请求
	req := &GoogleAIChatCompletionRequest{}

	// 设置模型
	if model, ok := args["model"].(string); ok {
		req.Model = model
	} else {
		req.Model = "gemini-pro" // 默认模型
	}

	// 设置消息
	messagesArray := args["messages"].([]interface{})
	req.Messages = make([]googleai.Message, len(messagesArray))
	for i, msg := range messagesArray {
		msgObj := msg.(map[string]interface{})
		req.Messages[i] = googleai.Message{
			Role:    msgObj["role"].(string),
			Content: msgObj["content"].(string),
		}
	}

	// 设置可选参数
	if maxTokens, ok := args["max_tokens"].(float64); ok {
		maxTokensInt := int(maxTokens)
		req.MaxTokens = &maxTokensInt
	}

	if temperature, ok := args["temperature"].(float64); ok {
		tempFloat32 := float32(temperature)
		req.Temperature = &tempFloat32
	}

	if topP, ok := args["top_p"].(float64); ok {
		topPFloat32 := float32(topP)
		req.TopP = &topPFloat32
	}

	if topK, ok := args["top_k"].(float64); ok {
		topKInt := int(topK)
		req.TopK = &topKInt
	}

	if stream, ok := args["stream"].(bool); ok {
		req.Stream = stream
	}

	// 调用 Google AI 服务
	resp, err := gct.googleaiService.ChatCompletion(ctx, req)
	if err != nil {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("Google AI API error: %v", err),
			}},
			IsError: true,
		}, nil
	}

	return &dto.MCPExecuteResponse{
		Content: []dto.MCPContent{{
			Type: "text",
			Data: resp,
		}},
	}, nil
}

// GoogleAIModelsTool Google AI 模型工具
type GoogleAIModelsTool struct {
	*BaseTool
	googleaiService GoogleAIService
}

// NewGoogleAIModelsTool 创建 Google AI 模型工具
func NewGoogleAIModelsTool(googleaiService GoogleAIService) *GoogleAIModelsTool {
	baseTool := &BaseTool{
		Name:        "googleai_models",
		Description: "List and manage Google AI models",
		InputSchema: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
	}
	
	return &GoogleAIModelsTool{
		BaseTool:        baseTool,
		googleaiService: googleaiService,
	}
}

// Execute 执行模型列表
func (gmt *GoogleAIModelsTool) Execute(ctx context.Context, args map[string]interface{}) (*dto.MCPExecuteResponse, error) {
	models, err := gmt.googleaiService.ListModels(ctx)
	if err != nil {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("Failed to list models: %v", err),
			}},
			IsError: true,
		}, nil
	}

	// 转换为响应格式
	modelList := make([]map[string]interface{}, 0, len(models))
	for name, config := range models {
		modelList = append(modelList, map[string]interface{}{
			"name":         name,
			"display_name": config.DisplayName,
			"max_tokens":   config.MaxTokens,
			"temperature":  config.Temperature,
			"top_p":        config.TopP,
			"top_k":        config.TopK,
			"enabled":      config.Enabled,
		})
	}

	return &dto.MCPExecuteResponse{
		Content: []dto.MCPContent{{
			Type: "text",
			Data: map[string]interface{}{
				"models": modelList,
				"count":  len(modelList),
			},
		}},
	}, nil
}

// GoogleAIConfigTool Google AI 配置工具
type GoogleAIConfigTool struct {
	*BaseTool
	googleaiService GoogleAIService
}

// NewGoogleAIConfigTool 创建 Google AI 配置工具
func NewGoogleAIConfigTool(googleaiService GoogleAIService) *GoogleAIConfigTool {
	baseTool := &BaseTool{
		Name:        "googleai_config",
		Description: "Configure Google AI models and settings",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"action": map[string]interface{}{
					"type":        "string",
					"description": "Action to perform",
					"enum":        []string{"get_model", "enable_model", "disable_model", "set_api_key", "validate_api_key"},
				},
				"model": map[string]interface{}{
					"type":        "string",
					"description": "Model name (required for model actions)",
				},
				"api_key": map[string]interface{}{
					"type":        "string",
					"description": "API key (required for set_api_key action)",
				},
			},
			"required": []string{"action"},
		},
	}
	
	return &GoogleAIConfigTool{
		BaseTool:        baseTool,
		googleaiService: googleaiService,
	}
}

// Validate 验证参数
func (gct *GoogleAIConfigTool) Validate(args map[string]interface{}) error {
	action, ok := args["action"].(string)
	if !ok {
		return fmt.Errorf("action parameter is required")
	}

	switch action {
	case "get_model", "enable_model", "disable_model":
		if _, ok := args["model"].(string); !ok {
			return fmt.Errorf("model parameter is required for %s action", action)
		}
	case "set_api_key":
		if _, ok := args["api_key"].(string); !ok {
			return fmt.Errorf("api_key parameter is required for set_api_key action")
		}
	case "validate_api_key":
		// 不需要额外参数
	default:
		return fmt.Errorf("invalid action: %s", action)
	}

	return nil
}

// Execute 执行配置操作
func (gct *GoogleAIConfigTool) Execute(ctx context.Context, args map[string]interface{}) (*dto.MCPExecuteResponse, error) {
	// 验证参数
	if err := gct.Validate(args); err != nil {
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{{
				Type: "text",
				Text: err.Error(),
			}},
			IsError: true,
		}, nil
	}

	action := args["action"].(string)

	switch action {
	case "get_model":
		model := args["model"].(string)
		config, err := gct.googleaiService.GetModelConfig(model)
		if err != nil {
			return &dto.MCPExecuteResponse{
				Content: []dto.MCPContent{{
					Type: "text",
					Text: fmt.Sprintf("Failed to get model config: %v", err),
				}},
				IsError: true,
			}, nil
		}

		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{{
				Type: "text",
				Data: map[string]interface{}{
					"name":         config.Name,
					"display_name": config.DisplayName,
					"max_tokens":   config.MaxTokens,
					"temperature":  config.Temperature,
					"top_p":        config.TopP,
					"top_k":        config.TopK,
					"enabled":      config.Enabled,
				},
			}},
		}, nil

	case "enable_model":
		model := args["model"].(string)
		err := gct.googleaiService.EnableModel(model)
		if err != nil {
			return &dto.MCPExecuteResponse{
				Content: []dto.MCPContent{{
					Type: "text",
					Text: fmt.Sprintf("Failed to enable model: %v", err),
				}},
				IsError: true,
			}, nil
		}

		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{{
				Type: "text",
				Data: map[string]interface{}{
					"message": fmt.Sprintf("Model %s enabled successfully", model),
				},
			}},
		}, nil

	case "disable_model":
		model := args["model"].(string)
		err := gct.googleaiService.DisableModel(model)
		if err != nil {
			return &dto.MCPExecuteResponse{
				Content: []dto.MCPContent{{
					Type: "text",
					Text: fmt.Sprintf("Failed to disable model: %v", err),
				}},
				IsError: true,
			}, nil
		}

		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{{
				Type: "text",
				Data: map[string]interface{}{
					"message": fmt.Sprintf("Model %s disabled successfully", model),
				},
			}},
		}, nil

	case "set_api_key":
		apiKey := args["api_key"].(string)
		err := gct.googleaiService.SetAPIKey(apiKey)
		if err != nil {
			return &dto.MCPExecuteResponse{
				Content: []dto.MCPContent{{
					Type: "text",
					Text: fmt.Sprintf("Failed to set API key: %v", err),
				}},
				IsError: true,
			}, nil
		}

		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{{
				Type: "text",
				Data: map[string]interface{}{
					"message": "API key set successfully",
				},
			}},
		}, nil

	case "validate_api_key":
		err := gct.googleaiService.ValidateAPIKey(ctx)
		if err != nil {
			return &dto.MCPExecuteResponse{
				Content: []dto.MCPContent{{
					Type: "text",
					Text: fmt.Sprintf("API key validation failed: %v", err),
				}},
				IsError: true,
			}, nil
		}

		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{{
				Type: "text",
				Data: map[string]interface{}{
					"message": "API key is valid",
				},
			}},
		}, nil

	default:
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{{
				Type: "text",
				Text: fmt.Sprintf("Unknown action: %s", action),
			}},
			IsError: true,
		}, nil
	}
}