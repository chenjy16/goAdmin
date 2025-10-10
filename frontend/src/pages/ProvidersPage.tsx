import React, { useEffect, useState } from 'react';
import {
  Card,
  Button,
  Switch,
  Tag,
  Space,
  Form,
  message,
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
  fetchAllModels,
  setAPIKey,
  validateAPIKey,
  toggleModel,
  fetchAPIKeyStatus,
} from '../store/slices/providersSlice';
import { setAPIKey as setSettingsAPIKey } from '../store/slices/settingsSlice';
import type { ProviderInfo, ModelInfo } from '../types/api';
import { SearchableTable, APIKeyForm } from '../components/common';

const ProvidersPage: React.FC = () => {
  const dispatch = useAppDispatch();
  const { providers, models, apiKeyStatus, isLoading, error } = useAppSelector(state => state.providers);

  const [selectedProvider, setSelectedProvider] = useState<string | null>(null);
  const [apiKeyModalVisible, setApiKeyModalVisible] = useState(false);
  const [form] = Form.useForm();

  useEffect(() => {
    dispatch(fetchProviders());
    dispatch(fetchAPIKeyStatus());
  }, [dispatch]);

  useEffect(() => {
    if (selectedProvider) {
      dispatch(fetchAllModels(selectedProvider));
    }
  }, [dispatch, selectedProvider]);

  const handleSetAPIKey = async (values: { apiKey: string }) => {
    if (!selectedProvider) return;

    try {
      await dispatch(setAPIKey({
        provider: selectedProvider,
        apiKey: values.apiKey,
      })).unwrap();
      
      // 同时更新本地settings状态
      dispatch(setSettingsAPIKey(selectedProvider, values.apiKey));
      
      await dispatch(validateAPIKey(selectedProvider)).unwrap();

      message.success('API密钥设置成功');
      setApiKeyModalVisible(false);
      form.resetFields();
      dispatch(fetchProviders()); // 刷新提供商状态
      dispatch(fetchAPIKeyStatus()); // 刷新API密钥状态
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
      dispatch(fetchAllModels(selectedProvider));
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
        const keyInfo = apiKeyStatus && apiKeyStatus[record.type];
        const hasKey = keyInfo?.has_key || false;
        const maskedKey = keyInfo?.masked_key;
        
        return (
          <div>
            <Tag color={hasKey ? 'green' : 'orange'}>
              {hasKey ? '已配置' : '未配置'}
            </Tag>
            {hasKey && maskedKey && (
              <div style={{ fontSize: '12px', color: '#666', marginTop: '4px' }}>
                {maskedKey}
              </div>
            )}
          </div>
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
          onChange={(checked) => handleToggleModel(selectedProvider!, record.name, checked)}
          checkedChildren="启用"
          unCheckedChildren="禁用"
        />
      ),
    },
  ];



  return (
    <div>
      <div style={{ marginBottom: '24px' }}>
        <h1>AI提供商管理</h1>
      </div>



      {/* 提供商列表 */}
      <SearchableTable<ProviderInfo>
        columns={providerColumns}
        dataSource={providers}
        rowKey="name"
        loading={isLoading}
        searchFields={['name', 'description']}
        searchPlaceholder="搜索提供商..."
        showRefresh={true}
        onRefresh={() => dispatch(fetchProviders())}
        refreshLoading={isLoading}
        title="AI提供商列表"
        onRow={(record) => ({
          onClick: () => {
            setSelectedProvider(record.name);
          },
        })}
      />

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
          <SearchableTable<ModelInfo>
            columns={modelColumns}
            dataSource={models[selectedProvider] || []}
            rowKey="id"
            loading={isLoading}
            searchFields={['name', 'description']}
            searchPlaceholder="搜索模型..."
            showRefresh={true}
            onRefresh={() => selectedProvider && dispatch(fetchAllModels(selectedProvider))}
            refreshLoading={isLoading}
          />
        </Card>
      )}

      {/* API密钥设置模态框 */}
      <APIKeyForm
        provider={selectedProvider || ''}
        open={apiKeyModalVisible}
        onCancel={() => setApiKeyModalVisible(false)}
        onFinish={handleSetAPIKey}
        form={form}
        loading={isLoading}
      />

      {error && (
        <div style={{ marginTop: '16px', padding: '16px', backgroundColor: '#fff2f0', border: '1px solid #ffccc7', borderRadius: '6px' }}>
          <span style={{ color: '#ff4d4f' }}>错误: {error}</span>
        </div>
      )}
    </div>
  );
};

export default ProvidersPage;