package repository

import (
	"database/sql"
	"fmt"

	"github.com/AlexHeadGU/todolist/internal/models"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Create создаёт новую задачу
func (r *TaskRepository) Create(title, description string, userID int) (*models.Task, error) {
	query := `
        INSERT INTO tasks (title, description, user_id, status, created_at, updated_at)
        VALUES ($1, $2, $3, 'pending', NOW(), NOW())
        RETURNING id, user_id, title, description, status, created_at, updated_at
    `

	task := &models.Task{}
	err := r.db.QueryRow(query, title, description, userID).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create• task: %w", err)
	}

	return task, nil
}

// GetByUserID возвращает все задачи пользователя
func (r *TaskRepository) GetByUserID(userID int) ([]models.Task, error) {
	query := `
        SELECT id, user_id, title, description, status, created_at, updated_at
        FROM tasks
        WHERE user_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetByID возвращает задачу по ее ID
func (r *TaskRepository) GetByID(taskID int) (*models.Task, error) {
	query := `
        SELECT id, user_id, title, description, status, created_at, updated_at
        FROM tasks
        WHERE id = $1
    `

	task := &models.Task{}
	err := r.db.QueryRow(query, taskID).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // задача не найдена (не ошибка)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return task, nil

}

// Update обновляет все поля задачи (PUT)
func (r *TaskRepository) Update(taskID, userID int, title, description, status string) (*models.Task, error) {
	query := `
        UPDATE tasks
        SET title = $1, description = $2, status = $3, updated_at = NOW()
        WHERE id = $4 AND user_id = $5
        RETURNING id, user_id, title, description, status, created_at, updated_at
    `

	task := &models.Task{}
	err := r.db.QueryRow(query, title, description, status, taskID, userID).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // задача не найдена или не принадлежит пользователю
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return task, nil
}

// Patch обновляет только переданные поля (PATCH)
func (r *TaskRepository) Patch(taskID, userID int, title, description, status *string) (*models.Task, error) {
	// Сначала получаем текущую задачу
	currentTask, err := r.GetByID(taskID)
	if err != nil {
		return nil, err
	}
	if currentTask == nil {
		return nil, nil
	}

	// Проверяем владельца
	if currentTask.UserID != userID {
		return nil, nil
	}

	// Обновляем только те поля, которые переданы
	if title != nil {
		currentTask.Title = *title
	}
	if description != nil {
		currentTask.Description = *description
	}
	if status != nil {
		currentTask.Status = *status
	}

	// Сохраняем обновлённую задачу
	query := `
        UPDATE tasks
        SET title = $1, description = $2, status = $3, updated_at = NOW()
        WHERE id = $4 AND user_id = $5
        RETURNING id, user_id, title, description, status, created_at, updated_at
    `

	updatedTask := &models.Task{}
	err = r.db.QueryRow(query, currentTask.Title, currentTask.Description, currentTask.Status, taskID, userID).Scan(
		&updatedTask.ID,
		&updatedTask.UserID,
		&updatedTask.Title,
		&updatedTask.Description,
		&updatedTask.Status,
		&updatedTask.CreatedAt,
		&updatedTask.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to patch task: %w", err)
	}

	return updatedTask, nil
}

// Delete удаляет задачу по ID и user_id (проверка владельца)
func (r *TaskRepository) Delete(taskID, userID int) error {
	query := `
        DELETE FROM tasks
        WHERE id = $1 AND user_id = $2
    `

	result, err := r.db.Exec(query, taskID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	// Проверяем, была ли удалена хотя бы одна запись
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found or access denied")
	}

	return nil
}
