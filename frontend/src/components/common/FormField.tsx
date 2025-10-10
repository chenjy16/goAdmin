import React from 'react';
import { Form, Input, Select, Switch, InputNumber, Slider } from 'antd';
import type { Rule } from 'antd/es/form';

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
  const { name, label, type, placeholder, rules, disabled, required, tooltip } = props;

  const baseRules = required ? [{ required: true, message: `请输入${label}` }, ...(rules || [])] : rules;

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
        return (
          <Switch
            disabled={disabled}
            checkedChildren={switchProps.checkedChildren}
            unCheckedChildren={switchProps.unCheckedChildren}
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
            style={{ width: '100%' }}
          />
        );

      case 'slider':
        const sliderProps = props as SliderFieldProps;
        return (
          <Slider
            disabled={disabled}
            min={sliderProps.min}
            max={sliderProps.max}
            step={sliderProps.step}
            marks={sliderProps.marks}
            range={sliderProps.range}
          />
        );

      default:
        return null;
    }
  };

  return (
    <Form.Item
      name={name}
      label={label}
      rules={baseRules}
      tooltip={tooltip}
      valuePropName={type === 'switch' ? 'checked' : 'value'}
    >
      {renderField()}
    </Form.Item>
  );
};

export default FormField;