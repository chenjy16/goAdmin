import React, { useEffect, useState } from 'react';
import {
  Card,
  Button,
  Space,
  Modal,
  Form,
  Input,
  Select,
  message,
  Tag,
  Divider,
  List,
  Typography,
  Collapse,
  Alert,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';

import {
  ToolOutlined,
  PlayCircleOutlined,
  SettingOutlined,
  ExclamationCircleOutlined,
  CodeOutlined,
} from '@ant-design/icons';
import { useAppDispatch, useAppSelector } from '../store';
import {
  initializeMCP,
  fetchMCPTools,
  executeMCPTool,
  fetchMCPLogs,
  checkMCPStatus,
} from '../store/slices/mcpSlice';
import type { MCPTool, MCPMessage } from '../types/api';
import { SearchableTable, FormModal } from '../components/common';
import { useAsyncOperation } from '../hooks';
import { 
  createTextColumn, 
  createActionColumn,
  mergeColumns 
} from '../utils/tableColumns';
import { useTranslation } from 'react-i18next';

const { TextArea } = Input;
const { Text, Paragraph } = Typography;

const MCPToolsPage: React.FC = () => {
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
  const {
    tools,
    logs,
    isInitialized,
    isLoading,
    error,
  } = useAppSelector(state => state.mcp);

  const [executeModalVisible, setExecuteModalVisible] = useState(false);
  const [selectedTool, setSelectedTool] = useState<MCPTool | null>(null);
  const [form] = Form.useForm();

  // 检查MCP状态并加载数据
  useEffect(() => {
    // 首先检查MCP状态
    dispatch(checkMCPStatus()).then((result) => {
      if (checkMCPStatus.fulfilled.match(result)) {
        // 如果已初始化，加载工具和日志
        if (result.payload.initialized) {
          dispatch(fetchMCPTools());
          dispatch(fetchMCPLogs());
        }
      }
    });
  }, [dispatch]);

  // 当MCP状态变为已初始化时，加载数据
  useEffect(() => {
    if (isInitialized) {
      dispatch(fetchMCPTools());
      dispatch(fetchMCPLogs());
    }
  }, [dispatch, isInitialized]);

  const initializeOperation = useAsyncOperation(
    () => dispatch(initializeMCP()),
    {
      successMessage: t('mcpTools.initializeSuccess'),
      errorMessage: t('mcpTools.initializeFailed')
    }
  );

  const executeToolOperation = useAsyncOperation(
    async (toolName: string, args: Record<string, any>) => {
      await dispatch(executeMCPTool({
        name: toolName,
        arguments: args,
      })).unwrap();
      
      setExecuteModalVisible(false);
      form.resetFields();
      dispatch(fetchMCPLogs()); // 刷新日志
    },
    {
      successMessage: t('mcpTools.executeSuccess'),
      errorMessage: t('mcpTools.executeFailed')
    }
  );

  const handleInitializeMCP = () => {
    initializeOperation.execute();
  };

  const handleExecuteTool = async (values: Record<string, any>) => {
    if (!selectedTool) return;
    await executeToolOperation.execute(selectedTool.name, values);
  };

  const renderInputSchema = (schema: any) => {
    if (!schema || !schema.properties) {
      return <Text type="secondary">{t('mcpTools.noParameters')}</Text>;
    }

    return (
      <Collapse 
        size="small"
        items={[
          {
            key: 'schema',
            label: t('mcpTools.parameterStructure'),
            children: (
              <pre style={{ fontSize: '12px', backgroundColor: '#f5f5f5', padding: '8px', borderRadius: '4px' }}>
                {JSON.stringify(schema, null, 2)}
              </pre>
            )
          }
        ]}
      />
    );
  };

  const renderExecutionForm = (schema: any) => {
    if (!schema || !schema.properties) {
      return <Alert message={t('mcpTools.noParametersRequired')} type="info" />;
    }

    const properties = schema.properties;
    const required = schema.required || [];

    return (
      <div>
        {Object.entries(properties).map(([key, prop]: [string, any]) => (
          <Form.Item
            key={key}
            label={prop.title || key}
            name={key}
            rules={[
              {
                required: required.includes(key),
                message: `${t('common.pleaseEnter')}${prop.title || key}`,
              },
            ]}
            tooltip={prop.description}
          >
            {prop.type === 'string' && prop.enum ? (
              <Select
                placeholder={`${t('common.pleaseSelect')}${prop.title || key}`}
                options={prop.enum.map((value: string) => ({
                  label: value,
                  value,
                }))}
              />
            ) : prop.type === 'boolean' ? (
              <Select
                placeholder={`${t('common.pleaseSelect')}${prop.title || key}`}
                options={[
                  { label: t('common.yes'), value: true },
                  { label: t('common.no'), value: false },
                ]}
              />
            ) : prop.type === 'number' || prop.type === 'integer' ? (
              <Input
                type="number"
                placeholder={prop.description || `${t('common.pleaseEnter')}${prop.title || key}`}
              />
            ) : (
              <TextArea
                placeholder={prop.description || `${t('common.pleaseEnter')}${prop.title || key}`}
                autoSize={{ minRows: 2, maxRows: 6 }}
              />
            )}
          </Form.Item>
        ))}
      </div>
    );
  };

  const toolColumns = mergeColumns<MCPTool>([
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
    {
      title: t('mcpTools.parameterStructure'),
      dataIndex: 'inputSchema',
      key: 'inputSchema',
      render: (schema: any) => renderInputSchema(schema),
    },
    createActionColumn({
      actions: [
        {
          key: 'execute',
          label: t('mcpTools.execute'),
          type: 'primary',
          icon: <PlayCircleOutlined />,
          onClick: (record: MCPTool) => {
            setSelectedTool(record);
            setExecuteModalVisible(true);
          }
        },
        {
          key: 'view-schema',
          label: t('mcpTools.viewStructure'),
          icon: <CodeOutlined />,
          onClick: (record: MCPTool) => {
            Modal.info({
              title: `${record.name} - ${t('mcpTools.parameterStructure')}`,
              content: (
                <pre style={{ fontSize: '12px', backgroundColor: '#f5f5f5', padding: '16px', borderRadius: '4px' }}>
                  {JSON.stringify(record.inputSchema, null, 2)}
                </pre>
              ),
              width: 600,
            });
          }
        }
      ]
    })
  ]);

  const getLogLevelColor = (level: string) => {
    switch (level) {
      case 'error':
        return 'red';
      case 'warn':
        return 'orange';
      case 'info':
        return 'blue';
      case 'debug':
        return 'gray';
      default:
        return 'default';
    }
  };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
        <h1>{t('mcpTools.systemTitle')}</h1>
        {!isInitialized && (
          <Button
            type="primary"
            icon={<SettingOutlined />}
            onClick={handleInitializeMCP}
            loading={isLoading}
          >
            {t('mcpTools.initializeMCP')}
          </Button>
        )}
      </div>



      {!isInitialized ? (
        <Card
          title={
            <Space>
              <ToolOutlined />
              {t('mcpTools.title')}
            </Space>
          }
          extra={
            <Button
              type="primary"
              icon={<SettingOutlined />}
              onClick={() => dispatch(initializeMCP())}
              loading={isLoading}
            >
              {t('mcpTools.initializeMCP')}
            </Button>
          }
        >
          <div style={{ textAlign: 'center', padding: '40px 0' }}>
            <ExclamationCircleOutlined style={{ fontSize: '48px', color: '#faad14', marginBottom: '16px' }} />
            <h3>{t('mcpTools.notInitialized')}</h3>
            <p style={{ color: '#666', marginBottom: '24px' }}>
              {t('mcpTools.initializePrompt')}
            </p>
            <Button
              type="primary"
              size="large"
              icon={<SettingOutlined />}
              onClick={() => dispatch(initializeMCP())}
              loading={isLoading}
            >
              {t('mcpTools.initializeNow')}
            </Button>
          </div>
        </Card>
      ) : (
        <>
          {/* 工具列表 */}
          <SearchableTable<MCPTool>
            columns={toolColumns}
            dataSource={tools}
            rowKey="name"
            loading={isLoading}
            searchFields={['name', 'description']}
            searchPlaceholder={t('mcpTools.searchTools')}
            showRefresh={true}
            onRefresh={() => dispatch(fetchMCPTools())}
            refreshLoading={isLoading}
            title={t('mcpTools.availableTools')}
          />

          {/* 执行日志 */}
          <Card title={t('mcpTools.executionLogs')} extra={
            <Button size="small" onClick={() => dispatch(fetchMCPLogs())}>
              {t('mcpTools.refreshLogs')}
            </Button>
          }>
            <List
              dataSource={Array.isArray(logs) ? logs.slice(0, 50) : []} // 只显示最近50条，确保logs是数组
              renderItem={(log: MCPMessage) => (
                <List.Item>
                  <List.Item.Meta
                    title={
                      <Space>
                        <Tag color={getLogLevelColor(log.level)}>{log.level.toUpperCase()}</Tag>
                        <Text>{log.message}</Text>
                      </Space>
                    }
                    description={
                      <div>
                        <Text type="secondary" style={{ fontSize: '12px' }}>
                          {new Date(log.timestamp).toLocaleString()}
                        </Text>
                        {log.data && (
                          <Collapse 
                            size="small" 
                            style={{ marginTop: '8px' }}
                            items={[
                              {
                                key: 'data',
                                label: t('mcpTools.detailedData'),
                                children: (
                                  <pre style={{ fontSize: '12px', backgroundColor: '#f5f5f5', padding: '8px', borderRadius: '4px' }}>
                                    {JSON.stringify(log.data, null, 2)}
                                  </pre>
                                )
                              }
                            ]}
                          />
                        )}
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
        </>
      )}

      {/* 工具执行模态框 */}
      <FormModal
        title={`${t('mcpTools.executeTool')}: ${selectedTool ? getLocalizedToolName(selectedTool.name) : ''}`}
        open={executeModalVisible}
        onCancel={() => {
          setExecuteModalVisible(false);
          form.resetFields();
          setSelectedTool(null);
        }}
        onFinish={handleExecuteTool}
        form={form}
        loading={isLoading}
        width={600}
        okText={t('mcpTools.executeTool')}
        cancelText={t('common.cancel')}
      >
        {selectedTool && (
          <div>
            <Paragraph>{selectedTool.description}</Paragraph>
            <Divider />
            {renderExecutionForm(selectedTool.inputSchema)}
          </div>
        )}
      </FormModal>

      {error && (
        <div style={{ marginTop: '16px', padding: '16px', backgroundColor: '#fff2f0', border: '1px solid #ffccc7', borderRadius: '6px' }}>
          <span style={{ color: '#ff4d4f' }}>{t('common.error')}: {error}</span>
        </div>
      )}
    </div>
  );
};

export default MCPToolsPage;