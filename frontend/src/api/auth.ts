import apiClient from './client'
import type { 
  LoginRequest, 
  LoginResponse, 
  RegisterRequest, 
  User,
  ApiResponse 
} from '@/types/api'

export const authApi = {
  // 用戶登錄
  async login(credentials: LoginRequest): Promise<ApiResponse<LoginResponse>> {
    return apiClient.post<LoginResponse>('/auth/login', credentials)
  },

  // 用戶註冊
  async register(userData: RegisterRequest): Promise<ApiResponse<User>> {
    return apiClient.post<User>('/auth/register', userData)
  },

  // 獲取用戶信息
  async getProfile(): Promise<ApiResponse<User>> {
    return apiClient.get<User>('/auth/profile')
  },

  // 用戶登出
  async logout(): Promise<ApiResponse<void>> {
    return apiClient.post<void>('/auth/logout')
  }
}
