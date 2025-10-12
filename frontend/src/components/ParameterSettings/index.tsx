import React from 'react';
import { Card, Slider, InputNumber, Row, Col, Space, Typography, Divider } from 'antd';
import { SettingOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';

const { Text, Title } = Typography;

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
      padding: '20px',
      backgroundColor: '#ffffff',
      borderRadius: '8px',
      border: '1px solid #e8e8e8',
      boxShadow: '0 1px 4px rgba(0, 0, 0, 0.04)',
      transition: 'all 0.2s ease',
      height: '100%',
      display: 'flex',
      flexDirection: 'column'
    }}>
      <div style={{ 
        marginBottom: '16px',
        textAlign: 'center',
        flex: '0 0 auto'
      }}>
        <Title level={5} style={{ 
          margin: 0, 
          marginBottom: '4px',
          fontSize: '16px', 
          color: '#262626',
          fontWeight: 600
        }}>
          {label}
        </Title>
        <Text type="secondary" style={{ 
          fontSize: '12px',
          lineHeight: '1.4'
        }}>
          {description}
        </Text>
      </div>
      
      <div style={{ flex: '1 1 auto', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
        <div style={{ marginBottom: '16px' }}>
          <Slider
            min={min}
            max={max}
            step={step}
            value={value}
            onChange={onChange}
            trackStyle={{ backgroundColor: color, height: '6px' }}
            handleStyle={{ 
              borderColor: color, 
              backgroundColor: color,
              width: '16px',
              height: '16px',
              marginTop: '-5px'
            }}
            railStyle={{ height: '6px', backgroundColor: '#f0f0f0' }}
          />
        </div>
        
        <div style={{ textAlign: 'center' }}>
          <InputNumber
            min={min}
            max={max}
            step={step}
            value={value}
            onChange={(val) => onChange(val || min)}
            style={{ 
              width: '100px',
              textAlign: 'center'
            }}
            size="middle"
            formatter={formatter}
            parser={parser}
            controls={false}
          />
        </div>
      </div>
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
  const { t } = useTranslation();
  return (
    <Card 
      className={className}
      title={
        <Space size="middle" style={{ alignItems: 'center' }}>
          <SettingOutlined style={{ 
            color: '#722ed1', 
            fontSize: '18px'
          }} />
          <span style={{ 
            fontSize: '16px', 
            fontWeight: 600,
            color: '#262626'
          }}>
{t('parameterSettings.title')}
          </span>
        </Space>
      }
      style={{ 
        borderRadius: '12px',
        border: '1px solid #e8e8e8',
        boxShadow: '0 2px 8px rgba(0, 0, 0, 0.06)',
        background: 'linear-gradient(135deg, #fafbfc 0%, #f8f9fa 100%)'
      }}
      bodyStyle={{ 
        padding: '24px',
        background: 'transparent'
      }}
      headStyle={{
        borderBottom: '1px solid #e8e8e8',
        background: 'rgba(255, 255, 255, 0.8)',
        borderRadius: '12px 12px 0 0'
      }}
    >
      <Row gutter={[24, 24]} style={{ minHeight: '200px' }}>
        <Col xs={24} sm={24} md={8} lg={8} xl={8}>
          <ParameterControl
            label={t('settings.temperature')}
            description={t('parameterSettings.temperatureDesc')}
            value={temperature}
            min={0}
            max={2}
            step={0.1}
            onChange={onTemperatureChange}
            color="#1890ff"
          />
        </Col>

        <Col xs={24} sm={24} md={8} lg={8} xl={8}>
          <ParameterControl
            label={t('settings.maxTokens')}
            description={t('parameterSettings.maxTokensDesc')}
            value={maxTokens}
            min={1}
            max={32768}
            step={1}
            onChange={onMaxTokensChange}
            color="#52c41a"
            parser={(value) => Number(value?.replace(/\$\s?|(,*)/g, ''))}
          />
        </Col>

        <Col xs={24} sm={24} md={8} lg={8} xl={8}>
          <ParameterControl
            label={t('settings.topP')}
            description={t('parameterSettings.topPDesc')}
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