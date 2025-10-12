import { useState, useCallback, useRef, useEffect } from 'react';
import { Modal } from 'antd';
import type { ModalProps } from 'antd';
import { useTranslation } from 'react-i18next';

// 模态框配置
export interface UseModalConfig extends Omit<ModalProps, 'open' | 'onOk' | 'onCancel'> {
  onOpen?: () => void;
  onClose?: () => void;
  onConfirm?: () => Promise<any> | any;
  closeOnConfirm?: boolean;
  resetOnClose?: boolean;
}

// 模态框状态
export interface ModalState {
  open: boolean;
  loading: boolean;
  data: any;
}

// 基础模态框hook
export function useModal(config: UseModalConfig = {}) {
  const {
    onOpen,
    onClose,
    onConfirm,
    closeOnConfirm = true,
    resetOnClose = false,
    ...modalProps
  } = config;

  const [state, setState] = useState<ModalState>({
    open: false,
    loading: false,
    data: null,
  });

  const dataRef = useRef<any>(null);

  // 打开模态框
  const open = useCallback((data?: any) => {
    setState(prev => ({
      ...prev,
      open: true,
      data,
    }));
    dataRef.current = data;
    onOpen?.();
  }, [onOpen]);

  // 关闭模态框
  const close = useCallback(() => {
    setState(prev => ({
      ...prev,
      open: false,
      loading: false,
      data: resetOnClose ? null : prev.data,
    }));
    
    if (resetOnClose) {
      dataRef.current = null;
    }
    
    onClose?.();
  }, [onClose, resetOnClose]);

  // 确认操作
  const confirm = useCallback(async () => {
    if (!onConfirm) {
      close();
      return;
    }

    setState(prev => ({
      ...prev,
      loading: true,
    }));

    try {
      await onConfirm();
      
      if (closeOnConfirm) {
        close();
      } else {
        setState(prev => ({
          ...prev,
          loading: false,
        }));
      }
    } catch (error) {
      setState(prev => ({
        ...prev,
        loading: false,
      }));
      throw error;
    }
  }, [onConfirm, closeOnConfirm, close]);

  // 切换模态框状态
  const toggle = useCallback((data?: any) => {
    if (state.open) {
      close();
    } else {
      open(data);
    }
  }, [state.open, open, close]);

  return {
    // 状态
    ...state,
    
    // 操作方法
    open,
    close,
    confirm,
    toggle,
    
    // 模态框属性
    modalProps: {
      ...modalProps,
      open: state.open,
      confirmLoading: state.loading,
      onOk: confirm,
      onCancel: close,
    },
    
    // 便捷属性
    data: state.data || dataRef.current,
  };
}

// 确认模态框hook
export function useConfirmModal(config: {
  title?: string;
  content?: string;
  onConfirm?: () => Promise<any> | any;
  okText?: string;
  cancelText?: string;
  type?: 'info' | 'success' | 'error' | 'warning' | 'confirm';
} = {}) {
  const { t } = useTranslation();
  const {
    title = t('modals.confirm'),
    content = t('modals.confirmAction'),
    onConfirm,
    okText = t('modals.ok'),
    cancelText = t('modals.cancel'),
    type = 'confirm',
  } = config;

  const [loading, setLoading] = useState(false);

  const show = useCallback(async (customConfig?: {
    title?: string;
    content?: string;
    onConfirm?: () => Promise<any> | any;
  }) => {
    const finalTitle = customConfig?.title || title;
    const finalContent = customConfig?.content || content;
    const finalOnConfirm = customConfig?.onConfirm || onConfirm;

    return new Promise<boolean>((resolve) => {
      const modal = Modal[type]({
        title: finalTitle,
        content: finalContent,
        okText,
        cancelText,
        onOk: async () => {
          if (!finalOnConfirm) {
            resolve(true);
            return;
          }

          setLoading(true);
          try {
            await finalOnConfirm();
            setLoading(false);
            resolve(true);
          } catch (error) {
            setLoading(false);
            resolve(false);
            throw error;
          }
        },
        onCancel: () => {
          resolve(false);
        },
      });
    });
  }, [title, content, onConfirm, okText, cancelText, type, loading]);

  return {
    show,
    loading,
  };
}

// 表单模态框hook
export function useFormModal<T = any>(config: UseModalConfig & {
  initialValues?: Partial<T>;
  onSubmit?: (values: T) => Promise<any> | any;
} = {}) {
  const { initialValues, onSubmit, ...modalConfig } = config;
  
  const [formData, setFormData] = useState<Partial<T>>(initialValues || {});
  
  const modal = useModal({
    ...modalConfig,
    onConfirm: async () => {
      if (onSubmit) {
        await onSubmit(formData as T);
      }
    },
    onOpen: () => {
      if (initialValues) {
        setFormData(initialValues);
      }
      modalConfig.onOpen?.();
    },
    onClose: () => {
      if (modalConfig.resetOnClose) {
        setFormData(initialValues || {});
      }
      modalConfig.onClose?.();
    },
  });

  const openWithData = useCallback((data: Partial<T>) => {
    setFormData(data);
    modal.open(data);
  }, [modal]);

  const updateFormData = useCallback((updates: Partial<T>) => {
    setFormData(prev => ({
      ...prev,
      ...updates,
    }));
  }, []);

  const setFormField = useCallback((field: keyof T, value: any) => {
    setFormData(prev => ({
      ...prev,
      [field]: value,
    }));
  }, []);

  return {
    ...modal,
    formData,
    updateFormData,
    setFormField,
    openWithData,
  };
}

// 列表操作模态框hook
export function useListModal<T = any>(config: {
  onAdd?: (item: T) => Promise<any> | any;
  onEdit?: (item: T) => Promise<any> | any;
  onDelete?: (item: T) => Promise<any> | any;
  onView?: (item: T) => void;
  addTitle?: string;
  editTitle?: string;
  deleteTitle?: string;
  viewTitle?: string;
} = {}) {
  const { t } = useTranslation();
  const {
    onAdd,
    onEdit,
    onDelete,
    onView,
    addTitle = t('modals.add'),
    editTitle = t('modals.edit'),
    deleteTitle = t('modals.delete'),
    viewTitle = t('modals.view'),
  } = config;

  const [mode, setMode] = useState<'add' | 'edit' | 'delete' | 'view'>('add');
  const [currentItem, setCurrentItem] = useState<T | null>(null);

  // 添加模态框
  const addModal = useFormModal<T>({
    title: addTitle,
    onSubmit: onAdd,
    resetOnClose: true,
  });

  // 编辑模态框
  const editModal = useFormModal<T>({
    title: editTitle,
    onSubmit: onEdit,
  });

  // 删除确认模态框
  const deleteModal = useConfirmModal({
    title: deleteTitle,
    content: t('modals.deleteConfirm'),
    onConfirm: async () => {
      if (currentItem && onDelete) {
        await onDelete(currentItem);
      }
    },
  });

  // 查看模态框
  const viewModal = useModal({
    title: viewTitle,
    footer: null,
  });

  // 打开添加模态框
  const openAdd = useCallback(() => {
    setMode('add');
    setCurrentItem(null);
    addModal.open();
  }, [addModal]);

  // 打开编辑模态框
  const openEdit = useCallback((item: T) => {
    setMode('edit');
    setCurrentItem(item);
    editModal.openWithData(item as Partial<T>);
  }, [editModal]);

  // 打开删除确认模态框
  const openDelete = useCallback((item: T) => {
    setMode('delete');
    setCurrentItem(item);
    deleteModal.show();
  }, [deleteModal]);

  // 打开查看模态框
  const openView = useCallback((item: T) => {
    setMode('view');
    setCurrentItem(item);
    viewModal.open(item);
    onView?.(item);
  }, [viewModal, onView]);

  return {
    // 当前状态
    mode,
    currentItem,
    
    // 模态框实例
    addModal,
    editModal,
    deleteModal,
    viewModal,
    
    // 操作方法
    openAdd,
    openEdit,
    openDelete,
    openView,
  };
}

// 批量操作模态框hook
export function useBatchModal<T = any>(config: {
  onBatchDelete?: (items: T[]) => Promise<any> | any;
  onBatchEdit?: (items: T[], updates: Partial<T>) => Promise<any> | any;
  onBatchExport?: (items: T[]) => Promise<any> | any;
} = {}) {
  const { t } = useTranslation();
  const { onBatchDelete, onBatchEdit, onBatchExport } = config;
  
  const [selectedItems, setSelectedItems] = useState<T[]>([]);

  // 批量删除确认模态框
  const batchDeleteModal = useConfirmModal({
    title: t('modals.batchDelete'),
    content: t('modals.batchDeleteConfirm', { count: selectedItems.length }),
    onConfirm: async () => {
      if (onBatchDelete) {
        await onBatchDelete(selectedItems);
      }
    },
  });

  // 批量编辑模态框
  const batchEditModal = useFormModal({
    title: t('modals.batchEdit'),
    onSubmit: async (updates: Partial<T>) => {
      if (onBatchEdit) {
        await onBatchEdit(selectedItems, updates);
      }
    },
  });

  // 打开批量删除确认
  const openBatchDelete = useCallback((items: T[]) => {
    setSelectedItems(items);
    batchDeleteModal.show({
      content: t('modals.batchDeleteConfirm', { count: items.length }),
    });
  }, [batchDeleteModal, t]);

  // 打开批量编辑
  const openBatchEdit = useCallback((items: T[]) => {
    setSelectedItems(items);
    batchEditModal.open();
  }, [batchEditModal]);

  // 批量导出
  const batchExport = useCallback(async (items: T[]) => {
    setSelectedItems(items);
    if (onBatchExport) {
      await onBatchExport(items);
    }
  }, [onBatchExport]);

  return {
    selectedItems,
    batchDeleteModal,
    batchEditModal,
    openBatchDelete,
    openBatchEdit,
    batchExport,
  };
}