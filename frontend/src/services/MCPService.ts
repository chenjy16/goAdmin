import { BaseService } from './base/BaseService';
import type { 
  BaseApiResponse, 
  IInitializable, 
  IValidatable, 
  ValidationResult 
} from '../types/base';
import type { 
  MCPTool, 
  MCPToolsResponse, 
  MCPExecuteRequest, 
  MCPExecuteResponse, 
  MCPMessage 
} from '../types/api';

/**
 * MCP服务类
 * 负责处理MCP工具相关的API调用
 */
export class MCPService extends BaseService implements IInitializable, IValidatable<MCPExecuteRequest> {
  private _initialized = false;

  /**
   * 初始化MCP系统
   */
  async initialize(): Promise<void> {
    try {
      const response = await this.post<BaseApiResponse>('/api/v1/mcp/initialize');
      if (response.code === 200) {
        this._initialized = true;
      } else {
        throw new Error(response.message || 'Failed to initialize MCP');
      }
    } catch (error) {
      this._initialized = false;
      throw error;
    }
  }

  /**
   * 检查是否已初始化
   */
  isInitialized(): boolean {
    return this._initialized;
  }

  /**
   * 获取MCP工具列表
   */
  async getTools(): Promise<MCPToolsResponse> {
    return this.get<MCPToolsResponse>('/api/v1/mcp/tools');
  }

  /**
   * 执行MCP工具
   */
  async executeTool(request: MCPExecuteRequest): Promise<MCPExecuteResponse> {
    const validationResult = this.validate(request);
    if (!validationResult.isValid) {
      throw new Error(`Validation failed: ${validationResult.errors.map(e => e.message).join(', ')}`);
    }
    return this.post<MCPExecuteResponse>('/api/v1/mcp/execute', request);
  }

  /**
   * 获取执行日志
   */
  async getLogs(): Promise<MCPMessage[]> {
    const response = await this.get<BaseApiResponse<MCPMessage[]>>('/api/v1/mcp/logs');
    return response.data || [];
  }

  /**
   * 获取单个执行日志
   */
  async getLog(id: string): Promise<MCPMessage> {
    const response = await this.get<BaseApiResponse<MCPMessage>>(`/api/v1/mcp/logs/${id}`);
    if (!response.data) {
      throw new Error('Log not found');
    }
    return response.data;
  }

  /**
   * 清除执行日志
   */
  async clearLogs(): Promise<BaseApiResponse> {
    return this.delete<BaseApiResponse>('/api/v1/mcp/logs');
  }

  /**
   * 创建SSE连接
   */
  createEventSource(): EventSource {
    const baseURL = this.baseURL || '';
    return new EventSource(`${baseURL}/api/v1/mcp/sse`);
  }

  /**
   * 获取MCP状态
   */
  async getStatus(): Promise<BaseApiResponse<{ initialized: boolean; toolCount: number; lastActivity?: string }>> {
    return this.get<BaseApiResponse<{ initialized: boolean; toolCount: number; lastActivity?: string }>>('/api/v1/mcp/status');
  }

  /**
   * 重启MCP系统
   */
  async restart(): Promise<BaseApiResponse> {
    try {
      const response = await this.post<BaseApiResponse>('/api/v1/mcp/restart');
      this._initialized = false;
      return response;
    } catch (error) {
      this._initialized = false;
      throw error;
    }
  }

  /**
   * 验证MCP执行请求
   */
  validate(request: MCPExecuteRequest): ValidationResult {
    const errors = [];

    if (!request.name || request.name.trim().length === 0) {
      errors.push({
        field: 'name',
        message: 'Tool name is required',
        code: 'REQUIRED'
      });
    }

    if (!request.arguments || typeof request.arguments !== 'object') {
      errors.push({
        field: 'arguments',
        message: 'Arguments must be an object',
        code: 'INVALID_TYPE'
      });
    }

    // 验证工具名称格式
    if (request.name && !/^[a-zA-Z][a-zA-Z0-9_-]*$/.test(request.name)) {
      errors.push({
        field: 'name',
        message: 'Tool name must start with a letter and contain only letters, numbers, underscores, and hyphens',
        code: 'INVALID_FORMAT'
      });
    }

    return {
      isValid: errors.length === 0,
      errors
    };
  }

  /**
   * 验证工具参数
   */
  validateToolArguments(tool: MCPTool, arguments_: Record<string, any>): ValidationResult {
    const errors = [];
    const schema = tool.inputSchema;

    if (!schema || !schema.properties) {
      return { isValid: true, errors: [] };
    }

    // 检查必需参数
    if (schema.required) {
      for (const requiredField of schema.required) {
        if (!(requiredField in arguments_)) {
          errors.push({
            field: requiredField,
            message: `Required parameter '${requiredField}' is missing`,
            code: 'REQUIRED'
          });
        }
      }
    }

    // 检查参数类型
    for (const [fieldName, fieldValue] of Object.entries(arguments_)) {
      const fieldSchema = schema.properties[fieldName];
      if (fieldSchema && fieldSchema.type) {
        const expectedType = fieldSchema.type;
        const actualType = typeof fieldValue;
        
        if (expectedType === 'string' && actualType !== 'string') {
          errors.push({
            field: fieldName,
            message: `Parameter '${fieldName}' must be a string`,
            code: 'INVALID_TYPE'
          });
        } else if (expectedType === 'number' && actualType !== 'number') {
          errors.push({
            field: fieldName,
            message: `Parameter '${fieldName}' must be a number`,
            code: 'INVALID_TYPE'
          });
        } else if (expectedType === 'boolean' && actualType !== 'boolean') {
          errors.push({
            field: fieldName,
            message: `Parameter '${fieldName}' must be a boolean`,
            code: 'INVALID_TYPE'
          });
        } else if (expectedType === 'array' && !Array.isArray(fieldValue)) {
          errors.push({
            field: fieldName,
            message: `Parameter '${fieldName}' must be an array`,
            code: 'INVALID_TYPE'
          });
        }
      }
    }

    return {
      isValid: errors.length === 0,
      errors
    };
  }

  /**
   * 格式化工具执行结果
   */
  formatExecutionResult(response: MCPExecuteResponse): string {
    if (response.isError) {
      return `Error: ${response.content.map(c => c.text || JSON.stringify(c.data)).join('\n')}`;
    }

    return response.content.map(content => {
      if (content.type === 'text') {
        return content.text || '';
      } else {
        return JSON.stringify(content.data, null, 2);
      }
    }).join('\n');
  }

  /**
   * 获取工具使用统计
   */
  async getToolUsageStats(): Promise<BaseApiResponse<Record<string, { count: number; lastUsed: string }>>> {
    return this.get<BaseApiResponse<Record<string, { count: number; lastUsed: string }>>>('/api/v1/mcp/stats');
  }
}

export const mcpService = new MCPService();
export default mcpService;