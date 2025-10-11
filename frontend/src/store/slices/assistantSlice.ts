import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';
import { v4 as uuidv4 } from 'uuid';
import apiService from '../../services/api';
import type { AssistantState, AssistantConversation } from '../../types/store';
import type { ChatRequest, ChatMessage } from '../../types/api';

// 异步thunks
export const initializeAssistant = createAsyncThunk(
  'assistant/initialize',
  async () => {
    const response = await apiService.initializeAssistant();
    return response;
  }
);

export const sendAssistantMessage = createAsyncThunk(
  'assistant/sendMessage',
  async (params: {
    conversationId: string;
    message: ChatMessage;
    model?: string;
    maxTokens?: number;
    temperature?: number;
    tools?: any[];
    useTools?: boolean;
    provider?: string;
    selectedTool?: string;
  }, { getState }) => {
    const { conversationId, message, model, maxTokens, temperature, tools, useTools, provider, selectedTool } = params;
    
    // 获取当前状态以获取完整的对话历史
    const state = getState() as any;
    const conversation = state.assistant.conversations?.find((c: any) => c.id === conversationId);
    
    // 构建包含完整对话历史的消息列表
    const messages: ChatMessage[] = [];
    if (conversation && conversation.messages) {
      // 添加历史消息
      messages.push(...conversation.messages);
    }
    // 添加当前用户消息
    messages.push(message);
    
    const request: ChatRequest = {
      messages,
      model,
      max_tokens: maxTokens,
      temperature,
      tools,
      use_tools: useTools,
      provider,
      selected_tool: selectedTool,
    };

    const response = await apiService.assistantChat(request);
    
    // 安全检查：确保response.choices存在且不为空
    if (!response.choices || response.choices.length === 0) {
      throw new Error('API响应格式错误：缺少choices数据');
    }
    
    return {
      conversationId,
      userMessage: message,
      assistantMessage: response.choices[0].message,
      usage: response.usage,
    };
  }
);

export const createAssistantConversation = createAsyncThunk(
  'assistant/createConversation',
  async (title: string) => {
    const conversation: AssistantConversation = {
      id: uuidv4(),
      title,
      messages: [],
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };
    return conversation;
  }
);

const initialState: AssistantState = {
  conversations: [],
  currentConversationId: null,
  isInitialized: false,
  isLoading: false,
  error: null,
};

const assistantSlice = createSlice({
  name: 'assistant',
  initialState,
  reducers: {
    setCurrentConversation: (state, action: PayloadAction<string>) => {
      state.currentConversationId = action.payload;
    },
    addMessage: (state, action: PayloadAction<{
      conversationId: string;
      message: ChatMessage;
    }>) => {
      const { conversationId, message } = action.payload;
      const conversation = (state.conversations || []).find(c => c.id === conversationId);
      if (conversation) {
        conversation.messages.push(message);
        conversation.updatedAt = new Date().toISOString();
      }
    },
    deleteConversation: (state, action: PayloadAction<string>) => {
      state.conversations = (state.conversations || []).filter(c => c.id !== action.payload);
      if (state.currentConversationId === action.payload) {
        state.currentConversationId = null;
      }
    },
    updateConversationTitle: (state, action: PayloadAction<{
      conversationId: string;
      title: string;
    }>) => {
      const { conversationId, title } = action.payload;
      const conversation = (state.conversations || []).find(c => c.id === conversationId);
      if (conversation) {
        conversation.title = title;
        conversation.updatedAt = new Date().toISOString();
      }
    },
    clearError: (state) => {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      // initializeAssistant
      .addCase(initializeAssistant.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(initializeAssistant.fulfilled, (state) => {
        state.isLoading = false;
        state.isInitialized = true;
      })
      .addCase(initializeAssistant.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || '助手初始化失败';
      })
      // createAssistantConversation
      .addCase(createAssistantConversation.fulfilled, (state, action) => {
        state.conversations.unshift(action.payload);
        state.currentConversationId = action.payload.id;
      })
      // sendAssistantMessage
      .addCase(sendAssistantMessage.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(sendAssistantMessage.fulfilled, (state, action) => {
        state.isLoading = false;
        const { conversationId, userMessage, assistantMessage } = action.payload;
        const conversation = (state.conversations || []).find(c => c.id === conversationId);
        if (conversation) {
          conversation.messages.push(userMessage, assistantMessage);
          conversation.updatedAt = new Date().toISOString();
        }
      })
      .addCase(sendAssistantMessage.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || '发送消息失败';
      });
  },
});

export const {
  setCurrentConversation,
  addMessage,
  deleteConversation,
  updateConversationTitle,
  clearError,
} = assistantSlice.actions;

export default assistantSlice.reducer;