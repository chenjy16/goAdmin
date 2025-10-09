import { createSlice } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';
import type { SettingsState } from '../../types/store';

const initialState: SettingsState = {
  theme: 'light',
  language: 'zh',
  apiKeys: {},
  defaultProvider: 'openai',
  defaultModel: {
    openai: 'gpt-3.5-turbo',
    googleai: 'gemini-1.5-flash',
  },
  chatSettings: {
    maxTokens: 2048,
    temperature: 0.7,
    streamResponse: true,
  },
};

const settingsSlice = createSlice({
  name: 'settings',
  initialState,
  reducers: {
    setTheme: (state, action: PayloadAction<'light' | 'dark'>) => {
      state.theme = action.payload;
    },
    setLanguage: (state, action: PayloadAction<'zh' | 'en'>) => {
      state.language = action.payload;
    },
    setAPIKey: (state, action: PayloadAction<{
      provider: string;
      apiKey: string;
    }>) => {
      const { provider, apiKey } = action.payload;
      state.apiKeys[provider] = apiKey;
    },
    removeAPIKey: (state, action: PayloadAction<string>) => {
      delete state.apiKeys[action.payload];
    },
    setDefaultProvider: (state, action: PayloadAction<string>) => {
      state.defaultProvider = action.payload;
    },
    setDefaultModel: (state, action: PayloadAction<{
      provider: string;
      model: string;
    }>) => {
      const { provider, model } = action.payload;
      state.defaultModel[provider] = model;
    },
    updateChatSettings: (state, action: PayloadAction<Partial<SettingsState['chatSettings']>>) => {
      state.chatSettings = { ...state.chatSettings, ...action.payload };
    },
    resetSettings: () => initialState,
    loadSettings: (state, action: PayloadAction<Partial<SettingsState>>) => {
      return { ...state, ...action.payload };
    },
  },
});

export const {
  setTheme,
  setLanguage,
  setAPIKey,
  removeAPIKey,
  setDefaultProvider,
  setDefaultModel,
  updateChatSettings,
  resetSettings,
  loadSettings,
} = settingsSlice.actions;

export default settingsSlice.reducer;