import { message, notification } from 'antd';

export interface ErrorInfo {
  code?: string;
  message: string;
  details?: any;
  timestamp: string;
  source?: 'api' | 'sse' | 'network' | 'validation' | 'unknown';
  severity?: 'low' | 'medium' | 'high' | 'critical';
}

export interface ErrorHandlerOptions {
  showNotification?: boolean;
  showMessage?: boolean;
  logToConsole?: boolean;
  reportToServer?: boolean;
  customHandler?: (error: ErrorInfo) => void;
}

class ErrorHandler {
  private errorLog: ErrorInfo[] = [];
  private maxLogSize = 100;

  // 处理错误的主要方法
  handle(error: any, options: ErrorHandlerOptions = {}): void {
    const errorInfo = this.normalizeError(error);
    
    // 记录错误
    this.logError(errorInfo);

    // 默认选项
    const defaultOptions: ErrorHandlerOptions = {
      showNotification: true,
      showMessage: false,
      logToConsole: true,
      reportToServer: false,
    };

    const finalOptions = { ...defaultOptions, ...options };

    // 控制台日志
    if (finalOptions.logToConsole) {
      this.logToConsole(errorInfo);
    }

    // 显示用户通知
    if (finalOptions.showNotification) {
      this.showNotification(errorInfo);
    }

    // 显示消息提示
    if (finalOptions.showMessage) {
      this.showMessage(errorInfo);
    }

    // 上报到服务器
    if (finalOptions.reportToServer) {
      this.reportToServer(errorInfo);
    }

    // 自定义处理器
    if (finalOptions.customHandler) {
      finalOptions.customHandler(errorInfo);
    }
  }

  // 标准化错误对象
  private normalizeError(error: any): ErrorInfo {
    const timestamp = new Date().toISOString();

    // 如果已经是ErrorInfo格式
    if (error && typeof error === 'object' && error.message && error.timestamp) {
      return error as ErrorInfo;
    }

    // 处理API错误
    if (error?.response) {
      return {
        code: error.response.status?.toString(),
        message: error.response.data?.message || error.message || '请求失败',
        details: error.response.data,
        timestamp,
        source: 'api',
        severity: this.getSeverityFromStatus(error.response.status),
      };
    }

    // 处理网络错误
    if (error?.code === 'NETWORK_ERROR' || error?.message?.includes('Network Error')) {
      return {
        code: 'NETWORK_ERROR',
        message: '网络连接失败，请检查网络设置',
        details: error,
        timestamp,
        source: 'network',
        severity: 'high',
      };
    }

    // 处理SSE错误
    if (error?.type === 'error' && error?.target instanceof EventSource) {
      return {
        code: 'SSE_ERROR',
        message: 'SSE连接错误',
        details: error,
        timestamp,
        source: 'sse',
        severity: 'medium',
      };
    }

    // 处理验证错误
    if (error?.name === 'ValidationError') {
      return {
        code: 'VALIDATION_ERROR',
        message: error.message || '数据验证失败',
        details: error,
        timestamp,
        source: 'validation',
        severity: 'low',
      };
    }

    // 处理JavaScript错误
    if (error instanceof Error) {
      return {
        code: error.name,
        message: error.message,
        details: { stack: error.stack },
        timestamp,
        source: 'unknown',
        severity: 'medium',
      };
    }

    // 处理字符串错误
    if (typeof error === 'string') {
      return {
        message: error,
        timestamp,
        source: 'unknown',
        severity: 'low',
      };
    }

    // 默认错误
    return {
      message: '未知错误',
      details: error,
      timestamp,
      source: 'unknown',
      severity: 'low',
    };
  }

  // 根据HTTP状态码确定严重程度
  private getSeverityFromStatus(status: number): ErrorInfo['severity'] {
    if (status >= 500) return 'critical';
    if (status >= 400) return 'high';
    if (status >= 300) return 'medium';
    return 'low';
  }

  // 记录错误到内存
  private logError(error: ErrorInfo): void {
    this.errorLog.unshift(error);
    if (this.errorLog.length > this.maxLogSize) {
      this.errorLog = this.errorLog.slice(0, this.maxLogSize);
    }
  }

  // 控制台日志
  private logToConsole(error: ErrorInfo): void {
    const logLevel = this.getLogLevel(error.severity);
    const logMessage = `[${error.source?.toUpperCase()}] ${error.message}`;
    
    console[logLevel](logMessage, error.details);
  }

  // 获取日志级别
  private getLogLevel(severity?: string): 'log' | 'warn' | 'error' {
    switch (severity) {
      case 'critical':
      case 'high':
        return 'error';
      case 'medium':
        return 'warn';
      default:
        return 'log';
    }
  }

  // 显示通知
  private showNotification(error: ErrorInfo): void {
    const { severity, message: errorMessage, code } = error;
    
    const config = {
      message: this.getNotificationTitle(severity),
      description: errorMessage,
      duration: this.getNotificationDuration(severity),
    };

    switch (severity) {
      case 'critical':
      case 'high':
        notification.error(config);
        break;
      case 'medium':
        notification.warning(config);
        break;
      default:
        notification.info(config);
        break;
    }
  }

  // 显示消息提示
  private showMessage(error: ErrorInfo): void {
    const { severity, message: errorMessage } = error;
    
    switch (severity) {
      case 'critical':
      case 'high':
        message.error(errorMessage);
        break;
      case 'medium':
        message.warning(errorMessage);
        break;
      default:
        message.info(errorMessage);
        break;
    }
  }

  // 获取通知标题
  private getNotificationTitle(severity?: string): string {
    switch (severity) {
      case 'critical':
        return '严重错误';
      case 'high':
        return '错误';
      case 'medium':
        return '警告';
      default:
        return '提示';
    }
  }

  // 获取通知持续时间
  private getNotificationDuration(severity?: string): number {
    switch (severity) {
      case 'critical':
        return 0; // 不自动关闭
      case 'high':
        return 8;
      case 'medium':
        return 5;
      default:
        return 3;
    }
  }

  // 上报错误到服务器
  private async reportToServer(error: ErrorInfo): Promise<void> {
    try {
      await fetch('/api/errors/report', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(error),
      });
    } catch (reportError) {
      console.error('上报错误失败:', reportError);
    }
  }

  // 获取错误日志
  getErrorLog(): ErrorInfo[] {
    return [...this.errorLog];
  }

  // 清空错误日志
  clearErrorLog(): void {
    this.errorLog = [];
  }

  // 获取特定类型的错误
  getErrorsBySource(source: ErrorInfo['source']): ErrorInfo[] {
    return this.errorLog.filter(error => error.source === source);
  }

  // 获取特定严重程度的错误
  getErrorsBySeverity(severity: ErrorInfo['severity']): ErrorInfo[] {
    return this.errorLog.filter(error => error.severity === severity);
  }

  // 检查是否有严重错误
  hasCriticalErrors(): boolean {
    return this.errorLog.some(error => error.severity === 'critical');
  }
}

// 创建全局错误处理器实例
export const errorHandler = new ErrorHandler();

// 便捷方法
export const handleError = (error: any, options?: ErrorHandlerOptions) => {
  errorHandler.handle(error, options);
};

export const handleAPIError = (error: any) => {
  errorHandler.handle(error, {
    showNotification: true,
    showMessage: false,
    logToConsole: true,
    reportToServer: true,
  });
};

export const handleSSEError = (error: any) => {
  errorHandler.handle(error, {
    showNotification: false,
    showMessage: true,
    logToConsole: true,
    reportToServer: false,
  });
};

export const handleValidationError = (error: any) => {
  errorHandler.handle(error, {
    showNotification: false,
    showMessage: true,
    logToConsole: false,
    reportToServer: false,
  });
};

// 全局错误监听器
window.addEventListener('error', (event) => {
  errorHandler.handle(event.error, {
    showNotification: true,
    reportToServer: true,
  });
});

window.addEventListener('unhandledrejection', (event) => {
  errorHandler.handle(event.reason, {
    showNotification: true,
    reportToServer: true,
  });
});

export default ErrorHandler;