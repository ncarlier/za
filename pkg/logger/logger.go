package logger

import (
	"log/slog"
	"os"
)

// Configure logger
func Configure(format, level string) {
	logLevel := slog.LevelInfo
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	}

	opts := slog.HandlerOptions{
		Level:     logLevel,
		AddSource: logLevel == slog.LevelDebug,
	}

	var logger *slog.Logger
	if format == "text" {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &opts))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &opts))
	}

	slog.SetDefault(logger)
}
