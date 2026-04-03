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
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}
