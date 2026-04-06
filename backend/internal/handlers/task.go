// обработчики регистрации и входа

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/AlexHeadGU/todolist/internal/models"
	"github.com/AlexHeadGU/todolist/internal/service"
	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// Create обрабатывает создание новой задачи
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	// 1. Извлекаем user_id из контекста (установлен middleware аутентификации)
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. Парсим тело запроса
	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 3. Базовая валидация
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// 4. Вызываем сервис для создания задачи
	task, err := h.taskService.Create(req.Title, req.Description, userID)
	if err != nil {
		// Обработка разных типов ошибок
		switch err.Error() {
		case "title cannot be empty":
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "Failed to create task", http.StatusInternalServerError)
		}
		return
	}

	// 5. Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(task)
}

// GetAll возвращает все задачи текущего пользователя
func (h *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Извлекаем user_id из контекста
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Получаем задачи
	tasks, err := h.taskService.GetUserTasks(userID)
	if err != nil {
		http.Error(w, "Failed to get tasks", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Получаем ID из URL
	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := h.taskService.GetTaskByID(taskID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// Update обрабатывает PUT /api/tasks/{id} (полное обновление)
func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Валидация
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	if req.Status != "pending" && req.Status != "done" {
		http.Error(w, "Status must be 'pending' or 'done'", http.StatusBadRequest)
		return
	}

	task, err := h.taskService.UpdateTask(taskID, userID, req.Title, req.Description, req.Status)
	if err != nil {
		if err.Error() == "task not found or access denied" {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update task", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// Patch обрабатывает PATCH /api/tasks/{id} (частичное обновление)
func (h *TaskHandler) Patch(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var req models.PatchTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task, err := h.taskService.PatchTask(taskID, userID, req.Title, req.Description, req.Status)
	if err != nil {
		if err.Error() == "task not found or access denied" {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// 1. Получаем user_id из контекста
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. Получаем task_id из URL
	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// 3. Вызываем сервис для удаления
	err = h.taskService.DeleteTask(taskID, userID)
	if err != nil {
		if err.Error() == "task not found or access denied" {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		}
		return
	}

	// 4. Успешное удаление (204 No Content)
	w.WriteHeader(http.StatusNoContent)
}
