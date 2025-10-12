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
import { useTranslation } from 'react-i18next';

const ProvidersPage: React.FC = () => {
  const { t } = useTranslation();
  const dispatch = useAppDispatch();
  const { providers, models, isLoading, error, apiKeyStatus } = useAppSelector(state => state.providers);
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
      successMessage: t('settings.settingsSaved'),
      errorMessage: t('errors.saveFailed')
    }
  );

  const toggleModelOperation = useAsyncOperation(
    async (provider: string, modelId: string, enabled: boolean) => {
      await dispatch(toggleModel({ provider, model: modelId, enabled })).unwrap();
    },
    {
      successMessage: '',
      errorMessage: t('errors.updateFailed')
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
      message.success(t('success.updated'));
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
      title: t('providers.providerName'),
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
      title: t('providers.description'),
      dataIndex: 'description',
    }),
    createStatusColumn<ProviderInfo>({
      title: t('providers.healthStatus'),
      dataIndex: 'healthy',
      statusMap: {
        true: { color: '#52c41a', text: t('providers.healthy'), icon: <CheckCircleOutlined /> },
        false: { color: '#ff4d4f', text: t('providers.unhealthy'), icon: <ExclamationCircleOutlined /> }
      }
    }),
    {
      title: t('providers.modelCount'),
      dataIndex: 'model_count',
      key: 'model_count',
      render: (count: number) => (
        <Statistic value={count} suffix={t('common.items')} />
      ),
    },
    {
      title: t('providers.apiKey'),
      key: 'api_key',
      render: (_: any, record: ProviderInfo) => {
        const status = apiKeyStatus[record.name];
        return (
          <div>
            {status?.configured ? (
              <Tag color="green">{t('providers.configured')}</Tag>
            ) : (
              <Tag color="red">{t('providers.notConfigured')}</Tag>
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
          label: t('providers.setKey'),
          type: 'primary',
          icon: <KeyOutlined />,
          onClick: (record) => {
            setSelectedProvider(record.name);
            setApiKeyModalVisible(true);
          }
        },
        {
          key: 'manage',
          label: t('providers.manageModels'),
          icon: <SettingOutlined />,
          onClick: (record) => setSelectedProvider(record.name)
        }
      ]
    })
  ]);

  const modelColumns = mergeColumns([
    {
      title: t('providers.modelName'),
      dataIndex: 'name',
      key: 'name',
      render: (name: string) => <span style={{ fontWeight: 'bold' }}>{name}</span>,
    },
    createTextColumn<ModelInfo>({
      title: t('providers.description'),
      dataIndex: 'description',
    }),
    {
      title: t('providers.maxTokens'),
      dataIndex: 'max_tokens',
      key: 'max_tokens',
      render: (tokens: number) => tokens.toLocaleString(),
    },
    {
      title: t('providers.contextWindow'),
      dataIndex: 'context_window',
      key: 'context_window',
      render: (window: number) => window ? window.toLocaleString() : '-',
    },
    {
      title: t('providers.inputCost'),
      dataIndex: 'input_cost',
      key: 'input_cost',
      render: (cost: number) => cost ? `$${cost}/1K tokens` : '-',
    },
    {
      title: t('providers.outputCost'),
      dataIndex: 'output_cost',
      key: 'output_cost',
      render: (cost: number) => cost ? `$${cost}/1K tokens` : '-',
    },
    createSwitchColumn<ModelInfo>({
      title: t('providers.status'),
      dataIndex: 'enabled',
      onChange: (checked, record) => handleToggleModel(selectedProvider!, record.name, checked)
    })
  ]);



  return (
    <div>
      <div style={{ marginBottom: '24px' }}>
        <h1>{t('providers.title')}</h1>
      </div>



      {/* 提供商列表 */}
      <SearchableTable<ProviderInfo>
        columns={providerColumns}
        dataSource={providers}
        rowKey="name"
        loading={isLoading}
        searchFields={['name', 'description']}
        searchPlaceholder={t('providers.searchProviders')}
        showRefresh={true}
        onRefresh={() => dispatch(fetchProviders())}
        refreshLoading={isLoading}
        title={t('providers.providerList')}
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
              <span>{t('providers.modelManagement')} - {selectedProvider}</span>
              <Tag color="blue">{models[selectedProvider]?.length || 0} {t('providers.modelsCount')}</Tag>
            </Space>
          }
          extra={
            <Button
              size="small"
              onClick={() => setSelectedProvider(null)}
            >
              {t('providers.close')}
            </Button>
          }
        >
          <SearchableTable<ModelInfo>
            columns={modelColumns}
            dataSource={models[selectedProvider] || []}
            rowKey="id"
            loading={isLoading}
            searchFields={['name', 'description']}
            searchPlaceholder={t('providers.searchModels')}
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
          <span style={{ color: '#ff4d4f' }}>{t('providers.error')}: {error}</span>
        </div>
      )}
    </div>
  );
};

export default ProvidersPage;