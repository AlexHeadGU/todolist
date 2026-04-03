package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/AlexHeadGU/todolist/internal/config"
	"github.com/AlexHeadGU/todolist/internal/handlers"
	"github.com/AlexHeadGU/todolist/internal/middleware"

	"github.com/AlexHeadGU/todolist/internal/repository"
	"github.com/AlexHeadGU/todolist/internal/service"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

// Встраиваем папку migrations в бинарник
//
//go:embed migrations/*.sql
var migrationsFS embed.FS

func main() {

	// Загружаем конфигурацию
	config.LoadConfig()

	// 1. Подключение к БД
	db, err := sql.Open("postgres", "postgres://headgus:15948256@144.31.69.223:5432/llm_db?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// 2. Проверяем соединение
	err = db.Ping()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	fmt.Println("Успешное подключение к базе данных!")

	// 3. Инициализация драйвера миграций для базы данных
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("Failed to create database driver:", err)
	}

	// 4. Создаём источник миграций из встроенных файлов
	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		log.Fatal("Failed to create source from embed:", err)
	}

	// 5. Создаём мигратор
	m, err := migrate.NewWithInstance(
		"iofs",     // имя источника
		source,     // экземпляр источника
		"postgres", // имя базы данных
		driver,     // экземпляр драйвера БД
	)
	if err != nil {
		log.Fatal("Failed to create migrator:", err)
	}

	// 6. Применяем миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Migrations applied successfully")

	// Инициализация репозиториев, сервисов, хендлеров
	userRepo := repository.NewUserRepository(db) // репозиторий пользователей
	taskRepo := repository.NewTaskRepository(db) // репозиторий задач

	// Сервисы
	authService := service.NewAuthService(userRepo) // сервис аутентификации
	taskService := service.NewTaskService(taskRepo) // сервис задач

	// Хендлеры
	authHandler := handlers.NewAuthHandler(authService) // хендлер аутентификации
	taskHandler := handlers.NewTaskHandler(taskService) // хендлер задач

	// 7. Настройка роутера
	r := middleware.Setup()

	// Публичные маршруты
	r.Post("/api/register", authHandler.Register)
	r.Post("/api/login", authHandler.Login)

	// Защищенные маршруты
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(authService))
		r.Post("/api/tasks", taskHandler.Create)
		r.Get("/api/tasks", taskHandler.GetAll)
		r.Get("/api/tasks/{id}", taskHandler.GetByID)
		// r.Put("/api/tasks/{id}", authHandler.ChangeTask)
		// r.Patch("/api/tasks/{id}", authHandler.EditTask)
		// r.Delete("/api/tasks/{id}", authHandler.DeleteTask)
	})

	// Тестовый маршрут
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// 8. Запуск сервера
	log.Println("Server starting on :3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal("Server error:", err)
	}
}
