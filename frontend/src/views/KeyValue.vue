<template>
  <div>
    <!-- Connection Tabs -->
    <ConnectionTabs :current-connection-id="connectionId" />
    
    <!-- Main Content -->
    <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <div class="px-4 py-6 sm:px-0">
        <div class="space-y-6">
          <!-- Header with Current Connection Details -->
          <div class="bg-white shadow sm:rounded-lg">
            <div class="px-4 py-5 sm:p-6">
              <div class="flex justify-between items-start">
                <div>
                  <h1 class="text-2xl font-bold text-gray-900">Key-Value Management</h1>
                  <div v-if="connectionName" class="flex items-center mt-2 space-x-3">
                    <div class="flex items-center space-x-2">
                      <div
                        :class="[
                          'w-3 h-3 rounded-full',
                          'bg-green-400'
                        ]"
                      ></div>
                      <p class="text-lg font-medium text-gray-700">
                        {{ connectionName }}
                      </p>
                    </div>
                    <span
                      v-if="connectionReadOnly"
                      class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-amber-100 text-amber-800"
                    >
                      <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                      </svg>
                      Read Only Connection
                    </span>
                  </div>
                  <div class="mt-1 text-sm text-gray-500">
                    Manage keys and values for this etcd connection
                  </div>
                </div>
                <div class="flex items-center space-x-3">
                  <button
                    @click="loadKeys"
                    :disabled="isLoading"
                    class="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
                  >
                    <svg v-if="isLoading" class="animate-spin -ml-1 mr-2 h-4 w-4 text-gray-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    <svg v-else class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                    </svg>
                    {{ isLoading ? 'Loading...' : 'Refresh' }}
                  </button>
                  <button
                    v-if="!connectionReadOnly"
                    @click="showCreateModal = true"
                    class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                  >
                    <svg class="-ml-1 mr-2 h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                    </svg>
                    Add Key
                  </button>
                  <div v-else class="text-sm text-amber-600 bg-amber-50 px-3 py-2 rounded-md">
                    <svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                    </svg>
                    Read-only mode: Cannot add keys
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Search and Filter -->
          <div class="bg-white shadow sm:rounded-lg">
            <div class="px-4 py-5 sm:p-6">
              <div class="flex space-x-4">
                <div class="flex-1">
                  <label for="prefix" class="sr-only">Prefix filter</label>
                  <div class="relative">
                    <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                      <svg class="h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                      </svg>
                    </div>
                    <input
                      id="prefix"
                      v-model="prefixFilter"
                      type="text"
                      class="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md leading-5 bg-white placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:ring-1 focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                      placeholder="Filter by prefix (e.g., /app/, config...)..."
                      @keyup.enter="loadKeys"
                    />
                  </div>
                </div>
                <button
                  @click="loadKeys"
                  :disabled="isLoading"
                  class="inline-flex items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
                >
                  <svg v-if="isLoading" class="animate-spin -ml-1 mr-2 h-4 w-4 text-gray-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  {{ isLoading ? 'Filtering...' : 'Apply Filter' }}
                </button>
              </div>
            </div>
          </div>

          <!-- Error State -->
          <div v-if="error" class="rounded-md bg-red-50 p-4">
            <div class="flex">
              <div class="flex-shrink-0">
                <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
                </svg>
              </div>
              <div class="ml-3">
                <h3 class="text-sm font-medium text-red-800">{{ error }}</h3>
              </div>
              <div class="ml-auto pl-3">
                <button
                  @click="error = null"
                  class="inline-flex rounded-md bg-red-50 p-1.5 text-red-500 hover:bg-red-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-red-50 focus:ring-red-600"
                >
                  <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                  </svg>
                </button>
              </div>
            </div>
          </div>

          <!-- Keys List -->
          <div v-if="keys.length > 0" class="bg-white shadow overflow-hidden sm:rounded-lg">
            <div class="px-4 py-5 sm:px-6 border-b border-gray-200">
              <div class="flex items-center justify-between">
                <div>
                  <h3 class="text-lg leading-6 font-medium text-gray-900">Keys</h3>
                  <p class="mt-1 max-w-2xl text-sm text-gray-500">
                    {{ keys.length }} key{{ keys.length !== 1 ? 's' : '' }} found
                    {{ prefixFilter ? ` matching "${prefixFilter}"` : '' }}
                  </p>
                </div>
                <div class="text-sm text-gray-500">
                  {{ connectionReadOnly ? 'Read-only view' : 'Manage keys and values' }}
                </div>
              </div>
            </div>
            <ul class="divide-y divide-gray-200">
              <li v-for="key in keys" :key="key" class="px-6 py-4 hover:bg-gray-50 transition-colors duration-150">
                <div class="flex items-center justify-between">
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center space-x-3">
                      <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m0 0a2 2 0 012 2v6a2 2 0 01-2 2h-6a2 2 0 01-2-2V9a2 2 0 012-2m0 0V7a2 2 0 012-2m-6 0h6M9 12h6m-6 4h6" />
                      </svg>
                      <p class="text-sm font-medium text-gray-900 truncate font-mono">
                        {{ key }}
                      </p>
                    </div>
                  </div>
                  <div class="flex items-center space-x-2 ml-4">
                    <button
                      @click="viewValue(key)"
                      class="inline-flex items-center px-3 py-1.5 border border-gray-300 shadow-sm text-xs font-medium rounded text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors duration-150"
                    >
                      <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                      </svg>
                      View
                    </button>
                    <button
                      v-if="!connectionReadOnly"
                      @click="editValue(key)"
                      class="inline-flex items-center px-3 py-1.5 border border-transparent text-xs font-medium rounded text-indigo-700 bg-indigo-100 hover:bg-indigo-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors duration-150"
                    >
                      <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                      </svg>
                      Edit
                    </button>
                    <button
                      v-if="!connectionReadOnly"
                      @click="deleteKey(key)"
                      class="inline-flex items-center px-3 py-1.5 border border-transparent text-xs font-medium rounded text-red-700 bg-red-100 hover:bg-red-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 transition-colors duration-150"
                    >
                      <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                      </svg>
                      Delete
                    </button>
                    <span
                      v-if="connectionReadOnly"
                      class="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium bg-amber-100 text-amber-700"
                    >
                      <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                      </svg>
                      Read-only
                    </span>
                  </div>
                </div>
              </li>
            </ul>
          </div>

          <!-- Empty State -->
          <div v-if="!isLoading && keys.length === 0 && !error" class="bg-white shadow sm:rounded-lg">
            <div class="text-center py-12">
              <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
              </svg>
              <h3 class="mt-2 text-sm font-medium text-gray-900">No keys found</h3>
              <p class="mt-1 text-sm text-gray-500">
                {{ prefixFilter ? 
                  `No keys match the filter "${prefixFilter}". Try a different prefix or clear the filter.` : 
                  'This etcd connection has no keys yet.' 
                }}
              </p>
              <div v-if="!connectionReadOnly" class="mt-6">
                <button
                  @click="showCreateModal = true"
                  class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                >
                  <svg class="-ml-1 mr-2 h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                  </svg>
                  {{ prefixFilter ? 'Add Key with this Prefix' : 'Add First Key' }}
                </button>
              </div>
              <div v-else class="mt-6">
                <div class="inline-flex items-center px-4 py-2 bg-amber-50 border border-amber-200 rounded-md">
                  <svg class="w-4 h-4 mr-2 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                  </svg>
                  <span class="text-sm text-amber-700">Read-only connection: Cannot add keys</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create/Edit Key Modal -->
    <div v-if="showCreateModal || showEditModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div class="relative top-10 mx-auto p-6 border w-4/5 max-w-6xl shadow-lg rounded-md bg-white min-h-[700px]">
        <div class="mt-3">
          <h3 class="text-lg font-medium text-gray-900 mb-6">
            {{ showCreateModal ? 'Add New Key' : 'Edit Key' }}
          </h3>
          
          <form @submit.prevent="saveKey" class="space-y-6">
            <div>
              <label for="key" class="block text-sm font-medium text-gray-700 mb-2">Key</label>
              <input
                id="key"
                v-model="currentKey"
                type="text"
                :disabled="showEditModal"
                required
                class="block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm disabled:bg-gray-100"
                placeholder="my-key"
              />
            </div>
            
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <!-- Editor -->
              <div>
                <div class="flex justify-between items-center mb-2">
                  <label for="value" class="block text-sm font-medium text-gray-700">Value Editor</label>
                  <div class="flex space-x-2">
                    <button
                      type="button"
                      @click="formatCurrentValue"
                      :disabled="!canFormatCurrentValue"
                      class="inline-flex items-center px-3 py-1.5 text-xs font-medium rounded border transition-colors duration-150 disabled:opacity-50 disabled:cursor-not-allowed"
                      :class="canFormatCurrentValue 
                        ? 'bg-purple-100 text-purple-700 border-purple-200 hover:bg-purple-200' 
                        : 'bg-gray-100 text-gray-400 border-gray-200'"
                    >
                      <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
                      </svg>
                      Format JSON
                    </button>
                  </div>
                </div>
                <textarea
                  id="value"
                  v-model="currentValue"
                  rows="20"
                  required
                  class="block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm font-mono"
                  placeholder="Enter JSON, string, or any value..."
                ></textarea>
                <p class="mt-1 text-xs text-gray-500">
                  Enter any value. Use "Format JSON" button to beautify JSON syntax. JSON objects will be automatically formatted and previewed.
                </p>
              </div>
              
              <!-- Preview -->
              <div>
                <div class="flex justify-between items-center mb-2">
                  <label class="block text-sm font-medium text-gray-700">Live Preview</label>
                  <div class="flex items-center space-x-3">
                    <div class="text-xs text-gray-500" v-if="previewSelectedPath">
                      Selected: {{ previewSelectedPath }}
                    </div>
                    <button
                      v-if="previewSelectedNodeData"
                      @click="copySelectedData(true)"
                      :class="[
                        'inline-flex items-center px-2 py-1 text-xs font-medium rounded transition-colors duration-150',
                        copySuccess 
                          ? 'bg-green-100 text-green-600 border border-green-200' 
                          : 'bg-blue-100 text-blue-600 border border-blue-200 hover:bg-blue-200'
                      ]"
                    >
                      <svg v-if="copySuccess" class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                      </svg>
                      <svg v-else class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 002 2v8a2 2 0 002 2z" />
                      </svg>
                      {{ copySuccess ? 'Copied!' : 'Copy' }}
                    </button>
                  </div>
                </div>
                <div class="border rounded-md bg-white h-[480px] overflow-auto">
                  <vue-json-pretty
                    v-if="isPreviewJson && parsedCurrentValue"
                    :data="parsedCurrentValue"
                    :deep="3"
                    :show-line="true"
                    :show-line-number="true"
                    :show-icon="true"
                    :collapsed-on-click-brackets="true"
                    :select-on-click-node="true"
                    :selectable-type="'single'"
                    :highlight-selected-node="true"
                    v-model:selected-value="previewSelectedPath"
                    theme="light"
                    class="p-4"
                    @node-click="handlePreviewNodeClick"
                  />
                  <div v-else-if="currentValue" class="p-4">
                    <pre class="text-sm text-gray-900 whitespace-pre-wrap font-mono">{{ currentValue }}</pre>
                  </div>
                  <div v-else class="p-4 text-gray-400 text-sm italic">
                    Enter a value to see the preview
                  </div>
                </div>
              </div>
            </div>

            <div class="flex justify-end space-x-3 pt-6 border-t">
              <button
                type="button"
                @click="cancelEdit"
                class="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="isSaving"
                class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
              >
                {{ isSaving ? 'Saving...' : 'Save' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- View Value Modal -->
    <div v-if="showViewModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div class="relative top-10 mx-auto p-6 border w-4/5 max-w-6xl shadow-lg rounded-md bg-white min-h-[600px]">
        <div class="mt-3">
          <div class="flex justify-between items-center mb-6">
            <h3 class="text-lg font-medium text-gray-900">View Key Value</h3>
            <button
              @click="showViewModal = false"
              class="text-gray-400 hover:text-gray-600"
            >
              <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          
          <div class="space-y-6">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Key</label>
              <div class="p-3 bg-gray-50 rounded-md border">
                <code class="text-sm text-gray-900 break-all">{{ viewingKey }}</code>
              </div>
            </div>
            
            <div>
              <div class="flex justify-between items-center mb-2">
                <label class="block text-sm font-medium text-gray-700">Value</label>
                <div class="flex items-center space-x-3">
                  <div class="text-xs text-gray-500" v-if="selectedPath">
                    Selected: {{ selectedPath }}
                  </div>
                  <button
                    v-if="selectedNodeData"
                    @click="copySelectedData(false)"
                    :class="[
                      'inline-flex items-center px-3 py-1.5 text-xs font-medium rounded-md transition-colors duration-150',
                      copySuccess 
                        ? 'bg-green-100 text-green-700 border border-green-200' 
                        : 'bg-blue-100 text-blue-700 border border-blue-200 hover:bg-blue-200'
                    ]"
                  >
                    <svg v-if="copySuccess" class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                    </svg>
                    <svg v-else class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                    </svg>
                    {{ copySuccess ? 'Copied!' : 'Copy Selected' }}
                  </button>
                </div>
              </div>
              <div class="border rounded-md bg-white min-h-[400px]">
                <vue-json-pretty
                  v-if="isJsonValue"
                  :data="parsedViewingValue"
                  :deep="3"
                  :show-line="true"
                  :show-line-number="true"
                  :show-icon="true"
                  :collapsed-on-click-brackets="true"
                  :select-on-click-node="true"
                  :selectable-type="'single'"
                  :highlight-selected-node="true"
                  v-model:selected-value="selectedPath"
                  theme="light"
                  class="p-4"
                  @node-click="handleNodeClick"
                />
                <pre v-else class="p-4 text-sm text-gray-900 whitespace-pre-wrap overflow-x-auto font-mono">{{ viewingValue }}</pre>
              </div>
            </div>
          </div>

          <div class="flex justify-between items-center pt-6">
            <div class="text-sm text-gray-500">
              <span v-if="isJsonValue">JSON format detected - Click nodes to select</span>
              <span v-else>Plain text format</span>
            </div>
            <div class="flex space-x-3">
              <button
                v-if="!connectionReadOnly"
                @click="editValue(viewingKey)"
                class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                Edit
              </button>
              <button
                @click="showViewModal = false"
                class="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                Close
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { useRoute } from 'vue-router'
import { kvApi } from '@/api/kv'
import { connectionApi } from '@/api/connection'
import ConnectionTabs from '@/components/ConnectionTabs.vue'
import { addActiveConnection, updateConnectionStatus } from '@/composables/useConnections'
import VueJsonPretty from 'vue-json-pretty'
import 'vue-json-pretty/lib/styles.css'

const route = useRoute()

const connectionId = ref<number>(parseInt(route.params.id as string))
const connectionName = ref<string>('')
const connectionReadOnly = ref<boolean>(false)
const keys = ref<string[]>([])
const isLoading = ref(false)
const error = ref<string | null>(null)
const prefixFilter = ref('')

// Modal states
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showViewModal = ref(false)
const isSaving = ref(false)

// Current key/value being edited
const currentKey = ref('')
const currentValue = ref('')
const viewingKey = ref('')
const viewingValue = ref('')

// JSON selection states
const selectedPath = ref('')
const previewSelectedPath = ref('')
const selectedNodeData = ref<any>(null)
const previewSelectedNodeData = ref<any>(null)
const copySuccess = ref(false)

// Computed properties for JSON handling
const isJsonValue = computed(() => {
  try {
    JSON.parse(viewingValue.value)
    return true
  } catch {
    return false
  }
})

const parsedViewingValue = computed(() => {
  try {
    return JSON.parse(viewingValue.value)
  } catch {
    return null
  }
})

const isPreviewJson = computed(() => {
  try {
    JSON.parse(currentValue.value)
    return true
  } catch {
    return false
  }
})

const parsedCurrentValue = computed(() => {
  try {
    return JSON.parse(currentValue.value)
  } catch {
    return null
  }
})

const canFormatCurrentValue = computed(() => {
  try {
    const parsed = JSON.parse(currentValue.value)
    // 檢查是否為對象或數組，且不為null
    return (typeof parsed === 'object' && parsed !== null)
  } catch {
    return false
  }
})

const loadConnection = async () => {
  try {
    const response = await connectionApi.getConnection(connectionId.value)
    if (response.status === 'success' && response.data) {
      connectionName.value = response.data.name
      connectionReadOnly.value = response.data.is_readonly || false
      
      // 添加到活跃连接列表
      addActiveConnection(response.data)
      updateConnectionStatus(connectionId.value, 'connected')
    }
  } catch (err) {
    console.error('Failed to load connection details:', err)
    updateConnectionStatus(connectionId.value, 'error', 'Failed to load connection')
  }
}

const loadKeys = async () => {
  isLoading.value = true
  error.value = null

  try {
    const response = await kvApi.getKeys(
      connectionId.value, 
      prefixFilter.value || undefined
    )
    
    if (response.status === 'success' && response.data) {
      keys.value = response.data.keys || []
    } else {
      error.value = response.message || 'Failed to load keys'
    }
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Failed to load keys'
  } finally {
    isLoading.value = false
  }
}

const viewValue = async (key: string) => {
  try {
    const response = await kvApi.getValue(connectionId.value, key)
    
    if (response.status === 'success' && response.data) {
      viewingKey.value = key
      viewingValue.value = formatValue(response.data.value)
      showViewModal.value = true
    } else {
      error.value = response.message || 'Failed to load value'
    }
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Failed to load value'
  }
}

const editValue = async (key: string) => {
  try {
    const response = await kvApi.getValue(connectionId.value, key)
    
    if (response.status === 'success' && response.data) {
      currentKey.value = key
      currentValue.value = formatValue(response.data.value)
      showEditModal.value = true
      showViewModal.value = false
    } else {
      error.value = response.message || 'Failed to load value'
    }
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Failed to load value'
  }
}

const saveKey = async () => {
  isSaving.value = true
  
  try {
    let value: any = currentValue.value
    
    // Try to parse as JSON if it looks like JSON
    if (typeof value === 'string' && (value.trim().startsWith('{') || value.trim().startsWith('['))) {
      try {
        value = JSON.parse(value)
        // 自動格式化編輯器中的JSON內容
        currentValue.value = JSON.stringify(value, null, 2)
      } catch {
        // Keep as string if not valid JSON
      }
    }

    const response = await kvApi.setValue(connectionId.value, currentKey.value, { value })
    
    if (response.status === 'success') {
      showCreateModal.value = false
      showEditModal.value = false
      await loadKeys()
      resetCurrentKeyValue()
    } else {
      error.value = response.message || 'Failed to save key'
    }
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Failed to save key'
  } finally {
    isSaving.value = false
  }
}

const deleteKey = async (key: string) => {
  if (!confirm(`Are you sure you want to delete key "${key}"?`)) {
    return
  }

  try {
    const response = await kvApi.deleteKey(connectionId.value, key)
    
    if (response.status === 'success') {
      await loadKeys()
    } else {
      error.value = response.message || 'Failed to delete key'
    }
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Failed to delete key'
  }
}

const cancelEdit = () => {
  showCreateModal.value = false
  showEditModal.value = false
  resetCurrentKeyValue()
}

const resetCurrentKeyValue = () => {
  currentKey.value = ''
  currentValue.value = ''
  selectedPath.value = ''
  previewSelectedPath.value = ''
  selectedNodeData.value = null
  previewSelectedNodeData.value = null
  copySuccess.value = false
}

const getNodeDataAtPath = (data: any, path: string): any => {
  if (!path || path === 'root') return data
  
  const pathParts = path.split('.').filter(p => p !== 'root')
  let current = data
  
  for (const part of pathParts) {
    if (current && typeof current === 'object') {
      if (Array.isArray(current)) {
        const index = parseInt(part)
        current = current[index]
      } else {
        current = current[part]
      }
    } else {
      return null
    }
  }
  
  return current
}

const handleNodeClick = (node: any) => {
  console.log('Node clicked:', node)
  selectedNodeData.value = node
  
  // 獲取節點完整數據
  if (node.path && parsedViewingValue.value) {
    const nodeData = getNodeDataAtPath(parsedViewingValue.value, node.path)
    selectedNodeData.value = {
      ...node,
      fullData: nodeData
    }
  }
}

const handlePreviewNodeClick = (node: any) => {
  console.log('Preview node clicked:', node)
  previewSelectedNodeData.value = node
  
  // 獲取節點完整數據
  if (node.path && parsedCurrentValue.value) {
    const nodeData = getNodeDataAtPath(parsedCurrentValue.value, node.path)
    previewSelectedNodeData.value = {
      ...node,
      fullData: nodeData
    }
  }
}

const copySelectedData = async (isPreview = false) => {
  const nodeData = isPreview ? previewSelectedNodeData.value : selectedNodeData.value
  
  if (!nodeData) {
    alert('請先選擇一個節點')
    return
  }
  
  let copyText = ''
  
  // 如果有完整數據，複製完整的JSON
  if (nodeData.fullData !== undefined) {
    if (typeof nodeData.fullData === 'object') {
      copyText = JSON.stringify(nodeData.fullData, null, 2)
    } else {
      copyText = String(nodeData.fullData)
    }
  } else if (nodeData.value !== undefined) {
    // 否則複製節點值
    if (typeof nodeData.value === 'object') {
      copyText = JSON.stringify(nodeData.value, null, 2)
    } else {
      copyText = String(nodeData.value)
    }
  } else {
    copyText = JSON.stringify(nodeData, null, 2)
  }
  
  try {
    await navigator.clipboard.writeText(copyText)
    copySuccess.value = true
    setTimeout(() => {
      copySuccess.value = false
    }, 2000)
  } catch (err) {
    console.error('Failed to copy:', err)
    // 回退方案
    const textArea = document.createElement('textarea')
    textArea.value = copyText
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
    copySuccess.value = true
    setTimeout(() => {
      copySuccess.value = false
    }, 2000)
  }
}

const formatCurrentValue = () => {
  try {
    const parsed = JSON.parse(currentValue.value)
    // 格式化JSON，使用2個空格縮進
    currentValue.value = JSON.stringify(parsed, null, 2)
  } catch (err) {
    // 如果解析失敗，顯示錯誤提示
    alert('無法格式化：不是有效的JSON格式')
  }
}

const formatValue = (value: any): string => {
  if (typeof value === 'object') {
    return JSON.stringify(value, null, 2)
  }
  return String(value)
}

// 监听路由参数变化
watch(() => route.params.id, async (newId) => {
  if (newId) {
    connectionId.value = parseInt(newId as string)
    // 重置状态
    error.value = null
    keys.value = []
    prefixFilter.value = ''
    // 重新加载数据
    await loadConnection()
    await loadKeys()
  }
})

onMounted(async () => {
  await loadConnection()
  await loadKeys()
})
</script>
