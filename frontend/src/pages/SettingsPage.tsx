import React, { useEffect, useState } from 'react';
import {
  Card,
  Form,
  Input,
  Switch,
  Select,
  Button,
  Space,
  Typography,
  InputNumber,
  message,
  Alert,
  Tabs,
  List,
  Tag,
  Table,
  Modal,
  Row,
  Col,
  Statistic,
  Divider,
  Collapse,
} from 'antd';
import {
  SettingOutlined,
  SaveOutlined,
  ReloadOutlined,
  CloudOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  KeyOutlined,
  ToolOutlined,
  PlayCircleOutlined,
  CodeOutlined,
  HistoryOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { useAppDispatch, useAppSelector } from '../store';
import {
  resetSettings,
  loadSettings,
} from '../store/slices/settingsSlice';
import {
  fetchProviders,
  fetchModels,
  setAPIKey,
  validateAPIKey,
  toggleModel,
} from '../store/slices/providersSlice';
import {
  initializeMCP,
  fetchMCPTools,
  executeMCPTool,
  fetchMCPLogs,
} from '../store/slices/mcpSlice';
import type { ProviderInfo, ModelInfo, MCPTool, MCPMessage } from '../types/api';

const { Title, Text, Paragraph } = Typography;
const { TextArea } = Input;
const { TabPane } = Tabs;

const SettingsPage: React.FC = () => {
  const dispatch = useAppDispatch();
  const settings = useAppSelector(state => state.settings);
  const { providers, models, isLoading: providersLoading } = useAppSelector(state => state.providers);
  const { 
    tools: mcpTools, 
    logs: mcpLogs, 
    isInitialized: mcpInitialized, 
    isLoading: mcpLoading, 
    error: mcpError 
  } = useAppSelector(state => state.mcp);
  const { apiKeys } = settings;
  const [form] = Form.useForm();
  const [hasChanges, setHasChanges] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  
  // 提供商管理相关状态
  const [selectedProvider, setSelectedProvider] = useState<string | null>(null);
  const [apiKeyModalVisible, setApiKeyModalVisible] = useState(false);
  const [apiKeyForm] = Form.useForm();

  // MCP工具相关状态
  const [executeModalVisible, setExecuteModalVisible] = useState(false);
  const [selectedTool, setSelectedTool] = useState<MCPTool | null>(null);
  const [mcpForm] = Form.useForm();

  useEffect(() => {
    form.setFieldsValue(settings);
  }, [settings, form]);

  // 初始化提供商数据
  useEffect(() => {
    dispatch(fetchProviders());
  }, [dispatch]);

  // 当选择的提供商改变时，获取对应的模型列表
  useEffect(() => {
    if (selectedProvider) {
      dispatch(fetchModels(selectedProvider));
    }
  }, [selectedProvider, dispatch]);

  // 根据provider type获取显示名称
  const getProviderDisplayName = (providerType: string): string => {
    const provider = providers.find(p => p.type === providerType);
    return provider ? provider.name : providerType;
  };

  // 初始化MCP数据
  useEffect(() => {
    if (!mcpInitialized) {
      dispatch(initializeMCP());
    } else {
      dispatch(fetchMCPTools());
      dispatch(fetchMCPLogs());
    }
  }, [dispatch, mcpInitialized]);

  // 提供商管理相关函数
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
      apiKeyForm.resetFields();
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

  const handleSave = async () => {
    try {
      setIsLoading(true);
      const values = await form.validateFields();
      dispatch(loadSettings(values));
      setHasChanges(false);
      message.success('设置已保存');
    } catch (err) {
      message.error('保存设置失败');
    } finally {
      setIsLoading(false);
    }
  };

  const handleReset = () => {
    dispatch(resetSettings());
    form.resetFields();
    setHasChanges(false);
    message.success('设置已重置为默认值');
  };



  const handleFormChange = () => {
    setHasChanges(true);
  };

  // MCP工具相关函数
  const handleInitializeMCP = () => {
    dispatch(initializeMCP());
  };

  const handleRefreshMCPTools = () => {
    dispatch(fetchMCPTools());
    dispatch(fetchMCPLogs());
  };

  const handleExecuteMCPTool = async (values: Record<string, any>) => {
    if (!selectedTool) return;

    try {
      await dispatch(executeMCPTool({
        name: selectedTool.name,
        arguments: values,
      })).unwrap();
      
      message.success('工具执行成功');
      setExecuteModalVisible(false);
      mcpForm.resetFields();
      dispatch(fetchMCPLogs()); // 刷新日志
    } catch (err) {
      message.error('工具执行失败');
    }
  };

  const renderExecutionForm = (inputSchema: any) => {
    if (!inputSchema || !inputSchema.properties) {
      return <div>该工具无需参数</div>;
    }

    const properties = inputSchema.properties;
    const required = inputSchema.required || [];

    return Object.entries(properties).map(([key, schema]: [string, any]) => {
      const isRequired = required.includes(key);
      const rules = isRequired ? [{ required: true, message: `请输入${schema.title || key}` }] : [];

      if (schema.type === 'string') {
        if (schema.enum) {
          return (
            <Form.Item
              key={key}
              name={key}
              label={schema.title || key}
              rules={rules}
            >
              <Select placeholder={schema.description}>
                {schema.enum.map((option: string) => (
                  <Select.Option key={option} value={option}>
                    {option}
                  </Select.Option>
                ))}
              </Select>
            </Form.Item>
          );
        } else {
          return (
            <Form.Item
              key={key}
              name={key}
              label={schema.title || key}
              rules={rules}
            >
              <Input placeholder={schema.description} />
            </Form.Item>
          );
        }
      } else if (schema.type === 'number' || schema.type === 'integer') {
        return (
          <Form.Item
            key={key}
            name={key}
            label={schema.title || key}
            rules={rules}
          >
            <InputNumber
              placeholder={schema.description}
              min={schema.minimum}
              max={schema.maximum}
              style={{ width: '100%' }}
            />
          </Form.Item>
        );
      } else if (schema.type === 'boolean') {
        return (
          <Form.Item
            key={key}
            name={key}
            label={schema.title || key}
            valuePropName="checked"
          >
            <Switch />
          </Form.Item>
        );
      } else {
        return (
          <Form.Item
            key={key}
            name={key}
            label={schema.title || key}
            rules={rules}
          >
            <TextArea placeholder={schema.description} rows={3} />
          </Form.Item>
        );
      }
    });
  };

  // 提供商表格列定义
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
      dataIndex: 'healthy',
      key: 'healthy',
      render: (healthy: boolean) => (
        <Tag color={healthy ? 'green' : 'red'} icon={healthy ? <CheckCircleOutlined /> : <ExclamationCircleOutlined />}>
          {healthy ? '正常' : '异常'}
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
              setSelectedProvider(record.type);
              setApiKeyModalVisible(true);
            }}
          >
            设置密钥
          </Button>
          <Button
            size="small"
            icon={<SettingOutlined />}
            onClick={() => setSelectedProvider(record.type)}
          >
            管理模型
          </Button>
        </Space>
      ),
    },
  ];

  // 模型表格列定义
  const modelColumns: ColumnsType<ModelInfo> = [
    {
      title: '模型名称',
      dataIndex: 'name',
      key: 'name',
      render: (name: string) => <span style={{ fontWeight: 'bold' }}>{name}</span>,
    },
    {
      title: '显示名称',
      dataIndex: 'display_name',
      key: 'display_name',
    },
    {
      title: '最大令牌',
      dataIndex: 'max_tokens',
      key: 'max_tokens',
      render: (tokens: number) => tokens.toLocaleString(),
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

  // MCP工具表格列定义
  const toolColumns: ColumnsType<MCPTool> = [
    {
      title: '工具名称',
      dataIndex: 'name',
      key: 'name',
      render: (name: string) => (
        <Space>
          <ToolOutlined />
          <span style={{ fontWeight: 'bold' }}>{name}</span>
        </Space>
      ),
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: '参数',
      dataIndex: 'inputSchema',
      key: 'inputSchema',
      render: (schema: any) => {
        if (!schema || !schema.properties) {
          return <Tag color="blue">无参数</Tag>;
        }
        const paramCount = Object.keys(schema.properties).length;
        return <Tag color="green">{paramCount} 个参数</Tag>;
      },
    },
    {
      title: '操作',
      key: 'actions',
      render: (_, record: MCPTool) => (
        <Button
          type="primary"
          size="small"
          icon={<PlayCircleOutlined />}
          onClick={() => {
            setSelectedTool(record);
            setExecuteModalVisible(true);
          }}
        >
          执行
        </Button>
      ),
    },
  ];





  const mcpSettings = (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
        <Title level={4} style={{ margin: 0 }}>
          <ToolOutlined style={{ marginRight: '8px' }} />
          MCP工具系统
        </Title>
        <Space>
          {!mcpInitialized && (
            <Button
              type="primary"
              icon={<SettingOutlined />}
              onClick={handleInitializeMCP}
              loading={mcpLoading}
            >
              初始化MCP
            </Button>
          )}
          <Button
            icon={<ReloadOutlined />}
            onClick={handleRefreshMCPTools}
            loading={mcpLoading}
            disabled={!mcpInitialized}
          >
            刷新工具
          </Button>
        </Space>
      </div>

      {/* 状态统计 */}
      <Row gutter={[16, 16]} style={{ marginBottom: '24px' }}>
        <Col xs={24} sm={8}>
          <Card>
            <Statistic
              title="初始化状态"
              value={mcpInitialized ? '已初始化' : '未初始化'}
              prefix={mcpInitialized ? <CheckCircleOutlined /> : <ExclamationCircleOutlined />}
              valueStyle={{ color: mcpInitialized ? '#3f8600' : '#cf1322' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card>
            <Statistic
              title="可用工具"
              value={(mcpTools || []).length}
              prefix={<ToolOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card>
            <Statistic
              title="日志条数"
              value={(mcpLogs || []).length}
              prefix={<HistoryOutlined />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
      </Row>

      {!mcpInitialized ? (
        <Card>
          <div style={{ textAlign: 'center', padding: '40px' }}>
            <ExclamationCircleOutlined style={{ fontSize: '48px', color: '#faad14', marginBottom: '16px' }} />
            <h3>MCP工具系统未初始化</h3>
            <p style={{ color: '#8c8c8c', marginBottom: '24px' }}>
              请先初始化MCP工具系统以使用相关功能
            </p>
            <Button
              type="primary"
              size="large"
              icon={<SettingOutlined />}
              onClick={handleInitializeMCP}
              loading={mcpLoading}
            >
              立即初始化
            </Button>
          </div>
        </Card>
      ) : (
        <>
          {/* 工具列表 */}
          <Card title="可用工具" style={{ marginBottom: '24px' }}>
            <Table
              columns={toolColumns}
              dataSource={mcpTools}
              rowKey="name"
              loading={mcpLoading}
              pagination={{
                pageSize: 10,
                showSizeChanger: true,
                showQuickJumper: true,
                showTotal: (total) => `共 ${total} 个工具`,
              }}
            />
          </Card>

          {/* 执行日志 */}
          <Card title="执行日志" style={{ marginBottom: '24px' }}>
            <List
              dataSource={mcpLogs}
              loading={mcpLoading}
              renderItem={(log: MCPMessage) => (
                <List.Item>
                  <List.Item.Meta
                    avatar={<CodeOutlined style={{ color: '#1890ff' }} />}
                    title={
                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <span>{log.message || '系统消息'}</span>
                        <Text type="secondary" style={{ fontSize: '12px' }}>
                          {new Date(log.timestamp).toLocaleString()}
                        </Text>
                      </div>
                    }
                    description={
                      <div>
                        <Paragraph ellipsis={{ rows: 2, expandable: true }}>
                          {log.data ? (typeof log.data === 'object' ? JSON.stringify(log.data, null, 2) : String(log.data)) : '-'}
                        </Paragraph>
                        <Tag color={log.level === 'error' ? 'red' : log.level === 'warn' ? 'orange' : 'blue'}>
                          {log.level.toUpperCase()}
                        </Tag>
                      </div>
                    }
                  />
                </List.Item>
              )}
              pagination={{
                pageSize: 10,
                showSizeChanger: false,
                showQuickJumper: false,
              }}
            />
          </Card>

          {/* 基础设置 */}
          <Card title="MCP工具设置" style={{ marginBottom: 16 }}>
            <Form.Item
              name={['mcp', 'autoInitialize']}
              label="自动初始化"
              valuePropName="checked"
              tooltip="启动时自动初始化MCP工具"
            >
              <Switch />
            </Form.Item>

            <Form.Item
              name={['mcp', 'timeout']}
              label="超时时间 (秒)"
              tooltip="MCP工具执行超时时间"
            >
              <InputNumber min={1} max={300} style={{ width: '100%' }} />
            </Form.Item>

            <Form.Item
              name={['mcp', 'maxRetries']}
              label="最大重试次数"
              tooltip="工具执行失败时的最大重试次数"
            >
              <InputNumber min={0} max={10} style={{ width: '100%' }} />
            </Form.Item>

            <Form.Item
              name={['mcp', 'enableLogging']}
              label="启用日志"
              valuePropName="checked"
              tooltip="记录MCP工具的执行日志"
            >
              <Switch />
            </Form.Item>
          </Card>
        </>
      )}

      {mcpError && (
        <Alert
          message="错误"
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
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
        <div>
          <Title level={4} style={{ margin: 0 }}>
            <CloudOutlined style={{ marginRight: '8px' }} />
            AI大模型管理
          </Title>
          <Text type="secondary">配置和管理AI服务提供商</Text>
        </div>
        <Button
          type="primary"
          icon={<ReloadOutlined />}
          onClick={handleRefreshProviders}
          loading={providersLoading}
        >
          刷新状态
        </Button>
      </div>



      {/* 提供商列表 */}
      <Card title="AI提供商列表" style={{ marginBottom: '24px' }}>
        <Table
          columns={providerColumns}
          dataSource={providers}
          rowKey="name"
          loading={providersLoading}
          pagination={false}
        />
      </Card>

      {/* 模型管理 */}
      {selectedProvider && (
        <Card
          title={
            <Space>
              <span>模型管理 - {getProviderDisplayName(selectedProvider)}</span>
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
            rowKey="name"
            loading={providersLoading}
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
          apiKeyForm.resetFields();
        }}
        footer={null}
      >
        <Form
          form={apiKeyForm}
          layout="vertical"
          onFinish={handleSetAPIKey}
        >
          <Form.Item
            name="apiKey"
            label="API密钥"
            rules={[{ required: true, message: '请输入API密钥' }]}
          >
            <Input.Password placeholder="请输入API密钥" />
          </Form.Item>
          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit">
                保存
              </Button>
              <Button onClick={() => {
                setApiKeyModalVisible(false);
                apiKeyForm.resetFields();
              }}>
                取消
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );

  const aboutInfo = (
    <Card title="关于" style={{ marginBottom: 16 }}>
      <List>
        <List.Item>
          <List.Item.Meta
            title="应用版本"
            description="Go-SpringAI v1.0.0"
          />
        </List.Item>
        <List.Item>
          <List.Item.Meta
            title="构建时间"
            description={new Date().toLocaleDateString()}
          />
        </List.Item>
        <List.Item>
          <List.Item.Meta
            title="支持的提供商"
            description={
              <Space wrap>
                <Tag color="blue">OpenAI</Tag>
                <Tag color="green">Anthropic</Tag>
                <Tag color="orange">Google</Tag>
                <Tag color="purple">Azure</Tag>
              </Space>
            }
          />
        </List.Item>
        <List.Item>
          <List.Item.Meta
            title="技术栈"
            description={
              <Space wrap>
                <Tag>React</Tag>
                <Tag>TypeScript</Tag>
                <Tag>Ant Design</Tag>
                <Tag>Redux Toolkit</Tag>
                <Tag>Go</Tag>
              </Space>
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
            设置
          </Title>
          <Text type="secondary">配置应用程序的各项设置</Text>
        </div>
        <Space>
          <Button
            icon={<ReloadOutlined />}
            onClick={handleReset}
            disabled={isLoading}
          >
            重置
          </Button>
          <Button
            type="primary"
            icon={<SaveOutlined />}
            onClick={handleSave}
            loading={isLoading}
            disabled={!hasChanges}
          >
            保存设置
          </Button>
        </Space>
      </div>



      <Form
        form={form}
        layout="vertical"
        onValuesChange={handleFormChange}
        initialValues={settings}
      >
        <Tabs defaultActiveKey="providers" type="card">
          <TabPane tab={<span><CloudOutlined />AI大模型管理</span>} key="providers">
            {providersManagement}
          </TabPane>
          <TabPane tab="MCP工具" key="mcp">
            {mcpSettings}
          </TabPane>
          <TabPane tab="关于" key="about">
            {aboutInfo}
          </TabPane>
        </Tabs>
      </Form>
    </div>
  );
};

export default SettingsPage;