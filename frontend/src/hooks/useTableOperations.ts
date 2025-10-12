import { useState, useCallback, useMemo } from 'react';
import { message } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { useTranslation } from 'react-i18next';

interface UseTableOperationsOptions<T> {
  searchFields?: (keyof T)[];
  defaultPageSize?: number;
  enableSelection?: boolean;
}

interface TableOperationsState {
  searchText: string;
  selectedRowKeys: React.Key[];
  currentPage: number;
  pageSize: number;
}

export function useTableOperations<T extends Record<string, any>>(
  dataSource: T[],
  options: UseTableOperationsOptions<T> = {}
) {
  const { t } = useTranslation();
  const {
    searchFields = [],
    defaultPageSize = 10,
    enableSelection = false,
  } = options;

  const [state, setState] = useState<TableOperationsState>({
    searchText: '',
    selectedRowKeys: [],
    currentPage: 1,
    pageSize: defaultPageSize,
  });

  // 过滤数据
  const filteredData = useMemo(() => {
    if (!state.searchText || searchFields.length === 0) {
      return dataSource;
    }

    return dataSource.filter(item =>
      searchFields.some(field => {
        const value = item[field];
        if (value == null) return false;
        return String(value).toLowerCase().includes(state.searchText.toLowerCase());
      })
    );
  }, [dataSource, state.searchText, searchFields]);

  // 分页数据
  const paginatedData = useMemo(() => {
    const startIndex = (state.currentPage - 1) * state.pageSize;
    const endIndex = startIndex + state.pageSize;
    return filteredData.slice(startIndex, endIndex);
  }, [filteredData, state.currentPage, state.pageSize]);

  // 搜索处理
  const handleSearch = useCallback((value: string) => {
    setState(prev => ({
      ...prev,
      searchText: value,
      currentPage: 1, // 重置到第一页
    }));
  }, []);

  // 清除搜索
  const clearSearch = useCallback(() => {
    setState(prev => ({
      ...prev,
      searchText: '',
      currentPage: 1,
    }));
  }, []);

  // 选择处理
  const handleSelectionChange = useCallback((newSelectedRowKeys: React.Key[]) => {
    if (!enableSelection) return;
    
    setState(prev => ({
      ...prev,
      selectedRowKeys: newSelectedRowKeys,
    }));
  }, [enableSelection]);

  // 全选处理
  const handleSelectAll = useCallback((selected: boolean) => {
    if (!enableSelection) return;
    
    if (selected) {
      const allKeys = filteredData.map((_, index) => index);
      setState(prev => ({
        ...prev,
        selectedRowKeys: allKeys,
      }));
    } else {
      setState(prev => ({
        ...prev,
        selectedRowKeys: [],
      }));
    }
  }, [enableSelection, filteredData]);

  // 批量操作
  const handleBatchOperation = useCallback((
    operation: (selectedItems: T[]) => Promise<void>,
    operationName: string = t('common.operation')
  ) => {
    if (state.selectedRowKeys.length === 0) {
      message.warning(t('common.pleaseSelectItems', { operation: operationName }));
      return;
    }

    const selectedItems = state.selectedRowKeys.map(key => 
      filteredData[key as number]
    ).filter(Boolean);

    return operation(selectedItems);
  }, [state.selectedRowKeys, filteredData]);

  // 分页处理
  const handlePageChange = useCallback((page: number, size?: number) => {
    setState(prev => ({
      ...prev,
      currentPage: page,
      pageSize: size || prev.pageSize,
    }));
  }, []);

  // 重置状态
  const reset = useCallback(() => {
    setState({
      searchText: '',
      selectedRowKeys: [],
      currentPage: 1,
      pageSize: defaultPageSize,
    });
  }, [defaultPageSize]);

  // 行选择配置
  const rowSelection = enableSelection ? {
    selectedRowKeys: state.selectedRowKeys,
    onChange: handleSelectionChange,
    onSelectAll: handleSelectAll,
  } : undefined;

  // 分页配置
  const pagination = {
    current: state.currentPage,
    pageSize: state.pageSize,
    total: filteredData.length,
    showSizeChanger: true,
    showQuickJumper: true,
    showTotal: (total: number, range: [number, number]) =>
      t('common.paginationTotal', { start: range[0], end: range[1], total }),
    onChange: handlePageChange,
  };

  return {
    // 状态
    searchText: state.searchText,
    selectedRowKeys: state.selectedRowKeys,
    currentPage: state.currentPage,
    pageSize: state.pageSize,
    
    // 数据
    filteredData,
    paginatedData,
    
    // 操作方法
    handleSearch,
    clearSearch,
    handleSelectionChange,
    handleSelectAll,
    handleBatchOperation,
    handlePageChange,
    reset,
    
    // 配置
    rowSelection,
    pagination,
  };
}

export default useTableOperations;