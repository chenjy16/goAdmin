import React from 'react';
import { Form, Input, Select, Switch, InputNumber, Slider } from 'antd';
import type { Rule } from 'antd/es/form';
import { colors, spacing, borderRadius, shadows, theme, combine } from '../../styles';

const { TextArea } = Input;
const { Option } = Select;

export type FieldType = 'input' | 'password' | 'textarea' | 'select' | 'switch' | 'number' | 'slider';

interface BaseFieldProps {
  name: string;
  label: string;
  type: FieldType;
  placeholder?: string;
  rules?: Rule[];
  disabled?: boolean;
  required?: boolean;
  tooltip?: string;
  className?: string;
  style?: React.CSSProperties;
  size?: 'small' | 'middle' | 'large';
  variant?: 'outlined' | 'filled' | 'borderless';
}

interface InputFieldProps extends BaseFieldProps {
  type: 'input' | 'password' | 'textarea';
  maxLength?: number;
  showCount?: boolean;
  autoSize?: boolean | { minRows?: number; maxRows?: number };
}

interface SelectFieldProps extends BaseFieldProps {
  type: 'select';
  options: Array<{ label: string; value: any; disabled?: boolean }>;
  mode?: 'multiple' | 'tags';
  allowClear?: boolean;
  showSearch?: boolean;
}

interface SwitchFieldProps extends BaseFieldProps {
  type: 'switch';
  checkedChildren?: React.ReactNode;
  unCheckedChildren?: React.ReactNode;
}

interface NumberFieldProps extends BaseFieldProps {
  type: 'number';
  min?: number;
  max?: number;
  step?: number;
  precision?: number;
}

interface SliderFieldProps extends BaseFieldProps {
  type: 'slider';
  min?: number;
  max?: number;
  step?: number;
  marks?: Record<number, React.ReactNode>;
  range?: boolean;
}

export type FormFieldProps = 
  | InputFieldProps 
  | SelectFieldProps 
  | SwitchFieldProps 
  | NumberFieldProps 
  | SliderFieldProps;

const FormField: React.FC<FormFieldProps> = (props) => {
  const { 
    name, 
    label, 
    type, 
    placeholder, 
    rules, 
    disabled, 
    required, 
    tooltip,
    className = '',
    style,
    size = 'middle',
    variant = 'outlined'
  } = props;

  const baseRules = required ? [{ required: true, message: `请输入${label}` }, ...(rules || [])] : rules;

  // 创建自定义样式
  const getFieldStyle = () => {
    const baseStyle = theme.getInputStyle();
    const customStyle = combine(
      baseStyle,
      {
        fontSize: size === 'small' ? '12px' : size === 'large' ? '16px' : '14px',
        ...(variant === 'filled' && {
          backgroundColor: colors.background('light'),
          border: 'none'
        }),
        ...(variant === 'borderless' && {
          border: 'none',
          boxShadow: 'none'
        }),
        ...(disabled && {
          backgroundColor: colors.background('light'),
          color: colors.text('disabled'),
          cursor: 'not-allowed'
        })
      },
      style
    );
    return customStyle;
  };

  const fieldStyle = getFieldStyle();

  const renderField = () => {
    switch (type) {
      case 'input':
        const inputProps = props as InputFieldProps;
        return (
          <Input
            placeholder={placeholder}
            disabled={disabled}
            maxLength={inputProps.maxLength}
            showCount={inputProps.showCount}
            size={size}
            variant={variant}
            style={fieldStyle}
            className={className}
          />
        );

      case 'password':
        const passwordProps = props as InputFieldProps;
        return (
          <Input.Password
            placeholder={placeholder}
            disabled={disabled}
            maxLength={passwordProps.maxLength}
            autoComplete="off"
            size={size}
            variant={variant}
            style={fieldStyle}
            className={className}
          />
        );

      case 'textarea':
        const textareaProps = props as InputFieldProps;
        return (
          <TextArea
            placeholder={placeholder}
            disabled={disabled}
            maxLength={textareaProps.maxLength}
            showCount={textareaProps.showCount}
            autoSize={textareaProps.autoSize}
            size={size}
            variant={variant}
            style={fieldStyle}
            className={className}
          />
        );

      case 'select':
        const selectProps = props as SelectFieldProps;
        return (
          <Select
            placeholder={placeholder}
            disabled={disabled}
            mode={selectProps.mode}
            allowClear={selectProps.allowClear}
            showSearch={selectProps.showSearch}
            size={size}
            variant={variant}
            style={fieldStyle}
            className={className}
          >
            {selectProps.options.map(option => (
              <Option 
                key={option.value} 
                value={option.value} 
                disabled={option.disabled}
              >
                {option.label}
              </Option>
            ))}
          </Select>
        );

      case 'switch':
        const switchProps = props as SwitchFieldProps;
        const switchStyle = combine(
          {
            ...(size === 'small' && { transform: 'scale(0.8)' }),
            ...(size === 'large' && { transform: 'scale(1.2)' })
          },
          style
        );
        return (
          <Switch
            disabled={disabled}
            checkedChildren={switchProps.checkedChildren}
            unCheckedChildren={switchProps.unCheckedChildren}
            size={size === 'small' ? 'small' : 'default'}
            style={switchStyle}
            className={className}
          />
        );

      case 'number':
        const numberProps = props as NumberFieldProps;
        return (
          <InputNumber
            placeholder={placeholder}
            disabled={disabled}
            min={numberProps.min}
            max={numberProps.max}
            step={numberProps.step}
            precision={numberProps.precision}
            size={size}
            variant={variant}
            style={combine(fieldStyle, { width: '100%' })}
            className={className}
          />
        );

      case 'slider':
        const sliderProps = props as SliderFieldProps;
        const sliderStyle = combine(
          spacing.marginY('sm'),
          style
        );
        return (
          <Slider
            disabled={disabled}
            min={sliderProps.min}
            max={sliderProps.max}
            step={sliderProps.step}
            marks={sliderProps.marks}
            range={sliderProps.range}
            style={sliderStyle}
            className={className}
          />
        );

      default:
        return null;
    }
  };

  // 创建Form.Item的样式
  const formItemStyle = combine(
    spacing.marginBottom('md'),
    style
  );

  return (
    <Form.Item
      name={name}
      label={label}
      rules={baseRules}
      tooltip={tooltip}
      valuePropName={type === 'switch' ? 'checked' : 'value'}
      style={formItemStyle}
      className={className}
    >
      {renderField()}
    </Form.Item>
  );
};

export default FormField;