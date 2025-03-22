package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
)

const (
	defaultLevel      = LevelInfo
	defaultAddSource  = false
	defaultIsJSON     = false
	defaultSetDefault = true
)

var (
	defaultWriter io.Writer = os.Stdout
	currentWriter io.Writer = os.Stdout
)

func NewLogger(opts ...Option) *Logger {
	cfg := &config{
		Level:      defaultLevel,
		AddSource:  defaultAddSource,
		IsJSON:     defaultIsJSON,
		SetDefault: defaultSetDefault,
		Writer:     defaultWriter,
	}

	logger := slog.New(cfg.createHandler(opts...))
	if cfg.SetDefault {
		slog.SetDefault(logger)
	}

	return logger
}

type config struct {
	Level      Level
	AddSource  bool
	IsJSON     bool
	SetDefault bool
	Writer     io.Writer
}

func (c *config) createHandler(opts ...Option) Handler {
	for _, opt := range opts {
		opt(c)
	}

	options := &HandlerOptions{
		Level:     c.Level,
		AddSource: c.AddSource,
	}

	var handler Handler
	if c.IsJSON {
		handler = slog.NewJSONHandler(c.Writer, options)
	} else {
		handler = slog.NewTextHandler(c.Writer, options)
	}

	return handler
}

func Writer() io.Writer {
	return currentWriter
}

func L(ctx context.Context) *Logger {
	return fromContext(ctx)
}
