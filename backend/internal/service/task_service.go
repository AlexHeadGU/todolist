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
