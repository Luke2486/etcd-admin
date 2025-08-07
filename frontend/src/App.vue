<template>
  <div id="app">
    <Layout v-if="isAuthenticated">
      <router-view />
    </Layout>
    <router-view v-else />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuth } from '@/composables/useAuth'
import Layout from '@/components/Layout.vue'

const { isAuthenticated, initAuth } = useAuth()

onMounted(() => {
  initAuth()
})
</script>

<style>
/* Global styles are handled by TailwindCSS */
</style>
