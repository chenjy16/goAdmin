import { useState, useCallback, useRef, useEffect } from 'react';
import { Form } from 'antd';
import type { FormInstance } from 'antd';
import { useAsyncOperation, type UseAsyncDataConfig } from './useAsyncData';

// 表单验证规则
export interface ValidationRule {
  required?: boolean;
  message?: string;
  pattern?: RegExp;
  min?: number;
  max?: number;
  validator?: (value: any) => Promise<void> | void;
}

// 表单字段配置
export interface FormFieldConfig {
  name: string;
  label?: string;
  rules?: ValidationRule[];
  initialValue?: any;
  dependencies?: string[];
}

// 表单配置
export interface UseFormConfig<T> {
  initialValues?: Partial<T>;
  fields?: FormFieldConfig[];
  onSubmit?: (values: T) => Promise<any> | any;
  onValuesChange?: (changedValues: Partial<T>, allValues: T) => void;
  validateOnChange?: boolean;
  resetOnSubmit?: boolean;
  errorConfig?: UseAsyncDataConfig<any>['errorConfig'];
}

// 表单状态
export interface FormState<T> {
  values: Partial<T>;
  errors: Record<string, string>;
  touched: Record<string, boolean>;
  dirty: boolean;
  submitting: boolean;
  submitted: boolean;
}

// 表单hook
export function useForm<T extends Record<string, any>>(config: UseFormConfig<T> = {}) {
  const {
    initialValues = {},
    fields = [],
    onSubmit,
    onValuesChange,
    validateOnChange = false,
    resetOnSubmit = false,
    errorConfig = {},
  } = config;

  const [form] = Form.useForm();
  const formRef = useRef<FormInstance>(form);

  const [formState, setFormState] = useState<FormState<T>>({
    values: initialValues,
    errors: {},
    touched: {},
    dirty: false,
    submitting: false,
    submitted: false,
  });

  // 异步提交操作
  const submitOperation = useAsyncOperation(
    async (values: T) => {
      if (onSubmit) {
        return await onSubmit(values);
      }
    },
    {
      errorConfig,
      onSuccess: () => {
        setFormState(prev => ({
          ...prev,
          submitted: true,
          submitting: false,
        }));
        
        if (resetOnSubmit) {
          reset();
        }
      },
      onError: () => {
        setFormState(prev => ({
          ...prev,
          submitting: false,
        }));
      },
    }
  );

  // 验证字段
  const validateField = useCallback(async (name: string, value: any): Promise<string | null> => {
    const fieldConfig = fields.find(field => field.name === name);
    if (!fieldConfig?.rules) return null;

    for (const rule of fieldConfig.rules) {
      // 必填验证
      if (rule.required && (value === undefined || value === null || value === '')) {
        return rule.message || `${fieldConfig.label || name}是必填项`;
      }

      // 正则验证
      if (rule.pattern && value && !rule.pattern.test(value)) {
        return rule.message || `${fieldConfig.label || name}格式不正确`;
      }

      // 最小值/长度验证
      if (rule.min !== undefined) {
        if (typeof value === 'string' && value.length < rule.min) {
          return rule.message || `${fieldConfig.label || name}最少${rule.min}个字符`;
        }
        if (typeof value === 'number' && value < rule.min) {
          return rule.message || `${fieldConfig.label || name}最小值为${rule.min}`;
        }
      }

      // 最大值/长度验证
      if (rule.max !== undefined) {
        if (typeof value === 'string' && value.length > rule.max) {
          return rule.message || `${fieldConfig.label || name}最多${rule.max}个字符`;
        }
        if (typeof value === 'number' && value > rule.max) {
          return rule.message || `${fieldConfig.label || name}最大值为${rule.max}`;
        }
      }

      // 自定义验证
      if (rule.validator) {
        try {
          await rule.validator(value);
        } catch (error) {
          return error instanceof Error ? error.message : rule.message || 'Validation failed';
        }
      }
    }

    return null;
  }, [fields]);

  // 验证所有字段
  const validateAll = useCallback(async (): Promise<boolean> => {
    const values = form.getFieldsValue();
    const errors: Record<string, string> = {};

    for (const field of fields) {
      const error = await validateField(field.name, values[field.name]);
      if (error) {
        errors[field.name] = error;
      }
    }

    setFormState(prev => ({
      ...prev,
      errors,
    }));

    return Object.keys(errors).length === 0;
  }, [fields, validateField, form]);

  // 设置字段值
  const setFieldValue = useCallback((name: string, value: any) => {
    form.setFieldValue(name, value);
    
    setFormState(prev => ({
      ...prev,
      values: {
        ...prev.values,
        [name]: value,
      },
      touched: {
        ...prev.touched,
        [name]: true,
      },
      dirty: true,
    }));

    // 实时验证
    if (validateOnChange) {
      validateField(name, value).then(error => {
        setFormState(prev => ({
          ...prev,
          errors: {
            ...prev.errors,
            [name]: error || '',
          },
        }));
      });
    }

    // 触发值变化回调
    if (onValuesChange) {
      const allValues = { ...formState.values, [name]: value };
      onValuesChange({ [name]: value } as Partial<T>, allValues as T);
    }
  }, [form, formState.values, validateOnChange, validateField, onValuesChange]);

  // 设置多个字段值
  const setFieldsValue = useCallback((values: Partial<T>) => {
    form.setFieldsValue(values);
    
    setFormState(prev => ({
      ...prev,
      values: {
        ...prev.values,
        ...values,
      },
      dirty: true,
    }));

    if (onValuesChange) {
      onValuesChange(values, { ...formState.values, ...values } as T);
    }
  }, [form, formState.values, onValuesChange]);

  // 获取字段值
  const getFieldValue = useCallback((name: string) => {
    return form.getFieldValue(name);
  }, [form]);

  // 获取所有字段值
  const getFieldsValue = useCallback(() => {
    return form.getFieldsValue() as T;
  }, [form]);

  // 重置表单
  const reset = useCallback(() => {
    form.resetFields();
    setFormState({
      values: initialValues,
      errors: {},
      touched: {},
      dirty: false,
      submitting: false,
      submitted: false,
    });
  }, [form, initialValues]);

  // 提交表单
  const submit = useCallback(async () => {
    setFormState(prev => ({
      ...prev,
      submitting: true,
    }));

    try {
      // 验证表单
      const isValid = await validateAll();
      if (!isValid) {
        setFormState(prev => ({
          ...prev,
          submitting: false,
        }));
        return;
      }

      // 获取表单值
      const values = getFieldsValue();
      
      // 执行提交操作
      await submitOperation.execute(values);
    } catch (error) {
      setFormState(prev => ({
        ...prev,
        submitting: false,
      }));
      throw error;
    }
  }, [validateAll, getFieldsValue, submitOperation]);

  // 清除字段错误
  const clearFieldError = useCallback((name: string) => {
    setFormState(prev => ({
      ...prev,
      errors: {
        ...prev.errors,
        [name]: '',
      },
    }));
  }, []);

  // 清除所有错误
  const clearErrors = useCallback(() => {
    setFormState(prev => ({
      ...prev,
      errors: {},
    }));
  }, []);

  // 初始化表单值
  useEffect(() => {
    if (initialValues) {
      form.setFieldsValue(initialValues);
      setFormState(prev => ({
        ...prev,
        values: initialValues,
      }));
    }
  }, [form, initialValues]);

  return {
    // 表单实例
    form,
    formRef,
    
    // 表单状态
    ...formState,
    
    // 提交状态
    submitting: formState.submitting || submitOperation.loading,
    submitError: submitOperation.error,
    
    // 操作方法
    setFieldValue,
    setFieldsValue,
    getFieldValue,
    getFieldsValue,
    validateField,
    validateAll,
    submit,
    reset,
    clearFieldError,
    clearErrors,
    
    // 便捷属性
    isValid: Object.keys(formState.errors).length === 0,
    hasErrors: Object.keys(formState.errors).length > 0,
    canSubmit: !formState.submitting && formState.dirty && Object.keys(formState.errors).length === 0,
  };
}

// 简化的表单hook（只处理基本状态）
export function useSimpleForm<T extends Record<string, any>>(initialValues: Partial<T> = {}) {
  const [values, setValues] = useState<Partial<T>>(initialValues);
  const [dirty, setDirty] = useState(false);

  const setValue = useCallback((name: keyof T, value: any) => {
    setValues(prev => ({
      ...prev,
      [name]: value,
    }));
    setDirty(true);
  }, []);

  const setMultipleValues = useCallback((newValues: Partial<T>) => {
    setValues(prev => ({
      ...prev,
      ...newValues,
    }));
    setDirty(true);
  }, []);

  const reset = useCallback(() => {
    setValues(initialValues);
    setDirty(false);
  }, [initialValues]);

  return {
    values,
    setValue,
    setValues: setMultipleValues,
    reset,
    dirty,
  };
}