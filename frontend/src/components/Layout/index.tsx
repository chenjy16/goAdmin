import React, { useEffect } from 'react';
import { Outlet, useLocation } from 'react-router-dom';
import { Layout as AntLayout, theme } from 'antd';
import { useAppDispatch } from '../../store';
import { setCurrentPage } from '../../store/slices/uiSlice';
import Header from './Header';
import Sidebar from './Sidebar';
import Footer from './Footer';
import './index.css';

const { Content } = AntLayout;

const Layout: React.FC = () => {
  const dispatch = useAppDispatch();
  const location = useLocation();
  const { token } = theme.useToken();

  // 根据路径更新当前页面
  useEffect(() => {
    const path = location.pathname.slice(1) || 'assistant';
    dispatch(setCurrentPage(path));
  }, [location.pathname, dispatch]);

  return (
    <AntLayout style={{ minHeight: '100vh' }}>
      {/* 侧边栏 */}
      <Sidebar />

      <AntLayout>
        {/* 顶部导航栏 */}
        <Header />

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
        <Footer />
      </AntLayout>
    </AntLayout>
  );
};

export default Layout;