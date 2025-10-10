import axios from 'axios';
import type { AxiosInstance, AxiosResponse, AxiosRequestConfig, InternalAxiosRequestConfig } from 'axios';

/**
 * 基础服务抽象类
 * 提供通用的HTTP请求封装和错误处理
 */
export abstract class BaseService {
  protected api: AxiosInstance;
  protected baseURL: string;

  constructor(baseURL?: string) {
    this.baseURL = baseURL || (import.meta.env.DEV ? '' : (import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'));
    this.api = this.createAxiosInstance();
    this.setupInterceptors();
  }

  /**
   * 创建axios实例
   */
  private createAxiosInstance(): AxiosInstance {
    return axios.create({
      baseURL: this.baseURL,
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    });
  }

  /**
   * 设置请求和响应拦截器
   */
  private setupInterceptors(): void {
    // 请求拦截器
    this.api.interceptors.request.use(
      (config) => {
        return this.onRequest(config);
      },
      (error) => {
        return this.onRequestError(error);
      }
    );

    // 响应拦截器
    this.api.interceptors.response.use(
      (response: AxiosResponse) => {
        return this.onResponse(response);
      },
      (error) => {
        return this.onResponseError(error);
      }
    );
  }

  /**
   * 请求前处理
   */
  protected onRequest(config: InternalAxiosRequestConfig): InternalAxiosRequestConfig {
    // 子类可以重写此方法添加认证token等
    return config;
  }

  /**
   * 请求错误处理
   */
  protected onRequestError(error: any): Promise<any> {
    return Promise.reject(error);
  }

  /**
   * 响应处理
   */
  protected onResponse(response: AxiosResponse): AxiosResponse {
    return response;
  }

  /**
   * 响应错误处理
   */
  protected onResponseError(error: any): Promise<any> {
    if (error.response) {
      // 服务器响应错误
      const { status, data } = error.response;
      return Promise.reject(new ServiceError(
        data?.message || `HTTP ${status} Error`,
        status,
        data
      ));
    } else if (error.request) {
      // 网络错误
      return Promise.reject(new ServiceError(
        'Network Error: Unable to connect to server',
        0,
        error.request
      ));
    } else {
      // 其他错误
      return Promise.reject(new ServiceError(
        error.message || 'Unknown Error',
        -1,
        error
      ));
    }
  }

  /**
   * GET请求
   */
  protected async get<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.api.get<T>(url, config);
    return response.data;
  }

  /**
   * POST请求
   */
  protected async post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.api.post<T>(url, data, config);
    return response.data;
  }

  /**
   * PUT请求
   */
  protected async put<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.api.put<T>(url, data, config);
    return response.data;
  }

  /**
   * DELETE请求
   */
  protected async delete<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.api.delete<T>(url, config);
    return response.data;
  }
}

/**
 * 服务错误类
 */
export class ServiceError extends Error {
  public readonly code: number;
  public readonly data?: any;

  constructor(message: string, code: number, data?: any) {
    super(message);
    this.name = 'ServiceError';
    this.code = code;
    this.data = data;
  }
}