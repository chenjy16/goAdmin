import { useState, useEffect, useCallback, useRef } from 'react';
import { handleApiError, type ErrorHandlerConfig } from '../utils/apiErrorHandler';

// 异步数据状态
export interface AsyncDataState<T> {
  data: T | null;
  loading: boolean;
  error: string | null;
  initialized: boolean;
}

// 异步数据配置
export interface UseAsyncDataConfig<T> {
  initialData?: T;
  immediate?: boolean;
  errorConfig?: ErrorHandlerConfig;
  onSuccess?: (data: T) => void;
  onError?: (error: any) => void;
  deps?: React.DependencyList;
}

// 异步数据hook
export function useAsyncData<T>(
  asyncFn: () => Promise<T>,
  config: UseAsyncDataConfig<T> = {}
) {
  const {
    initialData = null,
    immediate = true,
    errorConfig = {},
    onSuccess,
    onError,
    deps = [],
  } = config;

  const [state, setState] = useState<AsyncDataState<T>>({
    data: initialData,
    loading: false,
    error: null,
    initialized: false,
  });

  const mountedRef = useRef(true);
  const abortControllerRef = useRef<AbortController | null>(null);

  // 执行异步操作
  const execute = useCallback(async () => {
    // 取消之前的请求
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }

    // 创建新的AbortController
    abortControllerRef.current = new AbortController();

    setState(prev => ({
      ...prev,
      loading: true,
      error: null,
    }));

    try {
      const result = await asyncFn();
      
      if (mountedRef.current) {
        setState(prev => ({
          ...prev,
          data: result,
          loading: false,
          initialized: true,
          error: null,
        }));

        if (onSuccess) {
          onSuccess(result);
        }
      }

      return result;
    } catch (error) {
      if (mountedRef.current) {
        const apiError = handleApiError(error, errorConfig);
        
        setState(prev => ({
          ...prev,
          loading: false,
          error: apiError.message,
        }));

        if (onError) {
          onError(error);
        }
      }
      
      throw error;
    }
  }, [asyncFn, onSuccess, onError, errorConfig]);

  // 重置状态
  const reset = useCallback(() => {
    setState({
      data: initialData,
      loading: false,
      error: null,
      initialized: false,
    });
  }, [initialData]);

  // 清除错误
  const clearError = useCallback(() => {
    setState(prev => ({
      ...prev,
      error: null,
    }));
  }, []);

  // 自动执行
  useEffect(() => {
    if (immediate) {
      execute();
    }
  }, deps);

  // 清理
  useEffect(() => {
    return () => {
      mountedRef.current = false;
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
    };
  }, []);

  return {
    ...state,
    execute,
    reset,
    clearError,
    refresh: execute,
  };
}

// 异步操作hook（不自动执行）
export function useAsyncOperation<T, P = void>(
  asyncFn: (params: P) => Promise<T>,
  config: Omit<UseAsyncDataConfig<T>, 'immediate' | 'deps'> = {}
) {
  const [state, setState] = useState<AsyncDataState<T>>({
    data: config.initialData || null,
    loading: false,
    error: null,
    initialized: false,
  });

  const mountedRef = useRef(true);

  const execute = useCallback(async (params: P) => {
    setState(prev => ({
      ...prev,
      loading: true,
      error: null,
    }));

    try {
      const result = await asyncFn(params);
      
      if (mountedRef.current) {
        setState(prev => ({
          ...prev,
          data: result,
          loading: false,
          initialized: true,
          error: null,
        }));

        if (config.onSuccess) {
          config.onSuccess(result);
        }
      }

      return result;
    } catch (error) {
      if (mountedRef.current) {
        const apiError = handleApiError(error, config.errorConfig);
        
        setState(prev => ({
          ...prev,
          loading: false,
          error: apiError.message,
        }));

        if (config.onError) {
          config.onError(error);
        }
      }
      
      throw error;
    }
  }, [asyncFn, config.onSuccess, config.onError, config.errorConfig]);

  const reset = useCallback(() => {
    setState({
      data: config.initialData || null,
      loading: false,
      error: null,
      initialized: false,
    });
  }, [config.initialData]);

  const clearError = useCallback(() => {
    setState(prev => ({
      ...prev,
      error: null,
    }));
  }, []);

  useEffect(() => {
    return () => {
      mountedRef.current = false;
    };
  }, []);

  return {
    ...state,
    execute,
    reset,
    clearError,
  };
}

// 分页数据hook
export interface PaginationConfig {
  page?: number;
  pageSize?: number;
  total?: number;
}

export interface UsePaginatedDataConfig<T> extends UseAsyncDataConfig<T[]> {
  pagination?: PaginationConfig;
}

export function usePaginatedData<T>(
  asyncFn: (pagination: PaginationConfig) => Promise<{ data: T[]; total: number }>,
  config: UsePaginatedDataConfig<T> = {}
) {
  const [pagination, setPagination] = useState<PaginationConfig>({
    page: 1,
    pageSize: 10,
    total: 0,
    ...config.pagination,
  });

  const paginatedAsyncFn = useCallback(async () => {
    const result = await asyncFn(pagination);
    setPagination(prev => ({
      ...prev,
      total: result.total,
    }));
    return result.data;
  }, [asyncFn, pagination]);

  const asyncData = useAsyncData(paginatedAsyncFn, {
    ...config,
    deps: [pagination.page, pagination.pageSize, ...(config.deps || [])],
  });

  const changePage = useCallback((page: number) => {
    setPagination(prev => ({
      ...prev,
      page,
    }));
  }, []);

  const changePageSize = useCallback((pageSize: number) => {
    setPagination(prev => ({
      ...prev,
      page: 1,
      pageSize,
    }));
  }, []);

  return {
    ...asyncData,
    pagination,
    changePage,
    changePageSize,
  };
}

// 类型已在上面导出