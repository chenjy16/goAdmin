import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';
import { v4 as uuidv4 } from 'uuid';
import apiService from '../../services/api';
import type { ChatState, Conversation } from '../../types/store';
import type { UnifiedChatRequest, UnifiedMessage } from '../../types/api';

// 异步thunks
export const sendMessage = createAsyncThunk(
  'chat/sendMessage',
  async (params: {
    conversationId: string;
    message: UnifiedMessage;
    provider: string;
    model: string;
    maxTokens?: number;
    temperature?: number;
  }) => {
    const { conversationId, message, provider, model, maxTokens, temperature } = params;
    
    const request: UnifiedChatRequest = {
      messages: [message], // 这里应该包含完整的对话历史
      model,
      max_tokens: maxTokens,
      temperature,
      stream: false,
    };

    const response = await apiService.chatCompletion(provider, request);
    return {
      conversationId,
      userMessage: message,
      assistantMessage: response.choices[0].message,
      usage: response.usage,
    };
  }
);

export const createConversation = createAsyncThunk(
  'chat/createConversation',
  async (params: {
    title: string;
    provider: string;
    model: string;
  }) => {
    const { title, provider, model } = params;
    const conversation: Conversation = {
      id: uuidv4(),
      title,
      messages: [],
      provider,
      model,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };
    return conversation;
  }
);

const initialState: ChatState = {
  conversations: [],
  currentConversationId: null,
  isLoading: false,
  error: null,
  streamingMessage: '',
};

const chatSlice = createSlice({
  name: 'chat',
  initialState,
  reducers: {
    setCurrentConversation: (state, action: PayloadAction<string>) => {
      state.currentConversationId = action.payload;
    },
    addMessage: (state, action: PayloadAction<{
      conversationId: string;
      message: UnifiedMessage;
    }>) => {
      const { conversationId, message } = action.payload;
      const conversation = (state.conversations || []).find(c => c.id === conversationId);
      if (conversation) {
        conversation.messages.push(message);
        conversation.updatedAt = new Date().toISOString();
      }
    },
    updateStreamingMessage: (state, action: PayloadAction<string>) => {
      state.streamingMessage = action.payload;
    },
    clearStreamingMessage: (state) => {
      state.streamingMessage = '';
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
      // createConversation
      .addCase(createConversation.fulfilled, (state, action) => {
        state.conversations.unshift(action.payload);
        state.currentConversationId = action.payload.id;
      })
      // sendMessage
      .addCase(sendMessage.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(sendMessage.fulfilled, (state, action) => {
        state.isLoading = false;
        const { conversationId, userMessage, assistantMessage } = action.payload;
        const conversation = (state.conversations || []).find(c => c.id === conversationId);
        if (conversation) {
          conversation.messages.push(userMessage, assistantMessage);
          conversation.updatedAt = new Date().toISOString();
        }
      })
      .addCase(sendMessage.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || '发送消息失败';
      });
  },
});

export const {
  setCurrentConversation,
  addMessage,
  updateStreamingMessage,
  clearStreamingMessage,
  deleteConversation,
  updateConversationTitle,
  clearError,
} = chatSlice.actions;

export default chatSlice.reducer;