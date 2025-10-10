import { BaseService } from './base/BaseService';
import type { 
  BaseApiResponse, 
  IInitializable, 
  IValidatable, 
  ValidationResult 
} from '../types/base';
import type { 
  ChatRequest, 
  ChatResponse, 
  ChatMessage 
} from '../types/api';

/**
 * Assistant服务类
 * 负责处理AI助手相关的API调用
 */
export class AssistantService extends BaseService implements IInitializable, IValidatable<ChatRequest> {
  private _initialized = false;

  /**
   * 初始化Assistant系统
   */
  async initialize(): Promise<void> {
    try {
      const response = await this.post<BaseApiResponse>('/api/v1/assistant/initialize');
      if (response.code === 200) {
        this._initialized = true;
      } else {
        throw new Error(response.message || 'Failed to initialize Assistant');
      }
    } catch (error) {
      this._initialized = false;
      throw error;
    }
  }

  /**
   * 检查是否已初始化
   */
  isInitialized(): boolean {
    return this._initialized;
  }

  /**
   * 发送聊天消息
   */
  async chat(request: ChatRequest): Promise<ChatResponse> {
    const validationResult = this.validate(request);
    if (!validationResult.isValid) {
      throw new Error(`Validation failed: ${validationResult.errors.map(e => e.message).join(', ')}`);
    }

    return this.post<ChatResponse>('/api/v1/assistant/chat', request);
  }

  /**
   * 创建聊天流（SSE）
   */
  async createChatStream(request: ChatRequest): Promise<ReadableStream<ChatResponse>> {
    const validationResult = this.validate(request);
    if (!validationResult.isValid) {
      throw new Error(`Validation failed: ${validationResult.errors.map(e => e.message).join(', ')}`);
    }

    const response = await fetch(`${this.baseURL}/api/v1/assistant/chat/stream`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    if (!response.body) {
      throw new Error('Response body is null');
    }

    return response.body
      .pipeThrough(new TextDecoderStream())
      .pipeThrough(this.createSSETransformer());
  }

  /**
   * 创建SSE转换器
   */
  private createSSETransformer(): TransformStream<string, ChatResponse> {
    return new TransformStream({
      transform(chunk, controller) {
        const lines = chunk.split('\n');
        for (const line of lines) {
          if (line.startsWith('data: ')) {
            const data = line.slice(6);
            if (data === '[DONE]') {
              controller.terminate();
              return;
            }
            try {
              const parsed = JSON.parse(data);
              controller.enqueue(parsed);
            } catch (error) {
              console.warn('Failed to parse SSE data:', data);
            }
          }
        }
      }
    });
  }

  /**
   * 获取聊天历史
   */
  async getChatHistory(conversationId?: string): Promise<BaseApiResponse<ChatMessage[]>> {
    const url = conversationId 
      ? `/api/v1/assistant/conversations/${conversationId}/messages`
      : '/api/v1/assistant/messages';
    return this.get<BaseApiResponse<ChatMessage[]>>(url);
  }

  /**
   * 创建新对话
   */
  async createConversation(title?: string): Promise<BaseApiResponse<{ id: string; title: string }>> {
    return this.post<BaseApiResponse<{ id: string; title: string }>>('/api/v1/assistant/conversations', {
      title: title || `Conversation ${new Date().toLocaleString()}`
    });
  }

  /**
   * 获取对话列表
   */
  async getConversations(): Promise<BaseApiResponse<Array<{ id: string; title: string; createdAt: string; updatedAt: string }>>> {
    return this.get<BaseApiResponse<Array<{ id: string; title: string; createdAt: string; updatedAt: string }>>>('/api/v1/assistant/conversations');
  }

  /**
   * 删除对话
   */
  async deleteConversation(conversationId: string): Promise<BaseApiResponse> {
    return this.delete<BaseApiResponse>(`/api/v1/assistant/conversations/${conversationId}`);
  }

  /**
   * 更新对话标题
   */
  async updateConversationTitle(conversationId: string, title: string): Promise<BaseApiResponse> {
    return this.put<BaseApiResponse>(`/api/v1/assistant/conversations/${conversationId}`, { title });
  }

  /**
   * 验证聊天请求
   */
  validate(request: ChatRequest): ValidationResult {
    const errors = [];

    if (!request.messages || !Array.isArray(request.messages) || request.messages.length === 0) {
      errors.push({
        field: 'messages',
        message: 'Messages array is required and cannot be empty',
        code: 'REQUIRED'
      });
    }

    // 验证消息格式
    if (request.messages) {
      request.messages.forEach((message, index) => {
        if (!message.role || !['user', 'assistant', 'system'].includes(message.role)) {
          errors.push({
            field: `messages[${index}].role`,
            message: 'Message role must be one of: user, assistant, system',
            code: 'INVALID_VALUE'
          });
        }

        if (!message.content || typeof message.content !== 'string' || message.content.trim().length === 0) {
          errors.push({
            field: `messages[${index}].content`,
            message: 'Message content is required and must be a non-empty string',
            code: 'REQUIRED'
          });
        }

        // 检查内容长度
        if (message.content && message.content.length > 100000) {
          errors.push({
            field: `messages[${index}].content`,
            message: 'Message content is too long (max 100,000 characters)',
            code: 'MAX_LENGTH'
          });
        }
      });
    }

    // 验证可选参数
    if (request.max_tokens !== undefined) {
      if (typeof request.max_tokens !== 'number' || request.max_tokens < 1 || request.max_tokens > 100000) {
        errors.push({
          field: 'max_tokens',
          message: 'max_tokens must be a number between 1 and 100,000',
          code: 'INVALID_RANGE'
        });
      }
    }

    if (request.temperature !== undefined) {
      if (typeof request.temperature !== 'number' || request.temperature < 0 || request.temperature > 2) {
        errors.push({
          field: 'temperature',
          message: 'temperature must be a number between 0 and 2',
          code: 'INVALID_RANGE'
        });
      }
    }

    return {
      isValid: errors.length === 0,
      errors
    };
  }

  /**
   * 验证消息内容
   */
  validateMessage(message: ChatMessage): ValidationResult {
    const errors = [];

    if (!message.role || !['user', 'assistant', 'system'].includes(message.role)) {
      errors.push({
        field: 'role',
        message: 'Message role must be one of: user, assistant, system',
        code: 'INVALID_VALUE'
      });
    }

    if (!message.content || typeof message.content !== 'string' || message.content.trim().length === 0) {
      errors.push({
        field: 'content',
        message: 'Message content is required and must be a non-empty string',
        code: 'REQUIRED'
      });
    }

    if (message.content && message.content.length > 100000) {
      errors.push({
        field: 'content',
        message: 'Message content is too long (max 100,000 characters)',
        code: 'MAX_LENGTH'
      });
    }

    return {
      isValid: errors.length === 0,
      errors
    };
  }

  /**
   * 格式化消息用于显示
   */
  formatMessage(message: ChatMessage): string {
    const timestamp = message.timestamp ? new Date(message.timestamp).toLocaleString() : '';
    const role = message.role.charAt(0).toUpperCase() + message.role.slice(1);
    
    return `[${timestamp}] ${role}: ${message.content}`;
  }

  /**
   * 计算对话token数量（估算）
   */
  estimateTokenCount(messages: ChatMessage[]): number {
    // 简单的token估算：平均每4个字符为1个token
    const totalChars = messages.reduce((sum, msg) => sum + msg.content.length, 0);
    return Math.ceil(totalChars / 4);
  }

  /**
   * 获取系统状态
   */
  async getStatus(): Promise<BaseApiResponse<{ initialized: boolean; model?: string; provider?: string }>> {
    return this.get<BaseApiResponse<{ initialized: boolean; model?: string; provider?: string }>>('/api/v1/assistant/status');
  }
}

export const assistantService = new AssistantService();
export default assistantService;