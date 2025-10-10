import { useState, useCallback } from 'react';
import { message } from 'antd';

interface UseAsyncOperationOptions {
  successMessage?: string;
  errorMessage?: string;
  showSuccessMessage?: boolean;
  showErrorMessage?: boolean;
}

interface AsyncOperationState {
  loading: boolean;
  error: string | null;
}

export function useAsyncOperation<T extends any[], R>(
  asyncFunction: (...args: T) => Promise<R>,
  options: UseAsyncOperationOptions = {}
) {
  const {
    successMessage,
    errorMessage = '操作失败',
    showSuccessMessage = true,
    showErrorMessage = true,
  } = options;

  const [state, setState] = useState<AsyncOperationState>({
    loading: false,
    error: null,
  });

  const execute = useCallback(
    async (...args: T): Promise<R | null> => {
      setState({ loading: true, error: null });
      
      try {
        const result = await asyncFunction(...args);
        
        setState({ loading: false, error: null });
        
        if (showSuccessMessage && successMessage) {
          message.success(successMessage);
        }
        
        return result;
      } catch (error) {
        const errorMsg = error instanceof Error ? error.message : errorMessage;
        setState({ loading: false, error: errorMsg });
        
        if (showErrorMessage) {
          message.error(errorMsg);
        }
        
        return null;
      }
    },
    [asyncFunction, successMessage, errorMessage, showSuccessMessage, showErrorMessage]
  );

  const reset = useCallback(() => {
    setState({ loading: false, error: null });
  }, []);

  return {
    ...state,
    execute,
    reset,
  };
}

export default useAsyncOperation;