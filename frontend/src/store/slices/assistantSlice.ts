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
  }) => {
    const { conversationId, message, model, maxTokens, temperature, tools } = params;
    
    const request: ChatRequest = {
      messages: [message], // 这里应该包含完整的对话历史
      model,
      max_tokens: maxTokens,
      temperature,
      tools,
    };

    const response = await apiService.assistantChat(request);
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