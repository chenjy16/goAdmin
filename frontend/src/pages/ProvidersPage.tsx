import React, { useEffect, useState } from 'react';
import {
  Card,
  Button,
  Switch,
  Tag,
  Space,
  Form,
  Statistic,
  message,
} from 'antd';
import {
  CloudOutlined,
  SettingOutlined,
  KeyOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
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
import { useAsyncOperation } from '../hooks';
import { 
  createTextColumn, 
  createStatusColumn, 
  createSwitchColumn, 
  createActionColumn,
  commonActions,
  mergeColumns 
} from '../utils/tableColumns';

const ProvidersPage: React.FC = () => {
  const dispatch = useAppDispatch();
  const { providers, models, apiKeyStatus, isLoading, error } = useAppSelector(state => state.providers);

  const [selectedProvider, setSelectedProvider] = useState<string | null>(null);
  const [apiKeyModalVisible, setApiKeyModalVisible] = useState(false);
  const [form] = Form.useForm();

  // 使用通用异步操作hook
  const setAPIKeyOperation = useAsyncOperation(
    async (provider: string, apiKey: string) => {
      await dispatch(setAPIKey({ provider, apiKey })).unwrap();
      await dispatch(validateAPIKey(provider)).unwrap();
      dispatch(setSettingsAPIKey({ [provider]: apiKey }));
      dispatch(fetchProviders());
      dispatch(fetchAPIKeyStatus());
    },
    {
      successMessage: 'API密钥设置成功',
      errorMessage: 'API密钥设置失败'
    }
  );

  const toggleModelOperation = useAsyncOperation(
    async (provider: string, modelId: string, enabled: boolean) => {
      await dispatch(toggleModel({ provider, model: modelId, enabled })).unwrap();
    },
    {
      successMessage: '',
      errorMessage: '操作失败'
    }
  );

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

    const result = await setAPIKeyOperation.execute(selectedProvider, values.apiKey);
    if (result) {
      setApiKeyModalVisible(false);
      form.resetFields();
    }
  };

  const handleToggleModel = async (provider: string, modelId: string, enabled: boolean) => {
    // 动态设置成功消息
    const result = await toggleModelOperation.execute(provider, modelId, enabled);
    if (result) {
      message.success(`模型${enabled ? '启用' : '禁用'}成功`);
    }
  };

  const handleRefreshProviders = () => {
    dispatch(fetchProviders());
    if (selectedProvider) {
      dispatch(fetchAllModels(selectedProvider));
    }
  };

  const providerColumns = mergeColumns([
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
    createTextColumn<ProviderInfo>({
      title: '描述',
      dataIndex: 'description',
    }),
    createStatusColumn<ProviderInfo>({
      title: '健康状态',
      dataIndex: 'healthy',
      statusMap: {
        true: { color: '#52c41a', text: '健康', icon: <CheckCircleOutlined /> },
        false: { color: '#ff4d4f', text: '异常', icon: <ExclamationCircleOutlined /> }
      }
    }),
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
      key: 'api_key',
      render: (_: any, record: ProviderInfo) => {
        const status = apiKeyStatus[record.name];
        return (
          <div>
            {status?.configured ? (
              <Tag color="green">已配置</Tag>
            ) : (
              <Tag color="red">未配置</Tag>
            )}
            {status?.configured && status.masked_key && (
              <div style={{ fontSize: '12px', color: '#666', marginTop: '4px' }}>
                {status.masked_key}
              </div>
            )}
          </div>
        );
      },
    },
    createActionColumn<ProviderInfo>({
      actions: [
        {
          key: 'setKey',
          label: '设置密钥',
          type: 'primary',
          icon: <KeyOutlined />,
          onClick: (record) => {
            setSelectedProvider(record.name);
            setApiKeyModalVisible(true);
          }
        },
        {
          key: 'manageModels',
          label: '管理模型',
          icon: <SettingOutlined />,
          onClick: (record) => setSelectedProvider(record.name)
        }
      ]
    })
  ]);

  const modelColumns = mergeColumns([
    {
      title: '模型名称',
      dataIndex: 'name',
      key: 'name',
      render: (name: string) => <span style={{ fontWeight: 'bold' }}>{name}</span>,
    },
    createTextColumn<ModelInfo>({
      title: '描述',
      dataIndex: 'description',
    }),
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
    createSwitchColumn<ModelInfo>({
      title: '状态',
      dataIndex: 'enabled',
      onChange: (checked, record) => handleToggleModel(selectedProvider!, record.name, checked)
    })
  ]);



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