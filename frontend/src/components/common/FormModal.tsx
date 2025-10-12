import React from 'react';
import { Modal, Form, Button, Space } from 'antd';
import { useTranslation } from 'react-i18next';
import type { FormInstance } from 'antd/es/form';
import type { ModalProps } from 'antd/es/modal';

interface FormModalProps {
  title: string;
  open: boolean;
  onCancel: () => void;
  onFinish: (values: any) => void;
  form: FormInstance;
  loading?: boolean;
  width?: number;
  children: React.ReactNode;
  okText?: string;
  cancelText?: string;
  destroyOnHidden?: boolean;
}

const FormModal: React.FC<FormModalProps> = ({
  title,
  open,
  onCancel,
  onFinish,
  form,
  loading = false,
  width = 520,
  children,
  okText,
  cancelText,
  destroyOnHidden = true,
}) => {
  const { t } = useTranslation();
  const handleCancel = () => {
    form.resetFields();
    onCancel();
  };

  return (
    <Modal
      title={title}
      open={open}
      onCancel={handleCancel}
      width={width}
      destroyOnHidden={destroyOnHidden}
      footer={null}
    >
      <Form
        form={form}
        layout="vertical"
        onFinish={onFinish}
      >
        {children}
        
        <Form.Item style={{ marginBottom: 0, marginTop: 24 }}>
          <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
            <Button onClick={handleCancel}>
              {cancelText || t('common.cancel')}
            </Button>
            <Button type="primary" htmlType="submit" loading={loading}>
              {okText || t('common.confirm')}
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default FormModal;