import { BaseService } from './base/BaseService';
import { serviceFactory, getProviderService, getMCPService, getAssistantService } from './ServiceFactory';
import type {
  ApiResponse,
  ProvidersResponse,
  ModelsResponse,
  ModelConfigResponse,
  ValidateAPIKeyResponse,
  MCPToolsResponse,
  MCPExecuteRequest,
  MCPExecuteResponse,
  MCPMessage,
  ChatRequest,
  ChatResponse,
  HealthResponse,
  APIKeyInfo
} from '../types/api';

/**
 * 统一API服务类
 * 作为所有服务的统一入口点，保持向后兼容性
 * @deprecated 建议直接使用具体的服务类（ProviderService, MCPService, AssistantService）
 */
class ApiService extends BaseService {
  constructor() {
    super();
    // 确保服务工厂已初始化
    this.ensureServicesInitialized();
  }

  private async ensureServicesInitialized(): Promise<void> {
    if (!serviceFactory.isInitialized()) {
      try {
        await serviceFactory.initialize();
      } catch (error) {
        console.warn('Failed to initialize services:', error);
      }
    }
  }

  // 健康检查
  async healthCheck(): Promise<HealthResponse> {
    return this.get<HealthResponse>('/health');
  }

  // 提供商相关API - 委托给ProviderService
  async getProviders(): Promise<ProvidersResponse> {
    return getProviderService().getProviders();
  }

  async getModels(provider: string): Promise<ModelsResponse> {
    return getProviderService().getModels(provider);
  }

  async getAllModels(provider: string): Promise<ModelsResponse> {
    return getProviderService().getAllModels(provider);
  }

  async getModelConfig(provider: string, model: string): Promise<ModelConfigResponse> {
    return getProviderService().getModelConfig(provider, model);
  }

  async enableModel(provider: string, model: string): Promise<ApiResponse> {
    return getProviderService().enableModel(provider, model);
  }

  async disableModel(provider: string, model: string): Promise<ApiResponse> {
    return getProviderService().disableModel(provider, model);
  }

  async setAPIKey(provider: string, apiKey: string): Promise<ApiResponse> {
    return getProviderService().setAPIKey(provider, apiKey);
  }

  async validateAPIKey(provider: string): Promise<ValidateAPIKeyResponse> {
    return getProviderService().validateAPIKey(provider);
  }

  async getAPIKeyStatus(): Promise<ApiResponse<Record<string, APIKeyInfo>>> {
    return getProviderService().getAPIKeyStatus();
  }

  async getPlainAPIKey(provider: string): Promise<ApiResponse<{ provider: string; api_key: string }>> {
    return getProviderService().getPlainAPIKey(provider);
  }

  // MCP工具相关API - 委托给MCPService
  async initializeMCP(): Promise<ApiResponse> {
    const mcpService = getMCPService();
    await mcpService.initialize();
    return { message: 'MCP initialized successfully', status: 'success' };
  }

  async getMCPTools(): Promise<MCPToolsResponse> {
    return getMCPService().getTools();
  }

  async executeMCPTool(request: MCPExecuteRequest): Promise<MCPExecuteResponse> {
    return getMCPService().executeTool(request);
  }

  async getMCPLogs(): Promise<MCPMessage[]> {
    return getMCPService().getLogs();
  }

  async getMCPLog(id: string): Promise<MCPMessage> {
    return getMCPService().getLog(id);
  }

  createMCPEventSource(): EventSource {
    return getMCPService().createEventSource();
  }

  async getMCPStatus(): Promise<{ initialized: boolean; toolCount: number; lastActivity?: string }> {
    const response = await getMCPService().getStatus();
    return response.data || { initialized: false, toolCount: 0 };
  }

  // AI助手相关API - 委托给AssistantService
  async initializeAssistant(): Promise<ApiResponse> {
    const assistantService = getAssistantService();
    await assistantService.initialize();
    return { message: 'Assistant initialized successfully', status: 'success' };
  }

  async assistantChat(request: ChatRequest): Promise<ChatResponse> {
    return getAssistantService().chat(request);
  }
}

// 创建单例实例
export const apiService = new ApiService();
export default apiService;

// 同时导出新的服务架构
export {
  serviceFactory,
  getProviderService,
  getMCPService,
  getAssistantService
} from './ServiceFactory';