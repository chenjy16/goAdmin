import React, { useEffect, useState } from 'react';
import {
  Card,
  Table,
  Button,
  Switch,
  Tag,
  Space,
  Modal,
  Form,
  Input,
  message,
  Row,
  Col,
  Statistic,
} from 'antd';
import {
  CloudOutlined,
  SettingOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  KeyOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { useAppDispatch, useAppSelector } from '../store';
import {
  fetchProviders,
  fetchModels,
  setAPIKey,
  validateAPIKey,
  toggleModel,
} from '../store/slices/providersSlice';
import type { ProviderInfo, ModelInfo } from '../types/api';

const ProvidersPage: React.FC = () => {
  const dispatch = useAppDispatch();
  const { providers, models, isLoading, error } = useAppSelector(state => state.providers);
  const { apiKeys } = useAppSelector(state => state.settings);

  const [selectedProvider, setSelectedProvider] = useState<string | null>(null);
  const [apiKeyModalVisible, setApiKeyModalVisible] = useState(false);
  const [form] = Form.useForm();

  useEffect(() => {
    dispatch(fetchProviders());
  }, [dispatch]);

  useEffect(() => {
    if (selectedProvider) {
      dispatch(fetchModels(selectedProvider));
    }
  }, [dispatch, selectedProvider]);

  const handleSetAPIKey = async (values: { apiKey: string }) => {
    if (!selectedProvider) return;

    try {
      await dispatch(setAPIKey({
        provider: selectedProvider,
        apiKey: values.apiKey,
      })).unwrap();
      
      await dispatch(validateAPIKey(selectedProvider)).unwrap();

      message.success('API密钥设置成功');
      setApiKeyModalVisible(false);
      form.resetFields();
      dispatch(fetchProviders()); // 刷新提供商状态
    } catch (err) {
      message.error('API密钥设置失败');
    }
  };

  const handleToggleModel = async (provider: string, modelId: string, enabled: boolean) => {
    try {
      await dispatch(toggleModel({
        provider,
        model: modelId,
        enabled,
      })).unwrap();
      message.success(`模型${enabled ? '启用' : '禁用'}成功`);
    } catch (err) {
      message.error(`模型${enabled ? '启用' : '禁用'}失败`);
    }
  };

  const handleRefreshProviders = () => {
    dispatch(fetchProviders());
    if (selectedProvider) {
      dispatch(fetchModels(selectedProvider));
    }
  };

  const providerColumns: ColumnsType<ProviderInfo> = [
    {
      title: '提供商',
      dataIndex: 'name',
      key: 'name',
      render: (name: string, record: ProviderInfo) => (
        <Space>
          <CloudOutlined style={{ color: record.healthy ? '#52c41a' : '#ff4d4f' }} />
          <span style={{ fontWeight: 'bold' }}>{name}</span>
        </Space>
      ),
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: '状态',
      dataIndex: 'health',
      key: 'health',
      render: (health: boolean) => (
        <Tag color={health ? 'green' : 'red'} icon={health ? <CheckCircleOutlined /> : <ExclamationCircleOutlined />}>
          {health ? '正常' : '异常'}
        </Tag>
      ),
    },
    {
      title: '模型数量',
      dataIndex: 'model_count',
      key: 'model_count',
      render: (count: number) => (
        <Statistic value={count} suffix="个" />
      ),
    },
    {
      title: 'API密钥',
      key: 'apiKey',
      render: (_, record: ProviderInfo) => {
        const hasKey = apiKeys[record.name];
        return (
          <Tag color={hasKey ? 'green' : 'orange'}>
            {hasKey ? '已配置' : '未配置'}
          </Tag>
        );
      },
    },
    {
      title: '操作',
      key: 'actions',
      render: (_, record: ProviderInfo) => (
        <Space>
          <Button
            type="primary"
            size="small"
            icon={<KeyOutlined />}
            onClick={() => {
              setSelectedProvider(record.name);
              setApiKeyModalVisible(true);
            }}
          >
            设置密钥
          </Button>
          <Button
            size="small"
            icon={<SettingOutlined />}
            onClick={() => setSelectedProvider(record.name)}
          >
            管理模型
          </Button>
        </Space>
      ),
    },
  ];

  const modelColumns: ColumnsType<ModelInfo> = [
    {
      title: '模型名称',
      dataIndex: 'name',
      key: 'name',
      render: (name: string) => <span style={{ fontWeight: 'bold' }}>{name}</span>,
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: '最大令牌',
      dataIndex: 'max_tokens',
      key: 'max_tokens',
      render: (tokens: number) => tokens.toLocaleString(),
    },
    {
      title: '上下文窗口',
      dataIndex: 'context_window',
      key: 'context_window',
      render: (window: number) => window ? window.toLocaleString() : '-',
    },
    {
      title: '输入成本',
      dataIndex: 'input_cost',
      key: 'input_cost',
      render: (cost: number) => cost ? `$${cost}/1K tokens` : '-',
    },
    {
      title: '输出成本',
      dataIndex: 'output_cost',
      key: 'output_cost',
      render: (cost: number) => cost ? `$${cost}/1K tokens` : '-',
    },
    {
      title: '状态',
      dataIndex: 'enabled',
      key: 'enabled',
      render: (enabled: boolean, record: ModelInfo) => (
        <Switch
          checked={enabled}
          onChange={(checked) => handleToggleModel(selectedProvider!, record.id, checked)}
          checkedChildren="启用"
          unCheckedChildren="禁用"
        />
      ),
    },
  ];

  const healthyProviders = (providers || []).filter(p => p.healthy).length;
  const configuredProviders = (providers || []).filter(p => apiKeys[p.name]).length;

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
        <h1>AI提供商管理</h1>
        <Button
          type="primary"
          icon={<ReloadOutlined />}
          onClick={handleRefreshProviders}
          loading={isLoading}
        >
          刷新状态
        </Button>
      </div>

      {/* 统计信息 */}
      <Row gutter={[16, 16]} style={{ marginBottom: '24px' }}>
        <Col xs={24} sm={8}>
          <Card>
            <Statistic
              title="提供商总数"
              value={(providers || []).length}
              prefix={<CloudOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card>
            <Statistic
              title="健康状态"
              value={healthyProviders}
              suffix={`/ ${(providers || []).length}`}
              valueStyle={{ color: healthyProviders === (providers || []).length ? '#3f8600' : '#cf1322' }}
              prefix={<CheckCircleOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card>
            <Statistic
              title="已配置密钥"
              value={configuredProviders}
              suffix={`/ ${(providers || []).length}`}
              valueStyle={{ color: configuredProviders === (providers || []).length ? '#3f8600' : '#fa8c16' }}
              prefix={<KeyOutlined />}
            />
          </Card>
        </Col>
      </Row>

      {/* 提供商列表 */}
      <Card title="AI提供商列表" style={{ marginBottom: '24px' }}>
        <Table
          columns={providerColumns}
          dataSource={providers}
          rowKey="name"
          loading={isLoading}
          pagination={false}
        />
      </Card>

      {/* 模型管理 */}
      {selectedProvider && (
        <Card
          title={
            <Space>
              <span>模型管理 - {selectedProvider}</span>
              <Tag color="blue">{models[selectedProvider]?.length || 0} 个模型</Tag>
            </Space>
          }
          extra={
            <Button
              size="small"
              onClick={() => setSelectedProvider(null)}
            >
              关闭
            </Button>
          }
        >
          <Table
            columns={modelColumns}
            dataSource={models[selectedProvider] || []}
            rowKey="id"
            loading={isLoading}
            pagination={{
              pageSize: 10,
              showSizeChanger: true,
              showQuickJumper: true,
              showTotal: (total) => `共 ${total} 个模型`,
            }}
          />
        </Card>
      )}

      {/* API密钥设置模态框 */}
      <Modal
        title={`设置 ${selectedProvider} API密钥`}
        open={apiKeyModalVisible}
        onCancel={() => {
          setApiKeyModalVisible(false);
          form.resetFields();
        }}
        footer={null}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSetAPIKey}
        >
          <Form.Item
            label="API密钥"
            name="apiKey"
            rules={[
              { required: true, message: '请输入API密钥' },
              { min: 10, message: 'API密钥长度至少10个字符' },
            ]}
          >
            <Input.Password
              placeholder="请输入API密钥"
              autoComplete="off"
            />
          </Form.Item>
          
          <Form.Item style={{ marginBottom: 0 }}>
            <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
              <Button onClick={() => {
                setApiKeyModalVisible(false);
                form.resetFields();
              }}>
                取消
              </Button>
              <Button type="primary" htmlType="submit" loading={isLoading}>
                保存并验证
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>

      {error && (
        <div style={{ marginTop: '16px', padding: '16px', backgroundColor: '#fff2f0', border: '1px solid #ffccc7', borderRadius: '6px' }}>
          <span style={{ color: '#ff4d4f' }}>错误: {error}</span>
        </div>
      )}
    </div>
  );
};

export default ProvidersPage;