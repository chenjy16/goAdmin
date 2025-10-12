import React, { useState, useEffect, useCallback, useMemo } from 'react';
import type { ReactNode } from 'react';
import type { BaseState } from '../../types/base';
import { useTranslation } from 'react-i18next';

/**
 * 基础Hook组件属性接口
 */
export interface BaseHookComponentProps {
  className?: string;
  style?: React.CSSProperties;
  loading?: boolean;
  error?: string | null;
  onError?: (error: Error) => void;
  onLoading?: (loading: boolean) => void;
  children?: ReactNode;
}

/**
 * 基础状态Hook
 */
export function useBaseState(initialProps?: Partial<BaseState>) {
  const [state, setState] = useState<BaseState>({
    loading: initialProps?.loading || false,
    error: initialProps?.error || null,
    initialized: initialProps?.initialized || false
  });

  const setLoading = useCallback((loading: boolean) => {
    setState(prev => ({ ...prev, loading }));
  }, []);

  const setError = useCallback((error: string | null) => {
    setState(prev => ({ ...prev, error, loading: false }));
  }, []);

  const setInitialized = useCallback((initialized: boolean) => {
    setState(prev => ({ ...prev, initialized }));
  }, []);

  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }));
  }, []);

  const handleError = useCallback((error: Error, t?: (key: string) => string) => {
    const errorMessage = error.message || (t ? t('common.unknownError') : 'Unknown Error');
    setError(errorMessage);
    console.error('Component error:', error);
  }, [setError]);

  return {
    ...state,
    setLoading,
    setError,
    setInitialized,
    clearError,
    handleError
  };
}

/**
 * 异步操作Hook
 */
export function useAsyncOperation<T = any>() {
  const { loading, error, setLoading, clearError, handleError } = useBaseState();

  const execute = useCallback(async (operation: () => Promise<T>): Promise<T | null> => {
    try {
      setLoading(true);
      clearError();
      const result = await operation();
      setLoading(false);
      return result;
    } catch (error) {
      handleError(error as Error);
      return null;
    }
  }, [setLoading, clearError, handleError]);

  return {
    loading,
    error,
    execute,
    clearError
  };
}

/**
 * 配置Hook
 */
export function useConfig<T = Record<string, any>>(initialConfig?: T) {
  const [config, setConfig] = useState<T>(initialConfig || {} as T);

  const updateConfig = useCallback((newConfig: Partial<T>) => {
    setConfig(prev => ({ ...prev, ...newConfig }));
  }, []);

  const resetConfig = useCallback(() => {
    setConfig(initialConfig || {} as T);
  }, [initialConfig]);

  return {
    config,
    updateConfig,
    resetConfig
  };
}

/**
 * 初始化Hook
 */
export function useInitialization(initFunction: () => Promise<void>, deps: React.DependencyList = []) {
  const { initialized, setInitialized, handleError } = useBaseState();

  useEffect(() => {
    let mounted = true;

    const initialize = async () => {
      try {
        await initFunction();
        if (mounted) {
          setInitialized(true);
        }
      } catch (error) {
        if (mounted) {
          handleError(error as Error);
        }
      }
    };

    initialize();

    return () => {
      mounted = false;
    };
  }, deps); // eslint-disable-line react-hooks/exhaustive-deps

  return { initialized };
}

/**
 * 基础Hook组件高阶组件
 */
export function withBaseComponent<P extends Record<string, any>>(
  WrappedComponent: React.ComponentType<P>
) {
  return React.forwardRef<any, P & BaseHookComponentProps>((props, ref) => {
    const { className, style, loading: propLoading, error: propError, onError, onLoading, ...restProps } = props;
    const { loading, error, clearError } = useBaseState({ loading: propLoading, error: propError });

    // 同步外部状态
    useEffect(() => {
      if (propLoading !== undefined && propLoading !== loading) {
        onLoading?.(propLoading);
      }
    }, [propLoading, loading, onLoading]);

    useEffect(() => {
      if (propError !== undefined && propError !== error) {
        if (propError) {
          onError?.(new Error(propError));
        }
      }
    }, [propError, error, onError]);

    const { t } = useTranslation();

    const renderLoading = useCallback(() => (
      <div className="loading">{t('common.loading')}</div>
    ), [t]);

    const renderError = useCallback(() => (
      <div className="error">
        <p>{t('common.errorOccurred')}: {error}</p>
        <button onClick={clearError}>{t('common.retry')}</button>
      </div>
    ), [error, clearError, t]);

    const containerProps = useMemo(() => ({
      className: `base-component ${className || ''}`,
      style
    }), [className, style]);

    return (
      <div {...containerProps}>
        {error && renderError()}
        {loading && renderLoading()}
        {!loading && !error && (
          <WrappedComponent
            ref={ref}
            {...(restProps as unknown as P)}
          />
        )}
      </div>
    );
  });
}

/**
 * 基础Hook组件配置
 */
export interface BaseHookComponentConfig {
  showLoading?: boolean;
  showError?: boolean;
  autoRetry?: boolean;
  retryDelay?: number;
}

/**
 * 基础Hook组件
 */
export function BaseHookComponent({
  children,
  className,
  style,
  loading,
  error,
  config = {}
}: BaseHookComponentProps & { config?: BaseHookComponentConfig }) {
  const { showLoading = true, showError = true } = config;
  const { clearError } = useBaseState({ loading, error });
  const { t } = useTranslation();

  const renderLoading = () => showLoading && loading && (
    <div className="loading">{t('common.loading')}</div>
  );

  const renderError = () => showError && error && (
    <div className="error">
      <p>{t('common.errorOccurred')}: {error}</p>
      <button onClick={clearError}>{t('common.retry')}</button>
    </div>
  );

  return (
    <div className={`base-component ${className || ''}`} style={style}>
      {renderError()}
      {renderLoading()}
      {!loading && !error && children}
    </div>
  );
}

/**
 * 表单Hook
 */
export function useForm<T extends Record<string, any>>(initialValues: T) {
  const [values, setValues] = useState<T>(initialValues);
  const [errors, setErrors] = useState<Partial<Record<keyof T, string>>>({});
  const [touched, setTouched] = useState<Partial<Record<keyof T, boolean>>>({});

  const setValue = useCallback((field: keyof T, value: any) => {
    setValues(prev => ({ ...prev, [field]: value }));
    // 清除该字段的错误
    if (errors[field]) {
      setErrors(prev => ({ ...prev, [field]: undefined }));
    }
  }, [errors]);

  const setFieldError = useCallback((field: keyof T, error: string) => {
    setErrors(prev => ({ ...prev, [field]: error }));
  }, []);

  const setFieldTouched = useCallback((field: keyof T, touched: boolean = true) => {
    setTouched(prev => ({ ...prev, [field]: touched }));
  }, []);

  const resetForm = useCallback(() => {
    setValues(initialValues);
    setErrors({});
    setTouched({});
  }, [initialValues]);

  const isValid = useMemo(() => {
    return Object.keys(errors).length === 0;
  }, [errors]);

  return {
    values,
    errors,
    touched,
    setValue,
    setFieldError,
    setFieldTouched,
    resetForm,
    isValid
  };
}