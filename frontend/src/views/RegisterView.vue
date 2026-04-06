<template>
  <div>
    <h2>Register</h2>
    <form @submit.prevent="handleRegister">
      <div>
        <label>Email:</label>
        <input type="email" v-model="email" required />
      </div>
      <div>
        <label>Password:</label>
        <input type="password" v-model="password" required />
      </div>
      <div v-if="errorMessage" style="color: red">{{ errorMessage }}</div>
      <div v-if="successMessage" style="color: green">{{ successMessage }}</div>
      <button type="submit" :disabled="loading">
        {{ loading ? 'Loading...' : 'Register' }}
      </button>
    </form>
    <p>
      Already have an account? <router-link to="/login">Login</router-link>
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const handleRegister = async () => {
  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''
  
  const result = await authStore.register(email.value, password.value)
  
  if (result.success) {
    successMessage.value = 'Registration successful! Please login.'
    setTimeout(() => {
      router.push('/login')
    }, 2000)
  } else {
    errorMessage.value = result.error || 'Registration failed'
  }
  
  loading.value = false
}
</script>