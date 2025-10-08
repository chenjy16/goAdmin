# MCP 服务器

一个基于 Go 语言开发的高性能 MCP（Model Context Protocol）服务器，采用清洁架构设计，提供完整的工具系统、实时通信和结构化日志功能。

## ✨ 特性

- 🚀 **高性能**: 基于 Gin 框架，提供高性能的 HTTP 服务
- 🏗️ **清洁架构**: 采用分层架构设计，代码结构清晰，易于维护
- 🔧 **MCP 协议支持**: 完整实现 Model Context Protocol 规范
- 🛠️ **工具系统**: 内置可扩展的工具注册和执行系统
- 📡 **SSE 流式通信**: 支持 Server-Sent Events 实时事件推送
- 📊 **结构化日志**: 使用 Zap 提供详细的结构化日志记录
- 🗄️ **数据库支持**: 支持 SQLite 数据库，使用 SQLC 生成类型安全的数据库操作代码
- ⚡ **依赖注入**: 使用 Google Wire 进行依赖注入管理
- ✅ **数据验证**: 集成强大的数据验证功能
- 🔧 **配置管理**: 使用 Viper 进行灵活的配置管理
- 🛡️ **错误处理**: 统一的错误处理和安全日志记录
- 🔍 **监控支持**: 完整的请求/响应日志和性能监控

## 🛠️ 技术栈

### 核心框架
- **Go 1.24.0** - 编程语言
- **Gin v1.11.0** - HTTP Web 框架
- **SQLite3** - 轻量级数据库

### 主要依赖
- **SQLC** - 类型安全的 SQL 代码生成器
- **Google Wire v0.7.0** - 依赖注入框架
- **Zap v1.27.0** - 结构化日志库
- **Viper v1.17.0** - 配置管理
- **Validator v10.28.0** - 数据验证
- **Crypto** - 密码加密
- **Gorilla WebSocket** - WebSocket 和 SSE 支持

### 开发工具
- **Air** - 热重载开发工具（推荐）
- **Wire** - 依赖注入代码生成

## 📁 项目结构

```
admin/
├── cmd/                    # 应用程序入口
│   └── main.go
├── internal/               # 内部包（不对外暴露）
│   ├── config/            # 配置管理
│   ├── controllers/       # 控制器层
│   │   ├── base_controller.go
│   │   └── mcp_controller.go  # MCP 协议控制器
│   ├── database/          # 数据库相关
│   │   ├── connection.go  # 数据库连接
│   │   ├── curd/         # SQL 查询文件
│   │   └── generated/    # SQLC 生成的代码
│   ├── dto/              # 数据传输对象
│   │   ├── mcp.go        # MCP 协议相关 DTO
│   │   └── user.go       # 用户相关 DTO
│   ├── errors/           # 错误处理
│   ├── logger/           # 日志系统
│   │   ├── constants.go
│   │   └── logger.go
│   ├── mcp/              # MCP 工具系统
│   │   └── tool.go       # 工具定义和注册
│   ├── middleware/       # 中间件
│   │   ├── cors.go
│   │   ├── error_handler.go
│   │   ├── logger.go
│   │   ├── recovery.go
│   │   └── validation.go
│   ├── repository/       # 数据访问层
│   ├── response/         # 响应格式化
│   ├── route/           # 路由配置
│   ├── service/         # 业务逻辑层
│   │   └── mcp_service.go  # MCP 服务实现
│   ├── utils/           # 工具函数
│   └── wire/            # 依赖注入配置
├── docs/               # 文档目录
├── schemas/             # 数据库模式文件
│   └── users/
├── config.yaml         # 配置文件
├── sqlc.yaml          # SQLC 配置
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