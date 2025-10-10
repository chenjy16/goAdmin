import React from 'react';
import { Layout, theme } from 'antd';

const { Footer: AntFooter } = Layout;

const Footer: React.FC = () => {
  const { token } = theme.useToken();

  return (
    <AntFooter
      style={{
        textAlign: 'center',
        padding: '12px 16px',
        background: token.colorBgContainer,
        borderTop: `1px solid ${token.colorBorder}`,
      }}
    >
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <span>AI Chat Management Platform ©2024</span>
        <div style={{ display: 'flex', gap: '16px', fontSize: '12px', color: token.colorTextSecondary }}>
          <span>状态: 正常</span>
          <span>版本: v1.0.0</span>
        </div>
      </div>
    </AntFooter>
  );
};

export default Footer;