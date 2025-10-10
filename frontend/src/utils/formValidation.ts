import type { Rule } from 'antd/es/form';

// 常用验证规则
export const commonRules = {
  required: (message?: string): Rule => ({
    required: true,
    message: message || '此字段为必填项',
  }),

  email: (message?: string): Rule => ({
    type: 'email',
    message: message || '请输入有效的邮箱地址',
  }),

  url: (message?: string): Rule => ({
    type: 'url',
    message: message || '请输入有效的URL地址',
  }),

  minLength: (min: number, message?: string): Rule => ({
    min,
    message: message || `长度不能少于${min}个字符`,
  }),

  maxLength: (max: number, message?: string): Rule => ({
    max,
    message: message || `长度不能超过${max}个字符`,
  }),

  pattern: (pattern: RegExp, message: string): Rule => ({
    pattern,
    message,
  }),

  number: (message?: string): Rule => ({
    type: 'number',
    message: message || '请输入有效的数字',
  }),

  integer: (message?: string): Rule => ({
    type: 'integer',
    message: message || '请输入有效的整数',
  }),

  range: (min: number, max: number, message?: string): Rule => ({
    type: 'number',
    min,
    max,
    message: message || `请输入${min}到${max}之间的数字`,
  }),
};

// API Key 验证规则
export const apiKeyRules = {
  openai: [
    commonRules.required('请输入OpenAI API密钥'),
    commonRules.minLength(10, 'OpenAI API密钥长度至少10个字符'),
    commonRules.pattern(/^sk-/, 'OpenAI API密钥应以"sk-"开头'),
  ],

  anthropic: [
    commonRules.required('请输入Anthropic API密钥'),
    commonRules.minLength(10, 'Anthropic API密钥长度至少10个字符'),
    commonRules.pattern(/^sk-ant-/, 'Anthropic API密钥应以"sk-ant-"开头'),
  ],

  google: [
    commonRules.required('请输入Google API密钥'),
    commonRules.minLength(10, 'Google API密钥长度至少10个字符'),
  ],

  azure: [
    commonRules.required('请输入Azure API密钥'),
    commonRules.minLength(10, 'Azure API密钥长度至少10个字符'),
  ],

  generic: [
    commonRules.required('请输入API密钥'),
    commonRules.minLength(10, 'API密钥长度至少10个字符'),
  ],
};

// 获取提供商对应的API Key验证规则
export function getAPIKeyRules(provider: string): Rule[] {
  const normalizedProvider = provider.toLowerCase();
  return apiKeyRules[normalizedProvider as keyof typeof apiKeyRules] || apiKeyRules.generic;
}

// 模型名称验证规则
export const modelNameRules = [
  commonRules.required('请输入模型名称'),
  commonRules.minLength(1, '模型名称不能为空'),
  commonRules.maxLength(100, '模型名称不能超过100个字符'),
  commonRules.pattern(/^[a-zA-Z0-9\-_.]+$/, '模型名称只能包含字母、数字、连字符、下划线和点'),
];

// 工具名称验证规则
export const toolNameRules = [
  commonRules.required('请输入工具名称'),
  commonRules.minLength(1, '工具名称不能为空'),
  commonRules.maxLength(50, '工具名称不能超过50个字符'),
  commonRules.pattern(/^[a-zA-Z0-9\-_]+$/, '工具名称只能包含字母、数字、连字符和下划线'),
];

// 端口号验证规则
export const portRules = [
  commonRules.required('请输入端口号'),
  commonRules.integer('端口号必须是整数'),
  commonRules.range(1, 65535, '端口号必须在1-65535之间'),
];

// URL验证规则
export const urlRules = [
  commonRules.required('请输入URL'),
  commonRules.url('请输入有效的URL地址'),
];

// 自定义验证器
export const customValidators = {
  // 确认密码验证
  confirmPassword: (getFieldValue: (name: string) => any) => ({
    validator(_: any, value: string) {
      if (!value || getFieldValue('password') === value) {
        return Promise.resolve();
      }
      return Promise.reject(new Error('两次输入的密码不一致'));
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
        return Promise.reject(new Error('请输入有效的JSON格式'));
      }
    },
  },

  // 数组非空验证
  arrayNotEmpty: (message?: string) => ({
    validator(_: any, value: any[]) {
      if (Array.isArray(value) && value.length > 0) {
        return Promise.resolve();
      }
      return Promise.reject(new Error(message || '请至少选择一项'));
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
        return Promise.reject(new Error(`文件大小不能超过${maxSize}${unit}`));
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