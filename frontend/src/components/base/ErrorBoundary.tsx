import React, { Component } from 'react';
import type { ErrorInfo, ReactNode } from 'react';
import type { ILogger } from '../../types/base';
import { colors, spacing, borderRadius, shadows, theme, combine, layout } from '../../styles/utils';
import { withTranslation } from 'react-i18next';
import type { WithTranslation } from 'react-i18next';

/**
 * 错误边界属性接口
 */
export interface ErrorBoundaryProps extends WithTranslation {
  children: ReactNode;
  fallback?: ReactNode | ((error: Error, errorInfo: ErrorInfo) => ReactNode);
  onError?: (error: Error, errorInfo: ErrorInfo) => void;
  logger?: ILogger;
  resetOnPropsChange?: boolean;
  resetKeys?: Array<string | number>;
  isolate?: boolean;
}

/**
 * 错误边界状态接口
 */
export interface ErrorBoundaryState {
  hasError: boolean;
  error: Error | null;
  errorInfo: ErrorInfo | null;
  errorId: string | null;
}

/**
 * 简单日志记录器实现
 */
class SimpleLogger implements ILogger {
  debug(message: string, ...args: any[]): void {
    console.debug(`[DEBUG] ${message}`, ...args);
  }

  info(message: string, ...args: any[]): void {
    console.info(`[INFO] ${message}`, ...args);
  }

  warn(message: string, ...args: any[]): void {
    console.warn(`[WARN] ${message}`, ...args);
  }

  error(message: string, ...args: any[]): void {
    console.error(`[ERROR] ${message}`, ...args);
  }
}

/**
 * 错误边界组件
 * 用于捕获子组件中的JavaScript错误，记录错误并显示备用UI
 */
export class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  private resetTimeoutId: number | null = null;
  private logger: ILogger;

  constructor(props: ErrorBoundaryProps) {
    super(props);
    
    this.state = {
      hasError: false,
      error: null,
      errorInfo: null,
      errorId: null
    };

    this.logger = props.logger || new SimpleLogger();
  }

  static getDerivedStateFromError(error: Error): Partial<ErrorBoundaryState> {
    // 更新状态以显示错误UI
    return {
      hasError: true,
      error,
      errorId: `error_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
    };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    // 记录错误信息
    this.setState({ errorInfo });
    
    // 调用错误回调
    this.props.onError?.(error, errorInfo);
    
    // 记录到日志
    this.logger.error('React Error Boundary caught an error', {
      error: error.message,
      stack: error.stack,
      componentStack: errorInfo.componentStack,
      errorId: this.state.errorId
    });

    // 如果不是隔离模式，可以上报错误到监控系统
    if (!this.props.isolate) {
      this.reportError(error, errorInfo);
    }
  }

  componentDidUpdate(prevProps: ErrorBoundaryProps) {
    const { resetOnPropsChange, resetKeys } = this.props;
    const { hasError } = this.state;

    // 如果有错误且启用了属性变化重置
    if (hasError && resetOnPropsChange) {
      if (resetKeys) {
        // 检查重置键是否发生变化
        const hasResetKeyChanged = resetKeys.some((key, index) => 
          prevProps.resetKeys?.[index] !== key
        );
        if (hasResetKeyChanged) {
          this.resetErrorBoundary();
        }
      } else {
        // 检查所有属性是否发生变化
        if (prevProps.children !== this.props.children) {
          this.resetErrorBoundary();
        }
      }
    }
  }

  componentWillUnmount() {
    if (this.resetTimeoutId) {
      clearTimeout(this.resetTimeoutId);
    }
  }

  /**
   * 重置错误边界状态
   */
  resetErrorBoundary = () => {
    if (this.resetTimeoutId) {
      clearTimeout(this.resetTimeoutId);
    }

    this.setState({
      hasError: false,
      error: null,
      errorInfo: null,
      errorId: null
    });
  };

  /**
   * 延迟重置错误边界
   */
  resetErrorBoundaryDelayed = (delay: number = 100) => {
    this.resetTimeoutId = window.setTimeout(() => {
      this.resetErrorBoundary();
    }, delay);
  };

  /**
   * 上报错误到监控系统
   */
  private reportError(error: Error, errorInfo: ErrorInfo) {
    // 这里可以集成错误监控服务，如Sentry、Bugsnag等
    try {
      // 示例：发送到错误监控服务
      if (typeof window !== 'undefined' && (window as any).errorReporter) {
        (window as any).errorReporter.captureException(error, {
          extra: {
            componentStack: errorInfo.componentStack,
            errorId: this.state.errorId
          }
        });
      }
    } catch (reportingError) {
      this.logger.error('Failed to report error to monitoring service', reportingError);
    }
  }

  /**
   * 渲染默认错误UI
   */
  private renderDefaultErrorUI() {
    const { error, errorInfo, errorId } = this.state;
    const isDevelopment = typeof window !== 'undefined' && 
      (window as any).__DEV__ !== false;

    // 创建样式
    const errorBoundaryStyle = combine(
      layout.flexCenter(),
      {
        minHeight: '200px',
        backgroundColor: colors.background('light'),
        border: `1px solid ${colors.border()}`,
        ...borderRadius.get('md'),
        ...spacing.padding('lg')
      }
    );

    const containerStyle = combine(
      {
        textAlign: 'center' as const,
        maxWidth: '500px'
      }
    );

    const titleStyle = combine(
      {
        color: colors.error(),
        fontSize: '1.5rem',
        fontWeight: 600,
        ...spacing.marginBottom('md')
      }
    );

    const messageStyle = combine(
      {
        color: colors.text('secondary'),
        lineHeight: 1.5,
        ...spacing.marginBottom('lg')
      }
    );

    const actionsStyle = combine(
      layout.flex(),
      {
        justifyContent: 'center',
        gap: spacing.get('sm'),
        ...spacing.marginBottom('md')
      }
    );

    const buttonBaseStyle = combine(
      {
        border: 'none',
        cursor: 'pointer',
        fontSize: '14px',
        transition: 'background-color 0.2s',
        ...spacing.padding('md'),
        ...borderRadius.get('sm')
      }
    );

    const primaryButtonStyle = combine(
      buttonBaseStyle,
      {
        backgroundColor: colors.primary(),
        color: '#ffffff'
      }
    );

    const secondaryButtonStyle = combine(
      buttonBaseStyle,
      {
        backgroundColor: colors.text('secondary'),
        color: '#ffffff'
      }
    );

    const detailsStyle = combine(
      {
        textAlign: 'left' as const,
        ...spacing.marginTop('md')
      }
    );

    const errorInfoStyle = combine(
      {
        backgroundColor: colors.background('light'),
        fontFamily: 'monospace',
        fontSize: '12px',
        ...spacing.padding('sm'),
        ...borderRadius.get('sm')
      }
    );

    const stackStyle = combine(
      {
        backgroundColor: colors.background('dark'),
        overflowX: 'auto' as const,
        whiteSpace: 'pre-wrap' as const,
        ...spacing.padding('xs'),
        ...spacing.marginY('xs'),
        ...borderRadius.get('sm')
      }
    );
    
    return (
      <div style={errorBoundaryStyle}>
        <div style={containerStyle}>
          <h2 style={titleStyle}>{this.props.t('common.somethingWentWrong')}</h2>
          <p style={messageStyle}>
            {this.props.t('common.unexpectedErrorMessage')}
          </p>
          
          <div style={actionsStyle}>
            <button 
              style={primaryButtonStyle}
              onClick={this.resetErrorBoundary}
            >
              {this.props.t('common.retry')}
            </button>
            <button 
              style={secondaryButtonStyle}
              onClick={() => window.location.reload()}
            >
              {this.props.t('common.refreshPage')}
            </button>
          </div>

          {isDevelopment && (
            <details style={detailsStyle}>
              <summary>{this.props.t('common.errorDetails')}</summary>
              <div style={errorInfoStyle}>
                <p><strong>{this.props.t('common.errorId')}:</strong> {errorId}</p>
                <p><strong>{this.props.t('common.errorMessage')}:</strong> {error?.message}</p>
                <pre style={stackStyle}>
                  {error?.stack}
                </pre>
                {errorInfo && (
                  <pre style={stackStyle}>
                    {errorInfo.componentStack}
                  </pre>
                )}
              </div>
            </details>
          )}
        </div>
      </div>
    );
  }

  render() {
    const { hasError, error, errorInfo } = this.state;
    const { children, fallback } = this.props;

    if (hasError) {
      // 如果提供了自定义fallback
      if (fallback) {
        if (typeof fallback === 'function') {
          return fallback(error!, errorInfo!);
        }
        return fallback;
      }
      
      // 使用默认错误UI
      return this.renderDefaultErrorUI();
    }

    return children;
  }
}

/**
 * 错误边界Hook
 * 用于函数组件中的错误处理
 */
export function useErrorHandler() {
  const [error, setError] = React.useState<Error | null>(null);

  const resetError = React.useCallback(() => {
    setError(null);
  }, []);

  const captureError = React.useCallback((error: Error) => {
    setError(error);
  }, []);

  // 如果有错误，抛出它以便ErrorBoundary捕获
  React.useEffect(() => {
    if (error) {
      throw error;
    }
  }, [error]);

  return { captureError, resetError };
}

/**
 * 高阶组件：为组件添加错误边界
 */
export function withErrorBoundary<P extends Record<string, any>>(
  Component: React.ComponentType<P>,
  errorBoundaryProps?: Omit<ErrorBoundaryProps, 'children'>
) {
  const WrappedComponent = React.forwardRef<any, P>((props, ref) => (
    <ErrorBoundary {...errorBoundaryProps}>
      <Component {...(props as any)} ref={ref} />
    </ErrorBoundary>
  ));

  WrappedComponent.displayName = `withErrorBoundary(${Component.displayName || Component.name})`;

  return WrappedComponent;
}

/**
 * 错误边界样式（可以移到CSS文件中）
 */
export const errorBoundaryStyles = `
.error-boundary {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 200px;
  padding: 20px;
  background-color: #f8f9fa;
  border: 1px solid #dee2e6;
  border-radius: 8px;
}

.error-boundary__container {
  text-align: center;
  max-width: 500px;
}

.error-boundary__title {
  color: #dc3545;
  margin-bottom: 16px;
  font-size: 1.5rem;
  font-weight: 600;
}

.error-boundary__message {
  color: #6c757d;
  margin-bottom: 24px;
  line-height: 1.5;
}

.error-boundary__actions {
  display: flex;
  gap: 12px;
  justify-content: center;
  margin-bottom: 20px;
}

.error-boundary__button {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.2s;
}

.error-boundary__button--primary {
  background-color: #007bff;
  color: white;
}

.error-boundary__button--primary:hover {
  background-color: #0056b3;
}

.error-boundary__button--secondary {
  background-color: #6c757d;
  color: white;
}

.error-boundary__button--secondary:hover {
  background-color: #545b62;
}

.error-boundary__details {
  text-align: left;
  margin-top: 20px;
}

.error-boundary__error-info {
  background-color: #f8f9fa;
  padding: 12px;
  border-radius: 4px;
  font-family: monospace;
  font-size: 12px;
}

.error-boundary__stack,
.error-boundary__component-stack {
  background-color: #e9ecef;
  padding: 8px;
  border-radius: 4px;
  overflow-x: auto;
  white-space: pre-wrap;
  margin: 8px 0;
}
`;

// 导出国际化的ErrorBoundary组件
export default withTranslation()(ErrorBoundary);