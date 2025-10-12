import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import type { PayloadAction, CaseReducer } from '@reduxjs/toolkit';
import type { Draft, WritableDraft } from 'immer';

// 通用异步状态接口
export interface EnhancedAsyncState<T = any> {
  data: T;
  loading: boolean;
  error: string | null;
  initialized: boolean;
  lastUpdated: string | null;
}

// 异步操作配置
export interface AsyncOperationConfig<TData, TParams = void, TResult = TData> {
  name: string;
  operation: (params: TParams) => Promise<TResult>;
  transform?: (result: TResult) => TData;
  onSuccess?: (state: WritableDraft<EnhancedAsyncState<TData>>, result: TResult) => void;
  onError?: (state: WritableDraft<EnhancedAsyncState<TData>>, error: string) => void;
}

// 创建异步操作thunk的工厂函数
export function createAsyncOperation<TData, TParams = void, TResult = TData>(
  config: AsyncOperationConfig<TData, TParams, TResult>
) {
  const { name, operation, transform, onSuccess, onError } = config;

  return createAsyncThunk(
    name,
    async (params: TParams, { rejectWithValue }) => {
      try {
        const result = await operation(params);
        return transform ? transform(result) : (result as unknown as TData);
      } catch (error) {
        const message = error instanceof Error ? error.message : 'errors.operationFailed';
        return rejectWithValue(message);
      }
    }
  );
}

// 批量异步操作配置
export interface BatchAsyncConfig<T> {
  name: string;
  initialData: T;
  operations: Record<string, AsyncOperationConfig<T, any, any>>;
  extraReducers?: Record<string, (state: EnhancedAsyncState<T>, action: PayloadAction<any>) => void>;
}

// 创建批量异步操作slice的工厂函数
export function createBatchAsyncSlice<T>(config: BatchAsyncConfig<T>) {
  const { name, initialData, operations, extraReducers = {} } = config;

  // 创建所有异步thunks
  const asyncThunks: Record<string, any> = {};
  Object.entries(operations).forEach(([key, operationConfig]) => {
    asyncThunks[key] = createAsyncOperation({
      ...operationConfig,
      name: `${name}/${key}`,
    });
  });

  // 初始状态
  const initialState: EnhancedAsyncState<T> = {
    data: initialData,
    loading: false,
    error: null,
    initialized: false,
    lastUpdated: null,
  };

  // 创建slice
  const slice = createSlice({
    name,
    initialState,
    reducers: {
      // 设置数据
      setData: (state, action: PayloadAction<T>) => {
        state.data = action.payload as Draft<T>;
        state.initialized = true;
        state.lastUpdated = new Date().toISOString();
      },
      
      // 更新数据（部分更新）
      updateData: (state, action: PayloadAction<Partial<T>>) => {
        Object.assign(state.data as any, action.payload);
        state.lastUpdated = new Date().toISOString();
      },
      
      // 设置加载状态
      setLoading: (state, action: PayloadAction<boolean>) => {
        state.loading = action.payload;
      },
      
      // 设置错误
      setError: (state, action: PayloadAction<string | null>) => {
        state.error = action.payload;
      },
      
      // 清除错误
      clearError: (state) => {
        state.error = null;
      },
      
      // 重置状态
      reset: (state) => {
        state.data = initialData as Draft<T>;
        state.loading = false;
        state.error = null;
        state.initialized = false;
        state.lastUpdated = null;
      },
      
      // 额外的reducers
      ...extraReducers,
    },
    extraReducers: (builder) => {
      // 为每个异步操作添加处理器
      Object.entries(asyncThunks).forEach(([key, thunk]) => {
        const operationConfig = operations[key];
        
        builder
          .addCase(thunk.pending, (state) => {
            state.loading = true;
            state.error = null;
          })
          .addCase(thunk.fulfilled, (state, action) => {
            state.loading = false;
            state.data = action.payload as Draft<T>;
            state.initialized = true;
            state.lastUpdated = new Date().toISOString();
            state.error = null;
            
            // 执行自定义成功回调
            if (operationConfig.onSuccess) {
              operationConfig.onSuccess(state, action.payload);
            }
          })
          .addCase(thunk.rejected, (state, action) => {
            state.loading = false;
            state.error = action.payload as string || 'errors.operationFailed';
            
            // 执行自定义错误回调
            if (operationConfig.onError) {
              operationConfig.onError(state, state.error);
            }
          });
      });
    },
  });

  return {
    slice,
    actions: {
      ...slice.actions,
      ...asyncThunks,
    },
    reducer: slice.reducer,
    selectors: {
      selectData: (state: { [key: string]: EnhancedAsyncState<T> }) => state[name]?.data,
      selectLoading: (state: { [key: string]: EnhancedAsyncState<T> }) => state[name]?.loading,
      selectError: (state: { [key: string]: EnhancedAsyncState<T> }) => state[name]?.error,
      selectInitialized: (state: { [key: string]: EnhancedAsyncState<T> }) => state[name]?.initialized,
      selectLastUpdated: (state: { [key: string]: EnhancedAsyncState<T> }) => state[name]?.lastUpdated,
    },
  };
}

// 创建API资源slice的工厂函数
export interface ApiResourceConfig<T> {
  name: string;
  initialData: T;
  apiService: {
    fetch?: () => Promise<T>;
    create?: (data: Partial<T>) => Promise<T>;
    update?: (id: string | number, data: Partial<T>) => Promise<T>;
    delete?: (id: string | number) => Promise<void>;
  };
}

export function createApiResourceSlice<T>(config: ApiResourceConfig<T>) {
  const { name, initialData, apiService } = config;

  const operations: Record<string, AsyncOperationConfig<T, any, any>> = {};

  // 添加获取操作
  if (apiService.fetch) {
    operations.fetch = {
      name: 'fetch',
      operation: apiService.fetch,
    };
  }

  // 添加创建操作
  if (apiService.create) {
    operations.create = {
      name: 'create',
      operation: apiService.create,
    };
  }

  // 添加更新操作
  if (apiService.update) {
    operations.update = {
      name: 'update',
      operation: ({ id, data }: { id: string | number; data: Partial<T> }) => 
        apiService.update!(id, data),
    };
  }

  // 添加删除操作
  if (apiService.delete) {
    operations.delete = {
      name: 'delete',
      operation: apiService.delete,
      transform: () => initialData, // 删除后返回初始数据
    };
  }

  return createBatchAsyncSlice({
    name,
    initialData,
    operations,
  });
}

// 导出类型
export type { EnhancedAsyncState, AsyncOperationConfig, BatchAsyncConfig, ApiResourceConfig };