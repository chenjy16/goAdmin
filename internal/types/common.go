package types

// CommonMessage 通用聊天消息结构
type CommonMessage struct {
	Role    string `json:"role"`    // system, user, assistant
	Content string `json:"content"`
}

// CommonUsage 通用使用统计结构
type CommonUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// CommonChoice 通用响应选择结构
type CommonChoice struct {
	Index        int           `json:"index"`
	Message      CommonMessage `json:"message"`
	FinishReason string        `json:"finish_reason"`
}

// CommonChatRequest 通用聊天请求结构
type CommonChatRequest struct {
	Model       string                 `json:"model"`
	Messages    []CommonMessage        `json:"messages"`
	MaxTokens   *int                   `json:"max_tokens,omitempty"`
	Temperature *float32               `json:"temperature,omitempty"`
	TopP        *float32               `json:"top_p,omitempty"`
	TopK        *int                   `json:"top_k,omitempty"`
	Stream      bool                   `json:"stream,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
}

// CommonChatResponse 通用聊天响应结构
type CommonChatResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []CommonChoice `json:"choices"`
	Usage   CommonUsage    `json:"usage"`
}

// ProviderType 提供商类型
type ProviderType string

const (
	ProviderTypeOpenAI   ProviderType = "openai"
	ProviderTypeGoogleAI ProviderType = "googleai"
	ProviderTypeMock     ProviderType = "mock"
)

// CommonErrorResponse 通用错误响应
type CommonErrorResponse struct {
	Error CommonError `json:"error"`
}

// CommonError 通用错误信息
type CommonError struct {
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
	Code    string `json:"code,omitempty"`
}

// CommonProviderInfo 通用提供商信息
type CommonProviderInfo struct {
	Type        ProviderType `json:"type"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Healthy     bool         `json:"healthy"`
	ModelCount  int          `json:"model_count"`
}