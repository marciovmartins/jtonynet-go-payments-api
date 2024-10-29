package slogStrategy

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/jtonynet/go-payments-api/config"
)

/*
	font: https://betterstack.com/community/guides/logging/logging-in-go/
*/

var levelNameToValue = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

type Logger struct {
	instance *slog.Logger
}

func New(cfg config.Logger) (*Logger, error) {
	opts := &slog.HandlerOptions{
		AddSource: cfg.AddSource,
		Level:     levelNameToValue[cfg.Level],
	}

	var handler slog.Handler
	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case "text":
		handler = slog.NewTextHandler(os.Stdout, opts)
	default:
		return nil, fmt.Errorf("log strategy %s format: %s not suported", cfg.Strategy, cfg.Format)
	}

	instance := slog.New(handler)

	return &Logger{
		instance: instance,
	}, nil
}

func (l Logger) Info(msg string, args ...interface{}) {
	l.instance.Info(msg, args...)
}

func (l Logger) Debug(msg string, args ...interface{}) {
	l.instance.Debug(msg, args...)
}

func (l Logger) Warn(msg string, args ...interface{}) {
	l.instance.Warn(msg, args...)
}

func (l Logger) Error(msg string, args ...interface{}) {
	l.instance.Error(msg, args...)
}
