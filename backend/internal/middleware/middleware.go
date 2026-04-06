package middleware

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

// Setup настраивает и возвращает chi роутер со всеми необходимыми middleware
func Setup() *chi.Mux {
	r := chi.NewRouter()

	// CORS middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	// Стандартные middleware chi
	r.Use(middleware.Logger)    // логирование запросов
	r.Use(middleware.Recoverer) // восстановление после паники
	r.Use(corsMiddleware.Handler)

	return r
}
