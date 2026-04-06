export interface User {
  id: number
  email: string
  created_at: string
}

export interface Task {
  id: number
  user_id: number
  title: string
  description: string
  status: 'pending' | 'done'
  created_at: string
  updated_at: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
}

export interface CreateTaskRequest {
  title: string
  description?: string
}

export interface UpdateTaskRequest {
  title?: string
  description?: string
  status?: 'pending' | 'done'
}

export interface LoginResponse {
  token: string
  user: User
}