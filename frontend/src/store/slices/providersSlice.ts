import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';
import apiService from '../../services/api';
import type { ProviderInfo, ModelInfo, APIKeyInfo } from '../../types/api';

// 提供商状态接口
export interface ProvidersState {
  providers: ProviderInfo[];
  models: Record<string, ModelInfo[]>; // provider -> models
  apiKeyStatus: Record<string, APIKeyInfo>; // provider -> APIKeyInfo
  selectedProvider: string | null;
  isLoading: boolean;
  error: string | null;
}

// 异步thunks
export const fetchProviders = createAsyncThunk(
  'providers/fetchProviders',
  async () => {
    const response = await apiService.getProviders();
    return response.data.providers;
  }
);

// 自定义异步thunks
export const fetchModels = createAsyncThunk(
  'providers/fetchModels',
  async (provider: string) => {
    const response = await apiService.getModels(provider);
    const modelsArray = Object.values(response.data.models);
    return { provider, models: modelsArray };
  }
);

export const fetchAllModels = createAsyncThunk(
  'providers/fetchAllModels',
  async (provider: string) => {
    const response = await apiService.getAllModels(provider);
    const modelsArray = Object.values(response.data.models);
    return { provider, models: modelsArray };
  }
);

export const setAPIKey = createAsyncThunk(
  'providers/setAPIKey',
  async (params: { provider: string; apiKey: string }) => {
    const { provider, apiKey } = params;
    await apiService.setAPIKey(provider, apiKey);
    return { provider, apiKey };
  }
);

export const validateAPIKey = createAsyncThunk(
  'providers/validateAPIKey',
  async (provider: string) => {
    const response = await apiService.validateAPIKey(provider);
    return { provider, valid: response.valid };
  }
);

export const toggleModel = createAsyncThunk(
  'providers/toggleModel',
  async (params: { provider: string; model: string; enabled: boolean }) => {
    const { provider, model, enabled } = params;
    if (enabled) {
      await apiService.enableModel(provider, model);
    } else {
      await apiService.disableModel(provider, model);
    }
    return { provider, model, enabled };
  }
);

export const fetchAPIKeyStatus = createAsyncThunk<Record<string, APIKeyInfo>>(
  'providers/fetchAPIKeyStatus',
  async () => {
    const response = await apiService.getAPIKeyStatus();
    return response.data || {};
  }
);

export const fetchPlainAPIKey = createAsyncThunk(
  'providers/fetchPlainAPIKey',
  async (provider: string) => {
    const response = await apiService.getPlainAPIKey(provider);
    return { provider, apiKey: response.data.api_key };
  }
);

const initialState: ProvidersState = {
  providers: [],
  models: {},
  apiKeyStatus: {},
  selectedProvider: null,
  isLoading: false,
  error: null,
};

const providersSlice = createSlice({
  name: 'providers',
  initialState,
  reducers: {
    setSelectedProvider: (state, action: PayloadAction<string>) => {
      state.selectedProvider = action.payload;
    },
    clearError: (state) => {
      state.error = null;
    },
    updateProviderHealth: (state, action: PayloadAction<{
      provider: string;
      healthy: boolean;
    }>) => {
      const { provider, healthy } = action.payload;
      const providerInfo = (state.providers || []).find(p => p.name === provider);
      if (providerInfo) {
        providerInfo.healthy = healthy;
      }
    },
  },
  extraReducers: (builder) => {
    builder
      // fetchProviders
      .addCase(fetchProviders.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchProviders.fulfilled, (state, action) => {
        state.isLoading = false;
        state.providers = action.payload;
        if (!state.selectedProvider && (action.payload || []).length > 0 && (action.payload || [])[0]) {
          state.selectedProvider = (action.payload || [])[0].name;
        }
        
        // 检查每个提供商的健康状态
        (action.payload || []).forEach((provider: ProviderInfo) => {
          const providerInfo = (state.providers || []).find((p: ProviderInfo) => p.name === provider.name);
          if (providerInfo) {
            providerInfo.healthy = provider.healthy;
          }
        });
      })
      .addCase(fetchProviders.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || 'errors.loadFailed';
      })
      // fetchModels
      .addCase(fetchModels.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchModels.fulfilled, (state, action) => {
        state.isLoading = false;
        const { provider, models } = action.payload;
        state.models[provider] = models;
      })
      .addCase(fetchModels.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || 'errors.loadFailed';
      })
      // fetchAllModels
      .addCase(fetchAllModels.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchAllModels.fulfilled, (state, action) => {
        state.isLoading = false;
        const { provider, models } = action.payload;
        state.models[provider] = models;
      })
      .addCase(fetchAllModels.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || 'errors.loadFailed';
      })
      // setAPIKey
      .addCase(setAPIKey.fulfilled, () => {
        // API密钥设置成功，可以更新相关状态
      })
      .addCase(setAPIKey.rejected, (state, action) => {
        state.error = action.error.message || 'errors.saveFailed';
      })
      // validateAPIKey
      .addCase(validateAPIKey.fulfilled, (state, action) => {
        const { provider, valid } = action.payload;
        const providerInfo = (state.providers || []).find((p: ProviderInfo) => p.name === provider);
        if (providerInfo) {
          providerInfo.healthy = valid;
        }
      })
      .addCase(validateAPIKey.rejected, (state, action) => {
        state.error = action.error.message || 'errors.saveFailed';
      })
      // toggleModel
      .addCase(toggleModel.fulfilled, (state, action) => {
        const { provider, model, enabled } = action.payload;
        const models = state.models[provider];
        if (models) {
          const modelInfo = (models || []).find((m: ModelInfo) => m.name === model);
          if (modelInfo) {
            modelInfo.enabled = enabled;
          }
        }
      })
      .addCase(toggleModel.rejected, (state, action) => {
        state.error = action.error.message || 'errors.updateFailed';
      })
      // fetchAPIKeyStatus
      .addCase(fetchAPIKeyStatus.fulfilled, (state, action) => {
        state.apiKeyStatus = action.payload || {};
      })
      .addCase(fetchAPIKeyStatus.rejected, (state, action) => {
        state.error = action.error.message || 'errors.loadFailed';
      });
  },
});

export const {
  setSelectedProvider,
  clearError,
  updateProviderHealth,
} = providersSlice.actions;

export default providersSlice.reducer;