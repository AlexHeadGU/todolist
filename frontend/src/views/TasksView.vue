<template>
  <div>
    <h2>My Tasks</h2>
    
    <!-- Add Task Form -->
    <div>
      <h3>Add New Task</h3>
      <form @submit.prevent="handleAddTask">
        <input 
          type="text" 
          v-model="newTaskTitle" 
          placeholder="Task title" 
          required 
        />
        <input 
          type="text" 
          v-model="newTaskDescription" 
          placeholder="Description (optional)" 
        />
        <button type="submit" :disabled="tasksStore.loading">
          Add Task
        </button>
      </form>
    </div>

    <!-- Error Message -->
    <div v-if="tasksStore.error" style="color: red">
      {{ tasksStore.error }}
    </div>

    <!-- Loading -->
    <div v-if="tasksStore.loading">
      Loading...
    </div>

    <!-- Tasks List -->
    <div v-else>
      <div v-if="tasksStore.tasks.length === 0">
        <p>No tasks yet. Create your first task!</p>
      </div>
      
      <div v-for="task in tasksStore.tasks" :key="task.id">
        <div>
          <input 
            type="checkbox" 
            :checked="task.status === 'done'"
            @change="toggleTaskStatus(task)"
          />
          
          <div v-if="editingTaskId === task.id">
            <input 
              type="text" 
              v-model="editForm.title" 
              placeholder="Title"
            />
            <input 
              type="text" 
              v-model="editForm.description" 
              placeholder="Description"
            />
            <select v-model="editForm.status">
              <option value="pending">Pending</option>
              <option value="done">Done</option>
            </select>
            <button @click="saveEdit(task.id)">Save</button>
            <button @click="cancelEdit">Cancel</button>
          </div>
          
          <div v-else>
            <div>
              <strong>{{ task.title }}</strong>
              <span> ({{ task.status }})</span>
            </div>
            <div v-if="task.description">{{ task.description }}</div>
            <div>
              <small>Created: {{ formatDate(task.created_at) }}</small>
            </div>
            <div>
              <button @click="startEdit(task)">Edit</button>
              <button @click="deleteTask(task.id)">Delete</button>
            </div>
          </div>
        </div>
        <hr />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useTasksStore } from '../stores/tasks'
import type { Task, UpdateTaskRequest } from '../types'

const tasksStore = useTasksStore()

const newTaskTitle = ref('')
const newTaskDescription = ref('')

const editingTaskId = ref<number | null>(null)
const editForm = ref<UpdateTaskRequest>({})

onMounted(() => {
  tasksStore.fetchTasks()
})

const handleAddTask = async () => {
  if (!newTaskTitle.value.trim()) return
  
  const result = await tasksStore.addTask({
    title: newTaskTitle.value,
    description: newTaskDescription.value || undefined
  })
  
  if (result.success) {
    newTaskTitle.value = ''
    newTaskDescription.value = ''
  }
}

const toggleTaskStatus = (task: Task) => {
  tasksStore.toggleTaskStatus(task.id, task.status)
}

const startEdit = (task: Task) => {
  editingTaskId.value = task.id
  editForm.value = {
    title: task.title,
    description: task.description,
    status: task.status
  }
}

const saveEdit = async (taskId: number) => {
  const result = await tasksStore.editTask(taskId, editForm.value)
  if (result.success) {
    cancelEdit()
  }
}

const cancelEdit = () => {
  editingTaskId.value = null
  editForm.value = {}
}

const deleteTask = async (taskId: number) => {
  if (confirm('Are you sure you want to delete this task?')) {
    await tasksStore.removeTask(taskId)
  }
}

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString()
}
</script>