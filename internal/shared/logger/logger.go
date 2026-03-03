package logger

import (
	"log/slog"
	"os"
	"strings"
)

var Log *slog.Logger

// Setup inicializa el logger estructurado usando slog nativo de Go
func Setup(level, environment string) {
	var logLevel slog.Level

	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn", "warning":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	var handler slog.Handler
	if environment == "production" {
		// JSON estructurado para producción
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		// Formato legible para desarrollo
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	Log = slog.New(handler)
	slog.SetDefault(Log)
}

// Info registra un mensaje informativo
func Info(msg string, args ...any) {
	Log.Info(msg, args...)
}

// Debug registra un mensaje de depuración
func Debug(msg string, args ...any) {
	Log.Debug(msg, args...)
}

// Warn registra una advertencia
func Warn(msg string, args ...any) {
	Log.Warn(msg, args...)
}

// Error registra un error
func Error(msg string, args ...any) {
	Log.Error(msg, args...)
}

// With retorna un logger con atributos adicionales
func With(args ...any) *slog.Logger {
	return Log.With(args...)
}

// GetLogger retorna la instancia global del logger
func GetLogger() *slog.Logger {
	return Log
}
