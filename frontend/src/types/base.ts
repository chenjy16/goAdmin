/**
 * 基础API响应接口
 */
export interface BaseApiResponse<T = any> {
  code: number;
  message: string;
  data?: T;
  error?: string;
  status?: string;
}

/**
 * 基础实体接口
 */
export interface BaseEntity {
  id?: string;
  createdAt?: string;
  updatedAt?: string;
}

/**
 * 基础状态接口
 */
export interface BaseState {
  loading: boolean;
  error: string | null;
  initialized: boolean;
}

/**
 * 分页接口
 */
export interface Pagination {
  page: number;
  pageSize: number;
  total: number;
}

/**
 * 排序接口
 */
export interface Sort {
  field: string;
  order: 'asc' | 'desc';
}

/**
 * 过滤接口
 */
export interface Filter {
  field: string;
  operator: 'eq' | 'ne' | 'gt' | 'gte' | 'lt' | 'lte' | 'like' | 'in';
  value: any;
}

/**
 * 查询参数接口
 */
export interface QueryParams {
  pagination?: Pagination;
  sort?: Sort[];
  filters?: Filter[];
  search?: string;
}

/**
 * 服务接口
 */
export interface IService<T, CreateDTO = Partial<T>, UpdateDTO = Partial<T>> {
  getAll(params?: QueryParams): Promise<BaseApiResponse<T[]>>;
  getById(id: string): Promise<BaseApiResponse<T>>;
  create(data: CreateDTO): Promise<BaseApiResponse<T>>;
  update(id: string, data: UpdateDTO): Promise<BaseApiResponse<T>>;
  delete(id: string): Promise<BaseApiResponse<void>>;
}

/**
 * 配置接口
 */
export interface IConfigurable {
  configure(config: Record<string, any>): void;
  getConfig(): Record<string, any>;
}

/**
 * 可初始化接口
 */
export interface IInitializable {
  initialize(): Promise<void>;
  isInitialized(): boolean;
}

/**
 * 可验证接口
 */
export interface IValidatable<T = any> {
  validate(data: T): ValidationResult;
}

/**
 * 验证结果接口
 */
export interface ValidationResult {
  isValid: boolean;
  errors: ValidationError[];
}

/**
 * 验证错误接口
 */
export interface ValidationError {
  field: string;
  message: string;
  code?: string;
}

/**
 * 事件接口
 */
export interface IEventEmitter {
  on(event: string, listener: (...args: any[]) => void): void;
  off(event: string, listener: (...args: any[]) => void): void;
  emit(event: string, ...args: any[]): void;
}

/**
 * 缓存接口
 */
export interface ICache<T = any> {
  get(key: string): T | null;
  set(key: string, value: T, ttl?: number): void;
  delete(key: string): void;
  clear(): void;
  has(key: string): boolean;
}

/**
 * 日志级别
 */
export type LogLevel = 'debug' | 'info' | 'warn' | 'error';

/**
 * 日志接口
 */
export interface ILogger {
  debug(message: string, ...args: any[]): void;
  info(message: string, ...args: any[]): void;
  warn(message: string, ...args: any[]): void;
  error(message: string, ...args: any[]): void;
}

/**
 * 主题类型
 */
export type Theme = 'light' | 'dark';

/**
 * 语言类型
 */
export type Language = 'zh' | 'en';

/**
 * 通知类型
 */
export type NotificationType = 'success' | 'error' | 'warning' | 'info';

/**
 * 通知接口
 */
export interface Notification extends BaseEntity {
  type: NotificationType;
  title: string;
  message: string;
  duration?: number;
  timestamp: string;
}