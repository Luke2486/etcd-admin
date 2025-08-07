import { ref } from 'vue'

export interface Connection {
  id: number
  name: string
  endpoints: string
  username?: string
  password?: string
  tls_enabled: boolean
  cert_file?: string
  key_file?: string
  ca_file?: string
  description?: string
  is_active: boolean
  is_readonly: boolean
  created_at: string
  updated_at: string
}

export interface ActiveConnection extends Connection {
  status?: 'connected' | 'connecting' | 'disconnected' | 'error'
  error?: string
}

// 活躍連線管理
export const activeConnections = ref<Map<number, ActiveConnection>>(new Map())
export const currentConnectionId = ref<number | null>(null)

// 添加活躍連線
export function addActiveConnection(connection: Connection) {
  activeConnections.value.set(connection.id, {
    ...connection,
    status: 'connecting'
  })
  
  // 如果是第一個連線，設為當前連線
  if (currentConnectionId.value === null) {
    currentConnectionId.value = connection.id
  }
}

// 移除活躍連線
export function removeActiveConnection(connectionId: number) {
  activeConnections.value.delete(connectionId)
  
  // 如果移除的是當前連線，切換到第一個可用連線
  if (currentConnectionId.value === connectionId) {
    const firstConnection = Array.from(activeConnections.value.keys())[0]
    currentConnectionId.value = firstConnection || null
  }
}

// 設置當前連線
export function setCurrentConnection(connectionId: number) {
  if (activeConnections.value.has(connectionId)) {
    currentConnectionId.value = connectionId
  }
}

// 更新連線狀態
export function updateConnectionStatus(connectionId: number, status: ActiveConnection['status'], error?: string) {
  const connection = activeConnections.value.get(connectionId)
  if (connection) {
    connection.status = status
    if (error) {
      connection.error = error
    } else {
      delete connection.error
    }
  }
}

// 獲取當前連線
export function getCurrentConnection(): ActiveConnection | null {
  if (currentConnectionId.value === null) return null
  return activeConnections.value.get(currentConnectionId.value) || null
}

// 清除所有活躍連線
export function clearActiveConnections() {
  activeConnections.value.clear()
  currentConnectionId.value = null
}
