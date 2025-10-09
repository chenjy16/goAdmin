# Go-SpringAI Frontend

基于 React 19 + TypeScript + Vite 构建的现代化 AI 助手前端应用，为 Go-SpringAI 后端服务提供完整的用户界面。

## 🚀 项目特性

- **现代化技术栈**: React 19 + TypeScript + Vite 7.x
- **UI 组件库**: Ant Design v5 (已适配 React 19)
- **状态管理**: Redux Toolkit + React Redux
- **路由管理**: React Router v6
- **实时通信**: Server-Sent Events (SSE) 支持
- **响应式设计**: 支持桌面端和移动端
- **类型安全**: 完整的 TypeScript 类型定义

## 📁 项目结构

```
src/
├── components/          # 公共组件
│   └── Layout/         # 布局组件
├── pages/              # 页面组件
│   ├── DashboardPage.tsx    # 仪表板
│   ├── ChatPage.tsx         # AI 聊天
│   ├── ProvidersPage.tsx    # AI 提供商管理
│   ├── MCPToolsPage.tsx     # MCP 工具管理
│   ├── AssistantPage.tsx    # AI 助手
│   └── SettingsPage.tsx     # 系统设置
├── router/             # 路由配置
├── services/           # API 服务
│   ├── api.ts          # HTTP API 客户端
│   ├── sse.ts          # SSE 连接管理
│   └── errorHandler.ts # 错误处理
├── store/              # Redux 状态管理
│   ├── index.ts        # Store 配置
│   └── slices/         # Redux Slices
├── types/              # TypeScript 类型定义
└── assets/             # 静态资源
```

## 🛠️ 核心功能

### 1. 仪表板 (Dashboard)
- 系统状态概览
- AI 提供商健康状态
- MCP 工具统计
- 最近活动记录

### 2. AI 聊天 (Chat)
- 多提供商 AI 对话 (OpenAI, Google AI)
- 实时流式响应
- 对话历史管理
- 模型参数调节

### 3. 提供商管理 (Providers)
- AI 提供商配置
- API 密钥管理
- 模型启用/禁用
- 连接状态监控

### 4. MCP 工具 (Tools)
- Model Context Protocol 工具管理
- 工具执行日志
- 工具状态监控

### 5. AI 助手 (Assistant)
- 智能助手对话
- 上下文感知
- 任务执行

## 🔧 环境要求

- **Node.js**: v20.19.5+ (推荐使用 nvm 管理版本)
- **npm**: v10.8.2+
- **现代浏览器**: 支持 ES2020+ 特性

## 📦 安装与运行

### 1. 安装依赖

```bash
npm install
```

### 2. 启动开发服务器

```bash
npm run dev
```

应用将在 http://localhost:5173 启动

### 3. 构建生产版本

```bash
npm run build
```

### 4. 预览生产构建

```bash
npm run preview
```

### 5. 代码检查

```bash
npm run lint
```

## 🔗 后端集成

前端应用需要与 Go-SpringAI 后端服务配合使用：

- **后端地址**: http://localhost:8080
- **API 前缀**: `/api/v1`
- **支持的 AI 提供商**: OpenAI, Google AI
- **实时通信**: SSE (Server-Sent Events)

## 🎨 主要依赖

### 核心依赖
- `react`: ^19.1.1 - React 框架
- `react-dom`: ^19.1.1 - React DOM 渲染
- `typescript`: ~5.9.3 - TypeScript 支持
- `vite`: ^7.1.7 - 构建工具

### UI 与样式
- `antd`: ^5.27.4 - Ant Design 组件库
- `@ant-design/v5-patch-for-react-19`: ^1.0.3 - React 19 兼容补丁

### 状态管理与路由
- `@reduxjs/toolkit`: ^2.9.0 - Redux 状态管理
- `react-redux`: ^9.2.0 - React Redux 绑定
- `react-router-dom`: ^6.30.1 - 路由管理

### 网络请求
- `axios`: ^1.12.2 - HTTP 客户端
- `uuid`: ^13.0.0 - UUID 生成

## 🔧 开发配置

### TypeScript 配置
项目使用严格的 TypeScript 配置，包含：
- `tsconfig.json` - 主配置
- `tsconfig.app.json` - 应用配置
- `tsconfig.node.json` - Node.js 配置

### ESLint 配置
使用现代化的 ESLint 配置：
- TypeScript 支持
- React Hooks 规则
- React Refresh 支持

### Vite 配置
- React 插件支持
- 快速热重载 (HMR)
- 优化的构建输出

## 🚀 部署说明

1. 构建生产版本：`npm run build`
2. 构建产物位于 `dist/` 目录
3. 可部署到任何静态文件服务器
4. 确保后端 API 服务正常运行

## 🤝 开发指南

### 添加新页面
1. 在 `src/pages/` 创建新组件
2. 在 `src/router/index.tsx` 添加路由
3. 在布局组件中添加导航菜单

### 状态管理
使用 Redux Toolkit 管理应用状态：
- 在 `src/store/slices/` 创建新的 slice
- 在组件中使用 `useAppSelector` 和 `useAppDispatch`

### API 调用
- 使用 `src/services/api.ts` 中的 HTTP 客户端
- 对于实时数据，使用 `src/services/sse.ts` 的 SSE 连接

## 📄 许可证

本项目遵循 MIT 许可证。
