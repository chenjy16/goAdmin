import React, { useEffect, useState } from 'react';
import {
  Card,
  Form,
  Switch,
  Button,
  Space,
  Typography,
  message,
  Alert,
  Tabs,
  List,
  Tag,
  Statistic,
} from 'antd';
import {
  SettingOutlined,
  ReloadOutlined,
  CloudOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  KeyOutlined,
  ToolOutlined,
  PlayCircleOutlined,
  EyeOutlined,
  EyeInvisibleOutlined,
  GithubOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { useAppDispatch, useAppSelector, store } from '../store';

import {
  fetchProviders,
  fetchModels,
  fetchAllModels,
  setAPIKey,
  validateAPIKey,
  toggleModel,
  fetchAPIKeyStatus,
  fetchPlainAPIKey,
} from '../store/slices/providersSlice';
import { setAPIKey as setSettingsAPIKey } from '../store/slices/settingsSlice';
import {
  initializeMCP,
  fetchMCPTools,
  fetchMCPLogs,
  checkMCPStatus,
} from '../store/slices/mcpSlice';
import type { ProviderInfo, ModelInfo, MCPTool } from '../types/api';
import { SearchableTable, APIKeyForm } from '../components/common';
import { useTranslation } from 'react-i18next';


const { Title, Text } = Typography;

const SettingsPage: React.FC = () => {
  const { t } = useTranslation();
  const dispatch = useAppDispatch();

  // 工具名称映射函数
  const getToolNameKey = (name: string): string => {
    const nameMap: Record<string, string> = {
      '雅虎财经': 'yahoo_finance',
      '股票分析': 'stock_analysis', 
      '股票对比': 'stock_compare',
      '股票投资建议': 'stock_advice',
    };
    return nameMap[name] || name;
  };

  // 获取国际化的工具名称
  const getLocalizedToolName = (name: string): string => {
    const key = getToolNameKey(name);
    const translationKey = `mcpTools.toolNames.${key}`;
    const translated = t(translationKey);
    return translated !== translationKey ? translated : name;
  };

  // 获取国际化的工具描述
  const getLocalizedToolDescription = (name: string, originalDescription: string): string => {
    const key = getToolNameKey(name);
    const translationKey = `mcpTools.toolDescriptions.${key}`;
    const translated = t(translationKey);
    return translated !== translationKey ? translated : originalDescription;
  };
  const settings = useAppSelector(state => state.settings);
  const { providers, models, apiKeyStatus, isLoading: providersLoading } = useAppSelector(state => state.providers);
  const { 
    tools: mcpTools, 
    isInitialized: mcpInitialized, 
    isLoading: mcpLoading, 
    error: mcpError 
  } = useAppSelector(state => state.mcp);
  const [form] = Form.useForm();
  
  // 提供商管理相关状态
  const [selectedProvider, setSelectedProvider] = useState<string | null>(null);
  const [apiKeyModalVisible, setApiKeyModalVisible] = useState(false);
  const [apiKeyForm] = Form.useForm();

  // API密钥明文显示相关状态
  const [plainAPIKeys, setPlainAPIKeys] = useState<Record<string, string>>({});
  const [showingPlainKeys, setShowingPlainKeys] = useState<Record<string, boolean>>({});

  // MCP工具相关状态
  const [forceRender, setForceRender] = useState(0);

  useEffect(() => {
    form.setFieldsValue(settings);
  }, [settings, form]);

  // 初始化提供商数据
  useEffect(() => {
    dispatch(fetchProviders());
    dispatch(fetchAPIKeyStatus());
  }, [dispatch]);

  // 当选择的提供商改变时，获取对应的模型列表
  useEffect(() => {
    if (selectedProvider) {
      dispatch(fetchAllModels(selectedProvider));
    }
  }, [selectedProvider, dispatch]);

  // 根据provider type获取显示名称
  const getProviderDisplayName = (providerType: string): string => {
    if (!providers || providers.length === 0) {
      return providerType;
    }
    const provider = providers.find(p => p && p.type === providerType);
    return provider && provider.name ? provider.name : providerType;
  };

  // 检查MCP状态并加载数据
  useEffect(() => {
    // 首先检查MCP状态
    dispatch(checkMCPStatus()).then((result) => {
      if (result.type === 'mcp/checkStatus/fulfilled') {
        // 如果已初始化，加载工具和日志
        const payload = result.payload as { initialized: boolean; toolCount: number; lastActivity?: string };
        if (payload && payload.initialized) {
          dispatch(fetchMCPTools());
          dispatch(fetchMCPLogs());
        }
      }
    }).catch((error) => {
      console.error('Check MCP status failed:', error);
    });
  }, [dispatch]);

  // 监听MCP状态变化
  useEffect(() => {
    // 强制重新渲染以确保UI更新
    setForceRender(prev => prev + 1);
  }, [mcpInitialized, mcpLoading, mcpError]);

  // 提供商管理相关函数
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

      message.success(t('settings.settingsSaved'));
      setApiKeyModalVisible(false);
      apiKeyForm.resetFields();
      dispatch(fetchProviders()); // 刷新提供商状态
      dispatch(fetchAPIKeyStatus()); // 刷新API密钥状态
    } catch (err) {
      message.error(t('errors.saveFailed'));
    }
  };

  // 处理明文显示切换
  const handleTogglePlainText = async (providerType: string) => {
    const isCurrentlyShowing = showingPlainKeys[providerType];
    
    if (isCurrentlyShowing) {
      // 隐藏明文
      setShowingPlainKeys(prev => ({ ...prev, [providerType]: false }));
      setPlainAPIKeys(prev => {
        const newKeys = { ...prev };
        delete newKeys[providerType];
        return newKeys;
      });
    } else {
      // 显示明文
      try {
        const result = await dispatch(fetchPlainAPIKey(providerType)).unwrap();
        setPlainAPIKeys(prev => ({ ...prev, [providerType]: result.apiKey }));
        setShowingPlainKeys(prev => ({ ...prev, [providerType]: true }));
        
        // 5秒后自动隐藏明文（安全考虑）
        setTimeout(() => {
          setShowingPlainKeys(prev => ({ ...prev, [providerType]: false }));
          setPlainAPIKeys(prev => {
            const newKeys = { ...prev };
            delete newKeys[providerType];
            return newKeys;
          });
        }, 5000);
        
      } catch (err) {
        message.error(t('errors.loadFailed'));
      }
    }
  };

  const handleToggleModel = async (provider: string, modelId: string, enabled: boolean) => {
    try {
      await dispatch(toggleModel({
        provider,
        model: modelId,
        enabled,
      })).unwrap();
      message.success(enabled ? t('success.enabled') : t('success.disabled'));
    } catch (err) {
      message.error(enabled ? t('errors.updateFailed') : t('errors.updateFailed'));
    }
  };

  // 批量操作处理
  const handleBatchModelAction = async (action: string, selectedKeys: React.Key[]) => {
    if (!selectedProvider) return;
    
    const providerModels = models[selectedProvider] || [];
    const selectedModels = providerModels.filter((model: ModelInfo) => 
      selectedKeys.includes(model.name)
    );
    
    try {
      const promises = selectedModels.map((model: ModelInfo) => {
        switch (action) {
          case 'enable':
            return dispatch(toggleModel({
              provider: selectedProvider,
              model: model.name,
              enabled: true,
            })).unwrap();
          case 'disable':
            return dispatch(toggleModel({
              provider: selectedProvider,
              model: model.name,
              enabled: false,
            })).unwrap();
          default:
            return Promise.resolve();
        }
      });
      
      await Promise.all(promises);
      const actionText = action === 'enable' ? t('common.enable') : t('common.disable');
      message.success(`${t('common.success')} ${actionText} ${selectedModels.length} ${t('providers.models')}`);
    } catch (err) {
      message.error(t('errors.updateFailed'));
    }
  };

  const handleRefreshProviders = () => {
    dispatch(fetchProviders());
    if (selectedProvider) {
      dispatch(fetchAllModels(selectedProvider));
    }
  };



  // MCP工具相关函数
  const handleInitializeMCP = () => {
    dispatch(initializeMCP());
  };

  const handleRefreshMCPTools = () => {
    dispatch(fetchMCPTools());
    dispatch(fetchMCPLogs());
  };







  // 提供商表格列定义
  const providerColumns: ColumnsType<ProviderInfo> = [
    {
      title: t('providers.name'),
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
      title: t('common.description'),
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: t('common.status'),
      dataIndex: 'healthy',
      key: 'healthy',
      render: (healthy: boolean) => (
        <Tag color={healthy ? 'green' : 'red'} icon={healthy ? <CheckCircleOutlined /> : <ExclamationCircleOutlined />}>
          {healthy ? t('providers.healthy') : t('providers.unhealthy')}
        </Tag>
      ),
    },
    {
      title: t('providers.models'),
      dataIndex: 'model_count',
      key: 'model_count',
      render: (count: number) => (
        <Statistic value={count} />
      ),
    },
    {
      title: t('providers.apiKey'),
      key: 'apiKey',
      render: (_, record: ProviderInfo) => {
        const keyInfo = apiKeyStatus && apiKeyStatus[record.type];
        const hasKey = keyInfo?.has_key || false;
        const maskedKey = keyInfo?.masked_key;
        const isShowingPlain = showingPlainKeys[record.type];
        const plainKey = plainAPIKeys[record.type];
        
        return (
          <div>
            <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
              <Tag color={hasKey ? 'green' : 'orange'}>
                {hasKey ? t('providers.configured') : t('providers.notConfigured')}
              </Tag>
              {hasKey && (
                <Button
                  type="text"
                  size="small"
                  icon={isShowingPlain ? <EyeInvisibleOutlined /> : <EyeOutlined />}
                  onClick={() => handleTogglePlainText(record.type)}
                  title={isShowingPlain ? t('common.hide') : t('common.view')}
                />
              )}
            </div>
            {hasKey && (
              <div style={{ fontSize: '12px', color: '#666', marginTop: '4px' }}>
                {isShowingPlain ? (
                  <span style={{ 
                    fontFamily: 'monospace', 
                    backgroundColor: '#f5f5f5', 
                    padding: '2px 4px', 
                    borderRadius: '2px',
                    color: '#d32f2f'
                  }}>
                    {plainKey}
                  </span>
                ) : (
                  maskedKey
                )}
              </div>
            )}
          </div>
        );
      },
    },
    {
      title: t('common.action'),
      key: 'actions',
      render: (_, record: ProviderInfo) => (
        <Space>
          <Button
            type="primary"
            size="small"
            icon={<KeyOutlined />}
            onClick={() => {
              setSelectedProvider(record.type);
              setApiKeyModalVisible(true);
            }}
          >
            {t('settings.addApiKey')}
          </Button>
          <Button
            size="small"
            icon={<SettingOutlined />}
            onClick={() => setSelectedProvider(record.type)}
          >
            {t('providers.configure')}
          </Button>
        </Space>
      ),
    },
  ];

  // 模型表格列定义
  const modelColumns: ColumnsType<ModelInfo> = [
    {
      title: t('providers.modelName'),
      dataIndex: 'name',
      key: 'name',
      render: (name: string) => <span style={{ fontWeight: 'bold' }}>{name}</span>,
    },
    {
      title: t('providers.displayName'),
      dataIndex: 'display_name',
      key: 'display_name',
    },
    {
      title: t('providers.maxTokens'),
      dataIndex: 'max_tokens',
      key: 'max_tokens',
      render: (tokens: number) => tokens.toLocaleString(),
    },
    {
      title: t('common.status'),
      dataIndex: 'enabled',
      key: 'enabled',
      render: (enabled: boolean, record: ModelInfo) => (
        <Switch
          checked={enabled}
          onChange={(checked) => handleToggleModel(selectedProvider!, record.name, checked)}
          checkedChildren={t('common.enable')}
          unCheckedChildren={t('common.disable')}
        />
      ),
    },
  ];

  // MCP工具表格列定义
  const toolColumns: ColumnsType<MCPTool> = [
    {
      title: t('mcpTools.toolName'),
      dataIndex: 'name',
      key: 'name',
      render: (name: string) => (
        <Space>
          <ToolOutlined />
          <span style={{ fontWeight: 'bold' }}>{getLocalizedToolName(name)}</span>
        </Space>
      ),
    },
    {
      title: t('common.description'),
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
      render: (description: string, record: MCPTool) => getLocalizedToolDescription(record.name, description),
    },
  ];





  const mcpSettings = (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
        <div>
          <Title level={4} style={{ margin: 0 }}>
            <ToolOutlined style={{ marginRight: '8px' }} />
            {t('mcpTools.title')}
          </Title>
        </div>
        <Space>
          {!mcpInitialized && (
            <Button
              type="primary"
              icon={<SettingOutlined />}
              onClick={handleInitializeMCP}
              loading={mcpLoading}
            >
              {t('mcpTools.initialize')}
            </Button>
          )}
        </Space>
      </div>


      {!mcpInitialized ? (
        <Card key={`uninitialized-${forceRender}`}>
          <div style={{ textAlign: 'center', padding: '40px' }}>
            <ExclamationCircleOutlined style={{ fontSize: '48px', color: '#faad14', marginBottom: '16px' }} />
            <h3>{t('mcpTools.notInitialized')}</h3>
            <p style={{ color: '#8c8c8c', marginBottom: '24px' }}>
              {t('mcpTools.initializePrompt')}
            </p>
            <Button
              type="primary"
              size="large"
              icon={<SettingOutlined />}
              onClick={handleInitializeMCP}
              loading={mcpLoading}
            >
              {t('mcpTools.initializeNow')}
            </Button>
          </div>
        </Card>
      ) : (
        <div key={`initialized-${forceRender}`}>
          {/* 工具列表 */}
          <SearchableTable<MCPTool>
            columns={toolColumns}
            dataSource={mcpTools}
            rowKey="name"
            loading={mcpLoading}
            searchFields={['name', 'description']}
            searchPlaceholder={t('mcpTools.searchPlaceholder')}
            showRefresh={true}
            onRefresh={() => dispatch(fetchMCPTools())}
            refreshLoading={mcpLoading}
            title={t('mcpTools.availableTools')}
          />




        </div>
      )}

      {mcpError && (
        <Alert
          message={t('common.error')}
          description={mcpError}
          type="error"
          showIcon
          style={{ marginTop: '16px' }}
        />
      )}
    </div>
  );



  // 提供商管理组件
  const providersManagement = (
    <div>
      <div style={{ marginBottom: '24px' }}>
        <Title level={4} style={{ margin: 0 }}>
          <CloudOutlined style={{ marginRight: '8px' }} />
          {t('providers.title')}
        </Title>
        <Text type="secondary">{t('providers.description')}</Text>
      </div>



      {/* 提供商列表 */}
      <SearchableTable<ProviderInfo>
        columns={providerColumns}
        dataSource={providers}
        rowKey="name"
        loading={providersLoading}
        searchFields={['name', 'description']}
        searchPlaceholder={t('providers.searchPlaceholder')}
        showRefresh={true}
        onRefresh={handleRefreshProviders}
        refreshLoading={providersLoading}
        title={t('providers.list')}
      />

      {/* 模型管理 */}
      {selectedProvider && (
        <Card
          title={
            <Space>
              <span>{t('providers.modelManagement')} - {getProviderDisplayName(selectedProvider)}</span>
              <Tag color="blue">{models[selectedProvider]?.length || 0} {t('providers.models')}</Tag>
            </Space>
          }
          extra={
            <Button
              size="small"
              onClick={() => setSelectedProvider(null)}
            >
              {t('common.close')}
            </Button>
          }
        >
          <SearchableTable<ModelInfo>
            columns={modelColumns}
            dataSource={models[selectedProvider] || []}
            rowKey="name"
            loading={providersLoading}
            searchFields={['name', 'display_name']}
            searchPlaceholder={t('providers.searchModels')}
            showRefresh={true}
            onRefresh={() => selectedProvider && dispatch(fetchAllModels(selectedProvider))}
            refreshLoading={providersLoading}
            enableBatchSelection={true}
            batchActions={[
              { key: 'enable', label: t('providers.batchEnable') },
              { key: 'disable', label: t('providers.batchDisable') },
            ]}
            onBatchAction={handleBatchModelAction}
          />
        </Card>
      )}

      {/* API密钥设置模态框 */}
      <APIKeyForm
        provider={selectedProvider || ''}
        open={apiKeyModalVisible}
        onCancel={() => setApiKeyModalVisible(false)}
        onFinish={handleSetAPIKey}
        form={apiKeyForm}
        loading={providersLoading}
      />
    </div>
  );

  const aboutInfo = (
    <Card title={t('navigation.about')} style={{ marginBottom: 16 }}>
      <List>
        <List.Item>
          <List.Item.Meta
            title={t('settings.appVersion')}
            description="Go-SpringAI v1.0.0"
          />
        </List.Item>
        <List.Item>
          <List.Item.Meta
            title={t('settings.buildTime')}
            description={new Date().toLocaleDateString()}
          />
        </List.Item>
        <List.Item>
          <List.Item.Meta
            avatar={<GithubOutlined style={{ fontSize: '20px', color: '#1890ff' }} />}
            title={t('settings.sourceCode')}
            description={
              <a 
                href="https://github.com/chenjy16/go-springAi" 
                target="_blank" 
                rel="noopener noreferrer"
                style={{ color: '#1890ff' }}
              >
                https://github.com/chenjy16/go-springAi
              </a>
            }
          />
        </List.Item>
      </List>
    </Card>
  );

  return (
    <div style={{ padding: '24px', maxWidth: '1200px', margin: '0 auto' }}>
      <div style={{ marginBottom: '24px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <div>
          <Title level={2} style={{ margin: 0 }}>
            <SettingOutlined style={{ marginRight: '8px' }} />
            {t('navigation.settings')}
          </Title>
          <Text type="secondary">{t('settings.description')}</Text>
        </div>

      </div>



      <Form
        form={form}
        layout="vertical"
        initialValues={settings}
      >
        <Tabs 
          defaultActiveKey="providers" 
          type="card"
          items={[
            {
              key: 'providers',
              label: <span><CloudOutlined />{t('providers.title')}</span>,
              children: providersManagement,
            },
            {
              key: 'mcp',
              label: t('mcpTools.title'),
              children: mcpSettings,
            },
            {
              key: 'about',
              label: t('navigation.about'),
              children: aboutInfo,
            },
          ]}
        />
      </Form>
    </div>
  );
};

export default SettingsPage;