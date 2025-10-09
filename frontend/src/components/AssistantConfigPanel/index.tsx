import React, { useState, useEffect, useMemo } from 'react';
import { Card, Select, Radio, Checkbox, Slider, InputNumber, Row, Col, Space, Typography, Tag, Alert, Spin } from 'antd';
import { SettingOutlined, ApiOutlined, ToolOutlined } from '@ant-design/icons';
import type { ProviderInfo, ModelInfo, MCPTool } from '../../types/api';
import { configService, type ConfigData, type ConfigState } from '../../services/configService';
import { useAppDispatch, useAppSelector } from '../../store';
import { 
  loadConfigData, 
  selectConfig, 
  selectConfigData, 
  selectConfigLoading, 
  selectConfigError,
  updateConfig,
  setProvider,
  setModel,
  setTools,
  setTemperature,
  setMaxTokens,
  setTopP
} from '../../store/slices/configSlice';

const { Title, Text } = Typography;
const { Option } = Select;

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
    console.log('开始加载配置数据...');
    dispatch(loadConfigData());
  }, [dispatch]);

  // 当配置变化时通知父组件
  useEffect(() => {
    onConfigChange?.(config);
  }, [config, onConfigChange]);

  // 提供商变化处理
  const handleProviderChange = (provider: string) => {
    console.log('提供商改变:', provider);
    dispatch(setProvider(provider));
  };

  // 模型变化处理
  const handleModelChange = (model: string) => {
    console.log('模型改变:', model);
    dispatch(setModel(model));
  };

  // 工具变化处理
  const handleToolChange = (tools: string[]) => {
    console.log('工具改变:', tools);
    dispatch(setTools(tools));
  };

  // 参数变化处理
  const handleTemperatureChange = (value: number) => {
    dispatch(setTemperature(value));
  };

  const handleMaxTokensChange = (value: number) => {
    dispatch(setMaxTokens(value));
  };

  const handleTopPChange = (value: number) => {
    dispatch(setTopP(value));
  };

  // 获取当前选择的提供商的模型列表
  const currentProviderModels = useMemo(() => {
    if (!config.selectedProvider || !configData?.models) {
      console.log('没有选择提供商或没有模型数据:', { selectedProvider: config.selectedProvider, hasModels: !!configData?.models });
      return [];
    }
    const models = configData.models[config.selectedProvider] || [];
    console.log(`提供商 ${config.selectedProvider} 的模型:`, models);
    return models;
  }, [config.selectedProvider, configData?.models]);

  // 获取当前选择的模型信息
  const getCurrentModel = (): ModelInfo | undefined => {
    if (!config.selectedProvider || !config.selectedModel) return undefined;
    const providerModels = configData.models[config.selectedProvider] || [];
    return providerModels.find(m => m.name === config.selectedModel);
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

  const currentModel = getCurrentModel();

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
          <Card 
            type="inner" 
            title={
              <Space size="middle">
                <ApiOutlined style={{ color: '#52c41a' }} />
                <span style={{ fontSize: '15px', fontWeight: 500 }}>模型选择</span>
              </Space>
            }
            style={{ 
              borderRadius: '6px',
              border: '1px solid #e8f4fd'
            }}
            bodyStyle={{ padding: '20px' }}
          >
            <Row gutter={[20, 20]}>
              <Col span={24}>
                <div style={{ marginBottom: '12px' }}>
                  <Text strong style={{ fontSize: '14px', color: '#262626' }}>AI提供商</Text>
                  <Text type="secondary" style={{ marginLeft: '12px', fontSize: '13px' }}>
                    ({configData.providers.length} 个可用)
                  </Text>
                </div>
                <Select
                  style={{ width: '100%' }}
                  size="large"
                  placeholder={configData.providers.length > 0 ? "选择AI提供商" : "暂无可用提供商"}
                  value={config.selectedProvider}
                  onChange={handleProviderChange}
                  loading={loading}
                  notFoundContent={
                    loading ? (
                      <div style={{ textAlign: 'center', padding: '8px' }}>
                        <Spin size="small" /> 加载提供商中...
                      </div>
                    ) : (
                      <div style={{ textAlign: 'center', padding: '8px', color: '#8c8c8c' }}>
                        暂无可用的AI提供商
                      </div>
                    )
                  }
                >
                  {configData.providers.map(provider => (
                    <Option key={provider.name} value={provider.name} disabled={!provider.healthy}>
                      <div style={{ padding: '4px 0', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <div style={{ flex: 1 }}>
                          <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '2px' }}>
                            <span style={{ fontWeight: 500, fontSize: '14px', color: provider.healthy ? '#262626' : '#8c8c8c' }}>
                              {provider.name}
                            </span>
                            <Tag 
                              color={provider.healthy ? 'green' : 'red'} 
                              style={{ fontSize: '11px', lineHeight: '16px', margin: 0, padding: '0 4px' }}
                            >
                              {provider.healthy ? '健康' : '不可用'}
                            </Tag>
                          </div>
                          <div style={{ fontSize: '12px', color: '#8c8c8c', lineHeight: '16px' }}>
                            {provider.model_count || 0} 个模型
                            {provider.description && provider.description !== `${provider.name} AI Provider` && 
                              ` • ${provider.description}`
                            }
                          </div>
                        </div>
                      </div>
                    </Option>
                  ))}
                </Select>
              </Col>
              
              <Col span={24}>
                <div style={{ marginBottom: '12px' }}>
                  <Text strong style={{ fontSize: '14px', color: '#262626' }}>模型</Text>
                  {config.selectedProvider && (
                    <Text type="secondary" style={{ marginLeft: '12px', fontSize: '13px' }}>
                      ({(configData.models[config.selectedProvider] || []).length} 个可用)
                    </Text>
                  )}
                </div>
                <Select
                  style={{ width: '100%' }}
                  size="large"
                  placeholder={config.selectedProvider ? "选择模型" : "请先选择提供商"}
                  value={config.selectedModel}
                  onChange={handleModelChange}
                  disabled={!config.selectedProvider}
                  loading={loading}
                  notFoundContent={
                    loading ? (
                      <div style={{ textAlign: 'center', padding: '8px' }}>
                        <Spin size="small" /> 加载模型中...
                      </div>
                    ) : config.selectedProvider ? (
                      <div style={{ textAlign: 'center', padding: '8px', color: '#8c8c8c' }}>
                        该提供商暂无可用模型
                      </div>
                    ) : (
                      <div style={{ textAlign: 'center', padding: '8px', color: '#8c8c8c' }}>
                        请先选择提供商
                      </div>
                    )
                  }
                >
                  {currentProviderModels.map(model => (
                    <Option key={model.name} value={model.name}>
                      <div style={{ padding: '6px 0' }}>
                        <Space direction="vertical" size={4}>
                          <span style={{ fontWeight: 500, fontSize: '14px' }}>{model.name}</span>
                          <Text type="secondary" style={{ fontSize: '12px' }}>
                            最大 {model.max_tokens} tokens
                            {model.description && ` • ${model.description}`}
                          </Text>
                        </Space>
                      </div>
                    </Option>
                  ))}
                </Select>
              </Col>
            </Row>
            
            {/* 模型信息显示 */}
            {currentModel && (
              <div style={{ 
                marginTop: '16px', 
                padding: '12px 16px', 
                backgroundColor: '#f6ffed', 
                borderRadius: '6px',
                border: '1px solid #b7eb8f'
              }}>
                <Text style={{ fontSize: '13px', color: '#389e0d' }}>
                  <strong>模型信息:</strong> 最大tokens: {currentModel.max_tokens}, 
                  默认温度: {currentModel.temperature}
                  {currentModel.description && `, ${currentModel.description}`}
                </Text>
              </div>
            )}
          </Card>
        </Col>

        {/* 工具选择区域 */}
        <Col span={24}>
          <Card 
            type="inner" 
            title={
              <Space size="middle">
                <ToolOutlined style={{ color: '#fa8c16' }} />
                <span style={{ fontSize: '15px', fontWeight: 500 }}>工具选择</span>
              </Space>
            }
            style={{ 
              borderRadius: '6px',
              border: '1px solid #fff7e6'
            }}
            bodyStyle={{ padding: '20px' }}
          >
            <div style={{ marginBottom: '16px' }}>
              <Text strong style={{ fontSize: '14px', color: '#262626' }}>可用工具</Text>
              <Text type="secondary" style={{ marginLeft: '12px', fontSize: '13px' }}>
                (选择要使用的MCP工具)
              </Text>
            </div>
            <Checkbox.Group
              value={config.selectedTools}
              onChange={handleToolChange}
              style={{ width: '100%' }}
            >
              <Row gutter={[12, 12]}>
                {configData.tools.map(tool => (
                  <Col span={12} key={tool.name}>
                    <div style={{ 
                      padding: '12px',
                      border: '1px solid #f0f0f0',
                      borderRadius: '6px',
                      backgroundColor: '#fafafa'
                    }}>
                      <Checkbox value={tool.name}>
                        <Space direction="vertical" size={4}>
                          <Text strong style={{ fontSize: '14px' }}>{tool.name}</Text>
                          <Text type="secondary" style={{ fontSize: '12px' }}>
                            {tool.description}
                          </Text>
                        </Space>
                      </Checkbox>
                    </div>
                  </Col>
                ))}
              </Row>
            </Checkbox.Group>
          </Card>
        </Col>

        {/* 高级设置区域 */}
        <Col span={24}>
          <Card 
            type="inner" 
            title={
              <Space size="middle">
                <SettingOutlined style={{ color: '#722ed1' }} />
                <span style={{ fontSize: '15px', fontWeight: 500 }}>高级设置</span>
              </Space>
            }
            style={{ 
              borderRadius: '6px',
              border: '1px solid #f9f0ff'
            }}
            bodyStyle={{ padding: '20px' }}
          >
            <Row gutter={[24, 24]}>
              <Col span={8}>
                <div style={{ 
                  padding: '16px',
                  backgroundColor: '#fafafa',
                  borderRadius: '6px',
                  border: '1px solid #f0f0f0'
                }}>
                  <div style={{ marginBottom: '12px' }}>
                    <Text strong style={{ fontSize: '14px', color: '#262626' }}>Temperature</Text>
                    <Text type="secondary" style={{ marginLeft: '8px', fontSize: '12px' }}>
                      (创造性: 0-2)
                    </Text>
                  </div>
                  <Row gutter={12}>
                    <Col span={16}>
                      <Slider
                        min={0}
                        max={2}
                        step={0.1}
                        value={config.temperature}
                        onChange={(value) => handleTemperatureChange(value)}
                        trackStyle={{ backgroundColor: '#1890ff' }}
                        handleStyle={{ borderColor: '#1890ff' }}
                      />
                    </Col>
                    <Col span={8}>
                      <InputNumber
                        min={0}
                        max={2}
                        step={0.1}
                        value={config.temperature}
                        onChange={(value) => handleTemperatureChange(value || 0.7)}
                        style={{ width: '100%' }}
                        size="small"
                      />
                    </Col>
                  </Row>
                </div>
              </Col>

              <Col span={8}>
                <div style={{ 
                  padding: '16px',
                  backgroundColor: '#fafafa',
                  borderRadius: '6px',
                  border: '1px solid #f0f0f0'
                }}>
                  <div style={{ marginBottom: '12px' }}>
                    <Text strong style={{ fontSize: '14px', color: '#262626' }}>Max Tokens</Text>
                    <Text type="secondary" style={{ marginLeft: '8px', fontSize: '12px' }}>
                      (最大输出长度)
                    </Text>
                  </div>
                  <Row gutter={12}>
                    <Col span={16}>
                      <Slider
                        min={1}
                        max={32768}
                        value={config.maxTokens}
                        onChange={(value) => handleMaxTokensChange(value)}
                        trackStyle={{ backgroundColor: '#52c41a' }}
                        handleStyle={{ borderColor: '#52c41a' }}
                      />
                    </Col>
                    <Col span={8}>
                      <InputNumber
                        min={1}
                        max={32768}
                        value={config.maxTokens}
                        onChange={(value) => handleMaxTokensChange(value || 2048)}
                        style={{ width: '100%' }}
                        parser={(value) => Number(value?.replace(/\$\s?|(,*)/g, ''))}
                        size="small"
                      />
                    </Col>
                  </Row>
                </div>
              </Col>

              <Col span={8}>
                <div style={{ 
                  padding: '16px',
                  backgroundColor: '#fafafa',
                  borderRadius: '6px',
                  border: '1px solid #f0f0f0'
                }}>
                  <div style={{ marginBottom: '12px' }}>
                    <Text strong style={{ fontSize: '14px', color: '#262626' }}>Top-p</Text>
                    <Text type="secondary" style={{ marginLeft: '8px', fontSize: '12px' }}>
                      (多样性: 0-1)
                    </Text>
                  </div>
                  <Row gutter={12}>
                    <Col span={16}>
                      <Slider
                        min={0}
                        max={1}
                        step={0.1}
                        value={config.topP}
                        onChange={(value) => handleTopPChange(value)}
                        trackStyle={{ backgroundColor: '#fa8c16' }}
                        handleStyle={{ borderColor: '#fa8c16' }}
                      />
                    </Col>
                    <Col span={8}>
                      <InputNumber
                        min={0}
                        max={1}
                        step={0.1}
                        value={config.topP}
                        onChange={(value) => handleTopPChange(value || 1.0)}
                        style={{ width: '100%' }}
                        size="small"
                      />
                    </Col>
                  </Row>
                </div>
              </Col>
            </Row>
          </Card>
        </Col>
      </Row>
    </Card>
  );
};

export default AssistantConfigPanel;