import { useCallback } from 'react';
import { notification, message } from 'antd';
import type { NotificationArgsProps } from 'antd';

// 通知类型
export type NotificationType = 'success' | 'info' | 'warning' | 'error';

// 通知配置
export interface NotificationConfig extends Omit<NotificationArgsProps, 'type'> {
  type?: NotificationType;
  autoClose?: boolean;
  duration?: number;
}

// 消息配置
export interface MessageConfig {
  content: string;
  duration?: number;
  type?: NotificationType;
}

// 通知hook
export function useNotification() {
  // 显示通知
  const showNotification = useCallback((config: NotificationConfig) => {
    const {
      type = 'info',
      autoClose = true,
      duration = 4.5,
      ...notificationProps
    } = config;

    const finalConfig = {
      ...notificationProps,
      duration: autoClose ? duration : 0,
    };

    notification[type](finalConfig);
  }, []);

  // 成功通知
  const success = useCallback((config: Omit<NotificationConfig, 'type'>) => {
    showNotification({ ...config, type: 'success' });
  }, [showNotification]);

  // 信息通知
  const info = useCallback((config: Omit<NotificationConfig, 'type'>) => {
    showNotification({ ...config, type: 'info' });
  }, [showNotification]);

  // 警告通知
  const warning = useCallback((config: Omit<NotificationConfig, 'type'>) => {
    showNotification({ ...config, type: 'warning' });
  }, [showNotification]);

  // 错误通知
  const error = useCallback((config: Omit<NotificationConfig, 'type'>) => {
    showNotification({ ...config, type: 'error' });
  }, [showNotification]);

  // 关闭通知
  const close = useCallback((key: string) => {
    notification.destroy(key);
  }, []);

  // 关闭所有通知
  const closeAll = useCallback(() => {
    notification.destroy();
  }, []);

  return {
    showNotification,
    success,
    info,
    warning,
    error,
    close,
    closeAll,
  };
}

// 消息hook
export function useMessage() {
  // 显示消息
  const showMessage = useCallback((config: MessageConfig) => {
    const { type = 'info', ...messageProps } = config;
    return message[type](messageProps);
  }, []);

  // 成功消息
  const success = useCallback((content: string, duration?: number) => {
    return message.success(content, duration);
  }, []);

  // 信息消息
  const info = useCallback((content: string, duration?: number) => {
    return message.info(content, duration);
  }, []);

  // 警告消息
  const warning = useCallback((content: string, duration?: number) => {
    return message.warning(content, duration);
  }, []);

  // 错误消息
  const error = useCallback((content: string, duration?: number) => {
    return message.error(content, duration);
  }, []);

  // 加载消息
  const loading = useCallback((content: string, duration?: number) => {
    return message.loading(content, duration);
  }, []);

  // 关闭所有消息
  const destroy = useCallback(() => {
    message.destroy();
  }, []);

  return {
    showMessage,
    success,
    info,
    warning,
    error,
    loading,
    destroy,
  };
}

// 操作反馈hook
export function useOperationFeedback() {
  const notification = useNotification();
  const message = useMessage();

  // 操作成功反馈
  const operationSuccess = useCallback((operation: string, details?: string) => {
    notification.success({
      message: `${operation}成功`,
      description: details,
    });
  }, [notification]);

  // 操作失败反馈
  const operationError = useCallback((operation: string, error: any) => {
    const errorMessage = error?.message || error?.toString() || '操作失败';
    notification.error({
      message: `${operation}失败`,
      description: errorMessage,
    });
  }, [notification]);

  // 操作警告反馈
  const operationWarning = useCallback((operation: string, warning: string) => {
    notification.warning({
      message: `${operation}警告`,
      description: warning,
    });
  }, [notification]);

  // 简单成功消息
  const simpleSuccess = useCallback((text: string) => {
    message.success(text);
  }, [message]);

  // 简单错误消息
  const simpleError = useCallback((text: string) => {
    message.error(text);
  }, [message]);

  // 简单警告消息
  const simpleWarning = useCallback((text: string) => {
    message.warning(text);
  }, [message]);

  // 加载反馈
  const loadingFeedback = useCallback((text: string = '加载中...') => {
    return message.loading(text, 0);
  }, [message]);

  return {
    operationSuccess,
    operationError,
    operationWarning,
    simpleSuccess,
    simpleError,
    simpleWarning,
    loadingFeedback,
  };
}

// API操作反馈hook
export function useApiFeedback() {
  const feedback = useOperationFeedback();

  // 创建操作反馈
  const createFeedback = useCallback((resourceName: string) => ({
    success: (details?: string) => feedback.operationSuccess(`创建${resourceName}`, details),
    error: (error: any) => feedback.operationError(`创建${resourceName}`, error),
  }), [feedback]);

  // 更新操作反馈
  const updateFeedback = useCallback((resourceName: string) => ({
    success: (details?: string) => feedback.operationSuccess(`更新${resourceName}`, details),
    error: (error: any) => feedback.operationError(`更新${resourceName}`, error),
  }), [feedback]);

  // 删除操作反馈
  const deleteFeedback = useCallback((resourceName: string) => ({
    success: (details?: string) => feedback.operationSuccess(`删除${resourceName}`, details),
    error: (error: any) => feedback.operationError(`删除${resourceName}`, error),
  }), [feedback]);

  // 获取操作反馈
  const fetchFeedback = useCallback((resourceName: string) => ({
    success: (details?: string) => feedback.operationSuccess(`获取${resourceName}`, details),
    error: (error: any) => feedback.operationError(`获取${resourceName}`, error),
  }), [feedback]);

  // 批量操作反馈
  const batchFeedback = useCallback((operation: string, resourceName: string) => ({
    success: (count: number) => feedback.operationSuccess(`批量${operation}${resourceName}`, `成功处理 ${count} 项`),
    error: (error: any) => feedback.operationError(`批量${operation}${resourceName}`, error),
  }), [feedback]);

  return {
    createFeedback,
    updateFeedback,
    deleteFeedback,
    fetchFeedback,
    batchFeedback,
  };
}

// 表单验证反馈hook
export function useValidationFeedback() {
  const message = useMessage();

  // 验证失败反馈
  const validationError = useCallback((field: string, error: string) => {
    message.error(`${field}: ${error}`);
  }, [message]);

  // 表单提交失败反馈
  const submitError = useCallback((errors: Record<string, string>) => {
    const errorMessages = Object.entries(errors)
      .map(([field, error]) => `${field}: ${error}`)
      .join('; ');
    message.error(`表单验证失败: ${errorMessages}`);
  }, [message]);

  // 表单提交成功反馈
  const submitSuccess = useCallback((text: string = '提交成功') => {
    message.success(text);
  }, [message]);

  return {
    validationError,
    submitError,
    submitSuccess,
  };
}