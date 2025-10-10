import { useState, useCallback, useMemo } from 'react';
import type { TableProps } from 'antd';
import { useAsyncData, usePaginatedData, type UseAsyncDataConfig } from './useAsyncData';

// 表格配置
export interface UseTableConfig<T> extends UseAsyncDataConfig<T[]> {
  pagination?: boolean | {
    pageSize?: number;
    showSizeChanger?: boolean;
    showQuickJumper?: boolean;
    showTotal?: (total: number, range: [number, number]) => string;
  };
  selection?: {
    type?: 'checkbox' | 'radio';
    onChange?: (selectedRowKeys: React.Key[], selectedRows: T[]) => void;
    getCheckboxProps?: (record: T) => any;
  };
  sorting?: boolean;
  filtering?: boolean;
}

// 表格状态
export interface TableState<T> {
  selectedRowKeys: React.Key[];
  selectedRows: T[];
  sorter: any;
  filters: Record<string, any>;
}

// 基础表格hook
export function useTable<T extends Record<string, any>>(
  dataSource: T[] | (() => Promise<T[]>),
  config: UseTableConfig<T> = {}
) {
  const [tableState, setTableState] = useState<TableState<T>>({
    selectedRowKeys: [],
    selectedRows: [],
    sorter: {},
    filters: {},
  });

  // 处理数据源
  const asyncData = typeof dataSource === 'function' 
    ? useAsyncData(dataSource, config)
    : { data: dataSource, loading: false, error: null, initialized: true, execute: () => Promise.resolve(dataSource), reset: () => {}, clearError: () => {}, refresh: () => Promise.resolve(dataSource) };

  // 行选择配置
  const rowSelection = useMemo(() => {
    if (!config.selection) return undefined;

    return {
      type: config.selection.type || 'checkbox',
      selectedRowKeys: tableState.selectedRowKeys,
      onChange: (selectedRowKeys: React.Key[], selectedRows: T[]) => {
        setTableState(prev => ({
          ...prev,
          selectedRowKeys,
          selectedRows,
        }));
        
        if (config.selection?.onChange) {
          config.selection.onChange(selectedRowKeys, selectedRows);
        }
      },
      getCheckboxProps: config.selection.getCheckboxProps,
    };
  }, [config.selection, tableState.selectedRowKeys]);

  // 分页配置
  const pagination = useMemo(() => {
    if (config.pagination === false) return false;
    
    const paginationConfig = typeof config.pagination === 'object' ? config.pagination : {};
    
    return {
      pageSize: 10,
      showSizeChanger: true,
      showQuickJumper: true,
      showTotal: (total: number, range: [number, number]) => 
        `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
      ...paginationConfig,
    };
  }, [config.pagination]);

  // 表格变化处理
  const handleTableChange = useCallback((
    paginationInfo: any,
    filters: any,
    sorter: any
  ) => {
    setTableState(prev => ({
      ...prev,
      filters,
      sorter,
    }));
  }, []);

  // 清除选择
  const clearSelection = useCallback(() => {
    setTableState(prev => ({
      ...prev,
      selectedRowKeys: [],
      selectedRows: [],
    }));
  }, []);

  // 选择所有
  const selectAll = useCallback(() => {
    if (!asyncData.data) return;
    
    const allKeys = asyncData.data.map((item, index) => item.key || item.id || index);
    setTableState(prev => ({
      ...prev,
      selectedRowKeys: allKeys,
      selectedRows: asyncData.data || [],
    }));
  }, [asyncData.data]);

  // 反选
  const invertSelection = useCallback(() => {
    if (!asyncData.data) return;
    
    const allKeys = asyncData.data.map((item, index) => item.key || item.id || index);
    const unselectedKeys = allKeys.filter(key => !tableState.selectedRowKeys.includes(key));
    const unselectedRows = asyncData.data.filter((item, index) => {
      const key = item.key || item.id || index;
      return unselectedKeys.includes(key);
    });
    
    setTableState(prev => ({
      ...prev,
      selectedRowKeys: unselectedKeys,
      selectedRows: unselectedRows,
    }));
  }, [asyncData.data, tableState.selectedRowKeys]);

  // 表格属性
  const tableProps: TableProps<T> = {
    dataSource: asyncData.data || [],
    loading: asyncData.loading,
    pagination,
    rowSelection,
    onChange: handleTableChange,
  };

  return {
    // 数据相关
    ...asyncData,
    
    // 表格状态
    ...tableState,
    
    // 表格属性
    tableProps,
    
    // 操作方法
    clearSelection,
    selectAll,
    invertSelection,
    
    // 便捷方法
    hasSelection: tableState.selectedRowKeys.length > 0,
    selectedCount: tableState.selectedRowKeys.length,
  };
}

// 分页表格hook
export function usePaginatedTable<T extends Record<string, any>>(
  fetchData: (pagination: { page: number; pageSize: number; total?: number }) => Promise<{ data: T[]; total: number }>,
  config: UseTableConfig<T> = {}
) {
  const [tableState, setTableState] = useState<TableState<T>>({
    selectedRowKeys: [],
    selectedRows: [],
    sorter: {},
    filters: {},
  });

  // 创建适配器函数来处理类型不匹配
  const adaptedFetchData = useCallback(async (pagination: { page?: number; pageSize?: number; total?: number }) => {
    return fetchData({
      page: pagination.page || 1,
      pageSize: pagination.pageSize || 10,
      total: pagination.total,
    });
  }, [fetchData]);

  // 提取分页数据配置，排除表格特有的配置
  const { pagination: _, selection: __, sorting: ___, filtering: ____, ...paginatedConfig } = config;
  const paginatedData = usePaginatedData(adaptedFetchData, paginatedConfig);

  // 行选择配置
  const rowSelection = useMemo(() => {
    if (!config.selection) return undefined;

    return {
      type: config.selection.type || 'checkbox',
      selectedRowKeys: tableState.selectedRowKeys,
      onChange: (selectedRowKeys: React.Key[], selectedRows: T[]) => {
        setTableState(prev => ({
          ...prev,
          selectedRowKeys,
          selectedRows,
        }));
        
        if (config.selection?.onChange) {
          config.selection.onChange(selectedRowKeys, selectedRows);
        }
      },
      getCheckboxProps: config.selection.getCheckboxProps,
    };
  }, [config.selection, tableState.selectedRowKeys]);

  // 分页配置
  const pagination = useMemo(() => {
    if (config.pagination === false) return false;
    
    const paginationConfig = typeof config.pagination === 'object' ? config.pagination : {};
    
    return {
      current: paginatedData.pagination.page,
      pageSize: paginatedData.pagination.pageSize,
      total: paginatedData.pagination.total,
      showSizeChanger: true,
      showQuickJumper: true,
      showTotal: (total: number, range: [number, number]) => 
        `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
      onChange: paginatedData.changePage,
      onShowSizeChange: (_: number, size: number) => paginatedData.changePageSize(size),
      ...paginationConfig,
    };
  }, [config.pagination, paginatedData.pagination, paginatedData.changePage, paginatedData.changePageSize]);

  // 表格变化处理
  const handleTableChange = useCallback((
    paginationInfo: any,
    filters: any,
    sorter: any
  ) => {
    setTableState(prev => ({
      ...prev,
      filters,
      sorter,
    }));
  }, []);

  // 清除选择
  const clearSelection = useCallback(() => {
    setTableState(prev => ({
      ...prev,
      selectedRowKeys: [],
      selectedRows: [],
    }));
  }, []);

  // 表格属性
  const tableProps: TableProps<T> = {
    dataSource: paginatedData.data || [],
    loading: paginatedData.loading,
    pagination,
    rowSelection,
    onChange: handleTableChange,
  };

  return {
    // 数据相关
    ...paginatedData,
    
    // 表格状态
    ...tableState,
    
    // 表格属性
    tableProps,
    
    // 操作方法
    clearSelection,
    
    // 便捷方法
    hasSelection: tableState.selectedRowKeys.length > 0,
    selectedCount: tableState.selectedRowKeys.length,
  };
}