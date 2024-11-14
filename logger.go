package gorabbit

import (
	"fmt"
	"log/slog"
	"os"
)

const loggerPrefix = "[gorabbit]"

type ILogger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	fatal(msg string, args ...any)
}

type Logger struct {
	ILogger
}

func (l Logger) Info(msg string, args ...any) {
	slog.Info(fmt.Sprintf("%s %s", loggerPrefix, msg), args...)
}

func (l Logger) Warn(msg string, args ...any) {
	slog.Warn(fmt.Sprintf("%s %s", loggerPrefix, msg), args...)
}

func (l Logger) Error(msg string, args ...any) {
	slog.Error(fmt.Sprintf("%s %s", loggerPrefix, msg), args...)
}

func (l Logger) fatal(msg string, args ...any) {
	slog.Error(fmt.Sprintf("%s %s", loggerPrefix, msg), args...)
	os.Exit(1)
}
