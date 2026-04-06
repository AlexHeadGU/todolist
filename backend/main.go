package main

import (
	"database/sql"
	"embed"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/AlexHeadGU/todolist/internal/config"
	"github.com/AlexHeadGU/todolist/internal/handlers"
	"github.com/AlexHeadGU/todolist/internal/logger"
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
	// Инициализация логгера
	logger.InitLogger()
	logger.Info("Starting TodoList application...")

	// Загружаем конфигурацию
	config.LoadConfig()
	logger.Debug("Configuration loaded")

	// 1. Подключение к БД
	logger.Info("Connecting to database...")
	db, err := sql.Open("postgres", "postgres://headgus:15948256@144.31.69.223:5432/llm_db?sslmode=disable")
	if err != nil {
		logger.Error("Failed to open database connection", "error", err)
		panic(err)
	}
	defer func() {
		logger.Info("Closing database connection")
		db.Close()
	}()

	// 2. Проверяем соединение
	logger.Info("Pinging database...")
	err = db.Ping()
	if err != nil {
		logger.Error("Cannot ping database", "error", err)
		panic(err)
	}
	logger.Info("Database connected successfully")

	// 3. Инициализация драйвера миграций для базы данных
	logger.Info("Initializing migration driver...")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Error("Cannot create database driver", "error", err)
		panic(err)
	}
	logger.Debug("Database driver created")

	// 4. Создаём источник миграций из встроенных файлов
	logger.Info("Loading embedded migration files...")
	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		logger.Error("Cannot create source from embed", "error", err)
		panic(err)
	}
	logger.Debug("Migration source created")

	// 5. Создаём мигратор
	logger.Info("Creating migrator...")
	m, err := migrate.NewWithInstance(
		"iofs",     // имя источника
		source,     // экземпляр источника
		"postgres", // имя базы данных
		driver,     // экземпляр драйвера БД
	)
	if err != nil {
		logger.Error("Failed to create migrator", "error", err)
		panic(err)
	}
	logger.Debug("Migrator created")

	// 6. Применяем миграции
	logger.Info("Applying migrations...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Error("Migration failed", "error", err)
		panic(err)
	}
	if err == migrate.ErrNoChange {
		logger.Info("No new migrations to apply")
	} else {
		logger.Info("Migrations applied successfully")
	}

	// Инициализация репозиториев, сервисов, хендлеров
	logger.Info("Initializing repositories...")
	userRepo := repository.NewUserRepository(db) // репозиторий пользователей
	taskRepo := repository.NewTaskRepository(db) // репозиторий задач
	logger.Debug("Repositories created")

	logger.Info("Initializing services...")
	authService := service.NewAuthService(userRepo) // сервис аутентификации
	taskService := service.NewTaskService(taskRepo) // сервис задач
	logger.Debug("Services created")

	logger.Info("Initializing handlers...")
	authHandler := handlers.NewAuthHandler(authService) // хендлер аутентификации
	taskHandler := handlers.NewTaskHandler(taskService) // хендлер задач
	logger.Debug("Handlers created")

	// 7. Настройка роутера
	logger.Info("Setting up router...")
	r := middleware.Setup()
	logger.Debug("Router configured")

	// Публичные маршруты
	logger.Debug("Registering public routes")
	r.Post("/api/register", authHandler.Register)
	r.Post("/api/login", authHandler.Login)

	// Тестовый маршрут
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// Защищенные маршруты
	logger.Debug("Registering protected routes")
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(authService))

		r.Post("/api/tasks", taskHandler.Create)
		r.Get("/api/tasks", taskHandler.GetAll)
		r.Get("/api/tasks/{id}", taskHandler.GetByID)
		r.Put("/api/tasks/{id}", taskHandler.Update)
		r.Patch("/api/tasks/{id}", taskHandler.Patch)
		r.Delete("/api/tasks/{id}", taskHandler.Delete)
	})
	logger.Info("All routes registered")

	// 8. Запуск сервера
	logger.Info("Server starting on :3000")
	logger.Info("Ready to accept requests")
	if err := http.ListenAndServe(":3000", r); err != nil {
		logger.Error("Server error", "error", err)
		panic(err)
	}
}
