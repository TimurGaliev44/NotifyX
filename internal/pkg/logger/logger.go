package logger

import (
	"log/slog"
	"os"
)

func New(serviceName string, level slog.Level) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	return slog.New(handler).With("service", serviceName)
}
