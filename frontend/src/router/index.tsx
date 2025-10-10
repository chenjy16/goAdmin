import React from 'react';
import { createBrowserRouter, Navigate } from 'react-router-dom';
import Layout from '../components/Layout';
import SettingsPage from '../pages/SettingsPage';
import MCPToolsPage from '../pages/MCPToolsPage';
import AssistantPage from '../pages/AssistantPage';

const router = createBrowserRouter([
  {
    path: '/',
    element: <Layout />,
    children: [
      {
        index: true,
        element: <Navigate to="/assistant" replace />,
      },
      {
        path: 'settings',
        element: <SettingsPage />,
      },
      {
        path: 'mcp-tools',
        element: <MCPToolsPage />,
      },
      {
        path: 'assistant',
        element: <AssistantPage />,
      },
    ],
  },
]);

export default router;