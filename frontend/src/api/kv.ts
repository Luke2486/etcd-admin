import apiClient from './client'
import type { 
  KVListResponse,
  KVGetResponse,
  KVSetRequest,
  ApiResponse 
} from '@/types/api'

export const kvApi = {
  // 獲取所有鍵
  async getKeys(connectionId: number, prefix?: string): Promise<ApiResponse<KVListResponse>> {
    const params = prefix ? { prefix } : {}
    return apiClient.get<KVListResponse>(`/connections/${connectionId}/kv`, params)
  },

  // 獲取鍵值
  async getValue(connectionId: number, key: string): Promise<ApiResponse<KVGetResponse>> {
    // URL encode the key to handle special characters
    const encodedKey = encodeURIComponent(key)
    return apiClient.get<KVGetResponse>(`/connections/${connectionId}/kv/${encodedKey}`)
  },

  // 設置鍵值
  async setValue(connectionId: number, key: string, data: KVSetRequest): Promise<ApiResponse<KVGetResponse>> {
    const encodedKey = encodeURIComponent(key)
    return apiClient.put<KVGetResponse>(`/connections/${connectionId}/kv/${encodedKey}`, data)
  },

  // 刪除鍵
  async deleteKey(connectionId: number, key: string): Promise<ApiResponse<{ key: string }>> {
    const encodedKey = encodeURIComponent(key)
    return apiClient.delete<{ key: string }>(`/connections/${connectionId}/kv/${encodedKey}`)
  }
}
