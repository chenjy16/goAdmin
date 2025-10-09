# go-springAi

这是一个架构设计优秀、核心功能完整的AI应用框架，特别是在MCP协议支持和AI集成方面表现出色。主要需要补充用户认证、安全性和生产环境支持功能，就能成为一个完整的企业级AI应用。

## ✨ 特性

### 🤖 AI 集成能力
- 🚀 **多 AI 提供商支持**: 集成 OpenAI 和 Google AI，支持统一的 API 接口
- 🔄 **统一 AI API**: 提供统一的聊天完成、模型管理和配置接口
- 🧠 **AI 助手服务**: 内置智能助手，支持工具调用和上下文管理
- 🔑 **API 密钥管理**: 动态 API 密钥设置和验证功能
- 📊 **模型管理**: 支持多模型切换和配置管理

### 🔧 MCP 协议支持
- 🛠️ **完整 MCP 实现**: 完整实现 Model Context Protocol 规范
- 🔧 **AI 工具系统**: 内置 OpenAI 和 Google AI 工具，支持可扩展的工具注册和执行
- 📡 **SSE 流式通信**: 支持 Server-Sent Events 实时事件推送和流式响应
- 📝 **执行日志**: 完整的工具执行历史记录和性能监控
- 🔄 **动态工具注册**: 支持运行时工具发现和注册

### 🏗️ 架构设计
- 🚀 **高性能**: 基于 Gin 框架，提供高性能的 HTTP 服务
- 🏗️ **清洁架构**: 采用分层架构设计，代码结构清晰，易于维护
- ⚡ **依赖注入**: 使用 Google Wire 进行依赖注入管理
- 🔧 **配置管理**: 使用 Viper 进行灵活的配置管理
- 🛡️ **中间件支持**: 完整的 CORS、日志、错误处理和恢复中间件

### 🗄️ 数据持久化
- 🗄️ **数据库支持**: 支持 SQLite 数据库，使用 SQLC 生成类型安全的数据库操作代码
- 👤 **用户管理**: 完整的用户 CRUD 操作和认证系统
- 🔐 **JWT 认证**: 基于 JWT 的用户认证和授权
- 🔒 **密码安全**: 安全的密码哈希和验证机制

### 🛡️ 安全与监控
- ✅ **数据验证**: 集成强大的数据验证功能
- 🛡️ **错误处理**: 统一的错误处理和安全日志记录
- 📊 **结构化日志**: 使用 Zap 提供详细的结构化日志记录
- 🔍 **监控支持**: 完整的请求/响应日志和性能监控
- 🚨 **安全日志**: 记录潜在安全威胁和异常行为

### 🧪 开发与测试
- 🧪 **完整测试**: 包含单元测试和集成测试，确保代码质量
- 🔨 **构建工具**: 完整的 Makefile 支持多种开发任务
- 📚 **文档完善**: 详细的 API 文档和使用示例
- 🔄 **热重载**: 支持 Air 热重载开发

## 🛠️ 技术栈

### 核心框架
- **Go 1.24.0** - 编程语言
- **Gin v1.11.0** - HTTP Web 框架
- **SQLite3 v1.14.32** - 轻量级数据库

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
- **Makefile** - 构建和开发任务自动化

## 📁 项目结构

```
go-springAi/
├── cmd/                    # 应用程序入口
│   └── main.go
├── doc/                    # 项目文档
│   ├── ai_assistant_example.md      # AI助手使用示例
│   ├── mcp_sequence_diagram.svg     # MCP序列图
│   └── 项目功能组件关系流程图.svg    # 项目架构图
├── internal/               # 内部包（不对外暴露）
│   ├── config/            # 配置管理
│   │   └── config.go
│   ├── controllers/       # 控制器层
│   │   ├── base_controller.go          # 基础控制器
│   │   ├── ai_controller.go            # 统一AI控制器
│   │   ├── ai_assistant_controller.go  # AI助手控制器
│   │   ├── openai_controller.go        # OpenAI控制器
│   │   ├── googleai_controller.go      # Google AI控制器
│   │   ├── mcp_controller.go           # MCP协议控制器
│   │   └── *_test.go                  # 控制器测试文件
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
│   │   └── errors.go
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
│   │   ├── client.go         # MCP 客户端实现
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
│   │   └── user_repository_mock.go
│   ├── openai/           # OpenAI 集成
│   │   ├── client.go     # OpenAI 客户端
│   │   ├── config.go     # 配置管理
│   │   ├── key_manager.go    # API 密钥管理
│   │   ├── model_manager.go  # 模型管理
│   │   └── types.go      # 类型定义
│   ├── provider/         # AI 提供商抽象层
│   │   ├── manager.go            # 提供商管理器
│   │   ├── openai_provider.go    # OpenAI 提供商
│   │   ├── googleai_provider.go  # Google AI 提供商
│   │   └── types.go              # 提供商接口定义
│   ├── repository/       # 数据访问层
│   │   ├── manager.go
│   │   ├── user_interfaces.go
│   │   └── user_repository.go
│   ├── response/         # 响应格式化
│   │   └── response.go
│   ├── route/           # 路由配置
│   │   └── routes.go
│   ├── service/         # 业务逻辑层
│   │   ├── ai_assistant_service.go # AI助手服务
│   │   ├── mcp_service.go          # MCP 服务实现
│   │   ├── openai_service.go       # OpenAI 服务
│   │   ├── googleai_service.go     # Google AI 服务
│   │   └── user_service.go         # 用户服务
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
├── go.mod             # Go 模块文件
├── go.sum             # Go 模块校验和
└── .gitignore         # Git 忽略文件
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
curl -X POST http://localhost:8080/api/v1/ai/openai/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
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
curl -X POST http://localhost:8080/api/v1/ai/googleai/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
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
curl -X GET "http://localhost:8080/api/v1/ai/openai/models" | jq
curl -X GET "http://localhost:8080/api/v1/ai/googleai/models" | jq
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

#### 3. 调用 OpenAI 模型工具
```bash
curl -X POST http://localhost:8080/api/v1/mcp/execute \
  -H "Content-Type: application/json" \
  -d '{
    "name": "openai_models",
    "arguments": {}
  }' | jq
```

#### 4. 调用用户信息工具
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
POST /api/v1/ai/openai/chat/completions
POST /api/v1/ai/googleai/chat/completions
Content-Type: application/json

{
  "model": "gpt-3.5-turbo",  // 或 "gemini-pro"
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
GET /api/v1/openai/models
GET /api/v1/googleai/models
GET /api/v1/ai/openai/models
GET /api/v1/ai/googleai/models
```

#### 5. API 密钥管理
```http
POST /api/v1/openai/api-key
POST /api/v1/googleai/api-key
Content-Type: application/json

{
  "api_key": "your-api-key-here"
}
```

#### 6. 验证 API 密钥
```http
POST /api/v1/openai/validate
POST /api/v1/googleai/validate
```

#### 7. 模型配置管理
```http
GET /api/v1/openai/config/:model
GET /api/v1/googleai/config/:model
GET /api/v1/ai/openai/config/:model
GET /api/v1/ai/googleai/config/:model
```

#### 8. 模型启用/禁用
```http
PUT /api/v1/openai/models/:model/enable
PUT /api/v1/openai/models/:model/disable
PUT /api/v1/googleai/models/:model/enable
PUT /api/v1/googleai/models/:model/disable
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

## 🏗️ MCP 客户端架构分析

### 架构概览

本项目实现了一个高性能的 Model Context Protocol (MCP) 客户端系统，采用进程内通信设计，优化了性能和类型安全性。整个 MCP 系统分为以下几个核心层次：

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTP API 层                              │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │  MCP Controller │  │  AI Controller  │  │ OpenAI/Google│ │
│  │                 │  │                 │  │  Controllers │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                   业务逻辑层                                 │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │   MCP Service   │  │ AI Assistant    │  │   Provider   │ │
│  │                 │  │    Service      │  │   Services   │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                   MCP 核心层                                │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │ Internal MCP    │  │   Tool Registry │  │ MCP Client   │ │
│  │    Client       │  │                 │  │   Manager    │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                   适配器层                                   │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │ OpenAI Service  │  │ Google AI       │  │ User Service │ │
│  │    Adapter      │  │   Adapter       │  │   Adapter    │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### 核心接口设计

#### 1. MCP 服务接口 (`MCPServiceInterface`)
```go
type MCPServiceInterface interface {
    Initialize(request dto.MCPInitializeRequest) (*dto.MCPInitializeResponse, error)
    ListTools() (*dto.MCPListToolsResponse, error)
    ExecuteTool(request dto.MCPExecuteToolRequest) (*dto.MCPExecuteToolResponse, error)
    GetExecutionLogs() []dto.MCPExecutionLog
    AddSSEClient(clientID string, writer http.ResponseWriter)
    RemoveSSEClient(clientID string)
    BroadcastEvent(event dto.MCPEvent)
}
```

#### 2. 内部 MCP 客户端接口 (`InternalMCPClient`)
```go
type InternalMCPClient interface {
    Initialize() error
    ListTools() ([]dto.MCPTool, error)
    ExecuteTool(name string, arguments map[string]interface{}) (interface{}, error)
}
```

#### 3. 工具接口 (`Tool`)
```go
type Tool interface {
    GetDefinition() dto.MCPTool
    Execute(arguments map[string]interface{}) (interface{}, error)
    Validate(arguments map[string]interface{}) error
}
```

### MCP 客户端实现架构

#### 1. 核心实现类

**InternalMCPClientImpl**
- 实现 `InternalMCPClient` 接口
- 使用直接函数调用，避免 JSON-RPC 开销
- 提供类型安全的工具执行环境

**MCPServiceImpl**
- 实现 `MCPServiceInterface` 接口
- 管理工具注册表 (`ToolRegistry`)
- 集成外部服务适配器
- 支持 SSE 实时事件推送
- 维护执行日志记录

**MCPClientManager**
- 管理多个 MCP 客户端实例
- 提供客户端生命周期管理
- 支持客户端注册和注销

#### 2. 工具系统架构

**ToolRegistry**
- 中央化工具注册管理
- 支持动态工具注册和发现
- 提供工具验证和执行

**BaseTool**
- 提供工具的基础实现
- 统一的错误处理和日志记录
- 可扩展的工具基类

**具体工具实现**
- `EchoTool`: 基础回显工具
- `UserInfoTool`: 用户信息查询工具
- `OpenAI 工具集`: OpenAI API 集成工具
  - `OpenAIChatTool`: 聊天完成工具
  - `OpenAIModelsTool`: 模型列表工具
  - `OpenAIConfigTool`: 配置管理工具

### MCP 客户端调用关系

#### 主要调用路径

**1. HTTP API → MCP Controller → MCP Service**
```
HTTP Request
    ↓
MCPController.Initialize/ListTools/ExecuteTool
    ↓
MCPServiceImpl.Initialize/ListTools/ExecuteTool
    ↓
ToolRegistry.GetTool/ExecuteTool
    ↓
Tool.Execute (具体工具实现)
    ↓
External Service (OpenAI/Google AI/User Service)
```

**2. AI Assistant → MCP Client → MCP Service**
```
AIAssistantService.Chat
    ↓
InternalMCPClient.ListTools/ExecuteTool
    ↓
MCPServiceImpl.ListTools/ExecuteTool
    ↓
ToolRegistry.GetTool/ExecuteTool
    ↓
Tool.Execute (具体工具实现)
```

#### 详细调用流程

**工具执行流程**
1. **请求接收**: HTTP 请求到达 `MCPController.ExecuteTool`
2. **参数验证**: 验证工具名称和参数格式
3. **服务调用**: 调用 `MCPServiceImpl.ExecuteTool`
4. **工具查找**: 从 `ToolRegistry` 中查找对应工具
5. **参数验证**: 调用 `Tool.Validate` 验证参数
6. **工具执行**: 调用 `Tool.Execute` 执行具体逻辑
7. **结果处理**: 格式化执行结果
8. **日志记录**: 记录执行日志到 `executionLogs`
9. **事件广播**: 通过 SSE 广播执行事件
10. **响应返回**: 返回格式化的 HTTP 响应

**AI 助手集成流程**
1. **聊天请求**: `AIAssistantService.Chat` 接收聊天请求
2. **工具发现**: 调用 `mcpClient.ListTools()` 获取可用工具
3. **系统消息**: 构建包含工具信息的系统消息
4. **AI 调用**: 调用 OpenAI/Google AI API
5. **工具调用解析**: 解析 AI 响应中的工具调用
6. **工具执行**: 通过 `mcpClient.ExecuteTool` 执行工具
7. **结果整合**: 将工具执行结果整合到对话中
8. **最终响应**: 生成包含工具执行结果的最终响应

#### 数据流转换

**DTO 层数据结构**
```go
// HTTP 请求 → DTO
dto.MCPExecuteToolRequest {
    Name: "tool_name",
    Arguments: map[string]interface{}
}

// DTO → 内部调用
Tool.Execute(arguments map[string]interface{})

// 内部结果 → DTO
dto.MCPExecuteToolResponse {
    Content: []dto.MCPContent
}

// DTO → HTTP 响应
response.Success(result)
```

**服务间数据传递**
1. **HTTP 层**: JSON 格式的请求/响应
2. **DTO 层**: 结构化的数据传输对象
3. **服务层**: Go 原生类型和接口
4. **工具层**: `map[string]interface{}` 参数传递
5. **适配器层**: 特定服务的 API 调用格式

#### 依赖注入关系

**Wire 依赖图**
```
MCPController
    ↓ (依赖)
MCPServiceInterface (MCPServiceImpl)
    ↓ (依赖)
├── ToolRegistry
├── UserService
├── OpenAIService  
├── GoogleAIService
└── Logger

AIAssistantService
    ↓ (依赖)
InternalMCPClient (InternalMCPClientImpl)
    ↓ (依赖)
MCPServiceInterface (MCPServiceImpl)
```

**关键依赖关系**
- `MCPController` 依赖 `MCPServiceInterface`
- `MCPServiceImpl` 依赖各种外部服务适配器
- `AIAssistantService` 依赖 `InternalMCPClient`
- `InternalMCPClientImpl` 内部调用 `MCPServiceImpl`
- 所有组件都通过 Wire 进行依赖注入管理

### MCP 系统设计模式

#### 1. 适配器模式 (Adapter Pattern)
**目的**: 将不同的外部服务 API 适配为统一的内部接口

**实现**:
```go
// 统一的服务接口
type ServiceInterface interface {
    Chat(request ChatRequest) (ChatResponse, error)
    GetModels() ([]Model, error)
}

// OpenAI 适配器
type OpenAIServiceAdapter struct {
    client *openai.Client
}

// Google AI 适配器  
type GoogleAIServiceAdapter struct {
    client *googleai.Client
}
```

**优势**:
- 统一不同 AI 提供商的 API 接口
- 便于切换和扩展新的 AI 服务
- 降低业务逻辑与具体实现的耦合

#### 2. 注册表模式 (Registry Pattern)
**目的**: 集中管理和发现系统中的工具和客户端

**实现**:
```go
// 工具注册表
type ToolRegistry struct {
    tools map[string]Tool
    mutex sync.RWMutex
}

// 客户端管理器
type MCPClientManager struct {
    clients map[string]InternalMCPClient
    mutex   sync.RWMutex
}
```

**优势**:
- 动态工具注册和发现
- 线程安全的资源管理
- 支持运行时工具扩展

#### 3. 策略模式 (Strategy Pattern)
**目的**: 根据不同的 AI 提供商选择不同的处理策略

**实现**:
```go
// AI 提供商策略接口
type AIProviderStrategy interface {
    ProcessChat(request UnifiedChatRequest) (UnifiedChatResponse, error)
    GetModels() ([]UnifiedModel, error)
}

// 统一 AI 控制器
type UnifiedAIController struct {
    strategies map[string]AIProviderStrategy
}
```

**优势**:
- 支持多 AI 提供商的统一接口
- 便于添加新的 AI 提供商
- 运行时策略选择

#### 4. 观察者模式 (Observer Pattern)
**目的**: 实现 SSE 事件的实时推送和订阅

**实现**:
```go
// SSE 客户端管理
type SSEClientManager struct {
    clients map[string]http.ResponseWriter
    mutex   sync.RWMutex
}

// 事件广播
func (s *MCPServiceImpl) BroadcastEvent(event dto.MCPEvent) {
    for clientID, writer := range s.sseClients {
        // 推送事件到客户端
    }
}
```

**优势**:
- 实时事件推送
- 多客户端订阅支持
- 解耦事件生产者和消费者

#### 5. 依赖注入模式 (Dependency Injection)
**目的**: 管理组件间的依赖关系，提高可测试性

**实现**:
```go
// Wire 提供商定义
func ProvideMCPService(
    toolRegistry *ToolRegistry,
    userService service.UserServiceInterface,
    openaiService service.OpenAIServiceInterface,
    googleaiService service.GoogleAIServiceInterface,
    logger *zap.Logger,
) service.MCPServiceInterface {
    return service.NewMCPService(...)
}
```

**优势**:
- 自动依赖解析和注入
- 提高代码可测试性
- 降低组件间耦合

### MCP 系统核心特性

#### 1. 进程内通信优化
**特点**:
- 避免 JSON-RPC 序列化开销
- 直接函数调用，性能更高
- 类型安全的参数传递

**实现**:
```go
type InternalMCPClientImpl struct {
    mcpService service.MCPServiceInterface
}

func (c *InternalMCPClientImpl) ExecuteTool(name string, arguments map[string]interface{}) (interface{}, error) {
    // 直接调用服务方法，无需网络通信
    return c.mcpService.ExecuteTool(dto.MCPExecuteToolRequest{
        Name: name,
        Arguments: arguments,
    })
}
```

#### 2. 模块化工具系统
**特点**:
- 可插拔的工具架构
- 统一的工具接口
- 动态工具注册

**工具生命周期**:
1. **定义**: 实现 `Tool` 接口
2. **注册**: 添加到 `ToolRegistry`
3. **发现**: 通过 `ListTools` API 暴露
4. **执行**: 通过 `ExecuteTool` API 调用
5. **监控**: 记录执行日志和性能指标

#### 3. 类型安全保障
**特点**:
- 强类型接口定义
- 编译时类型检查
- 运行时参数验证

**实现层次**:
```go
// 接口层：强类型接口
type Tool interface {
    Execute(arguments map[string]interface{}) (interface{}, error)
}

// 实现层：具体类型实现
type OpenAIChatTool struct {
    openaiService service.OpenAIServiceInterface
}

// 验证层：参数类型验证
func (t *OpenAIChatTool) Validate(arguments map[string]interface{}) error {
    // 验证必需参数和类型
}
```

#### 4. 实时事件系统
**特点**:
- SSE 长连接支持
- 实时事件推送
- 多客户端订阅

**事件类型**:
- 工具执行开始/完成事件
- 系统状态变更事件
- 错误和异常事件
- 性能监控事件

#### 5. 完整的可观测性
**特点**:
- 结构化日志记录
- 执行性能监控
- 错误追踪和分析
- 实时状态监控

**监控维度**:
```go
type MCPExecutionLog struct {
    ID          string                 `json:"id"`
    ToolName    string                 `json:"tool_name"`
    Arguments   map[string]interface{} `json:"arguments"`
    Result      interface{}            `json:"result"`
    Error       string                 `json:"error,omitempty"`
    Timestamp   time.Time              `json:"timestamp"`
    DurationMs  int64                  `json:"duration_ms"`
}
```

#### 6. 安全性保障
**特点**:
- 参数验证和清理
- 错误信息脱敏
- 安全日志记录
- 访问控制支持

**安全措施**:
- 输入参数严格验证
- 敏感信息过滤
- 异常情况安全处理
- 潜在攻击检测和记录

### MCP 数据流和执行流程

#### 完整的工具执行流程图

```
┌─────────────────┐    HTTP Request     ┌─────────────────┐
│   HTTP Client   │ ──────────────────→ │  MCPController  │
└─────────────────┘                     └─────────────────┘
                                                 │
                                                 │ 1. 参数验证
                                                 ↓
                                        ┌─────────────────┐
                                        │   MCPService    │
                                        │   Interface     │
                                        └─────────────────┘
                                                 │
                                                 │ 2. 工具查找
                                                 ↓
                                        ┌─────────────────┐
                                        │  ToolRegistry   │
                                        └─────────────────┘
                                                 │
                                                 │ 3. 获取工具实例
                                                 ↓
                                        ┌─────────────────┐
                                        │   Tool.Execute  │
                                        └─────────────────┘
                                                 │
                                                 │ 4. 执行具体逻辑
                                                 ↓
                                        ┌─────────────────┐
                                        │ External Service│
                                        │ (OpenAI/Google) │
                                        └─────────────────┘
                                                 │
                                                 │ 5. 返回结果
                                                 ↓
                                        ┌─────────────────┐
                                        │  Result Format  │
                                        │   & Logging     │
                                        └─────────────────┘
                                                 │
                                                 │ 6. SSE 事件广播
                                                 ↓
                                        ┌─────────────────┐
                                        │  HTTP Response  │
                                        └─────────────────┘
```

#### AI 助手集成流程图

```
┌─────────────────┐   Chat Request    ┌─────────────────┐
│   AI Assistant  │ ─────────────────→│ AIAssistantSvc  │
│     Client      │                   └─────────────────┘
└─────────────────┘                            │
                                               │ 1. 初始化 MCP 客户端
                                               ↓
                                      ┌─────────────────┐
                                      │ InternalMCP     │
                                      │    Client       │
                                      └─────────────────┘
                                               │
                                               │ 2. 获取可用工具
                                               ↓
                                      ┌─────────────────┐
                                      │  ListTools()    │
                                      └─────────────────┘
                                               │
                                               │ 3. 构建系统消息
                                               ↓
                                      ┌─────────────────┐
                                      │   AI Provider   │
                                      │ (OpenAI/Google) │
                                      └─────────────────┘
                                               │
                                               │ 4. AI 响应解析
                                               ↓
                                      ┌─────────────────┐
                                      │  Tool Calls     │
                                      │   Detection     │
                                      └─────────────────┘
                                               │
                                               │ 5. 执行工具调用
                                               ↓
                                      ┌─────────────────┐
                                      │ ExecuteTool()   │
                                      └─────────────────┘
                                               │
                                               │ 6. 整合结果
                                               ↓
                                      ┌─────────────────┐
                                      │ Final Response  │
                                      └─────────────────┘
```

#### 数据转换流程

```
HTTP JSON Request
        │
        │ JSON 解析
        ↓
┌─────────────────┐
│      DTO        │ ← MCPExecuteToolRequest
│   (Data Transfer│   {
│     Object)     │     "name": "tool_name",
└─────────────────┘     "arguments": {...}
        │               }
        │ 结构体转换
        ↓
┌─────────────────┐
│   Service       │ ← map[string]interface{}
│   Layer         │   参数映射
└─────────────────┘
        │
        │ 接口调用
        ↓
┌─────────────────┐
│   Tool          │ ← Tool.Execute(arguments)
│   Layer         │   具体工具实现
└─────────────────┘
        │
        │ 外部 API 调用
        ↓
┌─────────────────┐
│  External       │ ← OpenAI/Google AI API
│  Service        │   HTTP 请求/响应
└─────────────────┘
        │
        │ 结果处理
        ↓
┌─────────────────┐
│   Response      │ ← MCPExecuteToolResponse
│     DTO         │   {
└─────────────────┘     "content": [...]
        │               }
        │ JSON 序列化
        ↓
HTTP JSON Response
```

#### SSE 事件流程图

```
┌─────────────────┐    SSE Connect     ┌─────────────────┐
│   Web Client    │ ──────────────────→│  MCPController  │
└─────────────────┘                    └─────────────────┘
        │                                       │
        │                                       │ 注册 SSE 客户端
        │                                       ↓
        │                              ┌─────────────────┐
        │                              │   MCPService    │
        │                              │  SSE Manager    │
        │                              └─────────────────┘
        │                                       │
        │ ←─────── 实时事件推送 ─────────────────┘
        │
        │ 事件类型:
        │ • tool_execution_start
        │ • tool_execution_complete  
        │ • tool_execution_error
        │ • system_status_change
        │
┌─────────────────┐
│   Event Data    │
│   {             │
│     "type": "tool_execution_complete",
│     "data": {   │
│       "tool": "openai_chat",
│       "duration": 1500,
│       "status": "success"
│     }           │
│   }             │
└─────────────────┘
```

#### 错误处理流程

```
┌─────────────────┐
│   Tool Execute  │
└─────────────────┘
        │
        │ 执行过程中发生错误
        ↓
┌─────────────────┐
│  Error Capture  │ ← 捕获异常
└─────────────────┘
        │
        │ 错误分类
        ↓
┌─────────────────┐
│ Error Analysis  │ ← 参数错误/网络错误/业务错误
└─────────────────┘
        │
        │ 错误处理
        ↓
┌─────────────────┐
│ Error Response  │ ← 格式化错误响应
│     Format      │   {
└─────────────────┘     "error": "TOOL_EXECUTION_FAILED",
        │               "message": "Tool execution failed",
        │               "details": {...}
        │             }
        │ 安全日志记录
        ↓
┌─────────────────┐
│  Security Log   │ ← 记录潜在安全威胁
└─────────────────┘
        │
        │ SSE 错误事件
        ↓
┌─────────────────┐
│ Error Broadcast │ ← 广播错误事件给订阅客户端
└─────────────────┘
```

#### 性能监控流程

```
┌─────────────────┐
│ Request Start   │ ← 记录开始时间
└─────────────────┘
        │
        │ 执行监控
        ↓
┌─────────────────┐
│ Execution       │ ← 监控执行过程
│   Monitoring    │   • CPU 使用率
└─────────────────┘   • 内存使用量
        │             • 网络延迟
        │
        │ 完成监控
        ↓
┌─────────────────┐
│ Performance     │ ← 计算性能指标
│   Metrics       │   • 执行时间
└─────────────────┘   • 成功率
        │             • 错误率
        │
        │ 日志记录
        ↓
┌─────────────────┐
│ Structured      │ ← 结构化性能日志
│    Logging      │   {
└─────────────────┘     "duration_ms": 1500,
        │               "tool_name": "openai_chat",
        │               "status": "success",
        │               "memory_usage": "45MB"
        │             }
        │
        │ 实时监控
        ↓
┌─────────────────┐
│ Real-time       │ ← SSE 性能事件推送
│  Monitoring     │
└─────────────────┘
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



## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。



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