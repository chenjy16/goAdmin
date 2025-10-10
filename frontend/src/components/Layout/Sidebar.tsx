import React from 'react';
import { Layout, Menu, theme } from 'antd';
import type { MenuProps } from 'antd';
import { useNavigate } from 'react-router-dom';
import {
  RobotOutlined,
  SettingOutlined,
} from '@ant-design/icons';
import { useAppDispatch, useAppSelector } from '../../store';
import { setCurrentPage } from '../../store/slices/uiSlice';

const { Sider } = Layout;

const Sidebar: React.FC = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const { token } = theme.useToken();
  const { sidebarCollapsed, currentPage } = useAppSelector(state => state.ui);

  // 菜单项配置
  const menuItems: MenuProps['items'] = [
    {
      key: 'assistant',
      icon: <RobotOutlined />,
      label: 'AI助手',
    },
    {
      key: 'settings',
      icon: <SettingOutlined />,
      label: '设置',
    },
  ];

  // 处理菜单点击
  const handleMenuClick = ({ key }: { key: string }) => {
    navigate(`/${key}`);
    dispatch(setCurrentPage(key));
  };

  return (
    <Sider
      trigger={null}
      collapsible
      collapsed={sidebarCollapsed}
      style={{
        background: token.colorBgContainer,
        borderRight: `1px solid ${token.colorBorder}`,
      }}
    >
      {/* Logo */}
      <div className="layout-logo">
        <RobotOutlined style={{ fontSize: '24px', color: token.colorPrimary }} />
        {!sidebarCollapsed && (
          <span style={{ marginLeft: '8px', fontWeight: 'bold' }}>
            AI Chat Platform
          </span>
        )}
      </div>

      {/* 导航菜单 */}
      <Menu
        mode="inline"
        selectedKeys={[currentPage]}
        items={menuItems}
        onClick={handleMenuClick}
        style={{ borderRight: 0 }}
      />
    </Sider>
  );
};

export default Sidebar;