import React, { useEffect } from 'react';
import { Card, Row, Col, Statistic, Progress, List, Tag, Button } from 'antd';
import {
  MessageOutlined,
  CloudOutlined,
  ToolOutlined,
  RobotOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  SyncOutlined,
  SettingOutlined,
} from '@ant-design/icons';
import { useAppDispatch, useAppSelector } from '../store';
import { fetchProviders } from '../store/slices/providersSlice';
import { fetchMCPTools } from '../store/slices/mcpSlice';

const DashboardPage: React.FC = () => {
  const dispatch = useAppDispatch();
  const { providers, isLoading: providersLoading } = useAppSelector(state => state.providers);
  const { tools, isInitialized: mcpInitialized } = useAppSelector(state => state.mcp);
  const { conversations } = useAppSelector(state => state.chat);
  const { conversations: assistantConversations } = useAppSelector(state => state.assistant);

  useEffect(() => {
    dispatch(fetchProviders());
    dispatch(fetchMCPTools());
  }, [dispatch]);

  // 计算统计数据
  const healthyProviders = (providers || []).filter(p => p.health).length;
  const totalModels = (providers || []).reduce((sum, p) => sum + p.model_count, 0);
  const totalConversations = (conversations || []).length + (assistantConversations || []).length;
  const availableTools = (tools || []).length;

  // 最近活动数据
  const recentActivities = [
    {
      title: '新建聊天对话',
      description: '与 GPT-4 开始新的对话',
      time: '2分钟前',
      type: 'chat',
    },
    {
      title: 'MCP工具执行',
      description: '执行了文件搜索工具',
      time: '5分钟前',
      type: 'tool',
    },
    {
      title: '提供商状态更新',
      description: 'OpenAI 连接状态正常',
      time: '10分钟前',
      type: 'provider',
    },
    {
      title: 'AI助手初始化',
      description: '助手服务启动完成',
      time: '15分钟前',
      type: 'assistant',
    },
  ];

  const getActivityIcon = (type: string) => {
    switch (type) {
      case 'chat':
        return <MessageOutlined style={{ color: '#1890ff' }} />;
      case 'tool':
        return <ToolOutlined style={{ color: '#52c41a' }} />;
      case 'provider':
        return <CloudOutlined style={{ color: '#722ed1' }} />;
      case 'assistant':
        return <RobotOutlined style={{ color: '#fa8c16' }} />;
      default:
        return <SyncOutlined style={{ color: '#8c8c8c' }} />;
    }
  };

  return (
    <div>
      <h1 style={{ marginBottom: '24px' }}>仪表板</h1>
      
      {/* 统计卡片 */}
      <Row gutter={[16, 16]} style={{ marginBottom: '24px' }}>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="AI提供商"
              value={(providers || []).length}
              prefix={<CloudOutlined />}
              suffix={
                <Tag color={healthyProviders === (providers || []).length ? 'green' : 'orange'}>
                  {healthyProviders}/{(providers || []).length} 健康
                </Tag>
              }
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="可用模型"
              value={totalModels}
              prefix={<RobotOutlined />}
              valueStyle={{ color: '#3f8600' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="聊天对话"
              value={totalConversations}
              prefix={<MessageOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="MCP工具"
              value={availableTools}
              prefix={<ToolOutlined />}
              suffix={
                <Tag color={mcpInitialized ? 'green' : 'red'}>
                  {mcpInitialized ? '已初始化' : '未初始化'}
                </Tag>
              }
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        {/* 系统状态 */}
        <Col xs={24} lg={12}>
          <Card title="系统状态" extra={<Button type="link">查看详情</Button>}>
            <div style={{ marginBottom: '16px' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '8px' }}>
                <span>API连接状态</span>
                <span>{healthyProviders}/{(providers || []).length}</span>
              </div>
              <Progress 
                percent={(providers || []).length > 0 ? Math.round((healthyProviders / (providers || []).length) * 100) : 0}
                status={healthyProviders === (providers || []).length ? 'success' : 'active'}
              />
            </div>
            
            <div style={{ marginBottom: '16px' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '8px' }}>
                <span>MCP工具状态</span>
                <span>{mcpInitialized ? '正常' : '异常'}</span>
              </div>
              <Progress 
                percent={mcpInitialized ? 100 : 0}
                status={mcpInitialized ? 'success' : 'exception'}
              />
            </div>

            <div>
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '8px' }}>
                <span>服务可用性</span>
                <span>99.9%</span>
              </div>
              <Progress percent={99.9} status="success" />
            </div>
          </Card>
        </Col>

        {/* 最近活动 */}
        <Col xs={24} lg={12}>
          <Card title="最近活动" extra={<Button type="link">查看全部</Button>}>
            <List
              dataSource={recentActivities}
              renderItem={(item) => (
                <List.Item>
                  <List.Item.Meta
                    avatar={getActivityIcon(item.type)}
                    title={item.title}
                    description={
                      <div>
                        <div>{item.description}</div>
                        <div style={{ color: '#8c8c8c', fontSize: '12px', marginTop: '4px' }}>
                          {item.time}
                        </div>
                      </div>
                    }
                  />
                </List.Item>
              )}
            />
          </Card>
        </Col>
      </Row>

      {/* 快速操作 */}
      <Row gutter={[16, 16]} style={{ marginTop: '24px' }}>
        <Col span={24}>
          <Card title="快速操作">
            <Row gutter={[16, 16]}>
              <Col xs={24} sm={12} md={6}>
                <Button 
                  type="primary" 
                  icon={<MessageOutlined />} 
                  size="large" 
                  block
                  onClick={() => window.location.href = '/chat'}
                >
                  开始聊天
                </Button>
              </Col>
              <Col xs={24} sm={12} md={6}>
                <Button 
                  icon={<CloudOutlined />} 
                  size="large" 
                  block
                  onClick={() => window.location.href = '/providers'}
                >
                  管理提供商
                </Button>
              </Col>
              <Col xs={24} sm={12} md={6}>
                <Button 
                  icon={<ToolOutlined />} 
                  size="large" 
                  block
                  onClick={() => window.location.href = '/tools'}
                >
                  使用工具
                </Button>
              </Col>
              <Col xs={24} sm={12} md={6}>
                <Button 
                  icon={<SettingOutlined />} 
                  size="large" 
                  block
                  onClick={() => window.location.href = '/settings'}
                >
                  系统设置
                </Button>
              </Col>
            </Row>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default DashboardPage;