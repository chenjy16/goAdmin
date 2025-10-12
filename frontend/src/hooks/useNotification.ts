import { useCallback } from 'react';
import { notification, message } from 'antd';
import type { NotificationArgsProps } from 'antd';
import { useTranslation } from 'react-i18next';

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
  const { t } = useTranslation();
  const notification = useNotification();
  const message = useMessage();

  // 操作成功反馈
  const operationSuccess = useCallback((operation: string, details?: string) => {
    notification.success({
      message: `${operation}${t('common.success')}`,
      description: details,
    });
  }, [notification, t]);

  // 操作失败反馈
  const operationError = useCallback((operation: string, error: any) => {
    const errorMessage = error?.message || error?.toString() || t('notifications.operationFailed');
    notification.error({
      message: `${operation}${t('common.failed')}`,
      description: errorMessage,
    });
  }, [notification, t]);

  // 操作警告反馈
  const operationWarning = useCallback((operation: string, warning: string) => {
    notification.warning({
      message: `${operation}${t('notifications.warning')}`,
      description: warning,
    });
  }, [notification, t]);

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
  const loadingFeedback = useCallback((text: string = t('notifications.loading')) => {
    return message.loading(text, 0);
  }, [message, t]);

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
export const useApiFeedback = () => {
  const { t } = useTranslation();
  
  const showApiError = useCallback((error: any, defaultMessage?: string) => {
    const message = error?.response?.data?.message || 
                   error?.message || 
                   defaultMessage || 
                   t('notifications.defaultError');
    
    notification.error({
      message: t('common.error'),
      description: message,
      placement: 'topRight',
      duration: 4,
    });
  }, [t]);

  const showApiSuccess = useCallback((message?: string, description?: string) => {
    notification.success({
      message: message || t('common.success'),
      description,
      placement: 'topRight',
      duration: 3,
    });
  }, [t]);

  return {
    showApiError,
    showApiSuccess,
  };
};

// 表单验证反馈hook
export const useValidationFeedback = () => {
  const { t } = useTranslation();
  
  const validationError = useCallback((message?: string) => {
    notification.error({
      message: t('notifications.formValidationFailed'),
      description: message || t('notifications.defaultWarning'),
      placement: 'topRight',
      duration: 4,
    });
  }, [t]);

  const submitError = useCallback((error: any) => {
    const message = error?.message || error?.toString() || t('notifications.defaultError');
    notification.error({
      message: t('common.error'),
      description: message,
      placement: 'topRight',
      duration: 4,
    });
  }, [t]);

  const submitSuccess = useCallback((message?: string) => {
    notification.success({
      message: t('notifications.submitSuccess'),
      description: message,
      placement: 'topRight',
      duration: 3,
    });
  }, [t]);

  return {
    validationError,
    submitError,
    submitSuccess,
  };
};