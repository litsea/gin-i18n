package i18n

import (
	"log/slog"
)

type Logger interface {
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type defaultLogger struct{}

func (l defaultLogger) Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

func (l defaultLogger) Error(msg string, args ...any) {
	slog.Error(msg, args...)
}
