<template>
  <div class="border-b border-gray-200 bg-white">
    <div class="flex items-center justify-between px-4 py-3">
      <div class="flex items-center space-x-2">
        <button
          @click="$router.push('/connections')"
          class="inline-flex items-center px-3 py-1.5 border border-gray-300 text-xs font-medium rounded text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        >
          <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
          Back to Connections
        </button>
        <span class="text-gray-400">|</span>
        <h2 class="text-sm font-medium text-gray-600">Active Connections</h2>
      </div>
      <div class="text-xs text-gray-500">
        {{ activeConnections.size }} connection{{ activeConnections.size !== 1 ? 's' : '' }} active
      </div>
    </div>
    
    <div v-if="activeConnections.size > 0" class="px-4">
      <nav class="-mb-px flex space-x-8 overflow-x-auto">
        <button
          v-for="[id, connection] in activeConnections"
          :key="id"
          @click="switchConnection(id)"
          :class="[
            'whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm transition-colors duration-200',
            currentConnectionId === id
              ? 'border-indigo-500 text-indigo-600'
              : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
          ]"
        >
          <div class="flex items-center space-x-2">
            <div
              :class="[
                'w-2 h-2 rounded-full',
                connection.status === 'connected' ? 'bg-green-400' : 'bg-red-400'
              ]"
            ></div>
            <span>{{ connection.name }}</span>
            <span
              v-if="connection.is_readonly"
              class="inline-flex items-center px-1.5 py-0.5 rounded-full text-xs font-medium bg-amber-100 text-amber-700"
            >
              <svg class="w-2 h-2 mr-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
              </svg>
              RO
            </span>
            <button
              @click.stop="closeConnection(id)"
              :class="[
                'ml-1 w-4 h-4 rounded-full flex items-center justify-center transition-colors duration-200',
                currentConnectionId === id
                  ? 'hover:bg-indigo-100 text-indigo-600'
                  : 'hover:bg-gray-100 text-gray-400'
              ]"
            >
              <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </button>
      </nav>
    </div>
    
    <div v-else class="px-4 py-8 text-center">
      <svg class="mx-auto h-8 w-8 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">No active connections</h3>
      <p class="mt-1 text-sm text-gray-500">Connect to an etcd instance to start managing keys.</p>
      <div class="mt-4">
        <button
          @click="$router.push('/connections')"
          class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        >
          Go to Connections
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { activeConnections, removeActiveConnection, setCurrentConnection } from '@/composables/useConnections'

interface Props {
  currentConnectionId?: number
}

const props = defineProps<Props>()
const router = useRouter()

const closeConnection = (connectionId: number) => {
  removeActiveConnection(connectionId)
  // 如果没有活跃连接了，回到连接列表
  if (activeConnections.value.size === 0) {
    router.push('/connections')
  }
}

const switchConnection = (connectionId: number) => {
  if (connectionId !== props.currentConnectionId) {
    // 更新当前连接状态
    setCurrentConnection(connectionId)
    // 跳转到新的连接页面
    router.push(`/kv/${connectionId}`)
  }
}
</script>
