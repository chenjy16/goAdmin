# go-springAi - Intelligent Stock Analysis AI Assistant

This is a professional stock analysis AI assistant platform with modern full-stack architecture design, focused on providing intelligent stock analysis services for investors and financial analysts. The project integrates multiple AI models and financial data sources, offering real-time stock analysis, investment advice, and market insights through MCP protocol support and advanced AI technology.


## ✨ Features

### 📈 Stock Analysis Core Features
- 🔍 **Intelligent Stock Analysis**: AI-powered deep stock analysis providing technical indicators, fundamental analysis, and market trend predictions
- 📊 **Real-time Stock Data**: Integrated with Yahoo Finance API for real-time stock quotes, historical data, and market information
- ⚖️ **Stock Comparison Analysis**: Support for horizontal comparison of multiple stocks, analyzing relative performance and investment value
- 💡 **Smart Investment Advice**: AI-driven investment recommendation system providing personalized investment strategies based on market analysis
- 📈 **Technical Indicator Calculation**: Built-in multiple technical indicator calculations including moving averages, RSI, MACD, etc.
- 🎯 **Risk Assessment**: Intelligent risk assessment model helping investors understand investment risks
- 📱 **Market Overview**: Comprehensive market overview and industry analysis
- 🔔 **Price Alerts**: Support for stock price monitoring and alert functionality

### 🤖 AI Integration Capabilities
- 🚀 **Multi-AI Provider Support**: Integration with OpenAI and Google AI, specifically optimized for stock analysis and financial data processing
- 🔄 **Unified AI API**: Provides unified chat completion, model management, and configuration interfaces with stock analysis-specific prompts
- 🧠 **Stock Analysis AI Assistant**: Built-in professional stock analysis assistant supporting financial tool calls and investment context management
- 📈 **Financial Data Understanding**: AI models specifically trained to understand and analyze financial data, market trends, and investment indicators
- 🔑 **API Key Management**: Dynamic API key setup and validation functionality supporting multiple financial data sources
- 📊 **Model Management**: Support for multi-model switching and configuration management optimized for different analysis scenarios
- 💬 **Intelligent Investment Dialogue**: Support for streaming responses and multi-turn investment consultation dialogues providing personalized advice

### 🔧 MCP Protocol Support
- 🛠️ **Complete MCP Implementation**: Full implementation of Model Context Protocol specification, specifically optimized for financial analysis tool integration
- 🔧 **Stock Analysis Tool System**: Built-in professional stock analysis tool set including stock analysis, comparison, advice, and data retrieval tools
- 📡 **SSE Streaming Communication**: Support for Server-Sent Events real-time stock data push and streaming analysis responses
- 📝 **Analysis Execution Logs**: Complete stock analysis execution history and performance monitoring
- 🔄 **Dynamic Tool Registration**: Support for runtime stock analysis tool discovery and registration
- 🎯 **Tool Management Interface**: Provides visual stock analysis tool management and execution interface

## 🛠️ Technology Stack

### 🏗️ Architecture Design
- **Architecture Pattern**: Frontend-Backend Separation + MCP Protocol Integration
- **Frontend**: React 19 + TypeScript + Vite + Ant Design
- **Backend**: Go + Gin + SQLite + Wire DI
- **AI Integration**: OpenAI + Google AI + Unified API Interface
- **Communication Protocol**: RESTful API + Server-Sent Events + WebSocket
- **Data Storage**: SQLite3 (Development) + Support for PostgreSQL/MySQL Extension

### Frontend Technology Stack
- **React 19.1.1** - Modern frontend framework supporting concurrent features and latest React 19 capabilities
- **TypeScript 5.9.3** - Type-safe JavaScript superset
- **Vite 7.1.7** - Fast frontend build tool and development server
- **Ant Design 5.27.4** - Enterprise-class UI design language and component library
  - **@ant-design/v5-patch-for-react-19** - React 19 compatibility patch
- **Redux Toolkit 2.9.0** - Modern Redux state management
- **React Router DOM 6.30.1** - Declarative routing management
- **Axios 1.12.2** - Promise-based HTTP client
- **ESLint 9.36.0** - Code quality and style checking tool
- **UUID 13.0.0** - Unique identifier generation

### Backend Technology Stack

#### Core Framework
- **Go 1.24.0** - Programming language (Latest version)
- **Gin v1.11.0** - High-performance HTTP web framework
- **SQLite3 v1.14.32** - Lightweight embedded database

#### AI Integration
- **Google AI SDK v1.28.0** - Google Gemini AI service integration
- **OpenAI API** - GPT series model integration (via HTTP client)
- **Unified AI Interface** - Unified calling interface supporting multiple AI providers

#### Core Dependencies
- **Google Wire v0.7.0** - Compile-time dependency injection framework
- **Zap v1.27.0** - High-performance structured logging library
- **Viper v1.17.0** - Configuration file management
- **Validator v10.28.0** - Request data validation
- **JWT v5.3.0** - JSON Web Token authentication
- **UUID v1.6.0** - Unique identifier generation
- **Crypto v0.42.0** - Password encryption and security features
- **Gorilla WebSocket v1.5.3** - WebSocket and SSE real-time communication

#### Development Tools
- **SQLC** - Type-safe SQL code generator
- **Testify v1.11.1** - Testing assertion and mocking framework
- **Go Mock v0.5.0** - Interface mock generator

## 🔧 Stock Analysis Tool System

### Core Analysis Tools

#### 📊 stock_analysis - Intelligent Stock Analysis Tool
- **Function**: AI-powered deep stock analysis providing comprehensive investment insights
- **Input Parameters**:
  - `symbol`: Stock symbol (e.g., AAPL, GOOGL, TSLA)
  - `period`: Analysis period (1d, 5d, 1mo, 3mo, 6mo, 1y, 2y, 5y, 10y, ytd, max)
- **Analysis Content**:
  - Technical indicator analysis (moving averages, RSI, MACD, etc.)
  - Fundamental analysis (P/E ratio, P/B ratio, earnings growth, etc.)
  - Market trend prediction
  - Risk assessment and investment advice

#### 📈 yahoo_finance - Real-time Financial Data Tool
- **Function**: Retrieve real-time stock data and historical price information
- **Data Source**: Yahoo Finance API
- **Provided Data**:
  - Real-time stock quotes
  - Historical price data
  - Trading volume information
  - Market indicators
  - Company basic information

#### ⚖️ stock_compare - Stock Comparison Analysis Tool
- **Function**: Horizontal comparison analysis of multiple stocks
- **Input Parameters**:
  - `symbols`: List of stock symbols (supports 2-10 stock comparison)
  - `metrics`: Comparison metrics (price performance, volatility, returns, etc.)
- **Comparison Dimensions**:
  - Price performance comparison
  - Risk-return ratio analysis
  - Technical indicator comparison
  - Industry position analysis

#### 💡 stock_advice - Smart Investment Advice Tool
- **Function**: Provide personalized investment advice based on AI analysis
- **Advice Types**:
  - Buy/Sell/Hold recommendations
  - Target price predictions
  - Risk level assessment
  - Investment timing analysis
  - Portfolio allocation suggestions

## 📁 Project Structure

```
go-springAi/
├── cmd/                    # Application entry point
│   └── main.go            # Main program entry
├── config.yaml            # Main configuration file
├── doc/                   # Project documentation
│   ├── ai_assistant_example.md      # AI assistant usage examples
│   ├── mcp_sequence_diagram.svg     # MCP sequence diagram
│   └── 项目功能组件关系流程图.svg    # Project architecture diagram
├── frontend/              # Frontend project
│   ├── src/              # Frontend source code
│   │   ├── components/   # React components
│   │   │   ├── AssistantConfigPanel/  # AI assistant configuration panel
│   │   │   ├── Layout/               # Layout components
│   │   │   ├── ModelSelector/        # Model selector
│   │   │   ├── ParameterSettings/    # Parameter settings component
│   │   │   └── ToolSelector/         # Tool selector
│   │   ├── pages/        # Page components
│   │   │   ├── AssistantPage.tsx     # AI assistant page
│   │   │   ├── MCPToolsPage.tsx      # MCP tools page
│   │   │   ├── ProvidersPage.tsx     # Provider management page
│   │   │   └── SettingsPage.tsx      # Settings page
│   │   ├── services/     # Service layer
│   │   ├── store/        # Redux state management
│   │   ├── hooks/        # Custom hooks
│   │   └── types/        # TypeScript type definitions
│   ├── package.json      # Frontend dependency configuration
│   └── vite.config.ts    # Vite build configuration
├── internal/              # Backend internal packages (not exposed externally)
│   ├── config/           # Configuration management
│   │   └── config.go
│   ├── controllers/      # Controller layer
│   │   ├── base_controller.go          # Base controller
│   │   ├── ai_controller.go            # Unified AI controller
│   │   ├── ai_assistant_controller.go  # AI assistant controller
│   │   └── mcp_controller.go           # MCP protocol controller
│   ├── database/         # Database related
│   │   ├── connection.go  # Database connection
│   │   ├── curd/         # SQL query files
│   │   └── generated/    # SQLC generated code
│   ├── dto/              # Data Transfer Objects
│   │   ├── mcp.go        # MCP protocol related DTOs
│   │   ├── openai.go     # OpenAI related DTOs
│   │   ├── googleai.go   # Google AI related DTOs
│   │   ├── unified.go    # Unified AI DTOs
│   │   ├── stock_analysis.go  # Stock analysis DTOs
│   │   └── user.go       # User related DTOs
│   ├── errors/           # Error handling
│   │   └── errors.go
│   ├── googleai/         # Google AI integration
│   │   ├── client.go     # Google AI client
│   │   ├── config.go     # Configuration management
│   │   ├── key_manager.go    # API key management
│   │   ├── model_manager.go  # Model management
│   │   ├── stream.go     # Streaming response handling
│   │   └── types.go      # Type definitions

│   ├── logger/           # Logging system
│   │   ├── constants.go  # Logging constants
│   │   └── logger.go     # Logging implementation
│   ├── mcp/              # MCP protocol implementation
│   │   ├── client.go     # MCP client implementation
│   │   ├── tool.go       # Base tool definition
│   │   ├── tools/        # Stock analysis tool set
│   │   └── tool_test.go  # Tool tests
│   ├── middleware/       # HTTP middleware
│   │   ├── auth.go       # Authentication middleware
│   │   ├── cors.go       # CORS handling
│   │   ├── error_handler.go  # Error handling
│   │   ├── logger.go     # Logging middleware
│   │   ├── recovery.go   # Recovery middleware
│   │   └── validation.go # Validation middleware
│   ├── mocks/            # Test mock objects
│   │   ├── generate.go
│   │   └── user_repository_mock.go
│   ├── openai/           # OpenAI integration
│   │   ├── client.go     # OpenAI client
│   │   ├── config.go     # Configuration management
│   │   ├── key_manager.go    # API key management
│   │   ├── model_manager.go  # Model management
│   │   └── types.go      # Type definitions
│   ├── provider/         # AI provider abstraction layer
│   │   ├── manager.go            # Provider manager
│   │   ├── openai_provider.go    # OpenAI provider
│   │   ├── googleai_provider.go  # Google AI provider
│   │   ├── mock_provider.go      # Mock provider
│   │   └── types.go              # Provider interface definitions
│   ├── repository/       # Data access layer
│   │   ├── manager.go            # Repository manager
│   │   ├── user_interfaces.go    # User interface definitions
│   │   ├── user_repository.go    # User repository implementation
│   │   ├── api_key_interfaces.go # API key interface definitions
│   │   └── api_key_repository.go # API key repository implementation
│   ├── response/         # Response handling
│   │   └── response.go   # Unified response format
│   ├── route/            # Route configuration
│   │   └── routes.go     # Route definitions
│   ├── service/          # Business logic layer
│   │   ├── ai_assistant_service.go    # AI assistant service
│   │   ├── ai_assistant_service_test.go # AI assistant service tests
│   │   ├── api_key_service.go         # API key service
│   │   ├── database_key_manager.go    # Database key management
│   │   ├── googleai_service.go        # Google AI service
│   │   ├── mcp_service.go             # MCP service
│   │   ├── openai_service.go          # OpenAI service
│   │   ├── stock_analysis_service.go  # Stock analysis service
│   │   └── user_service.go            # User service
│   ├── utils/            # Utility functions
│   │   ├── jwt.go        # JWT utilities
│   │   ├── password.go   # Password utilities
│   │   └── validator.go  # Validation utilities
│   └── wire/             # Dependency injection
│       ├── providers.go  # Provider definitions
│       ├── wire.go       # Wire configuration
│       └── wire_gen.go   # Wire generated code
├── schemas/              # Database schemas
│   ├── api_keys/         # API key table structure
│   └── users/            # User table structure
├── data/                 # Data directory
│   └── go-springAi.db    # SQLite database file
├── bin/                  # Compiled output
│   ├── admin             # Admin program
│   ├── go-springAi       # Main program
│   └── test              # Test program
├── go.mod                # Go module definition
├── go.sum                # Go module checksum
├── sqlc.yaml             # SQLC configuration
├── Makefile              # Build scripts
└── README.md             # Project documentation
```

## 🔧 Core Features

### Stock Analysis Tools (MCP Tools)

This project integrates multiple professional stock analysis tools through the MCP (Model Context Protocol) providing a unified interface:

#### 1. Stock Basic Analysis Tool (stock_analysis)
- **Function**: Retrieve stock basic information, real-time prices, technical indicators
- **Parameters**: Stock symbol (symbol)
- **Analysis Content**: 
  - Basic information (company name, industry, market cap, etc.)
  - Price data (current price, opening price, high, low)
  - Technical indicators (moving averages, RSI, MACD, etc.)
  - Volume analysis

#### 2. Yahoo Finance Data Tool (yahoo_finance)
- **Function**: Retrieve detailed stock financial data from Yahoo Finance
- **Parameters**: Stock symbol (symbol), data type (data_type)
- **Data Types**:
  - `info`: Basic information
  - `history`: Historical price data
  - `financials`: Financial statements
  - `balance_sheet`: Balance sheet
  - `cashflow`: Cash flow statement

#### 3. Stock Comparison Analysis Tool (stock_compare)
- **Function**: Comparative analysis of multiple stock performances
- **Parameters**: Stock symbol list (symbols), comparison period (period)
- **Comparison Dimensions**:
  - Price performance comparison
  - Volatility analysis
  - Correlation analysis
  - Risk-return comparison

#### 4. Stock Investment Advice Tool (stock_advice)
- **Function**: Provide investment advice based on technical and fundamental analysis
- **Parameters**: Stock symbol (symbol), analysis type (analysis_type)
- **Advice Types**:
  - `technical`: Technical analysis advice
  - `fundamental`: Fundamental analysis advice
  - `comprehensive`: Comprehensive analysis advice

### AI Assistant Features

This project provides intelligent stock analysis AI assistant supporting natural language interaction:

- **Multi-model Support**: Integration with OpenAI GPT and Google Gemini models
- **Real-time Analysis**: Combines real-time stock data for intelligent analysis
- **Personalized Recommendations**: Provides customized investment advice based on user preferences
- **Risk Assessment**: Intelligently assesses investment risks and provides warnings
- **Market Insights**: Analyzes market trends and hot sectors

## 🚀 Quick Start

### Environment Requirements

#### Backend Environment
- Go 1.24.0 or higher
- SQLite3

#### Frontend Environment
- Node.js 18.0 or higher
- npm 9.0 or higher

### Installation Steps

1. **Clone the Project**
   ```bash
   git clone https://github.com/yourusername/go-springAi.git
   cd go-springAi
   ```

2. **Backend Setup**
   ```bash
   # Install Go dependencies
   go mod download
   
   # Generate database code
   make sqlc-generate
   
   # Generate dependency injection code
   make wire-generate
   
   # Build the application
   make build
   ```

3. **Frontend Setup**
   ```bash
   cd frontend
   
   # Install Node.js dependencies
   npm install
   
   # Start development server
   npm run dev
   ```

4. **Configuration**
   
   Copy and modify the configuration file:
   ```bash
   cp config.yaml.example config.yaml
   ```
   
   Edit `config.yaml` and configure:
   - Database connection
   - AI provider API keys (OpenAI, Google AI)
   - Server ports and other settings

5. **Start Services**
   
   Start backend service:
   ```bash
   # Development mode
   make run
   
   # Or run the compiled binary
   ./bin/go-springAi
   ```
   
   Start frontend service (in another terminal):
   ```bash
   cd frontend
   npm run dev
   ```

6. **Access the Application**
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080
   - API Documentation: http://localhost:8080/swagger/index.html

## ⚙️ Configuration

### Backend Configuration (config.yaml)

```yaml
# Server configuration
server:
  port: 8080
  mode: debug  # debug, release, test

# Database configuration
database:
  driver: sqlite3
  dsn: "./data/go-springAi.db"

# AI Provider configuration
ai:
  providers:
    openai:
      enabled: true
      api_key: "your-openai-api-key"
      base_url: "https://api.openai.com/v1"
      models:
        - "gpt-4"
        - "gpt-3.5-turbo"
    
    googleai:
      enabled: true
      api_key: "your-google-ai-api-key"
      models:
        - "gemini-pro"
        - "gemini-pro-vision"

# MCP configuration
mcp:
  enabled: true
  tools:
    stock_analysis:
      enabled: true
    yahoo_finance:
      enabled: true
    stock_compare:
      enabled: true
    stock_advice:
      enabled: true

# Logging configuration
logging:
  level: info  # debug, info, warn, error
  format: json # json, text
  output: stdout # stdout, file
```

### Frontend Configuration

The frontend uses environment variables for configuration. Create a `.env` file in the `frontend` directory:

```env
# API base URL
VITE_API_BASE_URL=http://localhost:8080

# Application title
VITE_APP_TITLE=Stock Analysis AI Assistant

# Enable development features
VITE_DEV_MODE=true
```

## 🔑 API Key Management

### Setting Up AI Provider API Keys

1. **OpenAI API Key**
   - Visit [OpenAI Platform](https://platform.openai.com/api-keys)
   - Create a new API key
   - Add to configuration or set via the web interface

2. **Google AI API Key**
   - Visit [Google AI Studio](https://makersuite.google.com/app/apikey)
   - Create a new API key
   - Add to configuration or set via the web interface

3. **Dynamic Configuration**
   
   You can also set API keys through the web interface:
   - Navigate to Settings → Providers
   - Enter your API keys
   - Test the connection
   - Save the configuration

## 📖 Usage Guide

### Stock Analysis

1. **Basic Stock Analysis**
   ```bash
   # Using the web interface
   1. Navigate to "Stock Analysis" page
   2. Enter stock symbol (e.g., AAPL)
   3. Select analysis period
   4. Click "Analyze"
   
   # Using API directly
   curl -X POST http://localhost:8080/api/v1/mcp/tools/execute \
     -H "Content-Type: application/json" \
     -d '{
       "tool_name": "stock_analysis",
       "parameters": {
         "symbol": "AAPL",
         "period": "1mo"
       }
     }'
   ```

2. **Stock Comparison**
   ```bash
   # Compare multiple stocks
   curl -X POST http://localhost:8080/api/v1/mcp/tools/execute \
     -H "Content-Type: application/json" \
     -d '{
       "tool_name": "stock_compare",
       "parameters": {
         "symbols": ["AAPL", "GOOGL", "MSFT"],
         "period": "3mo"
       }
     }'
   ```

3. **AI Assistant Chat**
   ```bash
   # Chat with AI assistant
   curl -X POST http://localhost:8080/api/v1/ai/chat \
     -H "Content-Type: application/json" \
     -d '{
       "provider": "openai",
       "model": "gpt-4",
       "messages": [
         {
           "role": "user",
           "content": "Analyze AAPL stock and provide investment advice"
         }
       ],
       "tools_enabled": true
     }'
   ```

### MCP Tools Usage

The project implements several MCP tools for stock analysis:

1. **List Available Tools**
   ```bash
   curl http://localhost:8080/api/v1/mcp/tools
   ```

2. **Get Tool Definition**
   ```bash
   curl http://localhost:8080/api/v1/mcp/tools/stock_analysis
   ```

3. **Execute Tool**
   ```bash
   curl -X POST http://localhost:8080/api/v1/mcp/tools/execute \
     -H "Content-Type: application/json" \
     -d '{
       "tool_name": "yahoo_finance",
       "parameters": {
         "symbol": "TSLA",
         "data_type": "info"
       }
     }'
   ```

## 🏗️ Architecture Overview

### System Architecture

The project follows a clean architecture pattern with clear separation of concerns:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend       │    │   External      │
│   (React)       │    │   (Go/Gin)      │    │   Services      │
├─────────────────┤    ├─────────────────┤    ├─────────────────┤
│ • UI Components │    │ • Controllers   │    │ • OpenAI API    │
│ • State Mgmt    │◄──►│ • Services      │◄──►│ • Google AI     │
│ • API Clients   │    │ • Repositories  │    │ • Yahoo Finance │
│ • Routing       │    │ • MCP Tools     │    │ • SQLite DB     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### MCP Client Architecture Analysis

#### Core Interfaces

The MCP client system is built around three core interfaces:

```go
// MCPServiceInterface - Main service interface
type MCPServiceInterface interface {
    ListTools(ctx context.Context) ([]dto.MCPTool, error)
    GetTool(ctx context.Context, name string) (*dto.MCPTool, error)
    ExecuteTool(ctx context.Context, req *dto.MCPExecuteRequest) (*dto.MCPExecuteResponse, error)
    GetExecutionLogs(ctx context.Context, limit int) ([]dto.MCPExecutionLog, error)
}

// InternalMCPClient - Internal client interface
type InternalMCPClient interface {
    ListTools(ctx context.Context) ([]dto.MCPTool, error)
    GetTool(ctx context.Context, name string) (*dto.MCPTool, error)
    ExecuteTool(ctx context.Context, req *dto.MCPExecuteRequest) (*dto.MCPExecuteResponse, error)
}

// Tool - Tool interface
type Tool interface {
    GetDefinition() dto.MCPTool
    Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)
    Validate(params map[string]interface{}) error
}
```

#### Core Implementation Classes

1. **InternalMCPClientImpl** - Internal MCP client implementation
   - Responsible for in-process communication optimization
   - Direct method calls without network overhead
   - Manages tool registry and execution

2. **MCPServiceImpl** - MCP service implementation
   - Implements the MCPServiceInterface
   - Handles HTTP requests and responses
   - Manages execution logging and monitoring

3. **MCPClientManager** - Client manager
   - Manages multiple MCP client instances
   - Handles client lifecycle and configuration
   - Provides client discovery and routing

#### Tool System Architecture

1. **ToolRegistry** - Tool registration and discovery
   - Maintains a registry of available tools
   - Supports dynamic tool registration
   - Provides tool lookup and validation

2. **BaseTool** - Base tool implementation
   - Common functionality for all tools
   - Parameter validation and error handling
   - Execution context management

3. **Specific Tool Implementations**:
   - **StockAnalysisTool** - Stock analysis functionality
   - **YahooFinanceTool** - Yahoo Finance data retrieval
   - **StockCompareTool** - Stock comparison analysis

#### Main Call Paths

1. **HTTP API to MCP Service**:
   ```
   HTTP Request → MCPController → MCPService → InternalMCPClient → Tool → Response
   ```

2. **AI Assistant to MCP Service**:
   ```
   AI Request → AIAssistantService → MCPService → Tool Execution → AI Response
   ```

### Detailed Call Flow

#### Tool Execution Flow
1. **Request Reception**: HTTP request received by MCPController
2. **Parameter Validation**: Request parameters validated and parsed
3. **Service Call**: MCPService.ExecuteTool() called
4. **Tool Lookup**: Tool found in ToolRegistry
5. **Parameter Validation**: Tool-specific parameter validation
6. **Tool Execution**: Tool.Execute() method called
7. **Result Processing**: Execution result processed and formatted
8. **Logging**: Execution details logged to database
9. **Event Broadcasting**: Execution events broadcast via SSE
10. **Response Return**: Final response returned to client

#### AI Assistant Integration Flow
1. **Chat Request**: User sends chat message
2. **Tool Discovery**: Available tools discovered from MCP service
3. **System Message**: System message with tool definitions prepared
4. **AI Call**: AI provider called with tools and context
5. **Tool Call Parsing**: AI response parsed for tool calls
6. **Tool Execution**: Required tools executed via MCP service
7. **Result Integration**: Tool results integrated into AI context
8. **Final Response**: Complete response returned to user

### Data Flow Transformation

#### DTO Layer Data Structures
- **MCPTool**: Tool definition and metadata
- **MCPExecuteRequest**: Tool execution request
- **MCPExecuteResponse**: Tool execution response
- **MCPExecutionLog**: Execution history record

#### Inter-Service Data Transfer
- HTTP layer uses JSON serialization
- Internal calls use Go struct direct transfer
- Database layer uses SQLC generated types
- AI providers use provider-specific formats

### Dependency Injection Relationships

#### Wire Dependency Graph
```go
// Key dependency relationships
MCPController ← MCPService ← InternalMCPClient ← ToolRegistry
AIAssistantService ← MCPService
MCPService ← Logger ← Database
```

#### Key Dependency Relationships
- Controllers depend on Services
- Services depend on Repositories and External APIs
- MCP components depend on Tool implementations
- All components depend on Logger and Configuration

## MCP System Design Patterns

### Adapter Pattern
**Purpose**: Adapt different AI providers to a unified interface

**Implementation**:
```go
// AI Provider Adapter
type AIProvider interface {
    ChatCompletion(ctx context.Context, req ChatRequest) (ChatResponse, error)
    ListModels(ctx context.Context) ([]Model, error)
}

// Google AI Adapter
type GoogleAIAdapter struct {
    client *genai.Client
}

func (g *GoogleAIAdapter) ChatCompletion(ctx context.Context, req ChatRequest) (ChatResponse, error) {
    // Convert unified request to Google AI format
    // Call Google AI API
    // Convert response back to unified format
}
```

**Advantages**:
- Unified interface for different AI providers
- Easy to add new AI providers
- Consistent behavior across providers

### Registry Pattern
**Purpose**: Manage tool registration and discovery

**Implementation**:
```go
// Tool Registry
type ToolRegistry struct {
    tools map[string]Tool
    mutex sync.RWMutex
}

func (r *ToolRegistry) Register(name string, tool Tool) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    r.tools[name] = tool
    return nil
}

// Client Manager
type MCPClientManager struct {
    clients map[string]InternalMCPClient
    registry *ToolRegistry
}
```

**Advantages**:
- Dynamic tool registration
- Centralized tool management
- Easy tool discovery

### Strategy Pattern
**Purpose**: Support different AI provider strategies

**Implementation**:
```go
// AI Provider Strategy Interface
type AIProviderStrategy interface {
    Execute(ctx context.Context, req AIRequest) (AIResponse, error)
    GetCapabilities() []string
}

// AI Controller using Strategy
type AIController struct {
    strategies map[string]AIProviderStrategy
}

func (c *AIController) ProcessRequest(provider string, req AIRequest) (AIResponse, error) {
    strategy := c.strategies[provider]
    return strategy.Execute(ctx, req)
}
```

**Advantages**:
- Flexible provider switching
- Easy to add new strategies
- Runtime strategy selection

### Observer Pattern
**Purpose**: Real-time event notification system

**Implementation**:
```go
// SSE Client Management
type SSEClientManager struct {
    clients map[string]chan []byte
    mutex   sync.RWMutex
}

func (m *SSEClientManager) Broadcast(event Event) {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    
    data, _ := json.Marshal(event)
    for _, client := range m.clients {
        select {
        case client <- data:
        default:
            // Client buffer full, skip
        }
    }
}
```

**Advantages**:
- Real-time event broadcasting
- Decoupled event system
- Multiple subscriber support

### Dependency Injection Pattern
**Purpose**: Manage component dependencies

**Implementation**:
```go
// Wire Provider Definitions
func ProvideMCPService(
    client InternalMCPClient,
    logger *zap.Logger,
    db *sql.DB,
) MCPServiceInterface {
    return &MCPServiceImpl{
        client: client,
        logger: logger,
        db:     db,
    }
}
```

**Advantages**:
- Loose coupling between components
- Easy testing with mock dependencies
- Centralized dependency management

## Core Features

### In-Process Communication Optimization
- **Direct Method Calls**: `InternalMCPClientImpl` directly calls service methods without network overhead
- **Memory Sharing**: Shared data structures between components
- **Performance Benefits**: Eliminates serialization/deserialization overhead

### Modular Tool System
- **Tool Definition**: Clear interface for tool implementation
- **Tool Registration**: Dynamic tool registration at startup
- **Tool Discovery**: Runtime tool discovery and listing
- **Tool Execution**: Unified tool execution framework
- **Tool Monitoring**: Execution logging and performance monitoring

### Type Safety Guarantees
- **Interface Layer**: Strong typing through Go interfaces
- **Implementation Layer**: Type-safe implementations
- **Validation Layer**: Parameter validation and type checking

### Real-time Event System
- **SSE Long Connections**: Server-Sent Events for real-time communication
- **Event Types**: 
  - `tool_execution_start`: Tool execution started
  - `tool_execution_complete`: Tool execution completed
  - `tool_execution_error`: Tool execution error
  - `system_status_change`: System status change

### Complete Observability
- **Structured Logging**: Zap-based structured logging
- **Performance Monitoring**: Execution time and resource usage tracking
- **Error Tracking**: Comprehensive error logging and tracking
- **MCPExecutionLog Structure**:
  ```go
  type MCPExecutionLog struct {
      ID          int64     `json:"id"`
      ToolName    string    `json:"tool_name"`
      Parameters  string    `json:"parameters"`
      Result      string    `json:"result"`
      Error       string    `json:"error,omitempty"`
      Duration    int64     `json:"duration_ms"`
      ExecutedAt  time.Time `json:"executed_at"`
  }
  ```

### Security Guarantees
- **Parameter Validation**: Input parameter validation and sanitization
- **Sensitive Information Filtering**: Automatic filtering of sensitive data in logs
- **Exception Handling**: Comprehensive exception handling and error recovery
- **Attack Detection**: Basic attack pattern detection and prevention

#### Security Measures
1. **Input Parameter Validation**
   - Type validation for all input parameters
   - Range validation for numeric parameters
   - Format validation for string parameters
   - SQL injection prevention

2. **Sensitive Information Filtering**
   - API keys automatically masked in logs
   - Personal information redacted
   - Financial data anonymization

3. **Exception Handling**
   - Graceful error handling
   - Error message sanitization
   - Stack trace filtering

4. **Attack Detection**
   - Rate limiting
   - Suspicious pattern detection
   - Request size limits

## MCP Data Flow and Execution Flow

### Complete Tool Execution Flow Diagram

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Request  │    │  MCPController  │    │   MCPService    │
│                 │───►│                 │───►│                 │
│ POST /tools/    │    │ • Route Handler │    │ • Business      │
│ execute         │    │ • Input Valid.  │    │   Logic         │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
                                                        ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Tool Result   │    │InternalMCPClient│    │  ToolRegistry   │
│                 │◄───│                 │◄───│                 │
│ • Formatted     │    │ • Tool Lookup   │    │ • Tool Storage  │
│ • Logged        │    │ • Execution     │    │ • Tool Discovery│
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
                                                        ▼
                               ┌─────────────────┐    ┌─────────────────┐
                               │   Tool Execute  │    │   Specific Tool │
                               │                 │───►│                 │
                               │ • Parameter     │    │ • Stock Analysis│
                               │   Validation    │    │ • Yahoo Finance │
                               └─────────────────┘    └─────────────────┘
```

### AI Assistant Integration Flow Diagram

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   User Chat     │    │AIAssistantCtrl  │    │AIAssistantSvc   │
│                 │───►│                 │───►│                 │
│ "Analyze AAPL"  │    │ • Chat Handler  │    │ • Context Mgmt  │
│                 │    │ • SSE Support   │    │ • Tool Discovery│
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
                                                        ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Final Response│    │   AI Provider   │    │   MCP Service   │
│                 │◄───│                 │◄───│                 │
│ • Analysis      │    │ • OpenAI/Google │    │ • Tool Execution│
│ • Recommendations│   │ • Tool Calling  │    │ • Result Format │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### SSE Event Flow Diagram

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Client    │    │  MCPController  │    │SSEClientManager │
│                 │───►│                 │───►│                 │
│ GET /sse        │    │ • SSE Handler   │    │ • Client Mgmt   │
│                 │◄───│ • Event Stream  │◄───│ • Event Routing │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        ▲
                                                        │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Event Types   │    │   MCPService    │    │   Event Bus     │
│                 │    │                 │───►│                 │
│ • tool_execution│    │ • Event Trigger │    │ • Event Broadcast│
│ • system_status │    │ • Real-time Push│    │ • Multi-cast    │
└─────────────────┘    └─────────────────┘    └─────────────────┘

Event Types:
- tool_execution_start: {"tool_name": "stock_analysis", "timestamp": "..."}
- tool_execution_complete: {"tool_name": "stock_analysis", "result": "...", "duration": 1500}
- tool_execution_error: {"tool_name": "stock_analysis", "error": "...", "timestamp": "..."}
- system_status_change: {"status": "healthy", "timestamp": "..."}
```

## Development Guide

### Code Generation

The project uses code generation for database access and dependency injection:

1. **SQLC Generation**
   ```bash
   # Generate type-safe database code
   sqlc generate
   ```

2. **Wire Generation**
   ```bash
   # Generate dependency injection code
   wire
   ```

### Adding New MCP Tools

To add a new MCP tool, follow these steps:

1. **Define Tool Structure**
   ```go
   type MyTool struct {
       logger *zap.Logger
       // Add other dependencies
   }
   ```

2. **Implement Tool Interface**
   ```go
   func (t *MyTool) GetDefinition() dto.MCPTool {
       return dto.MCPTool{
           Name:        "my_tool",
           Description: "Description of my tool",
           Parameters: map[string]interface{}{
               "type": "object",
               "properties": map[string]interface{}{
                   "param1": map[string]interface{}{
                       "type":        "string",
                       "description": "Parameter description",
                   },
               },
               "required": []string{"param1"},
           },
       }
   }

   func (t *MyTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
       // Implement tool logic
       return result, nil
   }

   func (t *MyTool) Validate(params map[string]interface{}) error {
       // Implement parameter validation
       return nil
   }
   ```

3. **Register Tool**
   ```go
   // In registerDefaultTools() function
   registry.Register("my_tool", &MyTool{
       logger: logger,
   })
   ```

### Frontend Development Guide

#### Development Environment Setup

1. **Install Dependencies**
   ```bash
   cd frontend
   npm install
   ```

2. **Start Development Server**
   ```bash
   npm run dev
   ```

3. **Build for Production**
   ```bash
   npm run build
   ```

4. **Preview Production Build**
   ```bash
   npm run preview
   ```

#### Adding New Pages

1. **Create Component**
   ```typescript
   // src/pages/MyPage.tsx
   import React from 'react';
   import { Card, Typography } from 'antd';

   const MyPage: React.FC = () => {
     return (
       <Card title="My Page">
         <Typography.Text>Page content</Typography.Text>
       </Card>
     );
   };

   export default MyPage;
   ```

2. **Configure Route**
   ```typescript
   // src/router/index.tsx
   import MyPage from '../pages/MyPage';

   const routes = [
     // ... existing routes
     {
       path: '/my-page',
       element: <MyPage />,
     },
   ];
   ```

3. **Add Menu Item**
   ```typescript
   // src/components/Layout/index.tsx
   const menuItems = [
     // ... existing items
     {
       key: 'my-page',
       label: 'My Page',
       icon: <SomeIcon />,
     },
   ];
   ```

#### State Management

Using Redux Toolkit for state management:

```typescript
// src/store/slices/mySlice.ts
import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface MyState {
  data: any[];
  loading: boolean;
}

const initialState: MyState = {
  data: [],
  loading: false,
};

const mySlice = createSlice({
  name: 'my',
  initialState,
  reducers: {
    setData: (state, action: PayloadAction<any[]>) => {
      state.data = action.payload;
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
  },
});

export const { setData, setLoading } = mySlice.actions;
export default mySlice.reducer;
```

#### API Services

API services are located in `src/services/`:

```typescript
// src/services/myService.ts
import { api } from './api';

export const myService = {
  getData: async () => {
    const response = await api.get('/my-endpoint');
    return response.data;
  },
  
  postData: async (data: any) => {
    const response = await api.post('/my-endpoint', data);
    return response.data;
  },
};
```

#### Type Definitions

Type definitions are located in `src/types/`:

```typescript
// src/types/my.ts
export interface MyData {
  id: string;
  name: string;
  value: number;
}

export interface MyApiResponse {
  data: MyData[];
  total: number;
}
```

### Backend Development Guide

#### Project Architecture Description

The project follows clean architecture principles for the stock analysis AI assistant:

- **Stock Controllers**: Handle HTTP requests for stock analysis endpoints
- **Stock Services**: Implement business logic for stock analysis and AI integration
- **Financial Repository**: Data access layer for financial data and user management
- **Stock Models/DTO**: Data transfer objects for stock analysis requests and responses
- **Financial Middleware**: Authentication, logging, and validation for financial operations
- **Stock Analysis Tools**: MCP tools for stock analysis, comparison, and advice

The project implements complete MCP protocol support with:
- **Stock Tool Registration and Discovery**: Dynamic registration of stock analysis tools
- **Secure Stock Analysis Execution**: Safe execution of financial analysis tools
- **Real-time Stock Data Streaming**: SSE-based real-time stock data updates
- **Investment Decision Logging**: Complete audit trail of analysis and recommendations
- **Financial Error Handling**: Specialized error handling for financial operations

#### Unified Error Handling

The project uses a unified error handling system:

```go
// AppError structure
type AppError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

// MCP-specific errors
var (
    ErrToolNotFound     = &AppError{Code: "TOOL_NOT_FOUND", Message: "Tool not found"}
    ErrToolExecution    = &AppError{Code: "TOOL_EXECUTION_ERROR", Message: "Tool execution failed"}
    ErrInvalidParams    = &AppError{Code: "INVALID_PARAMETERS", Message: "Invalid parameters"}
)

// Error middleware
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            // Handle different error types
            // Return appropriate HTTP response
        }
    }
}
```

#### Structured Logging

The project uses Zap for structured logging:

```go
// API request/response logging
logger.Info("API request",
    zap.String("method", c.Request.Method),
    zap.String("path", c.Request.URL.Path),
    zap.String("user_id", userID),
    zap.Duration("duration", duration),
)

// MCP tool execution logging
logger.Info("MCP tool execution",
    zap.String("tool_name", toolName),
    zap.Any("parameters", params),
    zap.Duration("execution_time", executionTime),
    zap.String("result_status", status),
)

// Performance monitoring
logger.Info("Performance metrics",
    zap.String("operation", "stock_analysis"),
    zap.Duration("duration", duration),
    zap.Int("memory_usage", memUsage),
)

// Security events
logger.Warn("Security event",
    zap.String("event_type", "invalid_api_key"),
    zap.String("ip_address", clientIP),
    zap.String("user_agent", userAgent),
)

// Error and exception logging
logger.Error("Application error",
    zap.Error(err),
    zap.String("context", "stock_analysis"),
    zap.Any("request_data", requestData),
)
```

## Testing

### Using Makefile

The project includes a comprehensive Makefile for testing:

```bash
# Show available commands
make help

# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests
make test-integration

# Run tests with coverage
make test-coverage

# Run race condition tests
make test-race

# Generate mocks
make mock-gen

# Clean test artifacts
make clean
```

### Using Go Commands Directly

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run specific test
go test -run TestSpecificFunction ./internal/service

# Run tests with verbose output
go test -v ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Build and Deployment

### Frontend Build

```bash
cd frontend

# Install dependencies
npm install

# Build for production
npm run build
```

### Backend Build

```bash
# Build using Go
go build -o bin/go-springAi cmd/main.go

# Build using Makefile
make build

# Cross-compilation
GOOS=linux GOARCH=amd64 go build -o bin/go-springAi-linux-amd64 cmd/main.go
GOOS=windows GOARCH=amd64 go build -o bin/go-springAi-windows-amd64.exe cmd/main.go
GOOS=darwin GOARCH=arm64 go build -o bin/go-springAi-darwin-arm64 cmd/main.go
```

## Contributing

We welcome contributions! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure your code follows the project's coding standards and includes appropriate tests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

⭐ If you find this project helpful, please consider giving it a star!