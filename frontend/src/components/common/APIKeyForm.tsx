import React from 'react';
import type { FormInstance } from 'antd/es/form';
import FormModal from './FormModal';
import FormField from './FormField';
import { useTranslation } from 'react-i18next';

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
  const { t } = useTranslation();
  return (
    <FormModal
      title={t('providers.setAPIKey', { provider })}
      open={open}
      onCancel={onCancel}
      onFinish={onFinish}
      form={form}
      loading={loading}
      okText={t('common.save')}
    >
      <FormField
          name="apiKey"
          label={t('providers.apiKey')}
          required
          placeholder={t('providers.enterAPIKey')}
          type="password"
          rules={[
            { min: 10, message: t('providers.apiKeyMinLength') },
          ]}
          tooltip={t('providers.apiKeySecurity')}
        />
    </FormModal>
  );
};

export default APIKeyForm;