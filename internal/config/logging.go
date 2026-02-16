package config

import (
	"log/slog"
	"os"
)

// InitLogging configures structured JSON logging for Lambda environments.
func InitLogging() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(handler))
}
