package utils

import (
	"encoding/json"
	"net/http"

	"github.com/AlexHeadGU/todolist/internal/logger"
)

// ErrorResponse структура для единого формата ошибок
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Status  int    `json:"-"`
}

// SendError отправляет JSON ошибку клиенту
func SendError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}

// SendValidationError отправляет ошибку валидации
func SendValidationError(w http.ResponseWriter, message string) {
	SendError(w, http.StatusBadRequest, message)
}

// SendNotFoundError отправляет 404
func SendNotFoundError(w http.ResponseWriter, message string) {
	SendError(w, http.StatusNotFound, message)
}

// SendUnauthorizedError отправляет 401
func SendUnauthorizedError(w http.ResponseWriter, message string) {
	SendError(w, http.StatusUnauthorized, message)
}

// SendInternalError отправляет 500 с логированием
func SendInternalError(w http.ResponseWriter, err error, message string) {
	logger.Error(message, "error", err)
	SendError(w, http.StatusInternalServerError, "Internal server error")
}
