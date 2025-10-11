# go-springAi - 智能股票分析AI助手

这是一个专业的股票分析AI助手平台，采用现代化的全栈架构设计，专注于为投资者和金融分析师提供智能化的股票分析服务。项目集成了多种AI模型和金融数据源，通过MCP协议支持和先进的AI技术，为用户提供实时股票分析、投资建议和市场洞察。

## 🎯 项目愿景

打造一个功能强大、易于使用的股票分析AI助手，帮助投资者做出更明智的投资决策。通过整合多种AI模型和实时金融数据，提供准确的股票分析、比较和投资建议，让每个人都能享受到专业级的投资分析服务。

## ✨ 特性

### 📈 股票分析核心功能
- 🔍 **智能股票分析**: 基于AI的深度股票分析，提供技术指标、基本面分析和市场趋势预测
- 📊 **实时股票数据**: 集成Yahoo Finance API，获取实时股票报价、历史数据和市场信息
- ⚖️ **股票对比分析**: 支持多只股票的横向对比，分析相对表现和投资价值
- 💡 **智能投资建议**: AI驱动的投资建议系统，根据市场分析提供个性化投资策略
- 📈 **技术指标计算**: 内置多种技术指标计算，包括移动平均线、RSI、MACD等
- 🎯 **风险评估**: 智能风险评估模型，帮助投资者了解投资风险
- 📱 **市场概览**: 提供全面的市场概览和行业分析
- 🔔 **价格预警**: 支持股票价格监控和预警功能


### 🤖 AI 集成能力
- 🚀 **多 AI 提供商支持**: 集成 OpenAI 和 Google AI，专门优化用于股票分析和金融数据处理
- 🔄 **统一 AI API**: 提供统一的聊天完成、模型管理和配置接口，支持股票分析专用提示
- 🧠 **股票分析AI助手**: 内置专业的股票分析助手，支持金融工具调用和投资上下文管理
- 📈 **金融数据理解**: AI模型专门训练用于理解和分析金融数据、市场趋势和投资指标
- 🔑 **API 密钥管理**: 动态 API 密钥设置和验证功能，支持多个金融数据源
- 📊 **模型管理**: 支持多模型切换和配置管理，针对不同分析场景优化
- 💬 **智能投资对话**: 支持流式响应和多轮投资咨询对话，提供个性化建议

### 🔧 MCP 协议支持
- 🛠️ **完整 MCP 实现**: 完整实现 Model Context Protocol 规范，专门优化金融分析工具集成
- 🔧 **股票分析工具系统**: 内置专业的股票分析工具集，包括股票分析、对比、建议和数据获取工具
- 📡 **SSE 流式通信**: 支持 Server-Sent Events 实时股票数据推送和流式分析响应
- 📝 **分析执行日志**: 完整的股票分析执行历史记录和性能监控
- 🔄 **动态工具注册**: 支持运行时股票分析工具发现和注册
- 🎯 **工具管理界面**: 提供可视化的股票分析工具管理和执行界面




## 🛠️ 技术栈

### 前端技术栈
- **React 19.1.1** - 现代化前端框架，支持并发特性和最新React 19特性
- **TypeScript 5.9.3** - 类型安全的JavaScript超集
- **Vite 7.1.7** - 快速的前端构建工具和开发服务器
- **Ant Design 5.27.4** - 企业级UI设计语言和组件库
- **Redux Toolkit 2.9.0** - 现代化的Redux状态管理
- **React Router DOM 6.30.1** - 声明式路由管理
- **Axios 1.12.2** - Promise based HTTP客户端
- **ESLint 9.36.0** - 代码质量和风格检查工具

### 后端技术栈

#### 核心框架
- **Go 1.24.0** - 编程语言
- **Gin v1.11.0** - HTTP Web 框架
- **SQLite3 v1.14.32** - 轻量级数据库

#### AI 集成
- **Google AI SDK v1.28.0** - Google AI 服务集成
- **OpenAI API** - OpenAI 服务集成（通过 HTTP 客户端）

#### 主要依赖
- **SQLC** - 类型安全的 SQL 代码生成器
- **Google Wire v0.7.0** - 依赖注入框架
- **Zap v1.27.0** - 结构化日志库
- **Viper v1.17.0** - 配置管理
- **Validator v10.28.0** - 数据验证
- **JWT v5.3.0** - JWT 令牌处理
- **UUID v1.6.0** - UUID 生成
- **Crypto v0.42.0** - 密码加密
- **Gorilla WebSocket v1.5.3** - WebSocket 和 SSE 支持

#### 测试框架
- **Testify v1.11.1** - 测试断言和模拟框架
- **Go Mock** - 接口模拟生成



## 🔧 股票分析工具系统

### 核心分析工具

#### 📊 stock_analysis - 智能股票分析工具
- **功能**: 基于AI的深度股票分析，提供全面的投资洞察
- **输入参数**:
  - `symbol`: 股票代码 (如: AAPL, GOOGL, TSLA)
  - `period`: 分析周期 (1d, 5d, 1mo, 3mo, 6mo, 1y, 2y, 5y, 10y, ytd, max)
- **分析内容**:
  - 技术指标分析 (移动平均线、RSI、MACD等)
  - 基本面分析 (市盈率、市净率、收益增长等)
  - 市场趋势预测
  - 风险评估和投资建议

#### 📈 yahoo_finance - 实时金融数据工具
- **功能**: 获取实时股票数据和历史价格信息
- **数据源**: Yahoo Finance API
- **提供数据**:
  - 实时股票报价
  - 历史价格数据
  - 交易量信息
  - 市场指标
  - 公司基本信息

#### ⚖️ stock_compare - 股票对比分析工具
- **功能**: 多只股票的横向对比分析
- **输入参数**:
  - `symbols`: 股票代码列表 (支持2-10只股票对比)
  - `metrics`: 对比指标 (价格表现、波动率、收益率等)
- **对比维度**:
  - 价格表现对比
  - 风险收益比分析
  - 技术指标对比
  - 行业地位分析

#### 💡 stock_advice - 智能投资建议工具
- **功能**: 基于AI分析提供个性化投资建议
- **建议类型**:
  - 买入/卖出/持有建议
  - 目标价位预测
  - 风险等级评估
  - 投资时机分析
  - 组合配置建议



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
│   ├── mcp/              # MCP 股票分析工具系统
│   │   ├── client.go         # MCP 客户端实现
│   │   ├── tool.go           # 基础工具定义和股票分析工具
│   │   ├── stock_tools.go    # 股票分析核心工具集
│   │   │   ├── StockAnalysisTool    # 智能股票分析工具
│   │   │   ├── YahooFinanceTool     # Yahoo Finance数据工具
│   │   │   ├── StockCompareTool     # 股票对比分析工具
│   │   │   └── StockAdviceTool      # 智能投资建议工具
│   │   └── *_test.go        # 股票分析工具测试文件
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
│   │   ├── stock_service.go        # 股票分析核心服务
│   │   ├── ai_assistant_service.go # 股票分析AI助手服务
│   │   ├── mcp_service.go          # MCP 股票工具服务实现
│   │   ├── openai_service.go       # OpenAI 股票分析服务
│   │   ├── googleai_service.go     # Google AI 股票分析服务
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
├── frontend/           # 前端应用
│   ├── public/        # 静态资源
│   │   ├── vite.svg
│   │   └── index.html
│   ├── src/           # 源代码
│   │   ├── components/    # 可复用组件
│   │   │   ├── Layout/   # 布局组件
│   │   │   │   └── index.tsx
│   │   │   └── common/   # 通用组件
│   │   ├── pages/        # 页面组件
│   │   │   ├── StockDashboardPage.tsx   # 股票分析仪表板页面
│   │   │   ├── StockAnalysisPage.tsx    # 股票分析页面
│   │   │   ├── StockComparePage.tsx     # 股票对比页面
│   │   │   ├── MCPToolsPage.tsx         # 股票分析工具页面
│   │   │   ├── AssistantPage.tsx        # 股票分析AI助手页面
│   │   │   ├── ProvidersPage.tsx        # AI提供商管理页面
│   │   │   └── SettingsPage.tsx         # 设置页面
│   │   ├── store/        # Redux状态管理
│   │   │   ├── index.ts  # Store配置
│   │   │   └── slices/   # 状态切片
│   │   │       ├── stockSlice.ts     # 股票数据状态
│   │   │       ├── analysisSlice.ts  # 股票分析状态
│   │   │       ├── authSlice.ts      # 认证状态
│   │   │       ├── providersSlice.ts # AI提供商状态
│   │   │       ├── mcpSlice.ts       # 股票工具状态
│   │   │       └── assistantSlice.ts # 股票AI助手状态
│   │   ├── services/     # API服务
│   │   │   ├── api.ts        # API配置
│   │   │   ├── stock.ts      # 股票数据服务
│   │   │   ├── analysis.ts   # 股票分析服务
│   │   │   ├── auth.ts       # 认证服务
│   │   │   ├── providers.ts  # AI提供商服务
│   │   │   ├── mcp.ts        # 股票分析工具服务
│   │   │   └── assistant.ts  # 股票AI助手服务
│   │   ├── types/        # TypeScript类型定义
│   │   │   ├── api.ts        # API类型
│   │   │   ├── stock.ts      # 股票数据类型
│   │   │   ├── analysis.ts   # 股票分析类型
│   │   │   ├── auth.ts       # 认证类型
│   │   │   ├── providers.ts  # AI提供商类型
│   │   │   ├── mcp.ts        # 股票分析工具类型
│   │   │   └── assistant.ts  # 股票AI助手类型
│   │   ├── utils/        # 工具函数
│   │   │   ├── request.ts    # 请求工具
│   │   │   ├── storage.ts    # 存储工具
│   │   │   └── constants.ts  # 常量定义
│   │   ├── router/       # 路由配置
│   │   │   └── index.tsx
│   │   ├── App.tsx       # 根组件
│   │   ├── main.tsx      # 应用入口
│   │   └── vite-env.d.ts # Vite类型声明
│   ├── package.json      # 前端依赖配置
│   ├── package-lock.json # 依赖锁定文件
│   ├── tsconfig.json     # TypeScript配置
│   ├── tsconfig.app.json # 应用TypeScript配置
│   ├── tsconfig.node.json # Node TypeScript配置
│   ├── vite.config.ts    # Vite配置
│   ├── eslint.config.js  # ESLint配置
│   └── README.md         # 前端文档
├── config.yaml         # 后端配置文件
├── sqlc.yaml          # SQLC 配置
├── Makefile           # 构建脚本
├── go.mod             # Go 模块文件
├── go.sum             # Go 模块校验和
└── .gitignore         # Git 忽略文件
```

## 🚀 快速开始

### 环境要求

#### 后端环境
- Go 1.24.0 或更高版本
- SQLite3

#### 前端环境
- Node.js 18.0 或更高版本
- npm 9.0 或更高版本

### 安装步骤

1. **克隆项目**
   ```bash
   git clone <repository-url>
   cd go-springAi
   ```

2. **后端安装**
   
   **安装Go依赖**
   ```bash
   go mod download
   ```
   
   **安装开发工具**
   ```bash
   # 安装 SQLC（用于生成数据库代码）
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
   
   # 安装 Wire（用于依赖注入）
   go install github.com/google/wire/cmd/wire@latest
   
   # 安装 Air（可选，用于热重载开发）
   go install github.com/air-verse/air@latest
   ```
   
   **生成代码**
   ```bash
   # 生成数据库访问代码
   sqlc generate
   
   # 生成依赖注入代码
   cd internal/wire && wire && cd ../..
   ```

3. **前端安装**
   ```bash
   # 进入前端目录
   cd frontend
   
   # 安装前端依赖
   npm install
   
   # 返回项目根目录
   cd ..
   ```

4. **初始化数据库**
   ```bash
   # 创建数据目录
   mkdir -p data
   
   # 初始化数据库
   sqlite3 data/go-springAi.db < schemas/users/001_create_users_table.sql
   ```

5. **配置应用**
   
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
     dsn: "./data/go-springAi.db"
   
   jwt:
     secret: "your-secret-key-change-this-in-production"
     expire_time: 24  # hours
   
   openai:
     api_key: ""  # 设置你的 OpenAI API 密钥
     base_url: "https://api.openai.com/v1"
   
   googleai:
     api_key: ""  # 设置你的 Google AI API 密钥
   ```

6. **运行应用**
   
   **启动后端服务**
   ```bash
   # 开发模式（推荐使用 Air 进行热重载）
   air
   
   # 或者直接运行
   go run cmd/main.go
   ```
   
   **启动前端服务**（新开一个终端）
   ```bash
   # 进入前端目录
   cd frontend
   
   # 启动开发服务器
   npm run dev
   ```

服务启动后：
- 后端API服务：`http://localhost:8080`
- 前端Web应用：`http://localhost:5173`

### 验证安装

#### 后端验证
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





### 后端API使用示例

#### AI 聊天完成示例

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

#### 2. 调用股票分析工具
```bash
curl -X POST http://localhost:8080/api/v1/mcp/execute \
  -H "Content-Type: application/json" \
  -d '{
    "name": "stock_analysis",
    "arguments": {
      "symbol": "AAPL"
    }
  }' | jq
```

#### 3. 调用股票对比工具
```bash
curl -X POST http://localhost:8080/api/v1/mcp/execute \
  -H "Content-Type: application/json" \
  -d '{
    "name": "stock_compare",
    "arguments": {
      "symbols": ["AAPL", "GOOGL", "MSFT"],
      "metrics": ["price_performance", "volatility", "market_cap"]
    }
  }' | jq
```

#### 4. 调用投资建议工具
```bash
curl -X POST http://localhost:8080/api/v1/mcp/execute \
  -H "Content-Type: application/json" \
  -d '{
    "name": "stock_advice",
    "arguments": {
      "symbol": "TSLA",
      "risk_tolerance": "moderate",
      "investment_horizon": "long_term"
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

# 股票分析使用示例
client = MCPClient()

# 获取可用的股票分析工具
tools = client.get_tools()
print("可用股票分析工具:", json.dumps(tools, indent=2, ensure_ascii=False))

# 分析苹果公司股票
apple_analysis = client.execute_tool("stock_analysis", {
    "symbol": "AAPL", 
    "period": "1mo"
})
print("苹果股票分析:", json.dumps(apple_analysis, indent=2, ensure_ascii=False))

# 对比多只科技股
comparison = client.execute_tool("stock_compare", {
    "symbols": ["AAPL", "GOOGL", "MSFT", "TSLA"],
    "metrics": ["price_performance", "volatility", "market_cap"]
})
print("科技股对比分析:", json.dumps(comparison, indent=2, ensure_ascii=False))

# 获取特斯拉投资建议
advice = client.execute_tool("stock_advice", {
    "symbol": "TSLA",
    "risk_tolerance": "moderate",
    "investment_horizon": "long_term"
})
print("特斯拉投资建议:", json.dumps(advice, indent=2, ensure_ascii=False))

# 获取执行日志
logs = client.get_logs()
print("分析执行日志:", json.dumps(logs, indent=2, ensure_ascii=False))
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
      "name": "stock_analysis",
      "description": "Analyze stock performance and provide insights",
      "inputSchema": {
        "type": "object",
        "properties": {
          "symbol": {
            "type": "string",
            "description": "Stock symbol to analyze"
          }
        },
        "required": ["symbol"]
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
  "name": "stock_analysis",
  "arguments": {
    "symbol": "AAPL"
  }
}
```

响应：
```json
{
  "content": [
    {
      "type": "text",
      "text": "Stock Analysis for AAPL: Current price $150.00, showing positive trend..."
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
      "tool_name": "stock_analysis",
      "arguments": {"symbol": "AAPL"},
      "result": {"content": [{"type": "text", "text": "Stock Analysis for AAPL: Current price $150.00, showing positive trend..."}]},
      "timestamp": "2024-01-01T00:00:00Z",
      "duration_ms": 150
    }
  ]
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
- `StockAnalysisTool`: 股票分析工具
- `YahooFinanceTool`: Yahoo Finance 数据工具
- `StockCompareTool`: 股票对比工具


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
    ChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
    GetModels() ([]ModelInfo, error)
}

// AI 控制器
type AIController struct {
    providers map[string]AIProviderStrategy
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

### 前端开发指南

#### 开发环境设置
```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build

# 预览生产构建
npm run preview
```

#### 添加新页面
1. 在 `src/pages/` 目录下创建新的页面组件
2. 在 `src/router/index.tsx` 中添加路由配置
3. 在 `src/components/Layout/index.tsx` 中添加菜单项（如需要）

#### 状态管理
使用 Redux Toolkit 进行状态管理：
```typescript
// 创建新的 slice
import { createSlice, PayloadAction } from '@reduxjs/toolkit'

interface MyState {
  data: any[]
  loading: boolean
}

const mySlice = createSlice({
  name: 'my',
  initialState: { data: [], loading: false } as MyState,
  reducers: {
    setData: (state, action: PayloadAction<any[]>) => {
      state.data = action.payload
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload
    }
  }
})

export const { setData, setLoading } = mySlice.actions
export default mySlice.reducer
```

#### API 服务
在 `src/services/` 目录下创建 API 服务：
```typescript
import { api } from './api'

export const myService = {
  getData: () => api.get('/api/v1/my/data'),
  createData: (data: any) => api.post('/api/v1/my/data', data),
  updateData: (id: string, data: any) => api.put(`/api/v1/my/data/${id}`, data),
  deleteData: (id: string) => api.delete(`/api/v1/my/data/${id}`)
}
```

#### 类型定义
在 `src/types/` 目录下定义 TypeScript 类型：
```typescript
export interface MyData {
  id: string
  name: string
  createdAt: string
  updatedAt: string
}

export interface MyApiResponse {
  data: MyData[]
  total: number
  page: number
  pageSize: number
}
```

### 后端开发指南

#### 项目架构说明

本项目采用清洁架构（Clean Architecture）设计，专为股票分析AI助手优化：

- **Stock Controllers**: 处理股票分析 HTTP 请求和响应，包括股票数据获取和分析端点
- **Stock Services**: 股票分析业务逻辑层，包括智能分析算法和投资策略服务
- **Financial Repository**: 金融数据访问层，管理股票数据和用户投资组合
- **Stock Models/DTO**: 股票数据传输对象，包括股票报价、分析结果和投资建议结构
- **Financial Middleware**: 金融级中间件（数据验证、安全日志、错误处理、风险控制等）
- **Stock Analysis Tools**: 可扩展的股票分析工具系统，支持多种分析策略

### MCP 协议支持

项目完整实现了 Model Context Protocol 规范，专为股票分析优化：

- **股票工具注册和发现**: 动态股票分析工具注册系统
- **安全股票分析执行**: 安全的股票数据处理和分析执行环境
- **实时股票数据流**: SSE 流式通信推送实时股票价格和分析结果
- **投资决策日志**: 完整的股票分析和投资建议执行历史记录
- **金融错误处理**: 统一的股票分析 MCP 错误响应格式和风险提示

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

### 使用 Makefile（推荐）

```bash
# 查看所有可用命令
make help

# 运行所有测试
make test

# 运行单元测试
make test-unit

# 运行集成测试
make test-integration

# 生成测试覆盖率报告
make test-coverage

# 运行竞态检测测试
make test-race

# 生成 Mock 文件
make mock-gen

# 清理测试缓存和生成文件
make clean
```

### 直接使用 Go 命令

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

### 前端构建

```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install

# 构建生产版本
npm run build

# 构建产物在 dist/ 目录下
```

### 后端构建

```bash
# 构建二进制文件
go build -o bin/admin cmd/main.go

# 或使用 Makefile
make build

# 交叉编译（Linux）
GOOS=linux GOARCH=amd64 go build -o bin/admin-linux cmd/main.go
```

### 完整部署

#### 方式一：分离部署
```bash
# 1. 构建前端
cd frontend
npm install
npm run build

# 2. 构建后端
cd ..
go build -o bin/admin cmd/main.go

# 3. 部署前端到静态文件服务器（如 Nginx）
# 4. 运行后端服务
./bin/admin
```

#### 方式二：Docker 部署

创建 `Dockerfile`：
```dockerfile
# 多阶段构建
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

FROM golang:1.24-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o admin cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# 复制后端二进制文件和配置
COPY --from=backend-builder /app/admin .
COPY --from=backend-builder /app/config.yaml .

# 复制前端构建产物
COPY --from=frontend-builder /app/frontend/dist ./static

# 暴露端口
EXPOSE 8080 5173

CMD ["./admin"]
```

构建和运行：
```bash
docker build -t go-springai .
docker run -p 8080:8080 -p 5173:5173 go-springai
```

#### 方式三：使用 Docker Compose

创建 `docker-compose.yml`：
```yaml
version: '3.8'

services:
  backend:
    build:
      context: .
      dockerfile: Dockerfile.backend
    ports:
      - "8080:8080"
    environment:
      - GIN_MODE=release
    volumes:
      - ./config.yaml:/app/config.yaml
      - ./data:/app/data

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    ports:
      - "5173:80"
    depends_on:
      - backend

volumes:
  data:
```

运行：
```bash
docker-compose up -d
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

### 后端技术
- [Gin](https://github.com/gin-gonic/gin) - 高性能 HTTP Web 框架
- [SQLC](https://github.com/sqlc-dev/sqlc) - 类型安全的 SQL 代码生成器
- [Wire](https://github.com/google/wire) - 编译时依赖注入框架
- [Zap](https://github.com/uber-go/zap) - 高性能结构化日志库
- [Viper](https://github.com/spf13/viper) - 灵活的配置管理库

### 前端技术
- [React](https://github.com/facebook/react) - 用户界面构建库
- [Vite](https://github.com/vitejs/vite) - 快速的前端构建工具
- [Ant Design](https://github.com/ant-design/ant-design) - 企业级UI设计语言
- [Redux Toolkit](https://github.com/reduxjs/redux-toolkit) - 现代化Redux状态管理
- [React Router](https://github.com/remix-run/react-router) - 声明式路由
- [Axios](https://github.com/axios/axios) - Promise based HTTP客户端
- [TypeScript](https://github.com/microsoft/TypeScript) - 类型安全的JavaScript

### AI与协议
- [OpenAI API](https://openai.com/api/) - OpenAI 服务集成
- [Google AI](https://ai.google.dev/) - Google AI 服务集成
- [Model Context Protocol](https://modelcontextprotocol.io/) - MCP 协议规范

---

⭐ 如果这个项目对你有帮助，请给它一个星标！