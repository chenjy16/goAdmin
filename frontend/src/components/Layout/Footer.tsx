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
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', gap: '8px' }}>
        <a 
          href="https://github.com/chenjy16/go-springAi" 
          target="_blank" 
          rel="noopener noreferrer"
          style={{ 
            color: token.colorPrimary,
            textDecoration: 'none'
          }}
        >
          https://github.com/chenjy16/go-springAi
        </a>
        <span style={{ color: token.colorTextSecondary }}>©2025</span>
      </div>
    </AntFooter>
  );
};

export default Footer;