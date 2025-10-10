import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';
import apiService from '../../services/api';
import type { MCPTool, MCPMessage, MCPExecuteRequest } from '../../types/api';
import { createCrudSlice } from '../utils/asyncSliceFactory';

// MCP状态接口
export interface MCPState {
  tools: MCPTool[];
  logs: MCPMessage[];
  isInitialized: boolean;
  isLoading: boolean;
  error: string | null;
  executionResults: Record<string, any>;
}

// 异步thunks
export const initializeMCP = createAsyncThunk(
  'mcp/initialize',
  async () => {
    const response = await apiService.initializeMCP();
    return response;
  }
);

// 使用CRUD slice管理MCP工具
const mcpToolsConfig = createCrudSlice<MCPTool>({
  name: 'mcpTools',
  idField: 'name',
  initialItems: []
});

export const fetchMCPTools = createAsyncThunk(
  'mcp/fetchTools',
  async () => {
    const response = await apiService.getMCPTools();
    return response.data?.tools || [];
  }
);

// 导出工具管理相关的actions
export const {
  addItem: addMCPTool,
  updateItem: updateMCPTool,
  removeItem: removeMCPTool,
  setItems: setMCPTools,
  reset: resetMCPTools
} = mcpToolsConfig.actions;

export const executeMCPTool = createAsyncThunk(
  'mcp/executeTool',
  async (request: MCPExecuteRequest) => {
    const response = await apiService.executeMCPTool(request);
    return { request, response };
  }
);

export const fetchMCPLogs = createAsyncThunk(
  'mcp/fetchLogs',
  async () => {
    const logs = await apiService.getMCPLogs();
    return logs;
  }
);

export const checkMCPStatus = createAsyncThunk(
  'mcp/checkStatus',
  async (_, { rejectWithValue }) => {
    try {
      const data = await apiService.getMCPStatus();
      
      if (!data || typeof data.initialized !== 'boolean') {
        console.error('Invalid MCP status response:', data);
        return rejectWithValue('Invalid response format');
      }
      
      return data;
    } catch (error) {
      console.error('Failed to check MCP status:', error);
      return rejectWithValue(error instanceof Error ? error.message : 'Failed to check MCP status');
    }
  }
);

const initialState: MCPState = {
  tools: [],
  logs: [],
  isInitialized: false,
  isLoading: false,
  error: null,
  executionResults: {},
};

const mcpSlice = createSlice({
  name: 'mcp',
  initialState,
  reducers: {
    addLog: (state, action: PayloadAction<MCPMessage>) => {
      state.logs.unshift(action.payload);
      // 保持最新的100条日志
      if ((state.logs || []).length > 100) {
        state.logs = (state.logs || []).slice(0, 100);
      }
    },
    clearLogs: (state) => {
      state.logs = [];
    },
    clearError: (state) => {
      state.error = null;
    },
    setExecutionResult: (state, action: PayloadAction<{
      toolName: string;
      result: any;
    }>) => {
      const { toolName, result } = action.payload;
      state.executionResults[toolName] = result;
    },
    clearExecutionResults: (state) => {
      state.executionResults = {};
    },
    setState: (state, action: PayloadAction<Partial<MCPState>>) => {
      Object.assign(state, action.payload);
    },
  },
  extraReducers: (builder) => {
    builder
      // initializeMCP
      .addCase(initializeMCP.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(initializeMCP.fulfilled, (state) => {
        state.isLoading = false;
        state.isInitialized = true;
      })
      .addCase(initializeMCP.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || 'MCP初始化失败';
      })
      // fetchMCPTools
      .addCase(fetchMCPTools.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchMCPTools.fulfilled, (state, action) => {
        state.isLoading = false;
        state.tools = action.payload;
      })
      .addCase(fetchMCPTools.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || '获取工具列表失败';
      })
      // executeMCPTool
      .addCase(executeMCPTool.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(executeMCPTool.fulfilled, (state, action) => {
        state.isLoading = false;
        const { request, response } = action.payload;
        state.executionResults[request.name] = response;
      })
      .addCase(executeMCPTool.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || '工具执行失败';
      })
      // fetchMCPLogs
      .addCase(fetchMCPLogs.fulfilled, (state, action) => {
        state.logs = Array.isArray(action.payload) ? action.payload : [];
      })
      .addCase(fetchMCPLogs.rejected, (state, action) => {
        state.error = action.error.message || '获取日志失败';
      })
      // checkMCPStatus
      .addCase(checkMCPStatus.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(checkMCPStatus.fulfilled, (state, action) => {
        state.isInitialized = action.payload.initialized;
        state.isLoading = false;
      })
      .addCase(checkMCPStatus.rejected, (state, action) => {
        state.error = action.error.message || '检查MCP状态失败';
        state.isLoading = false;
      });
  },
});

export const {
  addLog,
  clearLogs,
  clearError,
  setExecutionResult,
  clearExecutionResults,
  setState,
} = mcpSlice.actions;

export default mcpSlice.reducer;