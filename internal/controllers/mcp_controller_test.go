package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"admin/internal/dto"
	"admin/internal/mcp"
	"admin/internal/testutil"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// MockMCPService 模拟MCP服务
type MockMCPService struct{}

func (m *MockMCPService) Initialize(ctx context.Context, req *dto.MCPInitializeRequest) (*dto.MCPInitializeResponse, error) {
	return &dto.MCPInitializeResponse{
		ProtocolVersion: "2024-11-05",
		ServerInfo: dto.MCPServerInfo{
			Name:    "Admin MCP Server",
			Version: "1.0.0",
		},
		Capabilities: dto.MCPCapabilities{
			Tools: &dto.MCPToolsCapability{
				ListChanged: true,
			},
		},
	}, nil
}

func (m *MockMCPService) ListTools(ctx context.Context) (*dto.MCPToolsResponse, error) {
	return &dto.MCPToolsResponse{
		Tools: []dto.MCPTool{
			{
				Name:        "echo",
				Description: "Echo tool for testing",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"message": map[string]interface{}{
							"type":        "string",
							"description": "Message to echo",
						},
					},
					"required": []string{"message"},
				},
			},
			{
				Name:        "get_user_info",
				Description: "Get user information",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"user_id": map[string]interface{}{
							"type":        "number",
							"description": "User ID",
						},
					},
					"required": []string{"user_id"},
				},
			},
		},
	}, nil
}

func (m *MockMCPService) ExecuteTool(ctx context.Context, req *dto.MCPExecuteRequest) (*dto.MCPExecuteResponse, error) {
	if req.Name == "echo" {
		message, ok := req.Arguments["message"].(string)
		if !ok {
			return nil, errors.New("Invalid parameters")
		}
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: message,
				},
			},
		}, nil
	}
	
	if req.Name == "get_user_info" {
		userID, ok := req.Arguments["user_id"].(float64)
		if !ok {
			return nil, errors.New("Invalid parameters")
		}
		return &dto.MCPExecuteResponse{
			Content: []dto.MCPContent{
				{
					Type: "text",
					Text: "User info for ID: " + string(rune(int(userID))),
				},
			},
		}, nil
	}
	
	return nil, errors.New("Method not found")
}

func (m *MockMCPService) RegisterTool(tool mcp.Tool) error {
	return nil
}

func (m *MockMCPService) GetExecutionLog(ctx context.Context, executionID string) (*dto.MCPToolExecutionLog, error) {
	now := time.Now()
	return &dto.MCPToolExecutionLog{
		ID:        executionID,
		ToolName:  "echo",
		StartTime: now,
		EndTime:   &now,
	}, nil
}

func (m *MockMCPService) ListExecutionLogs(ctx context.Context, toolName *string, limit int) ([]*dto.MCPToolExecutionLog, error) {
	now := time.Now()
	return []*dto.MCPToolExecutionLog{
		{
			ID:        "test-execution-1",
			ToolName:  "echo",
			StartTime: now,
			EndTime:   &now,
		},
		{
			ID:        "test-execution-2",
			ToolName:  "get_user_info",
			StartTime: now,
			EndTime:   &now,
		},
	}, nil
}

func setupMCPController() *MCPController {
	mockService := &MockMCPService{}
	logger := zap.NewNop()
	return NewMCPController(mockService, logger)
}

func TestMCPController_Initialize(t *testing.T) {
	testutil.SetupGinTest()
	controller := setupMCPController()

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful initialization",
			requestBody: dto.MCPInitializeRequest{
				ProtocolVersion: "2024-11-05",
				Capabilities: dto.MCPCapabilities{
					Tools: &dto.MCPToolsCapability{
						ListChanged: true,
					},
				},
				ClientInfo: dto.MCPClientInfo{
					Name:    "test-client",
					Version: "1.0.0",
				},
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Data dto.MCPInitializeResponse `json:"data"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				testutil.AssertNoError(t, err)
				testutil.AssertEqual(t, "2024-11-05", response.Data.ProtocolVersion)
				testutil.AssertEqual(t, "Admin MCP Server", response.Data.ServerInfo.Name)
			},
		},
		{
			name:           "invalid request body",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				testutil.AssertEqual(t, true, strings.Contains(w.Body.String(), "error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error
			
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				testutil.AssertNoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/mcp/initialize", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router := gin.New()
			// 添加错误处理中间件
			router.Use(func(c *gin.Context) {
				c.Next()
				if len(c.Errors) > 0 {
					c.JSON(http.StatusBadRequest, gin.H{"error": c.Errors.Last().Error()})
					c.Abort()
				}
			})
			router.POST("/mcp/initialize", controller.Initialize)
			router.ServeHTTP(w, req)

			testutil.AssertEqual(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

func TestMCPController_ListTools(t *testing.T) {
	testutil.SetupGinTest()
	controller := setupMCPController()

	req := httptest.NewRequest(http.MethodGet, "/mcp/tools", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.GET("/mcp/tools", controller.ListTools)
	router.ServeHTTP(w, req)

	testutil.AssertEqual(t, http.StatusOK, w.Code)

	var response struct {
		Data dto.MCPToolsResponse `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	testutil.AssertNoError(t, err)
	testutil.AssertEqual(t, 2, len(response.Data.Tools))
	testutil.AssertEqual(t, "echo", response.Data.Tools[0].Name)
	testutil.AssertEqual(t, "get_user_info", response.Data.Tools[1].Name)
}

func TestMCPController_ExecuteTool(t *testing.T) {
	testutil.SetupGinTest()
	controller := setupMCPController()

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "execute echo tool successfully",
			requestBody: dto.MCPExecuteRequest{
				Name: "echo",
				Arguments: map[string]interface{}{
					"message": "Hello, World!",
				},
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Data dto.MCPExecuteResponse `json:"data"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				testutil.AssertNoError(t, err)
				testutil.AssertEqual(t, 1, len(response.Data.Content))
				testutil.AssertEqual(t, "text", response.Data.Content[0].Type)
				testutil.AssertEqual(t, "Hello, World!", response.Data.Content[0].Text)
			},
		},
		{
			name: "execute get_user_info tool successfully",
			requestBody: dto.MCPExecuteRequest{
				Name: "get_user_info",
				Arguments: map[string]interface{}{
					"user_id": float64(1),
				},
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Data dto.MCPExecuteResponse `json:"data"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				testutil.AssertNoError(t, err)
				testutil.AssertEqual(t, 1, len(response.Data.Content))
				testutil.AssertEqual(t, "text", response.Data.Content[0].Type)
			},
		},
		{
			name: "execute non-existent tool",
			requestBody: dto.MCPExecuteRequest{
				Name: "non_existent_tool",
				Arguments: map[string]interface{}{
					"param": "value",
				},
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				testutil.AssertEqual(t, true, strings.Contains(w.Body.String(), "Method not found"))
			},
		},
		{
			name:           "invalid request body",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				testutil.AssertEqual(t, true, strings.Contains(w.Body.String(), "error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error
			
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				testutil.AssertNoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/mcp/execute", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router := gin.New()
			// 添加错误处理中间件
			router.Use(func(c *gin.Context) {
				c.Next()
				if len(c.Errors) > 0 {
					c.JSON(http.StatusBadRequest, gin.H{"error": c.Errors.Last().Error()})
					c.Abort()
				}
			})
			router.POST("/mcp/execute", controller.ExecuteTool)
			router.ServeHTTP(w, req)

			testutil.AssertEqual(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

func TestMCPController_GetExecutionLog(t *testing.T) {
	testutil.SetupGinTest()
	controller := setupMCPController()

	req := httptest.NewRequest(http.MethodGet, "/mcp/logs/test-execution-1", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.GET("/mcp/logs/:id", controller.GetExecutionLog)
	router.ServeHTTP(w, req)

	testutil.AssertEqual(t, http.StatusOK, w.Code)

	var response struct {
		Data dto.MCPToolExecutionLog `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	testutil.AssertNoError(t, err)
	testutil.AssertEqual(t, "test-execution-1", response.Data.ID)
	testutil.AssertEqual(t, "echo", response.Data.ToolName)
}

func TestMCPController_ListExecutionLogs(t *testing.T) {
	testutil.SetupGinTest()
	controller := setupMCPController()

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "list all execution logs",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Data struct {
						Logs []*dto.MCPToolExecutionLog `json:"logs"`
					} `json:"data"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				testutil.AssertNoError(t, err)
				testutil.AssertEqual(t, 2, len(response.Data.Logs))
				testutil.AssertEqual(t, "test-execution-1", response.Data.Logs[0].ID)
				testutil.AssertEqual(t, "test-execution-2", response.Data.Logs[1].ID)
			},
		},
		{
			name:           "list execution logs with limit",
			queryParams:    "?limit=1",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Data struct {
						Logs []*dto.MCPToolExecutionLog `json:"logs"`
					} `json:"data"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				testutil.AssertNoError(t, err)
				testutil.AssertEqual(t, 2, len(response.Data.Logs)) // Mock返回固定数量
			},
		},
		{
			name:           "list execution logs with tool filter",
			queryParams:    "?tool_name=echo",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Data struct {
						Logs []*dto.MCPToolExecutionLog `json:"logs"`
					} `json:"data"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				testutil.AssertNoError(t, err)
				testutil.AssertEqual(t, 2, len(response.Data.Logs)) // Mock返回固定数量
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/mcp/logs"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			router := gin.New()
			router.GET("/mcp/logs", controller.ListExecutionLogs)
			router.ServeHTTP(w, req)

			testutil.AssertEqual(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

// BenchmarkMCPController_ExecuteTool 性能测试
func BenchmarkMCPController_ExecuteTool(b *testing.B) {
	testutil.SetupGinTest()
	controller := setupMCPController()

	requestBody := dto.MCPExecuteRequest{
		Name: "echo",
		Arguments: map[string]interface{}{
			"message": "Benchmark test message",
		},
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		b.Fatalf("Failed to marshal request: %v", err)
	}

	router := gin.New()
	router.POST("/mcp/execute", controller.ExecuteTool)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/mcp/execute", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}