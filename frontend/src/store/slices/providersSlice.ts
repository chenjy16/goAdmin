import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';
import apiService from '../../services/api';
import type { ProvidersState } from '../../types/store';


// 异步thunks
export const fetchProviders = createAsyncThunk(
  'providers/fetchProviders',
  async () => {
    const response = await apiService.getProviders();
    return response.providers;
  }
);

export const fetchModels = createAsyncThunk(
  'providers/fetchModels',
  async (provider: string) => {
    const response = await apiService.getModels(provider);
    return { provider, models: response.models };
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

const initialState: ProvidersState = {
  providers: [],
  models: {},
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
      health: boolean;
    }>) => {
      const { provider, health } = action.payload;
      const providerInfo = (state.providers || []).find(p => p.name === provider);
      if (providerInfo) {
        providerInfo.health = health;
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
        if (!state.selectedProvider && (action.payload || []).length > 0) {
          state.selectedProvider = (action.payload || [])[0].name;
        }
      })
      .addCase(fetchProviders.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || '获取提供商列表失败';
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
        state.error = action.error.message || '获取模型列表失败';
      })
      // setAPIKey
      .addCase(setAPIKey.fulfilled, () => {
        // API密钥设置成功，可以更新相关状态
      })
      .addCase(setAPIKey.rejected, (state, action) => {
        state.error = action.error.message || '设置API密钥失败';
      })
      // validateAPIKey
      .addCase(validateAPIKey.fulfilled, (state, action) => {
        const { provider, valid } = action.payload;
        const providerInfo = (state.providers || []).find(p => p.name === provider);
        if (providerInfo) {
          providerInfo.health = valid;
        }
      })
      .addCase(validateAPIKey.rejected, (state, action) => {
        state.error = action.error.message || '验证API密钥失败';
      })
      // toggleModel
      .addCase(toggleModel.fulfilled, (state, action) => {
        const { provider, model, enabled } = action.payload;
        const models = state.models[provider];
        if (models) {
          const modelInfo = (models || []).find(m => m.id === model);
          if (modelInfo) {
            modelInfo.enabled = enabled;
          }
        }
      })
      .addCase(toggleModel.rejected, (state, action) => {
        state.error = action.error.message || '切换模型状态失败';
      });
  },
});

export const {
  setSelectedProvider,
  clearError,
  updateProviderHealth,
} = providersSlice.actions;

export default providersSlice.reducer;