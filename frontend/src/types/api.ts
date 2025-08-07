// API Response Types
export interface ApiResponse<T = any> {
  status: 'success' | 'error'
  message: string
  data?: T
  error?: string
}

// User Types
export interface User {
  id: number
  username: string
  email: string
  role: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

// Connection Types
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

export interface CreateConnectionRequest {
  name: string
  endpoints: string[]
  username?: string
  password?: string
  description?: string
  is_readonly?: boolean
}

export interface UpdateConnectionRequest {
  name: string
  endpoints: string[]
  username?: string
  password?: string
  description?: string
  is_readonly?: boolean
}

export interface TestConnectionResponse {
  status: 'success' | 'error'
  message: string
}

// Key-Value Types
export interface KVItem {
  key: string
  value: any
}

export interface KVListResponse {
  keys: string[]
}

export interface KVGetResponse {
  key: string
  value: any
}

export interface KVSetRequest {
  value: any
}

// Backup Types
export interface BackupData {
  connection_name: string
  connection_id: number
  export_time: string
  data: Record<string, any>
}

export interface ImportRequest {
  data: Record<string, any>
  overwrite: boolean
}

export interface ImportResponse {
  success_count: number
  error_count: number
  errors: string[]
}

// Transfer Types
export interface TransferRequest {
  source_connection_id: number
  target_connection_id: number
  keys?: string[]
  prefix?: string
  overwrite: boolean
  key_mapping?: boolean
  source_prefix?: string
  target_prefix?: string
}

export interface TransferResponse {
  success_count: number
  error_count: number
  skipped_count: number
  errors: string[]
  details: string[]
}
