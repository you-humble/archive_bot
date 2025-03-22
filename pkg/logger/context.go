package logger

import "context"

type loggerCtx struct{}

// ContextWithLogger adds logger to context.
func ContextWithLogger(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerCtx{}, logger)
}

// fromContext returns logger from context.
func fromContext(ctx context.Context) *Logger {
	if l, ok := ctx.Value(loggerCtx{}).(*Logger); ok {
		return l
	}

	return Default()
}
