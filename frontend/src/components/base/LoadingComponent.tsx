import React from 'react';
import type { ReactNode } from 'react';
import { colors, spacing, size as sizeUtils, layout, borderRadius, animation, combine } from '../../styles';
import { useTranslation } from 'react-i18next';

/**
 * 加载组件类型
 */
export type LoadingType = 'spinner' | 'dots' | 'pulse' | 'skeleton' | 'progress';

/**
 * 加载组件大小
 */
export type LoadingSize = 'small' | 'medium' | 'large';

/**
 * 加载组件属性接口
 */
export interface LoadingComponentProps {
  loading?: boolean;
  type?: LoadingType;
  size?: LoadingSize;
  color?: string;
  text?: string;
  overlay?: boolean;
  children?: ReactNode;
  className?: string;
  style?: React.CSSProperties;
  delay?: number;
  timeout?: number;
  onTimeout?: () => void;
}

/**
 * 骨架屏属性接口
 */
export interface SkeletonProps {
  width?: string | number;
  height?: string | number;
  rows?: number;
  avatar?: boolean;
  title?: boolean;
  paragraph?: boolean;
  className?: string;
}

/**
 * 进度条属性接口
 */
export interface ProgressProps {
  percent?: number;
  showPercent?: boolean;
  status?: 'normal' | 'success' | 'error';
  strokeWidth?: number;
  className?: string;
}

/**
 * 旋转器组件
 */
const Spinner: React.FC<{ size: LoadingSize; color?: string }> = ({ size, color = colors.primary() }) => {
  const sizeMap: Record<LoadingSize, number> = {
    small: 16,
    medium: 24,
    large: 32
  };

  const currentSize = sizeMap[size];
  const spinnerStyle = combine(
    sizeUtils.square(currentSize),
    borderRadius.get('round'),
    {
      border: '2px solid transparent',
      borderColor: `${color}20`,
      borderTopColor: color,
      animation: 'loading-spin 1s linear infinite'
    }
  );

  return (
    <div
      className={`loading-spinner loading-spinner--${size}`}
      style={spinnerStyle}
    />
  );
};

/**
 * 点状加载器组件
 */
const Dots: React.FC<{ size: LoadingSize; color?: string }> = ({ size, color = colors.primary() }) => {
  const sizeMap: Record<LoadingSize, number> = {
    small: 4,
    medium: 6,
    large: 8
  };

  const containerStyle = combine(
    layout.flex('row', 'center', 'center'),
    { gap: spacing.get('xs') }
  );

  return (
    <div className={`loading-dots loading-dots--${size}`} style={containerStyle}>
      {[0, 1, 2].map(i => (
        <div
          key={i}
          className="loading-dots__dot"
          style={combine(
             sizeUtils.square(sizeMap[size]),
             borderRadius.get('round'),
             {
               backgroundColor: color,
               animation: 'loading-bounce 1.4s ease-in-out infinite both',
               animationDelay: `${i * 0.1}s`
             }
           )}
        />
      ))}
    </div>
  );
};

/**
 * 脉冲加载器组件
 */
const Pulse: React.FC<{ size: LoadingSize; color?: string }> = ({ size, color = colors.primary() }) => {
  const sizeMap: Record<LoadingSize, number> = {
    small: 16,
    medium: 24,
    large: 32
  };

  const pulseStyle = combine(
    sizeUtils.square(sizeMap[size]),
    borderRadius.get('round'),
    {
      backgroundColor: color,
      animation: 'loading-pulse-animation 1.5s ease-in-out infinite'
    }
  );

  return (
    <div
      className={`loading-pulse loading-pulse--${size}`}
      style={pulseStyle}
    />
  );
};

/**
 * 骨架屏组件
 */
const Skeleton: React.FC<SkeletonProps> = ({
  width = '100%',
  height = '20px',
  rows = 3,
  avatar = false,
  title = false,
  paragraph = true,
  className = ''
}) => {
  const skeletonStyle = combine(
    layout.flex('row'),
    { gap: spacing.get('sm') }
  );

  const avatarStyle = combine(
    sizeUtils.square(40),
    borderRadius.get('round'),
    {
      background: 'linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%)',
      backgroundSize: '200% 100%',
      animation: 'loading-skeleton-animation 1.5s infinite'
    }
  );

  const titleStyle = combine(
    {
      height: 20,
      width: typeof width === 'number' ? `${width}px` : width,
      background: 'linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%)',
      backgroundSize: '200% 100%',
      animation: 'loading-skeleton-animation 1.5s infinite',
      marginBottom: spacing.get('xs')
    },
    borderRadius.get('sm')
  );

  const lineStyle = (isLast: boolean) => combine(
    {
      height: typeof height === 'number' ? `${height}px` : height,
      width: isLast ? '60%' : '100%',
      background: 'linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%)',
      backgroundSize: '200% 100%',
      animation: 'loading-skeleton-animation 1.5s infinite'
    },
    borderRadius.get('sm')
  );

  return (
    <div className={`loading-skeleton ${className}`} style={skeletonStyle}>
      {avatar && (
        <div className="loading-skeleton__avatar" style={avatarStyle} />
      )}
      <div className="loading-skeleton__content" style={{ flex: 1 }}>
        {title && (
          <div className="loading-skeleton__title" style={titleStyle} />
        )}
        {paragraph && (
          <div className="loading-skeleton__paragraph" style={combine(layout.flex('column'), { gap: spacing.get('xs') })}>
            {Array.from({ length: rows }).map((_, index) => (
               <div
                 key={index}
                 className="loading-skeleton__line"
                 style={lineStyle(index === rows - 1)}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

/**
 * 进度条组件
 */
const Progress: React.FC<ProgressProps> = ({
  percent = 0,
  showPercent = true,
  status = 'normal',
  strokeWidth = 8,
  className = ''
}) => {
  const getStatusColor = () => {
    switch (status) {
      case 'success': return colors.success();
      case 'error': return colors.error();
      default: return colors.primary();
    }
  };

  const progressStyle = combine(
    layout.flex('row'),
    { alignItems: 'center', gap: spacing.get('xs') }
  );

  const trackStyle = combine(
    {
      flex: 1,
      height: strokeWidth,
      backgroundColor: colors.background('light'),
      overflow: 'hidden'
    },
    borderRadius.get('sm')
  );

  const barStyle = combine(
    {
      height: '100%',
      width: `${Math.min(100, Math.max(0, percent))}%`,
      backgroundColor: getStatusColor(),
      transition: 'width 0.3s ease'
    },
    borderRadius.get('sm')
  );

  const textStyle = {
    fontSize: '12px',
    color: colors.text('secondary'),
    minWidth: '35px',
    textAlign: 'right' as const
  };

  return (
    <div className={`loading-progress ${className}`} style={progressStyle}>
      <div className="loading-progress__track" style={trackStyle}>
        <div className="loading-progress__bar" style={barStyle} />
      </div>
      {showPercent && (
        <span className="loading-progress__text" style={textStyle}>
          {Math.round(percent)}%
        </span>
      )}
    </div>
  );
};

/**
 * 延迟显示Hook
 */
const useDelayedLoading = (loading: boolean, delay: number = 200) => {
  const [delayedLoading, setDelayedLoading] = React.useState(false);

  React.useEffect(() => {
    let timeoutId: number;

    if (loading) {
      timeoutId = window.setTimeout(() => {
        setDelayedLoading(true);
      }, delay);
    } else {
      setDelayedLoading(false);
    }

    return () => {
      if (timeoutId) {
        clearTimeout(timeoutId);
      }
    };
  }, [loading, delay]);

  return delayedLoading;
};

/**
 * 超时处理Hook
 */
const useLoadingTimeout = (loading: boolean, timeout?: number, onTimeout?: () => void) => {
  React.useEffect(() => {
    if (!loading || !timeout || !onTimeout) return;

    const timeoutId = window.setTimeout(() => {
      onTimeout();
    }, timeout);

    return () => clearTimeout(timeoutId);
  }, [loading, timeout, onTimeout]);
};

/**
 * 主加载组件
 */
export const LoadingComponent: React.FC<LoadingComponentProps> = ({
  loading = false,
  type = 'spinner',
  size = 'medium',
  color,
  text,
  overlay = false,
  children,
  className = '',
  style,
  delay = 0,
  timeout,
  onTimeout
}) => {
  const delayedLoading = useDelayedLoading(loading, delay);
  useLoadingTimeout(delayedLoading, timeout, onTimeout);

  const renderLoadingIndicator = () => {
    const props = { size, color };

    switch (type) {
      case 'spinner':
        return <Spinner {...props} />;
      case 'dots':
        return <Dots {...props} />;
      case 'pulse':
        return <Pulse {...props} />;
      case 'skeleton':
        return <Skeleton />;
      case 'progress':
        return <Progress />;
      default:
        return <Spinner {...props} />;
    }
  };

  const contentStyle = combine(
    layout.flexCenter(),
    layout.flexColumn(),
    { gap: spacing.get('xs') },
    spacing.padding('lg'),
    style
  );

  const textStyle = {
    fontSize: '14px',
    color: colors.text('secondary'),
    textAlign: 'center' as const
  };

  const loadingContent = (
    <div className={`loading-content loading-content--${type} ${className}`} style={contentStyle}>
      {renderLoadingIndicator()}
      {text && <div className="loading-text" style={textStyle}>{text}</div>}
    </div>
  );

  if (!delayedLoading) {
    return children ? <>{children}</> : null;
  }

  if (overlay) {
    const overlayStyle = { position: 'relative' as const };
    
    const backdropStyle = combine(
      layout.absolute('0', '0', '0', '0'),
      layout.flexCenter(),
      {
        backgroundColor: 'rgba(255, 255, 255, 0.8)',
        zIndex: 1000
      }
    );

    const hiddenContentStyle = {
      opacity: 0.5,
      pointerEvents: 'none' as const
    };

    return (
      <div className="loading-overlay" style={overlayStyle}>
        <div 
          className={`loading-wrapper__content ${delayedLoading ? 'loading-wrapper__content--hidden' : ''}`}
          style={delayedLoading ? hiddenContentStyle : undefined}
        >
          {children}
        </div>
        {delayedLoading && (
          <div className="loading-overlay__backdrop" style={backdropStyle}>
            {loadingContent}
          </div>
        )}
      </div>
    );
  }

  return children ? (
    <div className="loading-wrapper">
      <div className={`loading-wrapper__content ${delayedLoading ? 'loading-wrapper__content--hidden' : ''}`}>
        {children}
      </div>
      {delayedLoading && (
        <div className="loading-wrapper__loading">
          {loadingContent}
        </div>
      )}
    </div>
  ) : loadingContent;
};

/**
 * 高阶组件：为组件添加加载状态
 */
export function withLoading<P extends Record<string, any>>(
  Component: React.ComponentType<P>,
  loadingProps?: Partial<LoadingComponentProps>
) {
  return React.forwardRef<any, P & { loading?: boolean }>((props, ref) => {
    const { loading, ...restProps } = props;

    return (
      <LoadingComponent loading={loading} {...loadingProps}>
        <Component {...(restProps as unknown as P)} ref={ref} />
      </LoadingComponent>
    );
  });
}

/**
 * 加载组件样式
 */
export const loadingStyles = `
/* 旋转器样式 */
.loading-spinner {
  border: 2px solid transparent;
  border-radius: 50%;
  animation: loading-spin 1s linear infinite;
}

@keyframes loading-spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 点状加载器样式 */
.loading-dots {
  display: flex;
  gap: 4px;
  align-items: center;
}

.loading-dots__dot {
  border-radius: 50%;
  animation: loading-bounce 1.4s ease-in-out infinite both;
}

@keyframes loading-bounce {
  0%, 80%, 100% {
    transform: scale(0);
  }
  40% {
    transform: scale(1);
  }
}

/* 脉冲加载器样式 */
.loading-pulse {
  border-radius: 50%;
  animation: loading-pulse-animation 1.5s ease-in-out infinite;
}

@keyframes loading-pulse-animation {
  0% {
    transform: scale(0);
    opacity: 1;
  }
  100% {
    transform: scale(1);
    opacity: 0;
  }
}

/* 骨架屏样式 */
.loading-skeleton {
  display: flex;
  gap: 12px;
}

.loading-skeleton__avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: loading-skeleton-animation 1.5s infinite;
}

.loading-skeleton__content {
  flex: 1;
}

.loading-skeleton__title {
  height: 20px;
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: loading-skeleton-animation 1.5s infinite;
  margin-bottom: 8px;
  border-radius: 4px;
}

.loading-skeleton__paragraph {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.loading-skeleton__line {
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: loading-skeleton-animation 1.5s infinite;
  border-radius: 4px;
}

@keyframes loading-skeleton-animation {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}

/* 进度条样式 */
.loading-progress {
  display: flex;
  align-items: center;
  gap: 8px;
}

.loading-progress__track {
  flex: 1;
  background-color: #f0f0f0;
  border-radius: 4px;
  overflow: hidden;
}

.loading-progress__bar {
  height: 100%;
  transition: width 0.3s ease;
  border-radius: 4px;
}

.loading-progress__text {
  font-size: 12px;
  color: #666;
  min-width: 35px;
  text-align: right;
}

/* 加载内容样式 */
.loading-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 20px;
}

.loading-text {
  font-size: 14px;
  color: #666;
  text-align: center;
}

/* 覆盖层样式 */
.loading-overlay {
  position: relative;
}

.loading-overlay__backdrop {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(255, 255, 255, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

/* 包装器样式 */
.loading-wrapper {
  position: relative;
}

.loading-wrapper__content--hidden {
  opacity: 0.5;
  pointer-events: none;
}

.loading-wrapper__loading {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 1000;
}
`;