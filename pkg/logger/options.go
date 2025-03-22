package logger

import (
	"context"
	"io"
)

type Option func(*config)

// WithWriter logger option sets the writer, if not set, the default writer is os.Stdout.
func WithWriter(w io.Writer) Option {
	currentWriter = w
	return func(o *config) {
		o.Writer = w
	}
}

// WithLevel logger option sets the log level, if not set, the default level is Info.
func WithLevel(level string) Option {
	return func(o *config) {
		var l Level
		if err := l.UnmarshalText([]byte(level)); err != nil {
			l = LevelInfo
		}

		o.Level = l
	}
}

// WithAddSource logger option sets the add source option, which will add source file and line number to the log record.
func WithAddSource(addSource bool) Option {
	return func(o *config) {
		o.AddSource = addSource
	}
}

// WithIsJSON logger option sets the is json option, which will set JSON format for the log record.
func WithIsJSON(isJSON bool) Option {
	return func(o *config) {
		o.IsJSON = isJSON
	}
}

// WithSetDefault logger option sets the set default option, which will set the created logger as default logger.
func WithSetDefault(setDefault bool) Option {
	return func(o *config) {
		o.SetDefault = setDefault
	}
}



// WithAttrs returns logger with attributes.
func WithAttrs(ctx context.Context, attrs ...Attr) *Logger {
	logger := L(ctx)
	for _, attr := range attrs {
		logger = logger.With(attr)
	}

	return logger
}

// WithDefaultAttrs returns logger with default attributes.
func WithDefaultAttrs(logger *Logger, attrs ...Attr) *Logger {
	for _, attr := range attrs {
		logger = logger.With(attr)
	}

	return logger
}
