import { ProviderService } from './ProviderService';
import { MCPService } from './MCPService';
import { AssistantService } from './AssistantService';
import type { IInitializable } from '../types/base';

/**
 * 服务工厂类
 * 统一管理所有服务实例，提供单例模式和依赖注入
 */
export class ServiceFactory {
  private static instance: ServiceFactory;
  private services: Map<string, any> = new Map();
  private initialized = false;

  private constructor() {}

  /**
   * 获取服务工厂单例
   */
  static getInstance(): ServiceFactory {
    if (!ServiceFactory.instance) {
      ServiceFactory.instance = new ServiceFactory();
    }
    return ServiceFactory.instance;
  }

  /**
   * 初始化所有服务
   */
  async initialize(): Promise<void> {
    if (this.initialized) {
      return;
    }

    try {
      // 创建服务实例
      this.services.set('provider', new ProviderService());
      this.services.set('mcp', new MCPService());
      this.services.set('assistant', new AssistantService());

      // 初始化需要初始化的服务
      const initializableServices = Array.from(this.services.values())
        .filter((service): service is IInitializable => 
          'initialize' in service && typeof service.initialize === 'function'
        );

      await Promise.all(
        initializableServices.map(service => 
          service.initialize().catch(error => {
            console.warn(`Failed to initialize service:`, error);
            // 不抛出错误，允许其他服务继续初始化
          })
        )
      );

      this.initialized = true;
    } catch (error) {
      console.error('Failed to initialize services:', error);
      throw error;
    }
  }

  /**
   * 获取Provider服务
   */
  getProviderService(): ProviderService {
    return this.getService<ProviderService>('provider');
  }

  /**
   * 获取MCP服务
   */
  getMCPService(): MCPService {
    return this.getService<MCPService>('mcp');
  }

  /**
   * 获取Assistant服务
   */
  getAssistantService(): AssistantService {
    return this.getService<AssistantService>('assistant');
  }

  /**
   * 获取指定服务
   */
  private getService<T>(name: string): T {
    const service = this.services.get(name);
    if (!service) {
      throw new Error(`Service '${name}' not found. Make sure to call initialize() first.`);
    }
    return service as T;
  }

  /**
   * 注册自定义服务
   */
  registerService<T>(name: string, service: T): void {
    this.services.set(name, service);
  }

  /**
   * 检查服务是否存在
   */
  hasService(name: string): boolean {
    return this.services.has(name);
  }

  /**
   * 获取所有服务名称
   */
  getServiceNames(): string[] {
    return Array.from(this.services.keys());
  }

  /**
   * 检查是否已初始化
   */
  isInitialized(): boolean {
    return this.initialized;
  }

  /**
   * 重置工厂（主要用于测试）
   */
  reset(): void {
    this.services.clear();
    this.initialized = false;
  }

  /**
   * 获取所有服务的健康状态
   */
  async getHealthStatus(): Promise<Record<string, boolean>> {
    const status: Record<string, boolean> = {};
    
    for (const [name, service] of this.services) {
      try {
        if ('isInitialized' in service && typeof service.isInitialized === 'function') {
          status[name] = service.isInitialized();
        } else {
          status[name] = true; // 假设没有初始化方法的服务是健康的
        }
      } catch (error) {
        status[name] = false;
      }
    }

    return status;
  }
}

// 导出单例实例
export const serviceFactory = ServiceFactory.getInstance();

// 导出便捷的服务获取函数
export const getProviderService = () => serviceFactory.getProviderService();
export const getMCPService = () => serviceFactory.getMCPService();
export const getAssistantService = () => serviceFactory.getAssistantService();

// 默认导出
export default serviceFactory;