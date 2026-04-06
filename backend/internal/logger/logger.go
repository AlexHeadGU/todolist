package logger

import (
	"log"
	"log/slog"
	"os"
)

var Log *slog.Logger

func InitLogger() {
	// Создаём JSON логгер (удобно для продакшена)
	Log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Можно также использовать текстовый формат для разработки:
	// Log = slog.New(slog.NewTextHandler(os.Stdout, nil))

	log.Println("Logger initialized")
}

// Вспомогательные функции для удобства
func Info(msg string, args ...interface{}) {
	Log.Info(msg, args...)
}

func Error(msg string, args ...interface{}) {
	Log.Error(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	Log.Debug(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	Log.Warn(msg, args...)
}
