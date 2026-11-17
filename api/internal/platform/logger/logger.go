// Package logger menyiapkan slog logger standar.
package logger

import (
	"log/slog"
	"os"
)

func Setup(level string) *slog.Logger {
	var lv slog.Level
	switch level {
	case "debug":
		lv = slog.LevelDebug
	case "warn":
		lv = slog.LevelWarn
	case "error":
		lv = slog.LevelError
	default:
		lv = slog.LevelInfo
	}
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lv})
	lg := slog.New(handler)
	slog.SetDefault(lg)
	return lg
}
