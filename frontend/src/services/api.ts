import axios from 'axios';
import type { AxiosInstance, AxiosResponse } from 'axios';
import type {
  ApiResponse,
  ProvidersResponse,
  ModelsResponse,
  ModelConfigResponse,
  SetAPIKeyRequest,
  ValidateAPIKeyResponse,
  MCPToolsResponse,
  MCPExecuteRequest,
  MCPExecuteResponse,
  MCPMessage,
  ChatRequest,
  ChatResponse,
  HealthResponse
} from '../types/api';

class ApiService {
  private api: AxiosInstance;
  private baseURL: string;

  constructor() {
    this.baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
    this.api = axios.create({
      baseURL: this.baseURL,
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // 请求拦截器
    this.api.interceptors.request.use(
      (config) => {
        // 可以在这里添加认证token等
        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );

    // 响应拦截器
    this.api.interceptors.response.use(
      (response: AxiosResponse) => {
        return response;
      },
      (error) => {
        // 统一错误处理
        const errorMessage = error.response?.data?.message || error.message || '请求失败';
        return Promise.reject(new Error(errorMessage));
      }
    );
  }

  // 健康检查
  async healthCheck(): Promise<HealthResponse> {
    const response = await this.api.get<HealthResponse>('/health');
    return response.data;
  }

  // 提供商相关API
  async getProviders(): Promise<ProvidersResponse> {
    const response = await this.api.get<ProvidersResponse>('/api/v1/ai/providers');
    return response.data;
  }

  // 模型相关API
  async getModels(provider: string): Promise<ModelsResponse> {
    const response = await this.api.get<ModelsResponse>(`/api/v1/ai/${provider}/models`);
    return response.data;
  }

  async getModelConfig(provider: string, model: string): Promise<ModelConfigResponse> {
    const response = await this.api.get<ModelConfigResponse>(`/api/v1/ai/${provider}/config/${model}`);
    return response.data;
  }

  async enableModel(provider: string, model: string): Promise<ApiResponse> {
    const response = await this.api.put<ApiResponse>(`/api/v1/ai/${provider}/models/${model}/enable`);
    return response.data;
  }

  async disableModel(provider: string, model: string): Promise<ApiResponse> {
    const response = await this.api.put<ApiResponse>(`/api/v1/ai/${provider}/models/${model}/disable`);
    return response.data;
  }

  // API密钥相关API
  async setAPIKey(provider: string, apiKey: string): Promise<ApiResponse> {
    const response = await this.api.post<ApiResponse>(`/api/v1/ai/${provider}/api-key`, {
      api_key: apiKey
    } as SetAPIKeyRequest);
    return response.data;
  }

  async validateAPIKey(provider: string): Promise<ValidateAPIKeyResponse> {
    const response = await this.api.post<ValidateAPIKeyResponse>(`/api/v1/ai/${provider}/validate`);
    return response.data;
  }



  // MCP工具相关API
  async initializeMCP(): Promise<ApiResponse> {
    const response = await this.api.post<ApiResponse>('/api/v1/mcp/initialize');
    return response.data;
  }

  async getMCPTools(): Promise<MCPToolsResponse> {
    const response = await this.api.get<MCPToolsResponse>('/api/v1/mcp/tools');
    return response.data;
  }

  async executeMCPTool(request: MCPExecuteRequest): Promise<MCPExecuteResponse> {
    const response = await this.api.post<MCPExecuteResponse>('/api/v1/mcp/execute', request);
    return response.data;
  }

  async getMCPLogs(): Promise<MCPMessage[]> {
    const response = await this.api.get<MCPMessage[]>('/api/v1/mcp/logs');
    return response.data;
  }

  async getMCPLog(id: string): Promise<MCPMessage> {
    const response = await this.api.get<MCPMessage>(`/api/v1/mcp/logs/${id}`);
    return response.data;
  }

  // MCP SSE事件流
  createMCPEventSource(): EventSource {
    return new EventSource(`${this.baseURL}/api/v1/mcp/sse`);
  }

  // AI助手相关API
  async initializeAssistant(): Promise<ApiResponse> {
    const response = await this.api.post<ApiResponse>('/api/v1/assistant/initialize');
    return response.data;
  }

  async assistantChat(request: ChatRequest): Promise<ChatResponse> {
    const response = await this.api.post<ChatResponse>('/api/v1/assistant/chat', request);
    return response.data;
  }
}

// 创建单例实例
export const apiService = new ApiService();
export default apiService;