package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/AlexHeadGU/todolist/internal/logger"
	"github.com/AlexHeadGU/todolist/internal/models"
	"github.com/AlexHeadGU/todolist/internal/service"
	"github.com/AlexHeadGU/todolist/internal/utils"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		utils.SendUnauthorizedError(w, "Unauthorized")
		return
	}

	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("Invalid task creation request", "error", err)
		utils.SendValidationError(w, "Invalid request body")
		return
	}

	if req.Title == "" {
		logger.Warn("Task creation failed - missing title", "user_id", userID)
		utils.SendValidationError(w, "Title is required")
		return
	}

	task, err := h.taskService.Create(req.Title, req.Description, userID)
	if err != nil {
		logger.Error("Failed to create task", "error", err, "user_id", userID)
		utils.SendInternalError(w, err, "Failed to create task")
		return
	}

	logger.Info("Task created", "task_id", task.ID, "user_id", userID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		utils.SendUnauthorizedError(w, "Unauthorized")
		return
	}

	tasks, err := h.taskService.GetUserTasks(userID)
	if err != nil {
		logger.Error("Failed to get tasks", "error", err, "user_id", userID)
		utils.SendInternalError(w, err, "Failed to get tasks")
		return
	}

	logger.Debug("Tasks retrieved", "count", len(tasks), "user_id", userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		utils.SendUnauthorizedError(w, "Unauthorized")
		return
	}

	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		logger.Warn("Invalid task ID", "id", taskIDStr, "user_id", userID)
		utils.SendValidationError(w, "Invalid task ID")
		return
	}

	task, err := h.taskService.GetTaskByID(taskID, userID)
	if err != nil {
		if err.Error() == "task not found" {
			logger.Warn("Task not found", "task_id", taskID, "user_id", userID)
			utils.SendNotFoundError(w, "Task not found")
		} else if err.Error() == "access denied" {
			logger.Warn("Access denied to task", "task_id", taskID, "user_id", userID)
			utils.SendNotFoundError(w, "Task not found")
		} else {
			logger.Error("Failed to get task", "error", err, "task_id", taskID)
			utils.SendInternalError(w, err, "Failed to get task")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		utils.SendUnauthorizedError(w, "Unauthorized")
		return
	}

	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		logger.Warn("Invalid task ID for update", "id", taskIDStr)
		utils.SendValidationError(w, "Invalid task ID")
		return
	}

	var req models.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("Invalid update request", "error", err)
		utils.SendValidationError(w, "Invalid request body")
		return
	}

	if req.Title == "" {
		utils.SendValidationError(w, "Title is required")
		return
	}
	if req.Status != "pending" && req.Status != "done" {
		utils.SendValidationError(w, "Status must be 'pending' or 'done'")
		return
	}

	task, err := h.taskService.UpdateTask(taskID, userID, req.Title, req.Description, req.Status)
	if err != nil {
		if err.Error() == "task not found or access denied" {
			logger.Warn("Task not found for update", "task_id", taskID, "user_id", userID)
			utils.SendNotFoundError(w, "Task not found")
		} else {
			logger.Error("Failed to update task", "error", err, "task_id", taskID)
			utils.SendInternalError(w, err, "Failed to update task")
		}
		return
	}

	logger.Info("Task updated", "task_id", taskID, "user_id", userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) Patch(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		utils.SendUnauthorizedError(w, "Unauthorized")
		return
	}

	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		logger.Warn("Invalid task ID for patch", "id", taskIDStr)
		utils.SendValidationError(w, "Invalid task ID")
		return
	}

	var req models.PatchTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Warn("Invalid patch request", "error", err)
		utils.SendValidationError(w, "Invalid request body")
		return
	}

	task, err := h.taskService.PatchTask(taskID, userID, req.Title, req.Description, req.Status)
	if err != nil {
		if err.Error() == "task not found or access denied" {
			logger.Warn("Task not found for patch", "task_id", taskID, "user_id", userID)
			utils.SendNotFoundError(w, "Task not found")
		} else {
			logger.Error("Failed to patch task", "error", err, "task_id", taskID)
			utils.SendInternalError(w, err, err.Error())
		}
		return
	}

	logger.Info("Task patched", "task_id", taskID, "user_id", userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		utils.SendUnauthorizedError(w, "Unauthorized")
		return
	}

	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		logger.Warn("Invalid task ID for delete", "id", taskIDStr)
		utils.SendValidationError(w, "Invalid task ID")
		return
	}

	err = h.taskService.DeleteTask(taskID, userID)
	if err != nil {
		if err.Error() == "task not found or access denied" {
			logger.Warn("Task not found for delete", "task_id", taskID, "user_id", userID)
			utils.SendNotFoundError(w, "Task not found")
		} else {
			logger.Error("Failed to delete task", "error", err, "task_id", taskID)
			utils.SendInternalError(w, err, "Failed to delete task")
		}
		return
	}

	logger.Info("Task deleted", "task_id", taskID, "user_id", userID)
	w.WriteHeader(http.StatusNoContent)
}
