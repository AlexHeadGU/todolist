package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/AlexHeadGU/todolist/internal/logger"
	"github.com/AlexHeadGU/todolist/internal/service"
	"github.com/AlexHeadGU/todolist/internal/utils"
)

func Auth(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Логируем входящий запрос
			logger.Debug("Auth middleware", "path", r.URL.Path, "method", r.Method)

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Warn("Missing authorization header", "path", r.URL.Path)
				utils.SendUnauthorizedError(w, "Authorization header required")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				logger.Warn("Invalid auth header format", "path", r.URL.Path)
				utils.SendUnauthorizedError(w, "Invalid authorization header format. Use: Bearer <token>")
				return
			}

			token := parts[1]

			userID, err := authService.ValidateToken(token)
			if err != nil {
				logger.Warn("Invalid token", "path", r.URL.Path, "error", err)
				utils.SendUnauthorizedError(w, "Invalid or expired token")
				return
			}

			logger.Debug("Token validated", "user_id", userID)
			ctx := context.WithValue(r.Context(), "user_id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
