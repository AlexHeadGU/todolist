//# бизнес-логика (хеширование, JWT)

package service

import (
	"fmt"

	"github.com/AlexHeadGU/todolist/internal/models"
	"github.com/AlexHeadGU/todolist/internal/repository"
)

type TaskService struct {
	taskRepo *repository.TaskRepository
}

func NewTaskService(taskRepo *repository.TaskRepository) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
	}
}

// Создает новую задачу
func (s *TaskService) Create(title, description string, userID int) (*models.Task, error) {
	// Создаём задачу
	task, err := s.taskRepo.Create(title, description, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}

// GetUserTasks возвращает все задачи пользователя
func (s *TaskService) GetUserTasks(userID int) ([]models.Task, error) {
	tasks, err := s.taskRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user tasks: %w", err)
	}
	return tasks, nil
}

// GetTaskByID возвращает задачу по ID
func (s *TaskService) GetTaskByID(taskID, userID int) (*models.Task, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("task not found")
	}
	// Проверяем, что задача принадлежит пользователю
	if task.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}
	return task, nil
}

// UpdateTask полностью обновляет задачу (PUT)
func (s *TaskService) UpdateTask(taskID, userID int, title, description, status string) (*models.Task, error) {
	// Валидация
	if title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}
	if status != "pending" && status != "done" {
		return nil, fmt.Errorf("status must be 'pending' or 'done'")
	}

	task, err := s.taskRepo.Update(taskID, userID, title, description, status)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("task not found or access denied")
	}

	return task, nil
}

// PatchTask частично обновляет задачу (PATCH)
func (s *TaskService) PatchTask(taskID, userID int, title, description, status *string) (*models.Task, error) {
	// Валидация
	if title != nil && *title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}
	if status != nil && *status != "pending" && *status != "done" {
		return nil, fmt.Errorf("status must be 'pending' or 'done'")
	}

	task, err := s.taskRepo.Patch(taskID, userID, title, description, status)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("task not found or access denied")
	}

	return task, nil
}

// DeleteTask удаляет задачу (с проверкой владельца)
func (s *TaskService) DeleteTask(taskID, userID int) error {
	return s.taskRepo.Delete(taskID, userID)
}
