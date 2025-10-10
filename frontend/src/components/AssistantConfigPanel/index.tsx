import React, { useEffect } from 'react';
import { Card, Row, Col, Space, Alert, Spin } from 'antd';
import { SettingOutlined } from '@ant-design/icons';
import { useAppDispatch, useAppSelector } from '../../store';
import { 
  loadConfigData, 
  selectConfig, 
  selectConfigData, 
  selectConfigLoading, 
  selectConfigError,
  setProvider,
  setModel,
  setTools,
  setTemperature,
  setMaxTokens,
  setTopP
} from '../../store/slices/configSlice';
import ModelSelector from '../ModelSelector';
import ToolSelector from '../ToolSelector';
import ParameterSettings from '../ParameterSettings';
import type { ConfigState } from '../../services/configService';

interface AssistantConfigPanelProps {
  onConfigChange?: (config: ConfigState) => void;
  className?: string;
}

const AssistantConfigPanel: React.FC<AssistantConfigPanelProps> = ({
  onConfigChange,
  className
}) => {
  // 使用Redux store
  const dispatch = useAppDispatch();
  const config = useAppSelector(selectConfig);
  const configData = useAppSelector(selectConfigData);
  const loading = useAppSelector(selectConfigLoading);
  const error = useAppSelector(selectConfigError);

  // 初始化数据
  useEffect(() => {
    // 清除缓存以确保获取最新数据
    import('../../services/configService').then(({ configService }) => {
      configService.clearCache();
    });
    dispatch(loadConfigData());
  }, [dispatch]);

  // 当配置变化时通知父组件
  useEffect(() => {
    onConfigChange?.(config);
  }, [config, onConfigChange]);

  // 事件处理函数
  const handleProviderChange = (provider: string) => {
    dispatch(setProvider(provider));
  };

  const handleModelChange = (model: string) => {
    dispatch(setModel(model));
  };

  const handleToolChange = (tools: string[]) => {
    dispatch(setTools(tools));
  };

  const handleTemperatureChange = (value: number) => {
    dispatch(setTemperature(value));
  };

  const handleMaxTokensChange = (value: number) => {
    dispatch(setMaxTokens(value));
  };

  const handleTopPChange = (value: number) => {
    dispatch(setTopP(value));
  };

  // 渲染加载状态
  if (loading) {
    return (
      <Card className={className} style={{ height: '100%' }}>
        <div style={{ textAlign: 'center', padding: '40px' }}>
          <Spin size="large" />
          <div style={{ marginTop: '16px', color: '#666', fontSize: '14px' }}>加载配置中...</div>
        </div>
      </Card>
    );
  }

  // 渲染错误状态
  if (error) {
    return (
      <Card className={className} style={{ height: '100%' }}>
        <Alert
          message="配置加载失败"
          description={error}
          type="error"
          showIcon
          style={{ margin: '20px' }}
        />
      </Card>
    );
  }

  // 确保configData存在
  if (!configData) {
    return (
      <Card className={className} style={{ height: '100%' }}>
        <div style={{ textAlign: 'center', padding: '40px' }}>
          <div style={{ color: '#666', fontSize: '14px' }}>配置数据不可用</div>
        </div>
      </Card>
    );
  }

  return (
    <Card
      className={className}
      title={
        <Space size="middle">
          <SettingOutlined style={{ color: '#1890ff' }} />
          <span style={{ fontSize: '16px', fontWeight: 600 }}>AI助手配置</span>
        </Space>
      }
      style={{ 
        marginBottom: '20px',
        boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
        borderRadius: '8px'
      }}
      bodyStyle={{ padding: '24px' }}
    >
      <Row gutter={[20, 20]}>
        {/* 模型选择区域 */}
        <Col span={24}>
          <ModelSelector
            providers={configData.providers}
            models={configData.models}
            selectedProvider={config.selectedProvider}
            selectedModel={config.selectedModel}
            loading={loading}
            onProviderChange={handleProviderChange}
            onModelChange={handleModelChange}
          />
        </Col>

        {/* 工具选择区域 */}
        <Col span={24}>
          <ToolSelector
            tools={configData.tools}
            selectedTools={config.selectedTools}
            onToolChange={handleToolChange}
          />
        </Col>

        {/* 高级设置区域 */}
        <Col span={24}>
          <ParameterSettings
            temperature={config.temperature}
            maxTokens={config.maxTokens}
            topP={config.topP}
            onTemperatureChange={handleTemperatureChange}
            onMaxTokensChange={handleMaxTokensChange}
            onTopPChange={handleTopPChange}
          />
        </Col>
      </Row>
    </Card>
  );
};

export default AssistantConfigPanel;