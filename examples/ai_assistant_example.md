# AI助手集成示例

本文档展示如何使用集成的AI助手功能，该功能结合了OpenAI大模型和MCP工具系统。

## 功能特性

- **智能工具调用**: AI助手可以自动识别并调用合适的MCP工具
- **进程内通信**: MCP客户端和服务端在同一进程中，通过直接函数调用通信
- **流式响应**: 支持实时的对话体验
- **工具链执行**: 支持多个工具的连续调用

## API端点

### 1. 初始化AI助手

```bash
POST /api/v1/assistant/initialize
```

**响应示例:**
```json
{
  "code": 200,
  "message": "AI assistant initialized successfully",
  "data": {
    "status": "initialized"
  }
}
```

### 2. 与AI助手对话

```bash
POST /api/v1/assistant/chat
```

**请求体:**
```json
{
  "message": "帮我查询当前用户信息，然后用echo工具重复一遍用户名",
  "model": "gpt-3.5-turbo",
  "temperature": 0.7,
  "max_tokens": 1000
}
```

**响应示例:**
```json
{
  "code": 200,
  "message": "Chat completed successfully",
  "data": {
    "content": "我已经为您查询了用户信息并使用echo工具重复了用户名。当前用户是admin，我已经通过echo工具确认了这个信息。",
    "finish_reason": "stop",
    "tool_calls": [
      {
        "id": "call_1",
        "type": "function",
        "function": {
          "name": "user_info",
          "arguments": "{}"
        },
        "result": {
          "user_id": "admin",
          "username": "admin",
          "role": "administrator"
        }
      },
      {
        "id": "call_2", 
        "type": "function",
        "function": {
          "name": "echo",
          "arguments": "{\"message\": \"admin\"}"
        },
        "result": {
          "message": "admin"
        }
      }
    ],
    "usage": {
      "prompt_tokens": 150,
      "completion_tokens": 80,
      "total_tokens": 230
    }
  }
}
```

## 使用场景示例

### 场景1: 信息查询和处理

```bash
curl -X POST http://localhost:8080/api/v1/assistant/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "获取用户信息并用echo工具确认用户名",
    "model": "gpt-3.5-turbo"
  }'
```

### 场景2: 多工具协作

```bash
curl -X POST http://localhost:8080/api/v1/assistant/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "先获取用户信息，然后使用Google AI聊天功能询问该用户的权限级别",
    "model": "gpt-4"
  }'
```

### 场景3: 配置管理

```bash
curl -X POST http://localhost:8080/api/v1/assistant/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "检查OpenAI模型配置，如果gpt-4模型被禁用，请启用它",
    "model": "gpt-3.5-turbo"
  }'
```

## 可用工具列表

AI助手可以调用以下MCP工具：

1. **echo** - 回显消息
2. **user_info** - 获取用户信息
3. **googleai_chat** - Google AI聊天
4. **googleai_models** - 获取Google AI模型列表
5. **googleai_config** - Google AI配置管理
6. **openai_chat** - OpenAI聊天
7. **openai_models** - 获取OpenAI模型列表
8. **openai_config** - OpenAI配置管理

## 技术实现

### 架构设计

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   HTTP Client   │───▶│ AIAssistant      │───▶│   OpenAI API    │
│                 │    │ Controller       │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌──────────────────┐
                       │ AIAssistant      │
                       │ Service          │
                       └──────────────────┘
                                │
                    ┌───────────┴───────────┐
                    ▼                       ▼
            ┌──────────────────┐    ┌──────────────────┐
            │ InternalMCP      │    │ OpenAI           │
            │ Client           │    │ Service          │
            └──────────────────┘    └──────────────────┘
                    │
                    ▼
            ┌──────────────────┐
            │ MCP Service      │
            │ (Tools Registry) │
            └──────────────────┘
```

### 通信机制

1. **HTTP请求**: 客户端通过HTTP API与AI助手交互
2. **直接函数调用**: MCP客户端直接调用MCP服务的方法，无需网络通信
3. **工具执行**: 通过工具注册表动态调用相应的工具实现
4. **响应聚合**: 将工具执行结果和AI响应合并返回

### 优势

- **高性能**: 进程内通信避免了网络开销
- **类型安全**: 直接函数调用提供编译时类型检查
- **易于调试**: 统一的错误处理和日志记录
- **可扩展**: 通过MCP工具注册表轻松添加新功能

## 错误处理

AI助手会处理以下类型的错误：

1. **工具调用错误**: 工具执行失败时的错误处理
2. **OpenAI API错误**: API调用失败的重试和错误报告
3. **参数验证错误**: 无效参数的验证和提示
4. **系统错误**: 内部服务错误的处理和恢复

## 配置要求

确保在配置文件中设置了OpenAI API密钥：

```yaml
openai:
  api_key: "your-openai-api-key"
  base_url: "https://api.openai.com/v1"
  default_model: "gpt-3.5-turbo"
```