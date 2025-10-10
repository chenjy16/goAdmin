import React, { useMemo } from 'react';
import { Card, Select, Space, Typography, Tag, Spin, Row, Col } from 'antd';
import { RobotOutlined } from '@ant-design/icons';
import type { ProviderInfo, ModelInfo } from '../../types/api';

const { Text } = Typography;
const { Option } = Select;

interface ModelSelectorProps {
  providers: ProviderInfo[];
  models: Record<string, ModelInfo[]>;
  selectedProvider?: string;
  selectedModel?: string;
  loading?: boolean;
  onProviderChange: (provider: string) => void;
  onModelChange: (model: string) => void;
  className?: string;
}

const ModelSelector: React.FC<ModelSelectorProps> = ({
  providers,
  models,
  selectedProvider,
  selectedModel,
  loading = false,
  onProviderChange,
  onModelChange,
  className
}) => {
  // 获取健康的提供商
  const healthyProviders = useMemo(() => {
    return (providers || []).filter((provider: ProviderInfo) => provider.healthy);
  }, [providers]);

  // 获取当前选择的提供商的模型列表（只显示启用的模型）
  const availableModels = useMemo(() => {
    if (!selectedProvider || !models[selectedProvider]) {
      return [];
    }
    
    const allModels = models[selectedProvider];
    const enabledModels = allModels.filter((model: ModelInfo) => model.enabled);
    
    return enabledModels;
  }, [selectedProvider, models]);

  // 获取当前选择的模型信息
  const getCurrentModel = (): ModelInfo | undefined => {
    if (!selectedProvider || !selectedModel) return undefined;
    const providerModels = models[selectedProvider] || [];
    return providerModels.find(m => m.name === selectedModel);
  };

  const currentModel = getCurrentModel();

  return (
    <Card
      title={
        <Space size="middle">
          <RobotOutlined style={{ color: '#1890ff' }} />
          <span style={{ fontSize: '16px', fontWeight: 600 }}>模型选择</span>
        </Space>
      }
      style={{ 
        marginBottom: '20px',
        boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
        borderRadius: '8px'
      }}
      bodyStyle={{ padding: '24px' }}
    >
      <Row gutter={[16, 16]}>
        {/* AI提供商选择 */}
        <Col span={24}>
          <div style={{ marginBottom: '8px' }}>
            <Text strong>AI提供商</Text>
            <Text type="secondary" style={{ marginLeft: '8px' }}>
              ({healthyProviders.length} 个可用)
            </Text>
          </div>
          <Select
            value={selectedProvider}
            onChange={onProviderChange}
            placeholder="选择AI提供商"
            style={{ width: '100%' }}
            loading={loading}
            disabled={loading || healthyProviders.length === 0}
            size="large"
          >
            {healthyProviders.map((provider: ProviderInfo) => (
              <Option key={provider.type} value={provider.type}>
                <Space>
                  <span>{provider.name}</span>
                  <Tag color="green">健康</Tag>
                </Space>
              </Option>
            ))}
          </Select>
        </Col>

        {/* 模型选择 */}
        <Col span={24}>
          <div style={{ marginBottom: '8px' }}>
            <Text strong>模型</Text>
            <Text type="secondary" style={{ marginLeft: '8px' }}>
              ({availableModels.length} 个可用)
            </Text>
          </div>
          <Select
            value={selectedModel}
            onChange={onModelChange}
            placeholder="选择模型"
            style={{ width: '100%' }}
            loading={loading}
            disabled={loading || !selectedProvider || availableModels.length === 0}
            size="large"
          >
            {availableModels.map((model: ModelInfo) => (
              <Option key={model.name} value={model.name}>
                <Space>
                  <span>{model.display_name || model.name}</span>
                  {model.enabled ? (
                    <Tag color="green">启用</Tag>
                  ) : (
                    <Tag color="red">禁用</Tag>
                  )}
                </Space>
              </Option>
            ))}
          </Select>
        </Col>
      </Row>
    </Card>
  );
};

export default ModelSelector;