import React from 'react';
import { Button, Space, Switch, Tag, Tooltip } from 'antd';
import { 
  EditOutlined, 
  DeleteOutlined, 
  EyeOutlined, 
  SettingOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined 
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';

// 通用列类型
export interface BaseColumnConfig<T> {
  dataIndex: keyof T;
  title: string;
  width?: number;
  sorter?: boolean;
  ellipsis?: boolean;
}

// 状态列配置
export interface StatusColumnConfig<T> extends BaseColumnConfig<T> {
  statusMap?: Record<string, { color: string; text: string; icon?: React.ReactNode }>;
}

// 操作列配置
export interface ActionColumnConfig<T> {
  width?: number;
  actions: Array<{
    key: string;
    label: string;
    icon?: React.ReactNode;
    type?: 'primary' | 'default' | 'dashed' | 'link' | 'text';
    danger?: boolean;
    onClick: (record: T) => void;
    visible?: (record: T) => boolean;
    disabled?: (record: T) => boolean;
  }>;
}

// 开关列配置
export interface SwitchColumnConfig<T> extends BaseColumnConfig<T> {
  onChange: (checked: boolean, record: T) => void;
  disabled?: (record: T) => boolean;
}

// 标签列配置
export interface TagColumnConfig<T> extends BaseColumnConfig<T> {
  tagMap?: Record<string, string>; // value -> color
  multiple?: boolean; // 是否支持多个标签
}

/**
 * 创建基础文本列
 */
export function createTextColumn<T>(config: BaseColumnConfig<T>) {
  return {
    title: config.title,
    dataIndex: config.dataIndex as string,
    key: config.dataIndex as string,
    width: config.width,
    sorter: config.sorter,
    ellipsis: config.ellipsis ? { showTitle: false } : false,
    render: config.ellipsis 
      ? (text: string) => (
          <Tooltip placement="topLeft" title={text}>
            {text}
          </Tooltip>
        )
      : undefined,
  };
}

/**
 * 创建状态列
 */
export function createStatusColumn<T>(config: StatusColumnConfig<T>) {
  const defaultStatusMap: Record<string, { color: string; text: string; icon?: React.ReactNode }> = {
    active: { color: 'green', text: '活跃', icon: <CheckCircleOutlined /> },
    inactive: { color: 'red', text: '非活跃', icon: <ExclamationCircleOutlined /> },
    enabled: { color: 'green', text: '启用', icon: <CheckCircleOutlined /> },
    disabled: { color: 'red', text: '禁用', icon: <ExclamationCircleOutlined /> },
  };

  const statusMap = { ...defaultStatusMap, ...config.statusMap };

  return {
    title: config.title,
    dataIndex: config.dataIndex as string,
    key: config.dataIndex as string,
    width: config.width || 100,
    sorter: config.sorter,
    render: (status: string) => {
      const statusConfig = statusMap[status] || { color: 'default', text: status };
      return (
        <Tag color={statusConfig.color} icon={statusConfig.icon}>
          {statusConfig.text}
        </Tag>
      );
    },
  };
}

/**
 * 创建开关列
 */
export function createSwitchColumn<T>(config: SwitchColumnConfig<T>) {
  return {
    title: config.title,
    dataIndex: config.dataIndex as string,
    key: config.dataIndex as string,
    width: config.width || 80,
    render: (checked: boolean, record: T) => (
      <Switch
        checked={checked}
        onChange={(value) => config.onChange(value, record)}
        disabled={config.disabled ? config.disabled(record) : false}
        size="small"
      />
    ),
  };
}

/**
 * 创建标签列
 */
export function createTagColumn<T>(config: TagColumnConfig<T>) {
  return {
    title: config.title,
    dataIndex: config.dataIndex as string,
    key: config.dataIndex as string,
    width: config.width,
    render: (value: string | string[]) => {
      if (!value) return null;
      
      const values = Array.isArray(value) ? value : [value];
      
      return (
        <Space size={[0, 4]} wrap>
          {values.map((val, index) => (
            <Tag 
              key={index} 
              color={config.tagMap?.[val] || 'default'}
            >
              {val}
            </Tag>
          ))}
        </Space>
      );
    },
  };
}

/**
 * 创建操作列
 */
export function createActionColumn<T>(config: ActionColumnConfig<T>) {
  return {
    title: '操作',
    key: 'actions',
    width: config.width || 150,
    render: (_: any, record: T) => (
      <Space size="small">
        {config.actions.map((action) => {
          const visible = action.visible ? action.visible(record) : true;
          const disabled = action.disabled ? action.disabled(record) : false;
          
          if (!visible) return null;
          
          return (
            <Button
              key={action.key}
              type={action.type || 'link'}
              size="small"
              icon={action.icon}
              danger={action.danger}
              disabled={disabled}
              onClick={() => action.onClick(record)}
            >
              {action.label}
            </Button>
          );
        })}
      </Space>
    ),
  };
}

/**
 * 创建常用的操作按钮配置
 */
export const commonActions = {
  view: (onClick: (record: any) => void) => ({
    key: 'view',
    label: '查看',
    icon: <EyeOutlined />,
    onClick,
  }),
  
  edit: (onClick: (record: any) => void) => ({
    key: 'edit',
    label: '编辑',
    icon: <EditOutlined />,
    onClick,
  }),
  
  delete: (onClick: (record: any) => void) => ({
    key: 'delete',
    label: '删除',
    icon: <DeleteOutlined />,
    danger: true,
    onClick,
  }),
  
  settings: (onClick: (record: any) => void) => ({
    key: 'settings',
    label: '设置',
    icon: <SettingOutlined />,
    onClick,
  }),
};

/**
 * 合并列配置
 */
export function mergeColumns<T>(...columnGroups: ColumnsType<T>[]): ColumnsType<T> {
  return columnGroups.flat();
}