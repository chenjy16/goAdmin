import { BaseService } from './base/BaseService';
import type { 
  BaseApiResponse, 
  IValidatable, 
  ValidationResult 
} from '../types/base';
import type { 
  ProvidersResponse, 
  ModelsResponse, 
  ModelConfigResponse,
  ValidateAPIKeyResponse,
  APIKeyInfo
} from '../types/api';

/**
 * Provider服务类
 * 负责处理AI提供商相关的API调用
 */
export class ProviderService extends BaseService implements IValidatable<string> {
  
  /**
   * 获取所有提供商
   */
  async getProviders(): Promise<ProvidersResponse> {
    return this.get<ProvidersResponse>('/api/v1/ai/providers');
  }

  /**
   * 获取指定提供商的模型列表（仅启用的）
   */
  async getModels(provider: string): Promise<ModelsResponse> {
    return this.get<ModelsResponse>(`/api/v1/ai/${provider}/models`);
  }

  /**
   * 获取指定提供商的所有模型列表（包括禁用的）
   */
  async getAllModels(provider: string): Promise<ModelsResponse> {
    return this.get<ModelsResponse>(`/api/v1/ai/${provider}/models/all`);
  }

  /**
   * 获取模型配置
   */
  async getModelConfig(provider: string, model: string): Promise<ModelConfigResponse> {
    return this.get<ModelConfigResponse>(`/api/v1/ai/${provider}/config/${model}`);
  }

  /**
   * 启用模型
   */
  async enableModel(provider: string, model: string): Promise<BaseApiResponse> {
    return this.put<BaseApiResponse>(`/api/v1/ai/${provider}/models/${model}/enable`);
  }

  /**
   * 禁用模型
   */
  async disableModel(provider: string, model: string): Promise<BaseApiResponse> {
    return this.put<BaseApiResponse>(`/api/v1/ai/${provider}/models/${model}/disable`);
  }

  /**
   * 设置API密钥
   */
  async setAPIKey(provider: string, apiKey: string): Promise<BaseApiResponse> {
    return this.post<BaseApiResponse>(`/api/v1/ai/${provider}/api-key`, {
      api_key: apiKey
    });
  }

  /**
   * 验证API密钥
   */
  async validateAPIKey(provider: string): Promise<ValidateAPIKeyResponse> {
    return this.post<ValidateAPIKeyResponse>(`/api/v1/ai/${provider}/validate`);
  }

  /**
   * 获取API密钥状态
   */
  async getAPIKeyStatus(): Promise<BaseApiResponse<Record<string, APIKeyInfo>>> {
    return this.get<BaseApiResponse<Record<string, APIKeyInfo>>>('/api/v1/ai/api-keys/status');
  }

  /**
   * 获取明文API密钥
   */
  async getPlainAPIKey(provider: string): Promise<BaseApiResponse<{ provider: string; api_key: string }>> {
    return this.get<BaseApiResponse<{ provider: string; api_key: string }>>(`/api/v1/ai/${provider}/api-key/plain`);
  }

  /**
   * 检查提供商健康状态
   */
  async checkProviderHealth(provider: string): Promise<BaseApiResponse<{ healthy: boolean }>> {
    return this.get<BaseApiResponse<{ healthy: boolean }>>(`/api/v1/ai/${provider}/health`);
  }

  /**
   * 批量检查所有提供商健康状态
   */
  async checkAllProvidersHealth(): Promise<BaseApiResponse<Record<string, boolean>>> {
    return this.get<BaseApiResponse<Record<string, boolean>>>('/api/v1/ai/providers/health');
  }

  /**
   * 验证提供商名称
   */
  validate(provider: string): ValidationResult {
    const errors = [];
    
    if (!provider || provider.trim().length === 0) {
      errors.push({
        field: 'provider',
        message: 'Provider name is required',
        code: 'REQUIRED'
      });
    }

    if (provider && provider.length > 50) {
      errors.push({
        field: 'provider',
        message: 'Provider name must be less than 50 characters',
        code: 'MAX_LENGTH'
      });
    }

    return {
      isValid: errors.length === 0,
      errors
    };
  }

  /**
   * 验证API密钥格式
   */
  validateAPIKeyFormat(provider: string, apiKey: string): ValidationResult {
    const errors = [];
    
    if (!apiKey || apiKey.trim().length === 0) {
      errors.push({
        field: 'apiKey',
        message: 'API key is required',
        code: 'REQUIRED'
      });
      return { isValid: false, errors };
    }

    // 根据不同提供商验证API密钥格式
    switch (provider.toLowerCase()) {
      case 'openai':
        if (!apiKey.startsWith('sk-')) {
          errors.push({
            field: 'apiKey',
            message: 'OpenAI API key must start with "sk-"',
            code: 'INVALID_FORMAT'
          });
        }
        break;
      case 'anthropic':
        if (!apiKey.startsWith('sk-ant-')) {
          errors.push({
            field: 'apiKey',
            message: 'Anthropic API key must start with "sk-ant-"',
            code: 'INVALID_FORMAT'
          });
        }
        break;
      case 'google':
        // Google API keys have different formats, basic length check
        if (apiKey.length < 20) {
          errors.push({
            field: 'apiKey',
            message: 'Google API key appears to be too short',
            code: 'INVALID_LENGTH'
          });
        }
        break;
    }

    return {
      isValid: errors.length === 0,
      errors
    };
  }

  /**
   * 获取提供商显示名称
   */
  getProviderDisplayName(providerType: string): string {
    const displayNames: Record<string, string> = {
      'openai': 'OpenAI',
      'anthropic': 'Anthropic',
      'google': 'Google',
      'azure': 'Azure OpenAI',
      'huggingface': 'Hugging Face',
      'cohere': 'Cohere',
      'local': 'Local Model'
    };
    
    return displayNames[providerType.toLowerCase()] || providerType;
  }

  /**
   * 获取模型显示名称
   */
  getModelDisplayName(provider: string, modelName: string): string {
    // 可以根据需要添加特殊的模型名称映射
    const specialNames: Record<string, Record<string, string>> = {
      'openai': {
        'gpt-4': 'GPT-4',
        'gpt-3.5-turbo': 'GPT-3.5 Turbo',
        'text-davinci-003': 'Davinci'
      },
      'anthropic': {
        'claude-3-opus': 'Claude 3 Opus',
        'claude-3-sonnet': 'Claude 3 Sonnet',
        'claude-3-haiku': 'Claude 3 Haiku'
      }
    };

    return specialNames[provider]?.[modelName] || modelName;
  }
}

export const providerService = new ProviderService();
export default providerService;