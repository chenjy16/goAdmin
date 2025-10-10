import React, { useState, useMemo } from 'react';
import { Table, Input, Button, Space, Dropdown, message } from 'antd';
import { SearchOutlined, ReloadOutlined, DownOutlined } from '@ant-design/icons';
import type { TableProps } from 'antd/es/table';
import type { MenuProps } from 'antd';

const { Search } = Input;

interface BatchAction {
  key: string;
  label: string;
  icon?: React.ReactNode;
  danger?: boolean;
}

interface SearchableTableProps<T> extends Omit<TableProps<T>, 'dataSource' | 'title'> {
  dataSource: T[];
  searchFields?: (keyof T)[];
  searchPlaceholder?: string;
  showRefresh?: boolean;
  onRefresh?: () => void;
  refreshLoading?: boolean;
  title?: React.ReactNode;
  batchActions?: BatchAction[];
  onBatchAction?: (action: string, selectedKeys: React.Key[]) => void;
  enableBatchSelection?: boolean;
}

function SearchableTable<T extends Record<string, any>>({
  dataSource,
  columns,
  searchFields = [],
  searchPlaceholder = '搜索...',
  showRefresh = false,
  onRefresh,
  refreshLoading = false,
  title,
  batchActions = [],
  onBatchAction,
  enableBatchSelection = false,
  ...tableProps
}: SearchableTableProps<T>) {
  const [searchText, setSearchText] = useState('');
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

  // 过滤数据
  const filteredData = useMemo(() => {
    if (!searchText || searchFields.length === 0) {
      return dataSource;
    }

    return dataSource.filter(item =>
      searchFields.some(field => {
        const value = item[field];
        if (value == null) return false;
        return String(value).toLowerCase().includes(searchText.toLowerCase());
      })
    );
  }, [dataSource, searchText, searchFields]);

  // 批量操作处理
  const handleBatchAction = (actionKey: string) => {
    if (selectedRowKeys.length === 0) {
      message.warning('请先选择要操作的项目');
      return;
    }
    
    if (onBatchAction) {
      onBatchAction(actionKey, selectedRowKeys);
    }
  };

  // 批量操作菜单
  const batchMenuItems: MenuProps['items'] = batchActions.map(action => ({
    key: action.key,
    label: action.label,
    icon: action.icon,
    danger: action.danger,
    onClick: () => handleBatchAction(action.key),
  }));

  // 行选择配置
  const rowSelection = enableBatchSelection ? {
    selectedRowKeys,
    onChange: (newSelectedRowKeys: React.Key[]) => {
      setSelectedRowKeys(newSelectedRowKeys);
    },
    onSelectAll: (selected: boolean, selectedRows: T[], changeRows: T[]) => {
      if (selected) {
        const allKeys = filteredData.map((item, index) => tableProps.rowKey ? 
          (typeof tableProps.rowKey === 'function' ? tableProps.rowKey(item, index) : item[tableProps.rowKey]) : 
          index
        );
        setSelectedRowKeys(allKeys);
      } else {
        setSelectedRowKeys([]);
      }
    },
  } : undefined;

  const tableTitle = () => (
    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
      <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
        <div>{title}</div>
        {enableBatchSelection && selectedRowKeys.length > 0 && (
          <div style={{ color: '#1890ff', fontSize: '14px' }}>
            已选择 {selectedRowKeys.length} 项
          </div>
        )}
      </div>
      <Space>
        {enableBatchSelection && batchActions.length > 0 && (
          <Dropdown 
            menu={{ items: batchMenuItems }} 
            disabled={selectedRowKeys.length === 0}
          >
            <Button>
              批量操作 <DownOutlined />
            </Button>
          </Dropdown>
        )}
        {searchFields.length > 0 && (
          <Search
            placeholder={searchPlaceholder}
            allowClear
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
            style={{ width: 200 }}
            prefix={<SearchOutlined />}
          />
        )}
        {showRefresh && (
          <Button
            icon={<ReloadOutlined />}
            onClick={onRefresh}
            loading={refreshLoading}
          >
            刷新
          </Button>
        )}
      </Space>
    </div>
  );

  return (
    <Table<T>
      {...tableProps}
      dataSource={filteredData}
      columns={columns}
      rowSelection={rowSelection}
      title={title || showRefresh || searchFields.length > 0 || enableBatchSelection ? tableTitle : undefined}
      pagination={{
        pageSize: 10,
        showSizeChanger: true,
        showQuickJumper: true,
        showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
        ...tableProps.pagination,
      }}
    />
  );
}

export default SearchableTable;