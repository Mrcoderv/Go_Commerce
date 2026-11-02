package logger

import (
	"log/slog"
	"os"
)

type sloglogger struct {
	logger *slog.Logger
}

func NewSlogLogger(logger *slog.Logger) *sloglogger {
	var level slog.Level
	switch {
	case config.LogLevel == "debug":
		level = slog.LevelDebug
	case config.LogLevel == "info":
		level = slog.LevelInfo
	case config.LogLevel == "warn":
		level = slog.LevelWarn
	case config.LogLevel == "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	return &sloglogger{
		logger: logger,
	}
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	return &sloglogger{
		logger: slog.New(handler),
	}
}

type logger interface {
	info(msg string, args ...any)
	error(msg string, args ...any)
	Debug(msg string, args ...any)
	warn(msg string, args ...any)
}

func (l *sloglogger) info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *sloglogger) error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *sloglogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *sloglogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}
