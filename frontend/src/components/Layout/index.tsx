import React, { useEffect } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import { Layout as AntLayout, Menu, Button, Avatar, Dropdown, theme } from 'antd';
import type { MenuProps } from 'antd';
import {
  CloudOutlined,
  ToolOutlined,
  RobotOutlined,
  SettingOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  UserOutlined,
  LogoutOutlined,
} from '@ant-design/icons';
import { useAppDispatch, useAppSelector } from '../../store';
import { toggleSidebar, setCurrentPage } from '../../store/slices/uiSlice';
import './index.css';

const { Header, Sider, Content, Footer } = AntLayout;

const Layout: React.FC = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const location = useLocation();
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

  // 用户下拉菜单
  const userMenuItems: MenuProps['items'] = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人资料',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
    },
  ];

  // 处理菜单点击
  const handleMenuClick = ({ key }: { key: string }) => {
    navigate(`/${key}`);
    dispatch(setCurrentPage(key));
  };

  // 处理用户菜单点击
  const handleUserMenuClick = ({ key }: { key: string }) => {
    switch (key) {
      case 'profile':
        // 处理个人资料
        break;
      case 'logout':
        // 处理退出登录
        break;
    }
  };

  // 根据路径更新当前页面
  useEffect(() => {
    const path = location.pathname.slice(1) || 'assistant';
    dispatch(setCurrentPage(path));
  }, [location.pathname, dispatch]);

  return (
    <AntLayout style={{ minHeight: '100vh' }}>
      {/* 侧边栏 */}
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

      <AntLayout>
        {/* 顶部导航栏 */}
        <Header
          style={{
            padding: '0 16px',
            background: token.colorBgContainer,
            borderBottom: `1px solid ${token.colorBorder}`,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
          }}
        >
          <Button
            type="text"
            icon={sidebarCollapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => dispatch(toggleSidebar())}
            style={{
              fontSize: '16px',
              width: 64,
              height: 64,
            }}
          />

          <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
            {/* 用户信息 */}
            <Dropdown
              menu={{
                items: userMenuItems,
                onClick: handleUserMenuClick,
              }}
              placement="bottomRight"
            >
              <div style={{ display: 'flex', alignItems: 'center', cursor: 'pointer' }}>
                <Avatar icon={<UserOutlined />} />
                <span style={{ marginLeft: '8px' }}>用户</span>
              </div>
            </Dropdown>
          </div>
        </Header>

        {/* 主内容区域 */}
        <Content
          style={{
            margin: '16px',
            padding: '24px',
            background: token.colorBgContainer,
            borderRadius: token.borderRadius,
            overflow: 'auto',
          }}
        >
          <Outlet />
        </Content>

        {/* 底部状态栏 */}
        <Footer
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
        </Footer>
      </AntLayout>
    </AntLayout>
  );
};

export default Layout;