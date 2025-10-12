import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

// 导入翻译资源
import zh from './locales/zh';
import en from './locales/en';

const resources = {
  zh: {
    translation: zh,
  },
  en: {
    translation: en,
  },
};

i18n
  .use(initReactI18next)
  .init({
    resources,
    lng: 'en', // 默认语言
    fallbackLng: 'zh', // 回退语言
    interpolation: {
      escapeValue: false, // React已经默认转义了
    },
  });

export default i18n;