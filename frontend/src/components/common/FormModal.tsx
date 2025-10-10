import React from 'react';
import { Modal, Form, Button, Space } from 'antd';
import type { FormInstance } from 'antd/es/form';

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
  okText = '确定',
  cancelText = '取消',
  destroyOnHidden = true,
}) => {
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
              {cancelText}
            </Button>
            <Button type="primary" htmlType="submit" loading={loading}>
              {okText}
            </Button>
          </Space>
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default FormModal;