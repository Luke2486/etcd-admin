import apiClient from './client'
import type { 
  Connection,
  CreateConnectionRequest,
  UpdateConnectionRequest,
  TestConnectionResponse,
  ApiResponse 
} from '@/types/api'

export const connectionApi = {
  // 獲取連接列表
  async getConnections(): Promise<ApiResponse<Connection[]>> {
    return apiClient.get<Connection[]>('/connections')
  },

  // 獲取單個連接
  async getConnection(id: number): Promise<ApiResponse<Connection>> {
    return apiClient.get<Connection>(`/connections/${id}`)
  },

  // 創建連接
  async createConnection(data: CreateConnectionRequest): Promise<ApiResponse<Connection>> {
    return apiClient.post<Connection>('/connections', data)
  },

  // 更新連接
  async updateConnection(id: number, data: UpdateConnectionRequest): Promise<ApiResponse<Connection>> {
    return apiClient.put<Connection>(`/connections/${id}`, data)
  },

  // 刪除連接
  async deleteConnection(id: number): Promise<ApiResponse<void>> {
    return apiClient.delete<void>(`/connections/${id}`)
  },

  // 測試連接
  async testConnection(id: number): Promise<ApiResponse<TestConnectionResponse>> {
    return apiClient.post<TestConnectionResponse>(`/connections/${id}/test`)
  }
}
