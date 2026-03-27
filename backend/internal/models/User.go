package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // не возвращаем в JSON
	CreatedAt    time.Time `json:"created_at"`
}
