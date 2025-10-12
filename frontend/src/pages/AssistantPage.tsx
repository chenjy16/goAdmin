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
  Drawer,
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
import { loadConfigData, selectConfig } from '../store/slices/configSlice';
import {
  fetchMCPTools,
  checkMCPStatus,
} from '../store/slices/mcpSlice';
import AssistantConfigPanel from '../components/AssistantConfigPanel';
import { useTranslation } from 'react-i18next';
import type { ChatMessage } from '../types/api';

const { TextArea } = Input;
const { Text, Paragraph } = Typography;
const { Content } = Layout;

const AssistantPage: React.FC = () => {
  const { t } = useTranslation();
  const dispatch = useAppDispatch();
  const {
    conversations,
    currentConversationId,
    isInitialized,
    isLoading,
    error,
  } = useAppSelector(state => state.assistant);

  const config = useAppSelector(selectConfig);
  const { tools: mcpTools, isInitialized: mcpInitialized } = useAppSelector(state => state.mcp);

  const [inputValue, setInputValue] = useState('');
  const [configDrawerVisible, setConfigDrawerVisible] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // 工具名称映射函数
  const getToolDisplayName = (name: string): string => {
    const nameMap: Record<string, string> = {
      '雅虎财经': 'yahoo_finance',
      '股票分析': 'stock_analysis', 
      '股票对比': 'stock_compare',
      '股票投资建议': 'stock_advice',
    };
    const key = nameMap[name] || name;
    return t(`mcpTools.toolNames.${key}`, name);
  };

  useEffect(() => {
    // 加载配置数据
    dispatch(loadConfigData());
    
    // 检查MCP状态并获取工具
    dispatch(checkMCPStatus()).then((result) => {
      if (checkMCPStatus.fulfilled.match(result) && result.payload.initialized) {
        dispatch(fetchMCPTools());
      }
    });
    
    if (!isInitialized) {
      dispatch(initializeAssistant());
    } else if (!currentConversationId) {
      dispatch(createAssistantConversation(t('assistant.newConversationName')));
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

    // 验证配置是否完整
    if (!config.selectedModel || !config.selectedProvider) {
      message.error(t('assistant.selectProviderModel'));
      return;
    }

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
        model: config.selectedModel,
        temperature: config.temperature,
        maxTokens: config.maxTokens,
        useTools: !!config.selectedTool,
        provider: config.selectedProvider,
        selectedTool: config.selectedTool,
      })).unwrap();
    } catch (err) {
      message.error(t('assistant.sendFailed'));
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
    message.success(t('assistant.copySuccess'));
  };

  const handleNewConversation = () => {
    dispatch(createAssistantConversation(t('assistant.newConversationName')));
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
              <Text strong>{isUser ? t('assistant.user') : t('assistant.ai')}</Text>
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
                  <Tag color="blue">{t('assistant.toolUsed')}</Tag>
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
          <h3>{t('assistant.notInitialized')}</h3>
          <p style={{ color: '#8c8c8c', marginBottom: '24px' }}>
            {t('assistant.initializePrompt')}
          </p>
          <Button
            type="primary"
            size="large"
            icon={<SettingOutlined />}
            onClick={handleInitializeAssistant}
            loading={isLoading}
          >
            {t('assistant.initializeButton')}
          </Button>
          {error && (
            <div style={{ marginTop: '16px', color: '#ff4d4f', fontSize: '14px' }}>
              {t('assistant.initializeFailed')}: {error}
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
                {t('assistant.title')}
              </Text>
              <Tag color="green" icon={<CheckCircleOutlined />}>
                {t('assistant.initialized')}
              </Tag>
            </Space>
            {currentConversation && (
              <Text type="secondary" style={{ marginLeft: '16px' }}>
                {currentConversation.messages.length} {t('assistant.messagesCount')}
              </Text>
            )}
          </div>
          <Space>
            <Button 
              icon={<SettingOutlined />}
              onClick={() => setConfigDrawerVisible(true)}
            >
              {t('assistant.configure')}
            </Button>
            <Button onClick={handleNewConversation}>{t('assistant.newConversation')}</Button>
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
              <Text type="secondary">{t('assistant.ready')}</Text>
              <div style={{ marginTop: '16px', textAlign: 'center' }}>
                <Text type="secondary" style={{ fontSize: '12px' }}>
                  {mcpInitialized && mcpTools.length > 0 ? t('assistant.toolsAvailable') : t('assistant.capabilities')}
                </Text>
                <div style={{ marginTop: '8px' }}>
                  {mcpInitialized && mcpTools.length > 0 ? (
                    mcpTools.slice(0, 4).map((tool) => (
                      <Tag key={tool.name} color="blue">
                        {getToolDisplayName(tool.name)}
                      </Tag>
                    ))
                  ) : (
                    <>
                      <Tag>{t('assistant.intelligentChat')}</Tag>
                      <Tag>{t('assistant.toolCalling')}</Tag>
                      <Tag>{t('assistant.codeGeneration')}</Tag>
                      <Tag>{t('assistant.questionAnswering')}</Tag>
                    </>
                  )}
                  {mcpInitialized && mcpTools.length > 4 && (
                    <Tag color="default">+{mcpTools.length - 4} {t('assistant.moreTools')}</Tag>
                  )}
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
              <Spin tip={t('assistant.thinking')} spinning={true}>
                <div style={{ height: '40px' }} />
              </Spin>
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
              placeholder={t('assistant.placeholder')}
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
              {t('assistant.send')}
            </Button>
          </div>
          {error && (
            <div style={{ marginTop: '8px', color: '#ff4d4f', fontSize: '12px' }}>
              {error}
            </div>
          )}
        </div>
      </Content>
      
      {/* 配置抽屉 */}
      <Drawer
        title={t('assistant.configTitle')}
        placement="right"
        width={600}
        open={configDrawerVisible}
        onClose={() => setConfigDrawerVisible(false)}
        destroyOnClose={false}
      >
        <AssistantConfigPanel />
      </Drawer>
    </Layout>
  );
};

export default AssistantPage;