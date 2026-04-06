import axios from 'axios'
import type { LoginRequest, RegisterRequest, CreateTaskRequest, UpdateTaskRequest, Task, LoginResponse, User } from '../types'

const api = axios.create({
  baseURL: (import.meta as any).env.VITE_API_URL || 'http://localhost:3000/api',
  headers: {
    'Content-Type': 'application/json'
  }
})

// Интерсептор для добавления токена
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Auth API
export const register = (data: RegisterRequest) => {
  return api.post<User>('/register', data)
}

export const login = (data: LoginRequest) => {
  return api.post<LoginResponse>('/login', data)
}

// Tasks API
export const getTasks = () => {
  return api.get<Task[]>('/tasks')
}

export const getTask = (id: number) => {
  return api.get<Task>(`/tasks/${id}`)
}

export const createTask = (data: CreateTaskRequest) => {
  return api.post<Task>('/tasks', data)
}

export const updateTask = (id: number, data: UpdateTaskRequest) => {
  return api.put<Task>(`/tasks/${id}`, data)
}

export const patchTask = (id: number, data: UpdateTaskRequest) => {
  return api.patch<Task>(`/tasks/${id}`, data)
}

export const deleteTask = (id: number) => {
  return api.delete(`/tasks/${id}`)
}