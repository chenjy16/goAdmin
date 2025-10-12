import type { Rule } from 'antd/es/form';
import { useTranslation } from 'react-i18next';
import { useMemo } from 'react';

// 创建国际化的常用验证规则
export function createCommonRules(t: (key: string, options?: any) => string) {
  return {
    required: (message?: string): Rule => ({
      required: true,
      message: message || t('formValidation.required'),
    }),

    email: (message?: string): Rule => ({
      type: 'email',
      message: message || t('formValidation.email'),
    }),

    url: (message?: string): Rule => ({
      type: 'url',
      message: message || t('formValidation.url'),
    }),

    minLength: (min: number, message?: string): Rule => ({
      min,
      message: message || t('formValidation.minLength', { min }),
    }),

    maxLength: (max: number, message?: string): Rule => ({
      max,
      message: message || t('formValidation.maxLength', { max }),
    }),

    pattern: (pattern: RegExp, message: string): Rule => ({
      pattern,
      message,
    }),

    number: (message?: string): Rule => ({
      type: 'number',
      message: message || t('formValidation.number'),
    }),

    integer: (message?: string): Rule => ({
      type: 'integer',
      message: message || t('formValidation.integer'),
    }),

    range: (min: number, max: number, message?: string): Rule => ({
      type: 'number',
      min,
      max,
      message: message || t('formValidation.range', { min, max }),
    }),
  };
}

// 保持向后兼容的常用验证规则（使用默认英文消息）
export const commonRules = {
  required: (message?: string): Rule => ({
    required: true,
    message: message || 'This field is required',
  }),

  email: (message?: string): Rule => ({
    type: 'email',
    message: message || 'Please enter a valid email address',
  }),

  url: (message?: string): Rule => ({
    type: 'url',
    message: message || 'Please enter a valid URL',
  }),

  minLength: (min: number, message?: string): Rule => ({
    min,
    message: message || `Minimum length is ${min} characters`,
  }),

  maxLength: (max: number, message?: string): Rule => ({
    max,
    message: message || `Maximum length is ${max} characters`,
  }),

  pattern: (pattern: RegExp, message: string): Rule => ({
    pattern,
    message,
  }),

  number: (message?: string): Rule => ({
    type: 'number',
    message: message || 'Please enter a valid number',
  }),

  integer: (message?: string): Rule => ({
    type: 'integer',
    message: message || 'Please enter a valid integer',
  }),

  range: (min: number, max: number, message?: string): Rule => ({
    type: 'number',
    min,
    max,
    message: message || `Please enter a number between ${min} and ${max}`,
  }),
};

// 创建国际化的 API Key 验证规则
export function createApiKeyRules(t: (key: string, options?: any) => string) {
  const commonRules = createCommonRules(t);
  
  return {
    openai: [
      commonRules.required(t('formValidation.apiKey.openai.required')),
      commonRules.minLength(10, t('formValidation.apiKey.openai.minLength')),
      commonRules.pattern(/^sk-/, t('formValidation.apiKey.openai.pattern')),
    ],

    anthropic: [
      commonRules.required(t('formValidation.apiKey.anthropic.required')),
      commonRules.minLength(10, t('formValidation.apiKey.anthropic.minLength')),
      commonRules.pattern(/^sk-ant-/, t('formValidation.apiKey.anthropic.pattern')),
    ],

    google: [
      commonRules.required(t('formValidation.apiKey.google.required')),
      commonRules.minLength(10, t('formValidation.apiKey.google.minLength')),
    ],

    azure: [
      commonRules.required(t('formValidation.apiKey.azure.required')),
      commonRules.minLength(10, t('formValidation.apiKey.azure.minLength')),
    ],

    generic: [
      commonRules.required(t('formValidation.apiKey.generic.required')),
      commonRules.minLength(10, t('formValidation.apiKey.generic.minLength')),
    ],
  };
}

// 保持向后兼容的 API Key 验证规则（使用默认中文消息）
export const apiKeyRules = {
  openai: [
    commonRules.required('Please enter OpenAI API key'),
    commonRules.minLength(10, 'OpenAI API key must be at least 10 characters'),
    commonRules.pattern(/^sk-/, 'OpenAI API key should start with "sk-"'),
  ],

  anthropic: [
    commonRules.required('Please enter Anthropic API key'),
    commonRules.minLength(10, 'Anthropic API key must be at least 10 characters'),
    commonRules.pattern(/^sk-ant-/, 'Anthropic API key should start with "sk-ant-"'),
  ],

  google: [
    commonRules.required('Please enter Google API key'),
    commonRules.minLength(10, 'Google API key must be at least 10 characters'),
  ],

  azure: [
    commonRules.required('Please enter Azure API key'),
    commonRules.minLength(10, 'Azure API key must be at least 10 characters'),
  ],

  generic: [
    commonRules.required('Please enter API key'),
    commonRules.minLength(10, 'API key must be at least 10 characters'),
  ],
};

// 创建国际化的获取提供商对应的API Key验证规则函数
export function createGetAPIKeyRules(t: (key: string, options?: any) => string) {
  return (provider: string): Rule[] => {
    const apiKeyRules = createApiKeyRules(t);
    const normalizedProvider = provider.toLowerCase();
    return apiKeyRules[normalizedProvider as keyof typeof apiKeyRules] || apiKeyRules.generic;
  };
}

// 创建国际化的模型名称验证规则
export function createModelNameRules(t: (key: string, options?: any) => string): Rule[] {
  const commonRules = createCommonRules(t);
  return [
    commonRules.required(t('formValidation.modelName.required')),
    commonRules.minLength(1, t('formValidation.modelName.minLength')),
    commonRules.maxLength(100, t('formValidation.modelName.maxLength')),
    commonRules.pattern(/^[a-zA-Z0-9\-_.]+$/, t('formValidation.modelName.pattern')),
  ];
}

// 创建国际化的工具名称验证规则
export function createToolNameRules(t: (key: string, options?: any) => string): Rule[] {
  const commonRules = createCommonRules(t);
  return [
    commonRules.required(t('formValidation.toolName.required')),
    commonRules.minLength(1, t('formValidation.toolName.minLength')),
    commonRules.maxLength(50, t('formValidation.toolName.maxLength')),
    commonRules.pattern(/^[a-zA-Z0-9\-_]+$/, t('formValidation.toolName.pattern')),
  ];
}

// 创建国际化的端口号验证规则
export function createPortRules(t: (key: string, options?: any) => string): Rule[] {
  const commonRules = createCommonRules(t);
  return [
    commonRules.required(t('formValidation.port.required')),
    commonRules.integer(t('formValidation.port.integer')),
    commonRules.range(1, 65535, t('formValidation.port.range')),
  ];
}

// 创建国际化的URL验证规则
export function createUrlRules(t: (key: string, options?: any) => string): Rule[] {
  const commonRules = createCommonRules(t);
  return [
    commonRules.required(t('formValidation.required')),
    commonRules.url(t('formValidation.url')),
  ];
}

// 保持向后兼容的函数和规则（使用默认中文消息）
export function getAPIKeyRules(provider: string): Rule[] {
  const normalizedProvider = provider.toLowerCase();
  return apiKeyRules[normalizedProvider as keyof typeof apiKeyRules] || apiKeyRules.generic;
}

export const modelNameRules = [
  commonRules.required('Please enter model name'),
  commonRules.minLength(1, 'Model name cannot be empty'),
  commonRules.maxLength(100, 'Model name cannot exceed 100 characters'),
  commonRules.pattern(/^[a-zA-Z0-9\-_.]+$/, 'Model name can only contain letters, numbers, hyphens, underscores and dots'),
];

export const toolNameRules = [
  commonRules.required('Please enter tool name'),
  commonRules.minLength(1, 'Tool name cannot be empty'),
  commonRules.maxLength(50, 'Tool name cannot exceed 50 characters'),
  commonRules.pattern(/^[a-zA-Z0-9\-_]+$/, 'Tool name can only contain letters, numbers, hyphens and underscores'),
];

export const portRules = [
  commonRules.required('Please enter port number'),
  commonRules.integer('Port number must be an integer'),
  commonRules.range(1, 65535, 'Port number must be between 1-65535'),
];

export const urlRules = [
  commonRules.required('Please enter URL'),
  commonRules.url('Please enter a valid URL'),
];

// 创建国际化的自定义验证器
export function createCustomValidators(t: (key: string, options?: any) => string) {
  return {
    // 确认密码验证
    confirmPassword: (getFieldValue: (name: string) => any) => ({
      validator(_: any, value: string) {
        if (!value || getFieldValue('password') === value) {
          return Promise.resolve();
        }
        return Promise.reject(new Error(t('formValidation.confirmPassword')));
      },
    }),

    // JSON格式验证
    json: {
      validator(_: any, value: string) {
        if (!value) return Promise.resolve();
        try {
          JSON.parse(value);
          return Promise.resolve();
        } catch {
          return Promise.reject(new Error(t('formValidation.json')));
        }
      },
    },

    // 数组非空验证
    arrayNotEmpty: (message?: string) => ({
      validator(_: any, value: any[]) {
        if (Array.isArray(value) && value.length > 0) {
          return Promise.resolve();
        }
        return Promise.reject(new Error(message || t('formValidation.arrayNotEmpty')));
      },
    }),

    // 文件大小验证
    fileSize: (maxSize: number, unit: 'KB' | 'MB' | 'GB' = 'MB') => {
      const multiplier = { KB: 1024, MB: 1024 * 1024, GB: 1024 * 1024 * 1024 }[unit];
      const maxBytes = maxSize * multiplier;
      
      return {
        validator(_: any, value: File) {
          if (!value || value.size <= maxBytes) {
            return Promise.resolve();
          }
          return Promise.reject(new Error(t('formValidation.fileSize', { maxSize, unit })));
        },
      };
    },
  };
}

// 保持向后兼容的自定义验证器（使用默认英文消息）
export const customValidators = {
  // 确认密码验证
  confirmPassword: (getFieldValue: (name: string) => any) => ({
    validator(_: any, value: string) {
      if (!value || getFieldValue('password') === value) {
        return Promise.resolve();
      }
      return Promise.reject(new Error('Passwords do not match'));
    },
  }),

  // JSON格式验证
  json: {
    validator(_: any, value: string) {
      if (!value) return Promise.resolve();
      try {
        JSON.parse(value);
        return Promise.resolve();
      } catch {
        return Promise.reject(new Error('Please enter valid JSON format'));
      }
    },
  },

  // 数组非空验证
  arrayNotEmpty: (message?: string) => ({
    validator(_: any, value: any[]) {
      if (Array.isArray(value) && value.length > 0) {
        return Promise.resolve();
      }
      return Promise.reject(new Error(message || 'Please select at least one item'));
    },
  }),

  // 文件大小验证
  fileSize: (maxSize: number, unit: 'KB' | 'MB' | 'GB' = 'MB') => {
    const multiplier = { KB: 1024, MB: 1024 * 1024, GB: 1024 * 1024 * 1024 }[unit];
    const maxBytes = maxSize * multiplier;
    
    return {
      validator(_: any, value: File) {
        if (!value || value.size <= maxBytes) {
          return Promise.resolve();
        }
        return Promise.reject(new Error(`File size cannot exceed ${maxSize}${unit}`));
      },
    };
  },
};

// 组合验证规则
export function combineRules(...ruleGroups: Rule[][]): Rule[] {
  return ruleGroups.flat();
}

// 条件验证规则
export function conditionalRules(condition: boolean, rules: Rule[]): Rule[] {
  return condition ? rules : [];
}

// 国际化表单验证 Hook
export function useFormValidation() {
  const { t } = useTranslation();

  const commonRules = useMemo(() => createCommonRules(t), [t]);
  const apiKeyRules = useMemo(() => createApiKeyRules(t), [t]);
  const customValidators = useMemo(() => createCustomValidators(t), [t]);
  
  const getAPIKeyRules = useMemo(() => createGetAPIKeyRules(t), [t]);
  const modelNameRules = useMemo(() => createModelNameRules(t), [t]);
  const toolNameRules = useMemo(() => createToolNameRules(t), [t]);
  const portRules = useMemo(() => createPortRules(t), [t]);
  const urlRules = useMemo(() => createUrlRules(t), [t]);

  return {
    commonRules,
    apiKeyRules,
    customValidators,
    getAPIKeyRules,
    modelNameRules,
    toolNameRules,
    portRules,
    urlRules,
    // 工具函数
    combineRules,
    conditionalRules,
  };
}