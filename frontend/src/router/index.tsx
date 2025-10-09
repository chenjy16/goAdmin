import { createBrowserRouter, Navigate } from 'react-router-dom';
import Layout from '../components/Layout';
import AssistantPage from '../pages/AssistantPage';
import SettingsPage from '../pages/SettingsPage';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <Layout />,
    children: [
      {
        index: true,
        element: <Navigate to="/assistant" replace />,
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