# AI API Gateway

一个基于 Go 语言开发的高性能 AI API 网关，集成 OpenAI 和 Google AI，支持 MCP（Model Context Protocol）协议，采用清洁架构设计，提供统一的 AI 服务接口、完整的工具系统、实时通信和结构化日志功能。

## ✨ 特性

- 🚀 **高性能**: 基于 Gin 框架，提供高性能的 HTTP 服务
- 🏗️ **清洁架构**: 采用分层架构设计，代码结构清晰，易于维护
- 🤖 **多 AI 提供商支持**: 集成 OpenAI 和 Google AI，支持统一的 API 接口
- 🔄 **统一 AI API**: 提供统一的聊天完成、模型管理和配置接口
- 🔧 **MCP 协议支持**: 完整实现 Model Context Protocol 规范
- 🛠️ **AI 工具系统**: 内置 OpenAI 和 Google AI 工具，支持可扩展的工具注册和执行
- 📡 **SSE 流式通信**: 支持 Server-Sent Events 实时事件推送和流式响应
- 🔑 **API 密钥管理**: 动态 API 密钥设置和验证功能
- 📊 **结构化日志**: 使用 Zap 提供详细的结构化日志记录
- 🗄️ **数据库支持**: 支持 SQLite 数据库，使用 SQLC 生成类型安全的数据库操作代码
- ⚡ **依赖注入**: 使用 Google Wire 进行依赖注入管理
- ✅ **数据验证**: 集成强大的数据验证功能
- 🔧 **配置管理**: 使用 Viper 进行灵活的配置管理
- 🛡️ **错误处理**: 统一的错误处理和安全日志记录
- 🔍 **监控支持**: 完整的请求/响应日志和性能监控
- 🧪 **完整测试**: 包含单元测试和集成测试，确保代码质量

## 🛠️ 技术栈

### 核心框架
- **Go 1.24.0** - 编程语言
- **Gin v1.11.0** - HTTP Web 框架
- **SQLite3** - 轻量级数据库

### AI 集成
- **Google AI SDK v1.28.0** - Google AI 服务集成
- **OpenAI API** - OpenAI 服务集成（通过 HTTP 客户端）

### 主要依赖
- **SQLC** - 类型安全的 SQL 代码生成器
- **Google Wire v0.7.0** - 依赖注入框架
- **Zap v1.27.0** - 结构化日志库
- **Viper v1.17.0** - 配置管理
- **Validator v10.28.0** - 数据验证
- **JWT v5.3.0** - JWT 令牌处理
- **UUID v1.6.0** - UUID 生成
- **Crypto v0.42.0** - 密码加密
- **Gorilla WebSocket v1.5.3** - WebSocket 和 SSE 支持

### 测试框架
- **Testify v1.11.1** - 测试断言和模拟框架
- **Go Mock** - 接口模拟生成

### 开发工具
- **Air** - 热重载开发工具（推荐）
- **Wire** - 依赖注入代码生成
- **SQLC** - SQL 代码生成

## 📁 项目结构

```
admin/
├── cmd/                    # 应用程序入口
│   └── main.go
├── internal/               # 内部包（不对外暴露）
│   ├── config/            # 配置管理
│   ├── controllers/       # 控制器层
│   │   ├── base_controller.go      # 基础控制器
│   │   ├── ai_controller.go        # 统一AI控制器
│   │   ├── openai_controller.go    # OpenAI控制器
│   │   ├── googleai_controller.go  # Google AI控制器
│   │   ├── mcp_controller.go       # MCP协议控制器
│   │   └── *_test.go              # 控制器测试文件
│   ├── database/          # 数据库相关
│   │   ├── connection.go  # 数据库连接
│   │   ├── curd/         # SQL 查询文件
│   │   └── generated/    # SQLC 生成的代码
│   ├── dto/              # 数据传输对象
│   │   ├── mcp.go        # MCP 协议相关 DTO
│   │   ├── openai.go     # OpenAI 相关 DTO
│   │   ├── googleai.go   # Google AI 相关 DTO
│   │   ├── unified.go    # 统一 AI DTO
│   │   └── user.go       # 用户相关 DTO
│   ├── errors/           # 错误处理
│   ├── googleai/         # Google AI 集成
│   │   ├── client.go     # Google AI 客户端
│   │   ├── config.go     # 配置管理
│   │   ├── key_manager.go    # API 密钥管理
│   │   ├── model_manager.go  # 模型管理
│   │   ├── stream.go     # 流式响应处理
│   │   └── types.go      # 类型定义
│   ├── logger/           # 日志系统
│   │   ├── constants.go
│   │   └── logger.go
│   ├── mcp/              # MCP 工具系统
│   │   ├── tool.go           # 基础工具定义
│   │   ├── openai_tool.go    # OpenAI 工具
│   │   ├── googleai_tool.go  # Google AI 工具
│   │   └── *_test.go        # 工具测试文件
│   ├── middleware/       # 中间件
│   │   ├── cors.go
│   │   ├── error_handler.go
│   │   ├── logger.go
│   │   ├── recovery.go
│   │   └── validation.go
│   ├── mocks/            # 测试模拟对象
│   │   ├── generate.go
│   │   └── *_mock.go
│   ├── openai/           # OpenAI 集成
│   │   ├── client.go     # OpenAI 客户端
│   │   ├── config.go     # 配置管理
│   │   ├── key_manager.go    # API 密钥管理
│   │   ├── model_manager.go  # 模型管理
│   │   └── types.go      # 类型定义
│   ├── provider/         # AI 提供商抽象层
│   │   ├── manager.go        # 提供商管理器
│   │   ├── openai_provider.go    # OpenAI 提供商
│   │   ├── googleai_provider.go  # Google AI 提供商
│   │   └── types.go      # 提供商接口定义
│   ├── repository/       # 数据访问层
│   │   ├── manager.go
│   │   ├── user_interfaces.go
│   │   ├── user_repository.go
│   │   └── *_test.go
│   ├── response/         # 响应格式化
│   ├── route/           # 路由配置
│   ├── service/         # 业务逻辑层
│   │   ├── mcp_service.go      # MCP 服务实现
│   │   ├── openai_service.go   # OpenAI 服务
│   │   ├── googleai_service.go # Google AI 服务
│   │   ├── user_service.go     # 用户服务
│   │   └── *_test.go          # 服务测试文件
│   ├── testutil/         # 测试工具
│   ├── utils/           # 工具函数
│   │   ├── jwt.go       # JWT 处理
│   │   ├── password.go  # 密码处理
│   │   └── validator.go # 验证器
│   └── wire/            # 依赖注入配置
│       ├── providers.go # 提供商定义
│       ├── wire.go      # Wire 配置
│       └── wire_gen.go  # Wire 生成代码
├── schemas/             # 数据库模式文件
│   └── users/
│       └── 001_create_users_table.sql
├── config.yaml         # 配置文件
├── sqlc.yaml          # SQLC 配置
├── Makefile           # 构建脚本
└── go.mod             # Go 模块文件
```

## 🚀 快速开始

### 环境要求

- Go 1.24.0 或更高版本
- SQLite3

### 安装步骤

1. **克隆项目**
   ```bash
   git clone <repository-url>
   cd admin
   ```

2. **安装依赖**
   ```bash
   go mod download
   ```

3. **安装开发工具**
   ```bash
   # 安装 SQLC（用于生成数据库代码）
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
   
   # 安装 Wire（用于依赖注入）
   go install github.com/google/wire/cmd/wire@latest
   
   # 安装 Air（可选，用于热重载开发）
   go install github.com/air-verse/air@latest
   ```

4. **生成代码**
   ```bash
   # 生成数据库访问代码
   sqlc generate
   
   # 生成依赖注入代码
   cd internal/wire && wire
   ```

5. **初始化数据库**
   ```bash
   # 创建数据目录
   mkdir -p data
   
   # 初始化数据库
   sqlite3 data/admin.db < schemas/users/001_create_users_table.sql
   ```

6. **配置应用**
   
   复制并修改配置文件：
   ```bash
   cp config.yaml config.local.yaml
   ```
   
   编辑 `config.local.yaml` 根据需要修改配置：
   ```yaml
   server:
     host: "localhost"
     port: "8080"
     mode: "debug"  # debug, release, test
   
   database:
     driver: "sqlite3"
     dsn: "./data/admin.db"
   
   jwt:
     secret: "your-secret-key-change-this-in-production"
     expire_time: 24  # hours
   
   openai:
     api_key: ""  # 设置你的 OpenAI API 密钥
     base_url: "https://api.openai.com/v1"
   
   googleai:
     api_key: ""  # 设置你的 Google AI API 密钥
   ```

7. **运行应用**
   ```bash
   # 开发模式（推荐使用 Air 进行热重载）
   air
   
   # 或者直接运行
   go run cmd/main.go
   ```

应用将在 `http://localhost:8080` 启动。

### 验证安装

访问健康检查端点：
```bash
curl http://localhost:8080/health
```

预期响应：
```json
{
  "status": "ok",
  "message": "Server is running"
}
```

## 💡 使用示例

### AI 聊天完成示例

#### 1. 使用统一 AI 接口
```bash
# OpenAI 聊天完成
curl -X POST http://localhost:8080/api/v1/ai/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "openai",
    "model": "gpt-3.5-turbo",
    "messages": [
      {
        "role": "user",
        "content": "你好，请介绍一下你自己"
      }
    ],
    "stream": false
  }' | jq

# Google AI 聊天完成
curl -X POST http://localhost:8080/api/v1/ai/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "googleai",
    "model": "gemini-pro",
    "messages": [
      {
        "role": "user",
        "content": "你好，请介绍一下你自己"
      }
    ],
    "stream": false
  }' | jq
```

#### 2. 设置 API 密钥
```bash
# 设置 OpenAI API 密钥
curl -X POST http://localhost:8080/api/v1/openai/api-key \
  -H "Content-Type: application/json" \
  -d '{
    "api_key": "sk-your-openai-api-key"
  }' | jq

# 设置 Google AI API 密钥
curl -X POST http://localhost:8080/api/v1/googleai/api-key \
  -H "Content-Type: application/json" \
  -d '{
    "api_key": "your-google-ai-api-key"
  }' | jq
```

#### 3. 获取可用模型
```bash
# 获取 OpenAI 模型列表
curl -X GET http://localhost:8080/api/v1/openai/models | jq

# 获取 Google AI 模型列表
curl -X GET http://localhost:8080/api/v1/googleai/models | jq

# 使用统一接口获取模型列表
curl -X GET "http://localhost:8080/api/v1/ai/models?provider=openai" | jq
```

### 基本 MCP 工具调用

#### 1. 获取可用工具列表
```bash
curl -X GET http://localhost:8080/api/v1/mcp/tools | jq
```

#### 2. 调用 Echo 工具
```bash
curl -X POST http://localhost:8080/api/v1/mcp/execute \
  -H "Content-Type: application/json" \
  -d '{
    "name": "echo",
    "arguments": {
      "message": "Hello MCP Server!"
    }
  }' | jq
```

#### 3. 调用用户信息工具
```bash
curl -X POST http://localhost:8080/api/v1/mcp/execute \
  -H "Content-Type: application/json" \
  -d '{
    "name": "get_user_info",
    "arguments": {
      "user_id": "12345"
    }
  }' | jq
```

### SSE 事件流连接

#### JavaScript 客户端示例
```javascript
// 建立 SSE 连接
const eventSource = new EventSource('http://localhost:8080/api/v1/mcp/sse');

// 监听消息事件
eventSource.onmessage = function(event) {
    const data = JSON.parse(event.data);
    console.log('收到事件:', data);
};

// 监听错误事件
eventSource.onerror = function(event) {
    console.error('SSE 连接错误:', event);
};

// 关闭连接
// eventSource.close();
```

#### curl 测试 SSE 连接
```bash
curl -N -H "Accept: text/event-stream" http://localhost:8080/api/v1/mcp/sse
```

### Python 客户端示例

```python
import requests
import json

class MCPClient:
    def __init__(self, base_url="http://localhost:8080"):
        self.base_url = base_url
        
    def get_tools(self):
        """获取可用工具列表"""
        response = requests.get(f"{self.base_url}/api/v1/mcp/tools")
        return response.json()
    
    def execute_tool(self, tool_name, arguments):
        """执行工具"""
        payload = {
            "name": tool_name,
            "arguments": arguments
        }
        response = requests.post(
            f"{self.base_url}/api/v1/mcp/execute",
            json=payload,
            headers={"Content-Type": "application/json"}
        )
        return response.json()
    
    def get_logs(self):
        """获取执行日志"""
        response = requests.get(f"{self.base_url}/api/v1/mcp/logs")
        return response.json()

# 使用示例
client = MCPClient()

# 获取工具列表
tools = client.get_tools()
print("可用工具:", json.dumps(tools, indent=2, ensure_ascii=False))

# 执行 echo 工具
result = client.execute_tool("echo", {"message": "Hello from Python!"})
print("执行结果:", json.dumps(result, indent=2, ensure_ascii=False))

# 获取执行日志
logs = client.get_logs()
print("执行日志:", json.dumps(logs, indent=2, ensure_ascii=False))
```

## 📚 API 文档

### 基础信息

- **Base URL**: `http://localhost:8080`
- **Content-Type**: `application/json`

### 健康检查 API

#### 服务器状态检查
```http
GET /health
```

响应：
```json
{
  "status": "ok",
  "message": "Server is running"
}
```

### MCP 协议 API

#### 1. MCP 初始化
```http
POST /api/v1/mcp/initialize
Content-Type: application/json

{
  "protocolVersion": "2024-11-05",
  "capabilities": {
    "tools": {}
  },
  "clientInfo": {
    "name": "example-client",
    "version": "1.0.0"
  }
}
```

#### 2. 获取可用工具列表
```http
GET /api/v1/mcp/tools
```

响应：
```json
{
  "tools": [
    {
      "name": "echo",
      "description": "Echo back the input message",
      "inputSchema": {
        "type": "object",
        "properties": {
          "message": {
            "type": "string",
            "description": "The message to echo back"
          }
        },
        "required": ["message"]
      }
    },
    {
      "name": "get_user_info",
      "description": "Get user information by user ID",
      "inputSchema": {
        "type": "object",
        "properties": {
          "user_id": {
            "type": "string",
            "description": "The user ID to get information for"
          }
        },
        "required": ["user_id"]
      }
    }
  ]
}
```

#### 3. 执行工具
```http
POST /api/v1/mcp/execute
Content-Type: application/json

{
  "name": "echo",
  "arguments": {
    "message": "Hello World"
  }
}
```

响应：
```json
{
  "content": [
    {
      "type": "text",
      "text": "Echo: Hello World"
    }
  ]
}
```

#### 4. SSE 事件流
```http
GET /api/v1/mcp/sse
```

建立 Server-Sent Events 连接，接收实时事件推送。

#### 5. 获取执行日志
```http
GET /api/v1/mcp/logs
```

响应：
```json
{
  "logs": [
    {
      "id": "log-id",
      "tool_name": "echo",
      "arguments": {"message": "Hello World"},
      "result": {"content": [{"type": "text", "text": "Echo: Hello World"}]},
      "timestamp": "2024-01-01T00:00:00Z",
      "duration_ms": 5
    }
  ]
}
```

### AI API

#### 1. 统一 AI 聊天完成
```http
POST /api/v1/ai/chat/completions
Content-Type: application/json

{
  "provider": "openai",  // 或 "googleai"
  "model": "gpt-3.5-turbo",
  "messages": [
    {
      "role": "user",
      "content": "Hello, how are you?"
    }
  ],
  "stream": false
}
```

#### 2. OpenAI 聊天完成
```http
POST /api/v1/openai/chat/completions
Content-Type: application/json

{
  "model": "gpt-3.5-turbo",
  "messages": [
    {
      "role": "user",
      "content": "Hello, how are you?"
    }
  ],
  "stream": false
}
```

#### 3. Google AI 聊天完成
```http
POST /api/v1/googleai/chat/completions
Content-Type: application/json

{
  "model": "gemini-pro",
  "messages": [
    {
      "role": "user",
      "content": "Hello, how are you?"
    }
  ],
  "stream": false
}
```

#### 4. 获取模型列表
```http
GET /api/v1/ai/models?provider=openai
GET /api/v1/openai/models
GET /api/v1/googleai/models
```

#### 5. API 密钥管理
```http
POST /api/v1/ai/api-key
POST /api/v1/openai/api-key
POST /api/v1/googleai/api-key
Content-Type: application/json

{
  "api_key": "your-api-key-here"
}
```

#### 6. 验证 API 密钥
```http
POST /api/v1/ai/validate?provider=openai
POST /api/v1/openai/validate
POST /api/v1/googleai/validate
```

### 响应格式

成功响应：
```json
{
  "code": 200,
  "message": "Success message",
  "data": {
    // 响应数据
  }
}
```

错误响应：
```json
{
  "code": 400,
  "message": "Error message",
  "error": "ERROR_CODE"
}
```

## 🔧 开发指南

### 代码生成

当修改数据库模式或查询时，需要重新生成代码：

```bash
# 重新生成数据库访问代码
sqlc generate

# 重新生成依赖注入代码
cd internal/wire && wire
```

### 添加新的 MCP 工具

1. 在 `internal/mcp/tool.go` 中定义新的工具结构体
2. 实现 `Tool` 接口的方法：
   - `GetDefinition()`: 返回工具定义
   - `Execute()`: 执行工具逻辑
   - `Validate()`: 验证输入参数
3. 在 `registerDefaultTools()` 函数中注册新工具
4. 重新生成依赖注入代码

示例工具实现：
```go
type MyTool struct {
    BaseTool
}

func (t *MyTool) GetDefinition() MCPTool {
    return MCPTool{
        Name:        "my_tool",
        Description: "My custom tool description",
        InputSchema: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "param": map[string]interface{}{
                    "type":        "string",
                    "description": "Parameter description",
                },
            },
            "required": []string{"param"},
        },
    }
}

func (t *MyTool) Execute(arguments map[string]interface{}) (interface{}, error) {
    // 实现工具逻辑
    return map[string]interface{}{
        "content": []map[string]interface{}{
            {
                "type": "text",
                "text": "Tool result",
            },
        },
    }, nil
}
```

### 项目架构说明

本项目采用清洁架构（Clean Architecture）设计：

- **Controllers**: 处理 HTTP 请求和响应，包括 MCP 协议端点
- **Services**: 业务逻辑层，包括 MCP 服务实现
- **Repository**: 数据访问层
- **Models/DTO**: 数据传输对象，包括 MCP 协议相关结构
- **Middleware**: 中间件（CORS、日志、错误处理、恢复等）
- **MCP Tools**: 可扩展的工具系统

### MCP 协议支持

项目完整实现了 Model Context Protocol 规范：

- **工具注册和发现**: 动态工具注册系统
- **工具执行**: 安全的工具执行环境
- **SSE 流式通信**: 实时事件推送
- **执行日志**: 完整的工具执行历史记录
- **错误处理**: 统一的 MCP 错误响应格式

### 错误处理

项目使用统一的错误处理机制：

- 自定义错误类型 `AppError`
- MCP 特定错误类型（工具未找到、执行失败等）
- 错误中间件自动处理和格式化错误响应
- 结构化错误日志记录
- 安全日志记录（记录潜在的安全威胁）

### 日志记录

使用 Zap 进行结构化日志记录：

- API 请求/响应日志
- MCP 工具执行日志
- 性能监控日志
- 安全事件日志
- 错误和异常日志
- 支持不同日志级别（DEBUG、INFO、WARN、ERROR）

## 🧪 测试

```bash
# 运行所有测试
go test ./...

# 运行测试并显示覆盖率
go test -cover ./...

# 生成测试覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📦 构建和部署

### 构建

```bash
# 构建二进制文件
go build -o bin/admin cmd/main.go

# 交叉编译（Linux）
GOOS=linux GOARCH=amd64 go build -o bin/admin-linux cmd/main.go
```

### Docker 部署

创建 `Dockerfile`：
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o admin cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/admin .
COPY --from=builder /app/config.yaml .
CMD ["./admin"]
```

构建和运行：
```bash
docker build -t admin-system .
docker run -p 8080:8080 admin-system
```

## 🤝 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

### 代码规范

- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 添加必要的注释和文档
- 编写单元测试

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 📞 联系方式

如有问题或建议，请通过以下方式联系：

- 提交 Issue
- 发送邮件至：[your-email@example.com]

## 🙏 致谢

感谢以下开源项目：

- [Gin](https://github.com/gin-gonic/gin) - 高性能 HTTP Web 框架
- [SQLC](https://github.com/sqlc-dev/sqlc) - 类型安全的 SQL 代码生成器
- [Wire](https://github.com/google/wire) - 编译时依赖注入框架
- [Zap](https://github.com/uber-go/zap) - 高性能结构化日志库
- [Viper](https://github.com/spf13/viper) - 灵活的配置管理库
- [Model Context Protocol](https://modelcontextprotocol.io/) - MCP 协议规范

---

⭐ 如果这个项目对你有帮助，请给它一个星标！