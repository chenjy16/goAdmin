import { getProviderService, getMCPService } from './ServiceFactory';
import type { ProviderInfo, ModelInfo, MCPTool } from '../types/api';
import type { IConfigurable, IValidatable, ValidationResult, ValidationError } from '../types/base';

export interface ConfigData {
  providers: ProviderInfo[];
  models: Record<string, ModelInfo[]>;
  tools: MCPTool[];
}

export interface ConfigState {
  selectedProvider: string;
  selectedModel: string;
  selectedTools: string[];
  temperature: number;
  maxTokens: number;
  topP: number;
}

/**
 * 配置服务类
 * 负责管理应用程序的配置数据，包括提供商、模型和工具
 */
class ConfigService implements IConfigurable, IValidatable<ConfigState> {
  private cache: {
    providers?: ProviderInfo[];
    models: Record<string, ModelInfo[]>;
    tools?: MCPTool[];
  } = {
    models: {}
  };

  private readonly CACHE_TTL = 5 * 60 * 1000; // 5分钟缓存
  private cacheTimestamps: Record<string, number> = {};

  /**
   * 获取所有配置数据
   */
  async getAllConfigData(): Promise<ConfigData> {
    const [providers, tools] = await Promise.all([
      this.getProviders(),
      this.getTools()
    ]);

    // 获取所有提供商的模型
    const models: Record<string, ModelInfo[]> = {};
    for (const provider of providers) {
      try {
        models[provider.type] = await this.getModels(provider.type);
      } catch (error) {
        console.warn(`Failed to load models for provider ${provider.name}:`, error);
        models[provider.type] = [];
      }
    }

    const result = {
      providers,
      models,
      tools
    };
    
    return result;
  }

  /**
   * 获取提供商列表
   */
  async getProviders(): Promise<ProviderInfo[]> {
    const cacheKey = 'providers';
    
    if (this.isCacheValid(cacheKey) && this.cache.providers) {
      return this.cache.providers;
    }

    try {
      const providerService = getProviderService();
      const response = await providerService.getProviders();
      this.cache.providers = response.data?.providers || [];
      this.setCacheTimestamp(cacheKey);
      return this.cache.providers;
    } catch (error) {
      console.error('Failed to fetch providers:', error);
      return [];
    }
  }

  /**
   * 获取指定提供商的模型列表（仅启用的）
   */
  async getModels(provider: string): Promise<ModelInfo[]> {
    const cacheKey = `models_${provider}`;
    
    if (this.isCacheValid(cacheKey) && this.cache.models[provider]) {
      return this.cache.models[provider];
    }

    try {
      const providerService = getProviderService();
      const response = await providerService.getModels(provider);
      
      // 将对象格式的模型数据转换为数组格式
      const modelsObject = response.data?.models || {};
      
      // 确保每个模型都有必要的字段，并添加缺失的字段
      const modelsArray = Object.values(modelsObject).map((model: any) => {
        const processedModel = {
          name: model.name || model.Name || '',
          display_name: model.display_name || model.DisplayName || model.name || model.Name || '',
          max_tokens: model.max_tokens || model.MaxTokens || 2048,
          temperature: model.temperature || model.Temperature || 0.7,
          top_p: model.top_p || model.TopP || 1.0,
          enabled: model.enabled !== undefined ? model.enabled : (model.Enabled !== undefined ? model.Enabled : true),
          description: model.description || model.Description || ''
        };
        return processedModel;
      });
      
      this.cache.models[provider] = modelsArray;
      this.setCacheTimestamp(cacheKey);
      return this.cache.models[provider];
    } catch (error) {
      console.error(`configService - getModels: Failed to fetch models for provider ${provider}:`, error);
      return [];
    }
  }

  /**
   * 获取指定提供商的所有模型列表（包括禁用的，用于模型管理）
   */
  async getAllModels(provider: string): Promise<ModelInfo[]> {
    const cacheKey = `all_models_${provider}`;
    
    if (this.isCacheValid(cacheKey) && this.cache.models[`all_${provider}`]) {
      return this.cache.models[`all_${provider}`];
    }

    try {
      const providerService = getProviderService();
      const response = await providerService.getAllModels(provider);
      
      // 将对象格式的模型数据转换为数组格式
      const modelsObject = response.data?.models || {};
      
      // 确保每个模型都有必要的字段，并添加缺失的字段
      const modelsArray = Object.values(modelsObject).map((model: any) => ({
        name: model.name || model.Name || '',
        display_name: model.display_name || model.DisplayName || model.name || model.Name || '',
        max_tokens: model.max_tokens || model.MaxTokens || 2048,
        temperature: model.temperature || model.Temperature || 0.7,
        top_p: model.top_p || model.TopP || 1.0,
        enabled: model.enabled !== undefined ? model.enabled : (model.Enabled !== undefined ? model.Enabled : true),
        description: model.description || model.Description || ''
      }));
      
      // 缓存数据
      this.cache.models[`all_${provider}`] = modelsArray;
      this.cacheTimestamps[cacheKey] = Date.now();
      
      return modelsArray;
    } catch (error) {
      console.error(`configService - getAllModels: Failed to fetch all models for provider ${provider}:`, error);
      return [];
    }
  }

  /**
   * 获取MCP工具列表
   */
  async getTools(): Promise<MCPTool[]> {
    const cacheKey = 'tools';
    
    if (this.isCacheValid(cacheKey) && this.cache.tools) {
      return this.cache.tools;
    }

    try {
      const mcpService = getMCPService();
      const response = await mcpService.getTools();
      this.cache.tools = response.data?.tools || [];
      this.setCacheTimestamp(cacheKey);
      return this.cache.tools;
    } catch (error) {
      console.error('Failed to fetch MCP tools:', error);
      return [];
    }
  }

  // ICache接口实现
  clearCache(): void {
    this.cache = { models: {} };
    this.cacheTimestamps = {};
  }

  private isCacheValid(key: string): boolean {
    const timestamp = this.cacheTimestamps[key];
    if (!timestamp) return false;
    return Date.now() - timestamp < this.CACHE_TTL;
  }

  private setCacheTimestamp(key: string): void {
    this.cacheTimestamps[key] = Date.now();
  }

  // IConfigurable接口实现
  getDefaultConfig(): ConfigState {
    return {
      selectedProvider: '',
      selectedModel: '',
      selectedTools: [],
      temperature: 0.7,
      maxTokens: 2048,
      topP: 1.0
    };
  }

  // IValidatable接口实现
  validate(config: ConfigState): ValidationResult {
    return this.validateConfig(config);
  }

  validateConfig(config: ConfigState, configData?: ConfigData): ValidationResult {
    const errors: ValidationError[] = [];

    if (!config.selectedProvider) {
      errors.push({ field: 'selectedProvider', message: '请选择AI提供商' });
    }

    if (!config.selectedModel) {
      errors.push({ field: 'selectedModel', message: '请选择模型' });
    }

    if (config.temperature < 0 || config.temperature > 2) {
      errors.push({ field: 'temperature', message: '温度值必须在0-2之间' });
    }

    if (config.maxTokens < 1 || config.maxTokens > 32000) {
      errors.push({ field: 'maxTokens', message: '最大令牌数必须在1-32000之间' });
    }

    if (config.topP < 0 || config.topP > 1) {
      errors.push({ field: 'topP', message: 'Top P值必须在0-1之间' });
    }

    // 如果提供了配置数据，进行更详细的验证
    if (configData) {
      const provider = configData.providers.find(p => p.type === config.selectedProvider);
      if (!provider) {
        errors.push({ field: 'selectedProvider', message: '选择的提供商不存在' });
      }

      const models = configData.models[config.selectedProvider] || [];
      const model = models.find(m => m.name === config.selectedModel);
      if (!model) {
        errors.push({ field: 'selectedModel', message: '选择的模型不存在' });
      } else if (!model.enabled) {
        errors.push({ field: 'selectedModel', message: '选择的模型未启用' });
      }

      // 验证选择的工具是否存在
      const availableTools = configData.tools.map(t => t.name);
      const invalidTools = config.selectedTools.filter(tool => !availableTools.includes(tool));
      if (invalidTools.length > 0) {
        errors.push({ field: 'selectedTools', message: `以下工具不存在: ${invalidTools.join(', ')}` });
      }
    }

    return {
      isValid: errors.length === 0,
      errors
    };
  }

  // IConfigurable接口实现
  configure(config: Record<string, any>): void {
    // 可以在这里设置服务的配置参数
    if (config.cacheTTL) {
      // 动态设置缓存TTL等
    }
  }

  getConfig(): Record<string, any> {
    return {
      cacheTTL: this.CACHE_TTL,
      cacheKeys: Object.keys(this.cacheTimestamps)
    };
  }

  saveConfig(config: ConfigState): void {
    try {
      localStorage.setItem('app_config', JSON.stringify(config));
    } catch (error) {
      console.error('Failed to save config:', error);
    }
  }

  loadConfig(): ConfigState {
    try {
      const saved = localStorage.getItem('app_config');
      if (saved) {
        const config = JSON.parse(saved);
        // 合并默认配置以确保所有字段都存在
        return { ...this.getDefaultConfig(), ...config };
      }
    } catch (error) {
      console.error('Failed to load config:', error);
    }
    return this.getDefaultConfig();
  }

  /**
   * 重置配置到默认值
   */
  resetConfig(): ConfigState {
    const defaultConfig = this.getDefaultConfig();
    this.saveConfig(defaultConfig);
    return defaultConfig;
  }

  /**
   * 获取配置的摘要信息
   */
  getConfigSummary(config: ConfigState): string {
    return `Provider: ${config.selectedProvider}, Model: ${config.selectedModel}, Tools: ${config.selectedTools.length}`;
  }
}

export const configService = new ConfigService();
export default configService;