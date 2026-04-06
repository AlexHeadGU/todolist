import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as loginApi, register as registerApi } from '../services/api'
import type { User } from '../types'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const user = ref<User | null>(null)

  const isAuthenticated = computed(() => !!token.value)

  const setAuth = (newToken: string, newUser: User) => {
    token.value = newToken
    user.value = newUser
    localStorage.setItem('token', newToken)
  }

  const clearAuth = () => {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
  }

  const login = async (email: string, password: string) => {
    try {
      const response = await loginApi({ email, password })
      setAuth(response.data.token, response.data.user)
      return { success: true, error: null }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message || 'Login failed' }
    }
  }

  const register = async (email: string, password: string) => {
    try {
      const response = await registerApi({ email, password })
      return { success: true, user: response.data, error: null }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message || 'Registration failed' }
    }
  }

  const logout = () => {
    clearAuth()
  }

  return {
    token,
    user,
    isAuthenticated,
    login,
    register,
    logout
  }
})