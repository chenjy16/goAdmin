// API响应基础类型
export interface ApiResponse<T = any> {
  data?: T;
  message?: string;
  status?: string;
  error?: string;
}



// 提供商相关类型
export interface ProviderInfo {
  type: string;
  name: string;
  description: string;
  healthy: boolean;
  model_count: number;
}

export interface ProvidersResponse {
  code: number;
  message: string;
  data: {
    providers: ProviderInfo[];
  };
}

// 模型相关类型
export interface ModelInfo {
  name: string;
  display_name: string;
  max_tokens: number;
  temperature: number;
  top_p: number;
  top_k?: number;
  enabled: boolean;
  description?: string;
}

export interface ModelsResponse {
  code: number;
  message: string;
  data: {
    models: Record<string, ModelInfo>;
    provider: string;
  };
}

export interface ModelConfigResponse {
  model: ModelInfo;
}

// API密钥相关类型
export interface SetAPIKeyRequest {
  api_key: string;
}

export interface ValidateAPIKeyResponse {
  valid: boolean;
  message?: string;
}

export interface APIKeyInfo {
  has_key: boolean;
  masked_key?: string;
}

export interface APIKeyStatusResponse {
  code: number;
  message: string;
  data: Record<string, APIKeyInfo>;
}

// MCP工具相关类型
export interface MCPTool {
  name: string;
  description: string;
  inputSchema: any;
}

export interface MCPToolsResponse {
  code: number;
  message: string;
  data: {
    tools: MCPTool[];
  };
}

export interface MCPExecuteRequest {
  name: string;
  arguments: Record<string, any>;
}

export interface MCPContent {
  type: string;
  text?: string;
  data?: any;
}

export interface MCPExecuteResponse {
  content: MCPContent[];
  isError?: boolean;
}

export interface MCPMessage {
  id: string;
  timestamp: string;
  level: 'info' | 'warn' | 'error' | 'debug';
  message: string;
  data?: any;
}

export interface MCPError {
  code: number;
  message: string;
  data?: any;
}

// AI助手相关类型
export interface ChatMessage {
  role: 'user' | 'assistant' | 'system';
  content: string;
  timestamp?: string;
  tool_calls?: any[];
}

export interface ChatRequest {
  messages: ChatMessage[];
  model?: string;
  max_tokens?: number;
  temperature?: number;
  tools?: any[];
  tool_choice?: string;
  use_tools?: boolean;
  provider?: string;
  selected_tool?: string;
}

export interface ChatResponse {
  id: string;
  object: string;
  created: number;
  model: string;
  choices: {
    index: number;
    message: ChatMessage;
    finish_reason: string;
  }[];
  usage?: {
    prompt_tokens: number;
    completion_tokens: number;
    total_tokens: number;
  };
}

// SSE事件类型
export interface SSEEvent {
  id?: string;
  event?: string;
  data: string;
  retry?: number;
}

// 健康检查类型
export interface HealthResponse {
  message: string;
  status: string;
}