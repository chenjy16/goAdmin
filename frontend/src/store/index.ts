import { configureStore } from '@reduxjs/toolkit';
import { useDispatch, useSelector } from 'react-redux';
import type { TypedUseSelectorHook } from 'react-redux';

import chatReducer from './slices/chatSlice';
import providersReducer from './slices/providersSlice';
import mcpReducer from './slices/mcpSlice';
import assistantReducer from './slices/assistantSlice';
import settingsReducer from './slices/settingsSlice';
import uiReducer from './slices/uiSlice';

import type { RootState } from '../types/store';

export const store = configureStore({
  reducer: {
    chat: chatReducer,
    providers: providersReducer,
    mcp: mcpReducer,
    assistant: assistantReducer,
    settings: settingsReducer,
    ui: uiReducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        ignoredActions: ['persist/PERSIST', 'persist/REHYDRATE'],
      },
    }),
  devTools: import.meta.env.DEV,
});

export type AppDispatch = typeof store.dispatch;

// 类型化的hooks
export const useAppDispatch: () => AppDispatch = useDispatch;
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;

export default store;