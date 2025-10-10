import React from 'react';
import { Layout, Button, Avatar, Dropdown, theme } from 'antd';
import type { MenuProps } from 'antd';
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  UserOutlined,
  LogoutOutlined,
} from '@ant-design/icons';
import { useAppDispatch, useAppSelector } from '../../store';
import { toggleSidebar } from '../../store/slices/uiSlice';

const { Header: AntHeader } = Layout;

const Header: React.FC = () => {
  const dispatch = useAppDispatch();
  const { token } = theme.useToken();
  const { sidebarCollapsed } = useAppSelector(state => state.ui);

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

  return (
    <AntHeader
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
    </AntHeader>
  );
};

export default Header;