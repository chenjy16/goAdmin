import { message } from 'antd';
import { useTranslation } from 'react-i18next';
import { useMemo, useCallback } from 'react';

// API错误类型
export interface ApiError {
  code?: string;
  message: string;
  details?: any;
  status?: number;
}

// 错误处理配置
export interface ErrorHandlerConfig {
  showMessage?: boolean;
  messageType?: 'error' | 'warning' | 'info';
  customMessage?: string;
  onError?: (error: ApiError) => void;
  silent?: boolean;
}

// 默认错误处理配置
const defaultConfig: ErrorHandlerConfig = {
  showMessage: true,
  messageType: 'error',
  silent: false,
};

// 错误消息映射函数
const getErrorMessageMap = (t: (key: string) => string): Record<string, string> => ({
  NETWORK_ERROR: t('apiErrors.networkError'),
  TIMEOUT: t('apiErrors.timeout'),
  UNAUTHORIZED: t('apiErrors.unauthorized'),
  FORBIDDEN: t('apiErrors.forbidden'),
  NOT_FOUND: t('apiErrors.notFound'),
  INTERNAL_SERVER_ERROR: t('apiErrors.internalServerError'),
  BAD_REQUEST: t('apiErrors.badRequest'),
  CONFLICT: t('apiErrors.conflict'),
  TOO_MANY_REQUESTS: t('apiErrors.tooManyRequests'),
});

// 根据HTTP状态码获取错误消息
export const getErrorMessageByStatus = (status: number, t: (key: string) => string): string => {
  const errorMessageMap = getErrorMessageMap(t);
  switch (status) {
    case 400:
      return errorMessageMap.BAD_REQUEST;
    case 401:
      return errorMessageMap.UNAUTHORIZED;
    case 403:
      return errorMessageMap.FORBIDDEN;
    case 404:
      return errorMessageMap.NOT_FOUND;
    case 409:
      return errorMessageMap.CONFLICT;
    case 429:
      return errorMessageMap.TOO_MANY_REQUESTS;
    case 500:
      return errorMessageMap.INTERNAL_SERVER_ERROR;
    default:
      return t('apiErrors.defaultError');
  }
};

// 解析错误对象
export function parseApiError(error: any, t: (key: string) => string): ApiError {
  const errorMessageMap = getErrorMessageMap(t);
  
  // 如果已经是ApiError格式
  if (error && typeof error === 'object' && 'message' in error) {
    return {
      code: error.code,
      message: error.message,
      details: error.details,
      status: error.status,
    };
  }

  // 处理Axios错误
  if (error?.response) {
    const { status, data } = error.response;
    return {
      status,
      code: data?.code || `HTTP_${status}`,
      message: data?.message || getErrorMessageByStatus(status, t),
      details: data?.details || data,
    };
  }

  // 处理网络错误
  if (error?.code === 'NETWORK_ERROR' || error?.message?.includes('Network Error')) {
    return {
      code: 'NETWORK_ERROR',
      message: errorMessageMap.NETWORK_ERROR,
    };
  }

  // 处理超时错误
  if (error?.code === 'ECONNABORTED' || error?.message?.includes('timeout')) {
    return {
      code: 'TIMEOUT',
      message: errorMessageMap.TIMEOUT,
    };
  }

  // 处理字符串错误
  if (typeof error === 'string') {
    return {
      message: error,
    };
  }

  // 处理Error对象
  if (error instanceof Error) {
    return {
      message: error.message,
    };
  }

  // 默认错误
  return {
    message: t('apiErrors.unknownError'),
    details: error,
  };
}

// Hook for API error handling with internationalization
export const useApiErrorHandler = () => {
  const { t } = useTranslation();

  const getErrorMessageByStatusWithI18n = useMemo(() => 
    (status: number) => getErrorMessageByStatus(status, t), 
    [t]
  );

  const parseApiErrorWithI18n = useMemo(() => 
    (error: any) => parseApiError(error, t), 
    [t]
  );

  const handleApiError = useCallback(
    (error: any, config: ErrorHandlerConfig = {}): ApiError => {
      const finalConfig = { ...defaultConfig, ...config };
      const apiError = parseApiErrorWithI18n(error);

      // 使用自定义消息
      if (finalConfig.customMessage) {
        apiError.message = finalConfig.customMessage;
      }

      // 显示错误消息
      if (finalConfig.showMessage && !finalConfig.silent) {
        const messageType = finalConfig.messageType || 'error';
        message[messageType](apiError.message);
      }

      // 执行自定义错误处理
      if (finalConfig.onError) {
        finalConfig.onError(apiError);
      }

      // 记录错误日志
      console.error('API Error:', {
        code: apiError.code,
        message: apiError.message,
        status: apiError.status,
        details: apiError.details,
        originalError: error,
      });

      return apiError;
    }, 
    [parseApiErrorWithI18n]
  );

  return {
    getErrorMessageByStatus: getErrorMessageByStatusWithI18n,
    parseApiError: parseApiErrorWithI18n,
    handleApiError,
  };
};

// 处理API错误 (保留原函数用于向后兼容)
export function handleApiError(error: any, config: ErrorHandlerConfig = {}): ApiError {
  const { t } = useTranslation();
  const finalConfig = { ...defaultConfig, ...config };
  const apiError = parseApiError(error, t);

  // 使用自定义消息
  if (finalConfig.customMessage) {
    apiError.message = finalConfig.customMessage;
  }

  // 显示错误消息
  if (finalConfig.showMessage && !finalConfig.silent) {
    const messageType = finalConfig.messageType || 'error';
    message[messageType](apiError.message);
  }

  // 执行自定义错误处理
  if (finalConfig.onError) {
    finalConfig.onError(apiError);
  }

  // 记录错误日志
  console.error('API Error:', {
    code: apiError.code,
    message: apiError.message,
    status: apiError.status,
    details: apiError.details,
    originalError: error,
  });

  return apiError;
}

// 创建错误处理装饰器
export function withErrorHandler<T extends (...args: any[]) => Promise<any>>(
  fn: T,
  config: ErrorHandlerConfig = {}
): T {
  return (async (...args: any[]) => {
    try {
      return await fn(...args);
    } catch (error) {
      const apiError = handleApiError(error, config);
      throw apiError;
    }
  }) as T;
}

// 创建重试装饰器
export interface RetryConfig {
  maxRetries?: number;
  retryDelay?: number;
  retryCondition?: (error: ApiError) => boolean;
  onRetry?: (attempt: number, error: ApiError) => void;
}

// 创建重试函数工厂 (需要在组件中使用 useApiErrorHandler)
export const createRetryFunction = (parseApiErrorFn: (error: any) => ApiError) => {
  return function withRetry<T extends (...args: any[]) => Promise<any>>(
    fn: T,
    config: RetryConfig = {}
  ): T {
    const {
      maxRetries = 3,
      retryDelay = 1000,
      retryCondition = (error) => error.status === 500 || error.code === 'NETWORK_ERROR',
      onRetry,
    } = config;

    return (async (...args: any[]) => {
      let lastError: ApiError;
      
      for (let attempt = 0; attempt <= maxRetries; attempt++) {
        try {
          return await fn(...args);
        } catch (error) {
          lastError = parseApiErrorFn(error);
          
          // 如果是最后一次尝试或不满足重试条件，直接抛出错误
          if (attempt === maxRetries || !retryCondition(lastError)) {
            throw lastError;
          }
          
          // 执行重试回调
          if (onRetry) {
            onRetry(attempt + 1, lastError);
          }
          
          // 等待重试延迟
          if (retryDelay > 0) {
            await new Promise(resolve => setTimeout(resolve, retryDelay));
          }
        }
      }
      
      throw lastError!;
    }) as T;
  };
};

// 保留原函数用于向后兼容 (但不推荐在组件中使用)
export function withRetry<T extends (...args: any[]) => Promise<any>>(
  fn: T,
  config: RetryConfig = {}
): T {
  const {
    maxRetries = 3,
    retryDelay = 1000,
    retryCondition = (error) => error.status === 500 || error.code === 'NETWORK_ERROR',
    onRetry,
  } = config;

  return (async (...args: any[]) => {
    let lastError: ApiError;
    
    for (let attempt = 0; attempt <= maxRetries; attempt++) {
      try {
        return await fn(...args);
      } catch (error) {
        // 注意：这里在非组件中使用 useTranslation 可能会有问题
        const { t } = useTranslation();
        lastError = parseApiError(error, t);
        
        // 如果是最后一次尝试或不满足重试条件，直接抛出错误
        if (attempt === maxRetries || !retryCondition(lastError)) {
          throw lastError;
        }
        
        // 执行重试回调
        if (onRetry) {
          onRetry(attempt + 1, lastError);
        }
        
        // 等待重试延迟
        if (retryDelay > 0) {
          await new Promise(resolve => setTimeout(resolve, retryDelay));
        }
      }
    }
    
    throw lastError!;
  }) as T;
}

// 组合错误处理和重试
export function withErrorHandlingAndRetry<T extends (...args: any[]) => Promise<any>>(
  fn: T,
  errorConfig: ErrorHandlerConfig = {},
  retryConfig: RetryConfig = {}
): T {
  return withErrorHandler(withRetry(fn, retryConfig), errorConfig);
}

// 类型已在上面导出