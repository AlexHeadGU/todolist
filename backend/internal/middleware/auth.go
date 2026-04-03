package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/AlexHeadGU/todolist/internal/service"
)

func Auth(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Получаем заголовок Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// 2. Проверяем формат "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format. Use: Bearer <token>", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			// 3. Валидируем токен
			userID, err := authService.ValidateToken(token)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// 4. Кладём user_id в контекст
			ctx := context.WithValue(r.Context(), "user_id", userID)

			// 5. Передаём запрос дальше
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
