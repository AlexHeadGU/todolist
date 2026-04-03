// обработчики регистрации и входа

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AlexHeadGU/todolist/internal/models"
	"github.com/AlexHeadGU/todolist/internal/service"
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
