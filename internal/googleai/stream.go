package googleai

import (
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"strings"
	"time"

	"google.golang.org/genai"
)

// StreamReader 流式响应读取器
type StreamReader struct {
	iter  iter.Seq2[*genai.GenerateContentResponse, error]
	model string
	done  bool
	buf   []byte
}

// NewStreamReader 创建新的流式读取器
func NewStreamReader(iter iter.Seq2[*genai.GenerateContentResponse, error], model string) *StreamReader {
	return &StreamReader{
		iter:  iter,
		model: model,
		done:  false,
	}
}

// Read 实现 io.Reader 接口
func (sr *StreamReader) Read(p []byte) (n int, err error) {
	if sr.done {
		return 0, io.EOF
	}

	// 如果缓冲区为空，从迭代器读取下一个响应
	if len(sr.buf) == 0 {
		for resp, err := range sr.iter {
			if err != nil {
				sr.done = true
				return 0, err
			}

			// 转换为流式响应格式
			streamResp := sr.convertToStreamResponse(resp)
			
			// 序列化为JSON
			data, err := json.Marshal(streamResp)
			if err != nil {
				sr.done = true
				return 0, fmt.Errorf("marshal stream response: %w", err)
			}

			// 添加SSE格式
			sr.buf = append(sr.buf, []byte("data: ")...)
			sr.buf = append(sr.buf, data...)
			sr.buf = append(sr.buf, []byte("\n\n")...)
			break
		}

		// 如果没有更多数据，发送结束标记
		if len(sr.buf) == 0 {
			sr.buf = []byte("data: [DONE]\n\n")
			sr.done = true
		}
	}

	// 复制数据到输出缓冲区
	n = copy(p, sr.buf)
	sr.buf = sr.buf[n:]

	return n, nil
}

// Close 实现 io.Closer 接口
func (sr *StreamReader) Close() error {
	sr.done = true
	return nil
}

// convertToStreamResponse 将Google AI响应转换为流式响应格式
func (sr *StreamReader) convertToStreamResponse(resp *genai.GenerateContentResponse) *StreamResponse {
	streamResp := &StreamResponse{
		ID:      fmt.Sprintf("chatcmpl-%d", time.Now().Unix()),
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   sr.model,
		Choices: make([]StreamChoice, 0),
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

			choice := StreamChoice{
				Index: i,
				Delta: struct {
					Role    string `json:"role,omitempty"`
					Content string `json:"content,omitempty"`
				}{
					Role:    "assistant",
					Content: content.String(),
				},
			}

			if candidate.FinishReason != "" {
				finishReason := string(candidate.FinishReason)
				choice.FinishReason = &finishReason
			}

			streamResp.Choices = append(streamResp.Choices, choice)
		}
	}

	return streamResp
}