import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';
import type { Draft } from 'immer';

// 通用异步状态接口
export interface AsyncState<T = any> {
  data: T;
  loading: boolean;
  error: string | null;
  initialized: boolean;
}

// 异步操作配置
export interface AsyncSliceConfig<T, P = void> {
  name: string;
  initialData: T;
  asyncThunk: (params: P) => Promise<T>;
  extraReducers?: Record<string, (state: AsyncState<T>, action: PayloadAction<any>) => void>;
}

// 创建异步slice的工厂函数
export function createAsyncSlice<T, P = void>(config: AsyncSliceConfig<T, P>) {
  const { name, initialData, asyncThunk, extraReducers = {} } = config;

  // 创建异步thunk
  const fetchData = createAsyncThunk(
    `${name}/fetchData`,
    async (params: P) => {
      return await asyncThunk(params);
    }
  );

  // 初始状态
  const initialState: AsyncState<T> = {
    data: initialData,
    loading: false,
    error: null,
    initialized: false,
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
      },
      
      // 设置加载状态
      setLoading: (state, action: PayloadAction<boolean>) => {
        state.loading = action.payload;
      },
      
      // 设置错误
      setError: (state, action: PayloadAction<string | null>) => {
        state.error = action.payload;
      },
      
      // 重置状态
      reset: (state) => {
        state.data = initialData as Draft<T>;
        state.loading = false;
        state.error = null;
        state.initialized = false;
      },
      
      // 清除错误
      clearError: (state) => {
        state.error = null;
      },
      
      // 额外的reducers
      ...extraReducers,
    },
    extraReducers: (builder) => {
      builder
        .addCase(fetchData.pending, (state) => {
          state.loading = true;
          state.error = null;
        })
        .addCase(fetchData.fulfilled, (state, action) => {
          state.loading = false;
          state.data = action.payload as Draft<T>;
          state.initialized = true;
          state.error = null;
        })
        .addCase(fetchData.rejected, (state, action) => {
          state.loading = false;
          state.error = action.error.message || '操作失败';
        });
    },
  });

  return {
    slice,
    actions: {
      ...slice.actions,
      fetchData,
    },
    reducer: slice.reducer,
    selectors: {
      selectData: (state: { [key: string]: AsyncState<T> }) => state[name]?.data,
      selectLoading: (state: { [key: string]: AsyncState<T> }) => state[name]?.loading,
      selectError: (state: { [key: string]: AsyncState<T> }) => state[name]?.error,
      selectInitialized: (state: { [key: string]: AsyncState<T> }) => state[name]?.initialized,
    },
  };
}

// 创建CRUD slice的工厂函数
export interface CrudSliceConfig<T> {
  name: string;
  initialItems: T[];
  idField: keyof T;
}

export function createCrudSlice<T>(config: CrudSliceConfig<T>) {
  const { name, initialItems, idField } = config;

  const initialState: AsyncState<T[]> = {
    data: initialItems,
    loading: false,
    error: null,
    initialized: false,
  };

  const slice = createSlice({
    name,
    initialState,
    reducers: {
      // 添加项目
      addItem: (state, action: PayloadAction<T>) => {
        (state.data as Draft<T>[]).push(action.payload as Draft<T>);
      },
      
      // 更新项目
      updateItem: (state, action: PayloadAction<T>) => {
        const index = state.data.findIndex(item => (item as any)[idField] === (action.payload as any)[idField]);
        if (index !== -1) {
          (state.data as Draft<T>[])[index] = action.payload as Draft<T>;
        }
      },
      
      // 删除项目
      removeItem: (state, action: PayloadAction<T[typeof idField]>) => {
        state.data = state.data.filter(item => (item as any)[idField] !== action.payload) as Draft<T[]>;
      },
      
      // 设置所有项目
      setItems: (state, action: PayloadAction<T[]>) => {
        state.data = action.payload as Draft<T[]>;
        state.initialized = true;
      },
      
      // 清空项目
      clearItems: (state) => {
        state.data = [];
      },
      
      // 设置加载状态
      setLoading: (state, action: PayloadAction<boolean>) => {
        state.loading = action.payload;
      },
      
      // 设置错误
      setError: (state, action: PayloadAction<string | null>) => {
        state.error = action.payload;
      },
      
      // 重置状态
      reset: (state) => {
        state.data = initialItems as Draft<T[]>;
        state.loading = false;
        state.error = null;
        state.initialized = false;
      },
    },
  });

  return {
    slice,
    actions: slice.actions,
    reducer: slice.reducer,
    selectors: {
      selectItems: (state: { [key: string]: AsyncState<T[]> }) => state[name]?.data || [],
      selectItemById: (id: T[typeof idField]) => (state: { [key: string]: AsyncState<T[]> }) => 
        state[name]?.data?.find(item => item[idField] === id),
      selectLoading: (state: { [key: string]: AsyncState<T[]> }) => state[name]?.loading,
      selectError: (state: { [key: string]: AsyncState<T[]> }) => state[name]?.error,
      selectInitialized: (state: { [key: string]: AsyncState<T[]> }) => state[name]?.initialized,
    },
  };
}