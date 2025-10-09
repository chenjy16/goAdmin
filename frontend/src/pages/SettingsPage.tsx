import React, { useEffect, useState } from 'react';
import {
  Card,
  Form,
  Input,
  Switch,
  Select,
  Button,
  Space,
  Divider,
  Typography,
  Row,
  Col,
  InputNumber,
  message,
  Alert,
  Tabs,
  List,
  Tag,
} from 'antd';
import {
  SettingOutlined,
  SaveOutlined,
  ReloadOutlined,
  ExportOutlined,
  ImportOutlined,
  DeleteOutlined,
  InfoCircleOutlined,
} from '@ant-design/icons';
import { useAppDispatch, useAppSelector } from '../store';
import {
  updateChatSettings,
  resetSettings,
  loadSettings,
} from '../store/slices/settingsSlice';

const { Title, Text, Paragraph } = Typography;
const { TextArea } = Input;
const { TabPane } = Tabs;

const SettingsPage: React.FC = () => {
  const dispatch = useAppDispatch();
  const settings = useAppSelector(state => state.settings);
  const [form] = Form.useForm();
  const [hasChanges, setHasChanges] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    form.setFieldsValue(settings);
  }, [settings, form]);

  const handleSave = async () => {
    try {
      setIsLoading(true);
      const values = await form.validateFields();
      if (values.chatSettings) {
        dispatch(updateChatSettings(values.chatSettings));
      }
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

  const handleExport = async () => {
    try {
      const dataStr = JSON.stringify(settings, null, 2);
      const dataBlob = new Blob([dataStr], { type: 'application/json' });
      const url = URL.createObjectURL(dataBlob);
      const link = document.createElement('a');
      link.href = url;
      link.download = 'settings.json';
      link.click();
      URL.revokeObjectURL(url);
      message.success('设置已导出');
    } catch (err) {
      message.error('导出设置失败');
    }
  };

  const handleImport = () => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = '.json';
    input.onchange = async (e) => {
      const file = (e.target as HTMLInputElement).files?.[0];
      if (file) {
        try {
          const text = await file.text();
          const importedSettings = JSON.parse(text);
          dispatch(loadSettings(importedSettings));
          form.setFieldsValue(importedSettings);
          message.success('设置已导入');
        } catch (err) {
          message.error('导入设置失败');
        }
      }
    };
    input.click();
  };

  const handleFormChange = () => {
    setHasChanges(true);
  };

  const generalSettings = (
    <Card title="通用设置" style={{ marginBottom: 16 }}>
      <Form.Item
        name="theme"
        label="主题"
        tooltip="选择应用程序的主题"
      >
        <Select>
          <Select.Option value="light">浅色主题</Select.Option>
          <Select.Option value="dark">深色主题</Select.Option>
          <Select.Option value="auto">跟随系统</Select.Option>
        </Select>
      </Form.Item>

      <Form.Item
        name="language"
        label="语言"
        tooltip="选择界面语言"
      >
        <Select>
          <Select.Option value="zh-CN">简体中文</Select.Option>
          <Select.Option value="en-US">English</Select.Option>
        </Select>
      </Form.Item>

      <Form.Item
        name="autoSave"
        label="自动保存"
        valuePropName="checked"
        tooltip="自动保存对话和设置"
      >
        <Switch />
      </Form.Item>

      <Form.Item
        name="showNotifications"
        label="显示通知"
        valuePropName="checked"
        tooltip="显示系统通知"
      >
        <Switch />
      </Form.Item>
    </Card>
  );

  const chatSettings = (
    <Card title="聊天设置" style={{ marginBottom: 16 }}>
      <Form.Item
        name={['chat', 'defaultProvider']}
        label="默认提供商"
        tooltip="新对话的默认AI提供商"
      >
        <Select placeholder="选择默认提供商">
          <Select.Option value="openai">OpenAI</Select.Option>
          <Select.Option value="anthropic">Anthropic</Select.Option>
          <Select.Option value="google">Google</Select.Option>
          <Select.Option value="azure">Azure OpenAI</Select.Option>
        </Select>
      </Form.Item>

      <Form.Item
        name={['chat', 'defaultModel']}
        label="默认模型"
        tooltip="新对话的默认AI模型"
      >
        <Input placeholder="例如: gpt-4" />
      </Form.Item>

      <Form.Item
        name={['chat', 'maxTokens']}
        label="最大令牌数"
        tooltip="单次响应的最大令牌数"
      >
        <InputNumber min={1} max={32000} style={{ width: '100%' }} />
      </Form.Item>

      <Form.Item
        name={['chat', 'temperature']}
        label="温度"
        tooltip="控制响应的随机性 (0-2)"
      >
        <InputNumber min={0} max={2} step={0.1} style={{ width: '100%' }} />
      </Form.Item>

      <Form.Item
        name={['chat', 'streamResponse']}
        label="流式响应"
        valuePropName="checked"
        tooltip="启用实时流式响应"
      >
        <Switch />
      </Form.Item>

      <Form.Item
        name={['chat', 'saveHistory']}
        label="保存历史"
        valuePropName="checked"
        tooltip="保存对话历史记录"
      >
        <Switch />
      </Form.Item>
    </Card>
  );

  const mcpSettings = (
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
  );

  const advancedSettings = (
    <Card title="高级设置" style={{ marginBottom: 16 }}>
      <Alert
        message="警告"
        description="修改高级设置可能影响应用程序的稳定性，请谨慎操作。"
        type="warning"
        showIcon
        style={{ marginBottom: 16 }}
      />

      <Form.Item
        name={['advanced', 'apiTimeout']}
        label="API超时时间 (毫秒)"
        tooltip="API请求的超时时间"
      >
        <InputNumber min={1000} max={60000} style={{ width: '100%' }} />
      </Form.Item>

      <Form.Item
        name={['advanced', 'maxConcurrentRequests']}
        label="最大并发请求数"
        tooltip="同时进行的最大API请求数"
      >
        <InputNumber min={1} max={20} style={{ width: '100%' }} />
      </Form.Item>

      <Form.Item
        name={['advanced', 'enableDebugMode']}
        label="调试模式"
        valuePropName="checked"
        tooltip="启用调试模式以获取详细日志"
      >
        <Switch />
      </Form.Item>

      <Form.Item
        name={['advanced', 'customApiEndpoint']}
        label="自定义API端点"
        tooltip="使用自定义的API端点"
      >
        <Input placeholder="https://api.example.com" />
      </Form.Item>

      <Form.Item
        name={['advanced', 'proxySettings']}
        label="代理设置"
        tooltip="网络代理配置"
      >
        <TextArea
          rows={3}
          placeholder="代理配置 (JSON格式)"
        />
      </Form.Item>
    </Card>
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
            icon={<ImportOutlined />}
            onClick={handleImport}
          >
            导入设置
          </Button>
          <Button
            icon={<ExportOutlined />}
            onClick={handleExport}
          >
            导出设置
          </Button>
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
        <Tabs defaultActiveKey="general" type="card">
          <TabPane tab="通用" key="general">
            {generalSettings}
          </TabPane>
          <TabPane tab="聊天" key="chat">
            {chatSettings}
          </TabPane>
          <TabPane tab="MCP工具" key="mcp">
            {mcpSettings}
          </TabPane>
          <TabPane tab="高级" key="advanced">
            {advancedSettings}
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