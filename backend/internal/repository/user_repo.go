//# работа с БД для пользователей

package repository

import (
	"database/sql"
	"fmt"

	"github.com/AlexHeadGU/todolist/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create создаёт нового пользователя
func (r *UserRepository) Create(email, passwordHash string) (*models.User, error) {
	query := `
        INSERT INTO users (email, password_hash, created_at)
        VALUES ($1, $2, NOW())
        RETURNING id, email, created_at
    `

	user := &models.User{}
	err := r.db.QueryRow(query, email, passwordHash).Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// FindByEmail ищет пользователя по email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	query := `
        SELECT id, email, password_hash, created_at
        FROM users
        WHERE email = $1
    `

	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil // пользователь не найден
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

// ExistsByEmail проверяет существование пользователя
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.db.QueryRow(query, email).Scan(&exists)
	return exists, err
}
