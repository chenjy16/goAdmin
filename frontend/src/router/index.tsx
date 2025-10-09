import { createBrowserRouter, Navigate } from 'react-router-dom';
import Layout from '../components/Layout';
import ChatPage from '../pages/ChatPage';
import ProvidersPage from '../pages/ProvidersPage';
import MCPToolsPage from '../pages/MCPToolsPage';
import AssistantPage from '../pages/AssistantPage';
import SettingsPage from '../pages/SettingsPage';
import DashboardPage from '../pages/DashboardPage';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <Layout />,
    children: [
      {
        index: true,
        element: <Navigate to="/dashboard" replace />,
      },
      {
        path: 'dashboard',
        element: <DashboardPage />,
      },
      {
        path: 'chat',
        element: <ChatPage />,
      },
      {
        path: 'providers',
        element: <ProvidersPage />,
      },
      {
        path: 'tools',
        element: <MCPToolsPage />,
      },
      {
        path: 'assistant',
        element: <AssistantPage />,
      },
      {
        path: 'settings',
        element: <SettingsPage />,
      },
    ],
  },
]);

export default router;