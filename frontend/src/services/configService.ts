import { apiService } from './api';
import type { ProviderInfo, ModelInfo, MCPTool } from '../types/api';

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

class ConfigService {
  private cache: {
    providers?: ProviderInfo[];
    models: Record<string, ModelInfo[]>;
    tools?: MCPTool[];
  } = {
    models: {}
  };

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
        models[provider.name] = await this.getModels(provider.type);
      } catch (error) {
        console.warn(`Failed to load models for provider ${provider.name}:`, error);
        models[provider.name] = [];
      }
    }

    return {
      providers,
      models,
      tools
    };
  }

  /**
   * 获取提供商列表
   */
  async getProviders(): Promise<ProviderInfo[]> {
    if (this.cache.providers) {
      return this.cache.providers;
    }

    try {
      const response = await apiService.getProviders();
      this.cache.providers = response.data?.providers || [];
      return this.cache.providers;
    } catch (error) {
      console.error('Failed to fetch providers:', error);
      return [];
    }
  }

  /**
   * 获取指定提供商的模型列表
   */
  async getModels(provider: string): Promise<ModelInfo[]> {
    if (this.cache.models[provider]) {
      return this.cache.models[provider];
    }

    try {
      const response = await apiService.getModels(provider);
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
      
      this.cache.models[provider] = modelsArray;
      return this.cache.models[provider];
    } catch (error) {
      console.error(`Failed to fetch models for provider ${provider}:`, error);
      return [];
    }
  }

  /**
   * 获取MCP工具列表
   */
  async getTools(): Promise<MCPTool[]> {
    if (this.cache.tools) {
      return this.cache.tools;
    }

    try {
      const response = await apiService.getMCPTools();
      this.cache.tools = response.data?.tools || [];
      return this.cache.tools;
    } catch (error) {
      console.error('Failed to fetch MCP tools:', error);
      return [];
    }
  }

  /**
   * 清除缓存
   */
  clearCache(): void {
    this.cache = { models: {} };
  }

  /**
   * 获取默认配置
   */
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

  /**
   * 验证配置是否有效
   */
  validateConfig(config: ConfigState, configData: ConfigData): {
    isValid: boolean;
    errors: string[];
  } {
    const errors: string[] = [];

    // 验证提供商
    if (!config.selectedProvider) {
      errors.push('请选择AI提供商');
    } else {
      const provider = configData.providers.find(p => p.name === config.selectedProvider);
      if (!provider) {
        errors.push('选择的提供商不存在');
      } else if (provider.healthy !== true) {
        errors.push('选择的提供商当前不可用');
      }
    }

    // 验证模型
    if (!config.selectedModel) {
      errors.push('请选择模型');
    } else if (config.selectedProvider) {
      const models = configData.models[config.selectedProvider] || [];
      const model = models.find(m => m.name === config.selectedModel);
      if (!model) {
        errors.push('选择的模型不存在');
      }
    }

    // 验证参数范围
    if (config.temperature < 0 || config.temperature > 2) {
      errors.push('Temperature必须在0-2之间');
    }

    if (config.maxTokens < 1 || config.maxTokens > 32768) {
      errors.push('Max Tokens必须在1-32768之间');
    }

    if (config.topP < 0 || config.topP > 1) {
      errors.push('Top-p必须在0-1之间');
    }

    return {
      isValid: errors.length === 0,
      errors
    };
  }

  /**
   * 保存配置到本地存储
   */
  saveConfig(config: ConfigState): void {
    try {
      localStorage.setItem('assistant_config', JSON.stringify(config));
    } catch (error) {
      console.error('Failed to save config:', error);
    }
  }

  /**
   * 从本地存储加载配置
   */
  loadConfig(): ConfigState {
    try {
      const saved = localStorage.getItem('assistant_config');
      if (saved) {
        const config = JSON.parse(saved);
        return { ...this.getDefaultConfig(), ...config };
      }
    } catch (error) {
      console.error('Failed to load config:', error);
    }
    return this.getDefaultConfig();
  }
}

export const configService = new ConfigService();
export default configService;