package googleai

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"google.golang.org/genai"
)

// HTTPClient Google AI HTTP 客户端实现
type HTTPClient struct {
	config *Config
	client *genai.Client
}

// NewHTTPClient 创建新的 HTTP 客户端
func NewHTTPClient(config *Config) (*HTTPClient, error) {
	return &HTTPClient{
		config: config,
		client: nil, // 延迟初始化
	}, nil
}

// ensureClient 确保客户端已初始化
func (c *HTTPClient) ensureClient(ctx context.Context) error {
	if c.client != nil {
		return nil
	}

	if c.config.APIKey == "" {
		return fmt.Errorf("Google AI API key is required")
	}

	// 创建 Google AI 客户端
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  c.config.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return fmt.Errorf("create Google AI client: %w", err)
	}

	c.client = client
	return nil
}

// ChatCompletion 实现聊天完成
func (c *HTTPClient) ChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// 确保客户端已初始化
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	// 设置默认值
	if req.Model == "" {
		req.Model = c.config.DefaultModel
	}

	// 构建内容
	var contents []*genai.Content
	for _, msg := range req.Messages {
		role := genai.RoleUser
		if msg.Role == "assistant" || msg.Role == "model" {
			role = genai.RoleModel
		}
		
		content := &genai.Content{
			Role:  role,
			Parts: []*genai.Part{{Text: msg.Content}},
		}
		contents = append(contents, content)
	}

	// 构建生成配置
	config := &genai.GenerateContentConfig{}
	if req.Temperature > 0 {
		config.Temperature = &req.Temperature
	}
	if req.TopP > 0 {
		config.TopP = &req.TopP
	}
	if req.TopK > 0 {
		topK := float32(req.TopK)
		config.TopK = &topK
	}
	if req.MaxTokens > 0 {
		config.MaxOutputTokens = int32(req.MaxTokens)
	}

	// 生成内容
	resp, err := c.client.Models.GenerateContent(ctx, req.Model, contents, config)
	if err != nil {
		return nil, fmt.Errorf("generate content: %w", err)
	}

	// 转换响应格式
	chatResp := &ChatResponse{
		ID:      fmt.Sprintf("chatcmpl-%d", time.Now().Unix()),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: make([]Choice, 0),
		Usage: Usage{
			PromptTokens:     int(resp.UsageMetadata.PromptTokenCount),
			CompletionTokens: int(resp.UsageMetadata.CandidatesTokenCount),
			TotalTokens:      int(resp.UsageMetadata.TotalTokenCount),
		},
	}

	// 处理候选响应
	for i, candidate := range resp.Candidates {
		if candidate.Content != nil {
			var content strings.Builder
			for _, part := range candidate.Content.Parts {
				if part.Text != "" {
					content.WriteString(part.Text)
				}
			}

			choice := Choice{
				Index: i,
				Message: Message{
					Role:    "assistant",
					Content: content.String(),
				},
				FinishReason: string(candidate.FinishReason),
			}
			chatResp.Choices = append(chatResp.Choices, choice)
		}
	}

	return chatResp, nil
}

// ChatCompletionStream 实现流式聊天完成
func (c *HTTPClient) ChatCompletionStream(ctx context.Context, req *ChatRequest) (io.ReadCloser, error) {
	// 确保客户端已初始化
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	// 设置默认值
	if req.Model == "" {
		req.Model = c.config.DefaultModel
	}

	// 构建内容
	var contents []*genai.Content
	for _, msg := range req.Messages {
		role := genai.RoleUser
		if msg.Role == "assistant" || msg.Role == "model" {
			role = genai.RoleModel
		}
		
		content := &genai.Content{
			Role:  role,
			Parts: []*genai.Part{{Text: msg.Content}},
		}
		contents = append(contents, content)
	}

	// 构建生成配置
	config := &genai.GenerateContentConfig{}
	if req.Temperature > 0 {
		config.Temperature = &req.Temperature
	}
	if req.TopP > 0 {
		config.TopP = &req.TopP
	}
	if req.TopK > 0 {
		topK := float32(req.TopK)
		config.TopK = &topK
	}
	if req.MaxTokens > 0 {
		config.MaxOutputTokens = int32(req.MaxTokens)
	}

	// 生成流式内容
	iter := c.client.Models.GenerateContentStream(ctx, req.Model, contents, config)
	
	// 创建流式读取器
	return NewStreamReader(iter, req.Model), nil
}

// ListModels 列出可用模型
func (c *HTTPClient) ListModels(ctx context.Context) ([]string, error) {
	// Google AI 支持的模型列表
	models := []string{
		"gemini-1.5-flash",
		"gemini-1.5-pro",
		"gemini-2.0-flash-exp",
	}
	return models, nil
}

// ValidateAPIKey 验证API密钥
func (c *HTTPClient) ValidateAPIKey(ctx context.Context) error {
	// 确保客户端已初始化
	if err := c.ensureClient(ctx); err != nil {
		return err
	}

	// 尝试列出模型来验证API密钥
	_, err := c.ListModels(ctx)
	return err
}

// Close 关闭客户端
func (c *HTTPClient) Close() error {
	// Google AI SDK 的客户端不需要显式关闭
	c.client = nil
	return nil
}