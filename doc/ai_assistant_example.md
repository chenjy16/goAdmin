# 股票分析AI助手集成示例

本文档展示如何使用智能股票分析AI助手功能，该功能结合了OpenAI大模型和专业的股票分析MCP工具系统，为投资者提供全面的股票分析服务。

## 功能特性

- **智能股票分析**: AI助手可以自动识别股票代码并调用相应的分析工具
- **实时金融数据**: 集成Yahoo Finance API，获取实时股票报价和历史数据
- **多维度分析**: 支持技术分析、基本面分析、风险评估和投资建议
- **股票对比**: 支持多只股票的横向对比分析
- **进程内通信**: MCP客户端和服务端在同一进程中，通过直接函数调用通信
- **流式响应**: 支持实时的对话体验
- **工具链执行**: 支持多个股票分析工具的连续调用

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
  "message": "请帮我分析一下AAPL股票的技术指标和投资建议",
  "model": "gpt-3.5-turbo",
  "temperature": 0.7,
  "max_tokens": 1500
}
```

**响应示例:**
```json
{
  "code": 200,
  "message": "Chat completed successfully",
  "data": {
    "content": "我已经为您完成了AAPL股票的综合分析。根据当前数据，AAPL股价为$175.43，技术指标显示股票处于上升趋势，RSI为65.2表明股票略微超买但仍在健康范围内。基于分析结果，建议适量买入并设置止损位。",
    "finish_reason": "stop",
    "tool_calls": [
      {
        "id": "call_1",
        "type": "function",
        "function": {
          "name": "yahoo_finance",
          "arguments": "{\"action\": \"quote\", \"symbol\": \"AAPL\"}"
        },
        "result": {
          "content": [
            {
              "type": "text",
              "text": "📈 AAPL 股票报价\n当前价格: $175.43\n涨跌: +2.15 (+1.24%)\n成交量: 45,678,900"
            }
          ]
        }
      },
      {
        "id": "call_2", 
        "type": "function",
        "function": {
          "name": "stock_analysis",
          "arguments": "{\"symbol\": \"AAPL\", \"analysis_type\": \"comprehensive\", \"period\": \"3mo\"}"
        },
        "result": {
          "content": [
            {
              "type": "text",
              "text": "📊 AAPL 综合分析报告\n\n🔍 技术分析:\n• RSI: 65.2 (略微超买)\n• MACD: 正向信号\n• 移动平均线: 上升趋势\n\n💰 基本面分析:\n• 市盈率: 28.5\n• 市净率: 6.8\n• 收益增长: 稳定\n\n⚠️ 风险评估:\n• 波动率: 中等\n• 流动性: 优秀\n• 市场风险: 低"
            }
          ]
        }
      },
      {
        "id": "call_3",
        "type": "function", 
        "function": {
          "name": "stock_advice",
          "arguments": "{\"symbol\": \"AAPL\", \"investment_horizon\": \"medium_term\", \"risk_tolerance\": \"moderate\"}"
        },
        "result": {
          "content": [
            {
              "type": "text",
              "text": "💡 AAPL 投资建议\n\n📈 投资评级: 买入\n🎯 目标价位: $185-190\n⏰ 投资期限: 3-6个月\n🛡️ 止损位: $165\n💰 建议仓位: 适中配置"
            }
          ]
        }
      }
    ],
    "usage": {
      "prompt_tokens": 180,
      "completion_tokens": 120,
      "total_tokens": 300
    }
  }
}
```

## 使用场景示例

### 场景1: 单只股票分析

```bash
curl -X POST http://localhost:8080/api/v1/assistant/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "请分析TSLA股票的技术指标和投资风险",
    "model": "gpt-3.5-turbo"
  }'
```

### 场景2: 多只股票对比分析

```bash
curl -X POST http://localhost:8080/api/v1/assistant/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "请对比分析AAPL、MSFT、GOOGL这三只科技股的投资价值",
    "model": "gpt-4"
  }'
```

### 场景3: 投资建议咨询

```bash
curl -X POST http://localhost:8080/api/v1/assistant/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "我有10000美元想投资NVDA股票，请给我一个详细的投资建议，我的风险承受能力是中等",
    "model": "gpt-3.5-turbo"
  }'
```

### 场景4: 实时股价查询

```bash
curl -X POST http://localhost:8080/api/v1/assistant/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "请查询AMZN的当前股价和今日表现",
    "model": "gpt-3.5-turbo"
  }'
```

### 场景5: 历史数据分析

```bash
curl -X POST http://localhost:8080/api/v1/assistant/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "分析META股票过去6个月的价格走势和交易量变化",
    "model": "gpt-4"
  }'
```

## 可用工具列表

AI助手当前支持以下股票分析MCP工具:

### 1. **yahoo_finance** - Yahoo Finance数据获取
- **功能**: 从Yahoo Finance获取实时股票数据
- **支持操作**: 
  - `quote`: 获取股票报价
  - `history`: 获取历史价格数据
  - `info`: 获取公司基本信息
- **参数**: 
  - `symbol`: 股票代码 (必需)
  - `action`: 操作类型 (必需)
  - `period`: 时间周期 (可选)
  - `interval`: 数据间隔 (可选)

**API调用示例:**
```bash
curl -X POST http://localhost:8080/api/v1/mcp/execute \
  -H "Content-Type: application/json" \
  -d '{
    "name": "yahoo_finance",
    "arguments": {
      "action": "quote",
      "symbol": "AAPL"
    }
  }'
```

**响应示例:**
```json
{
  "code": 200,
  "message": "Tool executed successfully",
  "data": {
    "content": [
      {
        "type": "text",
        "text": "📈 AAPL (AAPL) 股票报价\n\n💰 当前价格: $245.27\n📊 前收盘价: $254.04\n📈 今日开盘: $256.38\n🔺 今日最高: $256.38\n🔻 今日最低: $244.57\n📊 成交量: 61.2M\n🏢 市场: NMS\n💱 货币: USD\n⏰ 更新时间: 2025-10-11 04:00:02\n📉 涨跌: $-8.77 (-3.45%)"
      }
    ]
  }
}
```

### 2. **stock_analysis** - 股票技术分析
- **功能**: 提供股票的技术指标分析、基本面分析和风险评估
- **分析类型**:
  - `technical`: 技术分析 (RSI, MACD, 移动平均线等)
  - `fundamental`: 基本面分析 (市盈率, 市净率等)
  - `risk`: 风险评估
  - `comprehensive`: 综合分析
- **参数**:
  - `symbol`: 股票代码 (必需)
  - `analysis_type`: 分析类型 (可选，默认comprehensive)
  - `period`: 分析周期 (可选，默认3mo)

### 3. **stock_compare** - 股票对比分析
- **功能**: 对比多只股票的表现和投资价值
- **对比类型**:
  - `performance`: 表现对比
  - `valuation`: 估值对比
  - `risk`: 风险对比
  - `comprehensive`: 综合对比
- **参数**:
  - `symbols`: 股票代码列表 (必需)
  - `compare_type`: 对比类型 (可选，默认comprehensive)
  - `period`: 对比周期 (可选，默认3mo)

### 4. **stock_advice** - 投资建议
- **功能**: 基于分析结果提供个性化投资建议和风险提示
- **建议类型**:
  - 买入/卖出/持有建议
  - 目标价位设定
  - 止损位建议
  - 仓位配置建议
- **参数**:
  - `symbol`: 股票代码 (必需)
  - `investment_horizon`: 投资期限 (可选)
  - `risk_tolerance`: 风险承受能力 (可选)
  - `investment_amount`: 投资金额 (可选)

## 技术实现

### 架构设计

本项目采用专为股票分析优化的模块化架构：

- **股票分析AI服务层**: 负责处理股票相关对话和智能工具调用决策
- **金融数据MCP服务层**: 管理股票分析工具的注册、发现和执行
- **股票工具实现层**: 专业的股票分析功能工具实现
- **Yahoo Finance集成层**: 实时金融数据获取和处理

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

- **进程内通信**: AI助手与股票分析工具在同一进程内通信，确保实时性
- **异步执行**: 支持多只股票并行分析和结果聚合
- **流式响应**: 实时返回股票分析结果，提升交易决策效率
- **数据缓存**: 智能缓存股票数据，减少API调用次数

1. **HTTP请求**: 客户端通过HTTP API与AI助手交互
2. **直接函数调用**: MCP客户端直接调用MCP服务的方法，无需网络通信
3. **工具执行**: 通过工具注册表动态调用相应的工具实现
4. **响应聚合**: 将工具执行结果和AI响应合并返回

### 股票分析优势

1. **实时性**: 进程内通信，毫秒级响应股价变化
2. **专业性**: 专门针对股票分析场景优化的工具链
3. **准确性**: 基于Yahoo Finance的可靠数据源
4. **智能化**: AI驱动的多维度股票分析和投资建议
5. **可扩展**: 易于添加新的金融指标和分析模型

- **高性能**: 进程内通信避免了网络开销
- **类型安全**: 直接函数调用提供编译时类型检查
- **易于调试**: 统一的错误处理和日志记录
- **可扩展**: 通过MCP工具注册表轻松添加新功能

## 错误处理

- 股票数据获取失败时的优雅降级
- 详细的金融数据错误信息和调试日志
- 自动重试机制（针对Yahoo Finance API）
- 股票代码验证和错误提示

AI助手会处理以下类型的错误：

1. **工具调用错误**: 工具执行失败时的错误处理
2. **OpenAI API错误**: API调用失败的重试和错误报告
3. **参数验证错误**: 无效参数的验证和提示
4. **系统错误**: 内部服务错误的处理和恢复

## 配置要求

```yaml
# config/config.yaml
ai:
  default_model: "gpt-3.5-turbo"
  max_tokens: 2000
  temperature: 0.7

mcp:
  tools_enabled: true
  max_concurrent_tools: 5
  tool_timeout: 30s

stock_analysis:
  data_source: "yahoo_finance"
  cache_duration: "5m"
  max_symbols_per_request: 10
  default_period: "3mo"
  risk_free_rate: 0.02
```

确保在配置文件中设置了OpenAI API密钥：

```yaml
openai:
  api_key: "your-openai-api-key"
  base_url: "https://api.openai.com/v1"
  default_model: "gpt-3.5-turbo"
```