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

import {
  ToolOutlined,
  PlayCircleOutlined,
  SettingOutlined,
  ExclamationCircleOutlined,
  CodeOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
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

const { TextArea } = Input;
const { Text, Paragraph } = Typography;

const MCPToolsPage: React.FC = () => {
  const dispatch = useAppDispatch();
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

  const handleInitializeMCP = () => {
    dispatch(initializeMCP());
  };

  const handleExecuteTool = async (values: Record<string, any>) => {
    if (!selectedTool) return;

    try {
      await dispatch(executeMCPTool({
        name: selectedTool.name,
        arguments: values,
      })).unwrap();
      
      message.success('工具执行成功');
      setExecuteModalVisible(false);
      form.resetFields();
      dispatch(fetchMCPLogs()); // 刷新日志
    } catch (err) {
      message.error('工具执行失败');
    }
  };

  const renderInputSchema = (schema: any) => {
    if (!schema || !schema.properties) {
      return <Text type="secondary">无参数</Text>;
    }

    return (
      <Collapse 
        size="small"
        items={[
          {
            key: 'schema',
            label: '参数结构',
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
      return <Alert message="此工具无需参数" type="info" />;
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
                message: `请输入${prop.title || key}`,
              },
            ]}
            tooltip={prop.description}
          >
            {prop.type === 'string' && prop.enum ? (
              <Select
                placeholder={`选择${prop.title || key}`}
                options={prop.enum.map((value: string) => ({
                  label: value,
                  value,
                }))}
              />
            ) : prop.type === 'boolean' ? (
              <Select
                placeholder={`选择${prop.title || key}`}
                options={[
                  { label: '是', value: true },
                  { label: '否', value: false },
                ]}
              />
            ) : prop.type === 'number' || prop.type === 'integer' ? (
              <Input
                type="number"
                placeholder={prop.description || `输入${prop.title || key}`}
              />
            ) : (
              <TextArea
                placeholder={prop.description || `输入${prop.title || key}`}
                autoSize={{ minRows: 2, maxRows: 6 }}
              />
            )}
          </Form.Item>
        ))}
      </div>
    );
  };

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
      title: '参数结构',
      dataIndex: 'inputSchema',
      key: 'inputSchema',
      render: (schema: any) => renderInputSchema(schema),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_, record: MCPTool) => (
        <Space>
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
          <Button
            size="small"
            icon={<CodeOutlined />}
            onClick={() => {
              Modal.info({
                title: `${record.name} - 参数结构`,
                content: (
                  <pre style={{ fontSize: '12px', backgroundColor: '#f5f5f5', padding: '16px', borderRadius: '4px' }}>
                    {JSON.stringify(record.inputSchema, null, 2)}
                  </pre>
                ),
                width: 600,
              });
            }}
          >
            查看结构
          </Button>
        </Space>
      ),
    },
  ];

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
        <h1>MCP工具系统</h1>
        {!isInitialized && (
          <Button
            type="primary"
            icon={<SettingOutlined />}
            onClick={handleInitializeMCP}
            loading={isLoading}
          >
            初始化MCP
          </Button>
        )}
      </div>



      {!isInitialized ? (
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
              loading={isLoading}
            >
              立即初始化
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
            searchPlaceholder="搜索工具..."
            showRefresh={true}
            onRefresh={() => dispatch(fetchMCPTools())}
            refreshLoading={isLoading}
            title="可用工具"
          />

          {/* 执行日志 */}
          <Card title="执行日志" extra={
            <Button size="small" onClick={() => dispatch(fetchMCPLogs())}>
              刷新日志
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
                                label: '详细数据',
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
        title={`执行工具: ${selectedTool?.name}`}
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
        okText="执行工具"
        cancelText="取消"
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
          <span style={{ color: '#ff4d4f' }}>错误: {error}</span>
        </div>
      )}
    </div>
  );
};

export default MCPToolsPage;