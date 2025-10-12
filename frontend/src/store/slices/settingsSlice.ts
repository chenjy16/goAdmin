import type { SettingsState } from '../../types/store';
import { createConfigSlice, createLocalStorageManager } from '../utils/configManager';

// 默认设置配置
const defaultSettings: SettingsState = {
  theme: 'light',
  language: 'en',
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

// 设置验证器
const validateSettings = (config: any): config is SettingsState => {
  return (
    config &&
    typeof config === 'object' &&
    ['light', 'dark'].includes(config.theme) &&
    ['zh', 'en'].includes(config.language) &&
    typeof config.apiKeys === 'object' &&
    typeof config.defaultProvider === 'string' &&
    typeof config.defaultModel === 'object' &&
    typeof config.chatSettings === 'object' &&
    typeof config.chatSettings.maxTokens === 'number' &&
    typeof config.chatSettings.temperature === 'number' &&
    typeof config.chatSettings.streamResponse === 'boolean'
  );
};

// 创建本地存储管理器
const settingsManager = createLocalStorageManager(
  'app-settings',
  defaultSettings,
  validateSettings
);

// 创建设置slice
const settingsConfig = createConfigSlice('settings', settingsManager);

// 导出actions和reducer
export const {
  updateConfig: updateSettings,
  setConfig: setSettings,
  resetConfig: resetSettings,
  markSaved,
  clearError,
  setLoading,
  loadConfig: loadSettings,
  saveConfig: saveSettings,
} = settingsConfig.actions;

// 便捷的设置更新actions
export const setTheme = (theme: 'light' | 'dark') => 
  updateSettings({ theme });

export const setLanguage = (language: 'zh' | 'en') => 
  updateSettings({ language });

export const setAPIKey = (provider: string, apiKey: string) => 
  updateSettings({ 
    apiKeys: { 
      ...settingsManager.load().apiKeys, 
      [provider]: apiKey 
    } 
  });

export const removeAPIKey = (provider: string) => {
  const currentSettings = settingsManager.load();
  const { [provider]: removed, ...remainingKeys } = currentSettings.apiKeys;
  return updateSettings({ apiKeys: remainingKeys });
};

export const setDefaultProvider = (provider: string) => 
  updateSettings({ defaultProvider: provider });

export const setDefaultModel = (provider: string, model: string) => 
  updateSettings({ 
    defaultModel: { 
      ...settingsManager.load().defaultModel, 
      [provider]: model 
    } 
  });

export const updateChatSettings = (chatSettings: Partial<SettingsState['chatSettings']>) => 
  updateSettings({ 
    chatSettings: { 
      ...settingsManager.load().chatSettings, 
      ...chatSettings 
    } 
  });

// 导出selectors
export const {
  selectConfig: selectSettings,
  selectLoading: selectSettingsLoading,
  selectError: selectSettingsError,
  selectInitialized: selectSettingsInitialized,
  selectIsDirty: selectSettingsIsDirty,
} = settingsConfig.selectors;

export default settingsConfig.reducer;