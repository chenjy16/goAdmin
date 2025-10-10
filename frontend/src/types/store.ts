import type { ProviderInfo, ModelInfo, MCPTool, MCPMessage, ChatMessage } from './api';
import type { ConfigSliceState } from '../store/slices/configSlice';

// 提供商状态
export interface ProvidersState {
  providers: ProviderInfo[];
  models: Record<string, ModelInfo[]>; // provider -> models
  apiKeyStatus: Record<string, any>; // provider -> APIKeyInfo
  selectedProvider: string | null;
  isLoading: boolean;
  error: string | null;
}

// MCP工具状态
export interface MCPState {
  tools: MCPTool[];
  logs: MCPMessage[];
  isInitialized: boolean;
  isLoading: boolean;
  error: string | null;
  executionResults: Record<string, any>;
}

// AI助手状态
export interface AssistantState {
  conversations: AssistantConversation[];
  currentConversationId: string | null;
  isInitialized: boolean;
  isLoading: boolean;
  error: string | null;
}

export interface AssistantConversation {
  id: string;
  title: string;
  messages: ChatMessage[];
  createdAt: string;
  updatedAt: string;
}

// 应用设置状态
export interface SettingsState {
  theme: 'light' | 'dark';
  language: 'zh' | 'en';
  apiKeys: Record<string, string>; // provider -> api_key
  defaultProvider: string;
  defaultModel: Record<string, string>; // provider -> model
  chatSettings: {
    maxTokens: number;
    temperature: number;
    streamResponse: boolean;
  };
}

// UI状态
export interface UIState {
  sidebarCollapsed: boolean;
  currentPage: string;
  loading: Record<string, boolean>;
  notifications: Notification[];
}

export interface Notification {
  id: string;
  type: 'success' | 'error' | 'warning' | 'info';
  title: string;
  message: string;
  duration?: number;
  timestamp: string;
}

// 根状态
export interface RootState {
  providers: ProvidersState;
  mcp: MCPState;
  assistant: AssistantState;
  settings: SettingsState;
  ui: UIState;
  config: ConfigSliceState;
}

// 异步操作状态
export interface AsyncState {
  loading: boolean;
  error: string | null;
}

// 分页相关类型
export interface PaginationState {
  page: number;
  pageSize: number;
  total: number;
}

// 过滤和搜索状态
export interface FilterState {
  searchTerm: string;
  filters: Record<string, any>;
  sortBy: string;
  sortOrder: 'asc' | 'desc';
}