import React, { useState, useMemo } from 'react';
import { Table, Input, Button, Space, Dropdown, message } from 'antd';
import { SearchOutlined, ReloadOutlined, DownOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import type { TableProps, ColumnsType } from 'antd/es/table';
import type { MenuProps } from 'antd/es/menu';

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
  searchPlaceholder,
  showRefresh = false,
  onRefresh,
  refreshLoading = false,
  title,
  batchActions = [],
  onBatchAction,
  enableBatchSelection = false,
  ...tableProps
}: SearchableTableProps<T>) {
  const { t } = useTranslation();
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
      message.warning(t('common.pleaseSelectItems', { operation: t('common.operation') }));
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
  const rowSelectionConfig = enableBatchSelection ? {
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
            {t('common.selectedItems', { count: selectedRowKeys.length })}
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
              {t('common.batchOperation')} <DownOutlined />
            </Button>
          </Dropdown>
        )}
        {searchFields.length > 0 && (
          <Search
            placeholder={searchPlaceholder || t('common.search')}
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
            {t('common.refresh')}
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
      rowSelection={rowSelectionConfig}
      title={title || showRefresh || searchFields.length > 0 || enableBatchSelection ? tableTitle : undefined}
      pagination={{
        pageSize: 10,
        showSizeChanger: true,
        showQuickJumper: true,
        showTotal: (total, range) => t('common.paginationTotal', { start: range[0], end: range[1], total }),
        ...tableProps.pagination,
      }}
    />
  );
}

export default SearchableTable;