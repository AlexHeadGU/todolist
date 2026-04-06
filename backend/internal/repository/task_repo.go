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

// Update обновляет задачу по ее ID
// func (r *TaskRepository) Update(userID int) ([]models.Task, error) {}

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
