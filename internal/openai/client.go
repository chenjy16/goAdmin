package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// HTTPClient OpenAI HTTP 客户端实现
type HTTPClient struct {
	config     *Config
	keyManager KeyManager
	httpClient *http.Client
}

// NewHTTPClient 创建新的 HTTP 客户端
func NewHTTPClient(config *Config, keyManager KeyManager) *HTTPClient {
	return &HTTPClient{
		config:     config,
		keyManager: keyManager,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// ChatCompletion 实现聊天完成
func (c *HTTPClient) ChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// 设置默认值
	if req.Model == "" {
		req.Model = c.config.DefaultModel
	}
	
	// 从密钥管理器获取API密钥
	apiKey, err := c.keyManager.GetAPIKey()
	if err != nil {
		return nil, fmt.Errorf("get API key: %w", err)
	}
	
	// 序列化请求
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	
	// 创建 HTTP 请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+"/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	
	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	
	// 发送请求
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()
	
	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	
	// 检查错误响应
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
		}
		return nil, fmt.Errorf("OpenAI API error: %s", errResp.Error.Message)
	}
	
	// 解析成功响应
	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	
	return &chatResp, nil
}

// ChatCompletionStream 实现流式聊天完成
func (c *HTTPClient) ChatCompletionStream(ctx context.Context, req *ChatRequest) (io.ReadCloser, error) {
	// 设置流式模式
	req.Stream = true
	
	// 设置默认值
	if req.Model == "" {
		req.Model = c.config.DefaultModel
	}
	
	// 从密钥管理器获取API密钥
	apiKey, err := c.keyManager.GetAPIKey()
	if err != nil {
		return nil, fmt.Errorf("get API key: %w", err)
	}
	
	// 序列化请求
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	
	// 创建 HTTP 请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+"/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	
	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Accept", "text/event-stream")
	
	// 发送请求
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	
	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
		}
		return nil, fmt.Errorf("OpenAI API error: %s", errResp.Error.Message)
	}
	
	return resp.Body, nil
}

// ListModels 列出可用模型
func (c *HTTPClient) ListModels(ctx context.Context) ([]string, error) {
	// 创建 HTTP 请求
	httpReq, err := http.NewRequestWithContext(ctx, "GET", c.config.BaseURL+"/models", nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	
	// 设置请求头
	httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	
	// 发送请求
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()
	
	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	
	// 检查错误响应
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
		}
		return nil, fmt.Errorf("OpenAI API error: %s", errResp.Error.Message)
	}
	
	// 解析响应
	var modelsResp struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &modelsResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	
	// 提取模型名称
	models := make([]string, len(modelsResp.Data))
	for i, model := range modelsResp.Data {
		models[i] = model.ID
	}
	
	return models, nil
}

// ValidateAPIKey 验证 API 密钥
func (c *HTTPClient) ValidateAPIKey(ctx context.Context) error {
	// 从密钥管理器获取API密钥
	apiKey, err := c.keyManager.GetAPIKey()
	if err != nil {
		return fmt.Errorf("get API key: %w", err)
	}
	
	// 创建一个简单的请求来验证密钥
	httpReq, err := http.NewRequestWithContext(ctx, "GET", c.config.BaseURL+"/models", nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	
	// 设置请求头
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	
	// 发送请求
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()
	
	// 检查响应状态
	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("invalid API key")
	}
	
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API validation failed: HTTP %d: %s", resp.StatusCode, string(respBody))
	}
	
	return nil
}

// StreamReader 流式响应读取器
type StreamReader struct {
	reader  *bufio.Scanner
	closer  io.Closer
}

// NewStreamReader 创建流式读取器
func NewStreamReader(rc io.ReadCloser) *StreamReader {
	return &StreamReader{
		reader: bufio.NewScanner(rc),
		closer: rc,
	}
}

// Read 读取下一个流式响应
func (sr *StreamReader) Read() (*StreamResponse, error) {
	for sr.reader.Scan() {
		line := sr.reader.Text()
		
		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}
		
		// 处理 data: 行
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			
			// 检查结束标记
			if data == "[DONE]" {
				return nil, io.EOF
			}
			
			// 解析 JSON
			var resp StreamResponse
			if err := json.Unmarshal([]byte(data), &resp); err != nil {
				continue // 跳过无法解析的行
			}
			
			return &resp, nil
		}
	}
	
	if err := sr.reader.Err(); err != nil {
		return nil, err
	}
	
	return nil, io.EOF
}

// Close 关闭流式读取器
func (sr *StreamReader) Close() error {
	return sr.closer.Close()
}