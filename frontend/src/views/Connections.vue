<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <h1 class="text-2xl font-bold text-gray-900">Connection Management</h1>
      <button
        @click="showCreateModal = true"
        class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
      >
        <svg class="-ml-1 mr-2 h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        New Connection
      </button>
    </div>

    <!-- Active Connections Tabs -->
    <div v-if="Array.from(activeConnections.keys()).length > 0" class="border-b border-gray-200">
      <nav class="-mb-px flex space-x-8">
        <button
          v-for="connectionId in Array.from(activeConnections.keys())"
          :key="connectionId"
          @click="setCurrentConnection(connectionId)"
          :class="[
            'whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm relative',
            currentConnectionId === connectionId
              ? 'border-indigo-500 text-indigo-600'
              : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
          ]"
        >
          <div class="flex items-center space-x-2">
            <span>{{ activeConnections.get(connectionId)?.name }}</span>
            
            <!-- Connection Status Indicator -->
            <div
              :class="[
                'w-2 h-2 rounded-full',
                {
                  'bg-green-500': activeConnections.get(connectionId)?.status === 'connected',
                  'bg-yellow-500': activeConnections.get(connectionId)?.status === 'connecting',
                  'bg-red-500': activeConnections.get(connectionId)?.status === 'error',
                  'bg-gray-400': activeConnections.get(connectionId)?.status === 'disconnected'
                }
              ]"
            ></div>
            
            <!-- Read-only Indicator -->
            <svg
              v-if="activeConnections.get(connectionId)?.is_readonly"
              class="w-4 h-4 text-amber-500"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
            </svg>
            
            <!-- Close Tab Button -->
            <button
              @click.stop="removeActiveConnection(connectionId)"
              class="ml-2 text-gray-400 hover:text-gray-600"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </button>
      </nav>
    </div>

    <!-- Current Connection Details -->
    <div v-if="getCurrentConnection()" class="bg-white shadow rounded-lg p-6">
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-lg font-medium text-gray-900">
          {{ getCurrentConnection()?.name }}
          <span v-if="getCurrentConnection()?.is_readonly" class="ml-2 inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-amber-100 text-amber-800">
            <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
            </svg>
            Read Only
          </span>
        </h2>
        <div class="flex space-x-2">
          <button
            @click="testCurrentConnection"
            :disabled="testingConnections.has(getCurrentConnection()?.id || 0)"
            class="inline-flex items-center px-3 py-1.5 border border-gray-300 shadow-sm text-xs font-medium rounded text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
          >
            <svg v-if="testingConnections.has(getCurrentConnection()?.id || 0)" class="animate-spin -ml-1 mr-1 h-3 w-3 text-gray-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ testingConnections.has(getCurrentConnection()?.id || 0) ? 'Testing...' : 'Test Connection' }}
          </button>
          
          <button
            @click="viewKeysForCurrentConnection"
            class="inline-flex items-center px-3 py-1.5 border border-transparent text-xs font-medium rounded text-indigo-700 bg-indigo-100 hover:bg-indigo-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          >
            View Keys
          </button>
          
          <button
            @click="removeActiveConnection(getCurrentConnection()?.id || 0)"
            class="inline-flex items-center px-3 py-1.5 border border-transparent text-xs font-medium rounded text-red-700 bg-red-100 hover:bg-red-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
          >
            Disconnect
          </button>
        </div>
      </div>
      
      <div class="grid grid-cols-2 gap-4 text-sm">
        <div>
          <span class="font-medium text-gray-500">Endpoints:</span>
          <p class="mt-1">{{ parseEndpoints(getCurrentConnection()?.endpoints || '').join(', ') }}</p>
        </div>
        <div>
          <span class="font-medium text-gray-500">Status:</span>
          <span :class="[
            'mt-1 inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
            {
              'bg-green-100 text-green-800': getCurrentConnection()?.status === 'connected',
              'bg-yellow-100 text-yellow-800': getCurrentConnection()?.status === 'connecting',
              'bg-red-100 text-red-800': getCurrentConnection()?.status === 'error',
              'bg-gray-100 text-gray-800': getCurrentConnection()?.status === 'disconnected'
            }
          ]">
            {{ getCurrentConnection()?.status || 'unknown' }}
          </span>
        </div>
        <div v-if="getCurrentConnection()?.description" class="col-span-2">
          <span class="font-medium text-gray-500">Description:</span>
          <p class="mt-1">{{ getCurrentConnection()?.description }}</p>
        </div>
        <div v-if="getCurrentConnection()?.error" class="col-span-2">
          <span class="font-medium text-red-500">Error:</span>
          <p class="mt-1 text-red-600">{{ getCurrentConnection()?.error }}</p>
        </div>
      </div>
    </div>

    <!-- Available Connections -->
    <div class="bg-white shadow rounded-lg">
      <div class="px-4 py-5 sm:p-6">
        <h3 class="text-lg leading-6 font-medium text-gray-900">Available Connections</h3>
        <div class="mt-2 max-w-xl text-sm text-gray-500">
          <p>Connect to an etcd cluster to start managing key-value pairs.</p>
        </div>

        <!-- Loading State -->
        <div v-if="isLoading && connections.length === 0" class="text-center py-8">
          <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
          <p class="mt-2 text-gray-500">Loading connections...</p>
        </div>

        <!-- Error State -->
        <div v-if="error" class="mt-4 rounded-md bg-red-50 p-4">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
              </svg>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">{{ error }}</h3>
            </div>
          </div>
        </div>

        <!-- Connections List -->
        <div v-if="connections.length > 0" class="mt-5 space-y-3">
          <div
            v-for="connection in connections"
            :key="connection.id"
            class="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50"
          >
            <div class="flex-1">
              <div class="flex items-center">
                <h4 class="text-sm font-medium text-gray-900">
                  {{ connection.name }}
                </h4>
                <div class="ml-2 flex items-center space-x-2">
                  <span
                    :class="[
                      'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
                      connection.is_active 
                        ? 'bg-green-100 text-green-800' 
                        : 'bg-gray-100 text-gray-800'
                    ]"
                  >
                    {{ connection.is_active ? 'Active' : 'Inactive' }}
                  </span>
                  <span
                    v-if="connection.is_readonly"
                    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-amber-100 text-amber-800"
                  >
                    <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                    </svg>
                    Read Only
                  </span>
                </div>
              </div>
              <p class="mt-1 text-sm text-gray-500">
                {{ parseEndpoints(connection.endpoints).join(', ') }}
              </p>
              <p v-if="connection.description" class="mt-1 text-sm text-gray-500">
                {{ connection.description }}
              </p>
            </div>
            
            <div class="flex items-center space-x-2">
              <button
                v-if="!activeConnections.has(connection.id)"
                @click="connectToConnection(connection)"
                class="inline-flex items-center px-3 py-1.5 border border-transparent text-xs font-medium rounded text-indigo-700 bg-indigo-100 hover:bg-indigo-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                Connect
              </button>
              <span v-else class="text-xs text-green-600 font-medium">Connected</span>
              
              <button
                @click="editConnection(connection)"
                class="inline-flex items-center px-3 py-1.5 border border-gray-300 shadow-sm text-xs font-medium rounded text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                Edit
              </button>
              
              <button
                @click="deleteConnection(connection.id)"
                class="inline-flex items-center px-3 py-1.5 border border-transparent text-xs font-medium rounded text-red-700 bg-red-100 hover:bg-red-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
              >
                Delete
              </button>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <div v-if="!isLoading && connections.length === 0 && !error" class="text-center py-8">
          <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
          <h3 class="mt-2 text-sm font-medium text-gray-900">No connections</h3>
          <p class="mt-1 text-sm text-gray-500">Get started by creating a new etcd connection.</p>
          <div class="mt-6">
            <button
              @click="showCreateModal = true"
              class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              <svg class="-ml-1 mr-2 h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
              </svg>
              New Connection
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Connection Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
        <div class="mt-3">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Create New Connection</h3>
          
          <form @submit.prevent="createConnection" class="space-y-4">
            <div>
              <label for="name" class="block text-sm font-medium text-gray-700">Name</label>
              <input
                id="name"
                v-model="newConnection.name"
                type="text"
                required
                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                placeholder="My etcd cluster"
              />
            </div>
            
            <div>
              <label for="endpoints" class="block text-sm font-medium text-gray-700">Endpoints</label>
              <input
                id="endpoints"
                v-model="endpointsInput"
                type="text"
                required
                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                placeholder="localhost:2379,localhost:2380"
              />
              <p class="mt-1 text-xs text-gray-500">Comma-separated list of endpoints</p>
            </div>
            
            <div>
              <label for="username" class="block text-sm font-medium text-gray-700">Username (Optional)</label>
              <input
                id="username"
                v-model="newConnection.username"
                type="text"
                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                placeholder="admin"
              />
            </div>
            
            <div>
              <label for="password" class="block text-sm font-medium text-gray-700">Password (Optional)</label>
              <input
                id="password"
                v-model="newConnection.password"
                type="password"
                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              />
            </div>
            
            <div>
              <label for="description" class="block text-sm font-medium text-gray-700">Description (Optional)</label>
              <textarea
                id="description"
                v-model="newConnection.description"
                rows="2"
                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                placeholder="Production etcd cluster"
              ></textarea>
            </div>

            <div class="flex items-center">
              <input
                id="readonly"
                v-model="newConnection.is_readonly"
                type="checkbox"
                class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
              />
              <label for="readonly" class="ml-2 block text-sm text-gray-900">
                Read-only connection
              </label>
            </div>
            <p class="text-xs text-gray-500 -mt-2">
              Read-only connections cannot modify or delete keys
            </p>

            <div class="flex justify-end space-x-3 pt-4">
              <button
                type="button"
                @click="cancelCreate"
                class="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="isCreating"
                class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
              >
                {{ isCreating ? 'Creating...' : 'Create' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Edit Connection Modal -->
    <div v-if="showEditModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
        <div class="mt-3">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Edit Connection</h3>
          
          <form @submit.prevent="updateConnection" class="space-y-4">
            <div>
              <label for="edit-name" class="block text-sm font-medium text-gray-700">Name</label>
              <input
                id="edit-name"
                v-model="editingConnection.name"
                type="text"
                required
                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                placeholder="My etcd cluster"
              />
            </div>
            
            <div>
              <label for="edit-endpoints" class="block text-sm font-medium text-gray-700">Endpoints</label>
              <input
                id="edit-endpoints"
                v-model="editEndpointsInput"
                type="text"
                required
                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                placeholder="localhost:2379,localhost:2380"
              />
              <p class="mt-1 text-xs text-gray-500">Comma-separated list of endpoints</p>
            </div>
            
            <div>
              <label for="edit-username" class="block text-sm font-medium text-gray-700">Username (Optional)</label>
              <input
                id="edit-username"
                v-model="editingConnection.username"
                type="text"
                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                placeholder="admin"
              />
            </div>
            
            <div>
              <label for="edit-password" class="block text-sm font-medium text-gray-700">Password (Optional)</label>
              <input
                id="edit-password"
                v-model="editingConnection.password"
                type="password"
                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              />
            </div>
            
            <div>
              <label for="edit-description" class="block text-sm font-medium text-gray-700">Description (Optional)</label>
              <textarea
                id="edit-description"
                v-model="editingConnection.description"
                rows="2"
                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                placeholder="Production etcd cluster"
              ></textarea>
            </div>

            <div class="flex items-center">
              <input
                id="edit-readonly"
                v-model="editingConnection.is_readonly"
                type="checkbox"
                class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
              />
              <label for="edit-readonly" class="ml-2 block text-sm text-gray-900">
                Read-only connection
              </label>
            </div>
            <p class="text-xs text-gray-500 -mt-2">
              Read-only connections cannot modify or delete keys
            </p>

            <div class="flex justify-end space-x-3 pt-4">
              <button
                type="button"
                @click="cancelEdit"
                class="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="isUpdating"
                class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
              >
                {{ isUpdating ? 'Updating...' : 'Update' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { connectionApi } from '@/api/connection'
import type { Connection, CreateConnectionRequest, UpdateConnectionRequest } from '@/types/api'
import { 
  activeConnections, 
  currentConnectionId, 
  addActiveConnection, 
  removeActiveConnection, 
  setCurrentConnection, 
  updateConnectionStatus, 
  getCurrentConnection 
} from '@/composables/useConnections'

const router = useRouter()

const connections = ref<Connection[]>([])
const isLoading = ref(false)
const error = ref<string | null>(null)
const showCreateModal = ref(false)
const showEditModal = ref(false)
const isCreating = ref(false)
const isUpdating = ref(false)
const testingConnections = ref(new Set<number>())

const newConnection = reactive<CreateConnectionRequest>({
  name: '',
  endpoints: [],
  username: '',
  password: '',
  description: '',
  is_readonly: false
})

const editingConnection = reactive<UpdateConnectionRequest>({
  name: '',
  endpoints: [],
  username: '',
  password: '',
  description: '',
  is_readonly: false
})

const endpointsInput = ref('')
const editEndpointsInput = ref('')
let currentEditingId = 0

const parseEndpoints = (endpointsStr: string): string[] => {
  try {
    return JSON.parse(endpointsStr)
  } catch {
    return []
  }
}

const loadConnections = async () => {
  isLoading.value = true
  error.value = null

  try {
    const response = await connectionApi.getConnections()
    if (response.status === 'success' && response.data) {
      connections.value = response.data
    } else {
      error.value = response.message || 'Failed to load connections'
    }
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Failed to load connections'
  } finally {
    isLoading.value = false
  }
}

const createConnection = async () => {
  isCreating.value = true
  
  try {
    // Parse endpoints from input
    const endpoints = endpointsInput.value.split(',').map(e => e.trim()).filter(e => e)
    
    const connectionData: CreateConnectionRequest = {
      ...newConnection,
      endpoints
    }

    const response = await connectionApi.createConnection(connectionData)
    
    if (response.status === 'success') {
      showCreateModal.value = false
      await loadConnections()
      resetForm()
    } else {
      error.value = response.message || 'Failed to create connection'
    }
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Failed to create connection'
  } finally {
    isCreating.value = false
  }
}

const connectToConnection = async (connection: Connection) => {
  // Add to active connections
  addActiveConnection(connection)
  updateConnectionStatus(connection.id, 'connecting')
  
  try {
    // Test the connection
    const response = await connectionApi.testConnection(connection.id)
    
    if (response.status === 'success') {
      updateConnectionStatus(connection.id, 'connected')
    } else {
      updateConnectionStatus(connection.id, 'error', response.message)
    }
  } catch (err: any) {
    updateConnectionStatus(connection.id, 'error', err.response?.data?.message || 'Connection failed')
  }
}

const testCurrentConnection = async () => {
  const current = getCurrentConnection()
  if (!current) return
  
  testingConnections.value.add(current.id)
  updateConnectionStatus(current.id, 'connecting')
  
  try {
    const response = await connectionApi.testConnection(current.id)
    
    if (response.status === 'success') {
      updateConnectionStatus(current.id, 'connected')
      alert('Connection test successful!')
    } else {
      updateConnectionStatus(current.id, 'error', response.message)
      alert(`Connection test failed: ${response.message}`)
    }
  } catch (err: any) {
    const errorMsg = err.response?.data?.message || 'Unknown error'
    updateConnectionStatus(current.id, 'error', errorMsg)
    alert(`Connection test failed: ${errorMsg}`)
  } finally {
    testingConnections.value.delete(current.id)
  }
}

const viewKeysForCurrentConnection = () => {
  const current = getCurrentConnection()
  if (current) {
    // 确保连接在活跃连接列表中
    if (!activeConnections.value.has(current.id)) {
      addActiveConnection(current)
    }
    // 设置为当前连接
    setCurrentConnection(current.id)
    // 跳转到KV管理页面
    router.push(`/kv/${current.id}`)
  }
}

const editConnection = (connection: Connection) => {
  // 設置編輯連線的資料
  currentEditingId = connection.id
  editingConnection.name = connection.name
  editingConnection.endpoints = parseEndpoints(connection.endpoints)
  editingConnection.username = connection.username || ''
  editingConnection.password = connection.password || ''
  editingConnection.description = connection.description || ''
  editingConnection.is_readonly = connection.is_readonly || false
  
  // 設置端點輸入框的值
  editEndpointsInput.value = parseEndpoints(connection.endpoints).join(', ')
  
  // 顯示編輯模態框
  showEditModal.value = true
}

const updateConnection = async () => {
  isUpdating.value = true
  
  try {
    // Parse endpoints from input
    const endpoints = editEndpointsInput.value.split(',').map(e => e.trim()).filter(e => e)
    
    const connectionData: UpdateConnectionRequest = {
      ...editingConnection,
      endpoints
    }

    const response = await connectionApi.updateConnection(currentEditingId, connectionData)
    
    if (response.status === 'success') {
      showEditModal.value = false
      await loadConnections()
      
      // 如果更新的連線是活躍連線，更新其資料
      if (activeConnections.value.has(currentEditingId)) {
        const updatedConnection = response.data
        if (updatedConnection) {
          addActiveConnection(updatedConnection)
        }
      }
      
      resetEditForm()
    } else {
      error.value = response.message || 'Failed to update connection'
    }
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Failed to update connection'
  } finally {
    isUpdating.value = false
  }
}

const cancelEdit = () => {
  showEditModal.value = false
  resetEditForm()
}

const resetEditForm = () => {
  Object.assign(editingConnection, {
    name: '',
    endpoints: [],
    username: '',
    password: '',
    description: '',
    is_readonly: false
  })
  editEndpointsInput.value = ''
  currentEditingId = 0
}

const deleteConnection = async (id: number) => {
  if (!confirm('Are you sure you want to delete this connection?')) {
    return
  }

  try {
    const response = await connectionApi.deleteConnection(id)
    
    if (response.status === 'success') {
      // Remove from active connections if it exists
      removeActiveConnection(id)
      await loadConnections()
    } else {
      error.value = response.message || 'Failed to delete connection'
    }
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Failed to delete connection'
  }
}

const cancelCreate = () => {
  showCreateModal.value = false
  resetForm()
}

const resetForm = () => {
  Object.assign(newConnection, {
    name: '',
    endpoints: [],
    username: '',
    password: '',
    description: '',
    is_readonly: false
  })
  endpointsInput.value = ''
}

onMounted(() => {
  loadConnections()
})
</script>
