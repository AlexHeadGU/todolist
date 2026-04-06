import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getTasks, createTask, updateTask, patchTask, deleteTask } from '../services/api'
import type { Task, CreateTaskRequest, UpdateTaskRequest } from '../types'

export const useTasksStore = defineStore('tasks', () => {
  const tasks = ref<Task[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetchTasks = async () => {
    loading.value = true
    error.value = null
    try {
      const response = await getTasks()
      tasks.value = response.data
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Failed to fetch tasks'
    } finally {
      loading.value = false
    }
  }

  const addTask = async (taskData: CreateTaskRequest) => {
    loading.value = true
    error.value = null
    try {
      const response = await createTask(taskData)
      tasks.value.unshift(response.data)
      return { success: true, task: response.data, error: null }
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Failed to create task'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  const editTask = async (id: number, taskData: UpdateTaskRequest) => {
    loading.value = true
    error.value = null
    try {
      const response = await patchTask(id, taskData)
      const index = tasks.value.findIndex(t => t.id === id)
      if (index !== -1) {
        tasks.value[index] = response.data
      }
      return { success: true, task: response.data, error: null }
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Failed to update task'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  const removeTask = async (id: number) => {
    loading.value = true
    error.value = null
    try {
      await deleteTask(id)
      tasks.value = tasks.value.filter(t => t.id !== id)
      return { success: true, error: null }
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Failed to delete task'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  const toggleTaskStatus = async (id: number, currentStatus: string) => {
    const newStatus = currentStatus === 'pending' ? 'done' : 'pending'
    return editTask(id, { status: newStatus })
  }

  return {
    tasks,
    loading,
    error,
    fetchTasks,
    addTask,
    editTask,
    removeTask,
    toggleTaskStatus
  }
})