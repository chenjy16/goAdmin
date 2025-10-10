import React from 'react';
import { Card, Slider, InputNumber, Row, Col, Space, Typography } from 'antd';
import { SettingOutlined } from '@ant-design/icons';

const { Text } = Typography;

interface ParameterSettingsProps {
  temperature: number;
  maxTokens: number;
  topP: number;
  onTemperatureChange: (value: number) => void;
  onMaxTokensChange: (value: number) => void;
  onTopPChange: (value: number) => void;
  className?: string;
}

interface ParameterControlProps {
  label: string;
  description: string;
  value: number;
  min: number;
  max: number;
  step?: number;
  onChange: (value: number) => void;
  color: string;
  formatter?: (value: number | undefined) => string;
  parser?: (value: string | undefined) => number;
}

const ParameterControl: React.FC<ParameterControlProps> = ({
  label,
  description,
  value,
  min,
  max,
  step = 0.1,
  onChange,
  color,
  formatter,
  parser
}) => {
  return (
    <div style={{ 
      padding: '16px',
      backgroundColor: '#fafafa',
      borderRadius: '6px',
      border: '1px solid #f0f0f0'
    }}>
      <div style={{ marginBottom: '12px' }}>
        <Text strong style={{ fontSize: '14px', color: '#262626' }}>{label}</Text>
        <Text type="secondary" style={{ marginLeft: '8px', fontSize: '12px' }}>
          ({description})
        </Text>
      </div>
      <Row gutter={12}>
        <Col span={16}>
          <Slider
            min={min}
            max={max}
            step={step}
            value={value}
            onChange={onChange}
            trackStyle={{ backgroundColor: color }}
            handleStyle={{ borderColor: color }}
          />
        </Col>
        <Col span={8}>
          <InputNumber
            min={min}
            max={max}
            step={step}
            value={value}
            onChange={(val) => onChange(val || min)}
            style={{ width: '100%' }}
            size="small"
            formatter={formatter}
            parser={parser}
          />
        </Col>
      </Row>
    </div>
  );
};

const ParameterSettings: React.FC<ParameterSettingsProps> = ({
  temperature,
  maxTokens,
  topP,
  onTemperatureChange,
  onMaxTokensChange,
  onTopPChange,
  className
}) => {
  return (
    <Card 
      className={className}
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
          <ParameterControl
            label="Temperature"
            description="创造性: 0-2"
            value={temperature}
            min={0}
            max={2}
            step={0.1}
            onChange={onTemperatureChange}
            color="#1890ff"
          />
        </Col>

        <Col span={8}>
          <ParameterControl
            label="Max Tokens"
            description="最大输出长度"
            value={maxTokens}
            min={1}
            max={32768}
            step={1}
            onChange={onMaxTokensChange}
            color="#52c41a"
            parser={(value) => Number(value?.replace(/\$\s?|(,*)/g, ''))}
          />
        </Col>

        <Col span={8}>
          <ParameterControl
            label="Top-p"
            description="多样性: 0-1"
            value={topP}
            min={0}
            max={1}
            step={0.1}
            onChange={onTopPChange}
            color="#fa8c16"
          />
        </Col>
      </Row>
    </Card>
  );
};

export default ParameterSettings;