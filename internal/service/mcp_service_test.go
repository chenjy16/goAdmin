package service

import (
	"context"
	"testing"

	"admin/internal/dto"
	"admin/internal/mcp"
	"admin/internal/testutil"

	"go.uber.org/zap"
)

// MockUserService 模拟用户服务
type MockUserService struct{}

func (m *MockUserService) GetUser(ctx context.Context, id int64) (*dto.UserResponse, error) {
	return &dto.UserResponse{
		ID:       id,
		Username: "testuser",
		Email:    "test@example.com",
	}, nil
}

func (m *MockUserService) ListUsers(ctx context.Context, page, limit int64) ([]*dto.UserResponse, error) {
	return []*dto.UserResponse{
		{
			ID:       1,
			Username: "user1",
			Email:    "user1@example.com",
		},
		{
			ID:       2,
			Username: "user2",
			Email:    "user2@example.com",
		},
	}, nil
}

// 创建一个简单的mock GoogleAI服务
func createMockGoogleAIService() *GoogleAIService {
	// 这里我们创建一个真实的GoogleAIService实例，但使用mock的依赖
	// 在实际项目中，应该使用依赖注入和接口来更好地支持测试
	return &GoogleAIService{}
}

func TestMCPServiceImpl_Initialize(t *testing.T) {
	testutil.SetupGinTest()
	logger := testutil.SetupTestLogger(t)

	userService := &MockUserService{}
	googleaiService := createMockGoogleAIService()
	mcpService := NewMCPService(userService, googleaiService, logger)

	tests := []struct {
		name    string
		request *dto.MCPInitializeRequest
		wantErr bool
	}{
		{
			name: "successful initialization",
			request: &dto.MCPInitializeRequest{
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
			wantErr: false,
		},
		{
			name: "initialization with empty protocol version",
			request: &dto.MCPInitializeRequest{
				ProtocolVersion: "",
				ClientInfo: dto.MCPClientInfo{
					Name:    "test-client",
					Version: "1.0.0",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := testutil.TestContext()
			resp, err := mcpService.Initialize(ctx, tt.request)

			if tt.wantErr {
				testutil.AssertError(t, err)
				return
			}

			testutil.AssertNoError(t, err)
			testutil.AssertNotEqual(t, nil, resp)
			testutil.AssertEqual(t, "2024-11-05", resp.ProtocolVersion)
			testutil.AssertEqual(t, "Admin MCP Server", resp.ServerInfo.Name)
		})
	}
}

func TestMCPServiceImpl_ListTools(t *testing.T) {
	testutil.SetupGinTest()
	logger := testutil.SetupTestLogger(t)

	userService := &MockUserService{}
	googleaiService := createMockGoogleAIService()
	mcpService := NewMCPService(userService, googleaiService, logger)

	ctx := testutil.TestContext()
	resp, err := mcpService.ListTools(ctx)

	testutil.AssertNoError(t, err)
	testutil.AssertNotEqual(t, nil, resp)
	testutil.AssertNotEqual(t, 0, len(resp.Tools))

	// 检查默认工具是否存在
	toolNames := make([]string, len(resp.Tools))
	for i, tool := range resp.Tools {
		toolNames[i] = tool.Name
	}

	testutil.AssertContains(t, toolNames, "echo")
	testutil.AssertContains(t, toolNames, "get_user_info")
	testutil.AssertContains(t, toolNames, "googleai_chat")
}

func TestMCPServiceImpl_ExecuteTool(t *testing.T) {
	testutil.SetupGinTest()
	logger := testutil.SetupTestLogger(t)

	userService := &MockUserService{}
	googleaiService := createMockGoogleAIService()
	mcpService := NewMCPService(userService, googleaiService, logger)

	tests := []struct {
		name    string
		request *dto.MCPExecuteRequest
		wantErr bool
	}{
		{
			name: "execute echo tool successfully",
			request: &dto.MCPExecuteRequest{
				Name: "echo",
				Arguments: map[string]interface{}{
					"message": "Hello, World!",
				},
			},
			wantErr: false,
		},
		{
			name: "execute get_user_info tool successfully",
			request: &dto.MCPExecuteRequest{
				Name: "get_user_info",
				Arguments: map[string]interface{}{
					"user_id": float64(1), // JSON numbers are float64
				},
			},
			wantErr: false,
		},
		{
			name: "execute non-existent tool",
			request: &dto.MCPExecuteRequest{
				Name: "non_existent_tool",
				Arguments: map[string]interface{}{
					"param": "value",
				},
			},
			wantErr: true,
		},
		{
			name: "execute tool with invalid arguments",
			request: &dto.MCPExecuteRequest{
				Name: "get_user_info",
				Arguments: map[string]interface{}{
					"invalid_param": "value",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := testutil.TestContext()
			resp, err := mcpService.ExecuteTool(ctx, tt.request)

			if tt.wantErr {
				testutil.AssertError(t, err)
				return
			}

			testutil.AssertNoError(t, err)
			testutil.AssertNotEqual(t, nil, resp)

			if tt.request.Name == "echo" {
				testutil.AssertEqual(t, "text", resp.Content[0].Type)
			}
		})
	}
}

func TestMCPServiceImpl_RegisterTool(t *testing.T) {
	testutil.SetupGinTest()
	logger := testutil.SetupTestLogger(t)

	userService := &MockUserService{}
	googleaiService := createMockGoogleAIService()
	mcpService := NewMCPService(userService, googleaiService, logger)

	// 创建一个测试工具
	testTool := mcp.NewEchoTool()

	err := mcpService.RegisterTool(testTool)
	testutil.AssertNoError(t, err)

	// 验证工具已注册
	ctx := testutil.TestContext()
	resp, err := mcpService.ListTools(ctx)
	testutil.AssertNoError(t, err)

	found := false
	for _, tool := range resp.Tools {
		if tool.Name == "echo" {
			found = true
			break
		}
	}
	testutil.AssertEqual(t, true, found)
}

func TestMCPServiceImpl_GetExecutionLog(t *testing.T) {
	testutil.SetupGinTest()
	logger := testutil.SetupTestLogger(t)

	userService := &MockUserService{}
	googleaiService := createMockGoogleAIService()
	mcpService := NewMCPService(userService, googleaiService, logger)

	// 首先执行一个工具来创建执行日志
	ctx := testutil.TestContext()
	executeReq := &dto.MCPExecuteRequest{
		Name: "echo",
		Arguments: map[string]interface{}{
			"message": "Test message",
		},
	}

	_, err := mcpService.ExecuteTool(ctx, executeReq)
	testutil.AssertNoError(t, err)

	// 从服务实现中获取执行ID（这需要访问内部状态）
	// 为了测试目的，我们可以通过列出执行日志来获取最新的执行ID
	logs, err := mcpService.ListExecutionLogs(ctx, nil, 1)
	testutil.AssertNoError(t, err)
	testutil.AssertNotEqual(t, 0, len(logs))

	executionID := logs[0].ID

	// 获取执行日志
	log, err := mcpService.GetExecutionLog(ctx, executionID)
	testutil.AssertNoError(t, err)
	testutil.AssertNotEqual(t, nil, log)
	testutil.AssertEqual(t, executionID, log.ID)
	testutil.AssertEqual(t, "echo", log.ToolName)
}

func TestMCPServiceImpl_ListExecutionLogs(t *testing.T) {
	testutil.SetupGinTest()
	logger := testutil.SetupTestLogger(t)

	userService := &MockUserService{}
	googleaiService := createMockGoogleAIService()
	mcpService := NewMCPService(userService, googleaiService, logger)

	// 执行几个工具来创建执行日志
	ctx := testutil.TestContext()
	for i := 0; i < 3; i++ {
		executeReq := &dto.MCPExecuteRequest{
			Name: "echo",
			Arguments: map[string]interface{}{
				"message": "Test message",
			},
		}
		_, err := mcpService.ExecuteTool(ctx, executeReq)
		testutil.AssertNoError(t, err)
	}

	// 列出执行日志
	logs, err := mcpService.ListExecutionLogs(ctx, nil, 10)
	testutil.AssertNoError(t, err)
	testutil.AssertEqual(t, true, len(logs) >= 3)
}

// BenchmarkMCPService_ExecuteTool 性能测试
func BenchmarkMCPService_ExecuteTool(b *testing.B) {
	testutil.SetupGinTest()
	logger := zap.NewNop() // 使用无操作日志器以减少开销

	userService := &MockUserService{}
	googleaiService := createMockGoogleAIService()
	mcpService := NewMCPService(userService, googleaiService, logger)

	ctx := testutil.TestContext()
	req := &dto.MCPExecuteRequest{
		Name: "echo",
		Arguments: map[string]interface{}{
			"message": "Benchmark test message",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := mcpService.ExecuteTool(ctx, req)
		if err != nil {
			b.Fatalf("ExecuteTool failed: %v", err)
		}
	}
}