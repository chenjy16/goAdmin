import React, { useEffect, useState, useRef } from 'react';
import {
  Layout,
  Input,
  Button,
  List,
  Avatar,
  Card,
  Space,
  Typography,
  Spin,
  message,
  Tag,
} from 'antd';
import {
  SendOutlined,
  UserOutlined,
  RobotOutlined,
  SettingOutlined,
  CopyOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
} from '@ant-design/icons';
import { useAppDispatch, useAppSelector } from '../store';
import {
  initializeAssistant,
  sendAssistantMessage,
  createAssistantConversation,
} from '../store/slices/assistantSlice';
import type { ChatMessage } from '../types/api';

const { TextArea } = Input;
const { Text, Paragraph } = Typography;
const { Content } = Layout;

const AssistantPage: React.FC = () => {
  const dispatch = useAppDispatch();
  const {
    conversations,
    currentConversationId,
    isInitialized,
    isLoading,
    error,
  } = useAppSelector(state => state.assistant);

  const [inputValue, setInputValue] = useState('');
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!isInitialized) {
      dispatch(initializeAssistant());
    } else if (!currentConversationId) {
      dispatch(createAssistantConversation('新助手对话'));
    }
  }, [dispatch, isInitialized, currentConversationId]);

  useEffect(() => {
    scrollToBottom();
  }, [conversations]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const currentConversation = (conversations || []).find(c => c.id === currentConversationId);
  const currentMessages = currentConversation?.messages || [];

  const handleSendMessage = async () => {
    if (!inputValue.trim() || isLoading || !isInitialized) return;

    const messageContent = inputValue.trim();
    setInputValue('');

    try {
      await dispatch(sendAssistantMessage({
        conversationId: currentConversationId!,
        message: {
          role: 'user',
          content: messageContent,
          timestamp: new Date().toISOString(),
        },
      })).unwrap();
    } catch (err) {
      message.error('发送消息失败');
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  const handleCopyMessage = (content: string) => {
    navigator.clipboard.writeText(content);
    message.success('已复制到剪贴板');
  };

  const handleNewConversation = () => {
    dispatch(createAssistantConversation('新助手对话'));
  };

  const handleInitializeAssistant = () => {
    dispatch(initializeAssistant());
  };

  const renderMessage = (msg: ChatMessage, index: number) => {
    const isUser = msg.role === 'user';

    return (
      <List.Item
        key={index}
        style={{
          padding: '16px',
          border: 'none',
          backgroundColor: isUser ? '#f6f8fa' : '#ffffff',
        }}
      >
        <List.Item.Meta
          avatar={
            <Avatar
              icon={isUser ? <UserOutlined /> : <RobotOutlined />}
              style={{
                backgroundColor: isUser ? '#1890ff' : '#52c41a',
              }}
            />
          }
          title={
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <Text strong>{isUser ? '用户' : 'AI助手'}</Text>
              <Space>
                <Button
                  type="text"
                  size="small"
                  icon={<CopyOutlined />}
                  onClick={() => handleCopyMessage(msg.content)}
                />
                {msg.timestamp && (
                  <Text type="secondary" style={{ fontSize: '12px' }}>
                    {new Date(msg.timestamp).toLocaleTimeString()}
                  </Text>
                )}
              </Space>
            </div>
          }
          description={
            <div style={{ marginTop: '8px' }}>
              <Paragraph
                style={{
                  marginBottom: 0,
                  whiteSpace: 'pre-wrap',
                  wordBreak: 'break-word',
                }}
              >
                {msg.content}
              </Paragraph>
              {msg.tool_calls && msg.tool_calls.length > 0 && (
                <div style={{ marginTop: '8px' }}>
                  <Tag color="blue">使用了工具调用</Tag>
                </div>
              )}
            </div>
          }
        />
      </List.Item>
    );
  };

  if (!isInitialized) {
    return (
      <div style={{ height: 'calc(100vh - 64px)', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <Card style={{ textAlign: 'center', maxWidth: '400px' }}>
          <ExclamationCircleOutlined style={{ fontSize: '48px', color: '#faad14', marginBottom: '16px' }} />
          <h3>AI助手未初始化</h3>
          <p style={{ color: '#8c8c8c', marginBottom: '24px' }}>
            请先初始化AI助手服务以开始对话
          </p>
          <Button
            type="primary"
            size="large"
            icon={<SettingOutlined />}
            onClick={handleInitializeAssistant}
            loading={isLoading}
          >
            初始化助手
          </Button>
          {error && (
            <div style={{ marginTop: '16px', color: '#ff4d4f', fontSize: '14px' }}>
              初始化失败: {error}
            </div>
          )}
        </Card>
      </div>
    );
  }

  return (
    <Layout style={{ height: 'calc(100vh - 64px)' }}>
      <Content style={{ display: 'flex', flexDirection: 'column', padding: 0 }}>
        {/* 助手头部 */}
        <div
          style={{
            padding: '16px 24px',
            borderBottom: '1px solid #f0f0f0',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
          }}
        >
          <div>
            <Space>
              <RobotOutlined style={{ color: '#52c41a' }} />
              <Text strong style={{ fontSize: '16px' }}>
                AI智能助手
              </Text>
              <Tag color="green" icon={<CheckCircleOutlined />}>
                已初始化
              </Tag>
            </Space>
            {currentConversation && (
              <Text type="secondary" style={{ marginLeft: '16px' }}>
                {currentConversation.messages.length} 条消息
              </Text>
            )}
          </div>
          <Space>
            <Button onClick={handleNewConversation}>新建对话</Button>
          </Space>
        </div>

        {/* 消息列表 */}
        <div style={{ flex: 1, overflow: 'auto', padding: '0 24px' }}>
          {currentMessages.length === 0 ? (
            <div
              style={{
                height: '100%',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                flexDirection: 'column',
              }}
            >
              <RobotOutlined style={{ fontSize: '48px', color: '#d9d9d9', marginBottom: '16px' }} />
              <Text type="secondary">AI助手已准备就绪，开始对话吧！</Text>
              <div style={{ marginTop: '16px', textAlign: 'center' }}>
                <Text type="secondary" style={{ fontSize: '12px' }}>
                  AI助手具备以下能力：
                </Text>
                <div style={{ marginTop: '8px' }}>
                  <Tag>智能对话</Tag>
                  <Tag>工具调用</Tag>
                  <Tag>代码生成</Tag>
                  <Tag>问题解答</Tag>
                </div>
              </div>
            </div>
          ) : (
            <List
              dataSource={currentMessages}
              renderItem={renderMessage}
              style={{ padding: '16px 0' }}
            />
          )}
          {isLoading && (
            <div style={{ textAlign: 'center', padding: '16px' }}>
              <Spin tip="AI助手正在思考中..." />
            </div>
          )}
          <div ref={messagesEndRef} />
        </div>

        {/* 输入区域 */}
        <div
          style={{
            padding: '16px 24px',
            borderTop: '1px solid #f0f0f0',
            backgroundColor: '#fafafa',
          }}
        >
          <div style={{ display: 'flex', gap: '8px' }}>
            <TextArea
              value={inputValue}
              onChange={(e) => setInputValue(e.target.value)}
              onKeyPress={handleKeyPress}
              placeholder="向AI助手提问... (Shift+Enter 换行，Enter 发送)"
              autoSize={{ minRows: 1, maxRows: 4 }}
              style={{ flex: 1 }}
              disabled={isLoading || !isInitialized}
            />
            <Button
              type="primary"
              icon={<SendOutlined />}
              onClick={handleSendMessage}
              loading={isLoading}
              disabled={!inputValue.trim() || !isInitialized}
            >
              发送
            </Button>
          </div>
          {error && (
            <div style={{ marginTop: '8px', color: '#ff4d4f', fontSize: '12px' }}>
              {error}
            </div>
          )}
        </div>
      </Content>
    </Layout>
  );
};

export default AssistantPage;