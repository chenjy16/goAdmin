package service

import (
	"testing"

	"go.uber.org/zap"
)

func TestParseToolCalls(t *testing.T) {
	// 创建一个测试用的AIAssistantService实例
	logger := zap.NewNop()
	service := &AIAssistantService{
		logger: logger,
	}

	tests := []struct {
		name     string
		content  string
		expected int // 期望解析出的工具调用数量
	}{
		{
			name:     "Empty content",
			content:  "",
			expected: 0,
		},
		{
			name:     "Direct JSON tool call",
			content:  `{"name": "stock_analysis", "arguments": {"symbol": "AAPL"}}`,
			expected: 1,
		},
		{
			name:     "Wrapped tool call",
			content:  `{"tool_call": {"name": "stock_analysis", "arguments": {"symbol": "AAPL"}}}`,
			expected: 1,
		},
		{
			name:     "JSON code block",
			content:  "```json\n{\"name\": \"stock_analysis\", \"arguments\": {\"symbol\": \"AAPL\"}}\n```",
			expected: 1,
		},
		{
			name:     "Multiple tool calls",
			content:  `[{"name": "stock_analysis", "arguments": {"symbol": "AAPL"}}, {"name": "stock_compare", "arguments": {"symbols": ["AAPL", "GOOGL"]}}]`,
			expected: 2,
		},
		{
			name:     "Tool call with no arguments",
			content:  `{"name": "stock_analysis"}`,
			expected: 1,
		},
		{
			name:     "Invalid JSON",
			content:  `{"name": "stock_analysis", "arguments": {invalid json}`,
			expected: 0,
		},
		{
			name:     "Mixed content with tool call",
			content:  `Here is the analysis: {"tool_call": {"name": "stock_analysis", "arguments": {"symbol": "AAPL"}}} and some more text`,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.parseToolCalls(tt.content)
			if len(result) != tt.expected {
				t.Errorf("parseToolCalls() = %d tool calls, expected %d", len(result), tt.expected)
			}

			// 验证解析出的工具调用是否有效
			for _, toolCall := range result {
				if toolCall.Name == "" {
					t.Errorf("Tool call has empty name")
				}
				if toolCall.Arguments == nil {
					t.Errorf("Tool call has nil arguments")
				}
			}
		})
	}
}

func TestShouldRetryError(t *testing.T) {
	logger := zap.NewNop()
	service := &AIAssistantService{
		logger: logger,
	}

	tests := []struct {
		name     string
		errMsg   string
		expected bool
	}{
		{
			name:     "Timeout error",
			errMsg:   "connection timeout",
			expected: true,
		},
		{
			name:     "Connection refused",
			errMsg:   "connection refused",
			expected: true,
		},
		{
			name:     "Network unreachable",
			errMsg:   "network is unreachable",
			expected: true,
		},
		{
			name:     "Invalid argument error",
			errMsg:   "invalid argument provided",
			expected: false,
		},
		{
			name:     "Authentication error",
			errMsg:   "authentication failed",
			expected: false,
		},
		{
			name:     "No error",
			errMsg:   "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.errMsg != "" {
				err = &testError{msg: tt.errMsg}
			}
			
			result := service.shouldRetryError(err)
			if result != tt.expected {
				t.Errorf("shouldRetryError() = %v, expected %v for error: %s", result, tt.expected, tt.errMsg)
			}
		})
	}
}

func TestTruncateString(t *testing.T) {
	logger := zap.NewNop()
	service := &AIAssistantService{
		logger: logger,
	}

	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "Short string",
			input:    "hello",
			maxLen:   10,
			expected: "hello",
		},
		{
			name:     "Long string",
			input:    "this is a very long string that should be truncated",
			maxLen:   10,
			expected: "this is a ...",
		},
		{
			name:     "Exact length",
			input:    "exactly10c",
			maxLen:   10,
			expected: "exactly10c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.truncateString(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("truncateString() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// testError 是一个简单的错误实现用于测试
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}