package main

import (
	"net/http"

	logger "github.com/AlexHeadGU/todolist/internal/middleware"
)

func main() {
	r := logger.Setup()

	// Тестовый маршрут
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// Запуск сервера
	http.ListenAndServe(":3000", r)
}
