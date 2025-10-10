import type { SSEEvent } from '../types/api';

export interface SSEOptions {
  onMessage?: (event: SSEEvent) => void;
  onError?: (error: Event) => void;
  onOpen?: (event: Event) => void;
  onClose?: (event: Event) => void;
  reconnectInterval?: number;
  maxReconnectAttempts?: number;
}

export class SSEService {
  private eventSource: EventSource | null = null;
  private url: string;
  private options: SSEOptions;
  private reconnectAttempts = 0;
  private isConnecting = false;
  private shouldReconnect = true;

  constructor(url: string, options: SSEOptions = {}) {
    this.url = url;
    this.options = {
      reconnectInterval: 3000,
      maxReconnectAttempts: 5,
      ...options,
    };
  }

  connect(): void {
    if (this.isConnecting || this.eventSource?.readyState === EventSource.OPEN) {
      return;
    }

    this.isConnecting = true;
    this.eventSource = new EventSource(this.url);

    this.eventSource.onopen = (event) => {
      this.isConnecting = false;
      this.reconnectAttempts = 0;
      this.options.onOpen?.(event);
    };

    this.eventSource.onmessage = (event) => {
      try {
        const data: SSEEvent = JSON.parse(event.data);
        this.options.onMessage?.(data);
      } catch (error) {
        console.error('解析SSE消息失败:', error);
      }
    };

    this.eventSource.onerror = (event) => {
      console.error('SSE连接错误:', event);
      this.isConnecting = false;
      this.options.onError?.(event);

      if (this.shouldReconnect && this.reconnectAttempts < (this.options.maxReconnectAttempts || 5)) {
        this.reconnect();
      }
    };

    // 监听自定义事件类型
    this.eventSource.addEventListener('chat_message', (event) => {
      try {
        const data: SSEEvent = JSON.parse((event as MessageEvent).data);
        this.options.onMessage?.(data);
      } catch (error) {
        console.error('解析聊天消息失败:', error);
      }
    });

    this.eventSource.addEventListener('chat_error', (event) => {
      try {
        const data: SSEEvent = JSON.parse((event as MessageEvent).data);
        this.options.onMessage?.(data);
      } catch (error) {
        console.error('解析聊天错误失败:', error);
      }
    });

    this.eventSource.addEventListener('chat_done', (event) => {
      try {
        const data: SSEEvent = JSON.parse((event as MessageEvent).data);
        this.options.onMessage?.(data);
      } catch (error) {
        console.error('解析聊天完成事件失败:', error);
      }
    });
  }

  private reconnect(): void {
    if (!this.shouldReconnect) return;

    this.reconnectAttempts++;

    setTimeout(() => {
      if (this.shouldReconnect) {
        this.connect();
      }
    }, this.options.reconnectInterval);
  }

  disconnect(): void {
    this.shouldReconnect = false;
    if (this.eventSource) {
      this.eventSource.close();
      this.eventSource = null;
    }
    this.isConnecting = false;
    this.reconnectAttempts = 0;
  }

  isConnected(): boolean {
    return this.eventSource?.readyState === EventSource.OPEN;
  }

  getReadyState(): number | null {
    return this.eventSource?.readyState || null;
  }
}

// 创建全局SSE管理器
class SSEManager {
  private connections: Map<string, SSEService> = new Map();

  createConnection(id: string, url: string, options?: SSEOptions): SSEService {
    // 如果已存在连接，先断开
    if (this.connections.has(id)) {
      this.connections.get(id)?.disconnect();
    }

    const sseService = new SSEService(url, options);
    this.connections.set(id, sseService);
    return sseService;
  }

  getConnection(id: string): SSEService | undefined {
    return this.connections.get(id);
  }

  disconnectConnection(id: string): void {
    const connection = this.connections.get(id);
    if (connection) {
      connection.disconnect();
      this.connections.delete(id);
    }
  }

  disconnectAll(): void {
    this.connections.forEach((connection) => {
      connection.disconnect();
    });
    this.connections.clear();
  }

  getActiveConnections(): string[] {
    const active: string[] = [];
    this.connections.forEach((connection, id) => {
      if (connection.isConnected()) {
        active.push(id);
      }
    });
    return active;
  }
}

export const sseManager = new SSEManager();

// 聊天SSE连接辅助函数
export const createChatSSE = (
  conversationId: string,
  onMessage: (event: SSEEvent) => void,
  onError?: (error: Event) => void
): SSEService => {
  const url = `/api/v1/chat/${conversationId}/stream`;
  
  return sseManager.createConnection(`chat-${conversationId}`, url, {
    onMessage,
    onError,
    onOpen: () => {
      // Chat SSE connection established
    },
    onClose: () => {
      // Chat SSE connection closed
    },
    reconnectInterval: 2000,
    maxReconnectAttempts: 3,
  });
};

// MCP工具SSE连接辅助函数
export const createMCPSSE = (
  onMessage: (event: SSEEvent) => void,
  onError?: (error: Event) => void
): SSEService => {
  const url = '/api/v1/mcp/sse';
  
  return sseManager.createConnection('mcp-tools', url, {
    onMessage,
    onError,
    onOpen: () => {
      // MCP tools SSE connection established
    },
    onClose: () => {
      // MCP tools SSE connection closed
    },
    reconnectInterval: 3000,
    maxReconnectAttempts: 5,
  });
};

// 助手SSE连接辅助函数
export const createAssistantSSE = (
  conversationId: string,
  onMessage: (event: SSEEvent) => void,
  onError?: (error: Event) => void
): SSEService => {
  const url = `/api/v1/assistant/${conversationId}/stream`;
  
  return sseManager.createConnection(`assistant-${conversationId}`, url, {
    onMessage,
    onError,
    onOpen: () => {
      // Assistant SSE connection established
    },
    onClose: () => {
      // Assistant SSE connection closed
    },
    reconnectInterval: 2000,
    maxReconnectAttempts: 3,
  });
};

export default SSEService;