import React, { useEffect, useState } from 'react';
import {
  Card,
  Table,
  Button,
  Space,
  Modal,
  Form,
  Input,
  Select,
  message,
  Tag,
  Divider,
  Row,
  Col,
  Statistic,
  List,
  Typography,
  Collapse,
  Alert,
} from 'antd';
import {
  ToolOutlined,
  PlayCircleOutlined,
  ReloadOutlined,
  SettingOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  CodeOutlined,
  HistoryOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { useAppDispatch, useAppSelector } from '../store';
import {
  initializeMCP,
  fetchMCPTools,
  executeMCPTool,
  fetchMCPLogs,
} from '../store/slices/mcpSlice';
import type { MCPTool, MCPMessage } from '../types/api';

const { TextArea } = Input;
const { Text, Paragraph } = Typography;
const { Panel } = Collapse;

const MCPToolsPage: React.FC = () => {
  const dispatch = useAppDispatch();
  const {
    tools,
    logs,
    isInitialized,
    isLoading,
    error,
    executionResults,
  } = useAppSelector(state => state.mcp);

  const [executeModalVisible, setExecuteModalVisible] = useState(false);
  const [selectedTool, setSelectedTool] = useState<MCPTool | null>(null);
  const [form] = Form.useForm();

  useEffect(() => {
    if (!isInitialized) {
      dispatch(initializeMCP());
    } else {
      dispatch(fetchMCPTools());
      dispatch(fetchMCPLogs());
    }
  }, [dispatch, isInitialized]);

  const handleInitializeMCP = () => {
    dispatch(initializeMCP());
  };

  const handleRefreshTools = () => {
    dispatch(fetchMCPTools());
    dispatch(fetchMCPLogs());
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
      <Collapse size="small">
        <Panel header="参数结构" key="schema">
          <pre style={{ fontSize: '12px', backgroundColor: '#f5f5f5', padding: '8px', borderRadius: '4px' }}>
            {JSON.stringify(schema, null, 2)}
          </pre>
        </Panel>
      </Collapse>
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
      dataIndex: 'input_schema',
      key: 'input_schema',
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
                    {JSON.stringify(record.input_schema, null, 2)}
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
        <Space>
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
          <Button
            icon={<ReloadOutlined />}
            onClick={handleRefreshTools}
            loading={isLoading}
            disabled={!isInitialized}
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
              value={isInitialized ? '已初始化' : '未初始化'}
              prefix={isInitialized ? <CheckCircleOutlined /> : <ExclamationCircleOutlined />}
              valueStyle={{ color: isInitialized ? '#3f8600' : '#cf1322' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card>
            <Statistic
              title="可用工具"
              value={(tools || []).length}
              prefix={<ToolOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card>
            <Statistic
              title="日志条数"
              value={(logs || []).length}
              prefix={<HistoryOutlined />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
      </Row>

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
          <Card title="可用工具" style={{ marginBottom: '24px' }}>
            <Table
              columns={toolColumns}
              dataSource={tools}
              rowKey="name"
              loading={isLoading}
              pagination={{
                pageSize: 10,
                showSizeChanger: true,
                showQuickJumper: true,
                showTotal: (total) => `共 ${total} 个工具`,
              }}
            />
          </Card>

          {/* 执行日志 */}
          <Card title="执行日志" extra={
            <Button size="small" onClick={() => dispatch(fetchMCPLogs())}>
              刷新日志
            </Button>
          }>
            <List
              dataSource={logs.slice(0, 50)} // 只显示最近50条
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
                          <Collapse size="small" style={{ marginTop: '8px' }}>
                            <Panel header="详细数据" key="data">
                              <pre style={{ fontSize: '12px', backgroundColor: '#f5f5f5', padding: '8px', borderRadius: '4px' }}>
                                {JSON.stringify(log.data, null, 2)}
                              </pre>
                            </Panel>
                          </Collapse>
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
      <Modal
        title={`执行工具: ${selectedTool?.name}`}
        open={executeModalVisible}
        onCancel={() => {
          setExecuteModalVisible(false);
          form.resetFields();
          setSelectedTool(null);
        }}
        footer={null}
        width={600}
      >
        {selectedTool && (
          <div>
            <Paragraph>{selectedTool.description}</Paragraph>
            <Divider />
            <Form
              form={form}
              layout="vertical"
              onFinish={handleExecuteTool}
            >
              {renderExecutionForm(selectedTool.input_schema)}
              
              <Form.Item style={{ marginBottom: 0, marginTop: '24px' }}>
                <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
                  <Button onClick={() => {
                    setExecuteModalVisible(false);
                    form.resetFields();
                    setSelectedTool(null);
                  }}>
                    取消
                  </Button>
                  <Button type="primary" htmlType="submit" loading={isLoading}>
                    执行工具
                  </Button>
                </Space>
              </Form.Item>
            </Form>
          </div>
        )}
      </Modal>

      {error && (
        <div style={{ marginTop: '16px', padding: '16px', backgroundColor: '#fff2f0', border: '1px solid #ffccc7', borderRadius: '6px' }}>
          <span style={{ color: '#ff4d4f' }}>错误: {error}</span>
        </div>
      )}
    </div>
  );
};

export default MCPToolsPage;