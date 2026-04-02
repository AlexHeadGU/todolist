package middleware

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

// Setup настраивает и возвращает chi роутер со всеми необходимыми middleware
func Setup() *chi.Mux {
	r := chi.NewRouter()

	// Стандартные middleware chi
	r.Use(middleware.Logger)    // логирование запросов
	r.Use(middleware.Recoverer) // восстановление после паники

	// CORS middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // на время разработки
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	})
	r.Use(corsMiddleware.Handler)

	return r
}
