import React from 'react';
import type { FormInstance } from 'antd/es/form';
import FormModal from './FormModal';
import FormField from './FormField';

interface APIKeyFormProps {
  provider: string;
  open: boolean;
  onCancel: () => void;
  onFinish: (values: { apiKey: string }) => void;
  form: FormInstance;
  loading?: boolean;
}

const APIKeyForm: React.FC<APIKeyFormProps> = ({
  provider,
  open,
  onCancel,
  onFinish,
  form,
  loading = false,
}) => {
  return (
    <FormModal
      title={`设置 ${provider} API密钥`}
      open={open}
      onCancel={onCancel}
      onFinish={onFinish}
      form={form}
      loading={loading}
      okText="保存"
    >
      <FormField
        name="apiKey"
        label="API密钥"
        type="password"
        placeholder="请输入API密钥"
        required
        rules={[
          { min: 10, message: 'API密钥长度至少10个字符' },
        ]}
        tooltip="请确保API密钥的安全性，不要在不安全的环境中使用"
      />
    </FormModal>
  );
};

export default APIKeyForm;