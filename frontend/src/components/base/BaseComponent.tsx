import React, { Component } from 'react';
import type { ReactNode } from 'react';
import type { IInitializable, IConfigurable, BaseState } from '../../types/base';

/**
 * 基础组件属性接口
 */
export interface BaseComponentProps {
  className?: string;
  style?: React.CSSProperties;
  loading?: boolean;
  error?: string | null;
  onError?: (error: Error) => void;
  onLoading?: (loading: boolean) => void;
}

/**
 * 基础组件状态接口
 */
export interface BaseComponentState extends BaseState {
  mounted: boolean;
}

/**
 * 基础组件抽象类
 * 提供通用的组件功能和生命周期管理
 */
export abstract class BaseComponent<P extends BaseComponentProps = BaseComponentProps, S extends BaseComponentState = BaseComponentState> 
  extends Component<P, S> 
  implements IInitializable, IConfigurable {

  protected config: Record<string, any> = {};

  constructor(props: P) {
    super(props);
    this.state = this.getInitialState();
  }

  /**
   * 获取初始状态
   */
  protected getInitialState(): S {
    return {
      loading: this.props.loading || false,
      error: this.props.error || null,
      initialized: false,
      mounted: false
    } as S;
  }

  /**
   * 组件挂载后的初始化
   */
  async componentDidMount(): Promise<void> {
    this.setState({ mounted: true } as Pick<S, keyof S>);
    
    try {
      await this.initialize();
      this.setState({ initialized: true } as Pick<S, keyof S>);
    } catch (error) {
      this.handleError(error as Error);
    }
  }

  /**
   * 组件卸载前的清理
   */
  componentWillUnmount(): void {
    this.setState({ mounted: false } as Pick<S, keyof S>);
    this.cleanup();
  }

  /**
   * 属性更新时的处理
   */
  componentDidUpdate(prevProps: P): void {
    if (prevProps.loading !== this.props.loading) {
      this.setState({ loading: this.props.loading || false } as Pick<S, keyof S>);
    }
    
    if (prevProps.error !== this.props.error) {
      this.setState({ error: this.props.error || null } as Pick<S, keyof S>);
    }
  }

  /**
   * 错误边界处理
   */
  componentDidCatch(error: Error, errorInfo: React.ErrorInfo): void {
    console.error('Component error:', error, errorInfo);
    this.handleError(error);
  }

  // IInitializable接口实现
  abstract initialize(): Promise<void>;

  isInitialized(): boolean {
    return this.state.initialized;
  }

  // IConfigurable接口实现
  configure(config: Record<string, any>): void {
    this.config = { ...this.config, ...config };
  }

  getConfig(): Record<string, any> {
    return { ...this.config };
  }

  /**
   * 设置加载状态
   */
  protected setLoading(loading: boolean): void {
    this.setState({ loading } as Pick<S, keyof S>);
    this.props.onLoading?.(loading);
  }

  /**
   * 处理错误
   */
  protected handleError(error: Error): void {
    const errorMessage = error.message || '未知错误';
    this.setState({ 
      error: errorMessage, 
      loading: false 
    } as Pick<S, keyof S>);
    
    this.props.onError?.(error);
    console.error('Component error:', error);
  }

  /**
   * 清除错误
   */
  protected clearError(): void {
    this.setState({ error: null } as Pick<S, keyof S>);
  }

  /**
   * 安全的状态更新（检查组件是否已挂载）
   */
  protected safeSetState(state: Partial<S>): void {
    if (this.state.mounted) {
      this.setState(state as Pick<S, keyof S>);
    }
  }

  /**
   * 清理资源（子类可重写）
   */
  protected cleanup(): void {
    // 子类可以重写此方法进行资源清理
  }

  /**
   * 渲染加载状态
   */
  protected renderLoading(): ReactNode {
    return <div className="loading">加载中...</div>;
  }

  /**
   * 渲染错误状态
   */
  protected renderError(): ReactNode {
    return (
      <div className="error">
        <p>发生错误: {this.state.error}</p>
        <button onClick={() => this.clearError()}>重试</button>
      </div>
    );
  }

  /**
   * 抽象渲染方法（子类必须实现）
   */
  protected abstract renderContent(): ReactNode;

  /**
   * 主渲染方法
   */
  render(): ReactNode {
    const { className, style } = this.props;
    const { loading, error } = this.state;

    return (
      <div className={`base-component ${className || ''}`} style={style}>
        {error && this.renderError()}
        {loading && this.renderLoading()}
        {!loading && !error && this.renderContent()}
      </div>
    );
  }
}