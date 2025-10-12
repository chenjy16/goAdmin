import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import type { PayloadAction, Draft } from '@reduxjs/toolkit';

// 通用配置管理接口
export interface ConfigManager<T> {
  // 获取默认配置
  getDefault: () => T;
  // 加载配置
  load: () => T;
  // 保存配置
  save: (config: T) => void;
  // 验证配置
  validate?: (config: T) => boolean;
  // 迁移配置（版本升级时使用）
  migrate?: (oldConfig: any) => T;
}

// 配置状态接口
export interface ConfigState<T> {
  data: T;
  loading: boolean;
  error: string | null;
  initialized: boolean;
  isDirty: boolean; // 是否有未保存的更改
}

// 创建配置slice的工厂函数
export function createConfigSlice<T>(
  name: string,
  manager: ConfigManager<T>,
  asyncLoader?: () => Promise<Partial<T>>
) {
  // 异步加载配置
  const loadConfig = createAsyncThunk(
    `${name}/loadConfig`,
    async () => {
      let config = manager.load();
      
      // 如果有异步加载器，合并远程配置
      if (asyncLoader) {
        try {
          const remoteConfig = await asyncLoader();
          config = { ...config, ...remoteConfig };
        } catch (error) {
          console.warn(`Failed to load remote config for ${name}:`, error);
        }
      }
      
      // 验证配置
      if (manager.validate && !manager.validate(config)) {
        console.warn(`Invalid config for ${name}, using default`);
        config = manager.getDefault();
      }
      
      return config;
    }
  );

  // 异步保存配置
  const saveConfig = createAsyncThunk(
    `${name}/saveConfig`,
    async (config: T) => {
      manager.save(config);
      return config;
    }
  );

  const initialState: ConfigState<T> = {
    data: manager.load(), // 从存储中加载配置而不是使用默认配置
    loading: false,
    error: null,
    initialized: true, // 标记为已初始化，因为我们已经加载了配置
    isDirty: false,
  };

  const slice = createSlice({
    name,
    initialState,
    reducers: {
      // 更新配置（部分更新）
      updateConfig: (state, action: PayloadAction<Partial<T>>) => {
        state.data = { ...state.data, ...action.payload } as Draft<T>;
        state.isDirty = true;
      },
      
      // 设置完整配置
      setConfig: (state, action: PayloadAction<T>) => {
        state.data = action.payload as Draft<T>;
        state.isDirty = true;
      },
      
      // 重置为默认配置
      resetConfig: (state) => {
        state.data = manager.getDefault() as Draft<T>;
        state.isDirty = true;
      },
      
      // 标记为已保存
      markSaved: (state) => {
        state.isDirty = false;
      },
      
      // 清除错误
      clearError: (state) => {
        state.error = null;
      },
      
      // 设置加载状态
      setLoading: (state, action: PayloadAction<boolean>) => {
        state.loading = action.payload;
      },
    },
    extraReducers: (builder) => {
      builder
        // 加载配置
        .addCase(loadConfig.pending, (state) => {
          state.loading = true;
          state.error = null;
        })
        .addCase(loadConfig.fulfilled, (state, action) => {
          state.loading = false;
          state.data = action.payload as Draft<T>;
          state.initialized = true;
          state.isDirty = false;
          state.error = null;
        })
        .addCase(loadConfig.rejected, (state, action) => {
          state.loading = false;
          state.error = action.error.message || 'errors.loadFailed';
        })
        // 保存配置
        .addCase(saveConfig.pending, (state) => {
          state.loading = true;
          state.error = null;
        })
        .addCase(saveConfig.fulfilled, (state, action) => {
          state.loading = false;
          state.data = action.payload as Draft<T>;
          state.isDirty = false;
          state.error = null;
        })
        .addCase(saveConfig.rejected, (state, action) => {
          state.loading = false;
          state.error = action.error.message || 'errors.saveFailed';
        });
    },
  });

  return {
    slice,
    actions: {
      ...slice.actions,
      loadConfig,
      saveConfig,
    },
    reducer: slice.reducer,
    selectors: {
      selectConfig: (state: { [key: string]: ConfigState<T> }) => state[name]?.data,
      selectLoading: (state: { [key: string]: ConfigState<T> }) => state[name]?.loading,
      selectError: (state: { [key: string]: ConfigState<T> }) => state[name]?.error,
      selectInitialized: (state: { [key: string]: ConfigState<T> }) => state[name]?.initialized,
      selectIsDirty: (state: { [key: string]: ConfigState<T> }) => state[name]?.isDirty,
    },
  };
}

// 本地存储配置管理器
export function createLocalStorageManager<T>(
  key: string,
  defaultConfig: T,
  validator?: (config: any) => config is T
): ConfigManager<T> {
  return {
    getDefault: () => ({ ...defaultConfig }),
    
    load: () => {
      try {
        const stored = localStorage.getItem(key);
        if (!stored) return { ...defaultConfig };
        
        const parsed = JSON.parse(stored);
        
        // 验证配置
        if (validator && !validator(parsed)) {
          console.warn(`Invalid stored config for ${key}, using default`);
          return { ...defaultConfig };
        }
        
        // 合并默认配置，确保新字段有默认值
        return { ...defaultConfig, ...parsed };
      } catch (error) {
        console.error(`Failed to load config from localStorage for ${key}:`, error);
        return { ...defaultConfig };
      }
    },
    
    save: (config: T) => {
      try {
        localStorage.setItem(key, JSON.stringify(config));
      } catch (error) {
        console.error(`Failed to save config to localStorage for ${key}:`, error);
      }
    },
    
    validate: validator ? (config: T) => validator(config) : undefined,
  };
}

// 会话存储配置管理器
export function createSessionStorageManager<T>(
  key: string,
  defaultConfig: T,
  validator?: (config: any) => config is T
): ConfigManager<T> {
  return {
    getDefault: () => ({ ...defaultConfig }),
    
    load: () => {
      try {
        const stored = sessionStorage.getItem(key);
        if (!stored) return { ...defaultConfig };
        
        const parsed = JSON.parse(stored);
        
        if (validator && !validator(parsed)) {
          console.warn(`Invalid stored config for ${key}, using default`);
          return { ...defaultConfig };
        }
        
        return { ...defaultConfig, ...parsed };
      } catch (error) {
        console.error(`Failed to load config from sessionStorage for ${key}:`, error);
        return { ...defaultConfig };
      }
    },
    
    save: (config: T) => {
      try {
        sessionStorage.setItem(key, JSON.stringify(config));
      } catch (error) {
        console.error(`Failed to save config to sessionStorage for ${key}:`, error);
      }
    },
    
    validate: validator,
  };
}

// 内存配置管理器（用于临时配置）
export function createMemoryManager<T>(defaultConfig: T): ConfigManager<T> {
  let currentConfig = { ...defaultConfig };
  
  return {
    getDefault: () => ({ ...defaultConfig }),
    load: () => ({ ...currentConfig }),
    save: (config: T) => {
      currentConfig = { ...config };
    },
  };
}