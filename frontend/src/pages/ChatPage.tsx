import React, { useState, useEffect, useRef } from 'react';
import {
  Layout,
  Input,
  Button,
  List,
  Avatar,
  Card,
  Select,
  Slider,
  Switch,
  Divider,
  Space,
  Typography,
  Spin,
  message,
} from 'antd';
import {
  SendOutlined,
  UserOutlined,
  RobotOutlined,
  SettingOutlined,
  DeleteOutlined,
  CopyOutlined,
} from '@ant-design/icons';
import { useAppDispatch, useAppSelector } from '../store';
import { sendMessage, createConversation, addMessage, updateStreamingMessage, clearStreamingMessage, clearError } from '../store/slices/chatSlice';
import { fetchProviders, fetchModels } from '../store/slices/providersSlice';
import { createChatSSE, sseManager } from '../services/sse';
import { handleError, handleSSEError } from '../services/errorHandler';
import type { ChatMessage, SSEEvent } from '../types/api';

const { TextArea } = Input;
const { Text, Paragraph } = Typography;
const { Sider, Content } = Layout;

const ChatPage: React.FC = () => {
  const dispatch = useAppDispatch();
  const {
    conversations,
    currentConversationId,
    isLoading,
    error,
    streamingMessage,
  } = useAppSelector(state => state.chat);
  const { providers, selectedProvider, models } = useAppSelector(state => state.providers);
  const { defaultProvider, defaultModel, chatSettings } = useAppSelector(state => state.settings);

  const [inputValue, setInputValue] = useState('');
  const [showSettings, setShowSettings] = useState(false);
  const [localSettings, setLocalSettings] = useState({
    temperature: chatSettings.temperature,
    maxTokens: chatSettings.maxTokens,
    stream: chatSettings.streamResponse,
  });
  const [isStreaming, setIsStreaming] = useState(false);

  const messagesEndRef = useRef<HTMLDivElement>(null);
  const sseRef = useRef<any>(null);

  useEffect(() => {
    dispatch(fetchProviders());
    if (!currentConversationId) {
      dispatch(createConversation({ title: '新对话', provider: defaultProvider, model: defaultModel[defaultProvider] || '' }));
    }
  }, [dispatch, currentConversationId, defaultProvider, defaultModel]);

  useEffect(() => {
    if (selectedProvider) {
      dispatch(fetchModels(selectedProvider));
    }
  }, [dispatch, selectedProvider]);

  useEffect(() => {
    scrollToBottom();
  }, [conversations, streamingMessage]);

  // 清理SSE连接
  useEffect(() => {
    return () => {
      if (sseRef.current) {
        sseRef.current.disconnect();
      }
      sseManager.disconnectAll();
    };
  }, []);

  // 当对话切换时重新建立SSE连接
  useEffect(() => {
    if (currentConversationId && localSettings.stream) {
      setupSSEConnection();
    }
    return () => {
      if (sseRef.current) {
        sseRef.current.disconnect();
        sseRef.current = null;
      }
    };
  }, [currentConversationId, localSettings.stream]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  // 设置SSE连接
  const setupSSEConnection = () => {
    if (!currentConversationId) return;

    try {
      sseRef.current = createChatSSE(
        currentConversationId,
        handleSSEMessage,
        handleSSEError
      );
      sseRef.current.connect();
    } catch (error) {
      handleError(error, {
        showMessage: true,
        showNotification: false,
      });
    }
  };

  // 处理SSE消息
  const handleSSEMessage = (event: SSEEvent) => {
    try {
      const eventType = event.event || 'message';
      let data: any = {};
      
      try {
        data = JSON.parse(event.data);
      } catch {
        data = { content: event.data };
      }

      switch (eventType) {
        case 'chat_message':
          if (data?.content) {
            dispatch(updateStreamingMessage(data.content));
          }
          break;
        case 'chat_done':
          if (data?.message) {
            // 添加完整的助手消息
            dispatch(addMessage({
              conversationId: currentConversationId!,
              message: {
                role: 'assistant',
                content: data.message.content,
                timestamp: new Date().toISOString(),
              },
            }));
          }
          dispatch(clearStreamingMessage());
          setIsStreaming(false);
          break;
        case 'chat_error':
          handleError(data?.error || '聊天过程中发生错误', {
            showMessage: true,
            showNotification: false,
          });
          dispatch(clearStreamingMessage());
          setIsStreaming(false);
          break;
        default:
          console.log('未知SSE事件类型:', eventType);
      }
    } catch (error) {
      handleError(error, {
        showMessage: true,
        showNotification: false,
      });
    }
  };

  const currentConversation = (conversations || []).find(c => c.id === currentConversationId);
  const currentMessages = currentConversation?.messages || [];

  const handleSendMessage = async () => {
    if (!inputValue.trim() || isLoading || isStreaming) return;

    const messageContent = inputValue.trim();
    setInputValue('');

    try {
      if (localSettings.stream) {
        // 流式响应模式
        setIsStreaming(true);
        dispatch(clearStreamingMessage());
        
        // 设置SSE连接来接收流式响应
        setupSSEConnection();
      }

      await dispatch(sendMessage({
        conversationId: currentConversationId!,
        message: {
          role: 'user',
          content: messageContent,
          timestamp: new Date().toISOString(),
        },
        provider: selectedProvider || defaultProvider,
        model: defaultModel[selectedProvider || defaultProvider] || '',
        maxTokens: localSettings.maxTokens,
        temperature: localSettings.temperature,
      })).unwrap();
    } catch (err) {
      handleError(err, {
        showMessage: true,
        showNotification: false,
      });
      setIsStreaming(false);
      dispatch(clearStreamingMessage());
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  const handleCopyMessage = async (content: string) => {
    try {
      await navigator.clipboard.writeText(content);
      message.success('消息已复制到剪贴板');
    } catch (error) {
      // 降级方案：使用传统的复制方法
      const textArea = document.createElement('textarea');
      textArea.value = content;
      textArea.style.position = 'fixed';
      textArea.style.opacity = '0';
      document.body.appendChild(textArea);
      textArea.select();
      try {
        document.execCommand('copy');
        message.success('消息已复制到剪贴板');
      } catch (fallbackError) {
        message.error('复制失败，请手动复制');
      }
      document.body.removeChild(textArea);
    }
  };

  const handleNewConversation = () => {
    dispatch(createConversation({ 
      title: '新对话', 
      provider: selectedProvider || defaultProvider, 
      model: defaultModel[selectedProvider || defaultProvider] || '' 
    }));
  };

  const renderMessage = (msg: ChatMessage, index: number) => {
    const isUser = msg.role === 'user';
    const isAssistant = msg.role === 'assistant';

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

            </div>
          }
        />
      </List.Item>
    );
  };

  return (
    <Layout style={{ height: 'calc(100vh - 64px)' }}>
      <Content style={{ display: 'flex', flexDirection: 'column', padding: 0 }}>
        {/* 聊天头部 */}
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
            <Text strong style={{ fontSize: '16px' }}>
              AI聊天对话
            </Text>
            {currentConversation && (
              <Text type="secondary" style={{ marginLeft: '16px' }}>
                {currentConversation.messages.length} 条消息
              </Text>
            )}
          </div>
          <Space>
            <Button onClick={handleNewConversation}>新建对话</Button>
            <Button
              icon={<SettingOutlined />}
              onClick={() => setShowSettings(!showSettings)}
            >
              设置
            </Button>
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
              <Text type="secondary">开始与AI助手对话吧！</Text>
            </div>
          ) : (
            <>
              <List
                dataSource={currentMessages}
                renderItem={renderMessage}
                style={{ padding: '16px 0' }}
              />
              
              {/* 显示流式消息 */}
              {streamingMessage && (
                <div style={{ padding: '16px 0' }}>
                  <List.Item
                    style={{
                      padding: '16px',
                      border: 'none',
                      backgroundColor: '#ffffff',
                    }}
                  >
                    <List.Item.Meta
                      avatar={
                        <Avatar
                          icon={<RobotOutlined />}
                          style={{
                            backgroundColor: '#52c41a',
                          }}
                        />
                      }
                      title={
                        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                          <Text strong>AI助手</Text>
                          <Space>
                            <Button
                              type="text"
                              size="small"
                              icon={<CopyOutlined />}
                              onClick={() => handleCopyMessage(streamingMessage)}
                              disabled={!streamingMessage}
                            />
                            <Text type="secondary" style={{ fontSize: '12px' }}>
                              正在输入...
                            </Text>
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
                            {streamingMessage}
                            <span style={{ animation: 'blink 1s infinite' }}>|</span>
                          </Paragraph>
                        </div>
                      }
                    />
                  </List.Item>
                </div>
              )}
            </>
          )}
          {isLoading && (
            <div style={{ textAlign: 'center', padding: '16px' }}>
              <Spin tip="AI正在思考中..." />
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
              placeholder="输入您的消息... (Shift+Enter 换行，Enter 发送)"
              autoSize={{ minRows: 1, maxRows: 4 }}
              style={{ flex: 1 }}
              disabled={isLoading}
            />
            <Button
              type="primary"
              icon={<SendOutlined />}
              onClick={handleSendMessage}
              loading={isLoading}
              disabled={!inputValue.trim()}
            >
              发送
            </Button>
          </div>
          {error && (
            <div 
              style={{ 
                marginTop: '8px', 
                padding: '8px 12px',
                backgroundColor: '#fff2f0',
                border: '1px solid #ffccc7',
                borderRadius: '6px',
                color: '#ff4d4f', 
                fontSize: '13px',
                display: 'flex',
                alignItems: 'center',
                gap: '8px'
              }}
            >
              <span>⚠️</span>
              <span>{error}</span>
              <Button 
                type="text" 
                size="small" 
                onClick={() => dispatch(clearError())}
                style={{ marginLeft: 'auto', color: '#ff4d4f' }}
              >
                ✕
              </Button>
            </div>
          )}
        </div>
      </Content>

      {/* 设置侧边栏 */}
      {showSettings && (
        <Sider
          width={300}
          style={{
            backgroundColor: '#ffffff',
            borderLeft: '1px solid #f0f0f0',
          }}
        >
          <div style={{ padding: '24px' }}>
            <Text strong style={{ fontSize: '16px' }}>
              聊天设置
            </Text>
            <Divider />

            <div style={{ marginBottom: '24px' }}>
              <Text strong>AI提供商</Text>
              <Select
                style={{ width: '100%', marginTop: '8px' }}
                value={selectedProvider || defaultProvider}
                placeholder="选择AI提供商"
                options={(providers || []).map(p => ({
                  label: p.name,
                  value: p.name,
                  disabled: !p.health,
                }))}
              />
            </div>

            <div style={{ marginBottom: '24px' }}>
              <Text strong>模型</Text>
              <Select
                style={{ width: '100%', marginTop: '8px' }}
                value={defaultModel}
                placeholder="选择模型"
                options={models[selectedProvider || defaultProvider]?.map(m => ({
                  label: m.name,
                  value: m.name,
                  disabled: !m.enabled,
                })) || []}
              />
            </div>

            <div style={{ marginBottom: '24px' }}>
              <Text strong>温度 ({localSettings.temperature})</Text>
              <Slider
                min={0}
                max={2}
                step={0.1}
                value={localSettings.temperature}
                onChange={(value) =>
                  setLocalSettings(prev => ({ ...prev, temperature: value }))
                }
                style={{ marginTop: '8px' }}
              />
              <Text type="secondary" style={{ fontSize: '12px' }}>
                控制回答的随机性，值越高越有创意
              </Text>
            </div>

            <div style={{ marginBottom: '24px' }}>
              <Text strong>最大令牌数 ({localSettings.maxTokens})</Text>
              <Slider
                min={100}
                max={4000}
                step={100}
                value={localSettings.maxTokens}
                onChange={(value) =>
                  setLocalSettings(prev => ({ ...prev, maxTokens: value }))
                }
                style={{ marginTop: '8px' }}
              />
              <Text type="secondary" style={{ fontSize: '12px' }}>
                限制回答的最大长度
              </Text>
            </div>

            <div style={{ marginBottom: '24px' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <Text strong>流式响应</Text>
                <Switch
                  checked={localSettings.stream}
                  onChange={(checked) =>
                    setLocalSettings(prev => ({ ...prev, stream: checked }))
                  }
                />
              </div>
              <Text type="secondary" style={{ fontSize: '12px', marginTop: '4px' }}>
                实时显示AI回答过程
              </Text>
            </div>
          </div>
        </Sider>
      )}
    </Layout>
  );
};

export default ChatPage;