package logger

import (
	"log/slog"
	"os"
	"sync"
)

// singleton logger instance
var (
	log  *slog.Logger
	once sync.Once
)

// Logger modes
const (
	LOGGER_MODE_JSON = "json"
	LOGGER_MODE_TEXT = "text"
)

// Init initializes the logger.
// mode = "json" (production) or "text" (development)
func Init(mode string) {
	once.Do(func() {
		var handler slog.Handler
		opts := &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: true, // include file:line in logs
		}
		switch mode {
		case "text":
			handler = slog.NewTextHandler(os.Stdout, opts)
			opts.Level = slog.LevelDebug // dev: more verbose
		default:
			handler = slog.NewJSONHandler(os.Stdout, opts)
		}
		log = slog.New(handler)
	})
}

// Info logs an info-level message.
func Info(msg string, args ...any) {
	log.Info(msg, args...)
}

// Error logs an error-level message.
func Error(msg string, args ...any) {
	log.Error(msg, args...)
}

// Debug logs a debug-level message.
func Debug(msg string, args ...any) {
	log.Debug(msg, args...)
}

// Warn logs a warning-level message.
func Warn(msg string, args ...any) {
	log.Warn(msg, args...)
}

// GetLogger exposes the underlying slog.Logger if needed.
func GetLogger() *slog.Logger {

	return log
}
