import { createSlice } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';
import { v4 as uuidv4 } from 'uuid';
import type { UIState, Notification } from '../../types/store';

const initialState: UIState = {
  sidebarCollapsed: false,
  currentPage: 'chat',
  loading: {},
  notifications: [],
};

const uiSlice = createSlice({
  name: 'ui',
  initialState,
  reducers: {
    toggleSidebar: (state) => {
      state.sidebarCollapsed = !state.sidebarCollapsed;
    },
    setSidebarCollapsed: (state, action: PayloadAction<boolean>) => {
      state.sidebarCollapsed = action.payload;
    },
    setCurrentPage: (state, action: PayloadAction<string>) => {
      state.currentPage = action.payload;
    },
    setLoading: (state, action: PayloadAction<{
      key: string;
      loading: boolean;
    }>) => {
      const { key, loading } = action.payload;
      if (loading) {
        state.loading[key] = true;
      } else {
        delete state.loading[key];
      }
    },
    addNotification: (state, action: PayloadAction<Omit<Notification, 'id' | 'timestamp'>>) => {
      const notification: Notification = {
        ...action.payload,
        id: uuidv4(),
        timestamp: new Date().toISOString(),
      };
      state.notifications.unshift(notification);
      
      // 保持最新的50条通知
      if ((state.notifications || []).length > 50) {
        state.notifications = (state.notifications || []).slice(0, 50);
      }
    },
    removeNotification: (state, action: PayloadAction<string>) => {
      state.notifications = (state.notifications || []).filter(n => n.id !== action.payload);
    },
    clearNotifications: (state) => {
      state.notifications = [];
    },
    markNotificationAsRead: (state, action: PayloadAction<string>) => {
      const notification = (state.notifications || []).find(n => n.id === action.payload);
      if (notification) {
        // 可以添加已读状态字段
      }
    },
  },
});

export const {
  toggleSidebar,
  setSidebarCollapsed,
  setCurrentPage,
  setLoading,
  addNotification,
  removeNotification,
  clearNotifications,
  markNotificationAsRead,
} = uiSlice.actions;

export default uiSlice.reducer;