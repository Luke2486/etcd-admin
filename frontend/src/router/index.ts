import { createRouter, createWebHistory } from 'vue-router'
import { useAuth } from '@/composables/useAuth'

// Views
import Login from '@/views/Login.vue'
import Connections from '@/views/Connections.vue'
import KeyValue from '@/views/KeyValue.vue'

const routes = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresGuest: true }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Connections,
    meta: { requiresAuth: true }
  },
  {
    path: '/connections',
    name: 'Connections',
    component: Connections,
    meta: { requiresAuth: true }
  },
  {
    path: '/kv/:id',
    name: 'KeyValue',
    component: KeyValue,
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation guards
router.beforeEach((to, _from, next) => {
  const { isAuthenticated } = useAuth()

  if (to.meta.requiresAuth && !isAuthenticated.value) {
    next('/login')
  } else if (to.meta.requiresGuest && isAuthenticated.value) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
