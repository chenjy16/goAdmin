import React from 'react';
import { Card, Checkbox, Row, Col, Space, Typography } from 'antd';
import { ToolOutlined } from '@ant-design/icons';
import type { MCPTool } from '../../types/api';

const { Text } = Typography;

interface ToolSelectorProps {
  tools: MCPTool[];
  selectedTools: string[];
  onToolChange: (tools: string[]) => void;
  className?: string;
}

const ToolSelector: React.FC<ToolSelectorProps> = ({
  tools,
  selectedTools,
  onToolChange,
  className
}) => {
  return (
    <Card 
      className={className}
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
      
      {tools.length === 0 ? (
        <div style={{ 
          textAlign: 'center', 
          padding: '40px 20px',
          color: '#8c8c8c',
          backgroundColor: '#fafafa',
          borderRadius: '6px',
          border: '1px dashed #d9d9d9'
        }}>
          <ToolOutlined style={{ fontSize: '24px', marginBottom: '8px' }} />
          <div>暂无可用的MCP工具</div>
        </div>
      ) : (
        <Checkbox.Group
          value={selectedTools}
          onChange={onToolChange}
          style={{ width: '100%' }}
        >
          <Row gutter={[12, 12]} align="stretch">
            {tools.map(tool => (
              <Col span={12} key={tool.name} style={{ display: 'flex' }}>
                <div style={{ 
                  padding: '12px',
                  border: '1px solid #f0f0f0',
                  borderRadius: '6px',
                  backgroundColor: '#fafafa',
                  transition: 'all 0.2s ease',
                  cursor: 'pointer',
                  width: '100%',
                  minHeight: '80px',
                  display: 'flex',
                  alignItems: 'flex-start'
                }}>
                  <Checkbox value={tool.name} style={{ alignSelf: 'flex-start', marginTop: '2px' }}>
                    <div style={{ 
                      display: 'flex', 
                      flexDirection: 'column',
                      gap: '4px',
                      paddingLeft: '4px'
                    }}>
                      <Text strong style={{ 
                        fontSize: '14px',
                        lineHeight: '20px',
                        display: 'block'
                      }}>
                        {tool.name}
                      </Text>
                      <Text type="secondary" style={{ 
                        fontSize: '12px',
                        lineHeight: '16px',
                        display: 'block',
                        wordBreak: 'break-word'
                      }}>
                        {tool.description}
                      </Text>
                    </div>
                  </Checkbox>
                </div>
              </Col>
            ))}
          </Row>
        </Checkbox.Group>
      )}
    </Card>
  );
};

export default ToolSelector;