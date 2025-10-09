import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';
import { configService, type ConfigData, type ConfigState } from '../../services/configService';
import type { ProviderInfo, ModelInfo, MCPTool } from '../../types/api';

// 配置状态接口
export interface ConfigSliceState {
  // 配置数据
  data: ConfigData;
  // 当前配置
  config: ConfigState;
  // 加载状态
  loading: boolean;
  // 错误信息
  error: string | null;
  // 数据是否已初始化
  initialized: boolean;
}

// 异步thunks
export const loadConfigData = createAsyncThunk(
  'config/loadConfigData',
  async () => {
    const data = await configService.getAllConfigData();
    return data;
  }
);

export const loadProviders = createAsyncThunk(
  'config/loadProviders',
  async () => {
    const providers = await configService.getProviders();
    return providers;
  }
);

export const loadModels = createAsyncThunk(
  'config/loadModels',
  async (provider: string) => {
    const models = await configService.getModels(provider);
    return { provider, models };
  }
);

export const loadTools = createAsyncThunk(
  'config/loadTools',
  async () => {
    const tools = await configService.getTools();
    return tools;
  }
);

const initialState: ConfigSliceState = {
  data: {
    providers: [],
    models: {},
    tools: []
  },
  config: configService.getDefaultConfig(),
  loading: false,
  error: null,
  initialized: false,
};

const configSlice = createSlice({
  name: 'config',
  initialState,
  reducers: {
    // 更新配置
    updateConfig: (state, action: PayloadAction<Partial<ConfigState>>) => {
      state.config = { ...state.config, ...action.payload };
      configService.saveConfig(state.config);
    },
    
    // 设置完整配置
    setConfig: (state, action: PayloadAction<ConfigState>) => {
      state.config = action.payload;
      configService.saveConfig(state.config);
    },
    
    // 重置配置为默认值
    resetConfig: (state) => {
      state.config = configService.getDefaultConfig();
      configService.saveConfig(state.config);
    },
    
    // 加载保存的配置
    loadSavedConfig: (state) => {
      state.config = configService.loadConfig();
    },
    
    // 设置提供商
    setProvider: (state, action: PayloadAction<string>) => {
      const provider = action.payload;
      const providerModels = state.data.models[provider] || [];
      
      state.config.selectedProvider = provider;
      // 自动选择第一个可用模型
      state.config.selectedModel = providerModels.length > 0 && providerModels[0] ? providerModels[0].name : '';
      
      configService.saveConfig(state.config);
    },
    
    // 设置模型
    setModel: (state, action: PayloadAction<string>) => {
      state.config.selectedModel = action.payload;
      configService.saveConfig(state.config);
    },
    
    // 设置工具
    setTools: (state, action: PayloadAction<string[]>) => {
      state.config.selectedTools = action.payload;
      configService.saveConfig(state.config);
    },
    
    // 设置参数
    setTemperature: (state, action: PayloadAction<number>) => {
      state.config.temperature = action.payload;
      configService.saveConfig(state.config);
    },
    
    setMaxTokens: (state, action: PayloadAction<number>) => {
      state.config.maxTokens = action.payload;
      configService.saveConfig(state.config);
    },
    
    setTopP: (state, action: PayloadAction<number>) => {
      state.config.topP = action.payload;
      configService.saveConfig(state.config);
    },
    
    // 清除错误
    clearError: (state) => {
      state.error = null;
    },
    
    // 清除缓存
    clearCache: (state) => {
      configService.clearCache();
      state.data = {
        providers: [],
        models: {},
        tools: []
      };
      state.initialized = false;
    },
  },
  extraReducers: (builder) => {
    builder
      // 加载配置数据
      .addCase(loadConfigData.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(loadConfigData.fulfilled, (state, action) => {
        state.loading = false;
        state.data = action.payload;
        state.initialized = true;
        
        // 如果当前没有选择提供商，自动选择第一个健康的提供商
        if (!state.config.selectedProvider && action.payload.providers.length > 0) {
          const healthyProvider = action.payload.providers.find(p => p.healthy === true);
          if (healthyProvider) {
            state.config.selectedProvider = healthyProvider.name;
            
            // 自动选择第一个模型
            const providerModels = action.payload.models[healthyProvider.name] || [];
            if (providerModels.length > 0 && providerModels[0]) {
              state.config.selectedModel = providerModels[0].name;
            }
            
            configService.saveConfig(state.config);
          }
        }
      })
      .addCase(loadConfigData.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || '加载配置数据失败';
      })
      
      // 加载提供商
      .addCase(loadProviders.fulfilled, (state, action) => {
        state.data.providers = action.payload;
      })
      .addCase(loadProviders.rejected, (state, action) => {
        state.error = action.error.message || '加载提供商失败';
      })
      
      // 加载模型
      .addCase(loadModels.fulfilled, (state, action) => {
        const { provider, models } = action.payload;
        state.data.models[provider] = models;
      })
      .addCase(loadModels.rejected, (state, action) => {
        state.error = action.error.message || '加载模型失败';
      })
      
      // 加载工具
      .addCase(loadTools.fulfilled, (state, action) => {
        state.data.tools = action.payload;
      })
      .addCase(loadTools.rejected, (state, action) => {
        state.error = action.error.message || '加载工具失败';
      });
  },
});

export const {
  updateConfig,
  setConfig,
  resetConfig,
  loadSavedConfig,
  setProvider,
  setModel,
  setTools,
  setTemperature,
  setMaxTokens,
  setTopP,
  clearError,
  clearCache,
} = configSlice.actions;

// 选择器
export const selectConfigData = (state: { config: ConfigSliceState }) => state.config.data;
export const selectConfig = (state: { config: ConfigSliceState }) => state.config.config;
export const selectConfigLoading = (state: { config: ConfigSliceState }) => state.config.loading;
export const selectConfigError = (state: { config: ConfigSliceState }) => state.config.error;
export const selectConfigInitialized = (state: { config: ConfigSliceState }) => state.config.initialized;

// 复合选择器
export const selectCurrentProvider = (state: { config: ConfigSliceState }) => {
  const { data, config } = state.config;
  return data.providers.find(p => p.name === config.selectedProvider);
};

export const selectCurrentModel = (state: { config: ConfigSliceState }) => {
  const { data, config } = state.config;
  if (!config.selectedProvider || !config.selectedModel) return undefined;
  const providerModels = data.models[config.selectedProvider] || [];
  return providerModels.find(m => m.name === config.selectedModel);
};

export const selectAvailableModels = (state: { config: ConfigSliceState }) => {
  const { data, config } = state.config;
  return data.models[config.selectedProvider] || [];
};

export const selectSelectedTools = (state: { config: ConfigSliceState }) => {
  const { data, config } = state.config;
  return data.tools.filter(tool => config.selectedTools.includes(tool.name));
};

export const selectConfigValidation = (state: { config: ConfigSliceState }) => {
  const { data, config } = state.config;
  return configService.validateConfig(config, data);
};

export default configSlice.reducer;