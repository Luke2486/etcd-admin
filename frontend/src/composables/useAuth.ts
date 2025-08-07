import { ref, computed } from 'vue'
import { authApi } from '@/api/auth'
import type { User, LoginRequest, RegisterRequest } from '@/types/api'

// Global state
const user = ref<User | null>(null)
const token = ref<string | null>(null)
const isLoading = ref(false)
const error = ref<string | null>(null)

// Initialize from localStorage
const initAuth = () => {
  const storedToken = localStorage.getItem('auth_token')
  const storedUser = localStorage.getItem('user')
  
  if (storedToken) {
    token.value = storedToken
  }
  
  if (storedUser) {
    try {
      user.value = JSON.parse(storedUser)
    } catch (e) {
      console.error('Failed to parse stored user data:', e)
      localStorage.removeItem('user')
    }
  }
}

export const useAuth = () => {
  const isAuthenticated = computed(() => !!token.value && !!user.value)

  const login = async (credentials: LoginRequest) => {
    isLoading.value = true
    error.value = null

    try {
      console.log('Attempting login with:', credentials)
      const response = await authApi.login(credentials)
      console.log('Login response:', response)
      
      if (response.status === 'success' && response.data) {
        token.value = response.data.token
        user.value = response.data.user
        
        // Store in localStorage
        localStorage.setItem('auth_token', response.data.token)
        localStorage.setItem('user', JSON.stringify(response.data.user))
        
        return { success: true }
      } else {
        error.value = response.message || 'Login failed'
        console.error('Login failed:', response.message)
        return { success: false, error: error.value }
      }
    } catch (err: any) {
      console.error('Login error:', err)
      error.value = err.response?.data?.message || 'Login failed'
      return { success: false, error: error.value }
    } finally {
      isLoading.value = false
    }
  }

  const register = async (userData: RegisterRequest) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await authApi.register(userData)
      
      if (response.status === 'success') {
        return { success: true }
      } else {
        error.value = response.message || 'Registration failed'
        return { success: false, error: error.value }
      }
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Registration failed'
      return { success: false, error: error.value }
    } finally {
      isLoading.value = false
    }
  }

  const logout = async () => {
    isLoading.value = true
    
    try {
      await authApi.logout()
    } catch (err) {
      console.error('Logout API call failed:', err)
    } finally {
      // Clear local state regardless of API call result
      token.value = null
      user.value = null
      localStorage.removeItem('auth_token')
      localStorage.removeItem('user')
      isLoading.value = false
    }
  }

  const fetchProfile = async () => {
    if (!token.value) return

    isLoading.value = true
    error.value = null

    try {
      const response = await authApi.getProfile()
      
      if (response.status === 'success' && response.data) {
        user.value = response.data
        localStorage.setItem('user', JSON.stringify(response.data))
      }
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Failed to fetch profile'
      // If token is invalid, logout
      if (err.response?.status === 401) {
        await logout()
      }
    } finally {
      isLoading.value = false
    }
  }

  return {
    user: computed(() => user.value),
    token: computed(() => token.value),
    isAuthenticated,
    isLoading: computed(() => isLoading.value),
    error: computed(() => error.value),
    login,
    register,
    logout,
    fetchProfile,
    initAuth,
    clearError: () => { error.value = null }
  }
}
